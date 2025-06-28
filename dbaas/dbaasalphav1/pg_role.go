package dbaasalphav1

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

// PostgreSQL Role Operations

// CreatePgRole creates a new PostgreSQL role in a database cluster.
func (c *Client) CreatePgRole(ctx context.Context, dbClusterIdentity string, create CreatePgRoleRequest) error {
	if dbClusterIdentity == "" {
		return fmt.Errorf("database cluster identity is required")
	}
	if create.Name == "" {
		return fmt.Errorf("role name is required")
	}

	req := c.R().SetBody(create)
	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/roles", DbClusterEndpoint, dbClusterIdentity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

// UpdatePgRole updates an existing PostgreSQL role in a database cluster.
func (c *Client) UpdatePgRole(ctx context.Context, dbClusterIdentity string, roleName string, update UpdatePgRoleRequest) error {
	if dbClusterIdentity == "" {
		return fmt.Errorf("database cluster identity is required")
	}
	if roleName == "" {
		return fmt.Errorf("role name is required")
	}

	req := c.R().SetBody(update)
	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s/roles/%s", DbClusterEndpoint, dbClusterIdentity, roleName))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

// DeletePgRole deletes a PostgreSQL role from a database cluster.
func (c *Client) DeletePgRole(ctx context.Context, dbClusterIdentity string, roleName string) error {
	if dbClusterIdentity == "" {
		return fmt.Errorf("database cluster identity is required")
	}
	if roleName == "" {
		return fmt.Errorf("role name is required")
	}

	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s/roles/%s", DbClusterEndpoint, dbClusterIdentity, roleName))
	if err != nil {
		return err
	}
	return c.Check(resp)
}
