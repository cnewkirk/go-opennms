package opennms

// Tests for the SNMP configuration methods – /rest/snmpConfig.

import (
	"encoding/json"
	"testing"
)

const snmpConfigJSON = `{
	"version": "v2c",
	"port": 161,
	"community": "public",
	"timeout": 1800,
	"retries": 1,
	"maxVarsPerPdu": 10,
	"maxRepetitions": 2,
	"maxRequestSize": 65535,
	"proxyHost": null,
	"securityName": null,
	"location": "Default"
}`

func TestGetSnmpConfig(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/snmpConfig/192.168.1.1", snmpConfigJSON)
	result, err := c.GetSnmpConfig(t.Context(), "192.168.1.1", "")
	if err != nil {
		t.Fatal(err)
	}
	if result["version"] != "v2c" {
		t.Errorf("version = %v", result["version"])
	}
	if result["community"] != "public" {
		t.Errorf("community = %v", result["community"])
	}
	if result["port"].(float64) != 161 {
		t.Errorf("port = %v", result["port"])
	}
}

func TestGetSnmpConfigWithLocation(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/snmpConfig/192.168.1.1", snmpConfigJSON)
	if _, err := c.GetSnmpConfig(t.Context(), "192.168.1.1", "Remote-DC"); err != nil {
		t.Fatal(err)
	}
	if got := req.query.Get("location"); got != "Remote-DC" {
		t.Errorf("location = %q, want Remote-DC", got)
	}
}

func TestSetSnmpConfigV2c(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "PUT "+v1Path+"/snmpConfig/192.168.1.1", 204)
	newConfig := map[string]any{
		"version": "v2c", "community": "private", "port": 161,
	}
	if err := c.SetSnmpConfig(t.Context(), "192.168.1.1", newConfig); err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["version"] != "v2c" {
		t.Errorf("version = %v", body["version"])
	}
	if body["community"] != "private" {
		t.Errorf("community = %v", body["community"])
	}
	if got := req.header.Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q", got)
	}
}

func TestSetSnmpConfigV3(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "PUT "+v1Path+"/snmpConfig/10.0.0.1", 204)
	v3Config := map[string]any{
		"version":        "v3",
		"securityName":   "myUser",
		"authPassphrase": "authSecret",
		"authProtocol":   "SHA",
		"privPassphrase": "privSecret",
		"privProtocol":   "AES128",
		"securityLevel":  3,
	}
	if err := c.SetSnmpConfig(t.Context(), "10.0.0.1", v3Config); err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["version"] != "v3" {
		t.Errorf("version = %v", body["version"])
	}
	if body["authProtocol"] != "SHA" {
		t.Errorf("authProtocol = %v", body["authProtocol"])
	}
	if body["securityLevel"].(float64) != 3 {
		t.Errorf("securityLevel = %v", body["securityLevel"])
	}
}
