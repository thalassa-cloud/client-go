package iaas

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	RouteTableEndpoint = "/v1/route-tables"
)

// ListRouteTables lists all RouteTables for a given organisation.
func (c *Client) ListRouteTables(ctx context.Context) ([]RouteTable, error) {
	subnets := []RouteTable{}
	req := c.R().SetResult(&subnets)

	resp, err := c.Do(ctx, req, client.GET, RouteTableEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return subnets, err
	}
	return subnets, nil
}

// GetRouteTable retrieves a specific RouteTable by its identity.
func (c *Client) GetRouteTable(ctx context.Context, identity string) (*RouteTable, error) {
	var subnet *RouteTable
	req := c.R().SetResult(&subnet)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", RouteTableEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return subnet, err
	}
	return subnet, nil
}

// CreateRouteTable creates a new RouteTable.
func (c *Client) CreateRouteTable(ctx context.Context, create CreateRouteTable) (*RouteTable, error) {
	var subnet *RouteTable
	req := c.R().
		SetBody(create).SetResult(&subnet)

	resp, err := c.Do(ctx, req, client.POST, RouteTableEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return subnet, err
	}
	return subnet, nil
}

// UpdateRouteTable updates an existing RouteTable.
func (c *Client) UpdateRouteTable(ctx context.Context, identity string, update UpdateRouteTable) (*RouteTable, error) {
	var subnet *RouteTable
	req := c.R().
		SetBody(update).SetResult(&subnet)

	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", RouteTableEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return subnet, err
	}
	return subnet, nil
}

// DeleteRouteTable deletes a specific RouteTable by its identity.
func (c *Client) DeleteRouteTable(ctx context.Context, identity string) error {
	req := c.R()

	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", RouteTableEndpoint, identity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}
