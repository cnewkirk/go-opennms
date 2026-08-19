package opennms

// Requisitions (Provisioning) REST API – /rest/requisitions.

import (
	"context"
	"fmt"
	"net/url"
)

// requisitionForm converts a field map into form-encoded values (the
// requisitions API documents "PUT requests expect form-urlencoded
// data").
func requisitionForm(data map[string]string) url.Values {
	form := url.Values{}
	for k, v := range data {
		form.Set(k, v)
	}
	return form
}

// Requisitions

// GetRequisitions lists all active (pending or deployed)
// requisitions.
func (c *Client) GetRequisitions(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "requisitions", nil, false)
}

// GetRequisition returns the requisition with the given foreign
// source name.
func (c *Client) GetRequisition(ctx context.Context, name string) (map[string]any, error) {
	return c.getObject(ctx, "requisitions/"+url.PathEscape(name), nil, false)
}

// GetRequisitionCount returns the count of undeployed (pending)
// requisitions.
func (c *Client) GetRequisitionCount(ctx context.Context) (int, error) {
	return c.getCount(ctx, "requisitions/count", false)
}

// GetDeployedRequisitions lists all deployed requisitions.
func (c *Client) GetDeployedRequisitions(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "requisitions/deployed", nil, false)
}

// GetDeployedRequisitionCount returns the count of deployed
// requisitions.
func (c *Client) GetDeployedRequisitionCount(ctx context.Context) (int, error) {
	return c.getCount(ctx, "requisitions/deployed/count", false)
}

// CreateRequisition adds or replaces a requisition. The requisition
// must contain at least {"foreign-source": "name", "node": []}.
func (c *Client) CreateRequisition(ctx context.Context, requisition map[string]any) error {
	_, err := c.post(ctx, "requisitions", requisition, nil, false)
	return err
}

// ImportRequisition synchronises/imports the named requisition into
// the database. When rescanExisting is false only new/removed nodes
// are processed and existing nodes are not rescanned.
func (c *Client) ImportRequisition(ctx context.Context, name string, rescanExisting bool) error {
	var params url.Values
	if !rescanExisting {
		params = url.Values{"rescanExisting": {"false"}}
	}
	_, err := c.put(ctx,
		fmt.Sprintf("requisitions/%s/import", url.PathEscape(name)),
		nil, params, false)
	return err
}

// UpdateRequisition updates metadata on an existing requisition.
// data holds the requisition fields to change; it is sent
// form-encoded.
func (c *Client) UpdateRequisition(ctx context.Context, name string, data map[string]string) error {
	_, err := c.putForm(ctx, "requisitions/"+url.PathEscape(name),
		requisitionForm(data), nil, false)
	return err
}

// DeleteRequisition deletes a pending (not yet deployed)
// requisition.
func (c *Client) DeleteRequisition(ctx context.Context, name string) error {
	_, err := c.del(ctx, "requisitions/"+url.PathEscape(name), nil, nil, false, "")
	return err
}

// DeleteDeployedRequisition deletes a deployed requisition.
func (c *Client) DeleteDeployedRequisition(ctx context.Context, name string) error {
	_, err := c.del(ctx, "requisitions/deployed/"+url.PathEscape(name),
		nil, nil, false, "")
	return err
}

// Requisition Nodes

// GetRequisitionNodes lists nodes in the named requisition.
func (c *Client) GetRequisitionNodes(ctx context.Context, name string) (map[string]any, error) {
	return c.getObject(ctx,
		fmt.Sprintf("requisitions/%s/nodes", url.PathEscape(name)),
		nil, false)
}

// GetRequisitionNode returns a specific node in a requisition by its
// foreign ID.
func (c *Client) GetRequisitionNode(ctx context.Context, name, foreignID string) (map[string]any, error) {
	return c.getObject(ctx,
		fmt.Sprintf("requisitions/%s/nodes/%s",
			url.PathEscape(name), url.PathEscape(foreignID)),
		nil, false)
}

// CreateRequisitionNode adds or replaces a node in the named
// requisition.
func (c *Client) CreateRequisitionNode(ctx context.Context, name string, node map[string]any) error {
	_, err := c.post(ctx,
		fmt.Sprintf("requisitions/%s/nodes", url.PathEscape(name)),
		node, nil, false)
	return err
}

// UpdateRequisitionNode updates a node in the named requisition.
// node holds the fields to change; it is sent form-encoded — the
// requisitions API documents "PUT requests expect form-urlencoded
// data".
func (c *Client) UpdateRequisitionNode(ctx context.Context, name, foreignID string, node map[string]string) error {
	_, err := c.putForm(ctx,
		fmt.Sprintf("requisitions/%s/nodes/%s",
			url.PathEscape(name), url.PathEscape(foreignID)),
		requisitionForm(node), nil, false)
	return err
}

// DeleteRequisitionNode deletes a node from the named requisition
// (async – returns 202).
func (c *Client) DeleteRequisitionNode(ctx context.Context, name, foreignID string) error {
	_, err := c.del(ctx,
		fmt.Sprintf("requisitions/%s/nodes/%s",
			url.PathEscape(name), url.PathEscape(foreignID)),
		nil, nil, false, "")
	return err
}

// Requisition Node Interfaces

// GetRequisitionNodeInterfaces lists interfaces for a node in a
// requisition.
func (c *Client) GetRequisitionNodeInterfaces(ctx context.Context, name, foreignID string) (map[string]any, error) {
	return c.getObject(ctx,
		fmt.Sprintf("requisitions/%s/nodes/%s/interfaces",
			url.PathEscape(name), url.PathEscape(foreignID)),
		nil, false)
}

