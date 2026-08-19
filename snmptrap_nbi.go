package opennms

// SNMP Trap NBI Configuration REST API – /rest/config/snmptrap-nbi.

import (
	"context"
	"net/url"
	"strconv"
)

const snmptrapNbiBase = "config/snmptrap-nbi"

// GetSnmptrapNbiConfig returns the full SNMP trap NBI configuration.
func (c *Client) GetSnmptrapNbiConfig(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, snmptrapNbiBase, nil, false)
}

// GetSnmptrapNbiStatus returns the SNMP trap NBI forwarding status.
func (c *Client) GetSnmptrapNbiStatus(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, snmptrapNbiBase+"/status", nil, false)
}

// SetSnmptrapNbiStatus enables or disables SNMP trap NBI forwarding.
func (c *Client) SetSnmptrapNbiStatus(ctx context.Context, enabled bool) error {
	_, err := c.putForm(ctx, snmptrapNbiBase+"/status",
		url.Values{"enabled": {strconv.FormatBool(enabled)}}, nil, false)
	return err
}

// GetSnmptrapNbiTrapsinks lists all SNMP trap NBI trap sinks.
func (c *Client) GetSnmptrapNbiTrapsinks(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, snmptrapNbiBase+"/trapsinks", nil, false)
}

// GetSnmptrapNbiTrapsink returns the trap sink with the given name.
func (c *Client) GetSnmptrapNbiTrapsink(ctx context.Context, name string) (map[string]any, error) {
	return c.getObject(ctx,
		snmptrapNbiBase+"/trapsinks/"+url.PathEscape(name), nil, false)
}

// CreateSnmptrapNbiTrapsink creates a new SNMP trap NBI trap sink.
// data is a trap sink definition with keys such as "name",
// "ipAddress", "port", "community".
func (c *Client) CreateSnmptrapNbiTrapsink(ctx context.Context, data map[string]any) error {
	_, err := c.post(ctx, snmptrapNbiBase+"/trapsinks", data, nil, false)
	return err
}

// UpdateSnmptrapNbiTrapsink updates a trap sink with form-encoded
// key/value pairs.
func (c *Client) UpdateSnmptrapNbiTrapsink(ctx context.Context, name string, data map[string]string) error {
	form := url.Values{}
	for k, v := range data {
		form.Set(k, v)
	}
	_, err := c.putForm(ctx,
		snmptrapNbiBase+"/trapsinks/"+url.PathEscape(name), form, nil, false)
	return err
}

// DeleteSnmptrapNbiTrapsink deletes a trap sink.
func (c *Client) DeleteSnmptrapNbiTrapsink(ctx context.Context, name string) error {
	_, err := c.del(ctx,
		snmptrapNbiBase+"/trapsinks/"+url.PathEscape(name), nil, nil, false, "")
	return err
}

// UpdateSnmptrapNbiConfig updates the SNMP trap NBI configuration.
func (c *Client) UpdateSnmptrapNbiConfig(ctx context.Context, data map[string]any) error {
	_, err := c.post(ctx, snmptrapNbiBase, data, nil, false)
	return err
}
