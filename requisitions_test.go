package opennms

// Tests for the requisitions methods – /rest/requisitions.

import (
	"encoding/json"
	"strings"
	"testing"
)

const requisitionInterfaceJSON = `{
	"ip-addr": "192.168.1.1",
	"snmp-primary": "P",
	"status": 1,
	"monitored-service": [{"service-name": "ICMP"}, {"service-name": "SNMP"}]
}`

const requisitionNodeJSON = `{
	"foreign-id": "router01",
	"node-label": "router01.example.com",
	"location": "Default",
	"interface": [` + requisitionInterfaceJSON + `],
	"category": [{"name": "Production"}],
	"asset": [{"name": "manufacturer", "value": "Cisco"}],
	"meta-data": []
}`

const requisitionJSON = `{
	"foreign-source": "Routers",
	"date-stamp": "2024-06-01T10:00:00.000+0000",
	"node": [` + requisitionNodeJSON + `],
	"node-count": 1
}`

const requisitionListJSON = `{"requisition": [` + requisitionJSON + `], "count": 1}`

const requisitionNodeListJSON = `{
	"node": [` + requisitionNodeJSON + `],
	"totalCount": 1,
	"count": 1,
	"offset": 0
}`

const requisitionServiceJSON = `{"service-name": "ICMP"}`

const requisitionCategoryJSON = `{"name": "Production"}`

const requisitionAssetJSON = `{"name": "manufacturer", "value": "Cisco"}`

// requisitionUnmarshal decodes a fixture JSON string into a map for
// use as a request payload.
func requisitionUnmarshal(t *testing.T, fixture string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(fixture), &m); err != nil {
		t.Fatal(err)
	}
	return m
}

func TestGetRequisitions(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/requisitions", requisitionListJSON)
	result, err := c.GetRequisitions(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	req := result["requisition"].([]any)[0].(map[string]any)
	if req["foreign-source"] != "Routers" {
		t.Errorf("foreign-source = %v, want Routers", req["foreign-source"])
	}
}

func TestGetRequisition(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/requisitions/Routers", requisitionJSON)
	result, err := c.GetRequisition(t.Context(), "Routers")
	if err != nil {
		t.Fatal(err)
	}
	if result["foreign-source"] != "Routers" {
		t.Errorf("foreign-source = %v, want Routers", result["foreign-source"])
	}
	if result["node-count"].(float64) != 1 {
		t.Errorf("node-count = %v, want 1", result["node-count"])
	}
}

func TestGetRequisitionCount(t *testing.T) {
	mux, c := newTestClient(t)
	handleText(mux, "GET "+v1Path+"/requisitions/count", "3")
	count, err := c.GetRequisitionCount(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Errorf("count = %d, want 3", count)
	}
}

func TestGetDeployedRequisitions(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/requisitions/deployed", requisitionListJSON)
	result, err := c.GetDeployedRequisitions(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	first := result["requisition"].([]any)[0].(map[string]any)
	if first["foreign-source"] != "Routers" {
		t.Errorf("foreign-source = %v, want Routers", first["foreign-source"])
	}
	if !strings.Contains(req.path, "/deployed") {
		t.Errorf("path = %q, want /deployed", req.path)
	}
}

func TestGetDeployedRequisitionCount(t *testing.T) {
	mux, c := newTestClient(t)
	handleText(mux, "GET "+v1Path+"/requisitions/deployed/count", "2")
	count, err := c.GetDeployedRequisitionCount(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}
}

func TestCreateRequisition(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "POST "+v1Path+"/requisitions", requisitionJSON)
	err := c.CreateRequisition(t.Context(),
		map[string]any{"foreign-source": "Routers", "node": []any{}})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["foreign-source"] != "Routers" {
		t.Errorf("foreign-source = %v, want Routers", body["foreign-source"])
	}
	if got := req.header.Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
}

func TestImportRequisition(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "PUT "+v1Path+"/requisitions/Routers/import", 202)
	if err := c.ImportRequisition(t.Context(), "Routers", true); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(req.path, "/Routers/import") {
		t.Errorf("path = %q, want /Routers/import", req.path)
	}
	if req.query.Get("rescanExisting") != "" {
		t.Errorf("rescanExisting = %q, want absent", req.query.Get("rescanExisting"))
	}
}

func TestImportRequisitionNoRescan(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "PUT "+v1Path+"/requisitions/Routers/import", 202)
	if err := c.ImportRequisition(t.Context(), "Routers", false); err != nil {
		t.Fatal(err)
	}
	if got := req.query.Get("rescanExisting"); got != "false" {
		t.Errorf("rescanExisting = %q, want false", got)
	}
}

