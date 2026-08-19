package opennms

// Monitored Services (ifservices) REST API – /rest/ifservices + v2.

import (
	"context"
	"net/url"
)

// GetIfservices lists monitored services. params are optional query
// parameters such as "limit", "offset", "node.label",
// "ipInterface.ipAddress" (nil allowed).
func (c *Client) GetIfservices(ctx context.Context, params map[string]string) (map[string]any, error) {
	return c.getObject(ctx, "ifservices", mergeFilters(nil, params), false)
}

// UpdateIfservices bulk-updates monitored services. fields are
// form-encoded key/value pairs to update, e.g. "status": "A" plus
// filter parameters.
func (c *Client) UpdateIfservices(ctx context.Context, fields map[string]string) error {
	form := url.Values{}
	for k, v := range fields {
		form.Set(k, v)
	}
	_, err := c.putForm(ctx, "ifservices", form, nil, false)
	return err
}

// GetIfservicesV2 lists monitored services via the v2 API with an
// optional FIQL filter expression ("" lists all).
func (c *Client) GetIfservicesV2(ctx context.Context, fiql string, opts *ListOptions) (map[string]any, error) {
	params := listParams(opts)
	if fiql != "" {
		params.Set("_s", fiql)
	}
	return c.getObject(ctx, "ifservices", params, true)
}
