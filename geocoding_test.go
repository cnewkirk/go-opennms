package opennms

// Tests for the geocoding methods – /api/v2/geocoding.

import (
	"encoding/json"
	"reflect"
	"testing"
)

const geocodingConfigJSON = `{"activeGeocoderId": "nominatim"}`

const geocodingGeocoderListJSON = `[
	{
		"id": "nominatim",
		"active": true,
		"config": {"userAgent": "OpenNMS-Meridian/2025"}
	},
	{"id": "google", "active": false, "config": {}}
]`

func TestGetGeocodingConfig(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/geocoding/config", geocodingConfigJSON)
	result, err := c.GetGeocodingConfig(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result["activeGeocoderId"] != "nominatim" {
		t.Errorf("activeGeocoderId = %v, want nominatim", result["activeGeocoderId"])
	}
}

func TestGetGeocoders(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/geocoding/geocoders", geocodingGeocoderListJSON)
	result, err := c.GetGeocoders(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result[0].(map[string]any)["id"] != "nominatim" {
		t.Errorf("id = %v, want nominatim", result[0].(map[string]any)["id"])
	}
}

func TestGetGeocodersEmpty(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "GET "+v2Path+"/geocoding/geocoders", 204)
	result, err := c.GetGeocoders(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result != nil {
		t.Errorf("result = %v, want nil", result)
	}
}

func TestSetActiveGeocoder(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v2Path+"/geocoding/config", 202)
	if err := c.SetActiveGeocoder(t.Context(), "google"); err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(body, map[string]any{"activeGeocoderId": "google"}) {
		t.Errorf("body = %v", body)
	}
}

func TestConfigureGeocoder(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v2Path+"/geocoding/geocoders/google", 204)
	err := c.ConfigureGeocoder(t.Context(), "google",
		map[string]string{"apiKey": "AIza..."})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"config": map[string]any{"apiKey": "AIza..."}}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("body = %v, want %v", body, want)
	}
}

func TestResetGeocodingConfig(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v2Path+"/geocoding/config", 202)
	if err := c.ResetGeocodingConfig(t.Context()); err != nil {
		t.Fatal(err)
	}
}
