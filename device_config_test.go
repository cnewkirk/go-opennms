package opennms

// Tests for the device configuration methods – /rest/device-config.

import (
	"encoding/json"
	"strings"
	"testing"
)

const deviceConfigJSON = `{
	"id": 901,
	"ipInterfaceId": 101,
	"ipAddress": "192.168.1.1",
	"deviceName": "router01.example.com",
	"location": "Default",
	"serviceName": "DeviceConfig-default",
	"configType": "default",
	"createdTime": "2024-06-01T02:00:00.000+0000",
	"lastUpdatedTime": "2024-06-01T02:00:05.000+0000",
	"lastSucceeded": "2024-06-01T02:00:05.000+0000",
	"lastFailed": null,
	"failureReason": null,
	"isLatest": true,
	"backupStatus": "SUCCEEDED"
}`

const deviceConfigListJSON = `{
	"deviceConfigs": [` + deviceConfigJSON + `],
	"totalCount": 1, "count": 1, "offset": 0
}`

func TestGetDeviceConfigsDefaults(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/device-config", deviceConfigListJSON)
	result, err := c.GetDeviceConfigs(t.Context(), nil, "", "", "", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	cfg := result["deviceConfigs"].([]any)[0].(map[string]any)
	if cfg["id"].(float64) != 901 {
		t.Errorf("id = %v, want 901", cfg["id"])
	}
	if req.query.Get("limit") != "10" || req.query.Get("offset") != "0" {
		t.Errorf("query = %v", req.query)
	}
}

func TestGetDeviceConfigsFilters(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/device-config", deviceConfigListJSON)
	_, err := c.GetDeviceConfigs(t.Context(),
		&ListOptions{Limit: 5, Offset: 10, OrderBy: "createdTime", Order: "desc"},
		"router01", "192.168.1.1", "default", 1000000, 2000000)
	if err != nil {
		t.Fatal(err)
	}
	for param, want := range map[string]string{
		"limit": "5", "offset": "10", "orderBy": "createdTime",
		"order": "desc", "deviceName": "router01",
		"ipAddress": "192.168.1.1", "configType": "default",
		"createdAfter": "1000000", "createdBefore": "2000000",
	} {
		if got := req.query.Get(param); got != want {
			t.Errorf("%s = %q, want %q", param, got, want)
		}
	}
}

func TestGetDeviceConfig(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/device-config/901", deviceConfigJSON)
	result, err := c.GetDeviceConfig(t.Context(), 901)
	if err != nil {
		t.Fatal(err)
	}
	if result["id"].(float64) != 901 {
		t.Errorf("id = %v, want 901", result["id"])
	}
	if result["deviceName"] != "router01.example.com" {
		t.Errorf("deviceName = %v", result["deviceName"])
	}
	if req.path != v1Path+"/device-config/901" {
		t.Errorf("path = %q", req.path)
	}
}

func TestGetDeviceConfigByInterface(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/device-config/interface/101",
		deviceConfigListJSON)
	result, err := c.GetDeviceConfigByInterface(t.Context(), 101)
	if err != nil {
		t.Fatal(err)
	}
	cfg := result["deviceConfigs"].([]any)[0].(map[string]any)
	if cfg["ipInterfaceId"].(float64) != 101 {
		t.Errorf("ipInterfaceId = %v, want 101", cfg["ipInterfaceId"])
	}
	if req.path != v1Path+"/device-config/interface/101" {
		t.Errorf("path = %q", req.path)
	}
}

func TestGetLatestDeviceConfigsDefaults(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/device-config/latest",
		deviceConfigListJSON)
	result, err := c.GetLatestDeviceConfigs(t.Context(), nil, "", "")
	if err != nil {
		t.Fatal(err)
	}
	cfg := result["deviceConfigs"].([]any)[0].(map[string]any)
	if cfg["isLatest"] != true {
		t.Errorf("isLatest = %v, want true", cfg["isLatest"])
	}
	if req.query.Get("limit") != "10" || req.query.Get("offset") != "0" {
		t.Errorf("query = %v", req.query)
	}
}

func TestGetLatestDeviceConfigsFilters(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/device-config/latest",
		deviceConfigListJSON)
	_, err := c.GetLatestDeviceConfigs(t.Context(),
		&ListOptions{Limit: 25, OrderBy: "deviceName", Order: "asc"},
		"router", "SUCCEEDED")
	if err != nil {
		t.Fatal(err)
	}
	for param, want := range map[string]string{
		"limit": "25", "orderBy": "deviceName", "order": "asc",
		"search": "router", "status": "SUCCEEDED",
	} {
		if got := req.query.Get(param); got != want {
			t.Errorf("%s = %q, want %q", param, got, want)
		}
	}
}

func TestDownloadDeviceConfigs(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleText(mux, "GET "+v1Path+"/device-config/download",
		"! Cisco IOS config\nhostname router01\n")
	result, err := c.DownloadDeviceConfigs(t.Context(), []int{901, 902})
	if err != nil {
		t.Fatal(err)
	}
	if got := req.query.Get("id"); got != "901,902" {
		t.Errorf("id = %q, want 901,902", got)
	}
	if !strings.Contains(result, "router01") {
		t.Errorf("result = %q, want config text", result)
	}
}

func TestBackupDeviceConfig(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "POST "+v1Path+"/device-config/backup",
		"application/json", 202, `{"status": "submitted"}`)
	payload := []map[string]any{{
		"ipAddress":   "192.168.1.1",
		"location":    "Default",
		"serviceName": "DeviceConfig-default",
		"blocking":    false,
	}}
	result, err := c.BackupDeviceConfig(t.Context(), payload)
	if err != nil {
		t.Fatal(err)
	}
	var body []map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatalf("body = %q: %v", req.body, err)
	}
	if body[0]["ipAddress"] != "192.168.1.1" {
		t.Errorf("ipAddress = %v", body[0]["ipAddress"])
	}
	if body[0]["serviceName"] != "DeviceConfig-default" {
		t.Errorf("serviceName = %v", body[0]["serviceName"])
	}
	if got := req.header.Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q", got)
	}
	if result["status"] != "submitted" {
		t.Errorf("status = %v, want submitted", result["status"])
	}
}
