package opennms

// Tests for the outage timeline methods – /rest/timeline.

import (
	"bytes"
	"net/http"
	"strings"
	"testing"
)

var timelinePNG = []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR" +
	"\x00\x00\x00\x01\x00\x00\x00\x01\x08\x02\x00\x00\x00")

const timelineHTML = `<img src="/opennms/rest/timeline/image/1/10.0.0.1/4` +
	`/1700000000/1700086400/500" usemap="#timeline">`

// timelineHandle registers a handler answering 200 with the given
// content type and raw body (PNG images, HTML).
func timelineHandle(mux *http.ServeMux, pattern, contentType string, body []byte) {
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		w.Write(body)
	})
}

func TestGetTimelineHeader(t *testing.T) {
	mux, c := newTestClient(t)
	timelineHandle(mux, "GET "+v1Path+"/timeline/header/1700000000/1700086400/500",
		"image/png", timelinePNG)
	result, err := c.GetTimelineHeader(t.Context(), 1700000000, 1700086400, 500)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(result, timelinePNG) {
		t.Errorf("result = %q, want PNG fixture", result)
	}
	if !bytes.HasPrefix(result, []byte("\x89PNG")) {
		t.Errorf("result does not start with PNG magic: %q", result)
	}
}

func TestGetTimelineImage(t *testing.T) {
	mux, c := newTestClient(t)
	timelineHandle(mux,
		"GET "+v1Path+"/timeline/image/1/10.0.0.1/4/1700000000/1700086400/500",
		"image/png", timelinePNG)
	result, err := c.GetTimelineImage(t.Context(), 1, "10.0.0.1", 4,
		1700000000, 1700086400, 500)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(result, timelinePNG) {
		t.Errorf("result = %q, want PNG fixture", result)
	}
}

func TestGetTimelineEmpty(t *testing.T) {
	mux, c := newTestClient(t)
	timelineHandle(mux, "GET "+v1Path+"/timeline/empty/1700000000/1700086400/500",
		"image/png", timelinePNG)
	result, err := c.GetTimelineEmpty(t.Context(), 1700000000, 1700086400, 500)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(result, timelinePNG) {
		t.Errorf("result = %q, want PNG fixture", result)
	}
}

func TestGetTimelineHTML(t *testing.T) {
	mux, c := newTestClient(t)
	timelineHandle(mux,
		"GET "+v1Path+"/timeline/html/1/10.0.0.1/4/1700000000/1700086400/500",
		"text/html", []byte(timelineHTML))
	result, err := c.GetTimelineHTML(t.Context(), 1, "10.0.0.1", 4,
		1700000000, 1700086400, 500)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, "<img") {
		t.Errorf("result = %q, want <img present", result)
	}
}
