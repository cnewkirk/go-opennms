package opennms

// Tests for the geolocation methods – /api/v2/geolocation.

import (
	"encoding/json"
	"testing"
)

const geolocationConfigJSON = `{
	"tileServerName": "OpenStreetMap",
	"tileServerUrl": "https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png",
	"options": {"attribution": "© OpenStreetMap contributors"}
}`

const geolocationListJSON = `[
	{
		"nodeInfo": {"nodeId": 1, "nodeLabel": "router-01"},
		"coordinates": {"longitude": -79.06, "latitude": 35.72},
		"severityInfo": {"id": 5, "label": "Major"},
		"alarmUnackedCount": 3
	}
]`

func TestGetGeolocationConfig(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/geolocation/config", geolocationConfigJSON)
	result, err := c.GetGeolocationConfig(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result["tileServerName"] != "OpenStreetMap" {
		t.Errorf("tileServerName = %v, want OpenStreetMap", result["tileServerName"])
	}
}

func TestQueryGeolocations(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "POST "+v2Path+"/geolocation", geolocationListJSON)
	include := true
	result, err := c.QueryGeolocations(t.Context(), "Outages", "Major", &include)
	if err != nil {
		t.Fatal(err)
	}
	nodeInfo := result[0].(map[string]any)["nodeInfo"].(map[string]any)
	if nodeInfo["nodeLabel"] != "router-01" {
		t.Errorf("nodeLabel = %v, want router-01", nodeInfo["nodeLabel"])
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["strategy"] != "Outages" {
		t.Errorf("strategy = %v, want Outages", body["strategy"])
	}
	if body["severityFilter"] != "Major" {
		t.Errorf("severityFilter = %v, want Major", body["severityFilter"])
	}
	if body["includeAcknowledgedAlarms"] != true {
		t.Errorf("includeAcknowledgedAlarms = %v, want true", body["includeAcknowledgedAlarms"])
	}
}

func TestQueryGeolocationsDefaults(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v2Path+"/geolocation", 204)
	result, err := c.QueryGeolocations(t.Context(), "", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if result != nil {
		t.Errorf("result = %v, want nil", result)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if len(body) != 1 || body["strategy"] != "Alarms" {
		t.Errorf("body = %v, want only strategy=Alarms", body)
	}
}
