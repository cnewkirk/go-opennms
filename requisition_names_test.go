package opennms

// Tests for the requisition names method – /rest/requisitionNames.

import "testing"

const requisitionNamesJSON = `{
	"foreign-source": ["Routers", "Switches", "Servers"],
	"count": 3
}`

func TestGetRequisitionNames(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/requisitionNames", requisitionNamesJSON)
	result, err := c.GetRequisitionNames(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	names := result["foreign-source"].([]any)
	found := false
	for _, name := range names {
		if name == "Routers" {
			found = true
		}
	}
	if !found {
		t.Errorf("names = %v, want to contain Routers", names)
	}
	if result["count"].(float64) != 3 {
		t.Errorf("count = %v, want 3", result["count"])
	}
}
