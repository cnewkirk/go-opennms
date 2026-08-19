package opennms

// Tests for the user-defined links methods – /api/v2/userdefinedlinks.

import (
	"encoding/json"
	"testing"
)

const userDefinedLinkJSON = `{
	"id": 1,
	"nodeIdA": 1,
	"nodeIdZ": 2,
	"componentLabelA": "eth0",
	"componentLabelZ": "eth0",
	"linkId": "custom-link-1",
	"linkLabel": "Cross connect",
	"owner": "admin"
}`

const userDefinedLinkListJSON = `{
	"user-defined-link": [` + userDefinedLinkJSON + `],
	"totalCount": 1,
	"count": 1,
	"offset": 0
}`

func TestGetUserDefinedLinks(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/userdefinedlinks", userDefinedLinkListJSON)
	result, err := c.GetUserDefinedLinks(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	link := result["user-defined-link"].([]any)[0].(map[string]any)
	if link["linkLabel"] != "Cross connect" {
		t.Errorf("linkLabel = %v, want Cross connect", link["linkLabel"])
	}
}

func TestGetUserDefinedLink(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/userdefinedlinks/1", userDefinedLinkJSON)
	result, err := c.GetUserDefinedLink(t.Context(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if result["id"].(float64) != 1 {
		t.Errorf("id = %v, want 1", result["id"])
	}
}

func TestCreateUserDefinedLink(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "POST "+v2Path+"/userdefinedlinks",
		"application/json", 201, userDefinedLinkJSON)
	result, err := c.CreateUserDefinedLink(t.Context(),
		map[string]any{"nodeIdA": 1, "nodeIdZ": 2, "linkLabel": "Cross connect"})
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
	if body["nodeIdA"].(float64) != 1 {
		t.Errorf("nodeIdA = %v, want 1", body["nodeIdA"])
	}
}

func TestDeleteUserDefinedLink(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v2Path+"/userdefinedlinks/1", 204)
	if err := c.DeleteUserDefinedLink(t.Context(), 1); err != nil {
		t.Fatal(err)
	}
}
