package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

func setupTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/success":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		case "/not-found":
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}
	}))
}

func TestNewClient(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	tests := []struct {
		name        string
		options     []Option
		expectError bool
		errorMsg    string
	}{
		{
			name:        "missing base URL",
			options:     []Option{},
			expectError: true,
			errorMsg:    "base URL is required",
		},
		{
			name: "with base URL only",
			options: []Option{
				WithBaseURL(server.URL),
			},
			expectError: false,
		},
		{
			name: "with base URL and organization",
			options: []Option{
				WithBaseURL(server.URL),
				WithOrganisation("test-org"),
			},
			expectError: false,
		},
		{
			name: "with base URL and project",
			options: []Option{
				WithBaseURL(server.URL),
				WithProject("test-project"),
			},
			expectError: false,
		},
		{
			name: "with base URL and timeout",
			options: []Option{
				WithBaseURL(server.URL),
				WithTimeout(5 * time.Second),
			},
			expectError: false,
		},
		{
			name: "with base URL and retries",
			options: []Option{
				WithBaseURL(server.URL),
				WithRetries(3, 1*time.Second, 5*time.Second),
			},
			expectError: false,
		},
		{
			name: "with base URL and rate limit",
			options: []Option{
				WithBaseURL(server.URL),
				WithRateLimit(10, 5),
			},
			expectError: false,
		},
		{
			name: "with base URL and circuit breaker",
			options: []Option{
				WithBaseURL(server.URL),
				WithCircuitBreaker("test-breaker", gobreaker.Settings{
					Timeout: 5 * time.Second,
				}),
			},
			expectError: false,
		},
		{
			name: "with base URL and user agent",
			options: []Option{
				WithBaseURL(server.URL),
				WithUserAgent("test-user-agent"),
			},
			expectError: false,
		},
		{
			name: "with base URL and auth none",
			options: []Option{
				WithBaseURL(server.URL),
				WithAuthNone(),
			},
			expectError: false,
		},
		{
			name: "with base URL and auth custom",
			options: []Option{
				WithBaseURL(server.URL),
				WithAuthCustom(),
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.options...)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, client)
		})
	}
}

func TestClientWithAuthOptions(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	tests := []struct {
		name        string
		options     []Option
		expectError bool
		errorMsg    string
	}{
		{
			name: "with personal access token",
			options: []Option{
				WithBaseURL(server.URL),
				WithAuthPersonalToken("test-token"),
			},
			expectError: false,
		},
		{
			name: "with empty personal access token",
			options: []Option{
				WithBaseURL(server.URL),
				WithAuthPersonalToken(""),
			},
			expectError: true,
			errorMsg:    "personal access token cannot be empty",
		},
		{
			name: "with basic auth",
			options: []Option{
				WithBaseURL(server.URL),
				WithAuthBasic("username", "password"),
			},
			expectError: false,
		},
		{
			name: "with basic auth missing username",
			options: []Option{
				WithBaseURL(server.URL),
				WithAuthBasic("", "password"),
			},
			expectError: true,
			errorMsg:    "basic auth requires username/password",
		},
		{
			name: "with basic auth missing password",
			options: []Option{
				WithBaseURL(server.URL),
				WithAuthBasic("username", ""),
			},
			expectError: true,
			errorMsg:    "basic auth requires username/password",
		},
		{
			name: "with OIDC auth",
			options: []Option{
				WithBaseURL(server.URL),
				WithAuthOIDC("client-id", "client-secret", "http://example.com/token", "scope1", "scope2"),
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.options...)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, client)
		})
	}
}

func TestClientWithOptions(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	baseClient, err := NewClient(WithBaseURL(server.URL))
	require.NoError(t, err)
	require.NotNil(t, baseClient)

	// Test WithOptions
	client := baseClient.WithOptions(
		WithOrganisation("test-org"),
		WithProject("test-project"),
	)
	require.NotNil(t, client)

	// Verify the options were applied
	assert.Equal(t, "test-org", client.GetOrganisationIdentity())
	assert.Equal(t, server.URL, client.GetBaseURL())
}

