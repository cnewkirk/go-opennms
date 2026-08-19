package opennms

// Resources REST API – /rest/resources.

import (
	"context"
	"net/url"
	"strconv"
	"strings"
)

// resourcesDepthParams builds the depth query parameter both tree
// endpoints always send.
func resourcesDepthParams(depth int) url.Values {
	return url.Values{"depth": {strconv.Itoa(depth)}}
}

// GetResources returns the full resource tree (can be expensive on
// large systems). depth is the tree depth limit — the Python wrapper
// defaults to 1; pass -1 for unlimited depth.
func (c *Client) GetResources(ctx context.Context, depth int) (map[string]any, error) {
	return c.getObject(ctx, "resources", resourcesDepthParams(depth), false)
}

// GetResource returns the resource tree rooted at resourceID.
//
// resourceID examples: "node[1]", "node[Servers:router01]",
// "node[1].interfaceSnmp[eth0-04013f75f101]".
//
// depth is the tree depth limit — the Python wrapper defaults to -1,
// which returns the single resource.
func (c *Client) GetResource(ctx context.Context, resourceID string, depth int) (map[string]any, error) {
	return c.getObject(ctx, "resources/"+url.PathEscape(resourceID),
		resourcesDepthParams(depth), false)
}

// GetResourcesForNode returns all resources for a node.
// nodeCriteria is the node DB ID ("1") or "foreignSource:foreignId".
func (c *Client) GetResourcesForNode(ctx context.Context, nodeCriteria string) (map[string]any, error) {
	return c.getObject(ctx, "resources/fornode/"+url.PathEscape(nodeCriteria),
		nil, false)
}

// GetResourcesSelect returns a partial selection of the resource
// tree. nodes are node IDs or "FS:FID" strings; filterRules are
// filter rule strings; nodeSubresources are subresource names;
// stringProperties are string property names to include. Each nil or
// empty slice is omitted.
func (c *Client) GetResourcesSelect(ctx context.Context, nodes, filterRules, nodeSubresources, stringProperties []string) (map[string]any, error) {
	params := url.Values{}
	if len(nodes) > 0 {
		params.Set("nodes", strings.Join(nodes, ","))
	}
	if len(filterRules) > 0 {
		params.Set("filterRules", strings.Join(filterRules, ","))
	}
	if len(nodeSubresources) > 0 {
		params.Set("nodeSubresources", strings.Join(nodeSubresources, ","))
	}
	if len(stringProperties) > 0 {
		params.Set("stringProperties", strings.Join(stringProperties, ","))
	}
	return c.getObject(ctx, "resources/select", params, false)
}

// DeleteResource deletes a resource and all its child resources.
// resourceID is a resource ID string (see GetResource).
func (c *Client) DeleteResource(ctx context.Context, resourceID string) error {
	_, err := c.del(ctx, "resources/"+url.PathEscape(resourceID),
		nil, nil, false, "")
	return err
}
