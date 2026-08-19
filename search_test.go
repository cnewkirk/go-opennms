package opennms

// Tests for the search method – /api/v2/search.

import "testing"

const searchResultsJSON = `[{
	"context": "Node",
	"results": [{
		"identifier": "1",
		"label": "router-01",
		"url": "element/node.jsp?node=1",
		"matches": [{"id": "label", "label": "Label", "value": "router-01"}],
		"weight": 0
	}],
	"more": false
}]`

func TestSearch(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v2Path+"/search", searchResultsJSON)
	result, err := c.Search(t.Context(), "router", "", 0)
	if err != nil {
		t.Fatal(err)
	}
	group := result[0].(map[string]any)
	if group["context"] != "Node" {
		t.Errorf("context = %v, want Node", group["context"])
	}
	if got := req.query.Get("_s"); got != "router" {
		t.Errorf("_s = %q, want router", got)
	}
}

func TestSearchWithContextAndLimit(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v2Path+"/search", searchResultsJSON)
	if _, err := c.Search(t.Context(), "router", "Node", 5); err != nil {
		t.Fatal(err)
	}
	if got := req.query.Get("_c"); got != "Node" {
		t.Errorf("_c = %q, want Node", got)
	}
	if got := req.query.Get("_l"); got != "5" {
		t.Errorf("_l = %q, want 5", got)
	}
}

func TestSearchNoMatches(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "GET "+v2Path+"/search", 204)
	result, err := c.Search(t.Context(), "nothing-matches-this", "", 0)
	if err != nil {
		t.Fatal(err)
	}
	if result != nil {
		t.Errorf("result = %v, want nil", result)
	}
}
