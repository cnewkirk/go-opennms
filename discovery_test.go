package opennms

// Tests for the discovery method – /api/v2/discovery.

import (
	"encoding/json"
	"testing"
)

func discoveryBody(t *testing.T, raw string) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.Unmarshal([]byte(raw), &body); err != nil {
		t.Fatalf("parsing body %q: %v", raw, err)
	}
	return body
}

func TestDiscoverSpecifics(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v2Path+"/discovery", 204)
	err := c.Discover(t.Context(), map[string]any{
		"specifics": []map[string]any{{
			"ip":            "10.0.0.1",
			"location":      "Default",
			"retries":       1,
			"timeout":       2000,
			"foreignSource": "Routers",
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if req.method != "POST" {
		t.Errorf("method = %q, want POST", req.method)
	}
	if req.path != v2Path+"/discovery" {
		t.Errorf("path = %q", req.path)
	}
	body := discoveryBody(t, req.body)
	specific := body["specifics"].([]any)[0].(map[string]any)
	if specific["ip"] != "10.0.0.1" {
		t.Errorf("ip = %v, want 10.0.0.1", specific["ip"])
	}
	if specific["foreignSource"] != "Routers" {
		t.Errorf("foreignSource = %v, want Routers", specific["foreignSource"])
	}
	if got := req.header.Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
}

func TestDiscoverWithRanges(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v2Path+"/discovery", 204)
	err := c.Discover(t.Context(), map[string]any{
		"include_ranges": []map[string]any{{
			"begin": "10.0.1.1", "end": "10.0.1.254",
			"location": "Default", "retries": 1, "timeout": 2000,
		}},
		"exclude_ranges": []map[string]any{{
			"begin": "10.0.1.100", "end": "10.0.1.110",
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	body := discoveryBody(t, req.body)
	include := body["include_ranges"].([]any)[0].(map[string]any)
	if include["begin"] != "10.0.1.1" {
		t.Errorf("begin = %v, want 10.0.1.1", include["begin"])
	}
	exclude := body["exclude_ranges"].([]any)[0].(map[string]any)
	if exclude["end"] != "10.0.1.110" {
		t.Errorf("end = %v, want 10.0.1.110", exclude["end"])
	}
}

func TestDiscoverWithURL(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v2Path+"/discovery", 204)
	err := c.Discover(t.Context(), map[string]any{
		"include_urls": []map[string]any{{
			"url":      "file:/opt/opennms/etc/include.txt",
			"location": "Default",
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	body := discoveryBody(t, req.body)
	include := body["include_urls"].([]any)[0].(map[string]any)
	if include["location"] != "Default" {
		t.Errorf("location = %v, want Default", include["location"])
	}
}
