package opennms

// Tests for the scheduled outages methods – /rest/sched-outages.

import (
	"encoding/json"
	"testing"
)

const schedOutageJSON = `{
	"name": "Weekend-Maintenance",
	"type": "weekly",
	"time": [
		{"day": "saturday", "begins": "00:00:00", "ends": "23:59:59"},
		{"day": "sunday", "begins": "00:00:00", "ends": "23:59:59"}
	],
	"node": [{"id": 1}, {"id": 2}],
	"interface": [{"address": "192.168.0.1"}]
}`

const schedOutageListJSON = `{"scheduleOutage": [` + schedOutageJSON + `]}`

func TestGetSchedOutages(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/sched-outages", schedOutageListJSON)
	result, err := c.GetSchedOutages(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	outage := result["scheduleOutage"].([]any)[0].(map[string]any)
	if outage["name"] != "Weekend-Maintenance" {
		t.Errorf("name = %v, want Weekend-Maintenance", outage["name"])
	}
}

func TestGetSchedOutage(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/sched-outages/Weekend-Maintenance",
		schedOutageJSON)
	result, err := c.GetSchedOutage(t.Context(), "Weekend-Maintenance")
	if err != nil {
		t.Fatal(err)
	}
	if result["name"] != "Weekend-Maintenance" {
		t.Errorf("name = %v, want Weekend-Maintenance", result["name"])
	}
	if result["type"] != "weekly" {
		t.Errorf("type = %v, want weekly", result["type"])
	}
	if times := result["time"].([]any); len(times) != 2 {
		t.Errorf("time entries = %d, want 2", len(times))
	}
}

func TestCreateSchedOutage(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "POST "+v1Path+"/sched-outages", schedOutageJSON)
	var outage map[string]any
	if err := json.Unmarshal([]byte(schedOutageJSON), &outage); err != nil {
		t.Fatal(err)
	}
	if _, err := c.CreateSchedOutage(t.Context(), outage); err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["name"] != "Weekend-Maintenance" {
		t.Errorf("name = %v, want Weekend-Maintenance", body["name"])
	}
	if body["type"] != "weekly" {
		t.Errorf("type = %v, want weekly", body["type"])
	}
	if got := req.header.Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
}

func TestDeleteSchedOutage(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/sched-outages/Weekend-Maintenance", 204)
	if err := c.DeleteSchedOutage(t.Context(), "Weekend-Maintenance"); err != nil {
		t.Fatal(err)
	}
}

func TestAssociateSchedOutageCollectd(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "PUT "+v1Path+"/sched-outages/Weekend-Maintenance/collectd/default", 204)
	err := c.AssociateSchedOutageCollectd(t.Context(), "Weekend-Maintenance", "default")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDissociateSchedOutageCollectd(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/sched-outages/Weekend-Maintenance/collectd/default", 204)
	err := c.DissociateSchedOutageCollectd(t.Context(), "Weekend-Maintenance", "default")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAssociateSchedOutagePollerd(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "PUT "+v1Path+"/sched-outages/Weekend-Maintenance/pollerd/example1", 204)
	err := c.AssociateSchedOutagePollerd(t.Context(), "Weekend-Maintenance", "example1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDissociateSchedOutagePollerd(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/sched-outages/Weekend-Maintenance/pollerd/example1", 204)
	err := c.DissociateSchedOutagePollerd(t.Context(), "Weekend-Maintenance", "example1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAssociateSchedOutageThreshd(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "PUT "+v1Path+"/sched-outages/Weekend-Maintenance/threshd/default", 204)
	err := c.AssociateSchedOutageThreshd(t.Context(), "Weekend-Maintenance", "default")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDissociateSchedOutageThreshd(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/sched-outages/Weekend-Maintenance/threshd/default", 204)
	err := c.DissociateSchedOutageThreshd(t.Context(), "Weekend-Maintenance", "default")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAssociateSchedOutageNotifd(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "PUT "+v1Path+"/sched-outages/Weekend-Maintenance/notifd", 204)
	err := c.AssociateSchedOutageNotifd(t.Context(), "Weekend-Maintenance")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDissociateSchedOutageNotifd(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/sched-outages/Weekend-Maintenance/notifd", 204)
	err := c.DissociateSchedOutageNotifd(t.Context(), "Weekend-Maintenance")
	if err != nil {
		t.Fatal(err)
	}
}
