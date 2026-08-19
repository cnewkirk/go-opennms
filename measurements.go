package opennms

// Measurements REST API – /rest/measurements.

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// MeasurementsOptions carries the query parameters of GetMeasurements.
// Zero values fall back to the server-friendly defaults of the Python
// wrapper: Start -14400000 (4 hours before End), End 0 (now), Step
// 300000 (5 min), MaxRows 0 (no limit), Aggregation "AVERAGE". A nil
// *MeasurementsOptions uses all defaults.
type MeasurementsOptions struct {
	// Start is the start time in milliseconds since epoch. Negative
	// values are relative to End. 0 means the default -14400000.
	Start int
	// End is the end time in ms since epoch. 0 means now.
	End int
	// Step is the requested interval between rows in ms. 0 means the
	// default 300000.
	Step int
	// MaxRows is the maximum number of rows to return. 0 = no limit.
	MaxRows int
	// Aggregation is the consolidation function: "AVERAGE", "MIN", or
	// "MAX". "" means the default "AVERAGE".
	Aggregation string
	// FallbackAttribute is a secondary attribute used when the
	// primary doesn't exist.
	FallbackAttribute string
}

// GetMeasurements retrieves time-series values for a single
// attribute.
//
// resourceID is an OpenNMS resource ID string, e.g.
// "node[1].interfaceSnmp[eth0-04013f75f101]"; attribute is an RRD
// attribute name, e.g. "ifInOctets".
func (c *Client) GetMeasurements(ctx context.Context, resourceID, attribute string, opts *MeasurementsOptions) (map[string]any, error) {
	if opts == nil {
		opts = &MeasurementsOptions{}
	}
	start := opts.Start
	if start == 0 {
		start = -14400000
	}
	step := opts.Step
	if step == 0 {
		step = 300000
	}
	aggregation := opts.Aggregation
	if aggregation == "" {
		aggregation = "AVERAGE"
	}
	params := url.Values{}
	params.Set("start", strconv.Itoa(start))
	params.Set("end", strconv.Itoa(opts.End))
	params.Set("step", strconv.Itoa(step))
	params.Set("maxrows", strconv.Itoa(opts.MaxRows))
	params.Set("aggregation", aggregation)
	if opts.FallbackAttribute != "" {
		params.Set("fallback-attribute", opts.FallbackAttribute)
	}
	return c.getObject(ctx,
		fmt.Sprintf("measurements/%s/%s",
			url.PathEscape(resourceID), url.PathEscape(attribute)),
		params, false)
}

// GetMeasurementsMulti retrieves measurements for multiple attributes
// with JEXL expressions.
//
// query is a map with keys: "start" (int, ms epoch), "end" (int, ms
// epoch), "step" (int, ms interval), "maxrows" (int), "source" (list
// of source maps with keys "resourceId", "attribute", "label",
// "aggregation", and optionally "transient"), "expression" (list of
// maps with "label", "value" as JEXL expression, and "transient").
func (c *Client) GetMeasurementsMulti(ctx context.Context, query map[string]any) (map[string]any, error) {
	return asObject(c.post(ctx, "measurements", query, nil, false))
}
