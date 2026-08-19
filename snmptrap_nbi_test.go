package opennms

// Tests for the snmptrap_nbi methods – /rest/config/snmptrap-nbi.

import (
	"encoding/json"
	"testing"
)

const snmptrapNbiTrapsinkJSON = `{
	"name": "remote-nms",
	"ipAddress": "10.0.0.1",
	"port": 162,
	"community": "public"
}`

const snmptrapNbiConfigJSON = `{
	"enabled": true,
	"trapsinks": [` + snmptrapNbiTrapsinkJSON + `]
}`

const snmptrapNbiTrapsinkListJSON = `{"trapsink": [` + snmptrapNbiTrapsinkJSON + `]}`

func TestGetSnmptrapNbiConfig(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/config/snmptrap-nbi", snmptrapNbiConfigJSON)
	result, err := c.GetSnmptrapNbiConfig(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result["enabled"] != true {
		t.Errorf("enabled = %v, want true", result["enabled"])
	}
}

func TestGetSnmptrapNbiStatus(t *testing.T) {
	mux, c := newTestClient(t)
	handleText(mux, "GET "+v1Path+"/config/snmptrap-nbi/status", "true")
	result, err := c.GetSnmptrapNbiStatus(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result != "true" {
		t.Errorf("status = %q, want true", result)
	}
}

func TestSetSnmptrapNbiStatus(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "PUT "+v1Path+"/config/snmptrap-nbi/status", 204)
	if err := c.SetSnmptrapNbiStatus(t.Context(), false); err != nil {
		t.Fatal(err)
	}
	if req.body != "enabled=false" {
		t.Errorf("body = %q, want enabled=false", req.body)
	}
}

func TestGetSnmptrapNbiTrapsinks(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/config/snmptrap-nbi/trapsinks",
		snmptrapNbiTrapsinkListJSON)
	result, err := c.GetSnmptrapNbiTrapsinks(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	sink := result["trapsink"].([]any)[0].(map[string]any)
	if sink["name"] != "remote-nms" {
		t.Errorf("name = %v, want remote-nms", sink["name"])
	}
}

func TestGetSnmptrapNbiTrapsink(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/config/snmptrap-nbi/trapsinks/remote-nms",
		snmptrapNbiTrapsinkJSON)
	result, err := c.GetSnmptrapNbiTrapsink(t.Context(), "remote-nms")
	if err != nil {
		t.Fatal(err)
	}
	if result["ipAddress"] != "10.0.0.1" {
		t.Errorf("ipAddress = %v, want 10.0.0.1", result["ipAddress"])
	}
}

func TestCreateSnmptrapNbiTrapsink(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v1Path+"/config/snmptrap-nbi/trapsinks", 201)
	err := c.CreateSnmptrapNbiTrapsink(t.Context(), map[string]any{
		"name": "new-sink", "ipAddress": "10.0.0.2",
		"port": 162, "community": "public",
	})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["name"] != "new-sink" {
		t.Errorf("name = %v, want new-sink", body["name"])
	}
}

func TestUpdateSnmptrapNbiTrapsink(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux,
		"PUT "+v1Path+"/config/snmptrap-nbi/trapsinks/remote-nms", 204)
	err := c.UpdateSnmptrapNbiTrapsink(t.Context(), "remote-nms",
		map[string]string{"port": "163"})
	if err != nil {
		t.Fatal(err)
	}
	if req.body != "port=163" {
		t.Errorf("body = %q, want port=163", req.body)
	}
}

func TestDeleteSnmptrapNbiTrapsink(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux,
		"DELETE "+v1Path+"/config/snmptrap-nbi/trapsinks/remote-nms", 204)
	if err := c.DeleteSnmptrapNbiTrapsink(t.Context(), "remote-nms"); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateSnmptrapNbiConfig(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "POST "+v1Path+"/config/snmptrap-nbi", 204)
	err := c.UpdateSnmptrapNbiConfig(t.Context(), map[string]any{"enabled": false})
	if err != nil {
		t.Fatal(err)
	}
}
