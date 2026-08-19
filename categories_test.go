package opennms

// Tests for the categories methods – /rest/categories.

import (
	"encoding/json"
	"testing"
)

const categoryJSON = `{
	"id": 2,
	"name": "Production",
	"authorizedGroups": ["network-ops"]
}`

const categoryListJSON = `{
	"category": [` + categoryJSON + `],
	"totalCount": 1, "count": 1, "offset": 0
}`

func TestGetCategories(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/categories", categoryListJSON)
	result, err := c.GetCategories(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	category := result["category"].([]any)[0].(map[string]any)
	if category["name"] != "Production" {
		t.Errorf("name = %v, want Production", category["name"])
	}
}

func TestGetCategory(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/categories/Production", categoryJSON)
	result, err := c.GetCategory(t.Context(), "Production")
	if err != nil {
		t.Fatal(err)
	}
	if result["id"].(float64) != 2 {
		t.Errorf("id = %v, want 2", result["id"])
	}
	if result["name"] != "Production" {
		t.Errorf("name = %v, want Production", result["name"])
	}
}

func TestCreateCategory(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "POST "+v1Path+"/categories", categoryJSON)
	_, err := c.CreateCategory(t.Context(),
		map[string]any{"name": "Production", "authorizedGroups": []any{}})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["name"] != "Production" {
		t.Errorf("body name = %v, want Production", body["name"])
	}
}

func TestUpdateCategory(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "PUT "+v1Path+"/categories/Production", formType, 204, "")
	err := c.UpdateCategory(t.Context(), "Production",
		map[string]string{"description": "Prod nodes"})
	if err != nil {
		t.Fatal(err)
	}
	if req.body != "description=Prod+nodes" {
		t.Errorf("body = %q, want description=Prod+nodes", req.body)
	}
}

func TestAssociateCategoryWithNode(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "PUT "+v1Path+"/categories/Production/nodes/1",
		xmlType, 201, "")
	if err := c.AssociateCategoryWithNode(t.Context(), "Production", 1); err != nil {
		t.Fatal(err)
	}
	if req.body != `<category name="Production"/>` {
		t.Errorf("body = %q", req.body)
	}
}

func TestDeleteCategory(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/categories/Production", 204)
	if err := c.DeleteCategory(t.Context(), "Production"); err != nil {
		t.Fatal(err)
	}
}

func TestGetNodeCategoriesList(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/categories/nodes/1", categoryListJSON)
	result, err := c.GetNodeCategoriesList(t.Context(), 1)
	if err != nil {
		t.Fatal(err)
	}
	category := result["category"].([]any)[0].(map[string]any)
	if category["name"] != "Production" {
		t.Errorf("name = %v, want Production", category["name"])
	}
	if req.path != v1Path+"/categories/nodes/1" {
		t.Errorf("path = %q", req.path)
	}
}

func TestGetCategoryForNode(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/categories/Production/nodes/1", categoryJSON)
	result, err := c.GetCategoryForNode(t.Context(), "Production", 1)
	if err != nil {
		t.Fatal(err)
	}
	if result["name"] != "Production" {
		t.Errorf("name = %v, want Production", result["name"])
	}
}

func TestDissociateCategoryFromNode(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/categories/Production/nodes/1", 204)
	if err := c.DissociateCategoryFromNode(t.Context(), "Production", 1); err != nil {
		t.Fatal(err)
	}
}

func TestGetCategoriesForGroup(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/categories/groups/network-ops",
		categoryListJSON)
	_, err := c.GetCategoriesForGroup(t.Context(), "network-ops")
	if err != nil {
		t.Fatal(err)
	}
	if req.path != v1Path+"/categories/groups/network-ops" {
		t.Errorf("path = %q", req.path)
	}
}

func TestAssociateCategoryWithGroup(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "PUT "+v1Path+"/categories/Production/groups/network-ops", 204)
	err := c.AssociateCategoryWithGroup(t.Context(), "Production", "network-ops")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDissociateCategoryFromGroup(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/categories/Production/groups/network-ops", 204)
	err := c.DissociateCategoryFromGroup(t.Context(), "Production", "network-ops")
	if err != nil {
		t.Fatal(err)
	}
}
