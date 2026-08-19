package opennms

// Outages REST API – /rest/outages (read-only).

import (
	"context"
	"fmt"
)

// GetOutages lists outages (read-only).
//
// filters are additional Hibernate query filters passed directly as
// query parameters (e.g. "node.label": "myrouter").
func (c *Client) GetOutages(ctx context.Context, opts *ListOptions, filters map[string]string) (map[string]any, error) {
	params := mergeFilters(listParams(opts), filters)
	return c.getObject(ctx, "outages", params, false)
}

// GetOutage returns a single outage by ID.
func (c *Client) GetOutage(ctx context.Context, outageID int) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("outages/%d", outageID), nil, false)
}

// GetOutageCount returns the total number of outages.
func (c *Client) GetOutageCount(ctx context.Context) (int, error) {
	return c.getCount(ctx, "outages/count", false)
}

// GetNodeOutages returns all outages for the given node.
func (c *Client) GetNodeOutages(ctx context.Context, nodeID int) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("outages/forNode/%d", nodeID), nil, false)
}
