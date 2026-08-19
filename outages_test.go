package opennms

// Tests for the outages methods – /rest/outages.

import (
	"strings"
	"testing"
)

const outageJSON = `{
	"id": 501,
	"ifLostService": "2024-06-01T08:00:00.000+0000",
	"ifRegainedService": null,
	"ipAddress": "192.168.1.1",
	"serviceType": {"id": 1, "name": "ICMP"},
	"monitoredService": {
		"id": 201,
		"serviceName": "ICMP",
		"status": "A",
		"down": true,
		"serviceType": {"id": 1, "name": "ICMP"},
		"ipInterface": {"id": 101, "ipAddress": "192.168.1.1", "nodeId": 1}
	},
	"node": {"id": 1, "label": "router01.example.com"}
}`

const outageListJSON = `{
	"count": 1, "offset": 0, "totalCount": 1,
	"outage": [` + outageJSON + `]
}`

func TestGetOutagesDefault(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/outages", outageListJSON)
	result, err := c.GetOutages(t.Context(), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	outage := result["outage"].([]any)[0].(map[string]any)
	if outage["id"].(float64) != 501 {
		t.Errorf("id = %v, want 501", outage["id"])
	}
	if got := req.query.Get("limit"); got != "10" {
		t.Errorf("limit = %q, want 10", got)
	}
	if got := req.query.Get("offset"); got != "0" {
		t.Errorf("offset = %q, want 0", got)
	}
}

func TestGetOutagesWithFilters(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/outages", outageListJSON)
	_, err := c.GetOutages(t.Context(),
		&ListOptions{Limit: 25, OrderBy: "ifLostService", Order: "descending"},
		map[string]string{"ifRegainedService": "null"})
	if err != nil {
		t.Fatal(err)
	}
	for param, want := range map[string]string{
		"limit": "25", "orderBy": "ifLostService",
		"order": "descending", "ifRegainedService": "null",
	} {
		if got := req.query.Get(param); got != want {
			t.Errorf("%s = %q, want %q", param, got, want)
		}
	}
}

func TestGetOutage(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/outages/501", outageJSON)
	result, err := c.GetOutage(t.Context(), 501)
	if err != nil {
		t.Fatal(err)
	}
	if result["id"].(float64) != 501 {
		t.Errorf("id = %v, want 501", result["id"])
	}
	if result["ifRegainedService"] != nil {
		t.Errorf("ifRegainedService = %v, want nil", result["ifRegainedService"])
	}
}

func TestGetOutageCount(t *testing.T) {
	mux, c := newTestClient(t)
	handleText(mux, "GET "+v1Path+"/outages/count", "17")
	count, err := c.GetOutageCount(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if count != 17 {
		t.Errorf("count = %d, want 17", count)
	}
}

func TestGetNodeOutages(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/outages/forNode/1", outageListJSON)
	result, err := c.GetNodeOutages(t.Context(), 1)
	if err != nil {
		t.Fatal(err)
	}
	outage := result["outage"].([]any)[0].(map[string]any)
	if outage["node"].(map[string]any)["id"].(float64) != 1 {
		t.Errorf("node id = %v, want 1", outage["node"])
	}
	if !strings.HasSuffix(req.path, "/outages/forNode/1") {
		t.Errorf("path = %q", req.path)
	}
}
