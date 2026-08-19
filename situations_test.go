package opennms

// Tests for the situations methods – /api/v2/situations.

import (
	"encoding/json"
	"testing"
)

const situationJSON = `{
	"id": 99,
	"uei": "uei.opennms.org/alarms/situation",
	"nodeId": 1,
	"nodeLabel": "router01.example.com",
	"ipAddress": "192.168.1.1",
	"serviceType": {"id": 1, "name": "ICMP"},
	"reductionKey": "uei.opennms.org/alarms/situation::99",
	"clearKey": null,
	"alarmType": 1,
	"count": 3,
	"severity": "CRITICAL",
	"firstEventTime": "2024-06-01T08:00:00.000+0000",
	"lastEventTime": "2024-06-01T09:30:00.000+0000",
	"logMsg": "Node router01.example.com is down.",
	"description": "<p>Router is not responding.</p>",
	"operInstruct": null,
	"ackTime": "2024-06-01T09:45:00.000+0000",
	"ackUser": "admin",
	"x733ProbableCause": 0,
	"parameters": [
		{
			"parmName": "ifIndex",
			"value": {"content": "6", "type": "string", "encoding": "text"}
		}
	],
	"relatedAlarms": [
		{"id": 42, "reductionKey": "uei.opennms.org/nodes/nodeDown::1"},
		{"id": 43, "reductionKey": "uei.opennms.org/nodes/nodeDown::2"}
	],
	"lastEvent": {
		"id": 1001,
		"uei": "uei.opennms.org/nodes/nodeDown",
		"time": "2024-06-01T09:30:00.000+0000"
	},
	"isSituation": true,
	"affectedNodeCount": 2
}`

const situationListJSON = `{
	"alarm": [` + situationJSON + `],
	"totalCount": 1, "count": 1, "offset": 0
}`

func TestGetSituationsDefaults(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v2Path+"/situations", situationListJSON)
	result, err := c.GetSituations(t.Context(), nil)
	if err != nil {
		t.Fatal(err)
	}
	situation := result["alarm"].([]any)[0].(map[string]any)
	if situation["isSituation"].(bool) != true {
		t.Errorf("isSituation = %v, want true", situation["isSituation"])
	}
	if situation["id"].(float64) != 99 {
		t.Errorf("id = %v, want 99", situation["id"])
	}
	if got := req.query.Get("limit"); got != "10" {
		t.Errorf("limit = %q, want 10", got)
	}
	if got := req.query.Get("offset"); got != "0" {
		t.Errorf("offset = %q, want 0", got)
	}
}

func TestGetSituationsPagination(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v2Path+"/situations", situationListJSON)
	_, err := c.GetSituations(t.Context(), &ListOptions{Limit: 25, Offset: 50})
	if err != nil {
		t.Fatal(err)
	}
	if got := req.query.Get("limit"); got != "25" {
		t.Errorf("limit = %q, want 25", got)
	}
	if got := req.query.Get("offset"); got != "50" {
		t.Errorf("offset = %q, want 50", got)
	}
}

func TestCreateSituation(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "POST "+v2Path+"/situations/create",
		"application/json", 201, situationJSON)
	result, err := c.CreateSituation(t.Context(), []int{42, 43},
		"Correlated node-down event",
		"Both devices share the same uplink.")
	if err != nil {
		t.Fatal(err)
	}
	if result["isSituation"].(bool) != true {
		t.Errorf("isSituation = %v, want true", result["isSituation"])
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["alarmIdList"] != "42,43" {
		t.Errorf("alarmIdList = %v, want 42,43", body["alarmIdList"])
	}
	if body["description"] != "Correlated node-down event" {
		t.Errorf("description = %v", body["description"])
	}
	if body["diagnosticText"] != "Both devices share the same uplink." {
		t.Errorf("diagnosticText = %v", body["diagnosticText"])
	}
	if got := req.header.Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
}

func TestCreateSituationMinimal(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "POST "+v2Path+"/situations/create",
		"application/json", 201, situationJSON)
	if _, err := c.CreateSituation(t.Context(), []int{42}, "", ""); err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["alarmIdList"] != "42" {
		t.Errorf("alarmIdList = %v, want 42", body["alarmIdList"])
	}
	if _, ok := body["description"]; ok {
		t.Error("description should be absent")
	}
	if _, ok := body["diagnosticText"]; ok {
		t.Error("diagnosticText should be absent")
	}
}

func TestAddAlarmsToSituation(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "POST "+v2Path+"/situations/associateAlarm",
		situationJSON)
	result, err := c.AddAlarmsToSituation(t.Context(), 99, []int{44, 45},
		"Confirmed related")
	if err != nil {
		t.Fatal(err)
	}
	if result["id"].(float64) != 99 {
		t.Errorf("id = %v, want 99", result["id"])
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["situationId"].(float64) != 99 {
		t.Errorf("situationId = %v, want 99", body["situationId"])
	}
	if body["alarmIdList"] != "44,45" {
		t.Errorf("alarmIdList = %v, want 44,45", body["alarmIdList"])
	}
	if body["feedback"] != "Confirmed related" {
		t.Errorf("feedback = %v", body["feedback"])
	}
}

func TestAddAlarmsToSituationNoFeedback(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "POST "+v2Path+"/situations/associateAlarm",
		situationJSON)
	if _, err := c.AddAlarmsToSituation(t.Context(), 99, []int{44}, ""); err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if _, ok := body["feedback"]; ok {
		t.Error("feedback should be absent")
	}
}

func TestClearSituation(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v2Path+"/situations/clear", 204)
	if err := c.ClearSituation(t.Context(), 99); err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["situationId"].(float64) != 99 {
		t.Errorf("situationId = %v, want 99", body["situationId"])
	}
}

func TestClearSituationAlarms(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v2Path+"/situations/alarms/clear", 204)
	if err := c.ClearSituationAlarms(t.Context(), 99, []int{42, 43}); err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["situationId"].(float64) != 99 {
		t.Errorf("situationId = %v, want 99", body["situationId"])
	}
	if body["alarmIdList"] != "42,43" {
		t.Errorf("alarmIdList = %v, want 42,43", body["alarmIdList"])
	}
}

func TestAcceptSituation(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v2Path+"/situations/accepted/99", 204)
	if err := c.AcceptSituation(t.Context(), 99); err != nil {
		t.Fatal(err)
	}
	if req.path != v2Path+"/situations/accepted/99" {
		t.Errorf("path = %q", req.path)
	}
}

func TestRemoveAlarmsFromSituation(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "DELETE "+v2Path+"/situations/removeAlarm", 204)
	if err := c.RemoveAlarmsFromSituation(t.Context(), 99, []int{42, 43}); err != nil {
		t.Fatal(err)
	}
	if req.method != "DELETE" {
		t.Errorf("method = %q, want DELETE", req.method)
	}
	if got := req.query.Get("situationId"); got != "99" {
		t.Errorf("situationId = %q, want 99", got)
	}
	if got := req.query.Get("alarmIdList"); got != "42,43" {
		t.Errorf("alarmIdList = %q, want 42,43", got)
	}
}
