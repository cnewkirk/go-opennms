package opennms

// Tests for the data choices methods – /rest/datachoices.

import (
	"encoding/json"
	"testing"
)

const datachoicesUsageReportJSON = `{
	"systemId": "e0f0c9d8-1234-5678-9abc-def012345678",
	"version": "2025.1.0",
	"nodes": 42,
	"alarms": 7
}`

const datachoicesUsageStatusJSON = `{
	"enabled": true,
	"initialNoticeAcknowledged": true
}`

const datachoicesUsageMetaJSON = `{
	"systemId": "Unique identifier of this OpenNMS instance",
	"nodes": "Number of nodes in the database"
}`

const datachoicesProductUpdateStatusJSON = `{
	"optedIn": false,
	"noticeAcknowledged": true
}`

func TestGetUsageStatisticsReport(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/datachoices", datachoicesUsageReportJSON)
	result, err := c.GetUsageStatisticsReport(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result["nodes"].(float64) != 42 {
		t.Errorf("nodes = %v, want 42", result["nodes"])
	}
}

func TestGetUsageStatisticsStatus(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/datachoices/status", datachoicesUsageStatusJSON)
	result, err := c.GetUsageStatisticsStatus(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result["enabled"] != true {
		t.Errorf("enabled = %v, want true", result["enabled"])
	}
}

func TestSetUsageStatisticsStatus(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v1Path+"/datachoices/status", 202)
	enabled, acked := false, true
	if err := c.SetUsageStatisticsStatus(t.Context(), &enabled, &acked); err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if len(body) != 2 || body["enabled"] != false ||
		body["initialNoticeAcknowledged"] != true {
		t.Errorf("body = %v", body)
	}
}

func TestGetUsageStatisticsMeta(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/datachoices/meta", datachoicesUsageMetaJSON)
	result, err := c.GetUsageStatisticsMeta(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := result["systemId"]; !ok {
		t.Errorf("systemId missing from %v", result)
	}
}

func TestGetProductUpdateStatus(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/datachoices/productupdate/status",
		datachoicesProductUpdateStatusJSON)
	result, err := c.GetProductUpdateStatus(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result["optedIn"] != false {
		t.Errorf("optedIn = %v, want false", result["optedIn"])
	}
}

func TestSetProductUpdateStatus(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v1Path+"/datachoices/productupdate/status", 202)
	optedIn := true
	if err := c.SetProductUpdateStatus(t.Context(), &optedIn, nil); err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if len(body) != 1 || body["optedIn"] != true {
		t.Errorf("body = %v, want {optedIn: true}", body)
	}
}

func TestSubmitProductUpdateEnrollment(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v1Path+"/datachoices/productupdate/submit", 202)
	err := c.SubmitProductUpdateEnrollment(t.Context(),
		map[string]any{"consent": true, "email": "noc@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["consent"] != true {
		t.Errorf("consent = %v, want true", body["consent"])
	}
}
