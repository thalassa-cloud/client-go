package dns

import (
	"context"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

// GetDnssec returns DNSSEC status for a zone.
func (c *Client) GetDnssec(ctx context.Context, zoneIdentity string) (*DnsZoneDnssecStatus, error) {
	var status DnsZoneDnssecStatus
	r := c.R().SetResult(&status)
	resp, err := c.Do(ctx, r, client.GET, zonePath(zoneIdentity, "dnssec"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &status, err
	}
	return &status, nil
}

// UpsertDnssec enables or updates DNSSEC signing for a zone.
func (c *Client) UpsertDnssec(ctx context.Context, zoneIdentity string, upsert UpsertDnssecRequest) (*DnsZoneDnssecStatus, error) {
	var status DnsZoneDnssecStatus
	r := c.R().SetBody(upsert).SetResult(&status)
	resp, err := c.Do(ctx, r, client.PUT, zonePath(zoneIdentity, "dnssec"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &status, err
	}
	return &status, nil
}

// DeleteDnssec disables DNSSEC signing for a zone.
func (c *Client) DeleteDnssec(ctx context.Context, zoneIdentity string) error {
	resp, err := c.Do(ctx, c.R(), client.DELETE, zonePath(zoneIdentity, "dnssec"))
	if err != nil {
		return err
	}
	return c.Check(resp)
}
