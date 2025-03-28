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
	routeTables := []RouteTable{}
	req := c.R().SetResult(&routeTables)

	resp, err := c.Do(ctx, req, client.GET, RouteTableEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return routeTables, err
	}
	return routeTables, nil
}

// GetRouteTable retrieves a specific RouteTable by its identity.
func (c *Client) GetRouteTable(ctx context.Context, identity string) (*RouteTable, error) {
	var routeTable *RouteTable
	req := c.R().SetResult(&routeTable)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", RouteTableEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return routeTable, err
	}
	return routeTable, nil
}

// CreateRouteTable creates a new RouteTable.
func (c *Client) CreateRouteTable(ctx context.Context, create CreateRouteTable) (*RouteTable, error) {
	var routeTable *RouteTable
	req := c.R().
		SetBody(create).SetResult(&routeTable)

	resp, err := c.Do(ctx, req, client.POST, RouteTableEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return routeTable, err
	}
	return routeTable, nil
}

// UpdateRouteTable updates an existing RouteTable.
func (c *Client) UpdateRouteTable(ctx context.Context, identity string, update UpdateRouteTable) (*RouteTable, error) {
	var routeTable *RouteTable
	req := c.R().
		SetBody(update).SetResult(&routeTable)

	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", RouteTableEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return routeTable, err
	}
	return routeTable, nil
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

// UpdateRouteTableRoutes updates the routes for a specific RouteTable.
func (c *Client) UpdateRouteTableRoutes(ctx context.Context, identity string, update UpdateRouteTableRoutes) (*RouteTable, error) {
	var routeTable *RouteTable
	req := c.R().
		SetBody(update).SetResult(&routeTable)

	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s/routes", RouteTableEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return routeTable, err
	}
	return routeTable, nil
}
