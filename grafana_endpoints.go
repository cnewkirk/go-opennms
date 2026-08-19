package opennms

// Grafana Endpoints REST API – /rest/endpoints/grafana.

import (
	"context"
	"fmt"
	"net/url"
)

// GetGrafanaEndpoints returns all configured Grafana endpoints, or
// nil when none are configured (204 No Content).
func (c *Client) GetGrafanaEndpoints(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "endpoints/grafana", nil, false)
}

// GetGrafanaEndpoint returns a Grafana endpoint by its numeric
// identifier.
func (c *Client) GetGrafanaEndpoint(ctx context.Context, endpointID int) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("endpoints/grafana/%d", endpointID), nil, false)
}

// GetGrafanaDashboards returns the dashboards of the Grafana endpoint
// with the given uid (unique identifier of the endpoint).
func (c *Client) GetGrafanaDashboards(ctx context.Context, uid string) ([]any, error) {
	return c.getList(ctx,
		fmt.Sprintf("endpoints/grafana/%s/dashboards", url.PathEscape(uid)),
		nil, false)
}

// GetGrafanaDashboard returns a dashboard from the Grafana endpoint
// with the given uid. dashboardID identifies the dashboard to fetch.
func (c *Client) GetGrafanaDashboard(ctx context.Context, uid, dashboardID string) (map[string]any, error) {
	return c.getObject(ctx,
		fmt.Sprintf("endpoints/grafana/%s/dashboards/%s",
			url.PathEscape(uid), url.PathEscape(dashboardID)),
		nil, false)
}

// CreateGrafanaEndpoint creates a new Grafana endpoint. The endpoint
// definition requires "uid", "url", and "apiKey".
func (c *Client) CreateGrafanaEndpoint(ctx context.Context, endpoint map[string]any) error {
	_, err := c.post(ctx, "endpoints/grafana", endpoint, nil, false)
	return err
}

// VerifyGrafanaEndpoint verifies connectivity of a Grafana endpoint
// definition ("url" and "apiKey" at minimum).
func (c *Client) VerifyGrafanaEndpoint(ctx context.Context, endpoint map[string]any) error {
	_, err := c.post(ctx, "endpoints/grafana/verify", endpoint, nil, false)
	return err
}

// UpdateGrafanaEndpoint updates an existing Grafana endpoint. The
// updated endpoint definition must include "id".
func (c *Client) UpdateGrafanaEndpoint(ctx context.Context, endpointID int, endpoint map[string]any) error {
	_, err := c.put(ctx, fmt.Sprintf("endpoints/grafana/%d", endpointID),
		endpoint, nil, false)
	return err
}

// DeleteGrafanaEndpoints deletes all Grafana endpoints.
func (c *Client) DeleteGrafanaEndpoints(ctx context.Context) error {
	_, err := c.del(ctx, "endpoints/grafana", nil, nil, false, "")
	return err
}

// DeleteGrafanaEndpoint deletes a Grafana endpoint by its numeric
// identifier.
func (c *Client) DeleteGrafanaEndpoint(ctx context.Context, endpointID int) error {
	_, err := c.del(ctx, fmt.Sprintf("endpoints/grafana/%d", endpointID),
		nil, nil, false, "")
	return err
}
