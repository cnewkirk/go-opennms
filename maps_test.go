package opennms

// Tests for the maps methods – /rest/maps.

import (
	"encoding/json"
	"testing"
)

const mapJSON = `{
	"id": 1,
	"name": "Core Network",
	"mapWidth": 1024,
	"mapHeight": 768,
	"accessMode": "RW",
	"owner": "admin",
	"lastModifiedTime": "2024-05-01T10:00:00.000+0000",
	"createTime": "2024-01-01T08:00:00.000+0000"
}`

const mapListJSON = `{
	"map": [` + mapJSON + `],
	"totalCount": 1, "count": 1, "offset": 0
}`

const mapElementsJSON = `{
	"mapElement": [
		{
			"id": 1001,
			"mapId": 1,
			"elementId": 1,
			"type": "N",
			"label": "router01.example.com",
			"x": 200,
			"y": 150,
			"severity": "MAJOR"
		}
	]
}`

func TestGetMaps(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/maps", mapListJSON)
	result, err := c.GetMaps(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	m := result["map"].([]any)[0].(map[string]any)
	if m["id"].(float64) != 1 {
		t.Errorf("id = %v, want 1", m["id"])
	}
	if m["name"] != "Core Network" {
		t.Errorf("name = %v, want Core Network", m["name"])
	}
}

func TestGetMap(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/maps/1", mapJSON)
	result, err := c.GetMap(t.Context(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if result["id"].(float64) != 1 {
		t.Errorf("id = %v, want 1", result["id"])
	}
	if result["accessMode"] != "RW" {
		t.Errorf("accessMode = %v, want RW", result["accessMode"])
	}
}

func TestGetMapElements(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/maps/1/mapElements", mapElementsJSON)
	result, err := c.GetMapElements(t.Context(), 1)
	if err != nil {
		t.Fatal(err)
	}
	element := result["mapElement"].([]any)[0].(map[string]any)
	if element["elementId"].(float64) != 1 {
		t.Errorf("elementId = %v, want 1", element["elementId"])
	}
	if element["type"] != "N" {
		t.Errorf("type = %v, want N", element["type"])
	}
	if req.path != v1Path+"/maps/1/mapElements" {
		t.Errorf("path = %q", req.path)
	}
}

func TestCreateMap(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "POST "+v1Path+"/maps", mapJSON)
	newMap := map[string]any{
		"name": "Core Network", "mapWidth": 1024, "mapHeight": 768,
		"accessMode": "RW", "owner": "admin"}
	result, err := c.CreateMap(t.Context(), newMap)
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
	if body["name"] != "Core Network" {
		t.Errorf("body name = %v, want Core Network", body["name"])
	}
}

func TestUpdateMap(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "PUT "+v1Path+"/maps/1", 204)
	err := c.UpdateMap(t.Context(), 1, map[string]any{"name": "Updated Map Name"})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["name"] != "Updated Map Name" {
		t.Errorf("body name = %v, want Updated Map Name", body["name"])
	}
}

func TestDeleteMap(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/maps/1", 204)
	if err := c.DeleteMap(t.Context(), 1); err != nil {
		t.Fatal(err)
	}
}
