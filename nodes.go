package opennms

// Nodes REST API – /rest/nodes and sub-resources.

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

var nodesXMLEscaper = strings.NewReplacer(
	"&", "&amp;", "<", "&lt;", ">", "&gt;",
)

var nodesXMLAttrEscaper = strings.NewReplacer(
	"&", "&amp;", "<", "&lt;", ">", "&gt;",
	"\n", "&#10;", "\r", "&#13;", "\t", "&#9;",
)

// nodesQuoteAttr quotes an XML attribute value, mirroring Python's
// xml.sax.saxutils.quoteattr.
func nodesQuoteAttr(s string) string {
	s = nodesXMLAttrEscaper.Replace(s)
	if strings.Contains(s, `"`) {
		if strings.Contains(s, "'") {
			return `"` + strings.ReplaceAll(s, `"`, "&quot;") + `"`
		}
		return "'" + s + "'"
	}
	return `"` + s + `"`
}

// nodesXMLElement serializes data into a flat JAXB-shaped XML
// element: keys in attrKeys become attributes, keys in elemKeys
// become child elements, in the given order, skipping absent keys.
func nodesXMLElement(tag string, data map[string]any, attrKeys, elemKeys []string) string {
	var b strings.Builder
	b.WriteString("<" + tag)
	for _, k := range attrKeys {
		if v, ok := data[k]; ok {
			b.WriteString(" " + k + "=" + nodesQuoteAttr(fmt.Sprintf("%v", v)))
		}
	}
	b.WriteString(">")
	for _, k := range elemKeys {
		if v, ok := data[k]; ok {
			b.WriteString("<" + k + ">" + nodesXMLEscaper.Replace(fmt.Sprintf("%v", v)) + "</" + k + ">")
		}
	}
	b.WriteString("</" + tag + ">")
	return b.String()
}

// nodesFormValues converts a string map into form values.
func nodesFormValues(data map[string]string) url.Values {
	form := url.Values{}
	for k, v := range data {
		form.Set(k, v)
	}
	return form
}

// Nodes

// GetNodes lists nodes.
//
// filters are additional Hibernate query filters passed directly as
// query parameters (e.g. "label": "myrouter",
// "category": "Production").
func (c *Client) GetNodes(ctx context.Context, opts *ListOptions, filters map[string]string) (map[string]any, error) {
	params := mergeFilters(listParams(opts), filters)
	return c.getObject(ctx, "nodes", params, false)
}

// GetNode returns a single node by database ID or
// "foreignSource:foreignId".
func (c *Client) GetNode(ctx context.Context, nodeID string) (map[string]any, error) {
	return c.getObject(ctx, "nodes/"+url.PathEscape(nodeID), nil, false)
}

// GetNodeCount returns the total number of nodes.
//
// Uses the v2 API because the v1 /nodes endpoint does not expose a
// /count sub-resource.
func (c *Client) GetNodeCount(ctx context.Context) (int, error) {
	params := url.Values{"limit": {"1"}, "offset": {"0"}}
	result, err := c.getObject(ctx, "nodes", params, true)
	if err != nil {
		return 0, err
	}
	if n, ok := result["totalCount"].(float64); ok {
		return int(n), nil
	}
	return 0, nil
}

// CreateNode creates a node via the v1 API (POST /rest/nodes).
//
// node is a node attribute map. Common keys: "label", "type",
// "foreignSource", "foreignId", "location", "sysContact". The body
// is sent as XML — the nodes API documents "POST requires XML using
// application/xml as its Content-Type". For provisioning, prefer the
// requisitions API (CreateRequisitionNode).
func (c *Client) CreateNode(ctx context.Context, node map[string]any) (map[string]any, error) {
	xml := nodesXMLElement("node", node,
		[]string{"label", "type", "foreignSource", "foreignId"},
		[]string{"labelSource", "sysObjectId", "sysName",
			"sysDescription", "sysLocation", "sysContact",
			"location"})
	return asObject(c.postText(ctx, "nodes", xml, "application/xml", false, "", nil))
}

