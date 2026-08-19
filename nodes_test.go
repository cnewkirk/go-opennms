package opennms

// Tests for the nodes methods – /rest/nodes and sub-resources.

import (
	"encoding/json"
	"testing"
)

const nodeJSON = `{
	"id": 1,
	"label": "router01.example.com",
	"labelSource": "H",
	"foreignSource": "Routers",
	"foreignId": "router01",
	"location": "Default",
	"type": "A",
	"sysObjectId": ".1.3.6.1.4.1.9.1.1",
	"sysName": "router01",
	"sysDescription": "Cisco IOS Software, Version 15.1",
	"sysContact": "noc@example.com",
	"sysLocation": "DC1 Rack A1",
	"createTime": "2024-01-15T12:00:00.000+0000",
	"lastCapsdPoll": "2024-06-01T09:00:00.000+0000",
	"categories": [{"id": 2, "name": "Production"}]
}`

const nodeListJSON = `{
	"count": 1, "offset": 0, "totalCount": 1,
	"node": [` + nodeJSON + `]
}`

const nodeMonitoredServiceJSON = `{
	"id": 201,
	"serviceName": "ICMP",
	"status": "A",
	"lastGood": "2024-06-01T09:28:00.000+0000",
	"lastFail": "2024-06-01T08:00:00.000+0000",
	"qualifier": null,
	"source": "P",
	"respond": "Y",
	"down": false,
	"serviceType": {"id": 1, "name": "ICMP"}
}`

const nodeMonitoredServiceListJSON = `{
	"count": 1, "offset": 0, "totalCount": 1,
	"service": [` + nodeMonitoredServiceJSON + `]
}`

const nodeIpInterfaceJSON = `{
	"id": 101,
	"ipAddress": "192.168.1.1",
	"hostName": "router01.example.com",
	"isManaged": "M",
	"snmpPrimary": "P",
	"ipLastCapsdPoll": "2024-06-01T09:00:00.000+0000",
	"isDown": false,
	"nodeId": 1,
	"monitoredServices": [` + nodeMonitoredServiceJSON + `]
}`

const nodeIpInterfaceListJSON = `{
	"count": 1, "offset": 0, "totalCount": 1,
	"ipInterface": [` + nodeIpInterfaceJSON + `]
}`

const nodeSnmpInterfaceJSON = `{
	"id": 301,
	"ifIndex": 6,
	"ifName": "GigabitEthernet0/0",
	"ifDescr": "GigabitEthernet0/0",
	"ifAlias": "uplink-to-core",
	"ifType": 6,
	"ifOperStatus": 1,
	"ifAdminStatus": 1,
	"ifSpeed": 1000000000,
	"physAddr": "04:01:3f:75:f1:01",
	"collect": "C",
	"poll": "P",
	"nodeId": 1,
	"lastCapsdPoll": "2024-06-01T09:00:00.000+0000"
}`

const nodeSnmpInterfaceListJSON = `{
	"count": 1, "offset": 0, "totalCount": 1,
	"snmpInterface": [` + nodeSnmpInterfaceJSON + `]
}`

const nodeAssetRecordJSON = `{
	"id": 1,
	"category": "Routers",
	"manufacturer": "Cisco",
	"vendor": "Cisco Systems",
	"modelNumber": "ISR4331",
	"serialNumber": "FDO2147A0BC",
	"description": "Core distribution router",
	"operatingSystem": "IOS-XE 17.3",
	"rack": "A",
	"building": "HQ",
	"floor": "1",
	"room": "DC1",
	"country": "US",
	"lastModifiedBy": "admin",
	"lastModifiedDate": "2024-05-20T14:00:00.000+0000"
}`

const nodeHardwareEntityJSON = `{
	"id": 1,
	"entityPhysicalIndex": 1,
	"entPhysicalDescr": "Cisco ISR4331 chassis",
	"entPhysicalClass": 3,
	"entPhysicalName": "Chassis",
	"entPhysicalSerialNum": "FDO2147A0BC",
	"entPhysicalMfgName": "Cisco Systems",
	"entPhysicalModelName": "ISR4331",
	"entPhysicalIsFRU": true,
	"children": []
}`

