package opennms

// Tests for the flows methods – /rest/flows.

import (
	"testing"
)

const flowExporterJSON = `{
	"node": {
		"id": 1,
		"foreignSource": "Routers",
		"foreignId": "router01",
		"label": "router01.example.com",
		"location": "Default"
	},
	"snmpInterface": {
		"index": 6,
		"name": "GigabitEthernet0/0",
		"speed": 1000000000
	}
}`

const flowExporterListJSON = `[` + flowExporterJSON + `]`

const flowApplicationsJSON = `{
	"start": 1425580938256,
	"end": 1425588138256,
	"applications": [
		{"application": "HTTP", "bytesIn": 1048576, "bytesOut": 204800},
		{"application": "HTTPS", "bytesIn": 5242880, "bytesOut": 1048576}
	]
}`

const flowApplicationsEnumerateJSON = `{"label": ["HTTP", "HTTPS", "SSH"]}`

const flowSeriesJSON = `{
	"start": 1425580938256,
	"end": 1425588138256,
	"step": 300000,
	"timestamps": [1425581100000, 1425581400000],
	"columns": [
		{
			"label": "HTTP",
			"ingress": {"values": [139948.5, 199006.3]},
			"egress": {"values": [51661.2, 64741.0]}
		}
	]
}`

const flowConversationsJSON = `{
	"start": 1425580938256,
	"end": 1425588138256,
	"conversations": [
		{
			"location": "Default",
			"protocol": "TCP",
			"sourceIp": "10.0.0.1",
			"sourcePort": 45123,
			"destIp": "192.168.1.1",
			"destPort": 443,
			"bytesIn": 2097152,
			"bytesOut": 524288
		}
	]
}`

const flowHostsJSON = `{
	"start": 1425580938256,
	"end": 1425588138256,
	"hosts": [
		{"host": "192.168.1.1", "bytesIn": 5242880, "bytesOut": 2097152}
	]
}`

func TestGetFlowCount(t *testing.T) {
	mux, c := newTestClient(t)
	handleText(mux, "GET "+v1Path+"/flows/count", "1048576")
	count, err := c.GetFlowCount(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if count != 1048576 {
		t.Errorf("count = %d, want 1048576", count)
	}
}

func TestGetFlowExporters(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/flows/exporters", flowExporterListJSON)
	result, err := c.GetFlowExporters(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	node := result[0].(map[string]any)["node"].(map[string]any)
	if node["id"].(float64) != 1 {
		t.Errorf("node id = %v, want 1", node["id"])
	}
}

func TestGetFlowExporter(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/flows/exporters/1", flowExporterJSON)
	result, err := c.GetFlowExporter(t.Context(), "1")
	if err != nil {
		t.Fatal(err)
	}
	node := result["node"].(map[string]any)
	if node["foreignSource"] != "Routers" {
		t.Errorf("foreignSource = %v, want Routers", node["foreignSource"])
	}
}

// Applications

func TestGetFlowApplicationsDefaults(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/flows/applications", flowApplicationsJSON)
	result, err := c.GetFlowApplications(t.Context(), nil)
	if err != nil {
		t.Fatal(err)
	}
	app := result["applications"].([]any)[0].(map[string]any)
	if app["application"] != "HTTP" {
		t.Errorf("application = %v, want HTTP", app["application"])
	}
	for param, want := range map[string]string{
		"N": "10", "start": "-14400000", "end": "0", "includeOther": "false",
	} {
		if got := req.query.Get(param); got != want {
			t.Errorf("%s = %q, want %q", param, got, want)
		}
	}
}

func TestGetFlowApplicationsWithFilters(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/flows/applications", flowApplicationsJSON)
	_, err := c.GetFlowApplications(t.Context(), &FlowOptions{
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

func TestGetFlowApplicationsEnumerate(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/flows/applications/enumerate",
		flowApplicationsEnumerateJSON)
	result, err := c.GetFlowApplicationsEnumerate(t.Context(), nil)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, label := range result["label"].([]any) {
		if label == "HTTP" {
			found = true
		}
	}
	if !found {
		t.Errorf("label = %v, want to contain HTTP", result["label"])
	}
	if got := req.query.Get("limit"); got != "10" {
		t.Errorf("limit = %q, want 10", got)
	}
}

func TestGetFlowApplicationsSeries(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/flows/applications/series", flowSeriesJSON)
	result, err := c.GetFlowApplicationsSeries(t.Context(),
		&FlowOptions{TopN: 3, Step: 60000})
	if err != nil {
		t.Fatal(err)
	}
	column := result["columns"].([]any)[0].(map[string]any)
	if column["label"] != "HTTP" {
		t.Errorf("label = %v, want HTTP", column["label"])
	}
	if got := req.query.Get("N"); got != "3" {
		t.Errorf("N = %q, want 3", got)
	}
	if got := req.query.Get("step"); got != "60000" {
		t.Errorf("step = %q, want 60000", got)
	}
}

// Conversations

func TestGetFlowConversationsDefaults(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/flows/conversations", flowConversationsJSON)
	result, err := c.GetFlowConversations(t.Context(), nil)
	if err != nil {
		t.Fatal(err)
	}
	conv := result["conversations"].([]any)[0].(map[string]any)
	if conv["protocol"] != "TCP" {
		t.Errorf("protocol = %v, want TCP", conv["protocol"])
	}
	if got := req.query.Get("N"); got != "10" {
		t.Errorf("N = %q, want 10", got)
	}
}

func TestGetFlowConversationsEnumerate(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/flows/conversations/enumerate",
		`{"conversation": []}`)
	_, err := c.GetFlowConversationsEnumerate(t.Context(), &FlowOptions{Limit: 20})
	if err != nil {
		t.Fatal(err)
	}
	if got := req.query.Get("limit"); got != "20" {
		t.Errorf("limit = %q, want 20", got)
	}
}

func TestGetFlowConversationsSeries(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/flows/conversations/series", flowSeriesJSON)
	_, err := c.GetFlowConversationsSeries(t.Context(), &FlowOptions{TopN: 5})
	if err != nil {
		t.Fatal(err)
	}
	if got := req.query.Get("N"); got != "5" {
		t.Errorf("N = %q, want 5", got)
	}
}

// Hosts

func TestGetFlowHostsDefaults(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/flows/hosts", flowHostsJSON)
	result, err := c.GetFlowHosts(t.Context(), nil)
	if err != nil {
		t.Fatal(err)
	}
	host := result["hosts"].([]any)[0].(map[string]any)
	if host["host"] != "192.168.1.1" {
		t.Errorf("host = %v, want 192.168.1.1", host["host"])
	}
	if got := req.query.Get("N"); got != "10" {
		t.Errorf("N = %q, want 10", got)
	}
}

func TestGetFlowHostsEnumerate(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/flows/hosts/enumerate", `{"host": []}`)
	_, err := c.GetFlowHostsEnumerate(t.Context(), &FlowOptions{Limit: 25})
	if err != nil {
		t.Fatal(err)
	}
	if got := req.query.Get("limit"); got != "25" {
		t.Errorf("limit = %q, want 25", got)
	}
}

func TestGetFlowHostsSeries(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/flows/hosts/series", flowSeriesJSON)
	_, err := c.GetFlowHostsSeries(t.Context(), &FlowOptions{TopN: 3, Step: 300000})
	if err != nil {
		t.Fatal(err)
	}
	if got := req.query.Get("N"); got != "3" {
		t.Errorf("N = %q, want 3", got)
	}
	if got := req.query.Get("step"); got != "300000" {
		t.Errorf("step = %q, want 300000", got)
	}
}
