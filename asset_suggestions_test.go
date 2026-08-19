package opennms

// Tests for the asset suggestions method – /rest/assets/suggestions.

import "testing"

const assetSuggestionsJSON = `{
	"assetSuggestion": [
		{"column": "manufacturer", "values": ["Cisco", "Juniper"]},
		{"column": "operatingSystem", "values": ["IOS-XE", "JUNOS"]}
	]
}`

func TestGetAssetSuggestions(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/assets/suggestions", assetSuggestionsJSON)
	result, err := c.GetAssetSuggestions(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	suggestion := result["assetSuggestion"].([]any)[0].(map[string]any)
	if suggestion["column"] != "manufacturer" {
		t.Errorf("column = %v, want manufacturer", suggestion["column"])
	}
	values := suggestion["values"].([]any)
	if values[0] != "Cisco" {
		t.Errorf("values = %v, want Cisco first", values)
	}
}
