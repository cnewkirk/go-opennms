package opennms

// Tests for the metadata methods – /api/v2/nodes/{id}/metadata and
// sub-resources.

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

const metadataEntryJSON = `{
	"context": "X-OpenNMS-System",
	"key": "managedBy",
	"value": "ansible-tower"
}`

const metadataListJSON = `{"offset": 0, "count": 3, "totalCount": 3, "metaData": [
	` + metadataEntryJSON + `,
	{"context": "X-OpenNMS-System", "key": "environment", "value": "production"},
	{"context": "requisition", "key": "category", "value": "Routers"}
]}`

const (
	metadataNodeID = "1"
	metadataIP     = "192.168.1.1"
	metadataSvc    = "ICMP"
	metadataCtx    = "X-OpenNMS-System"
	metadataKey    = "managedBy"
	metadataVal    = "ansible-tower"
)

// Base paths for convenience
var (
	metadataNodeBase  = v2Path + "/nodes/" + metadataNodeID + "/metadata"
	metadataIfaceBase = v2Path + "/nodes/" + metadataNodeID +
		"/ipinterfaces/" + metadataIP + "/metadata"
	metadataSvcBase = v2Path + "/nodes/" + metadataNodeID +
		"/ipinterfaces/" + metadataIP + "/services/" + metadataSvc +
		"/metadata"
)

// metadataHandlePosts registers a POST handler that records each
// request body and answers 204.
func metadataHandlePosts(mux *http.ServeMux, pattern string) *[]string {
	var bodies []string
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		bodies = append(bodies, string(body))
		w.WriteHeader(204)
	})
	return &bodies
}

// metadataUnmarshal parses a captured JSON body.
func metadataUnmarshal(t *testing.T, raw string) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.Unmarshal([]byte(raw), &body); err != nil {
		t.Fatalf("parsing body %q: %v", raw, err)
	}
	return body
}

// Node metadata

func TestGetNodeMetadata(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+metadataNodeBase, metadataListJSON)
	result, err := c.GetNodeMetadata(t.Context(), metadataNodeID)
	if err != nil {
		t.Fatal(err)
	}
	entry := result["metaData"].([]any)[0].(map[string]any)
	if entry["context"] != metadataCtx {
		t.Errorf("context = %v, want %s", entry["context"], metadataCtx)
	}
	if entry["key"] != metadataKey {
		t.Errorf("key = %v, want %s", entry["key"], metadataKey)
	}
	if req.path != metadataNodeBase {
		t.Errorf("path = %q", req.path)
	}
}

func TestGetNodeMetadataContext(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+metadataNodeBase+"/"+metadataCtx,
		metadataListJSON)
	result, err := c.GetNodeMetadataContext(t.Context(), metadataNodeID,
		metadataCtx)
	if err != nil {
		t.Fatal(err)
	}
	if entries := result["metaData"].([]any); len(entries) != 3 {
		t.Errorf("len = %d, want 3", len(entries))
	}
	if req.path != metadataNodeBase+"/"+metadataCtx {
		t.Errorf("path = %q", req.path)
	}
}

func TestGetNodeMetadataValue(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux,
		"GET "+metadataNodeBase+"/"+metadataCtx+"/"+metadataKey,
		metadataEntryJSON)
	result, err := c.GetNodeMetadataValue(t.Context(), metadataNodeID,
		metadataCtx, metadataKey)
	if err != nil {
		t.Fatal(err)
	}
	if result["value"] != metadataVal {
		t.Errorf("value = %v, want %s", result["value"], metadataVal)
	}
	if req.path != metadataNodeBase+"/"+metadataCtx+"/"+metadataKey {
		t.Errorf("path = %q", req.path)
	}
}

func TestSetNodeMetadata(t *testing.T) {
	mux, c := newTestClient(t)
	bodies := metadataHandlePosts(mux, "POST "+metadataNodeBase)
	err := c.SetNodeMetadata(t.Context(), metadataNodeID, []map[string]any{
		{"context": metadataCtx, "key": metadataKey, "value": metadataVal},
		{"context": metadataCtx, "key": "k2", "value": metadataVal},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(*bodies) != 2 {
		t.Fatalf("calls = %d, want 2", len(*bodies))
	}
	first := metadataUnmarshal(t, (*bodies)[0])
	if first["context"] != metadataCtx {
		t.Errorf("context = %v, want %s", first["context"], metadataCtx)
	}
	if first["value"] != metadataVal {
		t.Errorf("value = %v, want %s", first["value"], metadataVal)
	}
	if second := metadataUnmarshal(t, (*bodies)[1]); second["key"] != "k2" {
		t.Errorf("key = %v, want k2", second["key"])
	}
}

func TestSetNodeMetadataValue(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux,
		"PUT "+metadataNodeBase+"/"+metadataCtx+"/"+metadataKey+"/"+metadataVal,
		204)
	err := c.SetNodeMetadataValue(t.Context(), metadataNodeID, metadataCtx,
		metadataKey, metadataVal)
	if err != nil {
		t.Fatal(err)
	}
	want := metadataNodeBase + "/" + metadataCtx + "/" + metadataKey +
		"/" + metadataVal
	if req.path != want {
		t.Errorf("path = %q, want %q", req.path, want)
	}
}

func TestDeleteNodeMetadataContext(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "DELETE "+metadataNodeBase+"/"+metadataCtx, 204)
	err := c.DeleteNodeMetadataContext(t.Context(), metadataNodeID,
		metadataCtx)
	if err != nil {
		t.Fatal(err)
	}
	if req.method != "DELETE" {
		t.Errorf("method = %q, want DELETE", req.method)
	}
	if req.path != metadataNodeBase+"/"+metadataCtx {
		t.Errorf("path = %q", req.path)
	}
}

