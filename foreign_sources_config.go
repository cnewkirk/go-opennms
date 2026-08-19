package opennms

// Foreign Sources Configuration REST API – /rest/foreignSourcesConfig.

import (
	"context"
	"net/url"
)

// GetForeignSourceConfigPolicies lists available provisioning
// policies.
func (c *Client) GetForeignSourceConfigPolicies(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "foreignSourcesConfig/policies", nil, false)
}

// GetForeignSourceConfigDetectors lists available service detectors.
func (c *Client) GetForeignSourceConfigDetectors(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "foreignSourcesConfig/detectors", nil, false)
}

// GetForeignSourceConfigServices lists services available for the
// given foreign source / provisioning group name.
func (c *Client) GetForeignSourceConfigServices(ctx context.Context, groupName string) (map[string]any, error) {
	return c.getObject(ctx,
		"foreignSourcesConfig/services/"+url.PathEscape(groupName),
		nil, false)
}

// GetForeignSourceConfigAssets lists available asset fields.
func (c *Client) GetForeignSourceConfigAssets(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "foreignSourcesConfig/assets", nil, false)
}

// GetForeignSourceConfigCategories lists available surveillance
// categories.
func (c *Client) GetForeignSourceConfigCategories(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "foreignSourcesConfig/categories", nil, false)
}
