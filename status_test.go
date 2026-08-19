package opennms

// Tests for the status methods – /api/v2/status.

import "testing"

const statusSummaryJSON = `[["Normal", 42], ["Warning", 3], ["Critical", 1]]`

const statusNodeListJSON = `{
	"count": 1, "offset": 0, "totalCount": 1,
	"nodes": [
		{"id": 1, "label": "router-01", "severity": "CRITICAL"}
	]
}`

const statusApplicationListJSON = `{
	"count": 1, "offset": 0, "totalCount": 1,
	"applications": [
		{"id": 7, "name": "Customer Portal", "severity": "WARNING"}
	]
}`

const statusBusinessServiceListJSON = `{
	"count": 1, "offset": 0, "totalCount": 1,
	"business-services": [
		{"id": 3, "name": "Email", "severity": "NORMAL"}
	]
}`

// statusCheckSummary asserts the parsed statusSummaryJSON shape.
func statusCheckSummary(t *testing.T, result []any) {
	t.Helper()
	if len(result) != 3 {
		t.Fatalf("len = %d, want 3", len(result))
	}
	pair := result[0].([]any)
	if pair[0] != "Normal" || pair[1].(float64) != 42 {
		t.Errorf("result[0] = %v, want [Normal 42]", pair)
	}
}

func TestGetStatusSummaryNodes(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/status/summary/nodes/alarms",
		statusSummaryJSON)
	result, err := c.GetStatusSummaryNodes(t.Context(), "")
	if err != nil {
		t.Fatal(err)
	}
	statusCheckSummary(t, result)
}

func TestGetStatusSummaryNodesOutages(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/status/summary/nodes/outages",
		statusSummaryJSON)
	result, err := c.GetStatusSummaryNodes(t.Context(), "outages")
	if err != nil {
		t.Fatal(err)
	}
	statusCheckSummary(t, result)
}

func TestGetStatusSummaryApplications(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/status/summary/applications",
		statusSummaryJSON)
	result, err := c.GetStatusSummaryApplications(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	statusCheckSummary(t, result)
}

func TestGetStatusSummaryBusinessServices(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/status/summary/business-services",
		statusSummaryJSON)
	result, err := c.GetStatusSummaryBusinessServices(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	statusCheckSummary(t, result)
}

func TestGetStatusNodes(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v2Path+"/status/nodes/alarms",
		statusNodeListJSON)
	result, err := c.GetStatusNodes(t.Context(), "",
		&StatusListOptions{Limit: 10, SeverityFilter: "CRITICAL"})
	if err != nil {
		t.Fatal(err)
	}
	if result["totalCount"].(float64) != 1 {
		t.Errorf("totalCount = %v, want 1", result["totalCount"])
	}
	if got := req.query.Get("limit"); got != "10" {
		t.Errorf("limit = %q, want 10", got)
	}
	if got := req.query.Get("severityFilter"); got != "CRITICAL" {
		t.Errorf("severityFilter = %q, want CRITICAL", got)
	}
}

func TestGetStatusApplications(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v2Path+"/status/applications",
		statusApplicationListJSON)
	result, err := c.GetStatusApplications(t.Context(),
		&StatusListOptions{OrderBy: "severity", Order: "desc"})
	if err != nil {
		t.Fatal(err)
	}
	app := result["applications"].([]any)[0].(map[string]any)
	if app["name"] != "Customer Portal" {
		t.Errorf("name = %v, want Customer Portal", app["name"])
	}
	if got := req.query.Get("orderBy"); got != "severity" {
		t.Errorf("orderBy = %q, want severity", got)
	}
	if got := req.query.Get("order"); got != "desc" {
		t.Errorf("order = %q, want desc", got)
	}
	if got := req.query.Get("limit"); got != "" {
		t.Errorf("limit = %q, want unset", got)
	}
}

func TestGetStatusBusinessServices(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v2Path+"/status/business-services",
		statusBusinessServiceListJSON)
	result, err := c.GetStatusBusinessServices(t.Context(), nil)
	if err != nil {
		t.Fatal(err)
	}
	svc := result["business-services"].([]any)[0].(map[string]any)
	if svc["name"] != "Email" {
		t.Errorf("name = %v, want Email", svc["name"])
	}
	if len(req.query) != 0 {
		t.Errorf("query = %v, want empty", req.query)
	}
}
