package dbaas

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

func TestCreatePgDatabase(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client, err := client.NewClient(client.WithBaseURL(server.URL))
	require.NoError(t, err)

	dbaasClient, err := New(client)
	require.NoError(t, err)

	tests := []struct {
		name              string
		dbClusterIdentity string
		request           CreatePgDatabaseRequest
		expectedError     string
	}{
		{
			name:              "successful database creation",
			dbClusterIdentity: "cluster-123",
			request: CreatePgDatabaseRequest{
				Name:  "testdb",
				Owner: "testuser",
			},
		},
		{
			name:              "missing cluster identity",
			dbClusterIdentity: "",
			request: CreatePgDatabaseRequest{
				Name:  "testdb",
				Owner: "testuser",
			},
			expectedError: "database cluster identity is required",
		},
		{
			name:              "missing database name",
			dbClusterIdentity: "cluster-123",
			request: CreatePgDatabaseRequest{
				Name:  "",
				Owner: "testuser",
			},
			expectedError: "database name is required",
		},
		{
			name:              "missing database owner",
			dbClusterIdentity: "cluster-123",
			request: CreatePgDatabaseRequest{
				Name:  "testdb",
				Owner: "",
			},
			expectedError: "database owner is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dbaasClient.CreatePgDatabase(context.Background(), tt.dbClusterIdentity, tt.request)

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdatePgDatabase(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client, err := client.NewClient(client.WithBaseURL(server.URL))
	require.NoError(t, err)

	dbaasClient, err := New(client)
	require.NoError(t, err)

	tests := []struct {
		name              string
		dbClusterIdentity string
		databaseName      string
		request           UpdatePgDatabaseRequest
		expectedError     string
	}{
		{
			name:              "successful database update",
			dbClusterIdentity: "cluster-123",
			databaseName:      "testdb",
			request: UpdatePgDatabaseRequest{
				AllowConnections: boolPtr(true),
			},
		},
		{
			name:              "missing cluster identity",
			dbClusterIdentity: "",
			databaseName:      "testdb",
			request:           UpdatePgDatabaseRequest{},
			expectedError:     "database cluster identity is required",
		},
		{
			name:              "missing database identity",
			dbClusterIdentity: "cluster-123",
			databaseName:      "",
			request:           UpdatePgDatabaseRequest{},
			expectedError:     "postgres database identity is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dbaasClient.UpdatePgDatabase(context.Background(), tt.dbClusterIdentity, tt.databaseName, tt.request)

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeletePgDatabase(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client, err := client.NewClient(client.WithBaseURL(server.URL))
	require.NoError(t, err)

	dbaasClient, err := New(client)
	require.NoError(t, err)

	tests := []struct {
		name              string
		dbClusterIdentity string
		databaseName      string
		immediate         bool
		expectedError     string
	}{
		{
			name:              "successful database deletion",
			dbClusterIdentity: "cluster-123",
			databaseName:      "testdb",
			immediate:         false,
		},
		{
			name:              "missing cluster identity",
			dbClusterIdentity: "",
			databaseName:      "testdb",
			expectedError:     "database cluster identity is required",
			immediate:         false,
		},
		{
			name:              "missing database identity",
			dbClusterIdentity: "cluster-123",
			databaseName:      "",
			expectedError:     "postgres database identity is required",
			immediate:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dbaasClient.DeletePgDatabase(context.Background(), tt.dbClusterIdentity, tt.databaseName, tt.immediate)

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
