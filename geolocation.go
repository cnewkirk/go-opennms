package opennms

// Geolocation REST API v2 – /api/v2/geolocation.

import "context"

// GetGeolocationConfig returns the map/tile-server configuration: a
// map with "tileServerName", "tileServerUrl", and "options" keys.
func (c *Client) GetGeolocationConfig(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "geolocation/config", nil, true)
}

// QueryGeolocations queries nodes with geographic locations and
// their status.
//
// strategy is the status computation strategy — "Alarms" or
// "Outages"; "" defaults to "Alarms". severityFilter only returns
// nodes at or above that severity (e.g. "Major"); "" applies no
// filter. includeAcknowledgedAlarms controls whether acknowledged
// alarms count toward node status; nil leaves the server default.
//
// Returns a list of geolocation info maps, or nil for an empty
// result set (204 No Content).
func (c *Client) QueryGeolocations(ctx context.Context, strategy, severityFilter string, includeAcknowledgedAlarms *bool) ([]any, error) {
	if strategy == "" {
		strategy = "Alarms"
	}
	body := map[string]any{"strategy": strategy}
	if severityFilter != "" {
		body["severityFilter"] = severityFilter
	}
	if includeAcknowledgedAlarms != nil {
		body["includeAcknowledgedAlarms"] = *includeAcknowledgedAlarms
	}
	return asList(c.post(ctx, "geolocation", body, nil, true))
}
