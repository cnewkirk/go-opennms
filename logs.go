package opennms

// Logs REST API – /rest/logs.

import (
	"context"
	"net/url"
	"strconv"
)

// GetLogFiles returns the names of the available OpenNMS log files
// (e.g. ["manager.log", ...]). Requires the ADMIN role.
func (c *Client) GetLogFiles(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "logs", nil, false)
}

// GetLogContents returns the contents of a single OpenNMS log file
// as raw text (the server labels the response application/json but
// sends plain log lines). Requires the ADMIN role.
//
// filename must end in ".log". lines is the maximum number of lines
// to return (server default 5000, range 1–10000); 0 sends no limit.
// reverse returns the most recent lines first (server default true);
// nil leaves the server default.
//
// Returns an empty string when the file does not exist (204).
func (c *Client) GetLogContents(ctx context.Context, filename string, lines int, reverse *bool) (string, error) {
	params := url.Values{"f": {filename}}
	if lines != 0 {
		params.Set("n", strconv.Itoa(lines))
	}
	if reverse != nil {
		// Python-style bool serialization ("True"/"False"), matching
		// the live-validated reference wrapper.
		if *reverse {
			params.Set("reverse", "True")
		} else {
			params.Set("reverse", "False")
		}
	}
	return c.getText(ctx, "logs/contents", params, false, "")
}
