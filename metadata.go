package opennms

// Metadata REST API v2 – /api/v2/nodes/{id}/metadata and
// sub-resources.

import (
	"context"
	"net/url"
)

// metadataNodePath builds the node metadata base path. nodeID is a
// node database ID or "foreignSource:foreignId".
func metadataNodePath(nodeID string) string {
	return "nodes/" + url.PathEscape(nodeID) + "/metadata"
}

// metadataInterfacePath builds the interface metadata base path.
func metadataInterfacePath(nodeID, ipInterface string) string {
	return "nodes/" + url.PathEscape(nodeID) + "/ipinterfaces/" +
		url.PathEscape(ipInterface) + "/metadata"
}

// metadataServicePath builds the service metadata base path.
func metadataServicePath(nodeID, ipInterface, service string) string {
	return "nodes/" + url.PathEscape(nodeID) + "/ipinterfaces/" +
		url.PathEscape(ipInterface) + "/services/" +
		url.PathEscape(service) + "/metadata"
}

// metadataSet posts each metadata entry individually — the endpoint
// accepts one entry per request.
func (c *Client) metadataSet(ctx context.Context, path string, metadata []map[string]any) error {
	for _, entry := range metadata {
		if _, err := c.post(ctx, path, entry, nil, true); err != nil {
			return err
		}
	}
	return nil
}

// Node metadata

// GetNodeMetadata returns all metadata for the node. nodeID is a node
// database ID or "foreignSource:foreignId".
func (c *Client) GetNodeMetadata(ctx context.Context, nodeID string) ([]any, error) {
	return c.getList(ctx, metadataNodePath(nodeID), nil, true)
}

// GetNodeMetadataContext returns all metadata for the node within the
// given context.
func (c *Client) GetNodeMetadataContext(ctx context.Context, nodeID, context string) ([]any, error) {
	return c.getList(ctx, metadataNodePath(nodeID)+"/"+url.PathEscape(context), nil, true)
}

// GetNodeMetadataValue returns a specific metadata value for the node.
func (c *Client) GetNodeMetadataValue(ctx context.Context, nodeID, context, key string) (map[string]any, error) {
	return c.getObject(ctx, metadataNodePath(nodeID)+"/"+
		url.PathEscape(context)+"/"+url.PathEscape(key), nil, true)
}

// SetNodeMetadata sets metadata entries for the node.
//
// The endpoint accepts one entry per request, so each entry is posted
// individually.
//
// nodeID is a node database ID or "foreignSource:foreignId". metadata
// is a list of maps with keys "context", "key", "value". Only
// user-defined contexts (prefixed "X-") can be modified.
func (c *Client) SetNodeMetadata(ctx context.Context, nodeID string, metadata []map[string]any) error {
	return c.metadataSet(ctx, metadataNodePath(nodeID), metadata)
}

// SetNodeMetadataValue sets a single metadata key/value for the node.
//
// nodeID is a node database ID or "foreignSource:foreignId". context
// must be user-defined (prefixed "X-").
func (c *Client) SetNodeMetadataValue(ctx context.Context, nodeID, context, key, value string) error {
	_, err := c.put(ctx, metadataNodePath(nodeID)+"/"+
		url.PathEscape(context)+"/"+url.PathEscape(key)+"/"+
		url.PathEscape(value), nil, nil, true)
	return err
}

// DeleteNodeMetadataContext deletes all metadata in the given context
// for the node.
func (c *Client) DeleteNodeMetadataContext(ctx context.Context, nodeID, context string) error {
	_, err := c.del(ctx, metadataNodePath(nodeID)+"/"+
		url.PathEscape(context), nil, nil, true, "")
	return err
}

// DeleteNodeMetadataKey deletes a specific metadata key for the node.
func (c *Client) DeleteNodeMetadataKey(ctx context.Context, nodeID, context, key string) error {
	_, err := c.del(ctx, metadataNodePath(nodeID)+"/"+
		url.PathEscape(context)+"/"+url.PathEscape(key), nil, nil, true, "")
	return err
}

// Interface metadata

// GetInterfaceMetadata returns all metadata for the interface
// ipInterface on the node.
func (c *Client) GetInterfaceMetadata(ctx context.Context, nodeID, ipInterface string) ([]any, error) {
	return c.getList(ctx, metadataInterfacePath(nodeID, ipInterface), nil, true)
}

// GetInterfaceMetadataContext returns metadata within the given
// context for the interface.
func (c *Client) GetInterfaceMetadataContext(ctx context.Context, nodeID, ipInterface, context string) ([]any, error) {
	return c.getList(ctx, metadataInterfacePath(nodeID, ipInterface)+"/"+
		url.PathEscape(context), nil, true)
}

// GetInterfaceMetadataValue returns a specific metadata value for the
// interface.
func (c *Client) GetInterfaceMetadataValue(ctx context.Context, nodeID, ipInterface, context, key string) (map[string]any, error) {
	return c.getObject(ctx, metadataInterfacePath(nodeID, ipInterface)+"/"+
		url.PathEscape(context)+"/"+url.PathEscape(key), nil, true)
}

