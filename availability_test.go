package opennms

// Tests for the availability methods – /rest/availability.

import "testing"

const availabilityJSON = `{
	"section": [{
		"name": "Production",
		"categories": [{
			"name": "Production",
			"normalPercent": 99.5,
			"warningPercent": 97.0
		}]
	}]
}`

const availabilityCategoryJSON = `{
	"name": "Production",
	"normalPercent": 99.5,
	"warningPercent": 97.0,
	"nodes": []
}`

const availabilityNodeJSON = `{
	"id": 1,
	"availability": 99.8,
	"serviceCount": 3,
	"serviceDownCount": 0
}`

func TestGetAvailability(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/availability", availabilityJSON)
	result, err := c.GetAvailability(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	section := result["section"].([]any)[0].(map[string]any)
	if section["name"] != "Production" {
		t.Errorf("name = %v, want Production", section["name"])
	}
}

func TestGetAvailabilityCategory(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/availability/categories/Production",
		availabilityCategoryJSON)
	result, err := c.GetAvailabilityCategory(t.Context(), "Production")
	if err != nil {
		t.Fatal(err)
	}
	if result["name"] != "Production" {
		t.Errorf("name = %v, want Production", result["name"])
	}
}

func TestGetAvailabilityCategoryNodes(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/availability/categories/Production/nodes",
		`{"node": [`+availabilityNodeJSON+`]}`)
	result, err := c.GetAvailabilityCategoryNodes(t.Context(), "Production")
	if err != nil {
		t.Fatal(err)
	}
	node := result["node"].([]any)[0].(map[string]any)
	if node["id"].(float64) != 1 {
		t.Errorf("id = %v, want 1", node["id"])
	}
}

func TestGetAvailabilityCategoryNode(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/availability/categories/Production/nodes/1",
		availabilityNodeJSON)
	result, err := c.GetAvailabilityCategoryNode(t.Context(), "Production", 1)
	if err != nil {
		t.Fatal(err)
	}
	if result["id"].(float64) != 1 {
		t.Errorf("id = %v, want 1", result["id"])
	}
}

func TestGetAvailabilityNode(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/availability/nodes/1", availabilityNodeJSON)
	result, err := c.GetAvailabilityNode(t.Context(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if result["id"].(float64) != 1 {
		t.Errorf("id = %v, want 1", result["id"])
	}
	if result["availability"].(float64) != 99.8 {
		t.Errorf("availability = %v, want 99.8", result["availability"])
	}
}
