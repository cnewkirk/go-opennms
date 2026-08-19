package opennms

// Situations REST API v2 – /api/v2/situations.

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// situationsAlarmIDList renders alarm IDs as the comma-separated
// string the situations endpoints expect.
func situationsAlarmIDList(alarmIDs []int) string {
	parts := make([]string, len(alarmIDs))
	for i, id := range alarmIDs {
		parts[i] = strconv.Itoa(id)
	}
	return strings.Join(parts, ",")
}

// GetSituations lists situations (v2).
func (c *Client) GetSituations(ctx context.Context, opts *ListOptions) (map[string]any, error) {
	return c.getObject(ctx, "situations", listParams(opts), true)
}

// CreateSituation creates a new situation from a list of alarm IDs.
// description and diagnosticText are optional; "" omits them.
func (c *Client) CreateSituation(ctx context.Context, alarmIDs []int, description, diagnosticText string) (map[string]any, error) {
	body := map[string]any{"alarmIdList": situationsAlarmIDList(alarmIDs)}
	if description != "" {
		body["description"] = description
	}
	if diagnosticText != "" {
		body["diagnosticText"] = diagnosticText
	}
	return asObject(c.post(ctx, "situations/create", body, nil, true))
}

// AddAlarmsToSituation links additional alarm IDs to an existing
// situation. feedback is an optional feedback string; "" omits it.
func (c *Client) AddAlarmsToSituation(ctx context.Context, situationID int, alarmIDs []int, feedback string) (map[string]any, error) {
	body := map[string]any{
		"situationId": situationID,
		"alarmIdList": situationsAlarmIDList(alarmIDs),
	}
	if feedback != "" {
		body["feedback"] = feedback
	}
	return asObject(c.post(ctx, "situations/associateAlarm", body, nil, true))
}

// ClearSituation clears a situation by ID.
func (c *Client) ClearSituation(ctx context.Context, situationID int) error {
	_, err := c.post(ctx, "situations/clear",
		map[string]any{"situationId": situationID}, nil, true)
	return err
}

// ClearSituationAlarms removes the given alarm IDs from a situation
// and clears them.
func (c *Client) ClearSituationAlarms(ctx context.Context, situationID int, alarmIDs []int) error {
	body := map[string]any{
		"situationId": situationID,
		"alarmIdList": situationsAlarmIDList(alarmIDs),
	}
	_, err := c.post(ctx, "situations/alarms/clear", body, nil, true)
	return err
}

// AcceptSituation accepts (acknowledges) a situation.
func (c *Client) AcceptSituation(ctx context.Context, situationID int) error {
	_, err := c.post(ctx, fmt.Sprintf("situations/accepted/%d", situationID),
		nil, nil, true)
	return err
}

// RemoveAlarmsFromSituation removes specific alarm IDs from a
// situation without clearing them.
func (c *Client) RemoveAlarmsFromSituation(ctx context.Context, situationID int, alarmIDs []int) error {
	params := url.Values{
		"situationId": {strconv.Itoa(situationID)},
		"alarmIdList": {situationsAlarmIDList(alarmIDs)},
	}
	_, err := c.del(ctx, "situations/removeAlarm", params, nil, true, "")
	return err
}
