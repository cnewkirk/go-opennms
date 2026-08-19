package opennms

// SNMP Metadata REST API v2 – /api/v2/snmpmetadata.

import (
	"context"
	"fmt"
)

// GetSnmpMetadata returns the SNMP metadata collected for the node
// with the given database ID.
func (c *Client) GetSnmpMetadata(ctx context.Context, nodeID int) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("snmpmetadata/%d", nodeID), nil, true)
}