func TestUpdateRequisition(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "PUT "+v1Path+"/requisitions/Routers", formType, 202, "")
	err := c.UpdateRequisition(t.Context(), "Routers",
		map[string]string{"date-stamp": "2024-06-01"})
	if err != nil {
		t.Fatal(err)
	}
	if req.body != "date-stamp=2024-06-01" {
		t.Errorf("body = %q, want date-stamp=2024-06-01", req.body)
	}
}

func TestDeleteRequisition(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/requisitions/Routers", 202)
	if err := c.DeleteRequisition(t.Context(), "Routers"); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteDeployedRequisition(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "DELETE "+v1Path+"/requisitions/deployed/Routers", 202)
	if err := c.DeleteDeployedRequisition(t.Context(), "Routers"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(req.path, "/deployed/Routers") {
		t.Errorf("path = %q, want /deployed/Routers", req.path)
	}
}

// Nodes

func TestGetRequisitionNodes(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/requisitions/Routers/nodes", requisitionNodeListJSON)
	result, err := c.GetRequisitionNodes(t.Context(), "Routers")
	if err != nil {
		t.Fatal(err)
	}
	node := result["node"].([]any)[0].(map[string]any)
	if node["foreign-id"] != "router01" {
		t.Errorf("foreign-id = %v, want router01", node["foreign-id"])
	}
}

func TestGetRequisitionNode(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/requisitions/Routers/nodes/router01",
		requisitionNodeJSON)
	result, err := c.GetRequisitionNode(t.Context(), "Routers", "router01")
	if err != nil {
		t.Fatal(err)
	}
	if result["node-label"] != "router01.example.com" {
		t.Errorf("node-label = %v", result["node-label"])
	}
}

func TestCreateRequisitionNode(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "POST "+v1Path+"/requisitions/Routers/nodes",
		requisitionNodeJSON)
	err := c.CreateRequisitionNode(t.Context(), "Routers",
		requisitionUnmarshal(t, requisitionNodeJSON))
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["foreign-id"] != "router01" {
		t.Errorf("foreign-id = %v, want router01", body["foreign-id"])
	}
}

func TestUpdateRequisitionNode(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "PUT "+v1Path+"/requisitions/Routers/nodes/router01",
		formType, 202, "")
	err := c.UpdateRequisitionNode(t.Context(), "Routers", "router01",
		map[string]string{"node-label": "updated"})
	if err != nil {
		t.Fatal(err)
	}
	if req.body != "node-label=updated" {
		t.Errorf("body = %q, want node-label=updated", req.body)
	}
}

func TestDeleteRequisitionNode(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/requisitions/Routers/nodes/router01", 202)
	if err := c.DeleteRequisitionNode(t.Context(), "Routers", "router01"); err != nil {
		t.Fatal(err)
	}
}

// Interfaces

func TestGetRequisitionNodeInterfaces(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/requisitions/Routers/nodes/router01/interfaces",
		`{"interface": [`+requisitionInterfaceJSON+`]}`)
	result, err := c.GetRequisitionNodeInterfaces(t.Context(), "Routers", "router01")
	if err != nil {
		t.Fatal(err)
	}
	iface := result["interface"].([]any)[0].(map[string]any)
	if iface["ip-addr"] != "192.168.1.1" {
		t.Errorf("ip-addr = %v, want 192.168.1.1", iface["ip-addr"])
	}
}

func TestCreateRequisitionNodeInterface(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "POST "+v1Path+"/requisitions/Routers/nodes/router01/interfaces",
		requisitionInterfaceJSON)
	err := c.CreateRequisitionNodeInterface(t.Context(), "Routers", "router01",
		requisitionUnmarshal(t, requisitionInterfaceJSON))
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["ip-addr"] != "192.168.1.1" {
		t.Errorf("ip-addr = %v, want 192.168.1.1", body["ip-addr"])
	}
}

func TestUpdateRequisitionNodeInterface(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux,
		"PUT "+v1Path+"/requisitions/Routers/nodes/router01/interfaces/192.168.1.1",
		formType, 202, "")
	err := c.UpdateRequisitionNodeInterface(t.Context(), "Routers", "router01",
		"192.168.1.1", map[string]string{"snmp-primary": "S"})
	if err != nil {
		t.Fatal(err)
	}
	if req.body != "snmp-primary=S" {
		t.Errorf("body = %q, want snmp-primary=S", req.body)
	}
}

