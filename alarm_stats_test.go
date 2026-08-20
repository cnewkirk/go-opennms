package opennms

// Tests for the alarm statistics methods – /rest/stats/alarms.

import "testing"

const alarmStatsAlarmFieldsJSON = `"uei": "uei.opennms.org/nodes/nodeDown",
	"nodeId": 1,
	"nodeLabel": "router01.example.com",
	"ipAddress": "192.168.1.1",
	"serviceType": {"id": 1, "name": "ICMP"},
	"reductionKey": "uei.opennms.org/nodes/nodeDown::1",
	"clearKey": null,
	"alarmType": 1,
	"count": 3,
	"severity": "MAJOR",
	"firstEventTime": "2024-06-01T08:00:00.000+0000",
	"lastEventTime": "2024-06-01T09:30:00.000+0000",
	"logMsg": "Node router01.example.com is down.",
	"description": "<p>Router is not responding.</p>",
	"operInstruct": null,
	"x733ProbableCause": 0,
	"parameters": [{"parmName": "ifIndex",
		"value": {"content": "6", "type": "string", "encoding": "text"}}],
	"relatedAlarms": [],
	"lastEvent": {"id": 1001, "uei": "uei.opennms.org/nodes/nodeDown",
		"time": "2024-06-01T09:30:00.000+0000"}`

const alarmStatsJSON = `{
	"totalCount": 15,
	"acknowledgedCount": 5,
	"unacknowledgedCount": 10,
	"newestAcknowledged": {"id": 42,
		"ackTime": "2024-06-01T09:45:00.000+0000", "ackUser": "admin",
		` + alarmStatsAlarmFieldsJSON + `},
	"newestUnacknowledged": {"id": 43, "ackTime": null, "ackUser": null,
		` + alarmStatsAlarmFieldsJSON + `},
	"oldestAcknowledged": {"id": 10,
		"ackTime": "2024-06-01T09:45:00.000+0000", "ackUser": "admin",
		` + alarmStatsAlarmFieldsJSON + `},
	"oldestUnacknowledged": {"id": 11, "ackTime": null, "ackUser": null,
		` + alarmStatsAlarmFieldsJSON + `}
}`

const alarmStatsSeverityCriticalJSON = `{"severity": "CRITICAL",
	"totalCount": 2, "acknowledgedCount": 0, "unacknowledgedCount": 2}`

const alarmStatsSeverityMajorJSON = `{"severity": "MAJOR",
	"totalCount": 5, "acknowledgedCount": 2, "unacknowledgedCount": 3}`

const alarmStatsSeverityMinorJSON = `{"severity": "MINOR",
	"totalCount": 4, "acknowledgedCount": 2, "unacknowledgedCount": 2}`

const alarmStatsBySeverityJSON = `{"alarmStatistics": [` +
	alarmStatsSeverityCriticalJSON + `,
	` + alarmStatsSeverityMajorJSON + `,
	` + alarmStatsSeverityMinorJSON + `]}`

const alarmStatsBySeverityFilteredJSON = `{"alarmStatistics": [` +
	alarmStatsSeverityCriticalJSON + `,
	` + alarmStatsSeverityMajorJSON + `]}`

func TestGetAlarmStats(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/stats/alarms", alarmStatsJSON)
	result, err := c.GetAlarmStats(t.Context(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if result["totalCount"].(float64) != 15 {
		t.Errorf("totalCount = %v, want 15", result["totalCount"])
	}
	if result["acknowledgedCount"].(float64) != 5 {
		t.Errorf("acknowledgedCount = %v, want 5", result["acknowledgedCount"])
	}
}

func TestGetAlarmStatsWithFilters(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/stats/alarms", alarmStatsJSON)
	_, err := c.GetAlarmStats(t.Context(), map[string]string{"severity": "MAJOR"})
	if err != nil {
		t.Fatal(err)
	}
	if got := req.query.Get("severity"); got != "MAJOR" {
		t.Errorf("severity = %q, want MAJOR", got)
	}
}

func TestGetAlarmStatsBySeverity(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/stats/alarms/by-severity",
		alarmStatsBySeverityJSON)
	result, err := c.GetAlarmStatsBySeverity(t.Context(), nil)
	if err != nil {
		t.Fatal(err)
	}
	stats := result["alarmStatistics"].([]any)
	if stats[0].(map[string]any)["severity"] != "CRITICAL" {
		t.Errorf("severity = %v, want CRITICAL", stats[0])
	}
}

func TestGetAlarmStatsBySeverityFiltered(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/stats/alarms/by-severity",
		alarmStatsBySeverityFilteredJSON)
	_, err := c.GetAlarmStatsBySeverity(t.Context(),
		[]string{"CRITICAL", "MAJOR"})
	if err != nil {
		t.Fatal(err)
	}
	if got := req.query.Get("severities"); got != "CRITICAL,MAJOR" {
		t.Errorf("severities = %q, want CRITICAL,MAJOR", got)
	}
}