// UpdateNode updates node properties (PUT /rest/nodes/{id}).
//
// node holds the node fields to change. Sent form-encoded — the
// nodes API documents "PUT requires form data".
func (c *Client) UpdateNode(ctx context.Context, nodeID string, node map[string]string) error {
	_, err := c.putForm(ctx, "nodes/"+url.PathEscape(nodeID),
		nodesFormValues(node), nil, false)
	return err
}

// DeleteNode deletes a node (async – returns 202 Accepted).
func (c *Client) DeleteNode(ctx context.Context, nodeID string) error {
	_, err := c.del(ctx, "nodes/"+url.PathEscape(nodeID), nil, nil, false, "")
	return err
}

// RescanNode triggers a capability scan of the node for new
// interfaces/services.
func (c *Client) RescanNode(ctx context.Context, nodeID string) error {
	_, err := c.putForm(ctx, "nodes/"+url.PathEscape(nodeID)+"/rescan",
		url.Values{}, nil, false)
	return err
}

// IP Interfaces

// GetNodeIpInterfaces lists IP interfaces for the node.
//
// filters are additional Hibernate query filters passed directly as
// query parameters.
func (c *Client) GetNodeIpInterfaces(ctx context.Context, nodeID string, opts *ListOptions, filters map[string]string) (map[string]any, error) {
	params := mergeFilters(listParams(opts), filters)
	return c.getObject(ctx, "nodes/"+url.PathEscape(nodeID)+"/ipinterfaces", params, false)
}

// GetNodeIpInterface returns a specific IP interface by IP address.
func (c *Client) GetNodeIpInterface(ctx context.Context, nodeID, ipAddress string) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("nodes/%s/ipinterfaces/%s",
		url.PathEscape(nodeID), url.PathEscape(ipAddress)), nil, false)
}

// CreateNodeIpInterface adds an IP interface to the node.
//
// iface is an interface map, e.g.
// {"ipAddress": "192.168.0.1", "snmpPrimary": "P", "isManaged": "M"}.
// The body is sent as XML (required by the v1 nodes API).
func (c *Client) CreateNodeIpInterface(ctx context.Context, nodeID string, iface map[string]any) (map[string]any, error) {
	xml := nodesXMLElement("ipInterface", iface,
		[]string{"ifIndex", "isManaged", "snmpPrimary"},
		[]string{"ipAddress", "hostName"})
	return asObject(c.postText(ctx,
		"nodes/"+url.PathEscape(nodeID)+"/ipinterfaces",
		xml, "application/xml", false, "", nil))
}

// UpdateNodeIpInterface updates an IP interface. iface holds the
// interface fields to change; sent form-encoded (required by the v1
// nodes API).
func (c *Client) UpdateNodeIpInterface(ctx context.Context, nodeID, ipAddress string, iface map[string]string) error {
	_, err := c.putForm(ctx, fmt.Sprintf("nodes/%s/ipinterfaces/%s",
		url.PathEscape(nodeID), url.PathEscape(ipAddress)),
		nodesFormValues(iface), nil, false)
	return err
}

// DeleteNodeIpInterface deletes an IP interface (async – returns 202
// Accepted).
func (c *Client) DeleteNodeIpInterface(ctx context.Context, nodeID, ipAddress string) error {
	_, err := c.del(ctx, fmt.Sprintf("nodes/%s/ipinterfaces/%s",
		url.PathEscape(nodeID), url.PathEscape(ipAddress)), nil, nil, false, "")
	return err
}

// Monitored Services

// GetNodeIpServices lists monitored services on the given IP address
// of the node.
func (c *Client) GetNodeIpServices(ctx context.Context, nodeID, ipAddress string) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("nodes/%s/ipinterfaces/%s/services",
		url.PathEscape(nodeID), url.PathEscape(ipAddress)), nil, false)
}

// GetNodeIpService returns a specific monitored service.
func (c *Client) GetNodeIpService(ctx context.Context, nodeID, ipAddress, service string) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("nodes/%s/ipinterfaces/%s/services/%s",
		url.PathEscape(nodeID), url.PathEscape(ipAddress),
		url.PathEscape(service)), nil, false)
}

// CreateNodeIpService adds a monitored service to the given IP
// address, e.g. service = {"serviceType": {"name": "HTTP"}}.
func (c *Client) CreateNodeIpService(ctx context.Context, nodeID, ipAddress string, service map[string]any) (map[string]any, error) {
	return asObject(c.post(ctx, fmt.Sprintf("nodes/%s/ipinterfaces/%s/services",
		url.PathEscape(nodeID), url.PathEscape(ipAddress)),
		service, nil, false))
}

