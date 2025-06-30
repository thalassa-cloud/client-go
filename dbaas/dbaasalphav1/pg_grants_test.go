package dbaasalphav1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

func setupGrantTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v1/dbaas/clusters/cluster-123/postgres-grants":
			if r.Method == "POST" {
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`{"message": "grant created"}`))
			}
		case "/v1/dbaas/clusters/cluster-123/postgres-grants/testgrant":
			if r.Method == "PUT" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"message": "grant updated"}`))
			} else if r.Method == "DELETE" {
				w.WriteHeader(http.StatusNoContent)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "not found"}`))
		}
	}))
}

func TestCreatePgGrant(t *testing.T) {
	server := setupGrantTestServer()
	defer server.Close()

	client, err := client.NewClient(client.WithBaseURL(server.URL))
	require.NoError(t, err)

	dbaasClient, err := New(client)
	require.NoError(t, err)

	tests := []struct {
		name              string
		dbClusterIdentity string
		request           CreatePgGrantRequest
		expectedError     string
	}{
		{
			name:              "successful grant creation",
			dbClusterIdentity: "cluster-123",
			request: CreatePgGrantRequest{
				Name:         "testgrant",
				RoleName:     "testrole",
				DatabaseName: "testdb",
				Read:         true,
				Write:        false,
			},
		},
		{
			name:              "missing cluster identity",
			dbClusterIdentity: "",
			request: CreatePgGrantRequest{
				Name:         "testgrant",
				RoleName:     "testrole",
				DatabaseName: "testdb",
			},
			expectedError: "database cluster identity is required",
		},
		{
			name:              "missing grant name",
			dbClusterIdentity: "cluster-123",
			request: CreatePgGrantRequest{
				Name:         "",
				RoleName:     "testrole",
				DatabaseName: "testdb",
			},
			expectedError: "grant name is required",
		},
		{
			name:              "missing role name",
			dbClusterIdentity: "cluster-123",
			request: CreatePgGrantRequest{
				Name:         "testgrant",
				RoleName:     "",
				DatabaseName: "testdb",
			},
			expectedError: "role name is required",
		},
		{
			name:              "missing database name",
			dbClusterIdentity: "cluster-123",
			request: CreatePgGrantRequest{
				Name:         "testgrant",
				RoleName:     "testrole",
				DatabaseName: "",
			},
			expectedError: "database name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dbaasClient.CreatePgGrant(context.Background(), tt.dbClusterIdentity, tt.request)

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdatePgGrant(t *testing.T) {
	server := setupGrantTestServer()
	defer server.Close()

	client, err := client.NewClient(client.WithBaseURL(server.URL))
	require.NoError(t, err)

	dbaasClient, err := New(client)
	require.NoError(t, err)

	tests := []struct {
		name              string
		dbClusterIdentity string
		grantName         string
		request           UpdatePgGrantRequest
		expectedError     string
	}{
		{
			name:              "successful grant update",
			dbClusterIdentity: "cluster-123",
			grantName:         "testgrant",
			request: UpdatePgGrantRequest{
				Read:  boolPtr(true),
				Write: boolPtr(false),
			},
		},
		{
			name:              "missing cluster identity",
			dbClusterIdentity: "",
			grantName:         "testgrant",
			request:           UpdatePgGrantRequest{},
			expectedError:     "database cluster identity is required",
		},
		{
			name:              "missing grant name",
			dbClusterIdentity: "cluster-123",
			grantName:         "",
			request:           UpdatePgGrantRequest{},
			expectedError:     "grant name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dbaasClient.UpdatePgGrant(context.Background(), tt.dbClusterIdentity, tt.grantName, tt.request)

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeletePgGrant(t *testing.T) {
	server := setupGrantTestServer()
	defer server.Close()

	client, err := client.NewClient(client.WithBaseURL(server.URL))
	require.NoError(t, err)

	dbaasClient, err := New(client)
	require.NoError(t, err)

	tests := []struct {
		name              string
		dbClusterIdentity string
		grantName         string
		expectedError     string
	}{
		{
			name:              "successful grant deletion",
			dbClusterIdentity: "cluster-123",
			grantName:         "testgrant",
		},
		{
			name:              "missing cluster identity",
			dbClusterIdentity: "",
			grantName:         "testgrant",
			expectedError:     "database cluster identity is required",
		},
		{
			name:              "missing grant name",
			dbClusterIdentity: "cluster-123",
			grantName:         "",
			expectedError:     "grant name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dbaasClient.DeletePgGrant(context.Background(), tt.dbClusterIdentity, tt.grantName)

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
