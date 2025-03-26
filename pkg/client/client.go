package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/sony/gobreaker"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/time/rate"
)

// Option is a function that modifies the Client.
type Option func(*thalassaCloudClient) error

var (
	ErrMissingBaseURL          = errors.New("missing base URL; use WithBaseURL(...)")
	ErrMissingOIDCConfig       = errors.New("OIDC configuration is missing")
	ErrEmptyPersonalToken      = errors.New("personal access token cannot be empty")
	ErrMissingBasicCredentials = errors.New("basic auth requires username/password")
	ErrUnsupportedHTTPMethod   = errors.New("unsupported HTTP method")
	ErrNotFound                = errors.New("not found")
)

type AuthenticationType int

const (
	AuthNone AuthenticationType = iota
	AuthOIDC
	AuthPersonalAccessToken
	AuthBasic
	AuthCustom
)

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

type Client interface {
	Do(ctx context.Context, req *resty.Request, method httpMethod, url string) (*resty.Response, error)
	Check(resp *resty.Response) error

	R() *resty.Request

	WithOptions(opts ...Option) Client
}

// NewClient applies all options, configures authentication, and returns the client.
func NewClient(opts ...Option) (Client, error) {
	c := &thalassaCloudClient{resty: resty.New()}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	if c.resty.BaseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}

	// Configure built-in authentication once we have all fields set.
	if err := c.configureAuth(); err != nil {
		return nil, err
	}

	return c, nil
}

type thalassaCloudClient struct {
	// Underlying resty client.
	resty *resty.Client

	organisationIdentity *string
	projectIdentity      *string

	// Authentication fields.
	authType AuthenticationType

	// OIDC (client credentials).
	oidcConfig *clientcredentials.Config
	oidcToken  *oauth2.Token // cached token

	// Personal Access Token.
	personalToken string

	// Basic Auth.
	basicUsername string
	basicPassword string

	// Rate limiting.
	limiter *rate.Limiter

	// Optional circuit breaker
	breaker *gobreaker.CircuitBreaker
}

func (c *thalassaCloudClient) WithOptions(opts ...Option) Client {
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *thalassaCloudClient) R() *resty.Request {
	return c.resty.R()
}