func TestDeleteNodeMetadataKey(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux,
		"DELETE "+metadataNodeBase+"/"+metadataCtx+"/"+metadataKey, 204)
	err := c.DeleteNodeMetadataKey(t.Context(), metadataNodeID, metadataCtx,
		metadataKey)
	if err != nil {
		t.Fatal(err)
	}
	if req.path != metadataNodeBase+"/"+metadataCtx+"/"+metadataKey {
		t.Errorf("path = %q", req.path)
	}
}

// Interface metadata

func TestGetInterfaceMetadata(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+metadataIfaceBase, metadataListJSON)
	result, err := c.GetInterfaceMetadata(t.Context(), metadataNodeID,
		metadataIP)
	if err != nil {
		t.Fatal(err)
	}
	if entries := result["metaData"].([]any); len(entries) != 3 {
		t.Errorf("len = %d, want 3", len(entries))
	}
	if req.path != metadataIfaceBase {
		t.Errorf("path = %q", req.path)
	}
}

func TestGetInterfaceMetadataContext(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+metadataIfaceBase+"/"+metadataCtx,
		metadataListJSON)
	result, err := c.GetInterfaceMetadataContext(t.Context(),
		metadataNodeID, metadataIP, metadataCtx)
	if err != nil {
		t.Fatal(err)
	}
	if entries := result["metaData"].([]any); len(entries) != 3 {
		t.Errorf("len = %d, want 3", len(entries))
	}
	if req.path != metadataIfaceBase+"/"+metadataCtx {
		t.Errorf("path = %q", req.path)
	}
}

func TestGetInterfaceMetadataValue(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux,
		"GET "+metadataIfaceBase+"/"+metadataCtx+"/"+metadataKey,
		metadataEntryJSON)
	result, err := c.GetInterfaceMetadataValue(t.Context(), metadataNodeID,
		metadataIP, metadataCtx, metadataKey)
	if err != nil {
		t.Fatal(err)
	}
	if result["key"] != metadataKey {
		t.Errorf("key = %v, want %s", result["key"], metadataKey)
	}
	if req.path != metadataIfaceBase+"/"+metadataCtx+"/"+metadataKey {
		t.Errorf("path = %q", req.path)
	}
}

func TestSetInterfaceMetadata(t *testing.T) {
	mux, c := newTestClient(t)
	bodies := metadataHandlePosts(mux, "POST "+metadataIfaceBase)
	err := c.SetInterfaceMetadata(t.Context(), metadataNodeID, metadataIP,
		[]map[string]any{
			{"context": metadataCtx, "key": metadataKey, "value": metadataVal},
		})
	if err != nil {
		t.Fatal(err)
	}
	if len(*bodies) != 1 {
		t.Fatalf("calls = %d, want 1", len(*bodies))
	}
	if body := metadataUnmarshal(t, (*bodies)[0]); body["key"] != metadataKey {
		t.Errorf("key = %v, want %s", body["key"], metadataKey)
	}
}

func TestSetInterfaceMetadataValue(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux,
		"PUT "+metadataIfaceBase+"/"+metadataCtx+"/"+metadataKey+"/"+metadataVal,
		204)
	err := c.SetInterfaceMetadataValue(t.Context(), metadataNodeID,
		metadataIP, metadataCtx, metadataKey, metadataVal)
	if err != nil {
		t.Fatal(err)
	}
	want := metadataIfaceBase + "/" + metadataCtx + "/" + metadataKey +
		"/" + metadataVal
	if req.path != want {
		t.Errorf("path = %q, want %q", req.path, want)
	}
}

