package opennms

// Tests for core client behavior: auth, retries, error mapping, and
// pagination.

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestBasicAuthSent(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/info", `{"version": "2025.1.0"}`)
	if _, err := c.get(t.Context(), "info", nil, false); err != nil {
		t.Fatal(err)
	}
	user, pass, _ := (&http.Request{Header: req.header}).BasicAuth()
	if user != "admin" || pass != "admin" {
		t.Errorf("basic auth = %q/%q", user, pass)
	}
}

func TestRetriesOn503(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if calls.Add(1) < 3 {
				w.WriteHeader(503)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"ok": true}`)
		}))
	defer srv.Close()
	c := NewClient(srv.URL, "admin", "admin", WithRetries(3))
	result, err := c.getObject(t.Context(), "whatever", nil, false)
	if err != nil {
		t.Fatal(err)
	}
	if result["ok"] != true {
		t.Errorf("result = %v", result)
	}
	if calls.Load() != 3 {
		t.Errorf("calls = %d, want 3", calls.Load())
	}
}

func TestNoRetriesWhenDisabled(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			calls.Add(1)
			w.WriteHeader(503)
		}))
	defer srv.Close()
	c := NewClient(srv.URL, "admin", "admin", WithRetries(0))
	_, err := c.get(t.Context(), "whatever", nil, false)
	if !errors.Is(err, ErrServer) {
		t.Fatalf("err = %v, want ErrServer", err)
	}
	if calls.Load() != 1 {
		t.Errorf("calls = %d, want 1", calls.Load())
	}
}

func TestErrorSentinels(t *testing.T) {
	for status, sentinel := range map[int]error{
		400: ErrBadRequest,
		401: ErrAuthentication,
		403: ErrForbidden,
		404: ErrNotFound,
		409: ErrConflict,
		500: ErrServer,
	} {
		mux, c := newTestClient(t)
		handleStatus(mux, "GET "+v1Path+"/thing", status)
		_, err := c.get(t.Context(), "thing", nil, false)
		if !errors.Is(err, sentinel) {
			t.Errorf("status %d: err = %v, want %v", status, err, sentinel)
		}
	}
}

func TestPaginate(t *testing.T) {
	pages := map[int]string{
		0: `{"totalCount": 3, "alarm": [{"id": 1}, {"id": 2}]}`,
		2: `{"totalCount": 3, "alarm": [{"id": 3}]}`,
	}
	mux, c := newTestClient(t)
	mux.HandleFunc("GET "+v1Path+"/alarms",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, pages[atoiOrZero(r.URL.Query().Get("offset"))])
		})
	var ids []int
	fetch := func(limit, offset int) (map[string]any, error) {
		return c.GetAlarms(t.Context(),
			&ListOptions{Limit: limit, Offset: offset}, nil)
	}
	for item, err := range Paginate(fetch, "alarm", 2) {
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, int(item.(map[string]any)["id"].(float64)))
	}
	want := []int{1, 2, 3}
	if len(ids) != len(want) {
		t.Fatalf("ids = %v, want %v", ids, want)
	}
	for i := range want {
		if ids[i] != want[i] {
			t.Fatalf("ids = %v, want %v", ids, want)
		}
	}
}

func atoiOrZero(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}
