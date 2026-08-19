package opennms

// Tests for the business services methods – /api/v2/business-services.

import (
	"encoding/json"
	"strings"
	"testing"
)

const businessServicesEdgeJSON = `{
	"id": 2001,
	"type": "IP_SERVICE",
	"reductionKeys": ["uei.opennms.org/nodes/nodeDown::1"],
	"mapFunction": {"type": "Identity", "properties": {}},
	"weight": 1,
	"operationalStatus": "MAJOR"
}`

const businessServiceJSON = `{
	"id": 1001,
	"name": "Core Network Availability",
	"attributes": {"dc": "us-east-1", "tier": "critical"},
	"reduceFunction": {"type": "HighestSeverity", "properties": {}},
	"operationalStatus": "MAJOR",
	"edges": [` + businessServicesEdgeJSON + `],
	"parentServices": [],
	"childServices": []
}`

const businessServiceListJSON = `{
	"business-services": [` + businessServiceJSON + `],
	"totalCount": 1, "count": 1, "offset": 0
}`

const businessServicesMapFunctionJSON = `{
	"name": "Identity", "type": "Identity", "properties": {}
}`

const businessServicesMapFunctionsJSON = `[` +
	businessServicesMapFunctionJSON +
	`, {"name": "Increase", "type": "Increase", "properties": {}}]`

const businessServicesReduceFunctionJSON = `{
	"name": "HighestSeverity", "type": "HighestSeverity", "properties": {}
}`

const businessServicesReduceFunctionsJSON = `[` +
	businessServicesReduceFunctionJSON +
	`, {"name": "Threshold", "type": "Threshold", "properties": {"threshold": "0.5"}}]`

func TestGetBusinessServices(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/business-services", businessServiceListJSON)
	result, err := c.GetBusinessServices(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	svc := result["business-services"].([]any)[0].(map[string]any)
	if svc["id"].(float64) != 1001 {
		t.Errorf("id = %v, want 1001", svc["id"])
	}
	if svc["name"] != "Core Network Availability" {
		t.Errorf("name = %v", svc["name"])
	}
}

func TestGetBusinessService(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v2Path+"/business-services/1001", businessServiceJSON)
	result, err := c.GetBusinessService(t.Context(), 1001)
	if err != nil {
		t.Fatal(err)
	}
	if result["id"].(float64) != 1001 {
		t.Errorf("id = %v, want 1001", result["id"])
	}
	if result["operationalStatus"] != "MAJOR" {
		t.Errorf("operationalStatus = %v", result["operationalStatus"])
	}
	if len(result["edges"].([]any)) != 1 {
		t.Errorf("edges = %v", result["edges"])
	}
	if !strings.Contains(req.path, "/business-services/1001") {
		t.Errorf("path = %q", req.path)
	}
}

func TestCreateBusinessService(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "POST "+v2Path+"/business-services",
		"application/json", 201, businessServiceJSON)
	payload := map[string]any{
		"name":           "Core Network Availability",
		"attributes":     map[string]any{"dc": "us-east-1", "tier": "critical"},
		"reduceFunction": map[string]any{"type": "HighestSeverity"},
	}
	result, err := c.CreateBusinessService(t.Context(), payload)
	if err != nil {
		t.Fatal(err)
	}
	if result["id"].(float64) != 1001 {
		t.Errorf("id = %v, want 1001", result["id"])
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["name"] != "Core Network Availability" {
		t.Errorf("name = %v", body["name"])
	}
	if body["reduceFunction"].(map[string]any)["type"] != "HighestSeverity" {
		t.Errorf("reduceFunction = %v", body["reduceFunction"])
	}
	if got := req.header.Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q", got)
	}
}

func TestUpdateBusinessService(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "PUT "+v2Path+"/business-services/1001", 204)
	var updated map[string]any
	if err := json.Unmarshal([]byte(businessServiceJSON), &updated); err != nil {
		t.Fatal(err)
	}
	updated["name"] = "Renamed Service"
	if err := c.UpdateBusinessService(t.Context(), 1001, updated); err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["name"] != "Renamed Service" {
		t.Errorf("name = %v", body["name"])
	}
	if !strings.Contains(req.path, "/business-services/1001") {
		t.Errorf("path = %q", req.path)
	}
}

