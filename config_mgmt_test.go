package opennms

// Tests for the configuration management methods – /rest/cm.

import (
	"encoding/json"
	"testing"
)

const configMgmtNamesJSON = `["provisiond", "poller", "collectd"]`

const configMgmtSchemasJSON = `{
	"schema": [
		{"name": "provisiond", "description": "Provisioning daemon config"}
	]
}`

const configMgmtSchemaJSON = `{
	"name": "provisiond",
	"description": "Provisioning daemon config",
	"properties": {}
}`

const configMgmtIdsJSON = `["default"]`

const configMgmtConfigJSON = `{
	"name": "provisiond",
	"id": "default",
	"properties": {"importThreads": 8}
}`

func TestGetConfigNames(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/cm", configMgmtNamesJSON)
	result, err := c.GetConfigNames(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, name := range result {
		if name == "provisiond" {
			found = true
		}
	}
	if !found {
		t.Errorf("result = %v, want to contain provisiond", result)
	}
}

func TestGetConfigSchemas(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/cm/schema", configMgmtSchemasJSON)
	result, err := c.GetConfigSchemas(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	schema := result["schema"].([]any)[0].(map[string]any)
	if schema["name"] != "provisiond" {
		t.Errorf("name = %v, want provisiond", schema["name"])
	}
}

func TestGetConfigSchema(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/cm/schema/provisiond", configMgmtSchemaJSON)
	result, err := c.GetConfigSchema(t.Context(), "provisiond")
	if err != nil {
		t.Fatal(err)
	}
	if result["name"] != "provisiond" {
		t.Errorf("name = %v, want provisiond", result["name"])
	}
}

func TestGetConfigIds(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/cm/provisiond", configMgmtIdsJSON)
	result, err := c.GetConfigIds(t.Context(), "provisiond")
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, id := range result {
		if id == "default" {
			found = true
		}
	}
	if !found {
		t.Errorf("result = %v, want to contain default", result)
	}
}

func TestGetConfig(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/cm/provisiond/default", configMgmtConfigJSON)
	result, err := c.GetConfig(t.Context(), "provisiond", "default")
	if err != nil {
		t.Fatal(err)
	}
	if result["name"] != "provisiond" {
		t.Errorf("name = %v, want provisiond", result["name"])
	}
}

func TestGetConfigPart(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/cm/provisiond/default/importThreads",
		`{"value": 8}`)
	result, err := c.GetConfigPart(t.Context(), "provisiond", "default",
		"importThreads")
	if err != nil {
		t.Fatal(err)
	}
	if result["value"].(float64) != 8 {
		t.Errorf("value = %v, want 8", result["value"])
	}
}

func TestCreateConfig(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v1Path+"/cm/provisiond/default", 201)
	err := c.CreateConfig(t.Context(), "provisiond", "default",
		map[string]any{"importThreads": 8})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["importThreads"].(float64) != 8 {
		t.Errorf("importThreads = %v, want 8", body["importThreads"])
	}
}

func TestUpdateConfig(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "PUT "+v1Path+"/cm/provisiond/default", 204)
	err := c.UpdateConfig(t.Context(), "provisiond", "default",
		map[string]any{"importThreads": 16})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteConfig(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/cm/provisiond/default", 204)
	if err := c.DeleteConfig(t.Context(), "provisiond", "default"); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteConfigPart(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/cm/provisiond/default/importThreads", 204)
	err := c.DeleteConfigPart(t.Context(), "provisiond", "default",
		"importThreads")
	if err != nil {
		t.Fatal(err)
	}
}
