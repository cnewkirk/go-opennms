package opennms

// Provisiond REST API v2 – /api/v2/provisiond.

import (
	"context"
	"net/url"
)

// GetProvisiondStatus returns the current status of the Provisiond
// daemon.
func (c *Client) GetProvisiondStatus(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "provisiond/status", nil, true)
}

// GetProvisiondJobStatus returns the status of a specific
// provisioning job.
func (c *Client) GetProvisiondJobStatus(ctx context.Context, jobID string) (map[string]any, error) {
	return c.getObject(ctx, "provisiond/status/"+url.PathEscape(jobID), nil, true)
}
