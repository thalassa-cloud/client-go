package quicklaunch

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

func TestQuickLaunchRequest_JSON(t *testing.T) {
	req := QuickLaunchRequest{
		Template:            QuickLaunchTemplateKubernetes,
		Name:                "demo",
		Description:         "test",
		CloudRegionIdentity: "nl-01",
		VpcCidr:             "10.0.0.0/16",
		SubnetCidrs:         []string{"10.0.1.0/24"},
		MachineType:         "pgp-small",
		Labels:              map[string]string{"env": "dev"},
		Annotations:         map[string]string{"note": "x"},
	}
	b, err := json.Marshal(req)
	require.NoError(t, err)
	var got map[string]any
	require.NoError(t, json.Unmarshal(b, &got))
	assert.Equal(t, "kubernetes", got["template"])
	assert.Equal(t, "demo", got["name"])
	assert.Equal(t, "nl-01", got["cloudRegionIdentity"])
	assert.Equal(t, "10.0.0.0/16", got["vpcCidr"])
}

func TestQuickLaunchResources_Add(t *testing.T) {
	var r QuickLaunchResources
	r.Add(QuickLaunchResource{Identity: "a", Type: "vpc", Name: "v"})
	r.Add(QuickLaunchResource{Identity: "a", Type: "vpc", Name: "dup"})
	r.Add(QuickLaunchResource{Identity: "b", Type: "subnet", Name: "s"})
	assert.Len(t, r, 2)
	assert.Equal(t, []QuickLaunchResource{
		{Identity: "a", Type: "vpc", Name: "v"},
		{Identity: "b", Type: "subnet", Name: "s"},
	}, r.ToSlice())
}

func TestQuickLaunch_unmarshal(t *testing.T) {
	raw := `{
		"identity": "ql-1",
		"name": "my-ql",
		"slug": "my-ql",
		"template": "vpc",
		"status": "running",
		"cloudRegionIdentity": "nl-01",
		"resources": [
			{"type": "vpc", "name": "v", "identity": "vpc-1", "lastStatus": "ready"}
		],
		"createdAt": "2025-01-01T12:00:00Z"
	}`
	var ql QuickLaunch
	require.NoError(t, json.Unmarshal([]byte(raw), &ql))
	assert.Equal(t, "ql-1", ql.Identity)
	assert.Equal(t, QuickLaunchTemplateVPC, ql.Template)
	require.Len(t, ql.Resources, 1)
	assert.Equal(t, "vpc", ql.Resources[0].Type)
	assert.Equal(t, "ready", ql.Resources[0].LastStatus)
}

func TestClient_QuickLaunchAPI(t *testing.T) {
	t.Parallel()

	const qlJSON = `{"identity":"ql-1","name":"n","slug":"n","template":"vpc","status":"pending","cloudRegionIdentity":"r1","createdAt":"2025-01-01T12:00:00Z"}`
	const listJSON = `[` + qlJSON + `]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/quick-launch":
			b, _ := io.ReadAll(r.Body)
			var body map[string]any
			require.NoError(t, json.Unmarshal(b, &body))
			assert.Equal(t, "demo", body["name"])
			assert.Equal(t, "r1", body["cloudRegionIdentity"])
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(qlJSON))
		case r.Method == http.MethodGet && r.URL.Path == "/v1/quick-launch":
			assert.Equal(t, "x", r.URL.Query().Get("name"))
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(listJSON))
		case r.Method == http.MethodGet && r.URL.Path == "/v1/quick-launch/ql-1":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(qlJSON))
		case r.Method == http.MethodGet && r.URL.Path == "/v1/quick-launch/ql-1/logs":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"lines":[]}`))
		case r.Method == http.MethodDelete && r.URL.Path == "/v1/quick-launch/ql-1":
			assert.Equal(t, string(QuickLaunchCascadeDelete), r.URL.Query().Get("cascade"))
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(srv.Close)

	baseClient, err := client.NewClient(client.WithBaseURL(srv.URL))
	require.NoError(t, err)
	qlc, err := New(baseClient)
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("CreateQuickLaunch validation", func(t *testing.T) {
		_, err := qlc.CreateQuickLaunch(ctx, QuickLaunchRequest{})
		require.Error(t, err)
		_, err = qlc.CreateQuickLaunch(ctx, QuickLaunchRequest{Name: "x"})
		require.Error(t, err)
	})

	created, err := qlc.CreateQuickLaunch(ctx, QuickLaunchRequest{Name: "demo", CloudRegionIdentity: "r1"})
	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Equal(t, "ql-1", created.Identity)

	list, err := qlc.ListQuickLaunches(ctx, &ListQuickLaunchesRequest{
		Filters: []filters.Filter{&filters.FilterKeyValue{Key: filters.FilterName, Value: "x"}},
	})
	require.NoError(t, err)
	require.Len(t, list, 1)

	got, err := qlc.GetQuickLaunch(ctx, "ql-1")
	require.NoError(t, err)
	assert.Equal(t, "pending", got.Status)

	logs, err := qlc.GetQuickLaunchLogs(ctx, "ql-1")
	require.NoError(t, err)
	assert.JSONEq(t, `{"lines":[]}`, string(logs))

	require.NoError(t, qlc.DeleteQuickLaunch(ctx, "ql-1", QuickLaunchCascadeDelete))
}

func TestClient_GetQuickLaunch_identityRequired(t *testing.T) {
	t.Parallel()
	baseClient, err := client.NewClient(client.WithBaseURL("http://example.com"))
	require.NoError(t, err)
	qlc, err := New(baseClient)
	require.NoError(t, err)
	_, err = qlc.GetQuickLaunch(context.Background(), "")
	require.Error(t, err)
}
