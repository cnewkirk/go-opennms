package opennms

// Tests for the notifications methods – /rest/notifications.

import (
	"strings"
	"testing"
)

const notificationJSON = `{
	"notifyId": 601,
	"textMsg": "Node router01.example.com is down.",
	"subject": "node down",
	"numericMsg": null,
	"pageTime": "2024-06-01T08:01:00.000+0000",
	"respondTime": "2024-06-01T08:10:00.000+0000",
	"answeredBy": "admin",
	"ipAddress": "192.168.1.1",
	"queueId": "default",
	"notifConfigName": "nodeDown",
	"eventUei": "uei.opennms.org/nodes/nodeDown",
	"ackUser": "admin",
	"ackTime": "2024-06-01T08:10:00.000+0000",
	"ackId": 701,
	"serviceType": {"id": 1, "name": "ICMP"},
	"node": {"id": 1, "label": "router01.example.com"}
}`

const notificationListJSON = `{
	"notification": [` + notificationJSON + `],
	"totalCount": 1, "count": 1, "offset": 0
}`

func TestGetNotificationsDefault(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/notifications", notificationListJSON)
	result, err := c.GetNotifications(t.Context(), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	notification := result["notification"].([]any)[0].(map[string]any)
	if notification["notifyId"].(float64) != 601 {
		t.Errorf("notifyId = %v, want 601", notification["notifyId"])
	}
	if got := req.query.Get("limit"); got != "10" {
		t.Errorf("limit = %q, want 10", got)
	}
	if got := req.query.Get("offset"); got != "0" {
		t.Errorf("offset = %q, want 0", got)
	}
}

func TestGetNotificationsWithFilters(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/notifications", notificationListJSON)
	_, err := c.GetNotifications(t.Context(),
		&ListOptions{Limit: 50, OrderBy: "pageTime", Order: "descending"},
		map[string]string{"answered": "false"})
	if err != nil {
		t.Fatal(err)
	}
	if got := req.query.Get("limit"); got != "50" {
		t.Errorf("limit = %q, want 50", got)
	}
	if got := req.query.Get("answered"); got != "false" {
		t.Errorf("answered = %q, want false", got)
	}
}

func TestGetNotification(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/notifications/601", notificationJSON)
	result, err := c.GetNotification(t.Context(), 601)
	if err != nil {
		t.Fatal(err)
	}
	if result["notifyId"].(float64) != 601 {
		t.Errorf("notifyId = %v, want 601", result["notifyId"])
	}
	if result["answeredBy"] != "admin" {
		t.Errorf("answeredBy = %v, want admin", result["answeredBy"])
	}
}

func TestGetNotificationCount(t *testing.T) {
	mux, c := newTestClient(t)
	handleText(mux, "GET "+v1Path+"/notifications/count", "8")
	count, err := c.GetNotificationCount(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if count != 8 {
		t.Errorf("count = %d, want 8", count)
	}
}

func TestTriggerDestinationPath(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux,
		"POST "+v1Path+"/notifications/destination-paths/Email/trigger", 200)
	if err := c.TriggerDestinationPath(t.Context(), "Email"); err != nil {
		t.Fatal(err)
	}
	if req.method != "POST" {
		t.Errorf("method = %q, want POST", req.method)
	}
	if !strings.Contains(req.path, "Email/trigger") {
		t.Errorf("path = %q, want to contain Email/trigger", req.path)
	}
}
