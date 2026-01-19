package dbaas

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

func TestCreatePgBackupSchedule(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client, err := client.NewClient(client.WithBaseURL(server.URL))
	require.NoError(t, err)

	dbaasClient, err := New(client)
	require.NoError(t, err)

	tests := []struct {
		name              string
		dbClusterIdentity string
		request           CreatePgBackupScheduleRequest
		expectedError     string
		expectedResult    *DbClusterBackupSchedule
	}{
		{
			name:              "successful backup schedule creation",
			dbClusterIdentity: "cluster-123",
			request: CreatePgBackupScheduleRequest{
				Name:            "daily-backup",
				Schedule:        "0 2 * * *",
				RetentionPolicy: "30d",
			},
			expectedResult: &DbClusterBackupSchedule{
				Identity: "schedule-123",
				Name:     "daily-backup",
				Status:   "ready",
			},
		},
		{
			name:              "missing cluster identity",
			dbClusterIdentity: "",
			request: CreatePgBackupScheduleRequest{
				Name:            "daily-backup",
				Schedule:        "0 2 * * *",
				RetentionPolicy: "30d",
			},
			expectedError: "database cluster identity is required",
		},
		{
			name:              "missing schedule name",
			dbClusterIdentity: "cluster-123",
			request: CreatePgBackupScheduleRequest{
				Name:            "",
				Schedule:        "0 2 * * *",
				RetentionPolicy: "30d",
			},
			expectedError: "backup schedule name is required",
		},
		{
			name:              "missing schedule",
			dbClusterIdentity: "cluster-123",
			request: CreatePgBackupScheduleRequest{
				Name:            "daily-backup",
				Schedule:        "",
				RetentionPolicy: "30d",
			},
			expectedError: "backup schedule is required",
		},
		{
			name:              "missing retention policy",
			dbClusterIdentity: "cluster-123",
			request: CreatePgBackupScheduleRequest{
				Name:            "daily-backup",
				Schedule:        "0 2 * * *",
				RetentionPolicy: "",
			},
			expectedError: "retention policy is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := dbaasClient.CreatePgBackupSchedule(context.Background(), tt.dbClusterIdentity, tt.request)

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if result != nil {
					assert.Equal(t, tt.expectedResult.Identity, result.Identity, "Identity")
					assert.Equal(t, tt.expectedResult.Name, result.Name, "Name")
					assert.Equal(t, tt.expectedResult.Status, result.Status, "Status")
				}
			}
		})
	}
}

func TestUpdatePgBackupSchedule(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client, err := client.NewClient(client.WithBaseURL(server.URL))
	require.NoError(t, err)

	dbaasClient, err := New(client)
	require.NoError(t, err)

	tests := []struct {
		name                   string
		dbClusterIdentity      string
		backupScheduleIdentity string
		request                UpdatePgBackupScheduleRequest
		expectedError          string
		expectedResult         *DbClusterBackupSchedule
	}{
		{
			name:                   "successful backup schedule update",
			dbClusterIdentity:      "cluster-123",
			backupScheduleIdentity: "schedule-123",
			request: UpdatePgBackupScheduleRequest{
				Name:            "updated-backup",
				Schedule:        "0 3 * * *",
				RetentionPolicy: "60d",
			},
			expectedResult: &DbClusterBackupSchedule{
				Identity: "schedule-123",
				Name:     "updated-backup",
				Status:   "ready",
			},
		},
		{
			name:                   "missing cluster identity",
			dbClusterIdentity:      "",
			backupScheduleIdentity: "schedule-123",
			request: UpdatePgBackupScheduleRequest{
				Name:            "updated-backup",
				Schedule:        "0 3 * * *",
				RetentionPolicy: "60d",
			},
			expectedError: "database cluster identity is required",
		},
		{
			name:                   "missing backup schedule identity",
			dbClusterIdentity:      "cluster-123",
			backupScheduleIdentity: "",
			request: UpdatePgBackupScheduleRequest{
				Name:            "updated-backup",
				Schedule:        "0 3 * * *",
				RetentionPolicy: "60d",
			},
			expectedError: "backup schedule identity is required",
		},
		{
			name:                   "missing schedule name",
			dbClusterIdentity:      "cluster-123",
			backupScheduleIdentity: "schedule-123",
			request: UpdatePgBackupScheduleRequest{
				Name:            "",
				Schedule:        "0 3 * * *",
				RetentionPolicy: "60d",
			},
			expectedError: "backup schedule name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := dbaasClient.UpdatePgBackupSchedule(context.Background(), tt.dbClusterIdentity, tt.backupScheduleIdentity, tt.request)

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if result != nil {
					assert.Equal(t, tt.expectedResult.Identity, result.Identity)
					assert.Equal(t, tt.expectedResult.Name, result.Name)
					assert.Equal(t, tt.expectedResult.Status, result.Status)
				}
			}
		})
	}
}

func TestDeletePgBackupSchedule(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client, err := client.NewClient(client.WithBaseURL(server.URL))
	require.NoError(t, err)

	dbaasClient, err := New(client)
	require.NoError(t, err)

	tests := []struct {
		name                   string
		dbClusterIdentity      string
		backupScheduleIdentity string
		expectedError          string
	}{
		{
			name:                   "successful backup schedule deletion",
			dbClusterIdentity:      "cluster-123",
			backupScheduleIdentity: "schedule-123",
		},
		{
			name:                   "missing cluster identity",
			dbClusterIdentity:      "",
			backupScheduleIdentity: "schedule-123",
			expectedError:          "database cluster identity is required",
		},
		{
			name:                   "missing backup schedule identity",
			dbClusterIdentity:      "cluster-123",
			backupScheduleIdentity: "",
			expectedError:          "backup schedule identity is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dbaasClient.DeletePgBackupSchedule(context.Background(), tt.dbClusterIdentity, tt.backupScheduleIdentity)

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
