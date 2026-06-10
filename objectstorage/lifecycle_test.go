package objectstorage

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

func int64Ptr(v int64) *int64 {
	return &v
}

func TestGetBucketLifecycle(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v1/object-storage/buckets/my-app-logs/lifecycle", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(BucketLifecycle{
			Rules: []BucketLifecycleRule{
				{
					ID:     "expire-logs",
					Prefix: "logs/",
					Status: BucketLifecycleRuleStatusEnabled,
					Expiration: &BucketLifecycleRuleExpiration{
						Days: int64Ptr(30),
					},
				},
			},
		}))
	}))
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	require.NoError(t, err)
	osClient, err := New(c)
	require.NoError(t, err)

	lifecycle, err := osClient.GetBucketLifecycle(context.Background(), "my-app-logs")
	require.NoError(t, err)
	require.Len(t, lifecycle.Rules, 1)
	assert.Equal(t, "expire-logs", lifecycle.Rules[0].ID)
	assert.Equal(t, int64(30), *lifecycle.Rules[0].Expiration.Days)
}

func TestGetBucketLifecycleEmptyBucketName(t *testing.T) {
	c, err := client.NewClient(
		client.WithBaseURL("http://example.com"),
		client.WithAuthCustom(),
	)
	require.NoError(t, err)
	osClient, err := New(c)
	require.NoError(t, err)

	_, err = osClient.GetBucketLifecycle(context.Background(), "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "bucket name is required")
}

func TestSetBucketLifecycle(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/v1/object-storage/buckets/my-app-logs/lifecycle", r.URL.Path)

		var body SetBucketLifecycleRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.Len(t, body.Rules, 1)
		assert.Equal(t, "abort-incomplete-uploads", body.Rules[0].ID)
		assert.Equal(t, int64(7), *body.Rules[0].AbortIncompleteMultipartUpload.DaysAfterInitiation)

		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(BucketLifecycle{
			Rules: body.Rules,
		}))
	}))
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	require.NoError(t, err)
	osClient, err := New(c)
	require.NoError(t, err)

	lifecycle, err := osClient.SetBucketLifecycle(context.Background(), "my-app-logs", SetBucketLifecycleRequest{
		Rules: []BucketLifecycleRule{
			{
				ID:     "abort-incomplete-uploads",
				Prefix: "uploads/",
				AbortIncompleteMultipartUpload: &BucketLifecycleRuleAbortIncompleteMultipartUpload{
					DaysAfterInitiation: int64Ptr(7),
				},
			},
		},
	})
	require.NoError(t, err)
	require.Len(t, lifecycle.Rules, 1)
}

func TestDeleteBucketLifecycle(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/v1/object-storage/buckets/my-app-logs/lifecycle", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	require.NoError(t, err)
	osClient, err := New(c)
	require.NoError(t, err)

	err = osClient.DeleteBucketLifecycle(context.Background(), "my-app-logs")
	require.NoError(t, err)
}

func TestCreateBucketWithLifecycle(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v1/object-storage/buckets", r.URL.Path)

		var body CreateBucketRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.NotNil(t, body.Lifecycle)
		require.Len(t, body.Lifecycle.Rules, 1)
		assert.Equal(t, "expire-logs", body.Lifecycle.Rules[0].ID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		require.NoError(t, json.NewEncoder(w).Encode(ObjectStorageBucket{
			Name: "my-app-logs",
			Lifecycle: &BucketLifecycle{
				Rules: body.Lifecycle.Rules,
			},
		}))
	}))
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	require.NoError(t, err)
	osClient, err := New(c)
	require.NoError(t, err)

	bucket, err := osClient.CreateBucket(context.Background(), CreateBucketRequest{
		BucketName: "my-app-logs",
		Region:     "nl-01",
		Versioning: ObjectStorageBucketVersioningEnabled,
		Lifecycle: &SetBucketLifecycleRequest{
			Rules: []BucketLifecycleRule{
				{
					ID:     "expire-logs",
					Prefix: "logs/",
					Status: BucketLifecycleRuleStatusEnabled,
					Expiration: &BucketLifecycleRuleExpiration{
						Days: int64Ptr(30),
					},
				},
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, bucket.Lifecycle)
	assert.Equal(t, "expire-logs", bucket.Lifecycle.Rules[0].ID)
}

func TestBucketLifecyclePath(t *testing.T) {
	assert.Equal(t, "/v1/object-storage/buckets/my-bucket/lifecycle", bucketLifecyclePath("my-bucket"))
}
