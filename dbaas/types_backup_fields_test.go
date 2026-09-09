package dbaas

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDbClusterBackupRetentionAndSizeJSON(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want DbClusterBackup
	}{
		{
			name: "retention expired with size",
			in: `{
				"retentionExpired": true,
				"sizeBytes": 1048576
			}`,
			want: DbClusterBackup{
				RetentionExpired: true,
				SizeBytes:        ptrInt64(1048576),
			},
		},
		{
			name: "defaults when omitted",
			in:   `{}`,
			want: DbClusterBackup{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got DbClusterBackup
			require.NoError(t, json.Unmarshal([]byte(tt.in), &got))
			assert.Equal(t, tt.want.RetentionExpired, got.RetentionExpired)
			assert.Equal(t, tt.want.SizeBytes, got.SizeBytes)
		})
	}
}

func ptrInt64(v int64) *int64 {
	return &v
}
