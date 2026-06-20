package secrets

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

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      string
		wantErr   bool
		errSubstr string
	}{
		{
			name:  "adds leading slash",
			input: "app/db",
			want:  "/app/db",
		},
		{
			name:  "preserves leading slash",
			input: "/app/prod/db/password",
			want:  "/app/prod/db/password",
		},
		{
			name:  "allows plus sign",
			input: "/app/a+b",
			want:  "/app/a+b",
		},
		{
			name:      "empty path",
			input:     "",
			wantErr:   true,
			errSubstr: "path cannot be empty",
		},
		{
			name:      "invalid characters",
			input:     "/app/secret?query=1",
			wantErr:   true,
			errSubstr: "invalid secret path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizePath(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errSubstr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSecretResourceURL(t *testing.T) {
	tests := []struct {
		name    string
		region  string
		path    string
		suffix  string
		want    string
		wantErr bool
	}{
		{
			name:   "secret metadata",
			region: "nl-01",
			path:   "/app/prod/db/password",
			want:   "/v1/secrets/nl-01/secret/app/prod/db/password",
		},
		{
			name:   "secret value",
			region: "nl-01",
			path:   "/app/prod/db/password",
			suffix: "/value",
			want:   "/v1/secrets/nl-01/secret/app/prod/db/password/value",
		},
		{
			name:   "path with plus",
			region: "nl-01",
			path:   "/app/a+b",
			suffix: "/policy",
			want:   "/v1/secrets/nl-01/secret/app/a+b/policy",
		},
		{
			name:    "invalid path",
			region:  "nl-01",
			path:    "bad path",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SecretResourceURL(tt.region, tt.path, tt.suffix)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBrowseSecrets(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v1/secrets/nl-01/secrets", r.URL.Path)
		assert.Equal(t, "/app/prod/", r.URL.Query().Get("path"))
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(BrowseSecretsResponse{
			Path:     "/app/prod/",
			Prefixes: []string{"/app/prod/db/"},
			Secrets: []Secret{
				{Path: "/app/prod/config", CurrentVersion: 1},
			},
		}))
	}))
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	require.NoError(t, err)
	secretsClient, err := New(c)
	require.NoError(t, err)

	browse, err := secretsClient.BrowseSecrets(context.Background(), "nl-01", "/app/prod/")
	require.NoError(t, err)
	assert.Equal(t, "/app/prod/", browse.Path)
	assert.Equal(t, []string{"/app/prod/db/"}, browse.Prefixes)
	require.Len(t, browse.Secrets, 1)
}

func TestListSecrets(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v1/secrets/nl-01/secrets", r.URL.Path)
		assert.Equal(t, "/app/prod/", r.URL.Query().Get("pathPrefix"))
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode([]Secret{
			{Path: "/app/prod/db/password", CurrentVersion: 2},
		}))
	}))
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	require.NoError(t, err)
	secretsClient, err := New(c)
	require.NoError(t, err)

	secrets, err := secretsClient.ListSecrets(context.Background(), "nl-01", "app/prod/")
	require.NoError(t, err)
	require.Len(t, secrets, 1)
	assert.Equal(t, "/app/prod/db/password", secrets[0].Path)
}

func TestCreateSecret(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v1/secrets/nl-01/secrets", r.URL.Path)

		var body CreateSecretRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "/app/prod/db/password", body.Path)
		assert.Equal(t, "kms-abc123", body.KmsKeyIdentity)
		assert.NotEmpty(t, body.SecretString)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		require.NoError(t, json.NewEncoder(w).Encode(Secret{
			Path:           body.Path,
			CurrentVersion: 1,
		}))
	}))
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	require.NoError(t, err)
	secretsClient, err := New(c)
	require.NoError(t, err)

	secret, err := secretsClient.CreateSecret(context.Background(), "nl-01", CreateSecretRequest{
		Path:           "/app/prod/db/password",
		KmsKeyIdentity: "kms-abc123",
		SecretString:   EncodeBytes([]byte("super-secret")),
	})
	require.NoError(t, err)
	assert.Equal(t, "/app/prod/db/password", secret.Path)
}

