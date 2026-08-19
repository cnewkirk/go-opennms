package opennms

// Tests for the KSC reports methods – /rest/ksc.

import (
	"strings"
	"testing"
)

const kscReportJSON = `{
	"id": 1,
	"label": "Core Bandwidth Report",
	"show_timespan_button": true,
	"show_graphtype_button": false,
	"graphs_per_line": 2,
	"graphs": [
		{
			"title": "Core Switch Bandwidth",
			"resourceId": "node[1].interfaceSnmp[eth0-04013f75f101]",
			"timespan": "7_day",
			"graphtype": "mib2.bits"
		}
	]
}`

const kscReportListJSON = `{
	"kscReport": [{"id": 1, "label": "Core Bandwidth Report"}],
	"totalCount": 1, "count": 1, "offset": 0
}`

// kscReportFixture mirrors the KSC_REPORT fixture as the map a caller
// would pass to CreateKscReport.
func kscReportFixture() map[string]any {
	return map[string]any{
		"id":                    1,
		"label":                 "Core Bandwidth Report",
		"show_timespan_button":  true,
		"show_graphtype_button": false,
		"graphs_per_line":       2,
		"graphs": []map[string]any{{
			"title":      "Core Switch Bandwidth",
			"resourceId": "node[1].interfaceSnmp[eth0-04013f75f101]",
			"timespan":   "7_day",
			"graphtype":  "mib2.bits",
		}},
	}
}

func TestGetKscReports(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/ksc", kscReportListJSON)
	result, err := c.GetKscReports(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	report := result["kscReport"].([]any)[0].(map[string]any)
	if report["id"].(float64) != 1 {
		t.Errorf("id = %v, want 1", report["id"])
	}
	if report["label"] != "Core Bandwidth Report" {
		t.Errorf("label = %v", report["label"])
	}
}

func TestGetKscReport(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/ksc/1", kscReportJSON)
	result, err := c.GetKscReport(t.Context(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if result["id"].(float64) != 1 {
		t.Errorf("id = %v, want 1", result["id"])
	}
	if result["label"] != "Core Bandwidth Report" {
		t.Errorf("label = %v", result["label"])
	}
	graphs := result["graphs"].([]any)
	if len(graphs) != 1 {
		t.Fatalf("len(graphs) = %d, want 1", len(graphs))
	}
	if graphs[0].(map[string]any)["timespan"] != "7_day" {
		t.Errorf("timespan = %v", graphs[0].(map[string]any)["timespan"])
	}
}

func TestGetKscReportCount(t *testing.T) {
	mux, c := newTestClient(t)
	handleText(mux, "GET "+v1Path+"/ksc/count", "5")
	count, err := c.GetKscReportCount(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if count != 5 {
		t.Errorf("count = %d, want 5", count)
	}
}

func TestCreateKscReport(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "POST "+v1Path+"/ksc", xmlType, 204, "")
	if err := c.CreateKscReport(t.Context(), kscReportFixture()); err != nil {
		t.Fatal(err)
	}
	wantPrefix := `<kscReport id="1" label="Core Bandwidth Report"` +
		` show_timespan_button="true" show_graphtype_button="false"` +
		` graphs_per_line="2">`
	if !strings.HasPrefix(req.body, wantPrefix) {
		t.Errorf("body = %q, want prefix %q", req.body, wantPrefix)
	}
	if !strings.Contains(req.body, `<kscGraph title="Core Switch Bandwidth"`) {
		t.Errorf("body = %q, missing kscGraph title", req.body)
	}
	if !strings.Contains(req.body, `graphtype="mib2.bits"`) {
		t.Errorf("body = %q, missing graphtype", req.body)
	}
}

func TestAddGraphToKscReport(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "PUT "+v1Path+"/ksc/1", 204)
	err := c.AddGraphToKscReport(t.Context(), 1, "mib2.bits",
		"node[1].interfaceSnmp[eth0]", "Bandwidth", "7_day")
	if err != nil {
		t.Fatal(err)
	}
	for param, want := range map[string]string{
		"reportName": "mib2.bits",
		"resourceId": "node[1].interfaceSnmp[eth0]",
		"title":      "Bandwidth",
		"timespan":   "7_day",
	} {
		if got := req.query.Get(param); got != want {
			t.Errorf("%s = %q, want %q", param, got, want)
		}
	}
}
