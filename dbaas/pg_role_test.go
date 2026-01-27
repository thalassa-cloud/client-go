package dbaas

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

func TestCreatePgRole(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client, err := client.NewClient(client.WithBaseURL(server.URL))
	require.NoError(t, err)

	dbaasClient, err := New(client)
	require.NoError(t, err)

	tests := []struct {
		name              string
		dbClusterIdentity string
		request           CreatePgRoleRequest
		expectedError     string
	}{
		{
			name:              "successful role creation",
			dbClusterIdentity: "cluster-123",
			request: CreatePgRoleRequest{
				Name:     "testrole",
				Login:    true,
				CreateDb: false,
			},
		},
		{
			name:              "missing cluster identity",
			dbClusterIdentity: "",
			request: CreatePgRoleRequest{
				Name: "testrole",
			},
			expectedError: "database cluster identity is required",
		},
		{
			name:              "missing role name",
			dbClusterIdentity: "cluster-123",
			request: CreatePgRoleRequest{
				Name: "",
			},
			expectedError: "role name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := dbaasClient.CreatePgRole(context.Background(), tt.dbClusterIdentity, tt.request)

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdatePgRole(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client, err := client.NewClient(client.WithBaseURL(server.URL))
	require.NoError(t, err)

	dbaasClient, err := New(client)
	require.NoError(t, err)

	tests := []struct {
		name              string
		dbClusterIdentity string
		roleName          string
		request           UpdatePgRoleRequest
		expectedError     string
	}{
		{
			name:              "successful role update",
			dbClusterIdentity: "cluster-123",
			roleName:          "testrole",
			request: UpdatePgRoleRequest{
				ConnectionLimit: 100,
			},
		},
		{
			name:              "missing cluster identity",
			dbClusterIdentity: "",
			roleName:          "testrole",
			request:           UpdatePgRoleRequest{},
			expectedError:     "database cluster identity is required",
		},
		{
			name:              "missing role name",
			dbClusterIdentity: "cluster-123",
			roleName:          "",
			request:           UpdatePgRoleRequest{},
			expectedError:     "postgres role identity is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := dbaasClient.UpdatePgRole(context.Background(), tt.dbClusterIdentity, tt.roleName, tt.request)

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeletePgRole(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client, err := client.NewClient(client.WithBaseURL(server.URL))
	require.NoError(t, err)

	dbaasClient, err := New(client)
	require.NoError(t, err)

	tests := []struct {
		name              string
		dbClusterIdentity string
		roleIdentity      string
		expectedError     string
	}{
		{
			name:              "successful role deletion",
			dbClusterIdentity: "cluster-123",
			roleIdentity:      "testrole",
		},
		{
			name:              "missing cluster identity",
			dbClusterIdentity: "",
			roleIdentity:      "testrole",
			expectedError:     "database cluster identity is required",
		},
		{
			name:              "missing role identity",
			dbClusterIdentity: "cluster-123",
			roleIdentity:      "",
			expectedError:     "postgres role identity is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dbaasClient.DeletePgRole(context.Background(), tt.dbClusterIdentity, tt.roleIdentity)

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
