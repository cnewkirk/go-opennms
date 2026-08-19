package opennms

// Minions REST API – /rest/minions.

import (
	"context"
	"net/url"
)

// GetMinions lists all minions.
func (c *Client) GetMinions(ctx context.Context, opts *ListOptions) (map[string]any, error) {
	return c.getObject(ctx, "minions", listParams(opts), false)
}

// GetMinion returns a specific minion by ID.
func (c *Client) GetMinion(ctx context.Context, minionID string) (map[string]any, error) {
	return c.getObject(ctx, "minions/"+url.PathEscape(minionID), nil, false)
}

// GetMinionCount returns the number of minions.
func (c *Client) GetMinionCount(ctx context.Context) (int, error) {
	return c.getCount(ctx, "minions/count", false)
}