const nodeCategoryJSON = `{
	"id": 2,
	"name": "Production",
	"authorizedGroups": ["network-ops"]
}`

// Nodes

func TestGetNodesDefault(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/nodes", nodeListJSON)
	result, err := c.GetNodes(t.Context(), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	node := result["node"].([]any)[0].(map[string]any)
	if node["id"].(float64) != 1 {
		t.Errorf("id = %v, want 1", node["id"])
	}
	if got := req.query.Get("limit"); got != "10" {
		t.Errorf("limit = %q, want 10", got)
	}
	if got := req.query.Get("offset"); got != "0" {
		t.Errorf("offset = %q, want 0", got)
	}
}

func TestGetNodesWithFilters(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/nodes", nodeListJSON)
	_, err := c.GetNodes(t.Context(),
		&ListOptions{Limit: 100, OrderBy: "label"},
		map[string]string{"foreignSource": "Routers"})
	if err != nil {
		t.Fatal(err)
	}
	for param, want := range map[string]string{
		"limit": "100", "foreignSource": "Routers", "orderBy": "label",
	} {
		if got := req.query.Get(param); got != want {
			t.Errorf("%s = %q, want %q", param, got, want)
		}
	}
}

func TestGetNode(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/nodes/1", nodeJSON)
	result, err := c.GetNode(t.Context(), "1")
	if err != nil {
		t.Fatal(err)
	}
	if result["id"].(float64) != 1 {
		t.Errorf("id = %v, want 1", result["id"])
	}
	if result["label"] != "router01.example.com" {
		t.Errorf("label = %v", result["label"])
	}
}

func TestGetNodeByFsFid(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/nodes/Routers:router01", nodeJSON)
	result, err := c.GetNode(t.Context(), "Routers:router01")
	if err != nil {
		t.Fatal(err)
	}
	if result["foreignSource"] != "Routers" {
		t.Errorf("foreignSource = %v, want Routers", result["foreignSource"])
	}
}

func TestGetNodeCount(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v2Path+"/nodes",
		`{"totalCount": 57, "count": 1, "offset": 0, "node": []}`)
	count, err := c.GetNodeCount(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if count != 57 {
		t.Errorf("count = %d, want 57", count)
	}
	if got := req.query.Get("limit"); got != "1" {
		t.Errorf("limit = %q, want 1", got)
	}
}

func TestCreateNode(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "POST "+v1Path+"/nodes", xmlType, 201, nodeJSON)
	result, err := c.CreateNode(t.Context(), map[string]any{
		"label": "newnode.example.com", "type": "A",
		"foreignSource": "Test", "foreignId": "newnode01",
		"location": "Default"})
	if err != nil {
		t.Fatal(err)
	}
	if result["id"].(float64) != 1 {
		t.Errorf("id = %v, want 1", result["id"])
	}
	want := `<node label="newnode.example.com" type="A"` +
		` foreignSource="Test" foreignId="newnode01">` +
		`<location>Default</location></node>`
	if req.body != want {
		t.Errorf("body = %q, want %q", req.body, want)
	}
}

func TestUpdateNode(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "PUT "+v1Path+"/nodes/1", formType, 204, "")
	err := c.UpdateNode(t.Context(), "1", map[string]string{"label": "updated-label"})
	if err != nil {
		t.Fatal(err)
	}
	if req.body != "label=updated-label" {
		t.Errorf("body = %q, want label=updated-label", req.body)
	}
}

func TestDeleteNode(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/nodes/1", 202)
	if err := c.DeleteNode(t.Context(), "1"); err != nil {
		t.Fatal(err)
	}
}

func TestRescanNode(t *testing.T) {
	mux, c := newTestClient(t)
	handleContract(mux, "PUT "+v1Path+"/nodes/1/rescan", formType, 204, "")
	if err := c.RescanNode(t.Context(), "1"); err != nil {
		t.Fatal(err)
	}
}

// IP Interfaces

