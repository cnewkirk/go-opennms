package opennms

// Tests for the applications methods – /api/v2/applications.

import (
	"encoding/json"
	"testing"
)

const applicationJSON = `{
	"id": 1,
	"name": "Web Services",
	"monitoredServices": [
		{"id": 201, "serviceName": "HTTP"}
	]
}`

const applicationListJSON = `{
	"application": [` + applicationJSON + `],
	"totalCount": 1, "count": 1, "offset": 0
}`

func TestGetApplications(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/applications", applicationListJSON)
	result, err := c.GetApplications(t.Context(), nil)
	if err != nil {
		t.Fatal(err)
	}
	app := result["application"].([]any)[0].(map[string]any)
	if app["name"] != "Web Services" {
		t.Errorf("name = %v, want Web Services", app["name"])
	}
}

func TestGetApplication(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/applications/1", applicationJSON)
	result, err := c.GetApplication(t.Context(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if result["id"].(float64) != 1 {
		t.Errorf("id = %v, want 1", result["id"])
	}
}

func TestCreateApplication(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "POST "+v2Path+"/applications", applicationJSON)
	result, err := c.CreateApplication(t.Context(),
		map[string]any{"name": "Web Services"})
	if err != nil {
		t.Fatal(err)
	}
	if result["id"].(float64) != 1 {
		t.Errorf("id = %v, want 1", result["id"])
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["name"] != "Web Services" {
		t.Errorf("body name = %v, want Web Services", body["name"])
	}
}

func TestDeleteApplication(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "DELETE "+v2Path+"/applications/1", 204)
	if err := c.DeleteApplication(t.Context(), 1); err != nil {
		t.Fatal(err)
	}
	if req.method != "DELETE" {
		t.Errorf("method = %q, want DELETE", req.method)
	}
}
