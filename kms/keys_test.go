package kms

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

func TestEncodeDecodeBytes(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		field     string
		encoded   string
		wantErr   bool
		errSubstr string
	}{
		{
			name:  "round trip",
			input: []byte("hello"),
		},
		{
			name:      "invalid base64",
			field:     "plaintext",
			encoded:   "not-valid-base64!!!",
			wantErr:   true,
			errSubstr: "plaintext must be valid base64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.encoded != "" {
				_, err := DecodeBytes(tt.field, tt.encoded)
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errSubstr)
				return
			}

			encoded := EncodeBytes(tt.input)
			decoded, err := DecodeBytes("plaintext", encoded)
			require.NoError(t, err)
			assert.Equal(t, tt.input, decoded)
		})
	}
}

func TestGetSummary(t *testing.T) {
	tests := []struct {
		name         string
		serverStatus int
		response     KmsSummary
		expectError  bool
	}{
		{
			name:         "feature enabled",
			serverStatus: http.StatusOK,
			response: KmsSummary{
				FeatureEnabled: true,
				Regions: []KmsSummaryRegion{
					{Slug: "nl-01", KmsAvailable: true},
				},
			},
		},
		{
			name:         "not found",
			serverStatus: http.StatusNotFound,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "/v1/kms/summary", r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.serverStatus)
				if tt.serverStatus == http.StatusOK {
					require.NoError(t, json.NewEncoder(w).Encode(tt.response))
				}
			}))
			defer server.Close()

			c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
			require.NoError(t, err)
			kmsClient, err := New(c)
			require.NoError(t, err)

			summary, err := kmsClient.GetSummary(context.Background())
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.response.FeatureEnabled, summary.FeatureEnabled)
			assert.Len(t, summary.Regions, len(tt.response.Regions))
		})
	}
}

func TestCreateKey(t *testing.T) {
	tests := []struct {
		name         string
		request      CreateKmsKeyRequest
		serverStatus int
		expectError  bool
	}{
		{
			name: "create symmetric key",
			request: CreateKmsKeyRequest{
				Name:    "app-secrets",
				KeyType: KmsKeyTypeAES256GCM96,
			},
			serverStatus: http.StatusCreated,
		},
		{
			name: "byok import",
			request: CreateKmsKeyRequest{
				Name:              "imported-key",
				KeyType:           KmsKeyTypeAES256GCM96,
				ImportKeyMaterial: "wrapped-key-material",
				HashFunction:      "sha256",
				AllowRotation:     true,
			},
			serverStatus: http.StatusCreated,
		},
		{
			name: "conflict",
			request: CreateKmsKeyRequest{
				Name:    "duplicate",
				KeyType: KmsKeyTypeAES256GCM96,
			},
			serverStatus: http.StatusConflict,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/v1/kms/nl-01/keys", r.URL.Path)

				var body CreateKmsKeyRequest
				require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
				assert.Equal(t, tt.request.Name, body.Name)
				assert.Equal(t, tt.request.KeyType, body.KeyType)
				if tt.request.ImportKeyMaterial != "" {
					assert.Equal(t, tt.request.ImportKeyMaterial, body.ImportKeyMaterial)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.serverStatus)
				if tt.serverStatus == http.StatusCreated {
					require.NoError(t, json.NewEncoder(w).Encode(KmsKey{
						Identity: "kms-abc123",
						Name:     body.Name,
						KeyType:  body.KeyType,
						Status:   KmsKeyStatusActive,
					}))
				}
			}))
			defer server.Close()

			c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
			require.NoError(t, err)
			kmsClient, err := New(c)
			require.NoError(t, err)

			key, err := kmsClient.CreateKey(context.Background(), "nl-01", tt.request)
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, "kms-abc123", key.Identity)
			assert.Equal(t, tt.request.Name, key.Name)
		})
	}
}

func TestGetKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v1/kms/nl-01/keys/kms-abc123", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(KmsKey{
			Identity: "kms-abc123",
			Name:     "app-secrets",
			Status:   KmsKeyStatusActive,
		}))
	}))
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	require.NoError(t, err)
	kmsClient, err := New(c)
	require.NoError(t, err)

	key, err := kmsClient.GetKey(context.Background(), "nl-01", "kms-abc123")
	require.NoError(t, err)
	assert.Equal(t, "kms-abc123", key.Identity)
}

func TestListKeys(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v1/kms/nl-01/keys", r.URL.Path)
		assert.Equal(t, "app-secrets", r.URL.Query().Get("name"))
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode([]KmsKey{
			{Identity: "kms-1", Name: "app-secrets"},
		}))
	}))
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	require.NoError(t, err)
	kmsClient, err := New(c)
	require.NoError(t, err)

	keys, err := kmsClient.ListKeys(context.Background(), "nl-01", &ListKeysRequest{
		Filters: []filters.Filter{
			&filters.FilterKeyValue{
				Key:   filters.FilterName,
				Value: "app-secrets",
			},
		},
	})
	require.NoError(t, err)
	require.Len(t, keys, 1)
	assert.Equal(t, "kms-1", keys[0].Identity)
}

func TestEncryptDecrypt(t *testing.T) {
	var capturedEncrypt EncryptRequest
	var capturedDecrypt DecryptRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v1/kms/nl-01/keys/kms-abc123/encrypt":
			require.NoError(t, json.NewDecoder(r.Body).Decode(&capturedEncrypt))
			require.NoError(t, json.NewEncoder(w).Encode(EncryptResponse{
				Ciphertext: "thalassa:v1:encrypted",
				KeyVersion: "1",
			}))
		case "/v1/kms/nl-01/keys/kms-abc123/decrypt":
			require.NoError(t, json.NewDecoder(r.Body).Decode(&capturedDecrypt))
			require.NoError(t, json.NewEncoder(w).Encode(DecryptResponse{
				Plaintext:  EncodeBytes([]byte("hello")),
				KeyVersion: "1",
			}))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	require.NoError(t, err)
	kmsClient, err := New(c)
	require.NoError(t, err)

	enc, err := kmsClient.EncryptBytes(context.Background(), "nl-01", "kms-abc123", []byte("hello"))
	require.NoError(t, err)
	assert.Equal(t, "thalassa:v1:encrypted", enc.Ciphertext)
	assert.Equal(t, EncodeBytes([]byte("hello")), capturedEncrypt.Plaintext)

	dec, err := kmsClient.Decrypt(context.Background(), "nl-01", "kms-abc123", DecryptRequest{
		Ciphertext: enc.Ciphertext,
	})
	require.NoError(t, err)
	assert.Equal(t, enc.Ciphertext, capturedDecrypt.Ciphertext)

	plain, err := kmsClient.DecryptBytes(context.Background(), "nl-01", "kms-abc123", enc.Ciphertext)
	require.NoError(t, err)
	assert.Equal(t, []byte("hello"), plain)
	assert.Equal(t, EncodeBytes([]byte("hello")), dec.Plaintext)
}

func TestDeleteKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/v1/kms/nl-01/keys/kms-abc123", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	require.NoError(t, err)
	kmsClient, err := New(c)
	require.NoError(t, err)

	err = kmsClient.DeleteKey(context.Background(), "nl-01", "kms-abc123")
	require.NoError(t, err)
}

func TestRegionPath(t *testing.T) {
	assert.Equal(t, "/v1/kms/nl-01/keys/kms-abc123/encrypt", regionPath("nl-01", "keys", "kms-abc123", "encrypt"))
}