// CreateRequisitionNodeInterface adds or replaces an interface on a
// requisition node. iface holds keys such as "ip-addr",
// "snmp-primary", "status", and "monitored-service".
func (c *Client) CreateRequisitionNodeInterface(ctx context.Context, name, foreignID string, iface map[string]any) error {
	_, err := c.post(ctx,
		fmt.Sprintf("requisitions/%s/nodes/%s/interfaces",
			url.PathEscape(name), url.PathEscape(foreignID)),
		iface, nil, false)
	return err
}

// UpdateRequisitionNodeInterface updates an interface on a
// requisition node. iface holds the fields to change; it is sent
// form-encoded — the requisitions API documents "PUT requests expect
// form-urlencoded data".
func (c *Client) UpdateRequisitionNodeInterface(ctx context.Context, name, foreignID, ipAddress string, iface map[string]string) error {
	_, err := c.putForm(ctx,
		fmt.Sprintf("requisitions/%s/nodes/%s/interfaces/%s",
			url.PathEscape(name), url.PathEscape(foreignID),
			url.PathEscape(ipAddress)),
		requisitionForm(iface), nil, false)
	return err
}

// DeleteRequisitionNodeInterface deletes an interface from a
// requisition node (async).
func (c *Client) DeleteRequisitionNodeInterface(ctx context.Context, name, foreignID, ipAddress string) error {
	_, err := c.del(ctx,
		fmt.Sprintf("requisitions/%s/nodes/%s/interfaces/%s",
			url.PathEscape(name), url.PathEscape(foreignID),
			url.PathEscape(ipAddress)),
		nil, nil, false, "")
	return err
}

// Requisition Node Interface Services

// GetRequisitionNodeServices lists monitored services on a
// requisition node interface.
func (c *Client) GetRequisitionNodeServices(ctx context.Context, name, foreignID, ipAddress string) (map[string]any, error) {
	return c.getObject(ctx,
		fmt.Sprintf("requisitions/%s/nodes/%s/interfaces/%s/services",
			url.PathEscape(name), url.PathEscape(foreignID),
			url.PathEscape(ipAddress)),
		nil, false)
}

// CreateRequisitionNodeService adds or replaces a service on a
// requisition node interface. Example service:
// {"service-name": "HTTP"}.
func (c *Client) CreateRequisitionNodeService(ctx context.Context, name, foreignID, ipAddress string, service map[string]any) error {
	_, err := c.post(ctx,
		fmt.Sprintf("requisitions/%s/nodes/%s/interfaces/%s/services",
			url.PathEscape(name), url.PathEscape(foreignID),
			url.PathEscape(ipAddress)),
		service, nil, false)
	return err
}

// DeleteRequisitionNodeService deletes a service from a requisition
// node interface (async).
func (c *Client) DeleteRequisitionNodeService(ctx context.Context, name, foreignID, ipAddress, serviceName string) error {
	_, err := c.del(ctx,
		fmt.Sprintf("requisitions/%s/nodes/%s/interfaces/%s/services/%s",
			url.PathEscape(name), url.PathEscape(foreignID),
			url.PathEscape(ipAddress), url.PathEscape(serviceName)),
		nil, nil, false, "")
	return err
}

// Requisition Node Categories

// GetRequisitionNodeCategories lists categories for a requisition
// node.
func (c *Client) GetRequisitionNodeCategories(ctx context.Context, name, foreignID string) (map[string]any, error) {
	return c.getObject(ctx,
		fmt.Sprintf("requisitions/%s/nodes/%s/categories",
			url.PathEscape(name), url.PathEscape(foreignID)),
		nil, false)
}

// AddRequisitionNodeCategory adds or replaces a category on a
// requisition node. Example category: {"name": "Production"}.
func (c *Client) AddRequisitionNodeCategory(ctx context.Context, name, foreignID string, category map[string]any) error {
	_, err := c.post(ctx,
		fmt.Sprintf("requisitions/%s/nodes/%s/categories",
			url.PathEscape(name), url.PathEscape(foreignID)),
		category, nil, false)
	return err
}

// DeleteRequisitionNodeCategory deletes a category from a
// requisition node (async).
func (c *Client) DeleteRequisitionNodeCategory(ctx context.Context, name, foreignID, category string) error {
	_, err := c.del(ctx,
		fmt.Sprintf("requisitions/%s/nodes/%s/categories/%s",
			url.PathEscape(name), url.PathEscape(foreignID),
			url.PathEscape(category)),
		nil, nil, false, "")
	return err
}

// Requisition Node Assets

// GetRequisitionNodeAssets lists asset fields for a requisition
// node.
func (c *Client) GetRequisitionNodeAssets(ctx context.Context, name, foreignID string) (map[string]any, error) {
	return c.getObject(ctx,
		fmt.Sprintf("requisitions/%s/nodes/%s/assets",
			url.PathEscape(name), url.PathEscape(foreignID)),
		nil, false)
}

// SetRequisitionNodeAsset adds or replaces an asset on a requisition
// node. Example asset: {"name": "serialNumber", "value": "SN-1234"}.
func (c *Client) SetRequisitionNodeAsset(ctx context.Context, name, foreignID string, asset map[string]any) error {
	_, err := c.post(ctx,
		fmt.Sprintf("requisitions/%s/nodes/%s/assets",
			url.PathEscape(name), url.PathEscape(foreignID)),
		asset, nil, false)
	return err
}

// DeleteRequisitionNodeAsset deletes an asset field from a
// requisition node (async).
func (c *Client) DeleteRequisitionNodeAsset(ctx context.Context, name, foreignID, field string) error {
	_, err := c.del(ctx,
		fmt.Sprintf("requisitions/%s/nodes/%s/assets/%s",
			url.PathEscape(name), url.PathEscape(foreignID),
			url.PathEscape(field)),
		nil, nil, false, "")
	return err
}
