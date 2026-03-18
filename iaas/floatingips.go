package iaas

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const FloatingIPEndpoint = "/v1/floating-ips"

// ListFloatingIPs lists floating IPs for the organisation (auth context).
func (c *Client) ListFloatingIPs(ctx context.Context, listRequest *ListFloatingIPsRequest) ([]FloatingIP, error) {
	var out []FloatingIP
	req := c.R().SetResult(&out)
	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}
	resp, err := c.Do(ctx, req, client.GET, FloatingIPEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return out, err
	}
	return out, nil
}

// CreateFloatingIP creates a floating IP (201).
func (c *Client) CreateFloatingIP(ctx context.Context, create CreateFloatingIpRequest) (*FloatingIP, error) {
	if create.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if create.Region == "" {
		return nil, fmt.Errorf("region is required")
	}
	var fip *FloatingIP
	req := c.R().SetBody(create).SetResult(&fip)
	resp, err := c.Do(ctx, req, client.POST, FloatingIPEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return fip, err
	}
	return fip, nil
}

// GetFloatingIP returns a floating IP by identity.
func (c *Client) GetFloatingIP(ctx context.Context, identity string) (*FloatingIP, error) {
	if identity == "" {
		return nil, fmt.Errorf("identity is required")
	}
	var fip *FloatingIP
	req := c.R().SetResult(&fip)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", FloatingIPEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return fip, err
	}
	return fip, nil
}

// UpdateFloatingIP updates name, description, labels, and annotations.
func (c *Client) UpdateFloatingIP(ctx context.Context, identity string, update UpdateFloatingIpRequest) (*FloatingIP, error) {
	if identity == "" {
		return nil, fmt.Errorf("identity is required")
	}
	if update.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	var fip *FloatingIP
	req := c.R().SetBody(update).SetResult(&fip)
	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", FloatingIPEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return fip, err
	}
	return fip, nil
}

// DeleteFloatingIP deletes a floating IP. If attached, the API disassociates first (204).
func (c *Client) DeleteFloatingIP(ctx context.Context, identity string) error {
	if identity == "" {
		return fmt.Errorf("identity is required")
	}
	resp, err := c.Do(ctx, c.R(), client.DELETE, fmt.Sprintf("%s/%s", FloatingIPEndpoint, identity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

// AssociateFloatingIP attaches the floating IP to a load balancer or NAT gateway.
func (c *Client) AssociateFloatingIP(ctx context.Context, identity string, body AssociateFloatingIpRequest) (*FloatingIP, error) {
	if identity == "" {
		return nil, fmt.Errorf("identity is required")
	}
	var fip *FloatingIP
	req := c.R().SetBody(body).SetResult(&fip)
	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/associate", FloatingIPEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return fip, err
	}
	return fip, nil
}

// DisassociateFloatingIP detaches the floating IP from its current target.
func (c *Client) DisassociateFloatingIP(ctx context.Context, identity string) (*FloatingIP, error) {
	if identity == "" {
		return nil, fmt.Errorf("identity is required")
	}
	var fip *FloatingIP
	req := c.R().SetResult(&fip)
	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/disassociate", FloatingIPEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return fip, err
	}
	return fip, nil
}

// ListFloatingIPsRequest holds query filters for ListFloatingIPs.
// Use filters.FilterKeyValue with filters.FilterFloatingIp, FilterName, FilterIdentity,
// FilterSlug, FilterStatus, FilterRegion, or LabelFilter as for other IaaS list APIs.
type ListFloatingIPsRequest struct {
	Filters []filters.Filter
}
