package opennms

// Whoami REST API – /rest/whoami.

import "context"

// GetWhoami returns information about the currently authenticated
// user.
func (c *Client) GetWhoami(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "whoami", nil, false)
}
