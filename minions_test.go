package opennms

// Tests for the minions methods – /rest/minions.

import (
	"testing"
)

const minionJSON = `{
	"id": "minion-01",
	"location": "Default",
	"status": "UP",
	"lastUpdated": "2024-06-01T09:00:00.000+0000"
}`

const minionListJSON = `{
	"minion": [` + minionJSON + `],
	"totalCount": 1, "count": 1, "offset": 0
}`

func TestGetMinions(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/minions", minionListJSON)
	result, err := c.GetMinions(t.Context(), nil)
	if err != nil {
		t.Fatal(err)
	}
	minion := result["minion"].([]any)[0].(map[string]any)
	if minion["id"] != "minion-01" {
		t.Errorf("id = %v, want minion-01", minion["id"])
	}
}

func TestGetMinion(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/minions/minion-01", minionJSON)
	result, err := c.GetMinion(t.Context(), "minion-01")
	if err != nil {
		t.Fatal(err)
	}
	if result["id"] != "minion-01" {
		t.Errorf("id = %v, want minion-01", result["id"])
	}
	if result["location"] != "Default" {
		t.Errorf("location = %v, want Default", result["location"])
	}
}

func TestGetMinionCount(t *testing.T) {
	mux, c := newTestClient(t)
	handleText(mux, "GET "+v1Path+"/minions/count", "2")
	count, err := c.GetMinionCount(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}
}
