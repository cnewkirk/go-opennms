package opennms

// Tests for the perspective poller methods – /api/v2/perspectivepoller.

import "testing"

const perspectivePollerStatusJSON = `{
	"applicationId": 1,
	"services": [
		{
			"serviceId": 201,
			"locations": [
				{"location": "Default", "status": "Up"}
			]
		}
	]
}`

func TestGetPerspectivePollerStatus(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/perspectivepoller/1",
		perspectivePollerStatusJSON)
	result, err := c.GetPerspectivePollerStatus(t.Context(), 1, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if result["applicationId"].(float64) != 1 {
		t.Errorf("applicationId = %v, want 1", result["applicationId"])
	}
}

func TestGetPerspectivePollerStatusWithTime(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v2Path+"/perspectivepoller/1",
		perspectivePollerStatusJSON)
	_, err := c.GetPerspectivePollerStatus(t.Context(), 1, 1000, 2000)
	if err != nil {
		t.Fatal(err)
	}
	if got := req.query.Get("start"); got != "1000" {
		t.Errorf("start = %q, want 1000", got)
	}
	if got := req.query.Get("end"); got != "2000" {
		t.Errorf("end = %q, want 2000", got)
	}
}

func TestGetPerspectivePollerServiceStatus(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/perspectivepoller/1/201",
		perspectivePollerStatusJSON)
	result, err := c.GetPerspectivePollerServiceStatus(t.Context(), 1, 201, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if result["applicationId"].(float64) != 1 {
		t.Errorf("applicationId = %v, want 1", result["applicationId"])
	}
}
