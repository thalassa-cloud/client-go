package dbaas

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

func TestListDbObjectStores(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL))
	require.NoError(t, err)
	dbaasClient, err := New(c)
	require.NoError(t, err)

	stores, err := dbaasClient.ListDbObjectStores(context.Background(), nil)
	require.NoError(t, err)
	assert.Empty(t, stores)
}

func TestCreateDbObjectStore(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL))
	require.NoError(t, err)
	dbaasClient, err := New(c)
	require.NoError(t, err)

	tests := []struct {
		name          string
		req           CreateDbObjectStoreRequest
		expectedError string
	}{
		{
			name: "success",
			req: CreateDbObjectStoreRequest{
				Name:        "backup-store",
				Description: "",
				Region:      "nl-01",
			},
		},
		{
			name: "missing name",
			req: CreateDbObjectStoreRequest{
				Region: "nl-01",
			},
			expectedError: "name is required",
		},
		{
			name: "missing region",
			req: CreateDbObjectStoreRequest{
				Name: "backup-store",
			},
			expectedError: "region is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dbaasClient.CreateDbObjectStore(context.Background(), tt.req)
			if tt.expectedError != "" {
				assert.Nil(t, got)
				assert.EqualError(t, err, tt.expectedError)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, "os-1", got.Identity)
			assert.Equal(t, "backup-store", got.Name)
		})
	}
}

func TestGetDbObjectStore(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL))
	require.NoError(t, err)
	dbaasClient, err := New(c)
	require.NoError(t, err)

	_, err = dbaasClient.GetDbObjectStore(context.Background(), "")
	assert.EqualError(t, err, "identity is required")

	got, err := dbaasClient.GetDbObjectStore(context.Background(), "os-1")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, ObjectStatusReady, got.Status)
}

func TestUpdateDbObjectStore(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL))
	require.NoError(t, err)
	dbaasClient, err := New(c)
	require.NoError(t, err)

	_, err = dbaasClient.UpdateDbObjectStore(context.Background(), "", UpdateDbObjectStoreRequest{Name: "x"})
	assert.EqualError(t, err, "identity is required")

	_, err = dbaasClient.UpdateDbObjectStore(context.Background(), "os-1", UpdateDbObjectStoreRequest{})
	assert.EqualError(t, err, "name is required")

	got, err := dbaasClient.UpdateDbObjectStore(context.Background(), "os-1", UpdateDbObjectStoreRequest{
		Name:             "backup-store",
		Description:      "updated",
		RetentionPolicy:  "60d",
		DeleteProtection: false,
	})
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "updated", got.Description)
	assert.False(t, got.DeleteProtection)
}

func TestDeleteDbObjectStore(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL))
	require.NoError(t, err)
	dbaasClient, err := New(c)
	require.NoError(t, err)

	err = dbaasClient.DeleteDbObjectStore(context.Background(), "")
	assert.EqualError(t, err, "identity is required")

	err = dbaasClient.DeleteDbObjectStore(context.Background(), "os-1")
	assert.NoError(t, err)
}