func TestDeleteRequisitionNodeInterface(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux,
		"DELETE "+v1Path+"/requisitions/Routers/nodes/router01/interfaces/192.168.1.1",
		202)
	err := c.DeleteRequisitionNodeInterface(t.Context(), "Routers", "router01",
		"192.168.1.1")
	if err != nil {
		t.Fatal(err)
	}
}

// Services

func TestGetRequisitionNodeServices(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux,
		"GET "+v1Path+"/requisitions/Routers/nodes/router01/interfaces/192.168.1.1/services",
		`{"monitored-service": [`+requisitionServiceJSON+`]}`)
	result, err := c.GetRequisitionNodeServices(t.Context(), "Routers", "router01",
		"192.168.1.1")
	if err != nil {
		t.Fatal(err)
	}
	svc := result["monitored-service"].([]any)[0].(map[string]any)
	if svc["service-name"] != "ICMP" {
		t.Errorf("service-name = %v, want ICMP", svc["service-name"])
	}
}

func TestCreateRequisitionNodeService(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux,
		"POST "+v1Path+"/requisitions/Routers/nodes/router01/interfaces/192.168.1.1/services",
		requisitionServiceJSON)
	err := c.CreateRequisitionNodeService(t.Context(), "Routers", "router01",
		"192.168.1.1", map[string]any{"service-name": "HTTP"})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["service-name"] != "HTTP" {
		t.Errorf("service-name = %v, want HTTP", body["service-name"])
	}
}

func TestDeleteRequisitionNodeService(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux,
		"DELETE "+v1Path+"/requisitions/Routers/nodes/router01/interfaces/192.168.1.1/services/ICMP",
		202)
	err := c.DeleteRequisitionNodeService(t.Context(), "Routers", "router01",
		"192.168.1.1", "ICMP")
	if err != nil {
		t.Fatal(err)
	}
}

// Categories

func TestGetRequisitionNodeCategories(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/requisitions/Routers/nodes/router01/categories",
		`{"category": [`+requisitionCategoryJSON+`]}`)
	result, err := c.GetRequisitionNodeCategories(t.Context(), "Routers", "router01")
	if err != nil {
		t.Fatal(err)
	}
	cat := result["category"].([]any)[0].(map[string]any)
	if cat["name"] != "Production" {
		t.Errorf("name = %v, want Production", cat["name"])
	}
}

func TestAddRequisitionNodeCategory(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "POST "+v1Path+"/requisitions/Routers/nodes/router01/categories",
		requisitionCategoryJSON)
	err := c.AddRequisitionNodeCategory(t.Context(), "Routers", "router01",
		map[string]any{"name": "Production"})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["name"] != "Production" {
		t.Errorf("name = %v, want Production", body["name"])
	}
}

func TestDeleteRequisitionNodeCategory(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux,
		"DELETE "+v1Path+"/requisitions/Routers/nodes/router01/categories/Production",
		202)
	err := c.DeleteRequisitionNodeCategory(t.Context(), "Routers", "router01",
		"Production")
	if err != nil {
		t.Fatal(err)
	}
}

// Assets

func TestGetRequisitionNodeAssets(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/requisitions/Routers/nodes/router01/assets",
		`{"asset": [`+requisitionAssetJSON+`]}`)
	result, err := c.GetRequisitionNodeAssets(t.Context(), "Routers", "router01")
	if err != nil {
		t.Fatal(err)
	}
	asset := result["asset"].([]any)[0].(map[string]any)
	if asset["name"] != "manufacturer" {
		t.Errorf("name = %v, want manufacturer", asset["name"])
	}
}

func TestSetRequisitionNodeAsset(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "POST "+v1Path+"/requisitions/Routers/nodes/router01/assets",
		requisitionAssetJSON)
	err := c.SetRequisitionNodeAsset(t.Context(), "Routers", "router01",
		map[string]any{"name": "manufacturer", "value": "Cisco"})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["name"] != "manufacturer" || body["value"] != "Cisco" {
		t.Errorf("body = %v", body)
	}
}

func TestDeleteRequisitionNodeAsset(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux,
		"DELETE "+v1Path+"/requisitions/Routers/nodes/router01/assets/manufacturer",
		202)
	err := c.DeleteRequisitionNodeAsset(t.Context(), "Routers", "router01",
		"manufacturer")
	if err != nil {
		t.Fatal(err)
	}
}
