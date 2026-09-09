package dbaas

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDbClusterHealthStatusJSON(t *testing.T) {
	transitionedAt := time.Date(2026, 3, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		in   string
		want DbClusterHealthStatus
	}{
		{
			name: "continuous archiving healthy",
			in: `{
				"continuousArchiving": {
					"healthy": true,
					"reason": "ArchivingSucceeding",
					"message": "WAL is being archived",
					"lastTransitionTime": "2026-03-15T12:00:00Z"
				}
			}`,
			want: DbClusterHealthStatus{
				ContinuousArchiving: &DbClusterHealthCondition{
					Healthy:            true,
					Reason:             "ArchivingSucceeding",
					Message:            "WAL is being archived",
					LastTransitionTime: &transitionedAt,
				},
			},
		},
		{
			name: "empty health status",
			in:   `{}`,
			want: DbClusterHealthStatus{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got DbClusterHealthStatus
			require.NoError(t, json.Unmarshal([]byte(tt.in), &got))
			assert.Equal(t, tt.want, got)
		})
	}
}
