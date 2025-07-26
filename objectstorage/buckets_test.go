package objectstorage

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/client-go/pkg/base"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

func TestListBuckets(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse []ObjectStorageBucket
		serverStatus   int
		expectError    bool
		errorMessage   string
	}{
		{
			name: "successful list buckets",
			serverResponse: []ObjectStorageBucket{
				{
					Identity: "bucket-1",
					Name:     "test-bucket-1",
					Public:   false,
					Status:   "active",
					Endpoint: "https://test-bucket-1.s3.thalasascloud.nl",
				},
				{
					Identity: "bucket-2",
					Name:     "test-bucket-2",
					Public:   true,
					Status:   "active",
					Endpoint: "https://test-bucket-2.s3.thalasascloud.nl",
				},
			},
			serverStatus: http.StatusOK,
			expectError:  false,
		},
		{
			name:           "empty list",
			serverResponse: []ObjectStorageBucket{},
			serverStatus:   http.StatusOK,
			expectError:    false,
		},
		{
			name:         "server error",
			serverStatus: http.StatusInternalServerError,
			expectError:  true,
			errorMessage: "server error",
		},
		{
			name:         "unauthorized",
			serverStatus: http.StatusUnauthorized,
			expectError:  true,
			errorMessage: "unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method)
				assert.Equal(t, "/v1/object-storage/buckets", r.URL.Path)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.serverStatus)
				if tt.serverStatus == http.StatusOK {
					json.NewEncoder(w).Encode(tt.serverResponse)
				} else {
					json.NewEncoder(w).Encode(map[string]string{"message": tt.errorMessage})
				}
			}))
			defer server.Close()

			c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
			require.NoError(t, err)

			osClient, err := New(c)
			require.NoError(t, err)

			buckets, err := osClient.ListBuckets(context.Background())

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
			} else {
				assert.NoError(t, err)
				assert.Len(t, buckets, len(tt.serverResponse))
				for i, expected := range tt.serverResponse {
					assert.Equal(t, expected.Identity, buckets[i].Identity)
					assert.Equal(t, expected.Name, buckets[i].Name)
					assert.Equal(t, expected.Public, buckets[i].Public)
					assert.Equal(t, expected.Status, buckets[i].Status)
					assert.Equal(t, expected.Endpoint, buckets[i].Endpoint)
				}
			}
		})
	}
}

func TestGetBucket(t *testing.T) {
	tests := []struct {
		name           string
		bucketName     string
		serverResponse ObjectStorageBucket
		serverStatus   int
		expectError    bool
		errorMessage   string
	}{
		{
			name:       "successful get bucket",
			bucketName: "test-bucket",
			serverResponse: ObjectStorageBucket{
				Identity: "bucket-1",
				Name:     "test-bucket",
				Public:   false,
				Status:   "active",
				Endpoint: "https://test-bucket.s3.thalasascloud.nl",
			},
			serverStatus: http.StatusOK,
			expectError:  false,
		},
		{
			name:         "bucket not found",
			bucketName:   "non-existent-bucket",
			serverStatus: http.StatusNotFound,
			expectError:  true,
			errorMessage: "bucket not found",
		},
		{
			name:         "server error",
			bucketName:   "test-bucket",
			serverStatus: http.StatusInternalServerError,
			expectError:  true,
			errorMessage: "server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method)
				assert.Equal(t, "/v1/object-storage/buckets/"+tt.bucketName, r.URL.Path)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.serverStatus)
				if tt.serverStatus == http.StatusOK {
					json.NewEncoder(w).Encode(tt.serverResponse)
				} else {
					json.NewEncoder(w).Encode(map[string]string{"message": tt.errorMessage})
				}
			}))
			defer server.Close()

			c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
			require.NoError(t, err)

			osClient, err := New(c)
			require.NoError(t, err)

			bucket, err := osClient.GetBucket(context.Background(), tt.bucketName)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, bucket)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, bucket)
				assert.Equal(t, tt.serverResponse.Identity, bucket.Identity)
				assert.Equal(t, tt.serverResponse.Name, bucket.Name)
				assert.Equal(t, tt.serverResponse.Public, bucket.Public)
				assert.Equal(t, tt.serverResponse.Status, bucket.Status)
				assert.Equal(t, tt.serverResponse.Endpoint, bucket.Endpoint)
			}
		})
	}
}

