package opennms

// IP Interfaces REST API v2 – /api/v2/ipinterfaces (read-only).

import "context"

// GetIpInterfaces lists IP interfaces using the v2 API with an
// optional FIQL filter string (e.g. "node.label==onms-prd-01",
// "ipAddress==192.168.32.140", "node.foreignSource==Servers");
// "" lists all.
//
// This is a global (cross-node) read-only view. For write operations
// use CreateNodeIpInterface / DeleteNodeIpInterface.
func (c *Client) GetIpInterfaces(ctx context.Context, fiql string, opts *ListOptions) (map[string]any, error) {
	params := listParams(opts)
	if fiql != "" {
		params.Set("_s", fiql)
	}
	return c.getObject(ctx, "ipinterfaces", params, true)
}
