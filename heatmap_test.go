package opennms

// Tests for the heatmap methods – /rest/heatmap (read-only GET).

import "testing"

const heatmapResponseJSON = `{
	"heatmapEntry": [
		{"id": 2, "label": "Production", "nodesTotal": 15,
		 "nodesWithOutages": 2, "nodesWithAlarms": 5,
		 "maximumSeverity": "MAJOR"}
	]
}`

func heatmapAssertEntry(t *testing.T, result map[string]any) {
	t.Helper()
	if _, ok := result["heatmapEntry"]; !ok {
		t.Errorf("result = %v, want heatmapEntry key", result)
	}
}

func TestGetHeatmapOutagesCategories(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/heatmap/outages/categories", heatmapResponseJSON)
	result, err := c.GetHeatmapOutagesCategories(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	entry := result["heatmapEntry"].([]any)[0].(map[string]any)
	if entry["label"] != "Production" {
		t.Errorf("label = %v, want Production", entry["label"])
	}
}

func TestGetHeatmapOutagesForeignSources(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/heatmap/outages/foreignSources", heatmapResponseJSON)
	result, err := c.GetHeatmapOutagesForeignSources(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	heatmapAssertEntry(t, result)
}

func TestGetHeatmapOutagesMonitoredServices(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/heatmap/outages/monitoredServices", heatmapResponseJSON)
	result, err := c.GetHeatmapOutagesMonitoredServices(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	heatmapAssertEntry(t, result)
}

func TestGetHeatmapOutagesNodesByCategory(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/heatmap/outages/nodesByCategory/Production",
		heatmapResponseJSON)
	result, err := c.GetHeatmapOutagesNodesByCategory(t.Context(), "Production")
	if err != nil {
		t.Fatal(err)
	}
	heatmapAssertEntry(t, result)
	if req.path != v1Path+"/heatmap/outages/nodesByCategory/Production" {
		t.Errorf("path = %q", req.path)
	}
}

func TestGetHeatmapOutagesNodesByForeignSource(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/heatmap/outages/nodesByForeignSource/Routers",
		heatmapResponseJSON)
	result, err := c.GetHeatmapOutagesNodesByForeignSource(t.Context(), "Routers")
	if err != nil {
		t.Fatal(err)
	}
	heatmapAssertEntry(t, result)
}

func TestGetHeatmapOutagesNodesByService(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/heatmap/outages/nodesByMonitoredService/ICMP",
		heatmapResponseJSON)
	result, err := c.GetHeatmapOutagesNodesByService(t.Context(), "ICMP")
	if err != nil {
		t.Fatal(err)
	}
	heatmapAssertEntry(t, result)
}

func TestGetHeatmapAlarmsCategories(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/heatmap/alarms/categories", heatmapResponseJSON)
	result, err := c.GetHeatmapAlarmsCategories(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	entry := result["heatmapEntry"].([]any)[0].(map[string]any)
	if entry["maximumSeverity"] != "MAJOR" {
		t.Errorf("maximumSeverity = %v, want MAJOR", entry["maximumSeverity"])
	}
}

func TestGetHeatmapAlarmsForeignSources(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/heatmap/alarms/foreignSources", heatmapResponseJSON)
	result, err := c.GetHeatmapAlarmsForeignSources(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	heatmapAssertEntry(t, result)
}

func TestGetHeatmapAlarmsMonitoredServices(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/heatmap/alarms/monitoredServices", heatmapResponseJSON)
	result, err := c.GetHeatmapAlarmsMonitoredServices(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	heatmapAssertEntry(t, result)
}

func TestGetHeatmapAlarmsNodesByCategory(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/heatmap/alarms/nodesByCategory/Production",
		heatmapResponseJSON)
	result, err := c.GetHeatmapAlarmsNodesByCategory(t.Context(), "Production")
	if err != nil {
		t.Fatal(err)
	}
	heatmapAssertEntry(t, result)
	if req.path != v1Path+"/heatmap/alarms/nodesByCategory/Production" {
		t.Errorf("path = %q", req.path)
	}
}

func TestGetHeatmapAlarmsNodesByForeignSource(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/heatmap/alarms/nodesByForeignSource/Routers",
		heatmapResponseJSON)
	result, err := c.GetHeatmapAlarmsNodesByForeignSource(t.Context(), "Routers")
	if err != nil {
		t.Fatal(err)
	}
	heatmapAssertEntry(t, result)
}

func TestGetHeatmapAlarmsNodesByService(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/heatmap/alarms/nodesByMonitoredService/ICMP",
		heatmapResponseJSON)
	result, err := c.GetHeatmapAlarmsNodesByService(t.Context(), "ICMP")
	if err != nil {
		t.Fatal(err)
	}
	heatmapAssertEntry(t, result)
}
