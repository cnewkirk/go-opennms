package opennms

// Monitoring Locations REST API – /rest/monitoringLocations.

import (
	"context"
	"net/url"
)

// GetMonitoringLocations lists all monitoring locations.
func (c *Client) GetMonitoringLocations(ctx context.Context, opts *ListOptions) (map[string]any, error) {
	return c.getObject(ctx, "monitoringLocations", listParams(opts), false)
}

// GetMonitoringLocation returns the monitoring location with the
// given name.
func (c *Client) GetMonitoringLocation(ctx context.Context, name string) (map[string]any, error) {
	return c.getObject(ctx, "monitoringLocations/"+url.PathEscape(name), nil, false)
}

// GetDefaultMonitoringLocation returns the default monitoring
// location.
func (c *Client) GetDefaultMonitoringLocation(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "monitoringLocations/default", nil, false)
}

// GetMonitoringLocationCount returns the number of monitoring
// locations.
func (c *Client) GetMonitoringLocationCount(ctx context.Context) (int, error) {
	return c.getCount(ctx, "monitoringLocations/count", false)
}

// CreateMonitoringLocation creates a new monitoring location from the
// given definition (keys like "location-name", "monitoring-area").
func (c *Client) CreateMonitoringLocation(ctx context.Context, location map[string]any) error {
	_, err := c.post(ctx, "monitoringLocations", location, nil, false)
	return err
}

// UpdateMonitoringLocation updates a monitoring location with the
// given form-encoded key/value pairs.
func (c *Client) UpdateMonitoringLocation(ctx context.Context, name string, data map[string]string) error {
	form := url.Values{}
	for k, v := range data {
		form.Set(k, v)
	}
	_, err := c.putForm(ctx, "monitoringLocations/"+url.PathEscape(name), form, nil, false)
	return err
}

// DeleteMonitoringLocation deletes a monitoring location.
func (c *Client) DeleteMonitoringLocation(ctx context.Context, name string) error {
	_, err := c.del(ctx, "monitoringLocations/"+url.PathEscape(name), nil, nil, false, "")
	return err
}
