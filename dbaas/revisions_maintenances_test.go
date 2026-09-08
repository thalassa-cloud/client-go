package dbaas

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

func TestListDbClusterRevisions(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	base, err := client.NewClient(client.WithBaseURL(server.URL))
	require.NoError(t, err)
	dbaasClient, err := New(base)
	require.NoError(t, err)

	revisions, err := dbaasClient.ListDbClusterRevisions(context.Background(), "cluster-123")
	require.NoError(t, err)
	require.Len(t, revisions, 1)
	assert.Equal(t, "rev-1", revisions[0].Identity)

	_, err = dbaasClient.ListDbClusterRevisions(context.Background(), "")
	assert.EqualError(t, err, "database cluster identity is required")
}

func TestScheduledMaintenances(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	base, err := client.NewClient(client.WithBaseURL(server.URL))
	require.NoError(t, err)
	dbaasClient, err := New(base)
	require.NoError(t, err)

	maintenances, err := dbaasClient.ListScheduledMaintenances(context.Background(), "cluster-123")
	require.NoError(t, err)
	require.Len(t, maintenances, 1)
	assert.Equal(t, "maint-1", maintenances[0].Identity)

	started, err := dbaasClient.StartScheduledMaintenance(context.Background(), "cluster-123", "maint-1")
	require.NoError(t, err)
	assert.Equal(t, DbClusterScheduledMaintenanceStatusInProgress, started.Status)

	postponed, err := dbaasClient.PostponeScheduledMaintenance(context.Background(), "cluster-123", "maint-1")
	require.NoError(t, err)
	assert.Equal(t, "maint-1", postponed.Identity)

	_, err = dbaasClient.StartScheduledMaintenance(context.Background(), "cluster-123", "")
	assert.EqualError(t, err, "maintenance identity is required")
}

func TestUpdateDbBackup(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	base, err := client.NewClient(client.WithBaseURL(server.URL))
	require.NoError(t, err)
	dbaasClient, err := New(base)
	require.NoError(t, err)

	protection := true
	backup, err := dbaasClient.UpdateDbBackup(context.Background(), "backup-123", UpdateDbClusterBackupRequest{DeleteProtection: &protection})
	require.NoError(t, err)
	assert.True(t, backup.DeleteProtection)

	_, err = dbaasClient.UpdateDbBackup(context.Background(), "backup-123", UpdateDbClusterBackupRequest{})
	assert.EqualError(t, err, "deleteProtection is required")
}
