package opennms

// Tests for the logs methods – /rest/logs.

import (
	"strings"
	"testing"
)

const logsFileListJSON = `["manager.log", "web.log", "provisiond.log"]`

const logsContents = "2026-08-17 06:00:01 INFO  [Main] Manager started\n" +
	"2026-08-17 06:00:02 INFO  [Main] Startup complete"

func TestGetLogFiles(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/logs", logsFileListJSON)
	result, err := c.GetLogFiles(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, f := range result {
		if f == "manager.log" {
			found = true
		}
	}
	if !found {
		t.Errorf("result = %v, want manager.log in it", result)
	}
}

func TestGetLogContents(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/logs/contents", logsContents)
	reverse := false
	result, err := c.GetLogContents(t.Context(), "manager.log", 100, &reverse)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, "Manager started") {
		t.Errorf("result = %q", result)
	}
	if got := req.query.Get("f"); got != "manager.log" {
		t.Errorf("f = %q", got)
	}
	if got := req.query.Get("n"); got != "100" {
		t.Errorf("n = %q", got)
	}
	if got := req.query.Get("reverse"); got != "False" {
		t.Errorf("reverse = %q, want False", got)
	}
}

func TestGetLogContentsMissingFile(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "GET "+v1Path+"/logs/contents", 204)
	result, err := c.GetLogContents(t.Context(), "nope.log", 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if result != "" {
		t.Errorf("result = %q, want empty", result)
	}
}
