package opennms

// Syslog NBI Configuration REST API – /rest/config/syslog-nbi.

import (
	"context"
	"fmt"
	"net/url"
)

const syslogNbiPath = "config/syslog-nbi"

// GetSyslogNbiConfig returns the full syslog NBI configuration.
func (c *Client) GetSyslogNbiConfig(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, syslogNbiPath, nil, false)
}

// GetSyslogNbiStatus returns the syslog NBI forwarding status.
func (c *Client) GetSyslogNbiStatus(ctx context.Context) (string, error) {
	return asString(c.get(ctx, syslogNbiPath+"/status", nil, false))
}

// SetSyslogNbiStatus enables (true) or disables (false) syslog NBI
// forwarding.
func (c *Client) SetSyslogNbiStatus(ctx context.Context, enabled bool) error {
	_, err := c.putForm(ctx, syslogNbiPath+"/status",
		url.Values{"enabled": {fmt.Sprintf("%t", enabled)}}, nil, false)
	return err
}

// GetSyslogNbiDestinations lists all syslog NBI destinations.
func (c *Client) GetSyslogNbiDestinations(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, syslogNbiPath+"/destinations", nil, false)
}

// GetSyslogNbiDestination returns a specific syslog destination by
// name.
func (c *Client) GetSyslogNbiDestination(ctx context.Context, name string) (map[string]any, error) {
	return c.getObject(ctx,
		syslogNbiPath+"/destinations/"+url.PathEscape(name), nil, false)
}

// CreateSyslogNbiDestination creates a new syslog NBI destination.
// The destination definition carries keys such as "name", "host",
// "port", "firstOccurrenceOnly", "filters".
func (c *Client) CreateSyslogNbiDestination(ctx context.Context, data map[string]any) error {
	_, err := c.post(ctx, syslogNbiPath+"/destinations", data, nil, false)
	return err
}

// UpdateSyslogNbiDestination updates the named syslog NBI destination
// with the given form-encoded key/value pairs.
func (c *Client) UpdateSyslogNbiDestination(ctx context.Context, name string, data map[string]string) error {
	form := url.Values{}
	for k, v := range data {
		form.Set(k, v)
	}
	_, err := c.putForm(ctx,
		syslogNbiPath+"/destinations/"+url.PathEscape(name), form, nil, false)
	return err
}

// DeleteSyslogNbiDestination deletes a syslog NBI destination.
func (c *Client) DeleteSyslogNbiDestination(ctx context.Context, name string) error {
	_, err := c.del(ctx,
		syslogNbiPath+"/destinations/"+url.PathEscape(name), nil, nil, false, "")
	return err
}

// UpdateSyslogNbiConfig updates the syslog NBI configuration with the
// given configuration data.
func (c *Client) UpdateSyslogNbiConfig(ctx context.Context, data map[string]any) error {
	_, err := c.post(ctx, syslogNbiPath, data, nil, false)
	return err
}
