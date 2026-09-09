package dbaas

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDbClusterStorageAutoScaleJSON(t *testing.T) {
	in := `{
		"storageAutoScaleEnabled": true,
		"storageAutoScaleMaxGB": 500,
		"storageAutoScaleThresholdPercent": 85,
		"storageAutoScaleIncreasePercent": 15,
		"storageAutoScaleMaxIncreaseGB": 50
	}`

	var cluster DbCluster
	require.NoError(t, json.Unmarshal([]byte(in), &cluster))
	assert.True(t, cluster.StorageAutoScaleEnabled)
	assert.Equal(t, uint64(500), cluster.StorageAutoScaleMaxGB)
	assert.Equal(t, uint(85), cluster.StorageAutoScaleThresholdPercent)
	assert.Equal(t, uint(15), cluster.StorageAutoScaleIncreasePercent)
	assert.Equal(t, uint64(50), cluster.StorageAutoScaleMaxIncreaseGB)

	create := CreateDbClusterRequest{
		StorageAutoScaleEnabled:          ptrBool(true),
		StorageAutoScaleMaxGB:            ptrUint64(500),
		StorageAutoScaleThresholdPercent: ptrUint(85),
		StorageAutoScaleIncreasePercent:  ptrUint(15),
		StorageAutoScaleMaxIncreaseGB:    ptrUint64(50),
	}
	raw, err := json.Marshal(create)
	require.NoError(t, err)
	assert.Contains(t, string(raw), `"storageAutoScaleEnabled":true`)
	assert.Contains(t, string(raw), `"storageAutoScaleMaxGB":500`)
}

func ptrBool(v bool) *bool       { return &v }
func ptrUint(v uint) *uint       { return &v }
func ptrUint64(v uint64) *uint64 { return &v }
