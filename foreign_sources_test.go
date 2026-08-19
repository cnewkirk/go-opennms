package opennms

// Tests for the foreign sources methods – /rest/foreignSources.

import (
	"encoding/json"
	"strings"
	"testing"
)

const foreignSourceDetectorJSON = `{
	"name": "ICMP",
	"class": "org.opennms.netmgt.provision.detector.icmp.IcmpDetector",
	"parameter": []
}`

const foreignSourcePolicyJSON = `{
	"name": "Do Not Persist Discovered IPs",
	"class": "org.opennms.netmgt.provision.persist.policies.MatchingIpInterfacePolicy",
	"parameter": [{"key": "action", "value": "DO_NOT_PERSIST"}]
}`

const foreignSourceJSON = `{
	"name": "Routers",
	"date-stamp": "2024-05-20T14:00:00.000+0000",
	"scan-interval": "1d",
	"detectors": [` + foreignSourceDetectorJSON + `],
	"policies": [` + foreignSourcePolicyJSON + `]
}`

const foreignSourceListJSON = `{
	"count": 1, "totalCount": 1,
	"foreignSource": [` + foreignSourceJSON + `]
}`

func TestGetForeignSources(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/foreignSources", foreignSourceListJSON)
	result, err := c.GetForeignSources(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	fs := result["foreignSource"].([]any)[0].(map[string]any)
	if fs["name"] != "Routers" {
		t.Errorf("name = %v, want Routers", fs["name"])
	}
}

func TestGetForeignSource(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/foreignSources/Routers", foreignSourceJSON)
	result, err := c.GetForeignSource(t.Context(), "Routers")
	if err != nil {
		t.Fatal(err)
	}
	if result["name"] != "Routers" {
		t.Errorf("name = %v, want Routers", result["name"])
	}
	if result["scan-interval"] != "1d" {
		t.Errorf("scan-interval = %v, want 1d", result["scan-interval"])
	}
}

func TestGetDefaultForeignSource(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/foreignSources/default", foreignSourceJSON)
	result, err := c.GetDefaultForeignSource(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result["name"] != "Routers" {
		t.Errorf("name = %v, want Routers", result["name"])
	}
	if !strings.Contains(req.path, "/foreignSources/default") {
		t.Errorf("path = %q", req.path)
	}
}

func TestGetDeployedForeignSources(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/foreignSources/deployed", foreignSourceListJSON)
	result, err := c.GetDeployedForeignSources(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	fs := result["foreignSource"].([]any)[0].(map[string]any)
	if fs["name"] != "Routers" {
		t.Errorf("name = %v, want Routers", fs["name"])
	}
}

func TestGetDeployedForeignSourceCount(t *testing.T) {
	mux, c := newTestClient(t)
	handleText(mux, "GET "+v1Path+"/foreignSources/deployed/count", "4")
	count, err := c.GetDeployedForeignSourceCount(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Errorf("count = %d, want 4", count)
	}
}

func TestCreateForeignSource(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "POST "+v1Path+"/foreignSources", foreignSourceJSON)
	var fs map[string]any
	if err := json.Unmarshal([]byte(foreignSourceJSON), &fs); err != nil {
		t.Fatal(err)
	}
	if _, err := c.CreateForeignSource(t.Context(), fs); err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["name"] != "Routers" {
		t.Errorf("name = %v, want Routers", body["name"])
	}
	if got := req.header.Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
}

func TestUpdateForeignSource(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "PUT "+v1Path+"/foreignSources/Routers", formType, 204, "")
	err := c.UpdateForeignSource(t.Context(), "Routers",
		map[string]any{"scan-interval": "12h"})
	if err != nil {
		t.Fatal(err)
	}
	if req.body != "scan-interval=12h" {
		t.Errorf("body = %q, want scan-interval=12h", req.body)
	}
}

func TestDeleteForeignSource(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/foreignSources/Routers", 202)
	if err := c.DeleteForeignSource(t.Context(), "Routers"); err != nil {
		t.Fatal(err)
	}
}

func TestGetForeignSourceDetectors(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/foreignSources/Routers/detectors",
		`{"detector": [`+foreignSourceDetectorJSON+`]}`)
	result, err := c.GetForeignSourceDetectors(t.Context(), "Routers")
	if err != nil {
		t.Fatal(err)
	}
	detector := result["detector"].([]any)[0].(map[string]any)
	if detector["name"] != "ICMP" {
		t.Errorf("name = %v, want ICMP", detector["name"])
	}
}

