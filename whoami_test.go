package opennms

// Tests for the whoami method – /rest/whoami.

import "testing"

const whoamiJSON = `{
	"fullName": "Administrator",
	"id": "admin",
	"roles": ["ROLE_ADMIN", "ROLE_REST"]
}`

func TestGetWhoami(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/whoami", whoamiJSON)
	result, err := c.GetWhoami(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result["id"] != "admin" {
		t.Errorf("id = %v, want admin", result["id"])
	}
	roles := result["roles"].([]any)
	found := false
	for _, r := range roles {
		if r == "ROLE_ADMIN" {
			found = true
		}
	}
	if !found {
		t.Errorf("roles = %v, want ROLE_ADMIN present", roles)
	}
}
