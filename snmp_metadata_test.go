package opennms

// Tests for the SNMP metadata methods – /api/v2/snmpmetadata.

import "testing"

const snmpMetadataEntryJSON = `{
	"nodeId": 1,
	"entries": [
		{
			"oid": ".1.3.6.1.2.1.1.1.0",
			"value": "Cisco IOS Software",
			"instance": "0"
		}
	]
}`

func TestGetSnmpMetadata(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v2Path+"/snmpmetadata/1", snmpMetadataEntryJSON)
	result, err := c.GetSnmpMetadata(t.Context(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if result["nodeId"].(float64) != 1 {
		t.Errorf("nodeId = %v, want 1", result["nodeId"])
	}
	entry := result["entries"].([]any)[0].(map[string]any)
	if entry["oid"] != ".1.3.6.1.2.1.1.1.0" {
		t.Errorf("oid = %v", entry["oid"])
	}
}