func TestDeleteInterfaceMetadataContext(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "DELETE "+metadataIfaceBase+"/"+metadataCtx, 204)
	err := c.DeleteInterfaceMetadataContext(t.Context(), metadataNodeID,
		metadataIP, metadataCtx)
	if err != nil {
		t.Fatal(err)
	}
	if req.method != "DELETE" {
		t.Errorf("method = %q, want DELETE", req.method)
	}
}

func TestDeleteInterfaceMetadataKey(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux,
		"DELETE "+metadataIfaceBase+"/"+metadataCtx+"/"+metadataKey, 204)
	err := c.DeleteInterfaceMetadataKey(t.Context(), metadataNodeID,
		metadataIP, metadataCtx, metadataKey)
	if err != nil {
		t.Fatal(err)
	}
	if req.path != metadataIfaceBase+"/"+metadataCtx+"/"+metadataKey {
		t.Errorf("path = %q", req.path)
	}
}

// Service metadata

func TestGetServiceMetadata(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+metadataSvcBase, metadataListJSON)
	result, err := c.GetServiceMetadata(t.Context(), metadataNodeID,
		metadataIP, metadataSvc)
	if err != nil {
		t.Fatal(err)
	}
	if entries := result["metaData"].([]any); len(entries) != 3 {
		t.Errorf("len = %d, want 3", len(entries))
	}
	if req.path != metadataSvcBase {
		t.Errorf("path = %q", req.path)
	}
}

func TestGetServiceMetadataContext(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+metadataSvcBase+"/"+metadataCtx,
		metadataListJSON)
	result, err := c.GetServiceMetadataContext(t.Context(), metadataNodeID,
		metadataIP, metadataSvc, metadataCtx)
	if err != nil {
		t.Fatal(err)
	}
	if entries := result["metaData"].([]any); len(entries) != 3 {
		t.Errorf("len = %d, want 3", len(entries))
	}
	if req.path != metadataSvcBase+"/"+metadataCtx {
		t.Errorf("path = %q", req.path)
	}
}

func TestGetServiceMetadataValue(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux,
		"GET "+metadataSvcBase+"/"+metadataCtx+"/"+metadataKey,
		metadataEntryJSON)
	result, err := c.GetServiceMetadataValue(t.Context(), metadataNodeID,
		metadataIP, metadataSvc, metadataCtx, metadataKey)
	if err != nil {
		t.Fatal(err)
	}
	if result["value"] != metadataVal {
		t.Errorf("value = %v, want %s", result["value"], metadataVal)
	}
	if req.path != metadataSvcBase+"/"+metadataCtx+"/"+metadataKey {
		t.Errorf("path = %q", req.path)
	}
}

func TestSetServiceMetadata(t *testing.T) {
	mux, c := newTestClient(t)
	bodies := metadataHandlePosts(mux, "POST "+metadataSvcBase)
	err := c.SetServiceMetadata(t.Context(), metadataNodeID, metadataIP,
		metadataSvc, []map[string]any{
			{"context": metadataCtx, "key": metadataKey, "value": metadataVal},
		})
	if err != nil {
		t.Fatal(err)
	}
	if len(*bodies) != 1 {
		t.Fatalf("calls = %d, want 1", len(*bodies))
	}
	if body := metadataUnmarshal(t, (*bodies)[0]); body["value"] != metadataVal {
		t.Errorf("value = %v, want %s", body["value"], metadataVal)
	}
}

func TestSetServiceMetadataValue(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux,
		"PUT "+metadataSvcBase+"/"+metadataCtx+"/"+metadataKey+"/"+metadataVal,
		204)
	err := c.SetServiceMetadataValue(t.Context(), metadataNodeID, metadataIP,
		metadataSvc, metadataCtx, metadataKey, metadataVal)
	if err != nil {
		t.Fatal(err)
	}
	want := metadataSvcBase + "/" + metadataCtx + "/" + metadataKey +
		"/" + metadataVal
	if req.path != want {
		t.Errorf("path = %q, want %q", req.path, want)
	}
}

func TestDeleteServiceMetadataContext(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "DELETE "+metadataSvcBase+"/"+metadataCtx, 204)
	err := c.DeleteServiceMetadataContext(t.Context(), metadataNodeID,
		metadataIP, metadataSvc, metadataCtx)
	if err != nil {
		t.Fatal(err)
	}
	if req.method != "DELETE" {
		t.Errorf("method = %q, want DELETE", req.method)
	}
}

func TestDeleteServiceMetadataKey(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux,
		"DELETE "+metadataSvcBase+"/"+metadataCtx+"/"+metadataKey, 204)
	err := c.DeleteServiceMetadataKey(t.Context(), metadataNodeID,
		metadataIP, metadataSvc, metadataCtx, metadataKey)
	if err != nil {
		t.Fatal(err)
	}
	if req.path != metadataSvcBase+"/"+metadataCtx+"/"+metadataKey {
		t.Errorf("path = %q", req.path)
	}
}
