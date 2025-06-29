package dbaasalphav1

import (
	"context"
	"errors"
	"strings"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	EngineVersionEndpoint = "/v1/dbaas/engine-versions"
)

// ListEngineVersions lists all engine versions for a given organisation.
func (c *Client) ListEngineVersions(ctx context.Context, engine DbClusterDatabaseEngine, listRequest *ListEngineVersionsRequest) ([]DbClusterEngineVersion, error) {
	if strings.TrimSpace(string(engine)) == "" {
		return nil, errors.New("engine is required")
	}

	engineVersions := []DbClusterEngineVersion{}
	req := c.R().SetResult(&engineVersions).
		SetQueryParam("engine", string(engine))

	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, EngineVersionEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return engineVersions, err
	}

	return engineVersions, nil
}

// ListEngineVersionsRequest is the request for the ListEngineVersions function.
type ListEngineVersionsRequest struct {
	Filters []filters.Filter
}
