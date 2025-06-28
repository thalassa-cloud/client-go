package dbaasalphav1

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

// PostgreSQL Database Operations

// CreatePgDatabase creates a new PostgreSQL database in a database cluster.
func (c *Client) CreatePgDatabase(ctx context.Context, dbClusterIdentity string, create CreatePgDatabaseRequest) error {
	if dbClusterIdentity == "" {
		return fmt.Errorf("database cluster identity is required")
	}
	if create.Name == "" {
		return fmt.Errorf("database name is required")
	}
	if create.Owner == "" {
		return fmt.Errorf("database owner is required")
	}

	req := c.R().SetBody(create)
	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/databases", DbClusterEndpoint, dbClusterIdentity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

// UpdatePgDatabase updates an existing PostgreSQL database in a database cluster.
func (c *Client) UpdatePgDatabase(ctx context.Context, dbClusterIdentity string, databaseName string, update UpdatePgDatabaseRequest) error {
	if dbClusterIdentity == "" {
		return fmt.Errorf("database cluster identity is required")
	}
	if databaseName == "" {
		return fmt.Errorf("database name is required")
	}

	req := c.R().SetBody(update)
	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s/databases/%s", DbClusterEndpoint, dbClusterIdentity, databaseName))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

// DeletePgDatabase deletes a PostgreSQL database from a database cluster.
func (c *Client) DeletePgDatabase(ctx context.Context, dbClusterIdentity string, databaseName string) error {
	if dbClusterIdentity == "" {
		return fmt.Errorf("database cluster identity is required")
	}
	if databaseName == "" {
		return fmt.Errorf("database name is required")
	}

	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s/databases/%s", DbClusterEndpoint, dbClusterIdentity, databaseName))
	if err != nil {
		return err
	}
	return c.Check(resp)
}
