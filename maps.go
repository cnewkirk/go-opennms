package opennms

// Maps REST API – /rest/maps.
//
// The maps REST API was removed upstream in OpenNMS Horizon 16
// (2015). These methods are retained for backwards compatibility
// with pre-16 servers only and do not function on OpenNMS Horizon
// 16+ or any Meridian release (the server responds 404).

import (
	"context"
	"fmt"
)

// GetMaps lists all maps.
//
// Backwards compatibility only — the maps API was removed in OpenNMS
// Horizon 16 and this call 404s on newer servers.
func (c *Client) GetMaps(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "maps", nil, false)
}

// GetMap returns a specific map by ID.
//
// Backwards compatibility only — the maps API was removed in OpenNMS
// Horizon 16 and this call 404s on newer servers.
func (c *Client) GetMap(ctx context.Context, mapID int) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("maps/%d", mapID), nil, false)
}

// GetMapElements returns nodes, links, and elements for the map.
//
// Backwards compatibility only — the maps API was removed in OpenNMS
// Horizon 16 and this call 404s on newer servers.
func (c *Client) GetMapElements(ctx context.Context, mapID int) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("maps/%d/mapElements", mapID),
		nil, false)
}

// CreateMap adds a new map from the given map definition.
//
// Backwards compatibility only — the maps API was removed in OpenNMS
// Horizon 16 and this call 404s on newer servers.
func (c *Client) CreateMap(ctx context.Context, mapData map[string]any) (map[string]any, error) {
	return asObject(c.post(ctx, "maps", mapData, nil, false))
}

// UpdateMap updates map properties. mapData is the map of fields to
// change.
//
// Backwards compatibility only — the maps API was removed in OpenNMS
// Horizon 16 and this call 404s on newer servers.
func (c *Client) UpdateMap(ctx context.Context, mapID int, mapData map[string]any) error {
	_, err := c.put(ctx, fmt.Sprintf("maps/%d", mapID), mapData,
		nil, false)
	return err
}

// DeleteMap deletes a map.
//
// Backwards compatibility only — the maps API was removed in OpenNMS
// Horizon 16 and this call 404s on newer servers.
func (c *Client) DeleteMap(ctx context.Context, mapID int) error {
	_, err := c.del(ctx, fmt.Sprintf("maps/%d", mapID),
		nil, nil, false, "")
	return err
}
