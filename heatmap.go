package opennms

// Heatmap REST API – /rest/heatmap (read-only).

import (
	"context"
	"net/url"
)

// Outage-based heatmap

// GetHeatmapOutagesCategories returns outage heatmap data grouped by
// category.
func (c *Client) GetHeatmapOutagesCategories(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "heatmap/outages/categories", nil, false)
}

// GetHeatmapOutagesForeignSources returns outage heatmap data
// grouped by foreign source.
func (c *Client) GetHeatmapOutagesForeignSources(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "heatmap/outages/foreignSources", nil, false)
}

// GetHeatmapOutagesMonitoredServices returns outage heatmap data
// grouped by monitored service.
func (c *Client) GetHeatmapOutagesMonitoredServices(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "heatmap/outages/monitoredServices", nil, false)
}

// GetHeatmapOutagesNodesByCategory returns outage heatmap node data
// for the category.
func (c *Client) GetHeatmapOutagesNodesByCategory(ctx context.Context, category string) (map[string]any, error) {
	return c.getObject(ctx,
		"heatmap/outages/nodesByCategory/"+url.PathEscape(category), nil, false)
}

// GetHeatmapOutagesNodesByForeignSource returns outage heatmap node
// data for the foreign source.
func (c *Client) GetHeatmapOutagesNodesByForeignSource(ctx context.Context, foreignSource string) (map[string]any, error) {
	return c.getObject(ctx,
		"heatmap/outages/nodesByForeignSource/"+url.PathEscape(foreignSource),
		nil, false)
}

// GetHeatmapOutagesNodesByService returns outage heatmap node data
// for the monitored service.
func (c *Client) GetHeatmapOutagesNodesByService(ctx context.Context, service string) (map[string]any, error) {
	return c.getObject(ctx,
		"heatmap/outages/nodesByMonitoredService/"+url.PathEscape(service),
		nil, false)
}

// Alarm-based heatmap

// GetHeatmapAlarmsCategories returns alarm heatmap data grouped by
// category.
func (c *Client) GetHeatmapAlarmsCategories(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "heatmap/alarms/categories", nil, false)
}

// GetHeatmapAlarmsForeignSources returns alarm heatmap data grouped
// by foreign source.
func (c *Client) GetHeatmapAlarmsForeignSources(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "heatmap/alarms/foreignSources", nil, false)
}

// GetHeatmapAlarmsMonitoredServices returns alarm heatmap data
// grouped by monitored service.
func (c *Client) GetHeatmapAlarmsMonitoredServices(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "heatmap/alarms/monitoredServices", nil, false)
}

// GetHeatmapAlarmsNodesByCategory returns alarm heatmap node data
// for the category.
func (c *Client) GetHeatmapAlarmsNodesByCategory(ctx context.Context, category string) (map[string]any, error) {
	return c.getObject(ctx,
		"heatmap/alarms/nodesByCategory/"+url.PathEscape(category), nil, false)
}

// GetHeatmapAlarmsNodesByForeignSource returns alarm heatmap node
// data for the foreign source.
func (c *Client) GetHeatmapAlarmsNodesByForeignSource(ctx context.Context, foreignSource string) (map[string]any, error) {
	return c.getObject(ctx,
		"heatmap/alarms/nodesByForeignSource/"+url.PathEscape(foreignSource),
		nil, false)
}

// GetHeatmapAlarmsNodesByService returns alarm heatmap node data for
// the monitored service.
func (c *Client) GetHeatmapAlarmsNodesByService(ctx context.Context, service string) (map[string]any, error) {
	return c.getObject(ctx,
		"heatmap/alarms/nodesByMonitoredService/"+url.PathEscape(service),
		nil, false)
}
