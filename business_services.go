package opennms

// Business Service Monitoring REST API v2 – /api/v2/business-services.

import (
	"context"
	"fmt"
	"net/url"
)

// GetBusinessServices lists all business services.
func (c *Client) GetBusinessServices(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "business-services", nil, true)
}

// GetBusinessService returns a specific business service by ID.
func (c *Client) GetBusinessService(ctx context.Context, serviceID int) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("business-services/%d", serviceID), nil, true)
}

// CreateBusinessService creates a new business service.
//
// service is the business service definition. Common keys: "name"
// (string), "attributes" (map of key/value pairs), "reduceFunction"
// (map with a "type" key), "edges" (list of edge configuration maps).
// Example:
//
//	map[string]any{
//		"name":           "My App",
//		"attributes":     map[string]any{"dc": "us-east-1"},
//		"reduceFunction": map[string]any{"type": "HighestSeverity"},
//	}
func (c *Client) CreateBusinessService(ctx context.Context, service map[string]any) (map[string]any, error) {
	return asObject(c.post(ctx, "business-services", service, nil, true))
}

// UpdateBusinessService updates the business service with the given
// database ID from the updated definition in service.
func (c *Client) UpdateBusinessService(ctx context.Context, serviceID int, service map[string]any) error {
	_, err := c.put(ctx, fmt.Sprintf("business-services/%d", serviceID), service, nil, true)
	return err
}

// DeleteBusinessService deletes a business service.
func (c *Client) DeleteBusinessService(ctx context.Context, serviceID int) error {
	_, err := c.del(ctx, fmt.Sprintf("business-services/%d", serviceID), nil, nil, true, "")
	return err
}

// Edges

// GetBusinessServiceEdge returns a specific business service edge by
// ID.
func (c *Client) GetBusinessServiceEdge(ctx context.Context, edgeID int) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("business-services/edges/%d", edgeID), nil, true)
}

// AddIpServiceEdge adds an IP-service edge to the business service
// with the given database ID. edge is the edge definition with keys
// such as "ipServiceId", "mapFunction", and "weight".
func (c *Client) AddIpServiceEdge(ctx context.Context, serviceID int, edge map[string]any) (map[string]any, error) {
	return asObject(c.post(ctx,
		fmt.Sprintf("business-services/%d/ip-service-edge", serviceID),
		edge, nil, true))
}

// AddReductionKeyEdge adds a reduction-key edge to the business
// service with the given database ID. edge is the edge definition
// with keys such as "reductionKey", "mapFunction", and "weight".
func (c *Client) AddReductionKeyEdge(ctx context.Context, serviceID int, edge map[string]any) (map[string]any, error) {
	return asObject(c.post(ctx,
		fmt.Sprintf("business-services/%d/reduction-key-edge", serviceID),
		edge, nil, true))
}

// AddChildEdge adds a child-service edge to the business service
// with the given database ID. edge is the edge definition with keys
// such as "childId", "mapFunction", and "weight".
func (c *Client) AddChildEdge(ctx context.Context, serviceID int, edge map[string]any) (map[string]any, error) {
	return asObject(c.post(ctx,
		fmt.Sprintf("business-services/%d/child-edge", serviceID),
		edge, nil, true))
}

// RemoveBusinessServiceEdge removes the edge with database ID edgeID
// from the business service with database ID serviceID.
func (c *Client) RemoveBusinessServiceEdge(ctx context.Context, serviceID, edgeID int) error {
	_, err := c.del(ctx,
		fmt.Sprintf("business-services/%d/edges/%d", serviceID, edgeID),
		nil, nil, true, "")
	return err
}

// Daemon

// ReloadBusinessServiceDaemon reloads the Business Service
// Monitoring daemon.
func (c *Client) ReloadBusinessServiceDaemon(ctx context.Context) error {
	_, err := c.post(ctx, "business-services/daemon/reload", nil, nil, true)
	return err
}

// Functions

// GetMapFunctions lists all available map functions.
func (c *Client) GetMapFunctions(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "business-services/functions/map", nil, true)
}

// GetMapFunction returns a specific map function by name.
func (c *Client) GetMapFunction(ctx context.Context, name string) (map[string]any, error) {
	return c.getObject(ctx, "business-services/functions/map/"+url.PathEscape(name), nil, true)
}

// GetReduceFunctions lists all available reduce functions.
func (c *Client) GetReduceFunctions(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "business-services/functions/reduce", nil, true)
}

// GetReduceFunction returns a specific reduce function by name.
func (c *Client) GetReduceFunction(ctx context.Context, name string) (map[string]any, error) {
	return c.getObject(ctx, "business-services/functions/reduce/"+url.PathEscape(name), nil, true)
}
