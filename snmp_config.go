package opennms

// SNMP Configuration REST API – /rest/snmpConfig.

import (
	"context"
	"net/url"
)

// GetSnmpConfig returns the effective SNMP configuration for
// ipAddress. location optionally names a monitoring location; ""
// omits it.
func (c *Client) GetSnmpConfig(ctx context.Context, ipAddress, location string) (map[string]any, error) {
	var params url.Values
	if location != "" {
		params = url.Values{"location": {location}}
	}
	return c.getObject(ctx, "snmpConfig/"+url.PathEscape(ipAddress), params, false)
}

// SetSnmpConfig adds or updates the SNMP configuration for
// ipAddress. config holds the SNMP configuration fields (e.g.
// "version", "community", "port", or the v3 security settings).
func (c *Client) SetSnmpConfig(ctx context.Context, ipAddress string, config map[string]any) error {
	_, err := c.put(ctx, "snmpConfig/"+url.PathEscape(ipAddress), config, nil, false)
	return err
}