func TestGetSecret(t *testing.T) {
	tests := []struct {
		name            string
		includeVersions bool
		wantQuery       string
	}{
		{
			name:      "metadata only",
			wantQuery: "",
		},
		{
			name:            "with versions",
			includeVersions: true,
			wantQuery:       "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "/v1/secrets/nl-01/secret/app/prod/db/password", r.URL.Path)
				assert.Equal(t, tt.wantQuery, r.URL.Query().Get("includeVersions"))
				w.Header().Set("Content-Type", "application/json")
				require.NoError(t, json.NewEncoder(w).Encode(Secret{
					Path:           "/app/prod/db/password",
					CurrentVersion: 1,
				}))
			}))
			defer server.Close()

			c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
			require.NoError(t, err)
			secretsClient, err := New(c)
			require.NoError(t, err)

			secret, err := secretsClient.GetSecret(context.Background(), "nl-01", "/app/prod/db/password", tt.includeVersions)
			require.NoError(t, err)
			assert.Equal(t, "/app/prod/db/password", secret.Path)
		})
	}
}

func TestPutAndGetSecretString(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/secrets/nl-01/secret/app/prod/db/password/versions":
			var body PutSecretValueRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			assert.Equal(t, "/app/prod/db/password", body.Path)
			assert.Equal(t, EncodeBytes([]byte("super-secret")), body.SecretString)
			require.NoError(t, json.NewEncoder(w).Encode(PutSecretValueResponse{
				Path:    "/app/prod/db/password",
				Version: 2,
			}))
		case r.Method == http.MethodPost && r.URL.Path == "/v1/secrets/nl-01/secret/app/prod/db/password/value":
			var body GetSecretValueRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			assert.Equal(t, "/app/prod/db/password", body.Path)
			require.NoError(t, json.NewEncoder(w).Encode(GetSecretValueResponse{
				Path:           "/app/prod/db/password",
				Version:        2,
				SecretString:   EncodeBytes([]byte("super-secret")),
				KmsKeyIdentity: "kms-abc123",
				KmsKeyVersion:  "1",
			}))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	require.NoError(t, err)
	secretsClient, err := New(c)
	require.NoError(t, err)

	putResp, err := secretsClient.PutSecretString(context.Background(), "nl-01", "/app/prod/db/password", []byte("super-secret"))
	require.NoError(t, err)
	assert.Equal(t, 2, putResp.Version)

	val, version, err := secretsClient.GetSecretString(context.Background(), "nl-01", "/app/prod/db/password", nil)
	require.NoError(t, err)
	assert.Equal(t, 2, version)
	assert.Equal(t, []byte("super-secret"), val)
}

func TestDeleteSecret(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/v1/secrets/nl-01/secret/app/prod/db/password", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	require.NoError(t, err)
	secretsClient, err := New(c)
	require.NoError(t, err)

	err = secretsClient.DeleteSecret(context.Background(), "nl-01", "/app/prod/db/password")
	require.NoError(t, err)
}

func TestDestroySecretVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/v1/secrets/nl-01/secret/app/prod/db/password/versions", r.URL.Path)
		assert.Equal(t, "3", r.URL.Query().Get("version"))
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	require.NoError(t, err)
	secretsClient, err := New(c)
	require.NoError(t, err)

	err = secretsClient.DestroySecretVersion(context.Background(), "nl-01", "/app/prod/db/password", 3)
	require.NoError(t, err)
}

func TestUpdateAccessPolicy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/v1/secrets/nl-01/secret/app/prod/db/password/policy", r.URL.Path)

		var body UpdateAccessPolicyRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.Len(t, body.AccessPolicy.Statements, 1)

		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(Secret{
			Path:         "/app/prod/db/password",
			AccessPolicy: &body.AccessPolicy,
		}))
	}))
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	require.NoError(t, err)
	secretsClient, err := New(c)
	require.NoError(t, err)

	secret, err := secretsClient.UpdateAccessPolicy(context.Background(), "nl-01", "/app/prod/db/password", UpdateAccessPolicyRequest{
		AccessPolicy: SecretPolicy{
			Statements: []SecretPolicyStatement{
				{Effect: "allow", Actions: []string{"secrets:get-value"}},
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, secret.AccessPolicy)
}

func TestEncodeDecodeBytes(t *testing.T) {
	encoded := EncodeBytes([]byte("hello"))
	decoded, err := DecodeBytes("secretString", encoded)
	require.NoError(t, err)
	assert.Equal(t, []byte("hello"), decoded)
}
