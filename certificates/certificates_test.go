package certificates

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

func TestCreateCertificateHasNoPrivateKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v1/certificates/nl-01/certificates", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"identity":"crt-1","name":"www","status":"issuing","selfSigned":false,"issuer":{"identity":"iss-1","name":"Public","kind":"acme","provider":"letsencrypt"},"dnsNames":["www.example.com"]}`))
	}))
	defer server.Close()

	base, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	require.NoError(t, err)
	api, err := New(base)
	require.NoError(t, err)

	created, err := api.CreateCertificate(context.Background(), "nl-01", CreateCertificateRequest{
		Name:           "www",
		DNSNames:       []string{"www.example.com"},
		IssuerIdentity: "iss-1",
	})
	require.NoError(t, err)
	assert.Equal(t, "crt-1", created.Identity)
	assert.Equal(t, "issuing", created.Status)

	encoded, err := json.Marshal(created)
	require.NoError(t, err)
	assert.NotContains(t, string(encoded), `"privateKey"`)
	assert.NotContains(t, strings.ToLower(string(encoded)), "private_key")

	for _, typ := range []reflect.Type{reflect.TypeOf(Certificate{}), reflect.TypeOf(PEM{}), reflect.TypeOf(CertificateAuthority{}), reflect.TypeOf(AuthorityPEM{})} {
		for i := 0; i < typ.NumField(); i++ {
			name := typ.Field(i).Name
			if name == "AllowPrivateKeyExport" {
				continue
			}
			assert.NotContains(t, strings.ToLower(name), "private")
		}
	}
}

func TestExportPrivateKey(t *testing.T) {
	pemText := "-----BEGIN PRIVATE KEY-----\nsecret\n-----END PRIVATE KEY-----\n"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v1/certificates/nl-01/certificates/crt-1/private-key", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"privateKey":"` + base64.StdEncoding.EncodeToString([]byte(pemText)) + `"}`))
	}))
	defer server.Close()

	base, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	require.NoError(t, err)
	api, err := New(base)
	require.NoError(t, err)

	got, err := api.ExportPrivateKey(context.Background(), "nl-01", "crt-1")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, pemText, got.PrivateKey)
}

func TestRotateAndCrossSignAuthority(t *testing.T) {
	pemText := "-----BEGIN CERTIFICATE-----\nMIIB\n-----END CERTIFICATE-----\n"
	encoded := base64.StdEncoding.EncodeToString([]byte(pemText))
	tests := []struct {
		name   string
		method string
		path   string
		call   func(api *Client) error
	}{
		{
			name:   "rotate",
			method: http.MethodPost,
			path:   "/v1/certificates/nl-01/issuers/ca-1/rotate",
			call: func(api *Client) error {
				item, err := api.RotateAuthority(context.Background(), "nl-01", "ca-1", RotateAuthorityRequest{})
				if err != nil {
					return err
				}
				assert.Equal(t, "root", item.Type)
				assert.Equal(t, 1, item.PathLen)
				require.Len(t, item.Certificates, 1)
				assert.Equal(t, "previous", item.Certificates[0].Kind)
				return nil
			},
		},
		{
			name:   "cross-sign",
			method: http.MethodPost,
			path:   "/v1/certificates/nl-01/issuers/ca-1/cross-sign",
			call: func(api *Client) error {
				_, err := api.CrossSignAuthority(context.Background(), "nl-01", "ca-1", CrossSignAuthorityRequest{IssuerAuthorityIdentity: "ca-2"})
				return err
			},
		},
		{
			name:   "version pem",
			method: http.MethodGet,
			path:   "/v1/certificates/nl-01/issuers/ca-1/certificates/cac-1/pem",
			call: func(api *Client) error {
				body, err := api.AuthorityCertificatePEM(context.Background(), "nl-01", "ca-1", "cac-1")
				if err != nil {
					return err
				}
				assert.Equal(t, pemText, body.Certificate)
				assert.NotContains(t, body.Certificate, "PRIVATE KEY")
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tt.method, r.Method)
				assert.Equal(t, tt.path, r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				if tt.name == "version pem" {
					_, _ = w.Write([]byte(`{"certificate":"` + encoded + `"}`))
					return
				}
				_, _ = w.Write([]byte(`{"identity":"ca-1","name":"root","status":"active","type":"root","pathLen":1,"certificates":[{"identity":"cac-1","kind":"previous","current":false}]}`))
			}))
			t.Cleanup(server.Close)
			base, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
			require.NoError(t, err)
			api, err := New(base)
			require.NoError(t, err)
			require.NoError(t, tt.call(api))
		})
	}
}
