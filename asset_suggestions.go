package opennms

// Asset Suggestions REST API – /rest/assets/suggestions.

import "context"

// GetAssetSuggestions returns asset field suggestions based on
// existing inventory data.
func (c *Client) GetAssetSuggestions(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "assets/suggestions", nil, false)
}