// DeleteNodeIpService deletes a monitored service (async – returns
// 202 Accepted).
func (c *Client) DeleteNodeIpService(ctx context.Context, nodeID, ipAddress, service string) error {
	_, err := c.del(ctx, fmt.Sprintf("nodes/%s/ipinterfaces/%s/services/%s",
		url.PathEscape(nodeID), url.PathEscape(ipAddress),
		url.PathEscape(service)), nil, nil, false, "")
	return err
}

// SNMP Interfaces

// GetNodeSnmpInterfaces lists SNMP interfaces for the node.
//
// filters are additional Hibernate query filters passed directly as
// query parameters.
func (c *Client) GetNodeSnmpInterfaces(ctx context.Context, nodeID string, opts *ListOptions, filters map[string]string) (map[string]any, error) {
	params := mergeFilters(listParams(opts), filters)
	return c.getObject(ctx, "nodes/"+url.PathEscape(nodeID)+"/snmpinterfaces", params, false)
}

// GetNodeSnmpInterface returns a specific SNMP interface by ifIndex.
func (c *Client) GetNodeSnmpInterface(ctx context.Context, nodeID string, ifIndex int) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("nodes/%s/snmpinterfaces/%d",
		url.PathEscape(nodeID), ifIndex), nil, false)
}

// CreateNodeSnmpInterface adds an SNMP interface to the node.
//
// The body is sent as XML (required by the v1 nodes API). The
// "collect" and "poll" keys map to the "collectFlag" and "pollFlag"
// XML attributes.
func (c *Client) CreateNodeSnmpInterface(ctx context.Context, nodeID string, iface map[string]any) (map[string]any, error) {
	payload := make(map[string]any, len(iface))
	for k, v := range iface {
		payload[k] = v
	}
	if v, ok := payload["collect"]; ok {
		payload["collectFlag"] = v
		delete(payload, "collect")
	}
	if v, ok := payload["poll"]; ok {
		payload["pollFlag"] = v
		delete(payload, "poll")
	}
	xml := nodesXMLElement("snmpInterface", payload,
		[]string{"ifIndex", "collectFlag", "pollFlag"},
		[]string{"ifDescr", "ifName", "ifAlias", "ifType",
			"ifSpeed", "ifAdminStatus", "ifOperStatus",
			"netMask", "physAddr"})
	return asObject(c.postText(ctx,
		"nodes/"+url.PathEscape(nodeID)+"/snmpinterfaces",
		xml, "application/xml", false, "", nil))
}

// UpdateNodeSnmpInterface updates an SNMP interface. iface holds the
// interface fields to change; sent form-encoded (required by the v1
// nodes API).
func (c *Client) UpdateNodeSnmpInterface(ctx context.Context, nodeID string, ifIndex int, iface map[string]string) error {
	_, err := c.putForm(ctx, fmt.Sprintf("nodes/%s/snmpinterfaces/%d",
		url.PathEscape(nodeID), ifIndex),
		nodesFormValues(iface), nil, false)
	return err
}

// DeleteNodeSnmpInterface deletes an SNMP interface (sync – returns
// 204).
func (c *Client) DeleteNodeSnmpInterface(ctx context.Context, nodeID string, ifIndex int) error {
	_, err := c.del(ctx, fmt.Sprintf("nodes/%s/snmpinterfaces/%d",
		url.PathEscape(nodeID), ifIndex), nil, nil, false, "")
	return err
}

// Categories

// GetNodeCategories lists surveillance categories for the node.
func (c *Client) GetNodeCategories(ctx context.Context, nodeID string) (map[string]any, error) {
	return c.getObject(ctx, "nodes/"+url.PathEscape(nodeID)+"/categories", nil, false)
}

// GetNodeCategory returns a specific category association for the
// node.
func (c *Client) GetNodeCategory(ctx context.Context, nodeID, category string) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("nodes/%s/categories/%s",
		url.PathEscape(nodeID), url.PathEscape(category)), nil, false)
}

