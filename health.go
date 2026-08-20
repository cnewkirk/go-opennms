package opennms

// Health REST API – /rest/health.

import (
	"context"
	"net/url"
)

// GetHealth returns the health status of the OpenNMS instance. Pass
// tag to filter health checks; "" returns all.
func (c *Client) GetHealth(ctx context.Context, tag string) (map[string]any, error) {
	var params url.Values
	if tag != "" {
		params = url.Values{"tag": {tag}}
	}
	return c.getObject(ctx, "health", params, false)
}

// GetHealthProbe returns a simple health probe (up/down) response.
func (c *Client) GetHealthProbe(ctx context.Context) (string, error) {
	return asString(c.get(ctx, "health/probe", nil, false))
}
