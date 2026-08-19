package opennms

// Tests for the foreign sources configuration methods –
// /rest/foreignSourcesConfig.

import "testing"

const foreignSourcesConfigPoliciesJSON = `{
	"plugin": [
		{
			"name": "Match IP Interface",
			"class": "org.opennms.netmgt.provision.persist.policies.MatchingIpInterfacePolicy"
		}
	]
}`

const foreignSourcesConfigDetectorsJSON = `{
	"plugin": [
		{
			"name": "ICMP",
			"class": "org.opennms.netmgt.provision.detector.icmp.IcmpDetector"
		}
	]
}`

const foreignSourcesConfigServicesJSON = `{"service": ["ICMP", "SNMP", "HTTP"]}`

const foreignSourcesConfigAssetsJSON = `{"asset": ["manufacturer", "vendor", "serialNumber"]}`

const foreignSourcesConfigCategoriesJSON = `{"category": ["Production", "Development"]}`

func foreignSourcesConfigContains(list []any, want string) bool {
	for _, v := range list {
		if v == want {
			return true
		}
	}
	return false
}

func TestGetForeignSourceConfigPolicies(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/foreignSourcesConfig/policies",
		foreignSourcesConfigPoliciesJSON)
	result, err := c.GetForeignSourceConfigPolicies(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	plugin := result["plugin"].([]any)[0].(map[string]any)
	if plugin["name"] != "Match IP Interface" {
		t.Errorf("name = %v, want Match IP Interface", plugin["name"])
	}
}

func TestGetForeignSourceConfigDetectors(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/foreignSourcesConfig/detectors",
		foreignSourcesConfigDetectorsJSON)
	result, err := c.GetForeignSourceConfigDetectors(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	plugin := result["plugin"].([]any)[0].(map[string]any)
	if plugin["name"] != "ICMP" {
		t.Errorf("name = %v, want ICMP", plugin["name"])
	}
}

func TestGetForeignSourceConfigServices(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/foreignSourcesConfig/services/Routers",
		foreignSourcesConfigServicesJSON)
	result, err := c.GetForeignSourceConfigServices(t.Context(), "Routers")
	if err != nil {
		t.Fatal(err)
	}
	if !foreignSourcesConfigContains(result["service"].([]any), "ICMP") {
		t.Errorf("service = %v, want ICMP present", result["service"])
	}
}

func TestGetForeignSourceConfigAssets(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/foreignSourcesConfig/assets",
		foreignSourcesConfigAssetsJSON)
	result, err := c.GetForeignSourceConfigAssets(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if !foreignSourcesConfigContains(result["asset"].([]any), "manufacturer") {
		t.Errorf("asset = %v, want manufacturer present", result["asset"])
	}
}

func TestGetForeignSourceConfigCategories(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/foreignSourcesConfig/categories",
		foreignSourcesConfigCategoriesJSON)
	result, err := c.GetForeignSourceConfigCategories(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if !foreignSourcesConfigContains(result["category"].([]any), "Production") {
		t.Errorf("category = %v, want Production present", result["category"])
	}
}