func TestGetNodeIpInterfaces(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/nodes/1/ipinterfaces",
		nodeIpInterfaceListJSON)
	result, err := c.GetNodeIpInterfaces(t.Context(), "1", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	iface := result["ipInterface"].([]any)[0].(map[string]any)
	if iface["ipAddress"] != "192.168.1.1" {
		t.Errorf("ipAddress = %v, want 192.168.1.1", iface["ipAddress"])
	}
	if req.query.Get("limit") != "10" || req.query.Get("offset") != "0" {
		t.Errorf("query = %v", req.query)
	}
}

func TestGetNodeIpInterface(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/nodes/1/ipinterfaces/192.168.1.1",
		nodeIpInterfaceJSON)
	result, err := c.GetNodeIpInterface(t.Context(), "1", "192.168.1.1")
	if err != nil {
		t.Fatal(err)
	}
	if result["ipAddress"] != "192.168.1.1" {
		t.Errorf("ipAddress = %v", result["ipAddress"])
	}
	if result["snmpPrimary"] != "P" {
		t.Errorf("snmpPrimary = %v, want P", result["snmpPrimary"])
	}
}

func TestCreateNodeIpInterface(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "POST "+v1Path+"/nodes/1/ipinterfaces",
		xmlType, 201, nodeIpInterfaceJSON)
	result, err := c.CreateNodeIpInterface(t.Context(), "1", map[string]any{
		"ipAddress": "10.0.0.1", "snmpPrimary": "P", "isManaged": "M"})
	if err != nil {
		t.Fatal(err)
	}
	if result["ipAddress"] != "192.168.1.1" {
		t.Errorf("ipAddress = %v", result["ipAddress"])
	}
	want := `<ipInterface isManaged="M" snmpPrimary="P">` +
		`<ipAddress>10.0.0.1</ipAddress></ipInterface>`
	if req.body != want {
		t.Errorf("body = %q, want %q", req.body, want)
	}
}

func TestUpdateNodeIpInterface(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "PUT "+v1Path+"/nodes/1/ipinterfaces/192.168.1.1",
		formType, 204, "")
	err := c.UpdateNodeIpInterface(t.Context(), "1", "192.168.1.1",
		map[string]string{"isManaged": "U"})
	if err != nil {
		t.Fatal(err)
	}
	if req.body != "isManaged=U" {
		t.Errorf("body = %q, want isManaged=U", req.body)
	}
}

func TestDeleteNodeIpInterface(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/nodes/1/ipinterfaces/192.168.1.1", 202)
	if err := c.DeleteNodeIpInterface(t.Context(), "1", "192.168.1.1"); err != nil {
		t.Fatal(err)
	}
}

// Monitored Services

func TestGetNodeIpServices(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/nodes/1/ipinterfaces/192.168.1.1/services",
		nodeMonitoredServiceListJSON)
	result, err := c.GetNodeIpServices(t.Context(), "1", "192.168.1.1")
	if err != nil {
		t.Fatal(err)
	}
	service := result["service"].([]any)[0].(map[string]any)
	if service["serviceName"] != "ICMP" {
		t.Errorf("serviceName = %v, want ICMP", service["serviceName"])
	}
}

func TestGetNodeIpService(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/nodes/1/ipinterfaces/192.168.1.1/services/ICMP",
		nodeMonitoredServiceJSON)
	result, err := c.GetNodeIpService(t.Context(), "1", "192.168.1.1", "ICMP")
	if err != nil {
		t.Fatal(err)
	}
	if result["serviceName"] != "ICMP" {
		t.Errorf("serviceName = %v, want ICMP", result["serviceName"])
	}
	if result["status"] != "A" {
		t.Errorf("status = %v, want A", result["status"])
	}
}

func TestCreateNodeIpService(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux,
		"POST "+v1Path+"/nodes/1/ipinterfaces/192.168.1.1/services",
		"application/json", 201, nodeMonitoredServiceJSON)
	_, err := c.CreateNodeIpService(t.Context(), "1", "192.168.1.1",
		map[string]any{"serviceType": map[string]any{"name": "HTTP"}})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["serviceType"].(map[string]any)["name"] != "HTTP" {
		t.Errorf("serviceType = %v", body["serviceType"])
	}
}

