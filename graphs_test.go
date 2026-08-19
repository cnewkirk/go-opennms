package opennms

// Tests for the graph / topology methods – /rest/graphs.

import (
	"encoding/json"
	"strings"
	"testing"
)

const graphJSON = `{
	"namespace": "nodes",
	"description": "The default OpenNMS node graph",
	"preferredLayout": "Grid Layout",
	"focus": {"type": "SELECTION", "vertices": []}
}`

const graphContainerJSON = `{
	"id": "nodes",
	"label": "Nodes",
	"graphs": [` + graphJSON + `]
}`

const graphContainerListJSON = `{"graphContainer": [` + graphContainerJSON + `]}`

const graphSuggestionsJSON = `{
	"suggestion": [
		{"label": "router01.example.com", "context": "nodes"}
	]
}`

const graphSearchResultsJSON = `{
	"searchResult": [
		{"id": "nodes:1", "label": "router01.example.com", "namespace": "nodes"}
	]
}`

func TestGetGraphContainers(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/graphs", graphContainerListJSON)
	result, err := c.GetGraphContainers(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	container := result["graphContainer"].([]any)[0].(map[string]any)
	if container["id"] != "nodes" {
		t.Errorf("id = %v, want nodes", container["id"])
	}
}

func TestGetGraphContainer(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/graphs/nodes", graphContainerJSON)
	result, err := c.GetGraphContainer(t.Context(), "nodes")
	if err != nil {
		t.Fatal(err)
	}
	if result["id"] != "nodes" {
		t.Errorf("id = %v, want nodes", result["id"])
	}
	if result["label"] != "Nodes" {
		t.Errorf("label = %v, want Nodes", result["label"])
	}
}

func TestGetGraph(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/graphs/nodes/nodes", graphJSON)
	result, err := c.GetGraph(t.Context(), "nodes", "nodes")
	if err != nil {
		t.Fatal(err)
	}
	if result["namespace"] != "nodes" {
		t.Errorf("namespace = %v, want nodes", result["namespace"])
	}
	if !strings.Contains(req.path, "/graphs/nodes/nodes") {
		t.Errorf("path = %q", req.path)
	}
}

func TestGetGraphViewDefault(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "POST "+v1Path+"/graphs/nodes/nodes", graphJSON)
	result, err := c.GetGraphView(t.Context(), "nodes", "nodes", 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	if result["namespace"] != "nodes" {
		t.Errorf("namespace = %v, want nodes", result["namespace"])
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["semanticZoomLevel"].(float64) != 1 {
		t.Errorf("semanticZoomLevel = %v, want 1", body["semanticZoomLevel"])
	}
	if vertices := body["verticesInFocus"].([]any); len(vertices) != 0 {
		t.Errorf("verticesInFocus = %v, want []", vertices)
	}
}

func TestGetGraphViewWithFocus(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "POST "+v1Path+"/graphs/nodes/nodes", graphJSON)
	_, err := c.GetGraphView(t.Context(), "nodes", "nodes", 2,
		[]string{"nodes:1", "nodes:2"})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["semanticZoomLevel"].(float64) != 2 {
		t.Errorf("semanticZoomLevel = %v, want 2", body["semanticZoomLevel"])
	}
	vertices := body["verticesInFocus"].([]any)
	if vertices[0] != "nodes:1" {
		t.Errorf("verticesInFocus = %v", vertices)
	}
}

func TestGetGraphSearchSuggestions(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/graphs/search/suggestions/nodes",
		graphSuggestionsJSON)
	result, err := c.GetGraphSearchSuggestions(t.Context(), "nodes", "router01")
	if err != nil {
		t.Fatal(err)
	}
	suggestion := result["suggestion"].([]any)[0].(map[string]any)
	if suggestion["label"] != "router01.example.com" {
		t.Errorf("label = %v", suggestion["label"])
	}
	if got := req.query.Get("s"); got != "router01" {
		t.Errorf("s = %q, want router01", got)
	}
}

func TestGetGraphSearchResults(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/graphs/search/results/nodes",
		graphSearchResultsJSON)
	result, err := c.GetGraphSearchResults(t.Context(),
		"nodes", "NodeSearchProvider", "router01", "")
	if err != nil {
		t.Fatal(err)
	}
	searchResult := result["searchResult"].([]any)[0].(map[string]any)
	if searchResult["namespace"] != "nodes" {
		t.Errorf("namespace = %v, want nodes", searchResult["namespace"])
	}
	if got := req.query.Get("providerId"); got != "NodeSearchProvider" {
		t.Errorf("providerId = %q", got)
	}
	if got := req.query.Get("criteria"); got != "router01" {
		t.Errorf("criteria = %q", got)
	}
}
