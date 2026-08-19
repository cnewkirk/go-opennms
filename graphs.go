package opennms

// Graph / Topology REST API – /rest/graphs.

import (
	"context"
	"net/url"
)

// GetGraphContainers lists all registered graph containers and their
// metadata.
func (c *Client) GetGraphContainers(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "graphs", nil, false)
}

// GetGraphContainer returns a specific graph container by ID.
func (c *Client) GetGraphContainer(ctx context.Context, containerID string) (map[string]any, error) {
	return c.getObject(ctx, "graphs/"+url.PathEscape(containerID), nil, false)
}

// GetGraph returns a graph by namespace from the given container.
func (c *Client) GetGraph(ctx context.Context, containerID, namespace string) (map[string]any, error) {
	path := "graphs/" + url.PathEscape(containerID) + "/" + url.PathEscape(namespace)
	return c.getObject(ctx, path, nil, false)
}

// GetGraphView returns a focused graph view via POST.
//
// semanticZoomLevel is the SZL (the server default is 1).
// verticesInFocus is a list of vertex ID strings in the form
// "namespace:vertex_id"; nil focuses nothing.
func (c *Client) GetGraphView(ctx context.Context, containerID, namespace string, semanticZoomLevel int, verticesInFocus []string) (map[string]any, error) {
	if verticesInFocus == nil {
		verticesInFocus = []string{}
	}
	body := map[string]any{
		"semanticZoomLevel": semanticZoomLevel,
		"verticesInFocus":   verticesInFocus,
	}
	path := "graphs/" + url.PathEscape(containerID) + "/" + url.PathEscape(namespace)
	return asObject(c.post(ctx, path, body, nil, false))
}

// GetGraphSearchSuggestions returns search suggestions for graph
// elements in the namespace, matching searchTerm against graph
// element labels.
func (c *Client) GetGraphSearchSuggestions(ctx context.Context, namespace, searchTerm string) (map[string]any, error) {
	params := url.Values{"s": {searchTerm}}
	return c.getObject(ctx, "graphs/search/suggestions/"+url.PathEscape(namespace), params, false)
}

// GetGraphSearchResults returns search results for graph elements in
// the namespace.
//
// providerID optionally restricts results to one search provider,
// criteria is an optional search criteria string, and context is an
// optional context string passed to the search provider ("" = absent
// for all three).
func (c *Client) GetGraphSearchResults(ctx context.Context, namespace, providerID, criteria, context string) (map[string]any, error) {
	params := url.Values{}
	if providerID != "" {
		params.Set("providerId", providerID)
	}
	if criteria != "" {
		params.Set("criteria", criteria)
	}
	if context != "" {
		params.Set("context", context)
	}
	return c.getObject(ctx, "graphs/search/results/"+url.PathEscape(namespace), params, false)
}

// Prefab graphs

// GetPrefabGraphNames lists all prefab graph names.
//
// Hits the same path as GetGraphContainers but documents the expected
// return type when the server is configured for prefab graphs (a list
// of name strings).
func (c *Client) GetPrefabGraphNames(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "graphs", nil, false)
}

// GetPrefabGraph returns a specific prefab graph definition by name.
func (c *Client) GetPrefabGraph(ctx context.Context, name string) (map[string]any, error) {
	return c.getObject(ctx, "graphs/"+url.PathEscape(name), nil, false)
}

// GetPrefabGraphsForResource lists prefab graphs available for a
// given resource ID, e.g.
// "node[1].interfaceSnmp[eth0-04013f75f101]".
func (c *Client) GetPrefabGraphsForResource(ctx context.Context, resourceID string) (map[string]any, error) {
	return c.getObject(ctx, "graphs/for/"+url.PathEscape(resourceID), nil, false)
}

// GetPrefabGraphsForNode lists prefab graphs available for a given
// node. nodeCriteria is the node DB ID or "foreignSource:foreignId".
func (c *Client) GetPrefabGraphsForNode(ctx context.Context, nodeCriteria string) (map[string]any, error) {
	return c.getObject(ctx, "graphs/fornode/"+url.PathEscape(nodeCriteria), nil, false)
}
