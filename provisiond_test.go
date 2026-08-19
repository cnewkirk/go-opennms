package opennms

// Tests for the Provisiond methods – /api/v2/provisiond.

import "testing"

const provisiondStatusJSON = `{
	"status": "RUNNING",
	"jobCount": 0,
	"jobs": []
}`

const provisiondJobStatusJSON = `{
	"id": "job-123",
	"status": "COMPLETED",
	"foreignSource": "Routers",
	"startTime": "2024-06-01T09:00:00.000+0000",
	"endTime": "2024-06-01T09:01:00.000+0000"
}`

func TestGetProvisiondStatus(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/provisiond/status", provisiondStatusJSON)
	result, err := c.GetProvisiondStatus(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result["status"] != "RUNNING" {
		t.Errorf("status = %v, want RUNNING", result["status"])
	}
}

func TestGetProvisiondJobStatus(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/provisiond/status/job-123",
		provisiondJobStatusJSON)
	result, err := c.GetProvisiondJobStatus(t.Context(), "job-123")
	if err != nil {
		t.Fatal(err)
	}
	if result["id"] != "job-123" {
		t.Errorf("id = %v, want job-123", result["id"])
	}
	if result["status"] != "COMPLETED" {
		t.Errorf("status = %v, want COMPLETED", result["status"])
	}
}