func TestDeleteNodeIpService(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux,
		"DELETE "+v1Path+"/nodes/1/ipinterfaces/192.168.1.1/services/ICMP", 202)
	err := c.DeleteNodeIpService(t.Context(), "1", "192.168.1.1", "ICMP")
	if err != nil {
		t.Fatal(err)
	}
}

// SNMP Interfaces

func TestGetNodeSnmpInterfaces(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/nodes/1/snmpinterfaces",
		nodeSnmpInterfaceListJSON)
	result, err := c.GetNodeSnmpInterfaces(t.Context(), "1", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	iface := result["snmpInterface"].([]any)[0].(map[string]any)
	if iface["ifIndex"].(float64) != 6 {
		t.Errorf("ifIndex = %v, want 6", iface["ifIndex"])
	}
}

func TestGetNodeSnmpInterface(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/nodes/1/snmpinterfaces/6",
		nodeSnmpInterfaceJSON)
	result, err := c.GetNodeSnmpInterface(t.Context(), "1", 6)
	if err != nil {
		t.Fatal(err)
	}
	if result["ifIndex"].(float64) != 6 {
		t.Errorf("ifIndex = %v, want 6", result["ifIndex"])
	}
	if result["ifOperStatus"].(float64) != 1 {
		t.Errorf("ifOperStatus = %v, want 1", result["ifOperStatus"])
	}
}

func TestCreateNodeSnmpInterface(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "POST "+v1Path+"/nodes/1/snmpinterfaces",
		xmlType, 201, nodeSnmpInterfaceJSON)
	_, err := c.CreateNodeSnmpInterface(t.Context(), "1", map[string]any{
		"ifIndex": 1, "ifName": "lo", "ifType": 24, "collect": "C"})
	if err != nil {
		t.Fatal(err)
	}
	want := `<snmpInterface ifIndex="1" collectFlag="C">` +
		`<ifName>lo</ifName><ifType>24</ifType></snmpInterface>`
	if req.body != want {
		t.Errorf("body = %q, want %q", req.body, want)
	}
}

func TestUpdateNodeSnmpInterface(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "PUT "+v1Path+"/nodes/1/snmpinterfaces/6",
		formType, 204, "")
	err := c.UpdateNodeSnmpInterface(t.Context(), "1", 6,
		map[string]string{"ifAlias": "new-alias"})
	if err != nil {
		t.Fatal(err)
	}
	if req.body != "ifAlias=new-alias" {
		t.Errorf("body = %q, want ifAlias=new-alias", req.body)
	}
}

func TestDeleteNodeSnmpInterface(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/nodes/1/snmpinterfaces/6", 204)
	if err := c.DeleteNodeSnmpInterface(t.Context(), "1", 6); err != nil {
		t.Fatal(err)
	}
}

// Categories

func TestGetNodeCategories(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/nodes/1/categories",
		`{"category": [`+nodeCategoryJSON+`]}`)
	result, err := c.GetNodeCategories(t.Context(), "1")
	if err != nil {
		t.Fatal(err)
	}
	category := result["category"].([]any)[0].(map[string]any)
	if category["name"] != "Production" {
		t.Errorf("name = %v, want Production", category["name"])
	}
}

func TestGetNodeCategory(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/nodes/1/categories/Production",
		nodeCategoryJSON)
	result, err := c.GetNodeCategory(t.Context(), "1", "Production")
	if err != nil {
		t.Fatal(err)
	}
	if result["name"] != "Production" {
		t.Errorf("name = %v, want Production", result["name"])
	}
}

func TestAddNodeCategory(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "POST "+v1Path+"/nodes/1/categories",
		xmlType, 201, nodeCategoryJSON)
	_, err := c.AddNodeCategory(t.Context(), "1",
		map[string]any{"name": "Production"})
	if err != nil {
		t.Fatal(err)
	}
	if req.body != `<category name="Production"></category>` {
		t.Errorf("body = %q", req.body)
	}
}

