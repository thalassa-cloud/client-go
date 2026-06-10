package kms_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/thalassa-cloud/client-go/kms"
	"github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/thalassa"
)

const exampleRegion = "nl-01"

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func ExampleClient_GetSummary() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, kms.KmsSummary{
			FeatureEnabled: true,
			Regions: []kms.KmsSummaryRegion{
				{Slug: exampleRegion, KmsAvailable: true},
			},
		})
	}))
	defer server.Close()

	c, _ := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	kmsClient, _ := kms.New(c)

	summary, _ := kmsClient.GetSummary(context.Background())
	fmt.Println(summary.FeatureEnabled)
	// Output: true
}

func ExampleClient_CreateKey() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusCreated, kms.KmsKey{
			Identity: "kms-abc123",
			Name:     "app-secrets",
			KeyType:  kms.KmsKeyTypeAES256GCM96,
			Status:   kms.KmsKeyStatusActive,
		})
	}))
	defer server.Close()

	c, _ := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	kmsClient, _ := kms.New(c)

	key, _ := kmsClient.CreateKey(context.Background(), exampleRegion, kms.CreateKmsKeyRequest{
		Name:    "app-secrets",
		KeyType: kms.KmsKeyTypeAES256GCM96,
	})
	fmt.Println(key.Identity)
	// Output: kms-abc123
}

func ExampleClient_EncryptBytes() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/kms/nl-01/keys/kms-abc123/encrypt":
			writeJSON(w, http.StatusOK, kms.EncryptResponse{
				Ciphertext: "thalassa:v1:encrypted",
			})
		case "/v1/kms/nl-01/keys/kms-abc123/decrypt":
			writeJSON(w, http.StatusOK, kms.DecryptResponse{
				Plaintext: kms.EncodeBytes([]byte("hello")),
			})
		}
	}))
	defer server.Close()

	c, _ := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	kmsClient, _ := kms.New(c)
	ctx := context.Background()

	enc, _ := kmsClient.EncryptBytes(ctx, exampleRegion, "kms-abc123", []byte("hello"))
	plain, _ := kmsClient.DecryptBytes(ctx, exampleRegion, "kms-abc123", enc.Ciphertext)
	fmt.Printf("%s %q\n", enc.Ciphertext, plain)
	// Output: thalassa:v1:encrypted "hello"
}

func ExampleEncodeBytes() {
	encoded := kms.EncodeBytes([]byte("hello"))
	decoded, _ := kms.DecodeBytes("plaintext", encoded)
	fmt.Printf("%s %q\n", encoded, decoded)
	// Output: aGVsbG8= "hello"
}

func Example() {
	// Full workflow using the Thalassa facade client.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/kms/summary":
			writeJSON(w, http.StatusOK, kms.KmsSummary{FeatureEnabled: true})
		case "/v1/kms/nl-01/keys":
			writeJSON(w, http.StatusCreated, kms.KmsKey{
				Identity: "kms-abc123",
				Name:     "demo-key",
				Status:   kms.KmsKeyStatusActive,
			})
		case "/v1/kms/nl-01/keys/kms-abc123/encrypt":
			writeJSON(w, http.StatusOK, kms.EncryptResponse{Ciphertext: "thalassa:v1:encrypted"})
		case "/v1/kms/nl-01/keys/kms-abc123/decrypt":
			writeJSON(w, http.StatusOK, kms.DecryptResponse{Plaintext: kms.EncodeBytes([]byte("hello"))})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	tc, _ := thalassa.NewClient(
		client.WithBaseURL(server.URL),
		client.WithAuthCustom(),
	)
	ctx := context.Background()
	kmsClient := tc.KMS()

	summary, _ := kmsClient.GetSummary(ctx)
	if !summary.FeatureEnabled {
		panic("KMS not enabled")
	}

	key, _ := kmsClient.CreateKey(ctx, exampleRegion, kms.CreateKmsKeyRequest{
		Name:    "demo-key",
		KeyType: kms.KmsKeyTypeAES256GCM96,
	})
	enc, _ := kmsClient.EncryptBytes(ctx, exampleRegion, key.Identity, []byte("hello"))
	plain, _ := kmsClient.DecryptBytes(ctx, exampleRegion, key.Identity, enc.Ciphertext)
	fmt.Printf("round-trip OK: %q\n", plain)
	// Output: round-trip OK: "hello"
}