func TestDeleteBusinessService(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "DELETE "+v2Path+"/business-services/1001", 204)
	if err := c.DeleteBusinessService(t.Context(), 1001); err != nil {
		t.Fatal(err)
	}
	if req.method != "DELETE" {
		t.Errorf("method = %q", req.method)
	}
	if !strings.Contains(req.path, "/business-services/1001") {
		t.Errorf("path = %q", req.path)
	}
}

func TestGetBusinessServiceEdge(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/business-services/edges/2001",
		businessServicesEdgeJSON)
	result, err := c.GetBusinessServiceEdge(t.Context(), 2001)
	if err != nil {
		t.Fatal(err)
	}
	if result["id"].(float64) != 2001 {
		t.Errorf("id = %v, want 2001", result["id"])
	}
	if result["type"] != "IP_SERVICE" {
		t.Errorf("type = %v", result["type"])
	}
}

func TestAddIpServiceEdge(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "POST "+v2Path+"/business-services/1001/ip-service-edge",
		businessServicesEdgeJSON)
	edge := map[string]any{
		"ipServiceId": 201,
		"mapFunction": map[string]any{"type": "Identity"},
		"weight":      1,
	}
	result, err := c.AddIpServiceEdge(t.Context(), 1001, edge)
	if err != nil {
		t.Fatal(err)
	}
	if result["id"].(float64) != 2001 {
		t.Errorf("id = %v, want 2001", result["id"])
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["ipServiceId"].(float64) != 201 {
		t.Errorf("ipServiceId = %v", body["ipServiceId"])
	}
}

func TestAddReductionKeyEdge(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "POST "+v2Path+"/business-services/1001/reduction-key-edge",
		businessServicesEdgeJSON)
	edge := map[string]any{
		"reductionKey": "uei/test::1",
		"mapFunction":  map[string]any{"type": "Identity"},
		"weight":       1,
	}
	if _, err := c.AddReductionKeyEdge(t.Context(), 1001, edge); err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["reductionKey"] != "uei/test::1" {
		t.Errorf("reductionKey = %v", body["reductionKey"])
	}
}

func TestAddChildEdge(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "POST "+v2Path+"/business-services/1001/child-edge",
		businessServicesEdgeJSON)
	edge := map[string]any{
		"childId":     1002,
		"mapFunction": map[string]any{"type": "Identity"},
		"weight":      1,
	}
	if _, err := c.AddChildEdge(t.Context(), 1001, edge); err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["childId"].(float64) != 1002 {
		t.Errorf("childId = %v", body["childId"])
	}
}

func TestRemoveBusinessServiceEdge(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v2Path+"/business-services/1001/edges/2001", 204)
	if err := c.RemoveBusinessServiceEdge(t.Context(), 1001, 2001); err != nil {
		t.Fatal(err)
	}
}

func TestReloadBusinessServiceDaemon(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "POST "+v2Path+"/business-services/daemon/reload", 204)
	if err := c.ReloadBusinessServiceDaemon(t.Context()); err != nil {
		t.Fatal(err)
	}
}

func TestGetMapFunctions(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/business-services/functions/map",
		businessServicesMapFunctionsJSON)
	result, err := c.GetMapFunctions(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result[0].(map[string]any)["name"] != "Identity" {
		t.Errorf("name = %v", result[0])
	}
}

func TestGetMapFunction(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/business-services/functions/map/Identity",
		businessServicesMapFunctionJSON)
	result, err := c.GetMapFunction(t.Context(), "Identity")
	if err != nil {
		t.Fatal(err)
	}
	if result["type"] != "Identity" {
		t.Errorf("type = %v", result["type"])
	}
}

func TestGetReduceFunctions(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/business-services/functions/reduce",
		businessServicesReduceFunctionsJSON)
	result, err := c.GetReduceFunctions(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result[0].(map[string]any)["name"] != "HighestSeverity" {
		t.Errorf("name = %v", result[0])
	}
}

func TestGetReduceFunction(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/business-services/functions/reduce/HighestSeverity",
		businessServicesReduceFunctionJSON)
	result, err := c.GetReduceFunction(t.Context(), "HighestSeverity")
	if err != nil {
		t.Fatal(err)
	}
	if result["type"] != "HighestSeverity" {
		t.Errorf("type = %v", result["type"])
	}
}