func TestUpdateNodeCategory(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "PUT "+v1Path+"/nodes/1/categories/Production",
		formType, 204, "")
	err := c.UpdateNodeCategory(t.Context(), "1", "Production",
		map[string]string{"name": "Production"})
	if err != nil {
		t.Fatal(err)
	}
	if req.body != "name=Production" {
		t.Errorf("body = %q, want name=Production", req.body)
	}
}

func TestDeleteNodeCategory(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/nodes/1/categories/Production", 204)
	if err := c.DeleteNodeCategory(t.Context(), "1", "Production"); err != nil {
		t.Fatal(err)
	}
}

// Asset Record

func TestGetNodeAssetRecord(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/nodes/1/assetRecord", nodeAssetRecordJSON)
	result, err := c.GetNodeAssetRecord(t.Context(), "1")
	if err != nil {
		t.Fatal(err)
	}
	if result["manufacturer"] != "Cisco" {
		t.Errorf("manufacturer = %v, want Cisco", result["manufacturer"])
	}
	if result["serialNumber"] != "FDO2147A0BC" {
		t.Errorf("serialNumber = %v", result["serialNumber"])
	}
}

func TestUpdateNodeAssetRecord(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "PUT "+v1Path+"/nodes/1/assetRecord",
		formType, 204, "")
	err := c.UpdateNodeAssetRecord(t.Context(), "1", map[string]string{
		"manufacturer": "Cisco", "modelNumber": "ISR4331"})
	if err != nil {
		t.Fatal(err)
	}
	form := req.formBody(t)
	if form.Get("manufacturer") != "Cisco" || form.Get("modelNumber") != "ISR4331" {
		t.Errorf("form = %v", form)
	}
}

// Hardware Inventory

func TestGetNodeHardwareInventory(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/nodes/1/hardwareInventory",
		nodeHardwareEntityJSON)
	result, err := c.GetNodeHardwareInventory(t.Context(), "1")
	if err != nil {
		t.Fatal(err)
	}
	if result["entPhysicalMfgName"] != "Cisco Systems" {
		t.Errorf("entPhysicalMfgName = %v", result["entPhysicalMfgName"])
	}
	if result["entityPhysicalIndex"].(float64) != 1 {
		t.Errorf("entityPhysicalIndex = %v, want 1", result["entityPhysicalIndex"])
	}
}

func TestGetNodeHardwareEntity(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/nodes/1/hardwareInventory/1",
		nodeHardwareEntityJSON)
	result, err := c.GetNodeHardwareEntity(t.Context(), "1", 1)
	if err != nil {
		t.Fatal(err)
	}
	if result["entPhysicalModelName"] != "ISR4331" {
		t.Errorf("entPhysicalModelName = %v, want ISR4331", result["entPhysicalModelName"])
	}
}

func TestAddNodeHardwareInventory(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "POST "+v1Path+"/nodes/1/hardwareInventory",
		"application/json", 201, nodeHardwareEntityJSON)
	var entity map[string]any
	if err := json.Unmarshal([]byte(nodeHardwareEntityJSON), &entity); err != nil {
		t.Fatal(err)
	}
	if _, err := c.AddNodeHardwareInventory(t.Context(), "1", entity); err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["entityPhysicalIndex"].(float64) != 1 {
		t.Errorf("entityPhysicalIndex = %v, want 1", body["entityPhysicalIndex"])
	}
}

func TestUpdateNodeHardwareEntity(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "PUT "+v1Path+"/nodes/1/hardwareInventory/1",
		formType, 204, "")
	err := c.UpdateNodeHardwareEntity(t.Context(), "1", 1,
		map[string]string{"entPhysicalAlias": "chassis"})
	if err != nil {
		t.Fatal(err)
	}
	if req.body != "entPhysicalAlias=chassis" {
		t.Errorf("body = %q, want entPhysicalAlias=chassis", req.body)
	}
}

func TestDeleteNodeHardwareEntity(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/nodes/1/hardwareInventory/1", 204)
	if err := c.DeleteNodeHardwareEntity(t.Context(), "1", 1); err != nil {
		t.Fatal(err)
	}
}
