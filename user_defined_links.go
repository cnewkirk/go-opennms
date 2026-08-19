package opennms

// User Defined Links REST API v2 – /api/v2/userdefinedlinks.

import (
	"context"
	"fmt"
)

// GetUserDefinedLinks lists all user-defined links.
func (c *Client) GetUserDefinedLinks(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "userdefinedlinks", nil, true)
}

// GetUserDefinedLink returns a specific user-defined link by ID.
func (c *Client) GetUserDefinedLink(ctx context.Context, linkID int) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("userdefinedlinks/%d", linkID), nil, true)
}

// CreateUserDefinedLink creates a new user-defined link. link is a
// link definition with keys such as "nodeIdA", "nodeIdZ",
// "componentLabelA", "componentLabelZ", "linkId", "linkLabel",
// "owner".
func (c *Client) CreateUserDefinedLink(ctx context.Context, link map[string]any) (map[string]any, error) {
	return asObject(c.post(ctx, "userdefinedlinks", link, nil, true))
}

// DeleteUserDefinedLink deletes a user-defined link.
func (c *Client) DeleteUserDefinedLink(ctx context.Context, linkID int) error {
	_, err := c.del(ctx, fmt.Sprintf("userdefinedlinks/%d", linkID),
		nil, nil, true, "")
	return err
}
