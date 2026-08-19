package opennms

// Status REST API v2 – /api/v2/status.

import (
	"context"
	"net/url"
	"strconv"
)

// StatusListOptions carries the pagination, ordering, and severity
// filtering parameters shared by the status list endpoints. Zero
// values are omitted from the query (matching the Python wrapper,
// which sends no parameters by default); a nil *StatusListOptions
// sends none at all.
type StatusListOptions struct {
	// Limit is the maximum number of results. 0 = omit.
	Limit int
	// Offset is the result offset for pagination. 0 = omit.
	Offset int
	// OrderBy is the property to order by.
	OrderBy string
	// Order is the sort order, "asc" or "desc".
	Order string
	// SeverityFilter restricts results to entries with this severity.
	SeverityFilter string
}

// statusListParams builds the shared query parameters for status list
// calls.
func statusListParams(opts *StatusListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		return params
	}
	if opts.Limit != 0 {
		params.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Offset != 0 {
		params.Set("offset", strconv.Itoa(opts.Offset))
	}
	if opts.OrderBy != "" {
		params.Set("orderBy", opts.OrderBy)
	}
	if opts.Order != "" {
		params.Set("order", opts.Order)
	}
	if opts.SeverityFilter != "" {
		params.Set("severityFilter", opts.SeverityFilter)
	}
	return params
}

// statusType defaults an empty status computation strategy to
// "alarms".
func statusType(t string) string {
	if t == "" {
		return "alarms"
	}
	return t
}

// GetStatusSummaryNodes returns a severity-count summary of nodes as
// a list of [severityLabel, count] pairs.
//
// typ is the status computation strategy — "alarms" or "outages";
// "" defaults to "alarms".
func (c *Client) GetStatusSummaryNodes(ctx context.Context, typ string) ([]any, error) {
	return c.getList(ctx, "status/summary/nodes/"+url.PathEscape(statusType(typ)), nil, true)
}

// GetStatusSummaryApplications returns a severity-count summary of
// all applications as a list of [severityLabel, count] pairs.
func (c *Client) GetStatusSummaryApplications(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "status/summary/applications", nil, true)
}

// GetStatusSummaryBusinessServices returns a severity-count summary
// of all business services as a list of [severityLabel, count] pairs.
func (c *Client) GetStatusSummaryBusinessServices(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "status/summary/business-services", nil, true)
}

// GetStatusNodes returns nodes with their computed status.
//
// typ is the status computation strategy — "alarms" or "outages";
// "" defaults to "alarms".
func (c *Client) GetStatusNodes(ctx context.Context, typ string, opts *StatusListOptions) (map[string]any, error) {
	return c.getObject(ctx, "status/nodes/"+url.PathEscape(statusType(typ)),
		statusListParams(opts), true)
}

// GetStatusApplications returns applications with their computed
// status.
func (c *Client) GetStatusApplications(ctx context.Context, opts *StatusListOptions) (map[string]any, error) {
	return c.getObject(ctx, "status/applications", statusListParams(opts), true)
}

// GetStatusBusinessServices returns business services with their
// computed status.
func (c *Client) GetStatusBusinessServices(ctx context.Context, opts *StatusListOptions) (map[string]any, error) {
	return c.getObject(ctx, "status/business-services", statusListParams(opts), true)
}