// SetInterfaceMetadata sets metadata for the interface.
//
// nodeID is a node database ID or "foreignSource:foreignId".
// ipInterface is the IP address of the interface. metadata is a list
// of maps with keys "context", "key", "value"; each entry is posted
// individually — the endpoint accepts one entry per request.
func (c *Client) SetInterfaceMetadata(ctx context.Context, nodeID, ipInterface string, metadata []map[string]any) error {
	return c.metadataSet(ctx, metadataInterfacePath(nodeID, ipInterface), metadata)
}

// SetInterfaceMetadataValue sets a single metadata key/value for the
// interface.
//
// nodeID is a node database ID or "foreignSource:foreignId".
// ipInterface is the IP address of the interface. context must be
// user-defined (prefixed "X-").
func (c *Client) SetInterfaceMetadataValue(ctx context.Context, nodeID, ipInterface, context, key, value string) error {
	_, err := c.put(ctx, metadataInterfacePath(nodeID, ipInterface)+"/"+
		url.PathEscape(context)+"/"+url.PathEscape(key)+"/"+
		url.PathEscape(value), nil, nil, true)
	return err
}

// DeleteInterfaceMetadataContext deletes all metadata in the given
// context for the interface.
func (c *Client) DeleteInterfaceMetadataContext(ctx context.Context, nodeID, ipInterface, context string) error {
	_, err := c.del(ctx, metadataInterfacePath(nodeID, ipInterface)+"/"+
		url.PathEscape(context), nil, nil, true, "")
	return err
}

// DeleteInterfaceMetadataKey deletes a specific metadata key for the
// interface.
func (c *Client) DeleteInterfaceMetadataKey(ctx context.Context, nodeID, ipInterface, context, key string) error {
	_, err := c.del(ctx, metadataInterfacePath(nodeID, ipInterface)+"/"+
		url.PathEscape(context)+"/"+url.PathEscape(key), nil, nil, true, "")
	return err
}

// Service metadata

// GetServiceMetadata returns all metadata for the service on
// ipInterface of the node.
func (c *Client) GetServiceMetadata(ctx context.Context, nodeID, ipInterface, service string) ([]any, error) {
	return c.getList(ctx, metadataServicePath(nodeID, ipInterface, service), nil, true)
}

// GetServiceMetadataContext returns metadata within the given context
// for the service.
func (c *Client) GetServiceMetadataContext(ctx context.Context, nodeID, ipInterface, service, context string) ([]any, error) {
	return c.getList(ctx, metadataServicePath(nodeID, ipInterface, service)+"/"+
		url.PathEscape(context), nil, true)
}

// GetServiceMetadataValue returns a specific metadata value for the
// service.
func (c *Client) GetServiceMetadataValue(ctx context.Context, nodeID, ipInterface, service, context, key string) (map[string]any, error) {
	return c.getObject(ctx, metadataServicePath(nodeID, ipInterface, service)+"/"+
		url.PathEscape(context)+"/"+url.PathEscape(key), nil, true)
}

// SetServiceMetadata sets metadata for the service.
//
// nodeID is a node database ID or "foreignSource:foreignId".
// ipInterface is the IP address of the interface; service the
// monitored service name. metadata is a list of maps with keys
// "context", "key", "value"; each entry is posted individually — the
// endpoint accepts one entry per request.
func (c *Client) SetServiceMetadata(ctx context.Context, nodeID, ipInterface, service string, metadata []map[string]any) error {
	return c.metadataSet(ctx, metadataServicePath(nodeID, ipInterface, service), metadata)
}

// SetServiceMetadataValue sets a single metadata key/value for the
// service.
//
// nodeID is a node database ID or "foreignSource:foreignId".
// ipInterface is the IP address of the interface; service the
// monitored service name. context must be user-defined (prefixed
// "X-").
func (c *Client) SetServiceMetadataValue(ctx context.Context, nodeID, ipInterface, service, context, key, value string) error {
	_, err := c.put(ctx, metadataServicePath(nodeID, ipInterface, service)+"/"+
		url.PathEscape(context)+"/"+url.PathEscape(key)+"/"+
		url.PathEscape(value), nil, nil, true)
	return err
}

// DeleteServiceMetadataContext deletes all metadata in the given
// context for the service.
func (c *Client) DeleteServiceMetadataContext(ctx context.Context, nodeID, ipInterface, service, context string) error {
	_, err := c.del(ctx, metadataServicePath(nodeID, ipInterface, service)+"/"+
		url.PathEscape(context), nil, nil, true, "")
	return err
}

// DeleteServiceMetadataKey deletes a specific metadata key for the
// service.
func (c *Client) DeleteServiceMetadataKey(ctx context.Context, nodeID, ipInterface, service, context, key string) error {
	_, err := c.del(ctx, metadataServicePath(nodeID, ipInterface, service)+"/"+
		url.PathEscape(context)+"/"+url.PathEscape(key), nil, nil, true, "")
	return err
}
