package opennms

// Tests for the classifications methods – /rest/classifications.

import (
	"encoding/json"
	"testing"
)

const classificationRuleJSON = `{
	"id": 1,
	"name": "HTTPS",
	"dstPort": "443",
	"protocol": "tcp",
	"groupId": 1,
	"position": 0
}`

const classificationRuleListJSON = `{
	"classificationRule": [` + classificationRuleJSON + `],
	"totalCount": 1, "count": 1, "offset": 0
}`

const classificationGroupJSON = `{
	"id": 1,
	"name": "default",
	"enabled": true,
	"readOnly": false,
	"position": 0
}`

const classificationGroupListJSON = `{
	"classificationGroup": [` + classificationGroupJSON + `],
	"totalCount": 1, "count": 1, "offset": 0
}`

const classificationClassifyResultJSON = `{"classification": "HTTPS"}`

const classificationProtocolsJSON = `["TCP", "UDP", "ICMP"]`

func TestGetClassificationRules(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/classifications", classificationRuleListJSON)
	result, err := c.GetClassificationRules(t.Context(), nil)
	if err != nil {
		t.Fatal(err)
	}
	rule := result["classificationRule"].([]any)[0].(map[string]any)
	if rule["name"] != "HTTPS" {
		t.Errorf("name = %v, want HTTPS", rule["name"])
	}
}

func TestGetClassificationRule(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/classifications/1", classificationRuleJSON)
	result, err := c.GetClassificationRule(t.Context(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if result["id"].(float64) != 1 {
		t.Errorf("id = %v, want 1", result["id"])
	}
}

func TestCreateClassificationRule(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v1Path+"/classifications", 201)
	rule := map[string]any{"name": "SSH", "dstPort": "22", "protocol": "tcp"}
	if err := c.CreateClassificationRule(t.Context(), rule); err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["name"] != "SSH" {
		t.Errorf("name = %v, want SSH", body["name"])
	}
}

func TestUpdateClassificationRule(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "PUT "+v1Path+"/classifications/1", 204)
	err := c.UpdateClassificationRule(t.Context(), 1,
		map[string]any{"name": "HTTPS", "dstPort": "443"})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["dstPort"] != "443" {
		t.Errorf("dstPort = %v, want 443", body["dstPort"])
	}
}

func TestDeleteClassificationRulesByGroup(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "DELETE "+v1Path+"/classifications", 204)
	if err := c.DeleteClassificationRules(t.Context(), 1); err != nil {
		t.Fatal(err)
	}
	if got := req.query.Get("groupId"); got != "1" {
		t.Errorf("groupId = %q, want 1", got)
	}
}

func TestDeleteClassificationRule(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/classifications/1", 204)
	if err := c.DeleteClassificationRule(t.Context(), 1); err != nil {
		t.Fatal(err)
	}
}

func TestClassify(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "POST "+v1Path+"/classifications/classify",
		classificationClassifyResultJSON)
	result, err := c.Classify(t.Context(),
		map[string]any{"dstPort": "443", "protocol": "tcp"})
	if err != nil {
		t.Fatal(err)
	}
	if result["classification"] != "HTTPS" {
		t.Errorf("classification = %v, want HTTPS", result["classification"])
	}
}

func TestGetClassificationGroups(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/classifications/groups",
		classificationGroupListJSON)
	result, err := c.GetClassificationGroups(t.Context(), nil)
	if err != nil {
		t.Fatal(err)
	}
	group := result["classificationGroup"].([]any)[0].(map[string]any)
	if group["name"] != "default" {
		t.Errorf("name = %v, want default", group["name"])
	}
}

func TestGetClassificationGroup(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/classifications/groups/1",
		classificationGroupJSON)
	result, err := c.GetClassificationGroup(t.Context(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if result["id"].(float64) != 1 {
		t.Errorf("id = %v, want 1", result["id"])
	}
}

func TestCreateClassificationGroup(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v1Path+"/classifications/groups", 201)
	err := c.CreateClassificationGroup(t.Context(),
		map[string]any{"name": "custom"})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["name"] != "custom" {
		t.Errorf("name = %v, want custom", body["name"])
	}
}

func TestUpdateClassificationGroup(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "PUT "+v1Path+"/classifications/groups/1", 204)
	err := c.UpdateClassificationGroup(t.Context(), 1,
		map[string]any{"name": "renamed"})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["name"] != "renamed" {
		t.Errorf("name = %v, want renamed", body["name"])
	}
}

func TestDeleteClassificationGroup(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/classifications/groups/1", 204)
	if err := c.DeleteClassificationGroup(t.Context(), 1); err != nil {
		t.Fatal(err)
	}
}

func TestImportClassificationRules(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v1Path+"/classifications/groups/1", 204)
	csv := "name;dstPort;protocol\\nSSH;22;tcp"
	if err := c.ImportClassificationRules(t.Context(), 1, csv); err != nil {
		t.Fatal(err)
	}
	if got := req.header.Get("Content-Type"); got != "text/comma-separated-values" {
		t.Errorf("Content-Type = %q, want text/comma-separated-values", got)
	}
	if req.body != csv {
		t.Errorf("body = %q, want %q", req.body, csv)
	}
}

func TestGetClassificationProtocols(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/classifications/protocols",
		classificationProtocolsJSON)
	result, err := c.GetClassificationProtocols(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, p := range result {
		if p == "TCP" {
			found = true
		}
	}
	if !found {
		t.Errorf("result = %v, want to contain TCP", result)
	}
}
