package opennms

// SNMP Interfaces REST API v2 – /api/v2/snmpinterfaces (read-only).

import "context"

// GetSnmpInterfaces lists SNMP interfaces using the v2 API with an
// optional FIQL filter string ("" lists all).
//
// This is a global (cross-node) read-only view. For write operations
// use CreateNodeSnmpInterface / DeleteNodeSnmpInterface.
//
// fiql examples: "node.label==onms-prd-01", "ifIndex==6",
// "node.foreignSource==Servers;ipInterfaces.ipAddress=127.0.0.1".
func (c *Client) GetSnmpInterfaces(ctx context.Context, fiql string, opts *ListOptions) (map[string]any, error) {
	params := listParams(opts)
	if fiql != "" {
		params.Set("_s", fiql)
	}
	return c.getObject(ctx, "snmpinterfaces", params, true)
}
