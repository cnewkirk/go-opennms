package opennms

// Event Configuration REST API v2 – /api/v2/eventconf.
//
// This API exists in OpenNMS Horizon 35+ only; Horizon 34 and
// earlier (including Meridian 2025) return 404.

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
)

// Filter

// GetEventconfFilter returns event configuration with optional
// filtering. uei is an optional UEI pattern to filter by, vendor an
// optional vendor name; "" means absent. params are additional query
// parameters (nil allowed).
func (c *Client) GetEventconfFilter(ctx context.Context, uei, vendor string, params map[string]string) (map[string]any, error) {
	p := mergeFilters(url.Values{}, params)
	if uei != "" {
		p.Set("uei", uei)
	}
	if vendor != "" {
		p.Set("vendor", vendor)
	}
	return c.getObject(ctx, "eventconf/filter", p, true)
}

// GetEventconfFilterSources returns event configuration sources.
// filter is an optional filter expression ("" = absent); params are
// additional query parameters (nil allowed).
func (c *Client) GetEventconfFilterSources(ctx context.Context, filter string, params map[string]string) (map[string]any, error) {
	p := mergeFilters(url.Values{}, params)
	if filter != "" {
		p.Set("filter", filter)
	}
	return c.getObject(ctx, "eventconf/filter/sources", p, true)
}

// GetEventconfFilterEvents returns events for a specific
// configuration source. params are additional query parameters (nil
// allowed).
func (c *Client) GetEventconfFilterEvents(ctx context.Context, sourceID string, params map[string]string) (map[string]any, error) {
	return c.getObject(ctx,
		fmt.Sprintf("eventconf/filter/%s/events", url.PathEscape(sourceID)),
		mergeFilters(url.Values{}, params), true)
}

// Sources

// GetEventconfSourceNames lists all event configuration source names.
func (c *Client) GetEventconfSourceNames(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "eventconf/sources/names", nil, true)
}

// GetEventconfSource returns a specific event configuration source.
func (c *Client) GetEventconfSource(ctx context.Context, sourceID string) (map[string]any, error) {
	return c.getObject(ctx, "eventconf/sources/"+url.PathEscape(sourceID), nil, true)
}

// DownloadEventconfEvents downloads events for a source as raw XML
// text.
func (c *Client) DownloadEventconfEvents(ctx context.Context, sourceID string) (string, error) {
	return c.getText(ctx,
		fmt.Sprintf("eventconf/sources/%s/events/download", url.PathEscape(sourceID)),
		nil, true, "")
}

// Vendor events

// GetEventconfVendorEvents returns events for a specific vendor.
func (c *Client) GetEventconfVendorEvents(ctx context.Context, vendorName string) (map[string]any, error) {
	return c.getObject(ctx,
		fmt.Sprintf("eventconf/vendors/%s/events", url.PathEscape(vendorName)),
		nil, true)
}

// CRUD

// CreateEventconfEvent creates a new event definition in a source.
func (c *Client) CreateEventconfEvent(ctx context.Context, sourceID string, event map[string]any) (map[string]any, error) {
	return asObject(c.post(ctx,
		fmt.Sprintf("eventconf/sources/%s/events", url.PathEscape(sourceID)),
		event, nil, true))
}

// UpdateEventconfEvent updates an event definition identified by
// eventID within the source.
func (c *Client) UpdateEventconfEvent(ctx context.Context, sourceID, eventID string, event map[string]any) error {
	_, err := c.put(ctx,
		fmt.Sprintf("eventconf/sources/%s/events/%s",
			url.PathEscape(sourceID), url.PathEscape(eventID)),
		event, nil, true)
	return err
}

// Upload

// UploadEventconf uploads raw event configuration bytes as multipart
// form data (filename "events.xml").
func (c *Client) UploadEventconf(ctx context.Context, content []byte) error {
	_, err := c.postFiles(ctx, "eventconf/upload",
		map[string]fileUpload{"file": {filename: "events.xml", content: content}},
		nil, true)
	return err
}

// UploadEventconfFile uploads an event configuration file from disk
// as multipart form data.
func (c *Client) UploadEventconfFile(ctx context.Context, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = c.postFiles(ctx, "eventconf/upload",
		map[string]fileUpload{"file": {filename: filepath.Base(path), content: content}},
		nil, true)
	return err
}

// Status (PATCH)

// SetEventconfSourcesStatus sets the enabled/disabled status of
// event configuration sources. payload maps source IDs to boolean
// status values.
func (c *Client) SetEventconfSourcesStatus(ctx context.Context, payload map[string]any) error {
	_, err := c.patch(ctx, "eventconf/sources/status", payload, nil, true)
	return err
}

// SetEventconfEventsStatus sets the enabled/disabled status of
// events within a source. payload maps event IDs to boolean status
// values.
func (c *Client) SetEventconfEventsStatus(ctx context.Context, sourceID string, payload map[string]any) error {
	_, err := c.patch(ctx,
		fmt.Sprintf("eventconf/sources/%s/events/status", url.PathEscape(sourceID)),
		payload, nil, true)
	return err
}

// Delete

// DeleteEventconfSources deletes event configuration sources.
// payload identifies the sources to delete.
func (c *Client) DeleteEventconfSources(ctx context.Context, payload map[string]any) error {
	_, err := c.del(ctx, "eventconf/sources", nil, payload, true, "")
	return err
}

// DeleteEventconfEvents deletes events from a source. payload
// identifies the events to delete.
func (c *Client) DeleteEventconfEvents(ctx context.Context, sourceID string, payload map[string]any) error {
	_, err := c.del(ctx,
		fmt.Sprintf("eventconf/sources/%s/events", url.PathEscape(sourceID)),
		nil, payload, true, "")
	return err
}