func TestCreateBucket(t *testing.T) {
	createRequest := CreateBucketRequest{
		BucketName: "new-bucket-",
		Public:     false,
		Region:     "us-east-1",
		PolicyDocument: &PolicyDocument{
			Version: "2012-10-17",
			Statement: []Statement{
				{
					Sid:    "PublicReadGetObject",
					Effect: "Allow",
					Principal: Principal{
						AWS: "*",
					},
					Action:   "s3:GetObject",
					Resource: []string{"arn:aws:s3:::new-bucket/*"},
				},
			},
		},
	}

	tests := []struct {
		name           string
		createRequest  CreateBucketRequest
		serverResponse ObjectStorageBucket
		serverStatus   int
		expectError    bool
		errorMessage   string
	}{
		{
			name:          "successful create bucket",
			createRequest: createRequest,
			serverResponse: ObjectStorageBucket{
				Identity: "bucket-new",
				Name:     "new-bucket",
				Public:   false,
				Status:   "creating",
				Endpoint: "https://new-bucket.s3.thalasascloud.nl",
			},
			serverStatus: http.StatusCreated,
			expectError:  false,
		},
		{
			name: "create bucket without policy",
			createRequest: CreateBucketRequest{
				BucketName: "simple-bucket",
				Public:     true,
				Region:     "us-west-2",
			},
			serverResponse: ObjectStorageBucket{
				Identity: "bucket-simple",
				Name:     "simple-bucket",
				Public:   true,
				Status:   "creating",
			},
			serverStatus: http.StatusCreated,
			expectError:  false,
		},
		{
			name: "bucket name conflict",
			createRequest: CreateBucketRequest{
				BucketName: "existing-bucket",
				Public:     false,
				Region:     "us-east-1",
			},
			serverStatus: http.StatusConflict,
			expectError:  true,
			errorMessage: "bucket already exists",
		},
		{
			name: "invalid region",
			createRequest: CreateBucketRequest{
				BucketName: "invalid-bucket",
				Public:     false,
				Region:     "invalid-region",
			},
			serverStatus: http.StatusBadRequest,
			expectError:  true,
			errorMessage: "invalid region",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, "/v1/object-storage/buckets", r.URL.Path)

				// Verify request body
				var receivedRequest CreateBucketRequest
				err := json.NewDecoder(r.Body).Decode(&receivedRequest)
				assert.NoError(t, err)
				assert.Equal(t, tt.createRequest, receivedRequest)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.serverStatus)
				if tt.serverStatus == http.StatusCreated {
					json.NewEncoder(w).Encode(tt.serverResponse)
				} else {
					json.NewEncoder(w).Encode(map[string]string{"message": tt.errorMessage})
				}
			}))
			defer server.Close()

			c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
			require.NoError(t, err)

			osClient, err := New(c)
			require.NoError(t, err)

			bucket, err := osClient.CreateBucket(context.Background(), tt.createRequest)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, bucket)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, bucket)
				assert.Equal(t, tt.serverResponse.Identity, bucket.Identity)
				assert.Equal(t, tt.serverResponse.Name, bucket.Name)
				assert.Equal(t, tt.serverResponse.Public, bucket.Public)
				assert.Equal(t, tt.serverResponse.Status, bucket.Status)
				assert.Equal(t, tt.serverResponse.Endpoint, bucket.Endpoint)
			}
		})
	}
}

func TestUpdateBucket(t *testing.T) {
	updateRequest := UpdateBucketRequest{
		Public: true,
		PolicyDocument: &PolicyDocument{
			Version: "2012-10-17",
			Statement: []Statement{
				{
					Sid:    "PublicReadGetObject",
					Effect: "Allow",
					Principal: Principal{
						AWS: "*",
					},
					Action:   "s3:GetObject",
					Resource: []string{"arn:aws:s3:::test-bucket/*"},
				},
			},
		},
	}

	tests := []struct {
		name           string
		bucketName     string
		updateRequest  UpdateBucketRequest
		serverResponse ObjectStorageBucket
		serverStatus   int
		expectError    bool
		errorMessage   string
	}{
		{
			name:          "successful update bucket",
			bucketName:    "test-bucket",
			updateRequest: updateRequest,
			serverResponse: ObjectStorageBucket{
				Identity: "bucket-1",
				Name:     "test-bucket",
				Public:   true,
				Status:   "active",
				Endpoint: "https://test-bucket.s3.thalasascloud.nl",
			},
			serverStatus: http.StatusOK,
			expectError:  false,
		},
		{
			name:       "update bucket visibility only",
			bucketName: "test-bucket",
			updateRequest: UpdateBucketRequest{
				Public: false,
			},
			serverResponse: ObjectStorageBucket{
				Identity: "bucket-1",
				Name:     "test-bucket",
				Public:   false,
				Status:   "active",
			},
			serverStatus: http.StatusOK,
			expectError:  false,
		},
		{
			name:       "bucket not found",
			bucketName: "non-existent-bucket",
			updateRequest: UpdateBucketRequest{
				Public: true,
			},
			serverStatus: http.StatusNotFound,
			expectError:  true,
			errorMessage: "bucket not found",
		},
		{
			name:       "invalid policy",
			bucketName: "test-bucket",
			updateRequest: UpdateBucketRequest{
				Public: true,
				PolicyDocument: &PolicyDocument{
					Version: "invalid",
				},
			},
			serverStatus: http.StatusBadRequest,
			expectError:  true,
			errorMessage: "invalid policy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "PUT", r.Method)
				assert.Equal(t, "/v1/object-storage/buckets/"+tt.bucketName, r.URL.Path)

				// Verify request body
				var receivedRequest UpdateBucketRequest
				err := json.NewDecoder(r.Body).Decode(&receivedRequest)
				assert.NoError(t, err)
				assert.Equal(t, tt.updateRequest, receivedRequest)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.serverStatus)
				if tt.serverStatus == http.StatusOK {
					json.NewEncoder(w).Encode(tt.serverResponse)
				} else {
					json.NewEncoder(w).Encode(map[string]string{"message": tt.errorMessage})
				}
			}))
			defer server.Close()

			c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
			require.NoError(t, err)

			osClient, err := New(c)
			require.NoError(t, err)

			bucket, err := osClient.UpdateBucket(context.Background(), tt.bucketName, tt.updateRequest)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, bucket)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, bucket)
				assert.Equal(t, tt.serverResponse.Identity, bucket.Identity)
				assert.Equal(t, tt.serverResponse.Name, bucket.Name)
				assert.Equal(t, tt.serverResponse.Public, bucket.Public)
				assert.Equal(t, tt.serverResponse.Status, bucket.Status)
				assert.Equal(t, tt.serverResponse.Endpoint, bucket.Endpoint)
			}
		})
	}
}

