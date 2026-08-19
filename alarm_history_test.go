package opennms

// Tests for the alarm history methods – /rest/alarms/history.

import (
	"strings"
	"testing"
)

const alarmHistoryAlarmJSON = `{
	"id": 42,
	"uei": "uei.opennms.org/nodes/nodeDown",
	"nodeId": 1,
	"nodeLabel": "router01.example.com",
	"severity": "MAJOR",
	"count": 3,
	"logMsg": "Node router01.example.com is down.",
	"ackUser": "admin"
}`

const alarmHistoryListJSON = `[` + alarmHistoryAlarmJSON + `]`

const alarmHistoryStateJSON = `{
	"id": "42:2024-06-01T09:30:00.000+0000",
	"alarmId": 42,
	"time": "2024-06-01T09:30:00.000+0000",
	"type": "ALARM_CREATED",
	"user": null,
	"alarm": ` + alarmHistoryAlarmJSON + `
}`

const alarmHistoryStatesListJSON = `[` + alarmHistoryStateJSON + `]`

func TestGetAlarmHistory(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/alarms/history", alarmHistoryListJSON)
	result, err := c.GetAlarmHistory(t.Context(), 0)
	if err != nil {
		t.Fatal(err)
	}
	if result[0].(map[string]any)["id"].(float64) != 42 {
		t.Errorf("id = %v, want 42", result[0].(map[string]any)["id"])
	}
}

func TestGetAlarmHistoryWithTimestamp(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/alarms/history", alarmHistoryListJSON)
	_, err := c.GetAlarmHistory(t.Context(), 1717228800000)
	if err != nil {
		t.Fatal(err)
	}
	if got := req.query.Get("at"); got != "1717228800000" {
		t.Errorf("at = %q, want 1717228800000", got)
	}
}

func TestGetAlarmHistoryAt(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/alarms/history/42", alarmHistoryStateJSON)
	result, err := c.GetAlarmHistoryAt(t.Context(), 42, 0)
	if err != nil {
		t.Fatal(err)
	}
	if result["alarmId"].(float64) != 42 {
		t.Errorf("alarmId = %v, want 42", result["alarmId"])
	}
	if result["type"] != "ALARM_CREATED" {
		t.Errorf("type = %v, want ALARM_CREATED", result["type"])
	}
}

func TestGetAlarmHistoryAtWithTimestamp(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/alarms/history/42", alarmHistoryStateJSON)
	_, err := c.GetAlarmHistoryAt(t.Context(), 42, 1717228800000)
	if err != nil {
		t.Fatal(err)
	}
	if got := req.query.Get("at"); got != "1717228800000" {
		t.Errorf("at = %q, want 1717228800000", got)
	}
}

func TestGetAlarmHistoryStates(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/alarms/history/42/states",
		alarmHistoryStatesListJSON)
	result, err := c.GetAlarmHistoryStates(t.Context(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if result[0].(map[string]any)["alarmId"].(float64) != 42 {
		t.Errorf("alarmId = %v, want 42", result[0].(map[string]any)["alarmId"])
	}
	if !strings.HasSuffix(req.path, "/alarms/history/42/states") {
		t.Errorf("path = %q", req.path)
	}
}
