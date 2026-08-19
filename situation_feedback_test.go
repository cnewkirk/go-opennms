package opennms

// Tests for the situation feedback methods – /rest/situation-feedback.

import (
	"encoding/json"
	"testing"
)

const situationFeedbackTagsJSON = `[{"name": "correct", "count": 5}]`

const situationFeedbackJSON = `[
	{
		"situationKey": "uei.opennms.org/alarms/situation::99",
		"alarmKey": "uei.opennms.org/nodes/nodeDown::1",
		"fingerprint": "abc123",
		"feedbackType": "CORRECT",
		"reason": "correlation is correct",
		"user": "admin",
		"timestamp": 1717234200000
	}
]`

func TestGetSituationFeedbackTags(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/situation-feedback/tags",
		situationFeedbackTagsJSON)
	result, err := c.GetSituationFeedbackTags(t.Context(), "")
	if err != nil {
		t.Fatal(err)
	}
	tag := result[0].(map[string]any)
	if tag["name"] != "correct" {
		t.Errorf("name = %v, want correct", tag["name"])
	}
}

func TestGetSituationFeedbackTagsWithPrefix(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleJSON(mux, "GET "+v1Path+"/situation-feedback/tags",
		situationFeedbackTagsJSON)
	_, err := c.GetSituationFeedbackTags(t.Context(), "cor")
	if err != nil {
		t.Fatal(err)
	}
	if got := req.query.Get("prefix"); got != "cor" {
		t.Errorf("prefix = %q, want cor", got)
	}
}

func TestGetSituationFeedback(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/situation-feedback/99",
		situationFeedbackJSON)
	result, err := c.GetSituationFeedback(t.Context(), 99)
	if err != nil {
		t.Fatal(err)
	}
	entry := result[0].(map[string]any)
	if entry["feedbackType"] != "CORRECT" {
		t.Errorf("feedbackType = %v, want CORRECT", entry["feedbackType"])
	}
}

func TestSubmitSituationFeedback(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v1Path+"/situation-feedback/99", 204)
	var feedback []any
	if err := json.Unmarshal([]byte(situationFeedbackJSON), &feedback); err != nil {
		t.Fatal(err)
	}
	if err := c.SubmitSituationFeedback(t.Context(), 99, feedback); err != nil {
		t.Fatal(err)
	}
	var body []any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	entry := body[0].(map[string]any)
	if entry["feedbackType"] != "CORRECT" {
		t.Errorf("feedbackType = %v, want CORRECT", entry["feedbackType"])
	}
}
