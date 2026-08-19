package opennms

// Tests for the server info method – /rest/info.

import (
	"strings"
	"testing"
)

const infoServerJSON = `{
	"displayVersion": "Horizon 35.0.0",
	"version": "35.0.0",
	"packageName": "opennms",
	"packageDescription": "OpenNMS",
	"ticketerConfig": {"plugin": null, "enabled": false},
	"datetimeformatConfig": {
		"zoneId": "America/New_York",
		"datetimeformat": "yyyy-MM-dd'T'HH:mm:ss.SSSZ"
	},
	"services": {
		"OpenNMS:Name=Pollerd": "running",
		"OpenNMS:Name=Collectd": "running",
		"OpenNMS:Name=Eventd": "running",
		"OpenNMS:Name=Alarmd": "running"
	}
}`

func TestGetInfo(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/info", infoServerJSON)
	result, err := c.GetInfo(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result["version"] != "35.0.0" {
		t.Errorf("version = %v", result["version"])
	}
	if result["displayVersion"] != "Horizon 35.0.0" {
		t.Errorf("displayVersion = %v", result["displayVersion"])
	}
	if result["packageName"] != "opennms" {
		t.Errorf("packageName = %v", result["packageName"])
	}
	if result["ticketerConfig"].(map[string]any)["enabled"] != false {
		t.Errorf("ticketerConfig = %v", result["ticketerConfig"])
	}
	services := result["services"].(map[string]any)
	if services["OpenNMS:Name=Pollerd"] != "running" {
		t.Errorf("services = %v", services)
	}
	if req.method != "GET" {
		t.Errorf("method = %q", req.method)
	}
	if !strings.Contains(req.path, "/rest/info") {
		t.Errorf("path = %q", req.path)
	}
}
