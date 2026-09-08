package observability

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

func setupTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/v1/observability/workspaces" && r.Method == http.MethodGet:
			_ = json.NewEncoder(w).Encode([]ObservabilityWorkspace{{
				Identity: "obsw-1",
				Name:     "ws-1",
				Status:   ObservabilityWorkspaceStatusReady,
			}})
		case r.URL.Path == "/v1/observability/workspaces" && r.Method == http.MethodPost:
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(ObservabilityWorkspace{
				Identity:          "obsw-new",
				Name:              "created",
				Status:            ObservabilityWorkspaceStatusProvisioning,
				PrometheusEnabled: true,
				LokiEnabled:       true,
			})
		case r.URL.Path == "/v1/observability/workspaces/obsw-1" && r.Method == http.MethodGet:
			_ = json.NewEncoder(w).Encode(ObservabilityWorkspace{
				Identity: "obsw-1",
				Name:     "ws-1",
				Status:   ObservabilityWorkspaceStatusReady,
			})
		case r.URL.Path == "/v1/observability/workspaces/obsw-1" && r.Method == http.MethodPut:
			_ = json.NewEncoder(w).Encode(ObservabilityWorkspace{
				Identity: "obsw-1",
				Name:     "updated",
				Status:   ObservabilityWorkspaceStatusReady,
			})
		case r.URL.Path == "/v1/observability/workspaces/obsw-1" && r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func newTestClient(t *testing.T, serverURL string) *Client {
	t.Helper()
	base, err := client.NewClient(client.WithBaseURL(serverURL), client.WithAuthCustom())
	require.NoError(t, err)
	obs, err := New(base)
	require.NoError(t, err)
	return obs
}

func TestListObservabilityWorkspaces(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	obs := newTestClient(t, server.URL)
	workspaces, err := obs.ListObservabilityWorkspaces(context.Background(), nil)
	require.NoError(t, err)
	require.Len(t, workspaces, 1)
	assert.Equal(t, "obsw-1", workspaces[0].Identity)
}

func TestGetObservabilityWorkspace(t *testing.T) {
	server := setupTestServer()
	defer server.Close()
	obs := newTestClient(t, server.URL)

	tests := []struct {
		name          string
		identity      string
		expectedError string
	}{
		{name: "success", identity: "obsw-1"},
		{name: "missing identity", identity: "", expectedError: "workspace identity is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ws, err := obs.GetObservabilityWorkspace(context.Background(), tt.identity)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, "obsw-1", ws.Identity)
		})
	}
}

func TestCreateObservabilityWorkspace(t *testing.T) {
	server := setupTestServer()
	defer server.Close()
	obs := newTestClient(t, server.URL)

	retentionTooHigh := 2000
	tests := []struct {
		name          string
		request       CreateObservabilityWorkspaceRequest
		expectedError string
	}{
		{
			name:    "success",
			request: CreateObservabilityWorkspaceRequest{Name: "created", RegionIdentity: "nl-01"},
		},
		{
			name:          "missing name",
			request:       CreateObservabilityWorkspaceRequest{},
			expectedError: "name is required",
		},
		{
			name:          "invalid retention",
			request:       CreateObservabilityWorkspaceRequest{Name: "x", RetentionDays: &retentionTooHigh},
			expectedError: "retentionDays cannot be greater than 1095 (3 years), got 2000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ws, err := obs.CreateObservabilityWorkspace(context.Background(), tt.request)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, "obsw-new", ws.Identity)
		})
	}
}

func TestUpdateObservabilityWorkspace(t *testing.T) {
	server := setupTestServer()
	defer server.Close()
	obs := newTestClient(t, server.URL)

	ws, err := obs.UpdateObservabilityWorkspace(context.Background(), "obsw-1", UpdateObservabilityWorkspaceRequest{Name: "updated"})
	require.NoError(t, err)
	assert.Equal(t, "updated", ws.Name)

	_, err = obs.UpdateObservabilityWorkspace(context.Background(), "", UpdateObservabilityWorkspaceRequest{Name: "x"})
	assert.EqualError(t, err, "workspace identity is required")

	_, err = obs.UpdateObservabilityWorkspace(context.Background(), "obsw-1", UpdateObservabilityWorkspaceRequest{})
	assert.EqualError(t, err, "name is required")
}

func TestDeleteObservabilityWorkspace(t *testing.T) {
	server := setupTestServer()
	defer server.Close()
	obs := newTestClient(t, server.URL)

	require.NoError(t, obs.DeleteObservabilityWorkspace(context.Background(), "obsw-1"))
	assert.EqualError(t, obs.DeleteObservabilityWorkspace(context.Background(), ""), "workspace identity is required")
}
