package opennms

// Acknowledgements REST API – /rest/acks.

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// GetAcks lists acknowledgements.
//
// filters are additional Hibernate query filters passed directly as
// query parameters (e.g. "ackUser": "admin").
func (c *Client) GetAcks(ctx context.Context, opts *ListOptions, filters map[string]string) (map[string]any, error) {
	params := mergeFilters(listParams(opts), filters)
	return c.getObject(ctx, "acks", params, false)
}

// GetAck returns the acknowledgement with the given ID.
func (c *Client) GetAck(ctx context.Context, ackID int) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("acks/%d", ackID), nil, false)
}

// GetAckCount returns the total number of acknowledgements.
func (c *Client) GetAckCount(ctx context.Context) (int, error) {
	return c.getCount(ctx, "acks/count", false)
}

// Write (form-encoded POST — the acks endpoint rejects JSON on every
// OpenNMS version)

// CreateAck creates or modifies an acknowledgement.
//
// action is one of "ack", "unack", "clear", "esc". alarmID and
// notificationID (0 = absent) are the mutually exclusive targets.
func (c *Client) CreateAck(ctx context.Context, action string, alarmID, notificationID int) (map[string]any, error) {
	form := url.Values{"action": {action}}
	if alarmID != 0 {
		form.Set("alarmId", strconv.Itoa(alarmID))
	}
	if notificationID != 0 {
		form.Set("notifId", strconv.Itoa(notificationID))
	}
	return asObject(c.postForm(ctx, "acks", form, nil, false))
}

// Convenience wrappers

// AckNotification acknowledges the given notification.
func (c *Client) AckNotification(ctx context.Context, notificationID int) (map[string]any, error) {
	return c.CreateAck(ctx, "ack", 0, notificationID)
}

// UnackNotification removes acknowledgement from the given
// notification.
func (c *Client) UnackNotification(ctx context.Context, notificationID int) (map[string]any, error) {
	return c.CreateAck(ctx, "unack", 0, notificationID)
}
