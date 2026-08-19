package opennms

// Tests for the ifservices methods – /rest/ifservices + v2.

import "testing"

const ifservicesListJSON = `{
	"service": [
		{"id": 201, "serviceType": {"id": 1, "name": "ICMP"},
		 "status": "A", "statusLong": "Active",
		 "ipInterfaceId": 101, "nodeId": 1}
	],
	"totalCount": 1, "count": 1, "offset": 0
}`

const ifservicesListV2JSON = ifservicesListJSON

func TestGetIfservices(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/ifservices", ifservicesListJSON)
	result, err := c.GetIfservices(t.Context(), nil)
	if err != nil {
		t.Fatal(err)
	}
	service := result["service"].([]any)[0].(map[string]any)
	if service["id"].(float64) != 201 {
		t.Errorf("id = %v, want 201", service["id"])
	}
}

func TestUpdateIfservices(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "PUT "+v1Path+"/ifservices", 204)
	err := c.UpdateIfservices(t.Context(), map[string]string{"status": "A"})
	if err != nil {
		t.Fatal(err)
	}
	if req.formBody(t).Get("status") != "A" {
		t.Errorf("body = %q, want status=A", req.body)
	}
}

func TestGetIfservicesV2(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v2Path+"/ifservices", ifservicesListV2JSON)
	result, err := c.GetIfservicesV2(t.Context(), "node.id==1", nil)
	if err != nil {
		t.Fatal(err)
	}
	service := result["service"].([]any)[0].(map[string]any)
	if service["id"].(float64) != 201 {
		t.Errorf("id = %v, want 201", service["id"])
	}
	if got := req.query.Get("_s"); got != "node.id==1" {
		t.Errorf("_s = %q, want node.id==1", got)
	}
}
