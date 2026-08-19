package opennms

// Flow REST API – /rest/flows (read-only).

import (
	"context"
	"net/url"
	"strconv"
)

// FlowOptions carries the query parameters shared by the flow
// endpoints. Zero-valued fields fall back to the Python wrapper's
// defaults, so a nil *FlowOptions requests the last four hours:
//
//   - Start: start time in ms since epoch; negative values are
//     relative to End (default -14400000 = 4 hours before End).
//   - End: end time in ms since epoch; 0 means now (the default).
//   - TopN: number of top entries to return (default 10); top-N and
//     series endpoints only.
//   - Step: requested interval between data points in ms (default
//     300000 = 5 min); series endpoints only.
//   - Limit: max number of results (default 10); enumerate
//     endpoints only.
//   - IfIndex: SNMP ifIndex to filter by interface; 0 = no filter.
//   - ExporterNode: node criteria (DB ID or
//     "foreignSource:foreignId") to filter by exporter; "" = none.
//   - IncludeOther: include an aggregated "Other" category for
//     traffic outside the top N.
type FlowOptions struct {
	TopN         int
	Start        int
	End          int
	Step         int
	Limit        int
	IfIndex      int
	ExporterNode string
	IncludeOther bool
}

// flowParams builds the common query parameter set shared by all
// flow API methods (start, end, ifIndex, exporterNode).
func flowParams(opts *FlowOptions) url.Values {
	if opts == nil {
		opts = &FlowOptions{}
	}
	start := opts.Start
	if start == 0 {
		start = -14400000
	}
	p := url.Values{}
	p.Set("start", strconv.Itoa(start))
	p.Set("end", strconv.Itoa(opts.End))
	if opts.IfIndex != 0 {
		p.Set("ifIndex", strconv.Itoa(opts.IfIndex))
	}
	if opts.ExporterNode != "" {
		p.Set("exporterNode", opts.ExporterNode)
	}
	return p
}

// flowTopNParams extends flowParams with N and includeOther.
func flowTopNParams(opts *FlowOptions) url.Values {
	if opts == nil {
		opts = &FlowOptions{}
	}
	p := flowParams(opts)
	topN := opts.TopN
	if topN == 0 {
		topN = 10
	}
	p.Set("N", strconv.Itoa(topN))
	p.Set("includeOther", strconv.FormatBool(opts.IncludeOther))
	return p
}

// flowSeriesParams extends flowTopNParams with step.
func flowSeriesParams(opts *FlowOptions) url.Values {
	if opts == nil {
		opts = &FlowOptions{}
	}
	p := flowTopNParams(opts)
	step := opts.Step
	if step == 0 {
		step = 300000
	}
	p.Set("step", strconv.Itoa(step))
	return p
}

// flowEnumerateParams extends flowParams with limit.
func flowEnumerateParams(opts *FlowOptions) url.Values {
	if opts == nil {
		opts = &FlowOptions{}
	}
	p := flowParams(opts)
	limit := opts.Limit
	if limit == 0 {
		limit = 10
	}
	p.Set("limit", strconv.Itoa(limit))
	return p
}

// GetFlowCount returns the number of flows available.
func (c *Client) GetFlowCount(ctx context.Context) (int, error) {
	return c.getCount(ctx, "flows/count", false)
}

// GetFlowExporters returns basic information for all exporter nodes
// that have flows.
func (c *Client) GetFlowExporters(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "flows/exporters", nil, false)
}

// GetFlowExporter returns details about a specific flow exporter
// node. nodeCriteria is a node DB ID or "foreignSource:foreignId".
func (c *Client) GetFlowExporter(ctx context.Context, nodeCriteria string) (map[string]any, error) {
	return c.getObject(ctx, "flows/exporters/"+url.PathEscape(nodeCriteria), nil, false)
}

// GetFlowApplications returns traffic stats for the top TopN
// applications.
func (c *Client) GetFlowApplications(ctx context.Context, opts *FlowOptions) (map[string]any, error) {
	return c.getObject(ctx, "flows/applications", flowTopNParams(opts), false)
}

// GetFlowApplicationsEnumerate lists application names that have
// flows in the given time range.
func (c *Client) GetFlowApplicationsEnumerate(ctx context.Context, opts *FlowOptions) (map[string]any, error) {
	return c.getObject(ctx, "flows/applications/enumerate", flowEnumerateParams(opts), false)
}

// GetFlowApplicationsSeries returns time-series data for the top
// TopN applications.
func (c *Client) GetFlowApplicationsSeries(ctx context.Context, opts *FlowOptions) (map[string]any, error) {
	return c.getObject(ctx, "flows/applications/series", flowSeriesParams(opts), false)
}

// GetFlowConversations returns traffic stats for the top TopN
// conversations.
func (c *Client) GetFlowConversations(ctx context.Context, opts *FlowOptions) (map[string]any, error) {
	return c.getObject(ctx, "flows/conversations", flowTopNParams(opts), false)
}

// GetFlowConversationsEnumerate lists conversations that have flows
// in the given time range.
func (c *Client) GetFlowConversationsEnumerate(ctx context.Context, opts *FlowOptions) (map[string]any, error) {
	return c.getObject(ctx, "flows/conversations/enumerate", flowEnumerateParams(opts), false)
}

// GetFlowConversationsSeries returns time-series data for the top
// TopN conversations.
func (c *Client) GetFlowConversationsSeries(ctx context.Context, opts *FlowOptions) (map[string]any, error) {
	return c.getObject(ctx, "flows/conversations/series", flowSeriesParams(opts), false)
}

// GetFlowHosts returns traffic stats for the top TopN hosts.
func (c *Client) GetFlowHosts(ctx context.Context, opts *FlowOptions) (map[string]any, error) {
	return c.getObject(ctx, "flows/hosts", flowTopNParams(opts), false)
}

// GetFlowHostsEnumerate lists hosts that have flows in the given
// time range.
func (c *Client) GetFlowHostsEnumerate(ctx context.Context, opts *FlowOptions) (map[string]any, error) {
	return c.getObject(ctx, "flows/hosts/enumerate", flowEnumerateParams(opts), false)
}

// GetFlowHostsSeries returns time-series data for the top TopN
// hosts.
func (c *Client) GetFlowHostsSeries(ctx context.Context, opts *FlowOptions) (map[string]any, error) {
	return c.getObject(ctx, "flows/hosts/series", flowSeriesParams(opts), false)
}

// GetFlowDscp returns traffic stats for the top TopN DSCP values.
func (c *Client) GetFlowDscp(ctx context.Context, opts *FlowOptions) (map[string]any, error) {
	return c.getObject(ctx, "flows/dscp", flowTopNParams(opts), false)
}

// GetFlowDscpEnumerate lists DSCP values that have flows in the
// given time range.
func (c *Client) GetFlowDscpEnumerate(ctx context.Context, opts *FlowOptions) (map[string]any, error) {
	return c.getObject(ctx, "flows/dscp/enumerate", flowEnumerateParams(opts), false)
}

// GetFlowDscpSeries returns time-series data for the top TopN DSCP
// values.
func (c *Client) GetFlowDscpSeries(ctx context.Context, opts *FlowOptions) (map[string]any, error) {
	return c.getObject(ctx, "flows/dscp/series", flowSeriesParams(opts), false)
}

// GetFlowGraphUrl returns the configured flow graph URL.
func (c *Client) GetFlowGraphUrl(ctx context.Context) (string, error) {
	return asString(c.get(ctx, "flows/flowGraphUrl", nil, false))
}
