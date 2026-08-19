package opennms

// Tests for the acknowledgement methods – /rest/acks.

import "testing"

const acksAckJSON = `{
	"id": 701,
	"ackTime": "2024-06-01T08:10:00.000+0000",
	"ackUser": "admin",
	"ackType": "ALARM",
	"ackAction": "ACKNOWLEDGE",
	"refId": 42
}`

const acksAckListJSON = `{
	"count": 1, "offset": 0, "totalCount": 1,
	"ack": [` + acksAckJSON + `]
}`

func TestGetAcksDefault(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/acks", acksAckListJSON)
	result, err := c.GetAcks(t.Context(), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	ack := result["ack"].([]any)[0].(map[string]any)
	if ack["id"].(float64) != 701 {
		t.Errorf("id = %v, want 701", ack["id"])
	}
	if got := req.query.Get("limit"); got != "10" {
		t.Errorf("limit = %q, want 10", got)
	}
	if got := req.query.Get("offset"); got != "0" {
		t.Errorf("offset = %q, want 0", got)
	}
}

func TestGetAcksWithFilters(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/acks", acksAckListJSON)
	_, err := c.GetAcks(t.Context(), &ListOptions{Limit: 25},
		map[string]string{"ackUser": "admin"})
	if err != nil {
		t.Fatal(err)
	}
	if got := req.query.Get("limit"); got != "25" {
		t.Errorf("limit = %q, want 25", got)
	}
	if got := req.query.Get("ackUser"); got != "admin" {
		t.Errorf("ackUser = %q, want admin", got)
	}
}

func TestGetAck(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/acks/701", acksAckJSON)
	result, err := c.GetAck(t.Context(), 701)
	if err != nil {
		t.Fatal(err)
	}
	if result["id"].(float64) != 701 {
		t.Errorf("id = %v, want 701", result["id"])
	}
	if result["ackAction"] != "ACKNOWLEDGE" {
		t.Errorf("ackAction = %v, want ACKNOWLEDGE", result["ackAction"])
	}
}

func TestGetAckCount(t *testing.T) {
	mux, c := newTestClient(t)
	handleText(mux, "GET "+v1Path+"/acks/count", "23")
	count, err := c.GetAckCount(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if count != 23 {
		t.Errorf("count = %d, want 23", count)
	}
}

func TestCreateAckAlarm(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "POST "+v1Path+"/acks", formType, 200, acksAckJSON)
	if _, err := c.CreateAck(t.Context(), "ack", 42, 0); err != nil {
		t.Fatal(err)
	}
	form := req.formBody(t)
	if form.Get("action") != "ack" || form.Get("alarmId") != "42" {
		t.Errorf("form = %v", form)
	}
}

func TestCreateAckNotification(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "POST "+v1Path+"/acks", formType, 200, acksAckJSON)
	if _, err := c.CreateAck(t.Context(), "ack", 0, 601); err != nil {
		t.Fatal(err)
	}
	form := req.formBody(t)
	if form.Get("action") != "ack" || form.Get("notifId") != "601" {
		t.Errorf("form = %v", form)
	}
}

func TestCreateAckEscalate(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "POST "+v1Path+"/acks", formType, 200, acksAckJSON)
	if _, err := c.CreateAck(t.Context(), "esc", 42, 0); err != nil {
		t.Fatal(err)
	}
	form := req.formBody(t)
	if form.Get("action") != "esc" || form.Get("alarmId") != "42" {
		t.Errorf("form = %v", form)
	}
}

func TestAckNotification(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "POST "+v1Path+"/acks", formType, 200, acksAckJSON)
	if _, err := c.AckNotification(t.Context(), 601); err != nil {
		t.Fatal(err)
	}
	form := req.formBody(t)
	if form.Get("action") != "ack" || form.Get("notifId") != "601" {
		t.Errorf("form = %v", form)
	}
}

func TestUnackNotification(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleContract(mux, "POST "+v1Path+"/acks", formType, 200, acksAckJSON)
	if _, err := c.UnackNotification(t.Context(), 601); err != nil {
		t.Fatal(err)
	}
	form := req.formBody(t)
	if form.Get("action") != "unack" || form.Get("notifId") != "601" {
		t.Errorf("form = %v", form)
	}
}
