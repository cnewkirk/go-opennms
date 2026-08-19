package opennms

// Tests for the eventconf methods – /api/v2/eventconf.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const eventconfFilterJSON = `{
	"events": [
		{"uei": "uei.opennms.org/nodes/nodeDown", "label": "Node Down",
		 "source": "default"}
	]
}`

const eventconfSourcesJSON = `{
	"sources": [{"id": "default", "label": "Default Events"}]
}`

const eventconfSourceJSON = `{
	"id": "default", "label": "Default Events", "eventCount": 150
}`

const eventconfSourceNamesJSON = `["default", "Cisco", "Juniper"]`

const eventconfVendorEventsJSON = `{
	"events": [
		{"uei": "uei.opennms.org/vendor/Cisco/traps/bgpUp", "label": "BGP Up"}
	]
}`

const eventconfEventJSON = `{
	"uei": "uei.opennms.org/custom/testEvent",
	"label": "Test Event",
	"descr": "A test event",
	"logmsg": {"content": "Test event fired", "dest": "logndisplay"},
	"severity": "Warning"
}`

func eventconfEventBody(t *testing.T) map[string]any {
	t.Helper()
	var event map[string]any
	if err := json.Unmarshal([]byte(eventconfEventJSON), &event); err != nil {
		t.Fatal(err)
	}
	return event
}

func TestGetEventconfFilter(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v2Path+"/eventconf/filter", eventconfFilterJSON)
	result, err := c.GetEventconfFilter(t.Context(), "nodeDown", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	event := result["events"].([]any)[0].(map[string]any)
	if event["uei"] != "uei.opennms.org/nodes/nodeDown" {
		t.Errorf("uei = %v", event["uei"])
	}
	if got := req.query.Get("uei"); got != "nodeDown" {
		t.Errorf("uei param = %q, want nodeDown", got)
	}
}

func TestGetEventconfFilterSources(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/eventconf/filter/sources", eventconfSourcesJSON)
	result, err := c.GetEventconfFilterSources(t.Context(), "", nil)
	if err != nil {
		t.Fatal(err)
	}
	source := result["sources"].([]any)[0].(map[string]any)
	if source["id"] != "default" {
		t.Errorf("id = %v, want default", source["id"])
	}
}

func TestGetEventconfFilterEvents(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/eventconf/filter/default/events",
		eventconfFilterJSON)
	result, err := c.GetEventconfFilterEvents(t.Context(), "default", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result["events"].([]any)) < 1 {
		t.Errorf("events = %v, want at least 1", result["events"])
	}
}

func TestGetEventconfSourceNames(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/eventconf/sources/names",
		eventconfSourceNamesJSON)
	result, err := c.GetEventconfSourceNames(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, name := range result {
		if name == "default" {
			found = true
		}
	}
	if !found {
		t.Errorf("result = %v, want to contain default", result)
	}
}

func TestGetEventconfSource(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/eventconf/sources/default", eventconfSourceJSON)
	result, err := c.GetEventconfSource(t.Context(), "default")
	if err != nil {
		t.Fatal(err)
	}
	if result["id"] != "default" {
		t.Errorf("id = %v, want default", result["id"])
	}
}

func TestDownloadEventconfEvents(t *testing.T) {
	mux, c := newTestClient(t)
	handleText(mux, "GET "+v2Path+"/eventconf/sources/default/events/download",
		"<events><event/></events>")
	result, err := c.DownloadEventconfEvents(t.Context(), "default")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, "<events>") {
		t.Errorf("result = %q, want to contain <events>", result)
	}
}

func TestGetEventconfVendorEvents(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/eventconf/vendors/Cisco/events",
		eventconfVendorEventsJSON)
	result, err := c.GetEventconfVendorEvents(t.Context(), "Cisco")
	if err != nil {
		t.Fatal(err)
	}
	event := result["events"].([]any)[0].(map[string]any)
	if event["label"] != "BGP Up" {
		t.Errorf("label = %v, want BGP Up", event["label"])
	}
}

func TestCreateEventconfEvent(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "POST "+v2Path+"/eventconf/sources/default/events",
		eventconfEventJSON)
	result, err := c.CreateEventconfEvent(t.Context(), "default",
		eventconfEventBody(t))
	if err != nil {
		t.Fatal(err)
	}
	if result["uei"] != "uei.opennms.org/custom/testEvent" {
		t.Errorf("uei = %v", result["uei"])
	}
}

func TestUpdateEventconfEvent(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "PUT "+v2Path+"/eventconf/sources/default/events/evt-1", 204)
	err := c.UpdateEventconfEvent(t.Context(), "default", "evt-1",
		eventconfEventBody(t))
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetEventconfSourcesStatus(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "PATCH "+v2Path+"/eventconf/sources/status", 204)
	err := c.SetEventconfSourcesStatus(t.Context(), map[string]any{"default": true})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["default"] != true {
		t.Errorf("body = %v, want default true", body)
	}
}

func TestSetEventconfEventsStatus(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "PATCH "+v2Path+"/eventconf/sources/default/events/status", 204)
	err := c.SetEventconfEventsStatus(t.Context(), "default",
		map[string]any{"evt-1": false})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteEventconfSources(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v2Path+"/eventconf/sources", 204)
	err := c.DeleteEventconfSources(t.Context(),
		map[string]any{"ids": []any{"custom"}})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteEventconfEvents(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v2Path+"/eventconf/sources/default/events", 204)
	err := c.DeleteEventconfEvents(t.Context(), "default",
		map[string]any{"ids": []any{"evt-1"}})
	if err != nil {
		t.Fatal(err)
	}
}

func TestUploadEventconfFromFile(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v2Path+"/eventconf/upload", 204)
	path := filepath.Join(t.TempDir(), "custom.xml")
	xml := "<events><event-file>custom.xml</event-file></events>"
	if err := os.WriteFile(path, []byte(xml), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := c.UploadEventconfFile(t.Context(), path); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(req.body, "custom.xml") {
		t.Errorf("body = %q, want to contain custom.xml", req.body)
	}
}

func TestUploadEventconfFromBytes(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v2Path+"/eventconf/upload", 204)
	xml := []byte("<events><event-file>custom.xml</event-file></events>")
	if err := c.UploadEventconf(t.Context(), xml); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(req.body, "custom.xml") {
		t.Errorf("body = %q, want to contain custom.xml", req.body)
	}
	if !strings.Contains(req.body, "events.xml") {
		t.Errorf("body = %q, want to contain events.xml", req.body)
	}
}
