package opennms

// Availability REST API – /rest/availability.

import (
	"context"
	"fmt"
	"net/url"
)

// GetAvailability returns the availability summary for all categories.
func (c *Client) GetAvailability(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "availability", nil, false)
}

// GetAvailabilityCategory returns the availability summary for the
// given surveillance category.
func (c *Client) GetAvailabilityCategory(ctx context.Context, category string) (map[string]any, error) {
	return c.getObject(ctx,
		fmt.Sprintf("availability/categories/%s", url.PathEscape(category)),
		nil, false)
}

// GetAvailabilityCategoryNodes returns per-node availability for the
// given surveillance category.
func (c *Client) GetAvailabilityCategoryNodes(ctx context.Context, category string) (map[string]any, error) {
	return c.getObject(ctx,
		fmt.Sprintf("availability/categories/%s/nodes", url.PathEscape(category)),
		nil, false)
}

// GetAvailabilityCategoryNode returns availability for a specific
// node (by database ID) within the given surveillance category.
func (c *Client) GetAvailabilityCategoryNode(ctx context.Context, category string, nodeID int) (map[string]any, error) {
	return c.getObject(ctx,
		fmt.Sprintf("availability/categories/%s/nodes/%d", url.PathEscape(category), nodeID),
		nil, false)
}

// GetAvailabilityNode returns the availability summary for a specific
// node by database ID.
func (c *Client) GetAvailabilityNode(ctx context.Context, nodeID int) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("availability/nodes/%d", nodeID), nil, false)
}
