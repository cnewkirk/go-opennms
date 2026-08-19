package opennms

// Shared test helpers mirroring the Python suite's conftest.py: a
// client wired to an httptest server, request capture, and the
// Content-Type contract helper.

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const (
	v1Path = "/opennms/rest"
	v2Path = "/opennms/api/v2"

	formType = "application/x-www-form-urlencoded"
	xmlType  = "application/xml"
)

// newTestClient returns a mux to register handlers on and a Client
// pointed at it. Retries are disabled so error-path tests stay fast.
func newTestClient(t *testing.T) (*http.ServeMux, *Client) {
	t.Helper()
	mux := http.NewServeMux()
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return mux, NewClient(srv.URL, "admin", "admin", WithRetries(0))
}

// capture records the last request a handler received.
type capture struct {
	method string
	path   string
	query  url.Values
	body   string
	header http.Header
}

func (cap *capture) record(r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	cap.method = r.Method
	cap.path = r.URL.Path
	cap.query = r.URL.Query()
	cap.body = string(body)
	cap.header = r.Header.Clone()
}

// formBody parses the captured body as form-encoded data.
func (cap *capture) formBody(t *testing.T) url.Values {
	t.Helper()
	form, err := url.ParseQuery(cap.body)
	if err != nil {
		t.Fatalf("parsing form body %q: %v", cap.body, err)
	}
	return form
}

// handleJSON registers a handler answering 200 with the given JSON
// document and returns a capture of the last request.
func handleJSON(mux *http.ServeMux, pattern, body string) *capture {
	cap := &capture{}
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		cap.record(r)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, body)
	})
	return cap
}

// handleText registers a handler answering 200 text/plain (the shape
// of the /count endpoints).
func handleText(mux *http.ServeMux, pattern, body string) *capture {
	cap := &capture{}
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		cap.record(r)
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, body)
	})
	return cap
}

// handleStatus registers a handler answering the given status with an
// empty body (204 for writes, error codes for failure tests).
func handleStatus(mux *http.ServeMux, pattern string, status int) *capture {
	cap := &capture{}
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		cap.record(r)
		w.WriteHeader(status)
	})
	return cap
}

// handleContract registers a handler that enforces the endpoint's
// documented Content-Type, answering 415 Unsupported Media Type for
// any other type exactly as OpenNMS does (see the Python suite's
// add_contract). jsonBody "" answers the bare status.
func handleContract(mux *http.ServeMux, pattern, consumes string, status int, jsonBody string) *capture {
	cap := &capture{}
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		cap.record(r)
		ct := strings.TrimSpace(strings.SplitN(r.Header.Get("Content-Type"), ";", 2)[0])
		if ct != consumes {
			http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
			return
		}
		if jsonBody != "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			io.WriteString(w, jsonBody)
			return
		}
		w.WriteHeader(status)
	})
	return cap
}
