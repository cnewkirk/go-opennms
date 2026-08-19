package opennms

// Tests for the monitoring locations methods –
// /rest/monitoringLocations.

import (
	"encoding/json"
	"testing"
)

const monitoringLocationJSON = `{
	"location-name": "Default",
	"monitoring-area": "default",
	"priority": 100,
	"tags": []
}`

const monitoringLocationListJSON = `{
	"count": 1, "offset": 0, "totalCount": 1,
	"location": [` + monitoringLocationJSON + `]
}`

func TestGetMonitoringLocations(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/monitoringLocations",
		monitoringLocationListJSON)
	result, err := c.GetMonitoringLocations(t.Context(), nil)
	if err != nil {
		t.Fatal(err)
	}
	location := result["location"].([]any)[0].(map[string]any)
	if location["location-name"] != "Default" {
		t.Errorf("location-name = %v, want Default", location["location-name"])
	}
}

func TestGetMonitoringLocation(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/monitoringLocations/Default",
		monitoringLocationJSON)
	result, err := c.GetMonitoringLocation(t.Context(), "Default")
	if err != nil {
		t.Fatal(err)
	}
	if result["location-name"] != "Default" {
		t.Errorf("location-name = %v, want Default", result["location-name"])
	}
}

func TestGetDefaultMonitoringLocation(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/monitoringLocations/default",
		monitoringLocationJSON)
	result, err := c.GetDefaultMonitoringLocation(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result["location-name"] != "Default" {
		t.Errorf("location-name = %v, want Default", result["location-name"])
	}
}

func TestGetMonitoringLocationCount(t *testing.T) {
	mux, c := newTestClient(t)
	handleText(mux, "GET "+v1Path+"/monitoringLocations/count", "3")
	count, err := c.GetMonitoringLocationCount(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Errorf("count = %d, want 3", count)
	}
}

func TestCreateMonitoringLocation(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v1Path+"/monitoringLocations", 201)
	err := c.CreateMonitoringLocation(t.Context(), map[string]any{
		"location-name": "Remote", "monitoring-area": "remote"})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["location-name"] != "Remote" {
		t.Errorf("location-name = %v, want Remote", body["location-name"])
	}
}

func TestUpdateMonitoringLocation(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "PUT "+v1Path+"/monitoringLocations/Default", 204)
	err := c.UpdateMonitoringLocation(t.Context(), "Default",
		map[string]string{"monitoring-area": "updated"})
	if err != nil {
		t.Fatal(err)
	}
	if req.body != "monitoring-area=updated" {
		t.Errorf("body = %q, want monitoring-area=updated", req.body)
	}
}

func TestDeleteMonitoringLocation(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/monitoringLocations/Remote", 204)
	if err := c.DeleteMonitoringLocation(t.Context(), "Remote"); err != nil {
		t.Fatal(err)
	}
}
