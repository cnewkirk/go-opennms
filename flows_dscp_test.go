package opennms

// Tests for the flows DSCP and flowGraphUrl methods.

import (
	"testing"
)

const flowDscpJSON = `{
	"start": 1425580938256,
	"end": 1425588138256,
	"headers": [
		{"dscp": 0, "bytesIn": 1048576, "bytesOut": 204800},
		{"dscp": 46, "bytesIn": 524288, "bytesOut": 102400}
	]
}`

const flowDscpEnumerateJSON = `{"dscp": [0, 8, 46]}`

const flowDscpSeriesJSON = `{
	"start": 1425580938256,
	"end": 1425588138256,
	"step": 300000,
	"timestamps": [1425581100000, 1425581400000],
	"columns": [
		{
			"label": "0",
			"ingress": {"values": [139948.5, 199006.3]},
			"egress": {"values": [51661.2, 64741.0]}
		}
	]
}`

func TestGetFlowDscpDefaults(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/flows/dscp", flowDscpJSON)
	result, err := c.GetFlowDscp(t.Context(), nil)
	if err != nil {
		t.Fatal(err)
	}
	header := result["headers"].([]any)[0].(map[string]any)
	if header["dscp"].(float64) != 0 {
		t.Errorf("dscp = %v, want 0", header["dscp"])
	}
	if got := req.query.Get("N"); got != "10" {
		t.Errorf("N = %q, want 10", got)
	}
	if got := req.query.Get("includeOther"); got != "false" {
		t.Errorf("includeOther = %q, want false", got)
	}
}

func TestGetFlowDscpWithFilters(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/flows/dscp", flowDscpJSON)
	_, err := c.GetFlowDscp(t.Context(), &FlowOptions{
		TopN: 5, IfIndex: 6, ExporterNode: "1", IncludeOther: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	for param, want := range map[string]string{
		"N": "5", "ifIndex": "6", "exporterNode": "1", "includeOther": "true",
	} {
		if got := req.query.Get(param); got != want {
			t.Errorf("%s = %q, want %q", param, got, want)
		}
	}
}

func TestGetFlowDscpEnumerate(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/flows/dscp/enumerate",
		flowDscpEnumerateJSON)
	result, err := c.GetFlowDscpEnumerate(t.Context(), nil)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, dscp := range result["dscp"].([]any) {
		if dscp.(float64) == 0 {
			found = true
		}
	}
	if !found {
		t.Errorf("dscp = %v, want to contain 0", result["dscp"])
	}
	if got := req.query.Get("limit"); got != "10" {
		t.Errorf("limit = %q, want 10", got)
	}
}

func TestGetFlowDscpSeries(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/flows/dscp/series", flowDscpSeriesJSON)
	result, err := c.GetFlowDscpSeries(t.Context(), &FlowOptions{TopN: 3, Step: 60000})
	if err != nil {
		t.Fatal(err)
	}
	column := result["columns"].([]any)[0].(map[string]any)
	if column["label"] != "0" {
		t.Errorf("label = %v, want 0", column["label"])
	}
	if got := req.query.Get("N"); got != "3" {
		t.Errorf("N = %q, want 3", got)
	}
	if got := req.query.Get("step"); got != "60000" {
		t.Errorf("step = %q, want 60000", got)
	}
}

func TestGetFlowGraphUrl(t *testing.T) {
	mux, c := newTestClient(t)
	handleText(mux, "GET "+v1Path+"/flows/flowGraphUrl",
		"http://grafana:3000/d/flows")
	result, err := c.GetFlowGraphUrl(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result != "http://grafana:3000/d/flows" {
		t.Errorf("result = %q", result)
	}
}
