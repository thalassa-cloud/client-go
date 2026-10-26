package certificates

import "github.com/thalassa-cloud/client-go/pkg/client"

// Client calls the Certificate Manager API.
type Client struct {
	client.Client
}

// New creates a Certificate Manager client.
func New(c client.Client, opts ...client.Option) (*Client, error) {
	c.WithOptions(opts...)
	return &Client{c}, nil
}
