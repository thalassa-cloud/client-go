package secrets_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/thalassa-cloud/client-go/kms"
	"github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/secrets"
	"github.com/thalassa-cloud/client-go/thalassa"
)

const exampleRegion = "nl-01"

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func ExampleSecretResourceURL() {
	url, _ := secrets.SecretResourceURL(exampleRegion, "/app/prod/db/password", "/value")
	fmt.Println(url)
	// Output: /v1/secrets/nl-01/secret/app/prod/db/password/value
}

func ExampleClient_CreateSecret() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusCreated, secrets.Secret{
			Path:           "/app/prod/db/password",
			KmsKeyIdentity: "kms-abc123",
			CurrentVersion: 1,
		})
	}))
	defer server.Close()

	c, _ := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	secretsClient, _ := secrets.New(c)

	secret, _ := secretsClient.CreateSecret(context.Background(), exampleRegion, secrets.CreateSecretRequest{
		Path:           "/app/prod/db/password",
		KmsKeyIdentity: "kms-abc123",
		SecretString:   secrets.EncodeBytes([]byte("super-secret")),
	})
	fmt.Printf("%s v%d\n", secret.Path, secret.CurrentVersion)
	// Output: /app/prod/db/password v1
}

func ExampleClient_GetSecretString() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, secrets.GetSecretValueResponse{
			Path:           "/app/prod/db/password",
			Version:        2,
			SecretString:   secrets.EncodeBytes([]byte("super-secret")),
			KmsKeyIdentity: "kms-abc123",
			KmsKeyVersion:  "1",
		})
	}))
	defer server.Close()

	c, _ := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	secretsClient, _ := secrets.New(c)

	val, version, _ := secretsClient.GetSecretString(context.Background(), exampleRegion, "/app/prod/db/password", nil)
	fmt.Printf("v%d %d bytes\n", version, len(val))
	// Output: v2 12 bytes
}

func ExampleClient_BrowseSecrets() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, secrets.BrowseSecretsResponse{
			Path:     "/app/prod/",
			Prefixes: []string{"/app/prod/db/"},
			Secrets: []secrets.Secret{
				{Path: "/app/prod/config", CurrentVersion: 1},
			},
		})
	}))
	defer server.Close()

	c, _ := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	secretsClient, _ := secrets.New(c)

	browse, _ := secretsClient.BrowseSecrets(context.Background(), exampleRegion, "/app/prod/")
	fmt.Printf("%d prefixes, %d secrets\n", len(browse.Prefixes), len(browse.Secrets))
	// Output: 1 prefixes, 1 secrets
}

func Example() {
	// End-to-end flow: create a KMS key, store a secret, read it back.
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/kms/summary", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, kms.KmsSummary{FeatureEnabled: true})
	})
	mux.HandleFunc("/v1/kms/nl-01/keys", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusCreated, kms.KmsKey{
			Identity: "kms-abc123",
			Name:     "app-secrets",
			Status:   kms.KmsKeyStatusActive,
		})
	})
	mux.HandleFunc("/v1/secrets/nl-01/secrets", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			writeJSON(w, http.StatusCreated, secrets.Secret{
				Path:           "/app/prod/db/password",
				CurrentVersion: 1,
			})
		}
	})
	mux.HandleFunc("/v1/secrets/nl-01/secret/app/prod/db/password/value", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, secrets.GetSecretValueResponse{
			Path:         "/app/prod/db/password",
			Version:      1,
			SecretString: secrets.EncodeBytes([]byte("initial-password")),
		})
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	tc, _ := thalassa.NewClient(
		client.WithBaseURL(server.URL),
		client.WithAuthCustom(),
	)
	ctx := context.Background()

	key, _ := tc.KMS().CreateKey(ctx, exampleRegion, kms.CreateKmsKeyRequest{
		Name:    "app-secrets",
		KeyType: kms.KmsKeyTypeAES256GCM96,
	})
	_, _ = tc.Secrets().CreateSecret(ctx, exampleRegion, secrets.CreateSecretRequest{
		Path:           "/app/prod/db/password",
		KmsKeyIdentity: key.Identity,
		SecretString:   secrets.EncodeBytes([]byte("initial-password")),
	})
	val, version, _ := tc.Secrets().GetSecretString(ctx, exampleRegion, "/app/prod/db/password", nil)
	fmt.Printf("read v%d (%d bytes)\n", version, len(val))
	// Output: read v1 (16 bytes)
}
