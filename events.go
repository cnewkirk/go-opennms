package opennms

// Events REST API – /rest/events.

import (
	"context"
	"fmt"
	"net/url"
)

// GetEvents lists events (v1).
//
// filters are additional Hibernate query filters passed directly as
// query parameters (e.g. "severity": "MAJOR").
func (c *Client) GetEvents(ctx context.Context, opts *ListOptions, filters map[string]string) (map[string]any, error) {
	params := mergeFilters(listParams(opts), filters)
	return c.getObject(ctx, "events", params, false)
}

// GetEvent returns the event with the given ID.
func (c *Client) GetEvent(ctx context.Context, eventID int) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("events/%d", eventID), nil, false)
}

// GetEventCount returns the total number of events.
func (c *Client) GetEventCount(ctx context.Context) (int, error) {
	return c.getCount(ctx, "events/count", false)
}

// CreateEvent publishes an event to the OpenNMS event bus.
//
// event is a map matching the OpenNMS event schema, e.g.:
//
//	client.CreateEvent(ctx, map[string]any{
//		"uei":       "uei.opennms.org/internal/test",
//		"source":    "my-script",
//		"severity":  "Normal",
//		"nodeId":    1,
//		"interface": "192.168.0.1",
//		"parms": []any{
//			map[string]any{"parmName": "key", "value": "val"},
//		},
//	})
func (c *Client) CreateEvent(ctx context.Context, event map[string]any) error {
	_, err := c.post(ctx, "events", event, nil, false)
	return err
}

// AckEvent acknowledges the event.
func (c *Client) AckEvent(ctx context.Context, eventID int) error {
	_, err := c.putForm(ctx, fmt.Sprintf("events/%d", eventID),
		url.Values{"ack": {"true"}}, nil, false)
	return err
}

// UnackEvent removes acknowledgement from the event.
func (c *Client) UnackEvent(ctx context.Context, eventID int) error {
	_, err := c.putForm(ctx, fmt.Sprintf("events/%d", eventID),
		url.Values{"ack": {"false"}}, nil, false)
	return err
}

// BulkAckEvents acknowledges all events matching the given filters
// (Hibernate query filters passed as form fields, e.g.
// "severity": "MAJOR").
func (c *Client) BulkAckEvents(ctx context.Context, filters map[string]string) error {
	form := url.Values{"ack": {"true"}}
	for k, v := range filters {
		form.Set(k, v)
	}
	_, err := c.putForm(ctx, "events", form, nil, false)
	return err
}

// BulkUnackEvents removes acknowledgement from all events matching
// the given filters.
func (c *Client) BulkUnackEvents(ctx context.Context, filters map[string]string) error {
	form := url.Values{"ack": {"false"}}
	for k, v := range filters {
		form.Set(k, v)
	}
	_, err := c.putForm(ctx, "events", form, nil, false)
	return err
}
