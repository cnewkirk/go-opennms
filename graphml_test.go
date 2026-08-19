package opennms

// Tests for the graphml methods – /rest/graphml.

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

const graphmlXML = `<?xml version="1.0" encoding="UTF-8"?>` +
	`<graphml xmlns="http://graphml.graphdrawing.org/xmlns">` +
	`<graph id="my-graph" edgedefault="undirected">` +
	`<node id="n0"/></graph></graphml>`

// graphmlHandleXML registers a handler answering 200 with an
// application/xml body and returns a capture of the last request.
func graphmlHandleXML(mux *http.ServeMux, pattern, body string) *capture {
	cap := &capture{}
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		cap.record(r)
		w.Header().Set("Content-Type", "application/xml")
		io.WriteString(w, body)
	})
	return cap
}

func TestGetGraphml(t *testing.T) {
	mux, c := newTestClient(t)
	graphmlHandleXML(mux, "GET "+v1Path+"/graphml/my-graph", graphmlXML)
	result, err := c.GetGraphml(t.Context(), "my-graph")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, "<graphml") {
		t.Errorf("result = %q, want to contain <graphml", result)
	}
}

func TestCreateGraphml(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v1Path+"/graphml/my-graph", 201)
	if err := c.CreateGraphml(t.Context(), "my-graph", graphmlXML); err != nil {
		t.Fatal(err)
	}
	if got := req.header.Get("Content-Type"); got != "application/xml" {
		t.Errorf("Content-Type = %q, want application/xml", got)
	}
	if req.body != graphmlXML {
		t.Errorf("body = %q, want the GraphML document", req.body)
	}
}

func TestDeleteGraphml(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/graphml/my-graph", 200)
	if err := c.DeleteGraphml(t.Context(), "my-graph"); err != nil {
		t.Fatal(err)
	}
}