func TestClientGetAuthToken(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Test with personal access token
	client, err := NewClient(
		WithBaseURL(server.URL),
		WithAuthPersonalToken("test-token"),
	)
	require.NoError(t, err)
	assert.Equal(t, "test-token", client.GetAuthToken())

	// Test with no auth
	client, err = NewClient(
		WithBaseURL(server.URL),
		WithAuthNone(),
	)
	require.NoError(t, err)
	assert.Empty(t, client.GetAuthToken())
}

func TestClientDo(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a client with base URL
	client, err := NewClient(
		WithBaseURL(server.URL),
		WithAuthPersonalToken("test-token"),
	)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Create a request
	req := client.R()
	require.NotNil(t, req)

	// Test successful request
	ctx := context.Background()
	resp, err := client.Do(ctx, req, http.MethodGet, "/success")
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, "success", string(resp.Body()))

	// Test not found request
	resp, err = client.Do(ctx, req, http.MethodGet, "/not-found")
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode())
	assert.Equal(t, "not found", string(resp.Body()))
}

// mockResponse is a mock implementation of resty.Response for testing
type mockResponse struct {
	statusCode int
	body       string
}

func (m *mockResponse) StatusCode() int {
	return m.statusCode
}

func (m *mockResponse) String() string {
	return m.body
}

func (m *mockResponse) IsError() bool {
	return m.statusCode >= 400
}

func TestClientCheck(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a client
	client, err := NewClient(
		WithBaseURL(server.URL),
		WithAuthPersonalToken("test-token"),
	)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Test successful request
	resp, err := client.R().Get("/success")
	require.NoError(t, err)
	err = client.Check(resp)
	assert.NoError(t, err)

	// Test not found request
	resp, err = client.R().Get("/not-found")
	require.NoError(t, err)
	err = client.Check(resp)
	assert.Error(t, err)
	assert.True(t, IsNotFound(err))
}

func TestClientDialWebsocket(t *testing.T) {
	// Create a client
	client, err := NewClient(
		WithBaseURL("http://localhost"),
		WithAuthPersonalToken("test-token"),
	)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Test DialWebsocket method
	ctx := context.Background()
	conn, err := client.DialWebsocket(ctx, "ws://localhost/ws")
	// This will fail because we're not actually running a websocket server
	assert.Error(t, err)
	assert.Nil(t, conn)
}

func TestClientSetOrganisation(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a client
	client, err := NewClient(
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Test SetOrganisation method
	client.SetOrganisation("test-org")
	assert.Equal(t, "test-org", client.GetOrganisationIdentity())
}

func TestClientWithMiddleware(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a client with middleware
	middlewareCalled := false
	client, err := NewClient(
		WithBaseURL(server.URL),
		WithMiddleware(func(c *resty.Client, req *resty.Request) error {
			middlewareCalled = true
			return nil
		}),
	)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Make a request to trigger the middleware
	resp, err := client.R().Get("/success")
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify the middleware was called
	assert.True(t, middlewareCalled)
}

func TestClientWithRateLimit(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a client with rate limit
	client, err := NewClient(
		WithBaseURL(server.URL),
		WithRateLimit(10, 5),
	)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Access the limiter through reflection to verify it's set
	clientStruct := client.(*thalassaCloudClient)
	assert.NotNil(t, clientStruct.limiter)
	assert.Equal(t, rate.Limit(10), clientStruct.limiter.Limit())
	assert.Equal(t, 5, clientStruct.limiter.Burst())
}

func TestClientWithCircuitBreaker(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a client with circuit breaker
	client, err := NewClient(
		WithBaseURL(server.URL),
		WithCircuitBreaker("test-breaker", gobreaker.Settings{
			Timeout: 5 * time.Second,
		}),
	)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Access the breaker through reflection to verify it's set
	clientStruct := client.(*thalassaCloudClient)
	assert.NotNil(t, clientStruct.breaker)
	assert.Equal(t, "test-breaker", clientStruct.breaker.Name())
}

func TestClientWithInsecure(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a client with insecure option
	client, err := NewClient(
		WithBaseURL(server.URL),
		WithInsecure(),
	)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Access the client through reflection to verify TLS config
	clientStruct := client.(*thalassaCloudClient)
	transport := clientStruct.resty.GetClient().Transport.(*http.Transport)
	assert.True(t, transport.TLSClientConfig.InsecureSkipVerify)
}