func TestGetForeignSourceDetector(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/foreignSources/Routers/detectors/ICMP",
		foreignSourceDetectorJSON)
	result, err := c.GetForeignSourceDetector(t.Context(), "Routers", "ICMP")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(result["class"].(string), "IcmpDetector") {
		t.Errorf("class = %v", result["class"])
	}
}

func TestAddForeignSourceDetector(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "POST "+v1Path+"/foreignSources/Routers/detectors",
		foreignSourceDetectorJSON)
	var detector map[string]any
	if err := json.Unmarshal([]byte(foreignSourceDetectorJSON), &detector); err != nil {
		t.Fatal(err)
	}
	if _, err := c.AddForeignSourceDetector(t.Context(), "Routers", detector); err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["name"] != "ICMP" {
		t.Errorf("name = %v, want ICMP", body["name"])
	}
}

func TestDeleteForeignSourceDetector(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/foreignSources/Routers/detectors/ICMP", 204)
	if err := c.DeleteForeignSourceDetector(t.Context(), "Routers", "ICMP"); err != nil {
		t.Fatal(err)
	}
}

func TestGetForeignSourcePolicies(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/foreignSources/Routers/policies",
		`{"policy": [`+foreignSourcePolicyJSON+`]}`)
	result, err := c.GetForeignSourcePolicies(t.Context(), "Routers")
	if err != nil {
		t.Fatal(err)
	}
	policy := result["policy"].([]any)[0].(map[string]any)
	if policy["name"] != "Do Not Persist Discovered IPs" {
		t.Errorf("name = %v", policy["name"])
	}
}

func TestGetForeignSourcePolicy(t *testing.T) {
	mux, c := newTestClient(t)
	// requests percent-encodes spaces in path segments as %20, not +;
	// the mux decodes %20 to spaces (a literal "+" would stay "+").
	req := handleJSON(mux, "GET "+v1Path+"/foreignSources/Routers/policies/{policy}",
		foreignSourcePolicyJSON)
	result, err := c.GetForeignSourcePolicy(t.Context(),
		"Routers", "Do Not Persist Discovered IPs")
	if err != nil {
		t.Fatal(err)
	}
	if result["name"] != "Do Not Persist Discovered IPs" {
		t.Errorf("name = %v", result["name"])
	}
	want := v1Path + "/foreignSources/Routers/policies/Do Not Persist Discovered IPs"
	if req.path != want {
		t.Errorf("path = %q, want %q", req.path, want)
	}
}

func TestAddForeignSourcePolicy(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "POST "+v1Path+"/foreignSources/Routers/policies",
		foreignSourcePolicyJSON)
	var policy map[string]any
	if err := json.Unmarshal([]byte(foreignSourcePolicyJSON), &policy); err != nil {
		t.Fatal(err)
	}
	if _, err := c.AddForeignSourcePolicy(t.Context(), "Routers", policy); err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["name"] != "Do Not Persist Discovered IPs" {
		t.Errorf("name = %v", body["name"])
	}
}

func TestDeleteForeignSourcePolicy(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "DELETE "+v1Path+"/foreignSources/Routers/policies/{policy}", 204)
	err := c.DeleteForeignSourcePolicy(t.Context(),
		"Routers", "Do Not Persist Discovered IPs")
	if err != nil {
		t.Fatal(err)
	}
	want := v1Path + "/foreignSources/Routers/policies/Do Not Persist Discovered IPs"
	if req.path != want {
		t.Errorf("path = %q, want %q", req.path, want)
	}
}
