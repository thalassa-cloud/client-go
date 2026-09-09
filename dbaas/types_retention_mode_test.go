package dbaas

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDbObjectStoreRetentionMode(t *testing.T) {
	tests := []struct {
		name                      string
		mode                      DbObjectStoreRetentionMode
		wantValid                 bool
		wantRetainsForPointInTime bool
	}{
		{
			name:                      "retain for point in time",
			mode:                      DbObjectStoreRetentionModeRetainForPointInTime,
			wantValid:                 true,
			wantRetainsForPointInTime: true,
		},
		{
			name:                      "force cleanup after expiry",
			mode:                      DbObjectStoreRetentionModeForceCleanupAfterExpiry,
			wantValid:                 true,
			wantRetainsForPointInTime: false,
		},
		{
			name:                      "empty defaults to retain for point in time",
			mode:                      "",
			wantValid:                 false,
			wantRetainsForPointInTime: true,
		},
		{
			name:                      "unknown mode",
			mode:                      "unknown",
			wantValid:                 false,
			wantRetainsForPointInTime: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantValid, tt.mode.IsValid())
			assert.Equal(t, tt.wantRetainsForPointInTime, tt.mode.RetainsForPointInTime())
		})
	}
}
