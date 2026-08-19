package opennms

// Tests for the EnLinkd methods – /api/v2/enlinkd.

import "testing"

const enlinkdLldpLinkJSON = `{
	"lldpLocalPort": "Gi0/0/0",
	"lldpLocalPortUrl": "/opennms/element/node.jsp?node=1",
	"lldpRemChassisId": "00:11:22:33:44:55",
	"lldpRemChassisIdUrl": "/opennms/element/node.jsp?node=2",
	"lldpRemInfo": "router02.example.com",
	"lldpRemPort": "Gi0/0/1",
	"lldpRemPortUrl": "/opennms/element/snmpinterface.jsp?node=2&ifindex=2",
	"lldpCreateTime": "2024-06-01T08:00:00.000+0000",
	"lldpLastPollTime": "2024-06-01T09:00:00.000+0000"
}`

const enlinkdCdpLinkJSON = `{
	"cdpLocalPort": "GigabitEthernet0/0",
	"cdpLocalPortUrl": "/opennms/element/node.jsp?node=1",
	"cdpCacheDevice": "router02.example.com",
	"cdpCacheDeviceUrl": "/opennms/element/node.jsp?node=2",
	"cdpCacheDevicePort": "GigabitEthernet0/1",
	"cdpCacheDevicePortUrl": "/opennms/element/snmpinterface.jsp?node=2&ifindex=2",
	"cdpCachePlatform": "Cisco IOS Software, ISR4331",
	"cdpCreateTime": "2024-06-01T08:00:00.000+0000",
	"cdpLastPollTime": "2024-06-01T09:00:00.000+0000"
}`

const enlinkdOspfLinkJSON = `{
	"ospfLocalPort": "192.168.1.1",
	"ospfLocalPortUrl": "/opennms/element/node.jsp?node=1",
	"ospfRemRouterId": "192.168.1.2",
	"ospfRemRouterUrl": "/opennms/element/node.jsp?node=2",
	"ospfRemPort": "192.168.1.2",
	"ospfRemPortUrl": "/opennms/element/snmpinterface.jsp?node=2&ifindex=2",
	"ospfLinkInfo": "point-to-point",
	"ospfLinkCreateTime": "2024-06-01T08:00:00.000+0000",
	"ospfLinkLastPollTime": "2024-06-01T09:00:00.000+0000"
}`

const enlinkdIsisLinkJSON = `{
	"isisCircIfIndex": 6,
	"isisCircAdminState": "on",
	"isisISAdjNeighSysID": "0100.0200.0300",
	"isisISAdjNeighSysType": "l1_l2IntermediateSystem",
	"isisISAdjNeighSysUrl": "/opennms/element/node.jsp?node=2",
	"isisISAdjNeighSNPAAddress": "00:11:22:33:44:55",
	"isisISAdjNeighPort": "Gi0/0/1",
	"isisISAdjState": "up",
	"isisISAdjNbrExtendedCircID": 1,
	"isisISAdjUrl": "/opennms/element/snmpinterface.jsp?node=2&ifindex=2",
	"isisLinkCreateTime": "2024-06-01T08:00:00.000+0000",
	"isisLinkLastPollTime": "2024-06-01T09:00:00.000+0000"
}`

const enlinkdBridgeLinkJSON = `{
	"bridgeLocalPort": "FastEthernet0/1",
	"bridgeLocalPortUrl": "/opennms/element/node.jsp?node=1",
	"BridgeLinkRemoteNodes": [{
		"bridgeRemote": "switch02.example.com",
		"bridgeRemoteUrl": "/opennms/element/node.jsp?node=2",
		"bridgeRemotePort": "FastEthernet0/2",
		"bridgeRemotePortUrl": "/opennms/element/snmpinterface.jsp?node=2&ifindex=2"
	}],
	"bridgeInfo": "STP designated bridge",
	"bridgeLinkCreateTime": "2024-06-01T08:00:00.000+0000",
	"bridgeLinkLastPollTime": "2024-06-01T09:00:00.000+0000"
}`

const enlinkdLldpElementJSON = `{
	"lldpChassisId": "00:11:22:33:44:55",
	"lldpSysName": "router01",
	"lldpCreateTime": "2024-06-01T08:00:00.000+0000",
	"lldpLastPollTime": "2024-06-01T09:00:00.000+0000"
}`

const enlinkdCdpElementJSON = `{
	"cdpGlobalRun": "true",
	"cdpGlobalDeviceId": "router01.example.com",
	"cdpGlobalDeviceIdFormat": "other",
	"cdpCreateTime": "2024-06-01T08:00:00.000+0000",
	"cdpLastPollTime": "2024-06-01T09:00:00.000+0000"
}`

