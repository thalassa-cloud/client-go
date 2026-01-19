package dbaas

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

// PostgreSQL Backup Schedule Operations

// ListPgBackupSchedules lists all PostgreSQL backup schedules for a database cluster.
func (c *Client) ListPgBackupSchedules(ctx context.Context, dbClusterIdentity string) ([]DbClusterBackupSchedule, error) {
	if dbClusterIdentity == "" {
		return nil, fmt.Errorf("database cluster identity is required")
	}

	backupSchedules := []DbClusterBackupSchedule{}
	req := c.R().SetResult(&backupSchedules)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/backup-schedules", DbClusterEndpoint, dbClusterIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return backupSchedules, nil
}

// CreatePgBackupSchedule creates a new PostgreSQL backup schedule for a database cluster.
func (c *Client) CreatePgBackupSchedule(ctx context.Context, dbClusterIdentity string, create CreatePgBackupScheduleRequest) (*DbClusterBackupSchedule, error) {
	if dbClusterIdentity == "" {
		return nil, fmt.Errorf("database cluster identity is required")
	}
	if create.Name == "" {
		return nil, fmt.Errorf("backup schedule name is required")
	}
	if create.Schedule == "" {
		return nil, fmt.Errorf("backup schedule is required")
	}
	if create.RetentionPolicy == "" {
		return nil, fmt.Errorf("retention policy is required")
	}

	var backupSchedule *DbClusterBackupSchedule
	req := c.R().SetBody(create).SetResult(&backupSchedule)
	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/backup-schedules", DbClusterEndpoint, dbClusterIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return backupSchedule, err
	}
	return backupSchedule, nil
}

// UpdatePgBackupSchedule updates an existing PostgreSQL backup schedule for a database cluster.
func (c *Client) UpdatePgBackupSchedule(ctx context.Context, dbClusterIdentity string, backupScheduleIdentity string, update UpdatePgBackupScheduleRequest) (*DbClusterBackupSchedule, error) {
	if dbClusterIdentity == "" {
		return nil, fmt.Errorf("database cluster identity is required")
	}
	if backupScheduleIdentity == "" {
		return nil, fmt.Errorf("backup schedule identity is required")
	}
	if update.Name == "" {
		return nil, fmt.Errorf("backup schedule name is required")
	}
	if update.Schedule == "" {
		return nil, fmt.Errorf("backup schedule is required")
	}
	if update.RetentionPolicy == "" {
		return nil, fmt.Errorf("retention policy is required")
	}

	var backupSchedule *DbClusterBackupSchedule
	req := c.R().SetBody(update).SetResult(&backupSchedule)
	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s/backup-schedules/%s", DbClusterEndpoint, dbClusterIdentity, backupScheduleIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return backupSchedule, err
	}
	return backupSchedule, nil
}

// GetPgBackupSchedule retrieves a specific PostgreSQL backup schedule for a database cluster.
func (c *Client) GetPgBackupSchedule(ctx context.Context, dbClusterIdentity string, backupScheduleIdentity string) (*DbClusterBackupSchedule, error) {
	if dbClusterIdentity == "" {
		return nil, fmt.Errorf("database cluster identity is required")
	}
	if backupScheduleIdentity == "" {
		return nil, fmt.Errorf("backup schedule identity is required")
	}

	var backupSchedule *DbClusterBackupSchedule
	req := c.R().SetResult(&backupSchedule)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/backup-schedules/%s", DbClusterEndpoint, dbClusterIdentity, backupScheduleIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return backupSchedule, err
	}
	return backupSchedule, nil
}

// DeletePgBackupSchedule deletes a PostgreSQL backup schedule from a database cluster.
func (c *Client) DeletePgBackupSchedule(ctx context.Context, dbClusterIdentity string, backupScheduleIdentity string) error {
	if dbClusterIdentity == "" {
		return fmt.Errorf("database cluster identity is required")
	}
	if backupScheduleIdentity == "" {
		return fmt.Errorf("backup schedule identity is required")
	}

	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s/backup-schedules/%s", DbClusterEndpoint, dbClusterIdentity, backupScheduleIdentity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

// ListPgBackupSchedulesForOrganisation lists all PostgreSQL backup schedules for the organisation.
func (c *Client) ListPgBackupSchedulesForOrganisation(ctx context.Context) ([]DbClusterBackupSchedule, error) {
	backupSchedules := []DbClusterBackupSchedule{}
	req := c.R().SetResult(&backupSchedules)
	resp, err := c.Do(ctx, req, client.GET, "/v1/dbaas/backup-schedules")
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return backupSchedules, nil
}
