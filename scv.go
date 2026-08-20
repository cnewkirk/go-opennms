package opennms

// Secure Credentials Vault REST API – /rest/scv.

import (
	"context"
	"net/url"
)

// GetCredentials lists all stored credentials.
func (c *Client) GetCredentials(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "scv", nil, false)
}

// GetCredential returns a specific credential by alias (the unique
// key).
func (c *Client) GetCredential(ctx context.Context, alias string) (map[string]any, error) {
	return c.getObject(ctx, "scv/"+url.PathEscape(alias), nil, false)
}

// CreateCredential creates a new credential entry. data is the
// credential definition map with keys "alias", "username",
// "password", and optionally "attributes".
func (c *Client) CreateCredential(ctx context.Context, data map[string]any) error {
	_, err := c.post(ctx, "scv", data, nil, false)
	return err
}

// UpdateCredential updates an existing credential by alias.
func (c *Client) UpdateCredential(ctx context.Context, alias string, data map[string]any) error {
	_, err := c.put(ctx, "scv/"+url.PathEscape(alias), data, nil, false)
	return err
}

// DeleteCredential deletes a credential.
func (c *Client) DeleteCredential(ctx context.Context, alias string) error {
	_, err := c.del(ctx, "scv/"+url.PathEscape(alias),
		nil, nil, false, "")
	return err
}
