package opennms

// Tests for the reports methods – /rest/reports.

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

const reportTemplateListJSON = `[
	{
		"id": "local_Early-Morning-Report",
		"name": "Early morning report",
		"description": "Global overview of outages, notifications, events",
		"allowAccess": false,
		"online": true
	}
]`

const reportTemplateDetailsJSON = `{
	"id": "local_Early-Morning-Report",
	"name": "Early morning report",
	"formats": [{"name": "PDF"}, {"name": "CSV"}],
	"parameters": [
		{"name": "range", "type": "date", "value": "2026-08-01"}
	],
	"categories": [],
	"timezones": ["US/Eastern", "UTC"]
}`

const reportPDF = "%PDF-1.4 fake report body"

const reportPersistedListJSON = `[
	{
		"id": 1,
		"reportId": "local_Early-Morning-Report",
		"title": "Early morning report",
		"location": "/opt/opennms/share/reports/report.pdf"
	}
]`

const reportScheduledJSON = `{
	"triggerName": "report_trigger_1",
	"reportId": "local_Early-Morning-Report",
	"cronExpression": "0 0 6 * * ?"
}`

const reportScheduledListJSON = `[` + reportScheduledJSON + `]`

// reportsHandlePDF registers a handler answering 200 with a rendered
// application/pdf body.
func reportsHandlePDF(mux *http.ServeMux, pattern, body string) *capture {
	cap := &capture{}
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		cap.record(r)
		w.Header().Set("Content-Type", "application/pdf")
		io.WriteString(w, body)
	})
	return cap
}

func TestGetReportTemplates(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/reports", reportTemplateListJSON)
	result, err := c.GetReportTemplates(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	template := result[0].(map[string]any)
	if template["id"] != "local_Early-Morning-Report" {
		t.Errorf("id = %v", template["id"])
	}
}

func TestGetReportTemplate(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/reports/local_Early-Morning-Report",
		reportTemplateDetailsJSON)
	result, err := c.GetReportTemplate(t.Context(),
		"local_Early-Morning-Report", "admin")
	if err != nil {
		t.Fatal(err)
	}
	format := result["formats"].([]any)[0].(map[string]any)
	if format["name"] != "PDF" {
		t.Errorf("formats[0].name = %v, want PDF", format["name"])
	}
	if got := req.query.Get("userId"); got != "admin" {
		t.Errorf("userId = %q, want admin", got)
	}
}

func TestRunReport(t *testing.T) {
	mux, c := newTestClient(t)
	req := reportsHandlePDF(mux,
		"POST "+v1Path+"/reports/local_Early-Morning-Report", reportPDF)
	params := []any{map[string]any{
		"name": "range", "type": "date", "value": "2026-08-01"}}
	result, err := c.RunReport(t.Context(),
		"local_Early-Morning-Report", "PDF", params)
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != reportPDF {
		t.Errorf("result = %q", result)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["format"] != "PDF" {
		t.Errorf("format = %v, want PDF", body["format"])
	}
	param := body["parameters"].([]any)[0].(map[string]any)
	if param["name"] != "range" || param["type"] != "date" ||
		param["value"] != "2026-08-01" {
		t.Errorf("parameters = %v", body["parameters"])
	}
}

func TestGetPersistedReports(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/reports/persisted", reportPersistedListJSON)
	result, err := c.GetPersistedReports(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result[0].(map[string]any)["id"].(float64) != 1 {
		t.Errorf("id = %v, want 1", result[0])
	}
}

func TestDeliverReport(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v1Path+"/reports/persisted", 202)
	err := c.DeliverReport(t.Context(),
		"local_Early-Morning-Report", "PDF", []any{},
		map[string]any{"instanceId": "morning", "persist": true})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["id"] != "local_Early-Morning-Report" {
		t.Errorf("id = %v", body["id"])
	}
	if body["deliveryOptions"].(map[string]any)["persist"] != true {
		t.Errorf("deliveryOptions = %v", body["deliveryOptions"])
	}
}

func TestDownloadReport(t *testing.T) {
	mux, c := newTestClient(t)
	req := reportsHandlePDF(mux, "GET "+v1Path+"/reports/download", reportPDF)
	result, err := c.DownloadReport(t.Context(), 1, "PDF")
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != reportPDF {
		t.Errorf("result = %q", result)
	}
	if got := req.query.Get("locatorId"); got != "1" {
		t.Errorf("locatorId = %q, want 1", got)
	}
	if got := req.query.Get("format"); got != "PDF" {
		t.Errorf("format = %q, want PDF", got)
	}
}

func TestDeletePersistedReports(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/reports/persisted", 202)
	if err := c.DeletePersistedReports(t.Context()); err != nil {
		t.Fatal(err)
	}
}

func TestDeletePersistedReport(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/reports/persisted/1", 202)
	if err := c.DeletePersistedReport(t.Context(), 1); err != nil {
		t.Fatal(err)
	}
}

func TestGetScheduledReports(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/reports/scheduled", reportScheduledListJSON)
	result, err := c.GetScheduledReports(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result[0].(map[string]any)["triggerName"] != "report_trigger_1" {
		t.Errorf("triggerName = %v", result[0])
	}
}

func TestGetScheduledReport(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/reports/scheduled/report_trigger_1",
		reportScheduledJSON)
	result, err := c.GetScheduledReport(t.Context(), "report_trigger_1")
	if err != nil {
		t.Fatal(err)
	}
	if result["cronExpression"] != "0 0 6 * * ?" {
		t.Errorf("cronExpression = %v", result["cronExpression"])
	}
}

func TestScheduleReport(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v1Path+"/reports/scheduled", 202)
	err := c.ScheduleReport(t.Context(),
		"local_Early-Morning-Report", "PDF", "0 0 6 * * ?", []any{},
		map[string]any{"instanceId": "morning", "sendMail": true,
			"mailTo": "noc@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["cronExpression"] != "0 0 6 * * ?" {
		t.Errorf("cronExpression = %v", body["cronExpression"])
	}
	if body["deliveryOptions"].(map[string]any)["mailTo"] != "noc@example.com" {
		t.Errorf("deliveryOptions = %v", body["deliveryOptions"])
	}
}

func TestUpdateScheduledReport(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "PUT "+v1Path+"/reports/scheduled/report_trigger_1", 202)
	err := c.UpdateScheduledReport(t.Context(), "report_trigger_1",
		map[string]any{"cronExpression": "0 0 7 * * ?"})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["cronExpression"] != "0 0 7 * * ?" {
		t.Errorf("cronExpression = %v", body["cronExpression"])
	}
}

func TestDeleteScheduledReports(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/reports/scheduled", 202)
	if err := c.DeleteScheduledReports(t.Context()); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteScheduledReport(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/reports/scheduled/report_trigger_1", 202)
	if err := c.DeleteScheduledReport(t.Context(), "report_trigger_1"); err != nil {
		t.Fatal(err)
	}
}
