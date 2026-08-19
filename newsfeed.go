package opennms

// News Feed REST API v2 – /api/v2/newsfeed.

import "context"

// GetNewsfeed returns the latest OpenNMS news feed items, including
// categories, tags, title, link, and descriptions.
func (c *Client) GetNewsfeed(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "newsfeed", nil, true)
}
