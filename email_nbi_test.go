package opennms

// Tests for the email NBI methods – /rest/config/email-nbi.

import (
	"encoding/json"
	"strings"
	"testing"
)

const emailNbiConfigJSON = `{"enabled": false, "destinations": []}`

const emailNbiDestinationJSON = `{
	"name": "ops-team",
	"firstOccurrenceOnly": true,
	"filters": []
}`

const emailNbiDestinationListJSON = `{
	"destination": [` + emailNbiDestinationJSON + `]
}`

func TestGetEmailNbiConfig(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/config/email-nbi", emailNbiConfigJSON)
	result, err := c.GetEmailNbiConfig(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result["enabled"].(bool) != false {
		t.Errorf("enabled = %v, want false", result["enabled"])
	}
}

func TestGetEmailNbiStatus(t *testing.T) {
	mux, c := newTestClient(t)
	handleText(mux, "GET "+v1Path+"/config/email-nbi/status", "false")
	result, err := c.GetEmailNbiStatus(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result != "false" {
		t.Errorf("status = %q, want false", result)
	}
}

func TestSetEmailNbiStatus(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "PUT "+v1Path+"/config/email-nbi/status", 204)
	if err := c.SetEmailNbiStatus(t.Context(), true); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(req.body, "enabled=true") {
		t.Errorf("body = %q, want to contain enabled=true", req.body)
	}
}

func TestGetEmailNbiDestinations(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/config/email-nbi/destinations",
		emailNbiDestinationListJSON)
	result, err := c.GetEmailNbiDestinations(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	dest := result["destination"].([]any)[0].(map[string]any)
	if dest["name"] != "ops-team" {
		t.Errorf("name = %v, want ops-team", dest["name"])
	}
}

func TestGetEmailNbiDestination(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/config/email-nbi/destinations/ops-team",
		emailNbiDestinationJSON)
	result, err := c.GetEmailNbiDestination(t.Context(), "ops-team")
	if err != nil {
		t.Fatal(err)
	}
	if result["name"] != "ops-team" {
		t.Errorf("name = %v, want ops-team", result["name"])
	}
}

func TestCreateEmailNbiDestination(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v1Path+"/config/email-nbi/destinations", 201)
	data := map[string]any{"name": "new-dest", "firstOccurrenceOnly": true}
	if err := c.CreateEmailNbiDestination(t.Context(), data); err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["name"] != "new-dest" {
		t.Errorf("name = %v, want new-dest", body["name"])
	}
}

func TestUpdateEmailNbiDestination(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "PUT "+v1Path+"/config/email-nbi/destinations/ops-team", 204)
	err := c.UpdateEmailNbiDestination(t.Context(), "ops-team",
		map[string]string{"firstOccurrenceOnly": "false"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteEmailNbiDestination(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/config/email-nbi/destinations/ops-team", 204)
	if err := c.DeleteEmailNbiDestination(t.Context(), "ops-team"); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateEmailNbiConfig(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "POST "+v1Path+"/config/email-nbi", 204)
	err := c.UpdateEmailNbiConfig(t.Context(), map[string]any{"enabled": true})
	if err != nil {
		t.Fatal(err)
	}
}
