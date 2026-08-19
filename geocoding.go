package opennms

// Geocoding REST API v2 – /api/v2/geocoding.

import (
	"context"
	"net/url"
)

// GetGeocodingConfig returns the geocoder service manager
// configuration: a map with the "activeGeocoderId" key (nil value
// when no geocoder is active).
func (c *Client) GetGeocodingConfig(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "geocoding/config", nil, true)
}

// GetGeocoders returns all registered geocoder services with their
// configuration, or nil when none are registered (204 No Content).
func (c *Client) GetGeocoders(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "geocoding/geocoders", nil, true)
}

// SetActiveGeocoder activates the geocoder service with the given
// ID, e.g. "nominatim" or "google".
func (c *Client) SetActiveGeocoder(ctx context.Context, geocoderID string) error {
	_, err := c.post(ctx, "geocoding/config",
		map[string]any{"activeGeocoderId": geocoderID}, nil, true)
	return err
}

// ConfigureGeocoder updates the configuration of a geocoder service.
// config holds provider-specific settings as string key/value pairs,
// e.g. {"apiKey": "..."}.
func (c *Client) ConfigureGeocoder(ctx context.Context, geocoderID string, config map[string]string) error {
	_, err := c.post(ctx, "geocoding/geocoders/"+url.PathEscape(geocoderID),
		map[string]any{"config": config}, nil, true)
	return err
}

// ResetGeocodingConfig resets the geocoder service manager to its
// defaults.
func (c *Client) ResetGeocodingConfig(ctx context.Context) error {
	_, err := c.del(ctx, "geocoding/config", nil, nil, true, "")
	return err
}
