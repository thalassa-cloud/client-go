package dbaas

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDbObjectStoreRecoveryOverviewJSON(t *testing.T) {
	in := `{
		"objectStoreIdentity": "dbos-1",
		"complete": true,
		"truncation": {"truncated": false},
		"summary": {
			"completedBackupCount": 2,
			"gapCount": 1,
			"timelines": ["00000001"]
		},
		"clusters": [
			{
				"clusterIdentity": "cluster-1",
				"clusterName": "prod",
				"clusterExists": true,
				"gapCount": 0,
				"backups": [
					{
						"backupIdentity": "bak-1",
						"status": "ready",
						"coverage": {
							"chainStatus": "continuous",
							"pitrThroughSource": "wal_object_last_modified",
							"continuousSegmentCount": 10
						}
					}
				]
			}
		]
	}`

	var got DbObjectStoreRecoveryOverview
	require.NoError(t, json.Unmarshal([]byte(in), &got))
	assert.Equal(t, "dbos-1", got.ObjectStoreIdentity)
	assert.True(t, got.Complete)
	require.Len(t, got.Clusters, 1)
	assert.Equal(t, "cluster-1", got.Clusters[0].ClusterIdentity)
	require.Len(t, got.Clusters[0].Backups, 1)
	assert.Equal(t, DbObjectStoreWalChainStatusContinuous, got.Clusters[0].Backups[0].Coverage.ChainStatus)
	assert.Equal(t, DbObjectStoreWalPitrSourceWalObjectLastModified, got.Clusters[0].Backups[0].Coverage.PitrThroughSource)
}

func TestListDbObjectStoreRecoveryWalSegmentsRequestValidation(t *testing.T) {
	c := &Client{}
	_, err := c.ListDbObjectStoreRecoveryWalSegments(t.Context(), "dbos-1", "cluster-1", nil)
	assert.Error(t, err)

	_, err = c.ListDbObjectStoreRecoveryWalSegments(t.Context(), "dbos-1", "cluster-1", &ListDbObjectStoreRecoveryWalSegmentsRequest{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timeline and log")

	_, err = c.ListDbObjectStoreRecoveryWalSegments(t.Context(), "", "cluster-1", &ListDbObjectStoreRecoveryWalSegmentsRequest{View: "flat"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "identity is required")
}
