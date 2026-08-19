package opennms

// Tests for the events methods – /rest/events.

import (
	"encoding/json"
	"testing"
)

const eventJSON = `{
	"id": 1001,
	"uei": "uei.opennms.org/nodes/nodeDown",
	"label": "Node Down",
	"time": "2024-06-01T09:30:00.000+0000",
	"createTime": "2024-06-01T09:30:01.000+0000",
	"source": "OpenNMS.Poller.Monitor.IcmpMonitor",
	"nodeId": 1,
	"nodeLabel": "router01.example.com",
	"ipAddress": "192.168.1.1",
	"serviceType": {"id": 1, "name": "ICMP"},
	"severity": "MAJOR",
	"logMsg": "Node router01.example.com is down.",
	"logMsgDest": "logndisplay",
	"eventDisplay": "Y",
	"ackUser": null,
	"ackTime": null,
	"parameters": [{
		"parmName": "ifIndex",
		"value": {"content": "6", "type": "string", "encoding": "text"}
	}]
}`

const eventListJSON = `{
	"count": 1, "offset": 0, "totalCount": 1,
	"event": [` + eventJSON + `]
}`

func TestGetEventsDefault(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/events", eventListJSON)
	result, err := c.GetEvents(t.Context(), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	event := result["event"].([]any)[0].(map[string]any)
	if event["id"].(float64) != 1001 {
		t.Errorf("id = %v, want 1001", event["id"])
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
}

func TestGetEventsWithFilters(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/events", eventListJSON)
	_, err := c.GetEvents(t.Context(),
		&ListOptions{Limit: 50, OrderBy: "eventTime", Order: "descending"},
		map[string]string{"uei": "uei.opennms.org/nodes/nodeDown"})
	if err != nil {
		t.Fatal(err)
	}
	for param, want := range map[string]string{
		"limit": "50", "uei": "uei.opennms.org/nodes/nodeDown",
		"orderBy": "eventTime", "order": "descending",
	} {
		if got := req.query.Get(param); got != want {
			t.Errorf("%s = %q, want %q", param, got, want)
		}
	}
}

func TestGetEvent(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/events/1001", eventJSON)
	result, err := c.GetEvent(t.Context(), 1001)
	if err != nil {
		t.Fatal(err)
	}
	if result["id"].(float64) != 1001 {
		t.Errorf("id = %v, want 1001", result["id"])
	}
	if result["severity"] != "MAJOR" {
		t.Errorf("severity = %v, want MAJOR", result["severity"])
	}
}

func TestGetEventCount(t *testing.T) {
	mux, c := newTestClient(t)
	handleText(mux, "GET "+v1Path+"/events/count", "1337")
	count, err := c.GetEventCount(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if count != 1337 {
		t.Errorf("count = %d, want 1337", count)
	}
}

func TestCreateEvent(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v1Path+"/events", 200)
	err := c.CreateEvent(t.Context(), map[string]any{
		"uei":      "uei.opennms.org/internal/test",
		"source":   "gotest",
		"severity": "Normal",
	})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["uei"] != "uei.opennms.org/internal/test" {
		t.Errorf("uei = %v", body["uei"])
	}
	if body["source"] != "gotest" {
		t.Errorf("source = %v", body["source"])
	}
	if got := req.header.Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
}

func TestAckEvent(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "PUT "+v1Path+"/events/1001", formType, 204, "")
	if err := c.AckEvent(t.Context(), 1001); err != nil {
		t.Fatal(err)
	}
	if req.body != "ack=true" {
		t.Errorf("body = %q, want ack=true", req.body)
	}
}

func TestUnackEvent(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "PUT "+v1Path+"/events/1001", formType, 204, "")
	if err := c.UnackEvent(t.Context(), 1001); err != nil {
		t.Fatal(err)
	}
	if req.body != "ack=false" {
		t.Errorf("body = %q, want ack=false", req.body)
	}
}

func TestBulkAckEvents(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "PUT "+v1Path+"/events", formType, 204, "")
	err := c.BulkAckEvents(t.Context(), map[string]string{"nodeId": "1"})
	if err != nil {
		t.Fatal(err)
	}
	form := req.formBody(t)
	if form.Get("ack") != "true" || form.Get("nodeId") != "1" {
		t.Errorf("form = %v", form)
	}
}

func TestBulkUnackEvents(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "PUT "+v1Path+"/events", formType, 204, "")
	if err := c.BulkUnackEvents(t.Context(), nil); err != nil {
		t.Fatal(err)
	}
	if req.body != "ack=false" {
		t.Errorf("body = %q, want ack=false", req.body)
	}
}
