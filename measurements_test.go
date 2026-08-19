package opennms

// Tests for the measurements methods – /rest/measurements.

import (
	"encoding/json"
	"testing"
)

const measurementsResourceID = "node[1].interfaceSnmp[eth0-04013f75f101]"
const measurementsAttribute = "ifInOctets"

const measurementsJSON = `{
	"start": 1425580938256,
	"end": 1425588138256,
	"step": 300000,
	"timestamps": [1425581100000, 1425581400000, 1425581700000],
	"labels": ["ifInOctets"],
	"columns": [
		{"values": [139948.5, 199006.3, 187234.8]}
	],
	"metadata": {
		"resources": [
			{
				"id": "node[1].interfaceSnmp[eth0-04013f75f101]",
				"label": "eth0 (04:01:3f:75:f1:01)"
			}
		],
		"nodes": [
			{"id": 1, "label": "router01.example.com"}
		]
	}
}`

func TestGetMeasurementsDefaults(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/measurements/"+
		measurementsResourceID+"/"+measurementsAttribute, measurementsJSON)
	result, err := c.GetMeasurements(t.Context(),
		measurementsResourceID, measurementsAttribute, nil)
	if err != nil {
		t.Fatal(err)
	}
	if result["labels"].([]any)[0] != "ifInOctets" {
		t.Errorf("labels = %v", result["labels"])
	}
	if len(result["timestamps"].([]any)) != 3 {
		t.Errorf("timestamps = %v", result["timestamps"])
	}
	for param, want := range map[string]string{
		"start": "-14400000", "end": "0", "step": "300000",
		"maxrows": "0", "aggregation": "AVERAGE",
	} {
		if got := req.query.Get(param); got != want {
			t.Errorf("%s = %q, want %q", param, got, want)
		}
	}
}

func TestGetMeasurementsCustomParams(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/measurements/"+
		measurementsResourceID+"/"+measurementsAttribute, measurementsJSON)
	_, err := c.GetMeasurements(t.Context(),
		measurementsResourceID, measurementsAttribute,
		&MeasurementsOptions{
			Start:             1425580938256,
			End:               1425588138256,
			Step:              60000,
			MaxRows:           800,
			Aggregation:       "MAX",
			FallbackAttribute: "ifHCInOctets",
		})
	if err != nil {
		t.Fatal(err)
	}
	for param, want := range map[string]string{
		"start": "1425580938256", "step": "60000", "maxrows": "800",
		"aggregation": "MAX", "fallback-attribute": "ifHCInOctets",
	} {
		if got := req.query.Get(param); got != want {
			t.Errorf("%s = %q, want %q", param, got, want)
		}
	}
}

func TestGetMeasurementsMulti(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "POST "+v1Path+"/measurements", measurementsJSON)
	query := map[string]any{
		"start": 1425580938256,
		"end":   1425588138256,
		"step":  300000,
		"source": []any{
			map[string]any{
				"resourceId":  measurementsResourceID,
				"attribute":   "ifInOctets",
				"label":       "octetsIn",
				"aggregation": "AVERAGE",
				"transient":   false,
			},
			map[string]any{
				"resourceId":  measurementsResourceID,
				"attribute":   "ifOutOctets",
				"label":       "octetsOut",
				"aggregation": "AVERAGE",
				"transient":   false,
			},
		},
		"expression": []any{
			map[string]any{
				"label":     "bitsIn",
				"value":     "octetsIn * 8",
				"transient": false,
			},
		},
	}
	result, err := c.GetMeasurementsMulti(t.Context(), query)
	if err != nil {
		t.Fatal(err)
	}
	if result["labels"].([]any)[0] != "ifInOctets" {
		t.Errorf("labels = %v", result["labels"])
	}
	if req.method != "POST" {
		t.Errorf("method = %q, want POST", req.method)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if len(body["source"].([]any)) != 2 {
		t.Errorf("source = %v", body["source"])
	}
	expr := body["expression"].([]any)[0].(map[string]any)
	if expr["label"] != "bitsIn" {
		t.Errorf("expression label = %v, want bitsIn", expr["label"])
	}
	if got := req.header.Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
}
