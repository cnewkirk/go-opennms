package opennms

// Notifications REST API – /rest/notifications.

import (
	"context"
	"fmt"
	"net/url"
)

// GetNotifications lists notifications.
//
// filters are additional Hibernate query filters passed directly as
// query parameters (e.g. "answered": "false").
func (c *Client) GetNotifications(ctx context.Context, opts *ListOptions, filters map[string]string) (map[string]any, error) {
	params := mergeFilters(listParams(opts), filters)
	return c.getObject(ctx, "notifications", params, false)
}

// GetNotification returns a single notification by ID.
func (c *Client) GetNotification(ctx context.Context, notificationID int) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("notifications/%d", notificationID), nil, false)
}

// GetNotificationCount returns the total number of notifications.
func (c *Client) GetNotificationCount(ctx context.Context) (int, error) {
	return c.getCount(ctx, "notifications/count", false)
}

// TriggerDestinationPath triggers the targets of the named
// destination path for testing.
func (c *Client) TriggerDestinationPath(ctx context.Context, destinationPathName string) error {
	_, err := c.post(ctx, "notifications/destination-paths/"+
		url.PathEscape(destinationPathName)+"/trigger", nil, nil, false)
	return err
}
