package opennms

// Requisition Names REST API – /rest/requisitionNames.

import "context"

// GetRequisitionNames lists all requisition (foreign source) names.
func (c *Client) GetRequisitionNames(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "requisitionNames", nil, false)
}
