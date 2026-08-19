package types

import (
	"reflect"
	"testing"
)

func TestMapDropsUnsetFields(t *testing.T) {
	m := Map(Node{Label: "web01"})
	if len(m) != 1 {
		t.Fatalf("expected 1 key, got %d: %v", len(m), m)
	}
	if m["label"] != "web01" {
		t.Fatalf("expected label=web01, got %v", m["label"])
	}
}

func TestMapKebabCaseTags(t *testing.T) {
	m := Map(RequisitionNode{
		ForeignID: "router01",
		NodeLabel: "router01.example.com",
	})
	want := map[string]any{
		"foreign-id": "router01",
		"node-label": "router01.example.com",
	}
	if !reflect.DeepEqual(m, want) {
		t.Fatalf("expected %v, got %v", want, m)
	}
}

func TestMapFalseBoolSurvives(t *testing.T) {
	enabled := false
	m := Map(ClassificationGroup{Name: "test", Enabled: &enabled})
	v, ok := m["enabled"]
	if !ok {
		t.Fatalf("expected enabled key, got %v", m)
	}
	if v != false {
		t.Fatalf("expected enabled=false, got %v", v)
	}
}

func TestMapNestedStructSlices(t *testing.T) {
	m := Map(RequisitionNode{
		ForeignID: "router01",
		Interface: []RequisitionInterface{
			{IPAddr: "10.0.0.1", SnmpPrimary: "P", Status: 1},
		},
	})
	ifaces, ok := m["interface"].([]any)
	if !ok {
		t.Fatalf("expected interface to be []any, got %T", m["interface"])
	}
	if len(ifaces) != 1 {
		t.Fatalf("expected 1 interface, got %d", len(ifaces))
	}
	iface, ok := ifaces[0].(map[string]any)
	if !ok {
		t.Fatalf("expected interface entry to be map[string]any, got %T",
			ifaces[0])
	}
	want := map[string]any{
		"ip-addr":      "10.0.0.1",
		"snmp-primary": "P",
		"status":       float64(1),
	}
	if !reflect.DeepEqual(iface, want) {
		t.Fatalf("expected %v, got %v", want, iface)
	}
}

func TestMapNilInput(t *testing.T) {
	if m := Map(nil); m != nil {
		t.Fatalf("expected nil for nil input, got %v", m)
	}
}
