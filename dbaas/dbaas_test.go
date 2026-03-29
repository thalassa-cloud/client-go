package dbaas

import (
	"net/http"
	"net/http/httptest"
)

func setupTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v1/dbaas/clusters/cluster-123/postgres-databases":
			if r.Method == "POST" {
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`{"message": "database created"}`))
			}
		case "/v1/dbaas/clusters/cluster-123/postgres-databases/testdb":
			if r.Method == "PUT" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"message": "database updated"}`))
			} else if r.Method == "DELETE" {
				w.WriteHeader(http.StatusNoContent)
			}
		case "/v1/dbaas/clusters/cluster-123/postgres-roles":
			if r.Method == "POST" {
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`{"message": "role created"}`))
			}
		case "/v1/dbaas/clusters/cluster-123/postgres-roles/testrole":
			if r.Method == "PUT" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"message": "role updated"}`))
			} else if r.Method == "DELETE" {
				w.WriteHeader(http.StatusNoContent)
			}
		case "/v1/dbaas/clusters/cluster-123/backup-schedules":
			if r.Method == "POST" {
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`{
					"identity": "schedule-123",
					"name": "daily-backup",
					"status": "ready",
					"statusMessage": "",
					"createdAt": "2023-01-01T00:00:00Z",
					"method": "snapshot",
					"schedule": "0 2 * * *",
					"retentionPolicy": "30d",
					"backupCount": 0,
					"suspended": false,
					"target": "primary"
				}`))
			}
		case "/v1/dbaas/clusters/cluster-123/backup-schedules/schedule-123":
			if r.Method == "PUT" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"identity": "schedule-123",
					"name": "updated-backup",
					"status": "ready",
					"statusMessage": "",
					"createdAt": "2023-01-01T00:00:00Z",
					"method": "snapshot",
					"schedule": "0 3 * * *",
					"retentionPolicy": "60d",
					"backupCount": 0,
					"suspended": false,
					"target": "primary"
				}`))
			} else if r.Method == "DELETE" {
				w.WriteHeader(http.StatusNoContent)
			}
		case "/v1/dbaas/object-stores":
			if r.Method == "GET" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[]`))
			} else if r.Method == "POST" {
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`{
					"identity": "os-1",
					"name": "backup-store",
					"description": "",
					"createdAt": "2023-01-01T00:00:00Z",
					"objectVersion": 1,
					"status": "creating",
					"deleteProtection": true,
					"retentionPolicy": "30d"
				}`))
			}
		case "/v1/dbaas/object-stores/os-1":
			if r.Method == "GET" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"identity": "os-1",
					"name": "backup-store",
					"description": "",
					"createdAt": "2023-01-01T00:00:00Z",
					"objectVersion": 1,
					"status": "ready",
					"deleteProtection": true,
					"retentionPolicy": "30d"
				}`))
			} else if r.Method == "PUT" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"identity": "os-1",
					"name": "backup-store",
					"description": "updated",
					"createdAt": "2023-01-01T00:00:00Z",
					"objectVersion": 2,
					"status": "ready",
					"deleteProtection": false,
					"retentionPolicy": "60d"
				}`))
			} else if r.Method == "DELETE" {
				w.WriteHeader(http.StatusNoContent)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "not found"}`))
		}
	}))
}

// Helper functions
func boolPtr(b bool) *bool {
	return &b
}
