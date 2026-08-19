package opennms

// Tests for the syslog NBI methods – /rest/config/syslog-nbi.

import (
	"encoding/json"
	"testing"
)

const syslogNbiConfigJSON = `{
	"enabled": false,
	"destinations": []
}`

const syslogNbiDestinationJSON = `{
	"name": "siem",
	"host": "10.0.0.2",
	"port": 514,
	"firstOccurrenceOnly": false,
	"filters": []
}`

const syslogNbiDestinationListJSON = `{
	"destination": [` + syslogNbiDestinationJSON + `]
}`

func TestGetSyslogNbiConfig(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/config/syslog-nbi", syslogNbiConfigJSON)
	result, err := c.GetSyslogNbiConfig(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result["enabled"] != false {
		t.Errorf("enabled = %v, want false", result["enabled"])
	}
}

func TestGetSyslogNbiStatus(t *testing.T) {
	mux, c := newTestClient(t)
	handleText(mux, "GET "+v1Path+"/config/syslog-nbi/status", "false")
	result, err := c.GetSyslogNbiStatus(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result != "false" {
		t.Errorf("status = %q, want false", result)
	}
}

func TestSetSyslogNbiStatus(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "PUT "+v1Path+"/config/syslog-nbi/status", 204)
	if err := c.SetSyslogNbiStatus(t.Context(), true); err != nil {
		t.Fatal(err)
	}
	if req.body != "enabled=true" {
		t.Errorf("body = %q, want enabled=true", req.body)
	}
}

func TestGetSyslogNbiDestinations(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/config/syslog-nbi/destinations",
		syslogNbiDestinationListJSON)
	result, err := c.GetSyslogNbiDestinations(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	destination := result["destination"].([]any)[0].(map[string]any)
	if destination["name"] != "siem" {
		t.Errorf("name = %v, want siem", destination["name"])
	}
}

func TestGetSyslogNbiDestination(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/config/syslog-nbi/destinations/siem",
		syslogNbiDestinationJSON)
	result, err := c.GetSyslogNbiDestination(t.Context(), "siem")
	if err != nil {
		t.Fatal(err)
	}
	if result["host"] != "10.0.0.2" {
		t.Errorf("host = %v, want 10.0.0.2", result["host"])
	}
}

func TestCreateSyslogNbiDestination(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "POST "+v1Path+"/config/syslog-nbi/destinations",
		"application/json", 201, "")
	data := map[string]any{"name": "new-syslog", "host": "10.0.0.3", "port": 514}
	if err := c.CreateSyslogNbiDestination(t.Context(), data); err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatalf("body = %q: %v", req.body, err)
	}
	if body["name"] != "new-syslog" {
		t.Errorf("name = %v, want new-syslog", body["name"])
	}
}

func TestUpdateSyslogNbiDestination(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "PUT "+v1Path+"/config/syslog-nbi/destinations/siem", 204)
	err := c.UpdateSyslogNbiDestination(t.Context(), "siem",
		map[string]string{"port": "515"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteSyslogNbiDestination(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/config/syslog-nbi/destinations/siem", 204)
	if err := c.DeleteSyslogNbiDestination(t.Context(), "siem"); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateSyslogNbiConfig(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "POST "+v1Path+"/config/syslog-nbi", 204)
	err := c.UpdateSyslogNbiConfig(t.Context(), map[string]any{"enabled": true})
	if err != nil {
		t.Fatal(err)
	}
}
