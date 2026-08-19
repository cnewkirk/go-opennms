package opennms

// Monitoring Systems REST API – /rest/monitoringSystems.

import "context"

// GetMonitoringSystem returns the main monitoring system information.
func (c *Client) GetMonitoringSystem(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "monitoringSystems/main", nil, false)
}
