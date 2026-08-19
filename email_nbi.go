package opennms

// Email NBI Configuration REST API – /rest/config/email-nbi.

import (
	"context"
	"net/url"
)

const emailNbiBase = "config/email-nbi"

// GetEmailNbiConfig returns the full email NBI configuration.
func (c *Client) GetEmailNbiConfig(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, emailNbiBase, nil, false)
}

// GetEmailNbiStatus returns the email NBI forwarding status.
func (c *Client) GetEmailNbiStatus(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, emailNbiBase+"/status", nil, false)
}

// SetEmailNbiStatus enables or disables email NBI forwarding.
func (c *Client) SetEmailNbiStatus(ctx context.Context, enabled bool) error {
	value := "false"
	if enabled {
		value = "true"
	}
	_, err := c.putForm(ctx, emailNbiBase+"/status",
		url.Values{"enabled": {value}}, nil, false)
	return err
}

// GetEmailNbiDestinations lists all email NBI destinations.
func (c *Client) GetEmailNbiDestinations(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, emailNbiBase+"/destinations", nil, false)
}

// GetEmailNbiDestination returns a specific email destination by name.
func (c *Client) GetEmailNbiDestination(ctx context.Context, name string) (map[string]any, error) {
	return c.getObject(ctx, emailNbiBase+"/destinations/"+url.PathEscape(name), nil, false)
}

// CreateEmailNbiDestination creates a new email NBI destination.
//
// data is the destination definition with keys such as "name",
// "firstOccurrenceOnly", "filters".
func (c *Client) CreateEmailNbiDestination(ctx context.Context, data map[string]any) error {
	_, err := c.post(ctx, emailNbiBase+"/destinations", data, nil, false)
	return err
}

// UpdateEmailNbiDestination updates the named email NBI destination
// with form-encoded key/value pairs.
func (c *Client) UpdateEmailNbiDestination(ctx context.Context, name string, data map[string]string) error {
	form := url.Values{}
	for k, v := range data {
		form.Set(k, v)
	}
	_, err := c.putForm(ctx, emailNbiBase+"/destinations/"+url.PathEscape(name),
		form, nil, false)
	return err
}

// DeleteEmailNbiDestination deletes an email NBI destination.
func (c *Client) DeleteEmailNbiDestination(ctx context.Context, name string) error {
	_, err := c.del(ctx, emailNbiBase+"/destinations/"+url.PathEscape(name),
		nil, nil, false, "")
	return err
}

// UpdateEmailNbiConfig updates the email NBI configuration with the
// updated configuration data.
func (c *Client) UpdateEmailNbiConfig(ctx context.Context, data map[string]any) error {
	_, err := c.post(ctx, emailNbiBase, data, nil, false)
	return err
}
