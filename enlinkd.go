package opennms

// EnLinkd REST API v2 – /api/v2/enlinkd.

import (
	"context"
	"fmt"
	"net/url"
)

// GetNodeEnlinkd returns all EnLinkd topology data for a node
// (aggregate). nodeCriteria is a node DB ID or
// "foreignSource:foreignId".
//
// The result carries lldpLinkNodes, cdpLinkNodes, ospfLinkNodes,
// isisLinkNodes, bridgeLinkNodes (lists), plus element maps/lists for
// each protocol.
func (c *Client) GetNodeEnlinkd(ctx context.Context, nodeCriteria string) (map[string]any, error) {
	return c.getObject(ctx,
		fmt.Sprintf("enlinkd/%s", url.PathEscape(nodeCriteria)), nil, true)
}

// GetNodeLldpLinks returns LLDP link topology for a node.
// nodeCriteria is a node DB ID or "foreignSource:foreignId".
func (c *Client) GetNodeLldpLinks(ctx context.Context, nodeCriteria string) ([]any, error) {
	return c.getList(ctx,
		fmt.Sprintf("enlinkd/lldp_links/%s", url.PathEscape(nodeCriteria)), nil, true)
}

// GetNodeCdpLinks returns CDP link topology for a node. nodeCriteria
// is a node DB ID or "foreignSource:foreignId".
func (c *Client) GetNodeCdpLinks(ctx context.Context, nodeCriteria string) ([]any, error) {
	return c.getList(ctx,
		fmt.Sprintf("enlinkd/cdp_links/%s", url.PathEscape(nodeCriteria)), nil, true)
}

// GetNodeOspfLinks returns OSPF link topology for a node.
// nodeCriteria is a node DB ID or "foreignSource:foreignId".
func (c *Client) GetNodeOspfLinks(ctx context.Context, nodeCriteria string) ([]any, error) {
	return c.getList(ctx,
		fmt.Sprintf("enlinkd/ospf_links/%s", url.PathEscape(nodeCriteria)), nil, true)
}

// GetNodeIsisLinks returns IS-IS link topology for a node.
// nodeCriteria is a node DB ID or "foreignSource:foreignId".
func (c *Client) GetNodeIsisLinks(ctx context.Context, nodeCriteria string) ([]any, error) {
	return c.getList(ctx,
		fmt.Sprintf("enlinkd/isis_links/%s", url.PathEscape(nodeCriteria)), nil, true)
}

// GetNodeBridgeLinks returns bridge link topology for a node.
// nodeCriteria is a node DB ID or "foreignSource:foreignId".
func (c *Client) GetNodeBridgeLinks(ctx context.Context, nodeCriteria string) ([]any, error) {
	return c.getList(ctx,
		fmt.Sprintf("enlinkd/bridge_links/%s", url.PathEscape(nodeCriteria)), nil, true)
}

// GetNodeLldpElement returns the LLDP element (chassis info) for a
// node. nodeCriteria is a node DB ID or "foreignSource:foreignId".
func (c *Client) GetNodeLldpElement(ctx context.Context, nodeCriteria string) (map[string]any, error) {
	return c.getObject(ctx,
		fmt.Sprintf("enlinkd/lldp_elems/%s", url.PathEscape(nodeCriteria)), nil, true)
}

// GetNodeCdpElement returns the CDP element (global config) for a
// node. nodeCriteria is a node DB ID or "foreignSource:foreignId".
func (c *Client) GetNodeCdpElement(ctx context.Context, nodeCriteria string) (map[string]any, error) {
	return c.getObject(ctx,
		fmt.Sprintf("enlinkd/cdp_elems/%s", url.PathEscape(nodeCriteria)), nil, true)
}

// GetNodeOspfElement returns the OSPF element (router info) for a
// node. nodeCriteria is a node DB ID or "foreignSource:foreignId".
func (c *Client) GetNodeOspfElement(ctx context.Context, nodeCriteria string) (map[string]any, error) {
	return c.getObject(ctx,
		fmt.Sprintf("enlinkd/ospf_elems/%s", url.PathEscape(nodeCriteria)), nil, true)
}

// GetNodeIsisElement returns the IS-IS element (system info) for a
// node. nodeCriteria is a node DB ID or "foreignSource:foreignId".
func (c *Client) GetNodeIsisElement(ctx context.Context, nodeCriteria string) (map[string]any, error) {
	return c.getObject(ctx,
		fmt.Sprintf("enlinkd/isis_elems/%s", url.PathEscape(nodeCriteria)), nil, true)
}

// GetNodeBridgeElements returns bridge elements (one per VLAN) for a
// node. nodeCriteria is a node DB ID or "foreignSource:foreignId".
func (c *Client) GetNodeBridgeElements(ctx context.Context, nodeCriteria string) ([]any, error) {
	return c.getList(ctx,
		fmt.Sprintf("enlinkd/bridge_elems/%s", url.PathEscape(nodeCriteria)), nil, true)
}
