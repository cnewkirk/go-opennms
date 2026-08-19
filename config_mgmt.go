package opennms

// Configuration Management REST API – /rest/cm.

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// configMgmtEscape percent-encodes a caller-supplied sub-path while
// preserving "/" separators, matching how requests encodes paths.
func configMgmtEscape(path string) string {
	segments := strings.Split(path, "/")
	for i, s := range segments {
		segments[i] = url.PathEscape(s)
	}
	return strings.Join(segments, "/")
}

// GetConfigNames lists all configuration names.
func (c *Client) GetConfigNames(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "cm", nil, false)
}

// GetConfigSchemas lists all configuration schemas.
func (c *Client) GetConfigSchemas(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "cm/schema", nil, false)
}

// GetConfigSchema returns the schema for the named configuration.
func (c *Client) GetConfigSchema(ctx context.Context, name string) (map[string]any, error) {
	return c.getObject(ctx, "cm/schema/"+url.PathEscape(name), nil, false)
}

// GetConfigIds lists configuration IDs for the named configuration.
func (c *Client) GetConfigIds(ctx context.Context, name string) ([]any, error) {
	return c.getList(ctx, "cm/"+url.PathEscape(name), nil, false)
}

// GetConfig returns a specific configuration by name and
// configuration identifier.
func (c *Client) GetConfig(ctx context.Context, name, configID string) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("cm/%s/%s",
		url.PathEscape(name), url.PathEscape(configID)), nil, false)
}

// GetConfigPart returns a sub-part of a configuration, addressed by
// the sub-path within the configuration.
func (c *Client) GetConfigPart(ctx context.Context, name, configID, path string) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("cm/%s/%s/%s",
		url.PathEscape(name), url.PathEscape(configID),
		configMgmtEscape(path)), nil, false)
}

// CreateConfig creates a new configuration with the given name and
// configuration identifier from the configuration data.
func (c *Client) CreateConfig(ctx context.Context, name, configID string, data map[string]any) error {
	_, err := c.post(ctx, fmt.Sprintf("cm/%s/%s",
		url.PathEscape(name), url.PathEscape(configID)), data, nil, false)
	return err
}

// UpdateConfig updates a configuration with the updated
// configuration data.
func (c *Client) UpdateConfig(ctx context.Context, name, configID string, data map[string]any) error {
	_, err := c.put(ctx, fmt.Sprintf("cm/%s/%s",
		url.PathEscape(name), url.PathEscape(configID)), data, nil, false)
	return err
}

// DeleteConfig deletes a configuration by name and configuration
// identifier.
func (c *Client) DeleteConfig(ctx context.Context, name, configID string) error {
	_, err := c.del(ctx, fmt.Sprintf("cm/%s/%s",
		url.PathEscape(name), url.PathEscape(configID)), nil, nil, false, "")
	return err
}

// DeleteConfigPart deletes a sub-part of a configuration, addressed
// by the sub-path within the configuration.
func (c *Client) DeleteConfigPart(ctx context.Context, name, configID, path string) error {
	_, err := c.del(ctx, fmt.Sprintf("cm/%s/%s/%s",
		url.PathEscape(name), url.PathEscape(configID),
		configMgmtEscape(path)), nil, nil, false, "")
	return err
}
