package opennms

// Server Info REST API – /rest/info.

import "context"

// GetInfo returns OpenNMS server version, package info, and running
// services.
//
// Example response keys: "displayVersion", "version", "packageName",
// "packageDescription", "ticketerConfig", "datetimeformatConfig",
// "services".
func (c *Client) GetInfo(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "info", nil, false)
}
