package opennms

// Tests for the health methods – /rest/health.

import "testing"

const healthJSON = `{
	"healthy": true,
	"responses": [
		{
			"description": "Ensuring installed bundles are started",
			"status": "Success"
		}
	]
}`

const healthProbeJSON = `{"status": "Everything is awesome"}`

func TestGetHealth(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/health", healthJSON)
	result, err := c.GetHealth(t.Context(), "")
	if err != nil {
		t.Fatal(err)
	}
	if result["healthy"] != true {
		t.Errorf("healthy = %v, want true", result["healthy"])
	}
}

func TestGetHealthWithTag(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/health", healthJSON)
	if _, err := c.GetHealth(t.Context(), "bundle"); err != nil {
		t.Fatal(err)
	}
	if got := req.query.Get("tag"); got != "bundle" {
		t.Errorf("tag = %q, want bundle", got)
	}
}

func TestGetHealthProbe(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/health/probe", healthProbeJSON)
	result, err := c.GetHealthProbe(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result["status"] != "Everything is awesome" {
		t.Errorf("status = %v", result["status"])
	}
}