const enlinkdOspfElementJSON = `{
	"ospfRouterId": "192.168.1.1",
	"ospfVersionNumber": 2,
	"ospfAdminStat": "enabled",
	"ospfCreateTime": "2024-06-01T08:00:00.000+0000",
	"ospfLastPollTime": "2024-06-01T09:00:00.000+0000"
}`

const enlinkdIsisElementJSON = `{
	"isisSysID": "0100.0200.0300",
	"isisSysAdminState": "on",
	"isisCreateTime": "2024-06-01T08:00:00.000+0000",
	"isisLastPollTime": "2024-06-01T09:00:00.000+0000"
}`

const enlinkdBridgeElementJSON = `{
	"baseBridgeAddress": "00:11:22:33:44:55",
	"baseNumPorts": 24,
	"baseType": "transparent-only",
	"stpProtocolSpecification": "ieee8021d",
	"stpPriority": 32768,
	"stpDesignatedRoot": "8000.001122334455",
	"stpRootCost": 0,
	"stpRootPort": 0,
	"vlan": 1,
	"vlanname": "default",
	"bridgeNodeCreateTime": "2024-06-01T08:00:00.000+0000",
	"bridgeNodeLastPollTime": "2024-06-01T09:00:00.000+0000"
}`

const enlinkdDtoJSON = `{
	"lldpLinkNodes": [` + enlinkdLldpLinkJSON + `],
	"cdpLinkNodes": [` + enlinkdCdpLinkJSON + `],
	"ospfLinkNodes": [` + enlinkdOspfLinkJSON + `],
	"isisLinkNodes": [` + enlinkdIsisLinkJSON + `],
	"bridgeLinkNodes": [` + enlinkdBridgeLinkJSON + `],
	"lldpElementNode": ` + enlinkdLldpElementJSON + `,
	"cdpElementNode": ` + enlinkdCdpElementJSON + `,
	"ospfElementNode": ` + enlinkdOspfElementJSON + `,
	"isisElementNode": ` + enlinkdIsisElementJSON + `,
	"bridgeElementNodes": [` + enlinkdBridgeElementJSON + `]
}`

func TestGetNodeEnlinkd(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/enlinkd/1", enlinkdDtoJSON)
	result, err := c.GetNodeEnlinkd(t.Context(), "1")
	if err != nil {
		t.Fatal(err)
	}
	lldp := result["lldpLinkNodes"].([]any)[0].(map[string]any)
	if lldp["lldpRemChassisId"] != "00:11:22:33:44:55" {
		t.Errorf("lldpRemChassisId = %v", lldp["lldpRemChassisId"])
	}
	cdp := result["cdpElementNode"].(map[string]any)
	if cdp["cdpGlobalDeviceId"] != "router01.example.com" {
		t.Errorf("cdpGlobalDeviceId = %v", cdp["cdpGlobalDeviceId"])
	}
}

func TestGetNodeLldpLinks(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/enlinkd/lldp_links/1",
		`[`+enlinkdLldpLinkJSON+`]`)
	result, err := c.GetNodeLldpLinks(t.Context(), "1")
	if err != nil {
		t.Fatal(err)
	}
	link := result[0].(map[string]any)
	if link["lldpLocalPort"] != "Gi0/0/0" {
		t.Errorf("lldpLocalPort = %v", link["lldpLocalPort"])
	}
	if link["lldpRemPort"] != "Gi0/0/1" {
		t.Errorf("lldpRemPort = %v", link["lldpRemPort"])
	}
}

func TestGetNodeCdpLinks(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/enlinkd/cdp_links/1",
		`[`+enlinkdCdpLinkJSON+`]`)
	result, err := c.GetNodeCdpLinks(t.Context(), "1")
	if err != nil {
		t.Fatal(err)
	}
	link := result[0].(map[string]any)
	if link["cdpCacheDevice"] != "router02.example.com" {
		t.Errorf("cdpCacheDevice = %v", link["cdpCacheDevice"])
	}
	if link["cdpCachePlatform"] != "Cisco IOS Software, ISR4331" {
		t.Errorf("cdpCachePlatform = %v", link["cdpCachePlatform"])
	}
}

func TestGetNodeOspfLinks(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/enlinkd/ospf_links/1",
		`[`+enlinkdOspfLinkJSON+`]`)
	result, err := c.GetNodeOspfLinks(t.Context(), "1")
	if err != nil {
		t.Fatal(err)
	}
	link := result[0].(map[string]any)
	if link["ospfRemRouterId"] != "192.168.1.2" {
		t.Errorf("ospfRemRouterId = %v", link["ospfRemRouterId"])
	}
	if link["ospfLinkInfo"] != "point-to-point" {
		t.Errorf("ospfLinkInfo = %v", link["ospfLinkInfo"])
	}
}

