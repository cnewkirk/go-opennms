package opennms

// Tests for the monitoring systems methods – /rest/monitoringSystems.

import "testing"

const monitoringSystemJSON = `{
	"id": "00000000-0000-0000-0000-000000000000",
	"label": "localhost",
	"location": "Default",
	"type": "OpenNMS",
	"status": "UP",
	"lastUpdated": "2024-06-01T09:00:00.000+0000"
}`

func TestGetMonitoringSystem(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/monitoringSystems/main",
		monitoringSystemJSON)
	result, err := c.GetMonitoringSystem(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result["type"] != "OpenNMS" {
		t.Errorf("type = %v, want OpenNMS", result["type"])
	}
	if result["location"] != "Default" {
		t.Errorf("location = %v, want Default", result["location"])
	}
}
