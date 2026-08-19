package opennms

// Alarm Statistics REST API – /rest/stats/alarms.

import (
	"context"
	"net/url"
	"strings"
)

// GetAlarmStats returns alarm statistics.
//
// filters are additional Hibernate query filters passed directly as
// query parameters (e.g. "severity": "MAJOR"). The same filter
// parameters as GetAlarms are supported.
func (c *Client) GetAlarmStats(ctx context.Context, filters map[string]string) (map[string]any, error) {
	return c.getObject(ctx, "stats/alarms", mergeFilters(nil, filters), false)
}

// GetAlarmStatsBySeverity returns alarm statistics grouped by
// severity. severities optionally restricts the result to the given
// severity strings, e.g. []string{"MAJOR", "CRITICAL"}; nil includes
// all severities.
func (c *Client) GetAlarmStatsBySeverity(ctx context.Context, severities []string) ([]any, error) {
	params := url.Values{}
	if len(severities) > 0 {
		params.Set("severities", strings.Join(severities, ","))
	}
	return c.getList(ctx, "stats/alarms/by-severity", params, false)
}
