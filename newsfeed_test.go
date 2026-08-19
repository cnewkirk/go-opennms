package opennms

// Tests for the news feed method – /api/v2/newsfeed.

import "testing"

const newsfeedJSON = `{
	"items": [
		{
			"title": "OpenNMS Meridian 2025 released",
			"link": "https://www.opennms.com/news/meridian-2025",
			"categories": ["Releases"],
			"tags": ["meridian"],
			"description": "Meridian 2025 is now available.",
			"shortDescription": "Meridian 2025 is out."
		}
	]
}`

func TestGetNewsfeed(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/newsfeed", newsfeedJSON)
	result, err := c.GetNewsfeed(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	item := result["items"].([]any)[0].(map[string]any)
	categories := item["categories"].([]any)
	if len(categories) != 1 || categories[0] != "Releases" {
		t.Errorf("categories = %v, want [Releases]", categories)
	}
}
