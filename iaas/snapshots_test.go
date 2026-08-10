package iaas

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSnapshotPolicyRequest_TtlJSON(t *testing.T) {
	tests := []struct {
		name string
		ttl  time.Duration
		want string
	}{
		{name: "one week", ttl: 168 * time.Hour, want: "168h0m0s"},
		{name: "one day", ttl: 24 * time.Hour, want: "24h0m0s"},
		{name: "sub hour", ttl: 90 * time.Minute, want: "1h30m0s"},
		{name: "zero", ttl: 0, want: "0s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			create, err := json.Marshal(CreateSnapshotPolicyRequest{Ttl: tt.ttl})
			require.NoError(t, err)
			var gotCreate map[string]any
			require.NoError(t, json.Unmarshal(create, &gotCreate))
			assert.Equal(t, tt.want, gotCreate["ttl"])

			update, err := json.Marshal(UpdateSnapshotPolicyRequest{Ttl: tt.ttl})
			require.NoError(t, err)
			var gotUpdate map[string]any
			require.NoError(t, json.Unmarshal(update, &gotUpdate))
			assert.Equal(t, tt.want, gotUpdate["ttl"])
		})
	}
}

func TestCreateSnapshotPolicyRequest_JSONKeepsOtherFields(t *testing.T) {
	keepCount := 10
	req := CreateSnapshotPolicyRequest{
		Name:        "daily",
		Description: "Daily snapshots",
		Region:      "nl-01",
		Ttl:         168 * time.Hour,
		KeepCount:   &keepCount,
		Enabled:     true,
		Schedule:    "0 3 * * *",
		Timezone:    "Europe/Amsterdam",
		Target: SnapshotPolicyTarget{
			Type:             SnapshotPolicyTargetTypeExplicit,
			VolumeIdentities: []string{"v-abc"},
		},
	}

	b, err := json.Marshal(req)
	require.NoError(t, err)
	var got map[string]any
	require.NoError(t, json.Unmarshal(b, &got))

	assert.Equal(t, "daily", got["name"])
	assert.Equal(t, "Daily snapshots", got["description"])
	assert.Equal(t, "nl-01", got["region"])
	assert.Equal(t, "168h0m0s", got["ttl"])
	assert.Equal(t, float64(10), got["keepCount"])
	assert.Equal(t, true, got["enabled"])
	assert.Equal(t, "0 3 * * *", got["schedule"])
	assert.Equal(t, "Europe/Amsterdam", got["timezone"])
	assert.NotNil(t, got["target"])
}
