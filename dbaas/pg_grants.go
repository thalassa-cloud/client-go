package dbaas

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

// CreatePgGrant creates a new PostgreSQL grant for a role on a database in a database cluster.
func (c *Client) CreatePgGrant(ctx context.Context, dbClusterIdentity string, create CreatePgGrantRequest) error {
	if dbClusterIdentity == "" {
		return fmt.Errorf("database cluster identity is required")
	}
	if create.Name == "" {
		return fmt.Errorf("grant name is required")
	}
	if create.RoleName == "" {
		return fmt.Errorf("role name is required")
	}
	if create.DatabaseName == "" {
		return fmt.Errorf("database name is required")
	}

	req := c.R().SetBody(create)
	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/postgres-grants", DbClusterEndpoint, dbClusterIdentity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

// UpdatePgGrant updates an existing PostgreSQL grant for a role on a database in a database cluster.
func (c *Client) UpdatePgGrant(ctx context.Context, dbClusterIdentity string, grantIdentity string, update UpdatePgGrantRequest) error {
	if dbClusterIdentity == "" {
		return fmt.Errorf("database cluster identity is required")
	}
	if grantIdentity == "" {
		return fmt.Errorf("postgres grant identity is required")
	}

	req := c.R().SetBody(update)
	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s/postgres-grants/%s", DbClusterEndpoint, dbClusterIdentity, grantIdentity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

// DeletePgGrant deletes a PostgreSQL grant from a database cluster.
func (c *Client) DeletePgGrant(ctx context.Context, dbClusterIdentity string, grantIdentity string) error {
	if dbClusterIdentity == "" {
		return fmt.Errorf("database cluster identity is required")
	}
	if grantIdentity == "" {
		return fmt.Errorf("postgres grant identity is required")
	}

	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s/postgres-grants/%s", DbClusterEndpoint, dbClusterIdentity, grantIdentity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}
