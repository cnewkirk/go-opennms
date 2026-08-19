package opennms

// KSC Reports REST API – /rest/ksc.

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

var kscGraphAttrs = []string{"title", "timespan", "graphtype",
	"resourceId", "nodeId", "nodeSource", "domain", "interfaceId",
	"extlink"}

var kscReportAttrs = []string{"id", "label", "show_timespan_button",
	"show_graphtype_button", "graphs_per_line"}

// kscQuoteAttr quotes a value as an XML attribute (the equivalent of
// Python's xml.sax.saxutils.quoteattr).
func kscQuoteAttr(s string) string {
	r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;",
		`"`, "&quot;")
	return `"` + r.Replace(s) + `"`
}

// kscAttrs serializes the given keys of data as XML attributes, in
// order, skipping keys that are absent.
func kscAttrs(data map[string]any, keys []string) string {
	var b strings.Builder
	for _, key := range keys {
		if value, ok := data[key]; ok {
			b.WriteString(" " + key + "=" + kscQuoteAttr(fmt.Sprintf("%v", value)))
		}
	}
	return b.String()
}

// kscGraphs normalizes the "graphs" entry of a report into a slice of
// graph maps.
func kscGraphs(v any) []map[string]any {
	switch graphs := v.(type) {
	case []map[string]any:
		return graphs
	case []any:
		out := make([]map[string]any, 0, len(graphs))
		for _, g := range graphs {
			if m, ok := g.(map[string]any); ok {
				out = append(out, m)
			}
		}
		return out
	}
	return nil
}

// GetKscReports lists all KSC reports (returns ID and label for each).
func (c *Client) GetKscReports(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "ksc", nil, false)
}

// GetKscReport returns a specific KSC report by ID.
func (c *Client) GetKscReport(ctx context.Context, reportID int) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("ksc/%d", reportID), nil, false)
}

// GetKscReportCount returns the total number of KSC reports.
func (c *Client) GetKscReportCount(ctx context.Context) (int, error) {
	return c.getCount(ctx, "ksc/count", false)
}

// CreateKscReport creates a new KSC report.
//
// The body is sent as XML — the KSC API documents "Create a report
// from an XML payload" with application/xml.
//
// report is the KSC report definition. Example:
//
//	map[string]any{
//		"id":                    0,
//		"label":                 "My Bandwidth Report",
//		"show_timespan_button":  false,
//		"show_graphtype_button": false,
//		"graphs_per_line":       1,
//		"graphs": []map[string]any{{
//			"title":      "Core Switch Bandwidth",
//			"resourceId": "node[1].interfaceSnmp[eth0-04013f75f101]",
//			"timespan":   "7_day",
//			"graphtype":  "mib2.bits",
//		}},
//	}
func (c *Client) CreateKscReport(ctx context.Context, report map[string]any) error {
	var graphs strings.Builder
	for _, graph := range kscGraphs(report["graphs"]) {
		graphs.WriteString("<kscGraph" + kscAttrs(graph, kscGraphAttrs) + "/>")
	}
	xml := "<kscReport" + kscAttrs(report, kscReportAttrs) + ">" +
		graphs.String() + "</kscReport>"
	_, err := c.postText(ctx, "ksc", xml, "application/xml", false, "", nil)
	return err
}

// AddGraphToKscReport adds a graph to an existing KSC report.
//
// Mirrors PUT /rest/ksc/{reportid}, which the KSC API documents as
// "Add a graph to the existing report with the given ID", built from
// query parameters.
//
// reportID is the database ID of the KSC report. reportName is the
// graph definition's report.name from snmp-graph.properties.d.
// resourceID is the time-series resource ID to graph. title is an
// optional graph title; timespan is an optional timespan (server
// default "7_day"). Pass "" to omit either.
func (c *Client) AddGraphToKscReport(ctx context.Context, reportID int, reportName, resourceID, title, timespan string) error {
	params := url.Values{
		"reportName": {reportName},
		"resourceId": {resourceID},
	}
	if title != "" {
		params.Set("title", title)
	}
	if timespan != "" {
		params.Set("timespan", timespan)
	}
	_, err := c.put(ctx, fmt.Sprintf("ksc/%d", reportID), nil, params, false)
	return err
}
