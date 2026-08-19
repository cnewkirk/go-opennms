package opennms

// Tests for the resources methods – /rest/resources.

import (
	"strings"
	"testing"
)

const resourcesTreeJSON = `{
	"resource": {
		"id": "node[1].interfaceSnmp[eth0-04013f75f101]",
		"label": "eth0 (04:01:3f:75:f1:01)",
		"name": "eth0-04013f75f101",
		"typeLabel": "SNMP Interface Data",
		"parentId": "node[1]",
		"stringPropertyAttributes": {
			"ifName": "eth0", "ifAlias": "uplink-to-core"
		},
		"externalValueAttributes": {},
		"rrdGraphAttributes": {
			"ifInOctets": {
				"name": "ifInOctets",
				"relativePath": "snmp/1/eth0-04013f75f101",
				"rrdFile": "ifInOctets.jrb"
			}
		},
		"children": {"resource": []}
	}
}`

const resourcesTestID = "node[1].interfaceSnmp[eth0-04013f75f101]"

func TestGetResources(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/resources", resourcesTreeJSON)
	result, err := c.GetResources(t.Context(), 1)
	if err != nil {
		t.Fatal(err)
	}
	resource := result["resource"].(map[string]any)
	if !strings.HasPrefix(resource["id"].(string), "node[1]") {
		t.Errorf("id = %v, want node[1] prefix", resource["id"])
	}
	if got := req.query.Get("depth"); got != "1" {
		t.Errorf("depth = %q, want 1", got)
	}
}

func TestGetResourcesCustomDepth(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/resources", resourcesTreeJSON)
	if _, err := c.GetResources(t.Context(), -1); err != nil {
		t.Fatal(err)
	}
	if got := req.query.Get("depth"); got != "-1" {
		t.Errorf("depth = %q, want -1", got)
	}
}

func TestGetResource(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/resources/"+resourcesTestID,
		resourcesTreeJSON)
	result, err := c.GetResource(t.Context(), resourcesTestID, -1)
	if err != nil {
		t.Fatal(err)
	}
	resource := result["resource"].(map[string]any)
	if resource["id"] != resourcesTestID {
		t.Errorf("id = %v, want %s", resource["id"], resourcesTestID)
	}
	if got := req.query.Get("depth"); got != "-1" {
		t.Errorf("depth = %q, want -1", got)
	}
}

func TestGetResourcesForNode(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/resources/fornode/1", resourcesTreeJSON)
	result, err := c.GetResourcesForNode(t.Context(), "1")
	if err != nil {
		t.Fatal(err)
	}
	resource := result["resource"].(map[string]any)
	if resource["parentId"] != "node[1]" {
		t.Errorf("parentId = %v, want node[1]", resource["parentId"])
	}
}

func TestGetResourcesSelect(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/resources/select", resourcesTreeJSON)
	_, err := c.GetResourcesSelect(t.Context(),
		[]string{"1", "2"}, nil, []string{"interfaceSnmp"}, []string{"ifAlias"})
	if err != nil {
		t.Fatal(err)
	}
	for param, want := range map[string]string{
		"nodes":            "1,2",
		"nodeSubresources": "interfaceSnmp",
		"stringProperties": "ifAlias",
	} {
		if got := req.query.Get(param); got != want {
			t.Errorf("%s = %q, want %q", param, got, want)
		}
	}
}

func TestDeleteResource(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "DELETE "+v1Path+"/resources/"+resourcesTestID, 204)
	if err := c.DeleteResource(t.Context(), resourcesTestID); err != nil {
		t.Fatal(err)
	}
	if req.method != "DELETE" {
		t.Errorf("method = %q, want DELETE", req.method)
	}
}