// AddNodeCategory adds a category to the node, e.g.
// category = {"name": "Production"}. The body is sent as XML
// (required by the v1 nodes API).
func (c *Client) AddNodeCategory(ctx context.Context, nodeID string, category map[string]any) (map[string]any, error) {
	xml := nodesXMLElement("category", category,
		[]string{"name"}, []string{"description"})
	return asObject(c.postText(ctx,
		"nodes/"+url.PathEscape(nodeID)+"/categories",
		xml, "application/xml", false, "", nil))
}

// UpdateNodeCategory updates a category association for the node.
// data holds the category fields to change; sent form-encoded
// (required by the v1 nodes API).
func (c *Client) UpdateNodeCategory(ctx context.Context, nodeID, category string, data map[string]string) error {
	_, err := c.putForm(ctx, fmt.Sprintf("nodes/%s/categories/%s",
		url.PathEscape(nodeID), url.PathEscape(category)),
		nodesFormValues(data), nil, false)
	return err
}

// DeleteNodeCategory removes a category from the node (sync –
// returns 204).
func (c *Client) DeleteNodeCategory(ctx context.Context, nodeID, category string) error {
	_, err := c.del(ctx, fmt.Sprintf("nodes/%s/categories/%s",
		url.PathEscape(nodeID), url.PathEscape(category)), nil, nil, false, "")
	return err
}

// Asset Record

// GetNodeAssetRecord returns the asset record for the node.
func (c *Client) GetNodeAssetRecord(ctx context.Context, nodeID string) (map[string]any, error) {
	return c.getObject(ctx, "nodes/"+url.PathEscape(nodeID)+"/assetRecord", nil, false)
}

// UpdateNodeAssetRecord updates the asset record for the node.
//
// asset holds the asset fields to change. Common fields:
// "description", "building", "floor", "room", "rack", "vendor",
// "modelNumber", "serialNumber", "operatingSystem". Sent
// form-encoded (required by the v1 nodes API).
func (c *Client) UpdateNodeAssetRecord(ctx context.Context, nodeID string, asset map[string]string) error {
	_, err := c.putForm(ctx, "nodes/"+url.PathEscape(nodeID)+"/assetRecord",
		nodesFormValues(asset), nil, false)
	return err
}

// Hardware Inventory

// GetNodeHardwareInventory returns the hardware inventory tree for
// the node.
func (c *Client) GetNodeHardwareInventory(ctx context.Context, nodeID string) (map[string]any, error) {
	return c.getObject(ctx, "nodes/"+url.PathEscape(nodeID)+"/hardwareInventory", nil, false)
}

// GetNodeHardwareEntity returns a specific hardware entity by
// entPhysicalIndex.
func (c *Client) GetNodeHardwareEntity(ctx context.Context, nodeID string, entPhysicalIndex int) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("nodes/%s/hardwareInventory/%d",
		url.PathEscape(nodeID), entPhysicalIndex), nil, false)
}

// AddNodeHardwareInventory adds a hardware inventory entry to the
// node.
func (c *Client) AddNodeHardwareInventory(ctx context.Context, nodeID string, data map[string]any) (map[string]any, error) {
	return asObject(c.post(ctx,
		"nodes/"+url.PathEscape(nodeID)+"/hardwareInventory",
		data, nil, false))
}

// UpdateNodeHardwareEntity updates a specific hardware entity
// (ENTITY-MIB entPhysicalIndex). data holds the entity fields to
// change; sent form-encoded (required by the v1 nodes API).
func (c *Client) UpdateNodeHardwareEntity(ctx context.Context, nodeID string, entPhysicalIndex int, data map[string]string) error {
	_, err := c.putForm(ctx, fmt.Sprintf("nodes/%s/hardwareInventory/%d",
		url.PathEscape(nodeID), entPhysicalIndex),
		nodesFormValues(data), nil, false)
	return err
}

// DeleteNodeHardwareEntity deletes a specific hardware entity.
func (c *Client) DeleteNodeHardwareEntity(ctx context.Context, nodeID string, entPhysicalIndex int) error {
	_, err := c.del(ctx, fmt.Sprintf("nodes/%s/hardwareInventory/%d",
		url.PathEscape(nodeID), entPhysicalIndex), nil, nil, false, "")
	return err
}
