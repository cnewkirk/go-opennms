package opennms

// Tests for the Secure Credentials Vault methods – /rest/scv.

import (
	"encoding/json"
	"testing"
)

const credentialJSON = `{
	"alias": "my-device",
	"username": "admin",
	"password": "secret",
	"attributes": {}
}`

const credentialListJSON = `{
	"credential": [` + credentialJSON + `],
	"totalCount": 1, "count": 1, "offset": 0
}`

func TestGetCredentials(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/scv", credentialListJSON)
	result, err := c.GetCredentials(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	credential := result["credential"].([]any)[0].(map[string]any)
	if credential["alias"] != "my-device" {
		t.Errorf("alias = %v, want my-device", credential["alias"])
	}
}

func TestGetCredential(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/scv/my-device", credentialJSON)
	result, err := c.GetCredential(t.Context(), "my-device")
	if err != nil {
		t.Fatal(err)
	}
	if result["alias"] != "my-device" {
		t.Errorf("alias = %v, want my-device", result["alias"])
	}
}

func TestCreateCredential(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v1Path+"/scv", 201)
	err := c.CreateCredential(t.Context(), map[string]any{
		"alias": "new-device", "username": "admin", "password": "pass"})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["alias"] != "new-device" {
		t.Errorf("alias = %v, want new-device", body["alias"])
	}
}

func TestUpdateCredential(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "PUT "+v1Path+"/scv/my-device", 204)
	err := c.UpdateCredential(t.Context(), "my-device",
		map[string]any{"username": "root", "password": "new"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteCredential(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/scv/my-device", 204)
	if err := c.DeleteCredential(t.Context(), "my-device"); err != nil {
		t.Fatal(err)
	}
}
