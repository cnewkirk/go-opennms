package opennms

// Tests for the v2 IP interfaces methods – /api/v2/ipinterfaces
// (read-only).

import "testing"

const ipInterfacesV2ListJSON = `{
	"ipInterface": [{
		"id": 101,
		"ipAddress": "192.168.1.1",
		"hostName": "router01.example.com",
		"isManaged": "M",
		"snmpPrimary": "P",
		"ipLastCapsdPoll": "2024-06-01T09:00:00.000+0000",
		"isDown": false,
		"nodeId": 1,
		"monitoredServices": [{
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
		}]
	}],
	"totalCount": 1, "count": 1, "offset": 0
}`

func TestGetIpInterfacesDefaults(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v2Path+"/ipinterfaces", ipInterfacesV2ListJSON)
	result, err := c.GetIpInterfaces(t.Context(), "", nil)
	if err != nil {
		t.Fatal(err)
	}
	iface := result["ipInterface"].([]any)[0].(map[string]any)
	if iface["ipAddress"] != "192.168.1.1" {
		t.Errorf("ipAddress = %v", iface["ipAddress"])
	}
	if req.query.Get("limit") != "10" || req.query.Get("offset") != "0" {
		t.Errorf("query = %v", req.query)
	}
	if req.query.Has("_s") {
		t.Errorf("_s = %q, want absent", req.query.Get("_s"))
	}
}

func TestGetIpInterfacesWithFiql(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v2Path+"/ipinterfaces", ipInterfacesV2ListJSON)
	result, err := c.GetIpInterfaces(t.Context(), "ipAddress==192.168.1.1", nil)
	if err != nil {
		t.Fatal(err)
	}
	iface := result["ipInterface"].([]any)[0].(map[string]any)
	if iface["nodeId"].(float64) != 1 {
		t.Errorf("nodeId = %v, want 1", iface["nodeId"])
	}
	if got := req.query.Get("_s"); got != "ipAddress==192.168.1.1" {
		t.Errorf("_s = %q", got)
	}
}

func TestGetIpInterfacesPagination(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v2Path+"/ipinterfaces", ipInterfacesV2ListJSON)
	_, err := c.GetIpInterfaces(t.Context(), "", &ListOptions{Limit: 25, Offset: 50})
	if err != nil {
		t.Fatal(err)
	}
	if req.query.Get("limit") != "25" || req.query.Get("offset") != "50" {
		t.Errorf("query = %v", req.query)
	}
}

func TestGetIpInterfacesFiqlNodeLabel(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v2Path+"/ipinterfaces", ipInterfacesV2ListJSON)
	_, err := c.GetIpInterfaces(t.Context(), "node.label==router01.example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := req.query.Get("_s"); got != "node.label==router01.example.com" {
		t.Errorf("_s = %q", got)
	}
}

func TestGetIpInterfacesUsesV2URL(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v2Path+"/ipinterfaces", ipInterfacesV2ListJSON)
	if _, err := c.GetIpInterfaces(t.Context(), "", nil); err != nil {
		t.Fatal(err)
	}
	if req.path != v2Path+"/ipinterfaces" {
		t.Errorf("path = %q", req.path)
	}
}
