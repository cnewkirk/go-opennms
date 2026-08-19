package opennms

// Applications REST API v2 – /api/v2/applications.

import (
	"context"
	"fmt"
)

// GetApplications lists all applications.
func (c *Client) GetApplications(ctx context.Context, opts *ListOptions) (map[string]any, error) {
	return c.getObject(ctx, "applications", listParams(opts), true)
}

// GetApplication returns a specific application by ID.
func (c *Client) GetApplication(ctx context.Context, appID int) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("applications/%d", appID), nil, true)
}

// CreateApplication creates a new application.
//
// app is the application definition with keys such as "name" and
// "monitoredServices".
func (c *Client) CreateApplication(ctx context.Context, app map[string]any) (map[string]any, error) {
	return asObject(c.post(ctx, "applications", app, nil, true))
}

// DeleteApplication deletes an application.
func (c *Client) DeleteApplication(ctx context.Context, appID int) error {
	_, err := c.del(ctx, fmt.Sprintf("applications/%d", appID), nil, nil, true, "")
	return err
}
