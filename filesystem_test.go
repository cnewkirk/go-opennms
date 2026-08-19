package opennms

// Tests for the filesystem methods – /rest/filesystem.

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

const filesystemFileListJSON = `[
	"discovery-configuration.xml",
	"eventconf.xml",
	"snmp-config.xml"
]`

const filesystemExtensionsJSON = `["cfg", "drl", "properties", "xml"]`

const filesystemHelpMD = "# discovery-configuration.xml\n\nControls ..."

const filesystemContentsXML = `<?xml version="1.0"?><discovery-configuration packets-per-second="1"/>`

func TestGetFilesystemFiles(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/filesystem", filesystemFileListJSON)
	result, err := c.GetFilesystemFiles(t.Context(), false)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, f := range result {
		if f == "discovery-configuration.xml" {
			found = true
		}
	}
	if !found {
		t.Errorf("result = %v, want discovery-configuration.xml in it", result)
	}
	if len(req.query) != 0 {
		t.Errorf("query = %v, want none", req.query)
	}
}

func TestGetFilesystemFilesChangedOnly(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/filesystem", filesystemFileListJSON)
	if _, err := c.GetFilesystemFiles(t.Context(), true); err != nil {
		t.Fatal(err)
	}
	if got := req.query.Get("changedFilesOnly"); got != "true" {
		t.Errorf("changedFilesOnly = %q, want true", got)
	}
}

func TestGetFilesystemExtensions(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/filesystem/extensions", filesystemExtensionsJSON)
	result, err := c.GetFilesystemExtensions(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, e := range result {
		if e == "xml" {
			found = true
		}
	}
	if !found {
		t.Errorf("result = %v, want xml in it", result)
	}
}

func TestGetFilesystemHelp(t *testing.T) {
	mux, c := newTestClient(t)
	req := &capture{}
	mux.HandleFunc("GET "+v1Path+"/filesystem/help",
		func(w http.ResponseWriter, r *http.Request) {
			req.record(r)
			w.Header().Set("Content-Type", "text/markdown")
			io.WriteString(w, filesystemHelpMD)
		})
	result, err := c.GetFilesystemHelp(t.Context(), "discovery-configuration.xml")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(result, "# discovery-configuration.xml") {
		t.Errorf("result = %q", result)
	}
	if got := req.query.Get("f"); got != "discovery-configuration.xml" {
		t.Errorf("f = %q", got)
	}
}

func TestGetFilesystemContents(t *testing.T) {
	mux, c := newTestClient(t)
	mux.HandleFunc("GET "+v1Path+"/filesystem/contents",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/xml")
			io.WriteString(w, filesystemContentsXML)
		})
	result, err := c.GetFilesystemContents(t.Context(), "discovery-configuration.xml")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, "<discovery-configuration") {
		t.Errorf("result = %q", result)
	}
}

func TestUploadFilesystemContents(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v1Path+"/filesystem/contents", 200)
	err := c.UploadFilesystemContents(t.Context(),
		"discovery-configuration.xml", []byte(filesystemContentsXML))
	if err != nil {
		t.Fatal(err)
	}
	if got := req.query.Get("f"); got != "discovery-configuration.xml" {
		t.Errorf("f = %q", got)
	}
	if got := req.header.Get("Content-Type"); !strings.HasPrefix(got, "multipart/form-data") {
		t.Errorf("Content-Type = %q", got)
	}
	if !strings.Contains(req.body, "<discovery-configuration") {
		t.Errorf("body = %q", req.body)
	}
}

func TestDeleteFilesystemFile(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "DELETE "+v1Path+"/filesystem/contents", 200)
	if err := c.DeleteFilesystemFile(t.Context(), "obsolete.xml"); err != nil {
		t.Fatal(err)
	}
	if got := req.query.Get("f"); got != "obsolete.xml" {
		t.Errorf("f = %q", got)
	}
}