func TestDeleteBucket(t *testing.T) {
	tests := []struct {
		name         string
		bucketName   string
		serverStatus int
		expectError  bool
		errorMessage string
	}{
		{
			name:         "successful delete bucket",
			bucketName:   "test-bucket",
			serverStatus: http.StatusNoContent,
			expectError:  false,
		},
		{
			name:         "bucket not found",
			bucketName:   "non-existent-bucket",
			serverStatus: http.StatusNotFound,
			expectError:  true,
			errorMessage: "bucket not found",
		},
		{
			name:         "bucket not empty",
			bucketName:   "non-empty-bucket",
			serverStatus: http.StatusConflict,
			expectError:  true,
			errorMessage: "bucket not empty",
		},
		{
			name:         "server error",
			bucketName:   "test-bucket",
			serverStatus: http.StatusInternalServerError,
			expectError:  true,
			errorMessage: "server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "DELETE", r.Method)
				assert.Equal(t, "/v1/object-storage/buckets/"+tt.bucketName, r.URL.Path)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.serverStatus)
				if tt.serverStatus != http.StatusNoContent {
					json.NewEncoder(w).Encode(map[string]string{"message": tt.errorMessage})
				}
			}))
			defer server.Close()

			c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
			require.NoError(t, err)

			osClient, err := New(c)
			require.NoError(t, err)

			err = osClient.DeleteBucket(context.Background(), tt.bucketName)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBucketWithFullData(t *testing.T) {
	// Test with a bucket that has all fields populated
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/object-storage/buckets/test-full-bucket", r.URL.Path)

		fullBucket := ObjectStorageBucket{
			Identity: "bucket-full-123",
			Organisation: &base.Organisation{
				Identity:      "org-123",
				Name:          "Test Organisation",
				Slug:          "test-org",
				CreatedAt:     time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				ObjectVersion: 1,
			},
			Name:     "test-full-bucket",
			Public:   true,
			Status:   "active",
			Endpoint: "https://test-full-bucket.s3.thalasascloud.nl",
			Region: &iaas.Region{
				Identity:      "region-123",
				Name:          "US East (N. Virginia)",
				Slug:          "us-east-1",
				Description:   "US East (N. Virginia) region",
				CreatedAt:     time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:     time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				ObjectVersion: 1,
			},
			Policy: PolicyDocument{
				Version: "2012-10-17",
				Statement: []Statement{
					{
						Effect: "Allow",
						Principal: Principal{
							AWS: "*",
						},
						Action:   "s3:GetObject",
						Resource: []string{"arn:aws:s3:::test-full-bucket/*"},
					},
				},
			},
			Usage: Usage{
				TotalSizeGB:  1024,
				TotalObjects: 10,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fullBucket)
	}))
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	require.NoError(t, err)

	osClient, err := New(c)
	require.NoError(t, err)

	bucket, err := osClient.GetBucket(context.Background(), "test-full-bucket")
	assert.NoError(t, err)
	assert.NotNil(t, bucket)

	// Verify all fields are populated correctly
	assert.Equal(t, "bucket-full-123", bucket.Identity)
	assert.Equal(t, "test-full-bucket", bucket.Name)
	assert.True(t, bucket.Public)
	assert.Equal(t, "active", bucket.Status)
	assert.Equal(t, "https://test-full-bucket.s3.thalasascloud.nl", bucket.Endpoint)

	// Verify organisation
	assert.NotNil(t, bucket.Organisation)
	assert.Equal(t, "org-123", bucket.Organisation.Identity)
	assert.Equal(t, "Test Organisation", bucket.Organisation.Name)

	// Verify region
	assert.NotNil(t, bucket.Region)
	assert.Equal(t, "region-123", bucket.Region.Identity)
	assert.Equal(t, "us-east-1", bucket.Region.Slug)

	// Verify JSON fields
	assert.NotNil(t, bucket.Policy)
	assert.NotNil(t, bucket.Usage)
}

func TestContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a slow response
		time.Sleep(100 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]ObjectStorageBucket{})
	}))
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	require.NoError(t, err)

	osClient, err := New(c)
	require.NoError(t, err)

	// Create a context that will be cancelled immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = osClient.ListBuckets(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")
}
