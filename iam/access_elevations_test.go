package iam

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

func TestListMyAccessElevations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v1/projects/iam/my-access-elevations", r.URL.Path)
		assertIAMScopedHeaders(t, r)

		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode([]IamAccessElevationRequest{
			{Identity: "ae-1", Status: IamAccessElevationStatusPending, Reason: "incident"},
		}))
	}))
	defer server.Close()

	requests, err := newTestIAMClient(t, server.URL).ListMyAccessElevations(context.Background())
	require.NoError(t, err)
	require.Len(t, requests, 1)
	assert.Equal(t, "ae-1", requests[0].Identity)
}

func TestCreateAccessElevation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v1/projects/iam/my-access-elevations", r.URL.Path)

		var body CreateAccessElevationRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "pol-1", body.PolicyIdentity)
		assert.Equal(t, "break glass", body.Reason)
		assert.Equal(t, "2h", body.Duration)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		require.NoError(t, json.NewEncoder(w).Encode(IamAccessElevationRequest{
			Identity:           "ae-1",
			Status:             IamAccessElevationStatusPending,
			Reason:             body.Reason,
			RequestedExpiresAt: time.Now().UTC().Add(2 * time.Hour),
		}))
	}))
	defer server.Close()

	request, err := newTestIAMClient(t, server.URL).CreateAccessElevation(context.Background(), CreateAccessElevationRequest{
		PolicyIdentity: "pol-1",
		Reason:         "break glass",
		Duration:       "2h",
	})
	require.NoError(t, err)
	assert.Equal(t, "ae-1", request.Identity)
	assert.Equal(t, IamAccessElevationStatusPending, request.Status)
}

func TestCreateAccessElevationValidation(t *testing.T) {
	c := newTestIAMClient(t, "http://example.com")

	_, err := c.CreateAccessElevation(context.Background(), CreateAccessElevationRequest{
		Reason:   "x",
		Duration: "1h",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "policyIdentity is required")

	_, err = c.CreateAccessElevation(context.Background(), CreateAccessElevationRequest{
		PolicyIdentity: "pol-1",
		Duration:       "1h",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "reason is required")

	_, err = c.CreateAccessElevation(context.Background(), CreateAccessElevationRequest{
		PolicyIdentity: "pol-1",
		Reason:         "x",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "duration or expiresAt is required")
}

func TestCancelAccessElevation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v1/projects/iam/my-access-elevations/ae-1/cancel", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(IamAccessElevationRequest{
			Identity: "ae-1",
			Status:   IamAccessElevationStatusCancelled,
		}))
	}))
	defer server.Close()

	request, err := newTestIAMClient(t, server.URL).CancelAccessElevation(context.Background(), "ae-1")
	require.NoError(t, err)
	assert.Equal(t, IamAccessElevationStatusCancelled, request.Status)
}

func TestGetMyAccessElevationNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		require.NoError(t, json.NewEncoder(w).Encode(map[string]string{"message": "not found"}))
	}))
	defer server.Close()

	_, err := newTestIAMClient(t, server.URL).GetMyAccessElevation(context.Background(), "missing")
	require.Error(t, err)
	assert.True(t, client.IsNotFound(err))
}

func TestListAccessElevations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/projects/iam/access-elevations", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode([]IamAccessElevationRequest{
			{Identity: "ae-2", Status: IamAccessElevationStatusPending},
		}))
	}))
	defer server.Close()

	requests, err := newTestIAMClient(t, server.URL).ListAccessElevations(context.Background())
	require.NoError(t, err)
	require.Len(t, requests, 1)
}

func TestApproveAccessElevation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v1/projects/iam/access-elevations/ae-1/approve", r.URL.Path)

		var body ReviewAccessElevationRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "ok", body.ReviewNote)

		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(IamAccessElevationRequest{
			Identity: "ae-1",
			Status:   IamAccessElevationStatusApproved,
		}))
	}))
	defer server.Close()

	request, err := newTestIAMClient(t, server.URL).ApproveAccessElevation(context.Background(), "ae-1", ReviewAccessElevationRequest{
		ReviewNote: "ok",
	})
	require.NoError(t, err)
	assert.Equal(t, IamAccessElevationStatusApproved, request.Status)
}

func TestRejectAccessElevation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/projects/iam/access-elevations/ae-1/reject", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(IamAccessElevationRequest{
			Identity: "ae-1",
			Status:   IamAccessElevationStatusRejected,
		}))
	}))
	defer server.Close()

	request, err := newTestIAMClient(t, server.URL).RejectAccessElevation(context.Background(), "ae-1", ReviewAccessElevationRequest{})
	require.NoError(t, err)
	assert.Equal(t, IamAccessElevationStatusRejected, request.Status)
}

func TestRevokeAccessElevation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/projects/iam/access-elevations/ae-1/revoke", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(IamAccessElevationRequest{
			Identity: "ae-1",
			Status:   IamAccessElevationStatusRevoked,
		}))
	}))
	defer server.Close()

	request, err := newTestIAMClient(t, server.URL).RevokeAccessElevation(context.Background(), "ae-1", ReviewAccessElevationRequest{
		ReviewNote: "done",
	})
	require.NoError(t, err)
	assert.Equal(t, IamAccessElevationStatusRevoked, request.Status)
}

func TestGetAccessElevation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/projects/iam/access-elevations/ae-1", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(IamAccessElevationRequest{Identity: "ae-1"}))
	}))
	defer server.Close()

	request, err := newTestIAMClient(t, server.URL).GetAccessElevation(context.Background(), "ae-1")
	require.NoError(t, err)
	assert.Equal(t, "ae-1", request.Identity)
}