func TestGetNodeIsisLinks(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/enlinkd/isis_links/1",
		`[`+enlinkdIsisLinkJSON+`]`)
	result, err := c.GetNodeIsisLinks(t.Context(), "1")
	if err != nil {
		t.Fatal(err)
	}
	link := result[0].(map[string]any)
	if link["isisCircIfIndex"].(float64) != 6 {
		t.Errorf("isisCircIfIndex = %v, want 6", link["isisCircIfIndex"])
	}
	if link["isisISAdjState"] != "up" {
		t.Errorf("isisISAdjState = %v", link["isisISAdjState"])
	}
}

func TestGetNodeBridgeLinks(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/enlinkd/bridge_links/1",
		`[`+enlinkdBridgeLinkJSON+`]`)
	result, err := c.GetNodeBridgeLinks(t.Context(), "1")
	if err != nil {
		t.Fatal(err)
	}
	link := result[0].(map[string]any)
	if link["bridgeLocalPort"] != "FastEthernet0/1" {
		t.Errorf("bridgeLocalPort = %v", link["bridgeLocalPort"])
	}
	remote := link["BridgeLinkRemoteNodes"].([]any)[0].(map[string]any)
	if remote["bridgeRemote"] != "switch02.example.com" {
		t.Errorf("bridgeRemote = %v", remote["bridgeRemote"])
	}
}

func TestGetNodeLldpElement(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/enlinkd/lldp_elems/1", enlinkdLldpElementJSON)
	result, err := c.GetNodeLldpElement(t.Context(), "1")
	if err != nil {
		t.Fatal(err)
	}
	if result["lldpChassisId"] != "00:11:22:33:44:55" {
		t.Errorf("lldpChassisId = %v", result["lldpChassisId"])
	}
	if result["lldpSysName"] != "router01" {
		t.Errorf("lldpSysName = %v", result["lldpSysName"])
	}
}

func TestGetNodeCdpElement(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/enlinkd/cdp_elems/1", enlinkdCdpElementJSON)
	result, err := c.GetNodeCdpElement(t.Context(), "1")
	if err != nil {
		t.Fatal(err)
	}
	if result["cdpGlobalDeviceId"] != "router01.example.com" {
		t.Errorf("cdpGlobalDeviceId = %v", result["cdpGlobalDeviceId"])
	}
	if result["cdpGlobalRun"] != "true" {
		t.Errorf("cdpGlobalRun = %v", result["cdpGlobalRun"])
	}
}

func TestGetNodeOspfElement(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/enlinkd/ospf_elems/1", enlinkdOspfElementJSON)
	result, err := c.GetNodeOspfElement(t.Context(), "1")
	if err != nil {
		t.Fatal(err)
	}
	if result["ospfRouterId"] != "192.168.1.1" {
		t.Errorf("ospfRouterId = %v", result["ospfRouterId"])
	}
	if result["ospfVersionNumber"].(float64) != 2 {
		t.Errorf("ospfVersionNumber = %v, want 2", result["ospfVersionNumber"])
	}
}

func TestGetNodeIsisElement(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/enlinkd/isis_elems/1", enlinkdIsisElementJSON)
	result, err := c.GetNodeIsisElement(t.Context(), "1")
	if err != nil {
		t.Fatal(err)
	}
	if result["isisSysID"] != "0100.0200.0300" {
		t.Errorf("isisSysID = %v", result["isisSysID"])
	}
	if result["isisSysAdminState"] != "on" {
		t.Errorf("isisSysAdminState = %v", result["isisSysAdminState"])
	}
}

func TestGetNodeBridgeElements(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/enlinkd/bridge_elems/1",
		`[`+enlinkdBridgeElementJSON+`]`)
	result, err := c.GetNodeBridgeElements(t.Context(), "1")
	if err != nil {
		t.Fatal(err)
	}
	elem := result[0].(map[string]any)
	if elem["baseBridgeAddress"] != "00:11:22:33:44:55" {
		t.Errorf("baseBridgeAddress = %v", elem["baseBridgeAddress"])
	}
	if elem["baseNumPorts"].(float64) != 24 {
		t.Errorf("baseNumPorts = %v, want 24", elem["baseNumPorts"])
	}
	if elem["vlan"].(float64) != 1 {
		t.Errorf("vlan = %v, want 1", elem["vlan"])
	}
}
