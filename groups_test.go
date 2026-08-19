package opennms

// Tests for the groups methods – /rest/groups.

import "testing"

const groupJSON = `{
	"name": "network-ops",
	"comments": "Network operations team",
	"users": ["admin", "jsmith"],
	"categories": ["Production"]
}`

const groupListJSON = `{
	"group": [` + groupJSON + `],
	"totalCount": 1, "count": 1, "offset": 0
}`

const groupUserListJSON = `{"users": ["admin", "jsmith"]}`

const groupCategoryListJSON = `{"categories": ["Production", "Routers"]}`

func TestGetGroups(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/groups", groupListJSON)
	result, err := c.GetGroups(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	group := result["group"].([]any)[0].(map[string]any)
	if group["name"] != "network-ops" {
		t.Errorf("name = %v, want network-ops", group["name"])
	}
}

func TestGetGroup(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/groups/network-ops", groupJSON)
	result, err := c.GetGroup(t.Context(), "network-ops")
	if err != nil {
		t.Fatal(err)
	}
	if result["name"] != "network-ops" {
		t.Errorf("name = %v, want network-ops", result["name"])
	}
	if result["comments"] != "Network operations team" {
		t.Errorf("comments = %v", result["comments"])
	}
}

func TestCreateGroup(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "POST "+v1Path+"/groups", xmlType, 201, groupJSON)
	_, err := c.CreateGroup(t.Context(), map[string]any{
		"name": "network-ops", "comments": "Network operations team"})
	if err != nil {
		t.Fatal(err)
	}
	want := "<group><name>network-ops</name>" +
		"<comments>Network operations team</comments></group>"
	if req.body != want {
		t.Errorf("body = %q, want %q", req.body, want)
	}
}

func TestCreateGroupWithMembers(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "POST "+v1Path+"/groups", xmlType, 201, "")
	_, err := c.CreateGroup(t.Context(), map[string]any{
		"name": "ops", "user": []string{"alice", "bob"}})
	if err != nil {
		t.Fatal(err)
	}
	want := "<group><name>ops</name>" +
		"<user>alice</user><user>bob</user></group>"
	if req.body != want {
		t.Errorf("body = %q, want %q", req.body, want)
	}
}

func TestUpdateGroup(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "PUT "+v1Path+"/groups/network-ops", formType, 204, "")
	err := c.UpdateGroup(t.Context(), "network-ops",
		map[string]string{"comments": "Updated comment"})
	if err != nil {
		t.Fatal(err)
	}
	if req.body != "comments=Updated+comment" {
		t.Errorf("body = %q, want comments=Updated+comment", req.body)
	}
}

func TestDeleteGroup(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/groups/network-ops", 202)
	if err := c.DeleteGroup(t.Context(), "network-ops"); err != nil {
		t.Fatal(err)
	}
}

func TestGetGroupUsers(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/groups/network-ops/users", groupUserListJSON)
	result, err := c.GetGroupUsers(t.Context(), "network-ops")
	if err != nil {
		t.Fatal(err)
	}
	users := result["users"].([]any)
	if users[0] != "admin" {
		t.Errorf("users = %v, want admin first", users)
	}
}

func TestAddUserToGroup(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "PUT "+v1Path+"/groups/network-ops/users/jsmith", 204)
	err := c.AddUserToGroup(t.Context(), "network-ops", "jsmith")
	if err != nil {
		t.Fatal(err)
	}
	if req.path != v1Path+"/groups/network-ops/users/jsmith" {
		t.Errorf("path = %q", req.path)
	}
}

func TestRemoveUserFromGroup(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/groups/network-ops/users/jsmith", 204)
	err := c.RemoveUserFromGroup(t.Context(), "network-ops", "jsmith")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetGroupCategories(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/groups/network-ops/categories",
		groupCategoryListJSON)
	result, err := c.GetGroupCategories(t.Context(), "network-ops")
	if err != nil {
		t.Fatal(err)
	}
	categories := result["categories"].([]any)
	if categories[0] != "Production" {
		t.Errorf("categories = %v, want Production first", categories)
	}
}

func TestAddCategoryToGroup(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "PUT "+v1Path+"/groups/network-ops/categories/Production", 204)
	err := c.AddCategoryToGroup(t.Context(), "network-ops", "Production")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRemoveCategoryFromGroup(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/groups/network-ops/categories/Production", 204)
	err := c.RemoveCategoryFromGroup(t.Context(), "network-ops", "Production")
	if err != nil {
		t.Fatal(err)
	}
}
