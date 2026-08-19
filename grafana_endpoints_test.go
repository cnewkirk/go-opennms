package opennms

// Tests for the Grafana endpoints methods – /rest/endpoints/grafana.

import (
	"encoding/json"
	"testing"
)

const grafanaEndpointUID = "e2c9dc32-7c8f-4d9e-9b0a-2f0d3c1a4b5c"

const grafanaEndpointJSON = `{
	"id": 1,
	"uid": "` + grafanaEndpointUID + `",
	"url": "https://grafana.example.com",
	"apiKey": "eyJrIjoi...",
	"description": "Production Grafana",
	"connectTimeout": 3000,
	"readTimeout": 3000
}`

const grafanaEndpointListJSON = `[` + grafanaEndpointJSON + `]`

const grafanaDashboardJSON = `{
	"id": 25,
	"uid": "b0d92dk4z",
	"title": "Node performance"
}`

const grafanaDashboardListJSON = `[` + grafanaDashboardJSON + `]`

func TestGetGrafanaEndpoints(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/endpoints/grafana", grafanaEndpointListJSON)
	result, err := c.GetGrafanaEndpoints(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	endpoint := result[0].(map[string]any)
	if endpoint["url"] != "https://grafana.example.com" {
		t.Errorf("url = %v", endpoint["url"])
	}
}

func TestGetGrafanaEndpointsEmpty(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "GET "+v1Path+"/endpoints/grafana", 204)
	result, err := c.GetGrafanaEndpoints(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result != nil {
		t.Errorf("result = %v, want nil", result)
	}
}

func TestGetGrafanaEndpoint(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/endpoints/grafana/1", grafanaEndpointJSON)
	result, err := c.GetGrafanaEndpoint(t.Context(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if result["id"].(float64) != 1 {
		t.Errorf("id = %v, want 1", result["id"])
	}
}

func TestGetGrafanaDashboards(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux,
		"GET "+v1Path+"/endpoints/grafana/"+grafanaEndpointUID+"/dashboards",
		grafanaDashboardListJSON)
	result, err := c.GetGrafanaDashboards(t.Context(), grafanaEndpointUID)
	if err != nil {
		t.Fatal(err)
	}
	dashboard := result[0].(map[string]any)
	if dashboard["title"] != "Node performance" {
		t.Errorf("title = %v", dashboard["title"])
	}
}

func TestGetGrafanaDashboard(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux,
		"GET "+v1Path+"/endpoints/grafana/"+grafanaEndpointUID+"/dashboards/b0d92dk4z",
		grafanaDashboardJSON)
	result, err := c.GetGrafanaDashboard(t.Context(), grafanaEndpointUID, "b0d92dk4z")
	if err != nil {
		t.Fatal(err)
	}
	if result["uid"] != "b0d92dk4z" {
		t.Errorf("uid = %v, want b0d92dk4z", result["uid"])
	}
}

func TestCreateGrafanaEndpoint(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "POST "+v1Path+"/endpoints/grafana",
		"application/json", 202, "")
	var endpoint map[string]any
	if err := json.Unmarshal([]byte(grafanaEndpointJSON), &endpoint); err != nil {
		t.Fatal(err)
	}
	if err := c.CreateGrafanaEndpoint(t.Context(), endpoint); err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatalf("body = %q: %v", req.body, err)
	}
	if body["uid"] != grafanaEndpointUID {
		t.Errorf("uid = %v", body["uid"])
	}
}

func TestVerifyGrafanaEndpoint(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "POST "+v1Path+"/endpoints/grafana/verify",
		"application/json", 200, "")
	err := c.VerifyGrafanaEndpoint(t.Context(),
		map[string]any{"url": "https://grafana.example.com", "apiKey": "abc"})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatalf("body = %q: %v", req.body, err)
	}
	if body["apiKey"] != "abc" {
		t.Errorf("apiKey = %v, want abc", body["apiKey"])
	}
}

func TestUpdateGrafanaEndpoint(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "PUT "+v1Path+"/endpoints/grafana/1",
		"application/json", 202, "")
	var endpoint map[string]any
	if err := json.Unmarshal([]byte(grafanaEndpointJSON), &endpoint); err != nil {
		t.Fatal(err)
	}
	if err := c.UpdateGrafanaEndpoint(t.Context(), 1, endpoint); err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatalf("body = %q: %v", req.body, err)
	}
	if body["id"].(float64) != 1 {
		t.Errorf("id = %v, want 1", body["id"])
	}
}

func TestDeleteGrafanaEndpoints(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/endpoints/grafana", 202)
	if err := c.DeleteGrafanaEndpoints(t.Context()); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteGrafanaEndpoint(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/endpoints/grafana/1", 202)
	if err := c.DeleteGrafanaEndpoint(t.Context(), 1); err != nil {
		t.Fatal(err)
	}
}
