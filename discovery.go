package opennms

// Discovery REST API v2 – /api/v2/discovery.

import "context"

// Discover submits a one-time discovery scan configuration (v2).
//
// config is the discovery configuration. Supported keys (all lists
// default to empty):
//
//   - "specifics" ([]): individual IPs to scan. Each entry:
//     {"ip": "192.168.0.1", "location": "Default", "retries": 1,
//     "timeout": 2000, "foreignSource": "FS"}
//   - "include_ranges" ([]): IP ranges to scan. Each entry:
//     {"begin": "192.168.0.1", "end": "192.168.0.254",
//     "location": "Default", "retries": 1, "timeout": 2000}
//   - "exclude_ranges" ([]): IP ranges to exclude. Each entry:
//     {"begin": "192.168.0.100", "end": "192.168.0.110"}
//   - "include_urls" ([]): URLs with newline-delimited IPs. Each
//     entry: {"url": "http://example.com/ips.txt", "location": "Default"}
//
// Note: the v2 discovery endpoint is documented as XML-only. This
// method sends JSON; if your OpenNMS version rejects it with HTTP 415
// please open an issue — the workaround is to submit the request
// manually with an XML body matching the discovery-configuration
// schema.
func (c *Client) Discover(ctx context.Context, config map[string]any) error {
	_, err := c.post(ctx, "discovery", config, nil, true)
	return err
}
