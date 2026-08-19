package opennms

// Tests for the SNMP interfaces v2 method – /api/v2/snmpinterfaces
// (read-only).

import "testing"

const snmpInterfaceListJSON = `{
	"snmpInterface": [
		{"id": 301, "ifIndex": 6, "ifName": "GigabitEthernet0/0",
		 "ifDescr": "GigabitEthernet0/0", "ifAlias": "uplink-to-core",
		 "ifType": 6, "ifOperStatus": 1, "ifAdminStatus": 1,
		 "ifSpeed": 1000000000, "physAddr": "04:01:3f:75:f1:01",
		 "collect": "C", "poll": "P", "nodeId": 1,
		 "lastCapsdPoll": "2024-06-01T09:00:00.000+0000"}
	],
	"totalCount": 1, "count": 1, "offset": 0
}`

func TestGetSnmpInterfacesDefaults(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v2Path+"/snmpinterfaces", snmpInterfaceListJSON)
	result, err := c.GetSnmpInterfaces(t.Context(), "", nil)
	if err != nil {
		t.Fatal(err)
	}
	iface := result["snmpInterface"].([]any)[0].(map[string]any)
	if iface["ifIndex"].(float64) != 6 {
		t.Errorf("ifIndex = %v, want 6", iface["ifIndex"])
	}
	if iface["ifName"] != "GigabitEthernet0/0" {
		t.Errorf("ifName = %v", iface["ifName"])
	}
	if req.query.Get("limit") != "10" || req.query.Get("offset") != "0" {
		t.Errorf("query = %v", req.query)
	}
	if req.query.Has("_s") {
		t.Errorf("_s should be absent, query = %v", req.query)
	}
}

func TestGetSnmpInterfacesWithFiql(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v2Path+"/snmpinterfaces", snmpInterfaceListJSON)
	result, err := c.GetSnmpInterfaces(t.Context(), "ifIndex==6", nil)
	if err != nil {
		t.Fatal(err)
	}
	iface := result["snmpInterface"].([]any)[0].(map[string]any)
	if iface["ifIndex"].(float64) != 6 {
		t.Errorf("ifIndex = %v, want 6", iface["ifIndex"])
	}
	if got := req.query.Get("_s"); got != "ifIndex==6" {
		t.Errorf("_s = %q, want ifIndex==6", got)
	}
}

func TestGetSnmpInterfacesPagination(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v2Path+"/snmpinterfaces", snmpInterfaceListJSON)
	_, err := c.GetSnmpInterfaces(t.Context(), "",
		&ListOptions{Limit: 50, Offset: 100})
	if err != nil {
		t.Fatal(err)
	}
	if req.query.Get("limit") != "50" || req.query.Get("offset") != "100" {
		t.Errorf("query = %v", req.query)
	}
}

func TestGetSnmpInterfacesFiqlNodeLabel(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v2Path+"/snmpinterfaces", snmpInterfaceListJSON)
	_, err := c.GetSnmpInterfaces(t.Context(), "node.label==onms-prd-01", nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := req.query.Get("_s"); got != "node.label==onms-prd-01" {
		t.Errorf("_s = %q", got)
	}
}

func TestGetSnmpInterfacesUsesV2URL(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v2Path+"/snmpinterfaces", snmpInterfaceListJSON)
	if _, err := c.GetSnmpInterfaces(t.Context(), "", nil); err != nil {
		t.Fatal(err)
	}
	if req.path != v2Path+"/snmpinterfaces" {
		t.Errorf("path = %q", req.path)
	}
}
