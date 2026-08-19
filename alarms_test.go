package opennms

// Tests for the alarms methods – /rest/alarms and /api/v2/alarms.

import (
	"errors"
	"testing"
)

const alarmJSON = `{
	"id": 42,
	"uei": "uei.opennms.org/nodes/nodeDown",
	"nodeId": 1,
	"nodeLabel": "router01.example.com",
	"severity": "MAJOR",
	"count": 3,
	"logMsg": "Node router01.example.com is down.",
	"ackUser": "admin"
}`

const alarmListJSON = `{
	"count": 1, "offset": 0, "totalCount": 1,
	"alarm": [` + alarmJSON + `]
}`

func TestGetAlarmsDefault(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/alarms", alarmListJSON)
	result, err := c.GetAlarms(t.Context(), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	alarm := result["alarm"].([]any)[0].(map[string]any)
	if alarm["id"].(float64) != 42 {
		t.Errorf("id = %v, want 42", alarm["id"])
	}
	if result["totalCount"].(float64) != 1 {
		t.Errorf("totalCount = %v, want 1", result["totalCount"])
	}
	if got := req.query.Get("limit"); got != "10" {
		t.Errorf("limit = %q, want 10", got)
	}
	if got := req.query.Get("offset"); got != "0" {
		t.Errorf("offset = %q, want 0", got)
	}
	if got := req.header.Get("Accept"); got != "application/json, text/plain;q=0.9" {
		t.Errorf("Accept = %q", got)
	}
}

func TestGetAlarmsWithFilters(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/alarms", alarmListJSON)
	_, err := c.GetAlarms(t.Context(),
		&ListOptions{Limit: 25, Offset: 10, OrderBy: "lastEventTime", Order: "descending"},
		map[string]string{"severity": "MAJOR"})
	if err != nil {
		t.Fatal(err)
	}
	for param, want := range map[string]string{
		"limit": "25", "offset": "10", "orderBy": "lastEventTime",
		"order": "descending", "severity": "MAJOR",
	} {
		if got := req.query.Get(param); got != want {
			t.Errorf("%s = %q, want %q", param, got, want)
		}
	}
}

func TestGetAlarm(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/alarms/42", alarmJSON)
	result, err := c.GetAlarm(t.Context(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if result["severity"] != "MAJOR" {
		t.Errorf("severity = %v, want MAJOR", result["severity"])
	}
	if result["uei"] != "uei.opennms.org/nodes/nodeDown" {
		t.Errorf("uei = %v", result["uei"])
	}
}

func TestGetAlarmNotFound(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "GET "+v1Path+"/alarms/99999", 404)
	_, err := c.GetAlarm(t.Context(), 99999)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 404 {
		t.Fatalf("err = %v, want *APIError with StatusCode 404", err)
	}
}

func TestGetAlarmCount(t *testing.T) {
	mux, c := newTestClient(t)
	handleText(mux, "GET "+v1Path+"/alarms/count", "42")
	count, err := c.GetAlarmCount(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if count != 42 {
		t.Errorf("count = %d, want 42", count)
	}
}

func TestAckAlarm(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "PUT "+v1Path+"/alarms/42", formType, 204, "")
	if err := c.AckAlarm(t.Context(), 42, ""); err != nil {
		t.Fatal(err)
	}
	if req.body != "ack=true" {
		t.Errorf("body = %q, want ack=true", req.body)
	}
}

func TestAckAlarmWithUser(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "PUT "+v1Path+"/alarms/42", formType, 204, "")
	if err := c.AckAlarm(t.Context(), 42, "jsmith"); err != nil {
		t.Fatal(err)
	}
	form := req.formBody(t)
	if form.Get("ack") != "true" || form.Get("ackUser") != "jsmith" {
		t.Errorf("form = %v", form)
	}
}

func TestUnackAlarm(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "PUT "+v1Path+"/alarms/42", formType, 204, "")
	if err := c.UnackAlarm(t.Context(), 42); err != nil {
		t.Fatal(err)
	}
	if req.body != "ack=false" {
		t.Errorf("body = %q, want ack=false", req.body)
	}
}

