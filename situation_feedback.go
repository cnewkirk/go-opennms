package opennms

// Situation Feedback REST API – /rest/situation-feedback.

import (
	"context"
	"fmt"
	"net/url"
)

// GetSituationFeedbackTags lists situation feedback tags. Pass a
// non-empty prefix to filter tags.
func (c *Client) GetSituationFeedbackTags(ctx context.Context, prefix string) ([]any, error) {
	var params url.Values
	if prefix != "" {
		params = url.Values{"prefix": {prefix}}
	}
	return c.getList(ctx, "situation-feedback/tags", params, false)
}

// GetSituationFeedback returns the feedback for a specific situation,
// identified by the alarm ID of the situation.
func (c *Client) GetSituationFeedback(ctx context.Context, situationID int) ([]any, error) {
	return c.getList(ctx, fmt.Sprintf("situation-feedback/%d", situationID), nil, false)
}

// SubmitSituationFeedback submits feedback for the situation with the
// given alarm ID.
//
// feedback is a list of feedback entries, each with keys such as
// "alarmKey", "fingerprint", "feedbackType", "reason", "user",
// "timestamp".
func (c *Client) SubmitSituationFeedback(ctx context.Context, situationID int, feedback []any) error {
	_, err := c.post(ctx, fmt.Sprintf("situation-feedback/%d", situationID),
		feedback, nil, false)
	return err
}
