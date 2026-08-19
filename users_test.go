package opennms

// Tests for the users methods – /rest/users.

import (
	"strings"
	"testing"
)

const usersUserJSON = `{
	"user-id": "jsmith",
	"full-name": "Jane Smith",
	"user-comments": "Senior Network Engineer",
	"email": "jsmith@example.com",
	"password": null,
	"password-salt": true,
	"duty-schedule": ["MoTuWeThFrSaSu800-2300"],
	"roles": ["ROLE_USER"]
}`

const usersUserListJSON = `{
	"user": [` + usersUserJSON + `],
	"totalCount": 1, "count": 1, "offset": 0
}`

func usersNewUser() map[string]any {
	return map[string]any{
		"user-id":   "newuser",
		"full-name": "New User",
		"password":  "s3cret",
		"email":     "newuser@example.com",
	}
}

func TestGetUsers(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/users", usersUserListJSON)
	result, err := c.GetUsers(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result["user"].([]any)[0].(map[string]any)["user-id"] != "jsmith" {
		t.Errorf("result = %v", result)
	}
}

func TestGetUser(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/users/jsmith", usersUserJSON)
	result, err := c.GetUser(t.Context(), "jsmith")
	if err != nil {
		t.Fatal(err)
	}
	if result["user-id"] != "jsmith" {
		t.Errorf("user-id = %v", result["user-id"])
	}
	if result["full-name"] != "Jane Smith" {
		t.Errorf("full-name = %v", result["full-name"])
	}
}

func TestCreateUser(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "POST "+v1Path+"/users", xmlType, 201, usersUserJSON)
	if _, err := c.CreateUser(t.Context(), usersNewUser(), false); err != nil {
		t.Fatal(err)
	}
	want := "<user><user-id>newuser</user-id>" +
		"<full-name>New User</full-name>" +
		"<email>newuser@example.com</email>" +
		"<password>s3cret</password></user>"
	if req.body != want {
		t.Errorf("body = %q, want %q", req.body, want)
	}
}

func TestCreateUserHashPassword(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "POST "+v1Path+"/users", xmlType, 201, usersUserJSON)
	if _, err := c.CreateUser(t.Context(), usersNewUser(), true); err != nil {
		t.Fatal(err)
	}
	if got := req.query.Get("hashPassword"); got != "true" {
		t.Errorf("hashPassword = %q, want true", got)
	}
}

func TestCreateUserRolesAndSalt(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "POST "+v1Path+"/users", xmlType, 201, "")
	_, err := c.CreateUser(t.Context(), map[string]any{
		"user-id":      "oncall",
		"password":     "hashed",
		"passwordSalt": true,
		"role":         []string{"ROLE_ADMIN"},
	}, false)
	if err != nil {
		t.Fatal(err)
	}
	want := "<user><user-id>oncall</user-id>" +
		"<password>hashed</password>" +
		"<passwordSalt>true</passwordSalt>" +
		"<role>ROLE_ADMIN</role></user>"
	if req.body != want {
		t.Errorf("body = %q, want %q", req.body, want)
	}
}

func TestUpdateUser(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "PUT "+v1Path+"/users/jsmith", formType, 204, "")
	err := c.UpdateUser(t.Context(), "jsmith",
		map[string]string{"fullName": "Jane A. Smith"})
	if err != nil {
		t.Fatal(err)
	}
	if req.body != "fullName=Jane+A.+Smith" {
		t.Errorf("body = %q, want fullName=Jane+A.+Smith", req.body)
	}
}

func TestDeleteUser(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/users/jsmith", 202)
	if err := c.DeleteUser(t.Context(), "jsmith"); err != nil {
		t.Fatal(err)
	}
}

func TestAssignRoleToUser(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "PUT "+v1Path+"/users/jsmith/roles/ROLE_ADMIN", 204)
	if err := c.AssignRoleToUser(t.Context(), "jsmith", "ROLE_ADMIN"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(req.path, "/users/jsmith/roles/ROLE_ADMIN") {
		t.Errorf("path = %q", req.path)
	}
}

func TestRevokeRoleFromUser(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/users/jsmith/roles/ROLE_ADMIN", 204)
	if err := c.RevokeRoleFromUser(t.Context(), "jsmith", "ROLE_ADMIN"); err != nil {
		t.Fatal(err)
	}
}