func TestClearAlarm(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "PUT "+v1Path+"/alarms/42", formType, 204, "")
	if err := c.ClearAlarm(t.Context(), 42); err != nil {
		t.Fatal(err)
	}
	if req.body != "clear=true" {
		t.Errorf("body = %q, want clear=true", req.body)
	}
}

func TestEscalateAlarm(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "PUT "+v1Path+"/alarms/42", formType, 204, "")
	if err := c.EscalateAlarm(t.Context(), 42); err != nil {
		t.Fatal(err)
	}
	if req.body != "escalate=true" {
		t.Errorf("body = %q, want escalate=true", req.body)
	}
}

func TestBulkAckAlarms(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "PUT "+v1Path+"/alarms", formType, 204, "")
	err := c.BulkAckAlarms(t.Context(), map[string]string{"severity": "MAJOR"})
	if err != nil {
		t.Fatal(err)
	}
	form := req.formBody(t)
	if form.Get("ack") != "true" || form.Get("severity") != "MAJOR" {
		t.Errorf("form = %v", form)
	}
}

func TestBulkUnackAlarms(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "PUT "+v1Path+"/alarms", formType, 204, "")
	err := c.BulkUnackAlarms(t.Context(),
		map[string]string{"nodeLabel": "router01.example.com"})
	if err != nil {
		t.Fatal(err)
	}
	form := req.formBody(t)
	if form.Get("ack") != "false" || form.Get("nodeLabel") != "router01.example.com" {
		t.Errorf("form = %v", form)
	}
}

func TestBulkClearAlarms(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "PUT "+v1Path+"/alarms", formType, 204, "")
	err := c.BulkClearAlarms(t.Context(), map[string]string{"severity": "CLEARED"})
	if err != nil {
		t.Fatal(err)
	}
	if req.formBody(t).Get("clear") != "true" {
		t.Errorf("body = %q", req.body)
	}
}

func TestBulkEscalateAlarms(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "PUT "+v1Path+"/alarms", formType, 204, "")
	if err := c.BulkEscalateAlarms(t.Context(), nil); err != nil {
		t.Fatal(err)
	}
	if req.formBody(t).Get("escalate") != "true" {
		t.Errorf("body = %q", req.body)
	}
}

func TestGetAlarmsV2Default(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v2Path+"/alarms", alarmListJSON)
	result, err := c.GetAlarmsV2(t.Context(), "", nil)
	if err != nil {
		t.Fatal(err)
	}
	alarm := result["alarm"].([]any)[0].(map[string]any)
	if alarm["id"].(float64) != 42 {
		t.Errorf("id = %v, want 42", alarm["id"])
	}
	if req.query.Get("limit") != "10" || req.query.Get("offset") != "0" {
		t.Errorf("query = %v", req.query)
	}
}

func TestGetAlarmsV2WithFiql(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v2Path+"/alarms", alarmListJSON)
	_, err := c.GetAlarmsV2(t.Context(), "alarm.severity==MAJOR", &ListOptions{Limit: 5})
	if err != nil {
		t.Fatal(err)
	}
	if got := req.query.Get("_s"); got != "alarm.severity==MAJOR" {
		t.Errorf("_s = %q", got)
	}
	if got := req.query.Get("limit"); got != "5" {
		t.Errorf("limit = %q, want 5", got)
	}
}

func TestGetAlarmV2(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v2Path+"/alarms/42", alarmJSON)
	result, err := c.GetAlarmV2(t.Context(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if result["id"].(float64) != 42 {
		t.Errorf("id = %v, want 42", result["id"])
	}
	if req.path != v2Path+"/alarms/42" {
		t.Errorf("path = %q", req.path)
	}
}
