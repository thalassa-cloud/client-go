package dbaas

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDbObjectStoreRecoveryWindowPitrAvailable(t *testing.T) {
	first := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	last := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name                   string
		window                 DbObjectStoreRecoveryWindow
		objectStoreReady       bool
		backupObjectStoreReady bool
		want                   bool
	}{
		{
			name: "ready with both points",
			window: DbObjectStoreRecoveryWindow{
				FirstRecoverabilityPoint: &first,
				LastSuccessfulBackupTime: &last,
			},
			objectStoreReady:       true,
			backupObjectStoreReady: true,
			want:                   true,
		},
		{
			name: "missing last successful backup",
			window: DbObjectStoreRecoveryWindow{
				FirstRecoverabilityPoint: &first,
			},
			objectStoreReady:       true,
			backupObjectStoreReady: true,
			want:                   false,
		},
		{
			name: "object store not ready",
			window: DbObjectStoreRecoveryWindow{
				FirstRecoverabilityPoint: &first,
				LastSuccessfulBackupTime: &last,
			},
			objectStoreReady:       false,
			backupObjectStoreReady: true,
			want:                   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.window.PitrAvailable(tt.objectStoreReady, tt.backupObjectStoreReady))
		})
	}
}

func TestDbObjectStoreServerRecoveryWindowJSON(t *testing.T) {
	first := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	in := `{
		"serverRecoveryWindow": {
			"cluster-1": {
				"firstRecoverabilityPoint": "2026-01-01T00:00:00Z"
			}
		}
	}`

	var store DbObjectStore
	require.NoError(t, json.Unmarshal([]byte(in), &store))
	require.Contains(t, store.ServerRecoveryWindow, "cluster-1")
	assert.Equal(t, &first, store.ServerRecoveryWindow["cluster-1"].FirstRecoverabilityPoint)
}
