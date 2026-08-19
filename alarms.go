package opennms

// Alarms REST API – /rest/alarms and /api/v2/alarms.

import (
	"context"
	"fmt"
	"net/url"
)

// GetAlarms lists alarms (v1).
//
// filters are additional Hibernate query filters passed directly as
// query parameters (e.g. "severity": "MAJOR"). Pass "comparator" to
// change the match type (eq/ilike/…).
func (c *Client) GetAlarms(ctx context.Context, opts *ListOptions, filters map[string]string) (map[string]any, error) {
	params := mergeFilters(listParams(opts), filters)
	return c.getObject(ctx, "alarms", params, false)
}

// GetAlarm returns the alarm with the given ID.
func (c *Client) GetAlarm(ctx context.Context, alarmID int) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("alarms/%d", alarmID), nil, false)
}

// GetAlarmCount returns the total number of alarms.
func (c *Client) GetAlarmCount(ctx context.Context) (int, error) {
	return c.getCount(ctx, "alarms/count", false)
}

// Single-alarm actions (v1 PUT; the API requires the flags as
// form-encoded body data, not query parameters)

// AckAlarm acknowledges the alarm. Pass ackUser to acknowledge on
// behalf of that user (requires admin role); "" acknowledges as the
// authenticated user.
func (c *Client) AckAlarm(ctx context.Context, alarmID int, ackUser string) error {
	form := url.Values{"ack": {"true"}}
	if ackUser != "" {
		form.Set("ackUser", ackUser)
	}
	_, err := c.putForm(ctx, fmt.Sprintf("alarms/%d", alarmID), form, nil, false)
	return err
}

// UnackAlarm removes acknowledgement from the alarm.
func (c *Client) UnackAlarm(ctx context.Context, alarmID int) error {
	_, err := c.putForm(ctx, fmt.Sprintf("alarms/%d", alarmID),
		url.Values{"ack": {"false"}}, nil, false)
	return err
}

// ClearAlarm clears the alarm (sets severity to CLEARED).
func (c *Client) ClearAlarm(ctx context.Context, alarmID int) error {
	_, err := c.putForm(ctx, fmt.Sprintf("alarms/%d", alarmID),
		url.Values{"clear": {"true"}}, nil, false)
	return err
}

// EscalateAlarm escalates the severity of the alarm by one step.
func (c *Client) EscalateAlarm(ctx context.Context, alarmID int) error {
	_, err := c.putForm(ctx, fmt.Sprintf("alarms/%d", alarmID),
		url.Values{"escalate": {"true"}}, nil, false)
	return err
}

// Bulk alarm actions

func (c *Client) bulkAlarmAction(ctx context.Context, flag string, filters map[string]string) error {
	form := url.Values{flag: {"true"}}
	for k, v := range filters {
		form.Set(k, v)
	}
	_, err := c.putForm(ctx, "alarms", form, nil, false)
	return err
}

// BulkAckAlarms acknowledges all alarms matching the given filters
// (Hibernate query filters passed as form fields, e.g.
// "severity": "MAJOR").
func (c *Client) BulkAckAlarms(ctx context.Context, filters map[string]string) error {
	return c.bulkAlarmAction(ctx, "ack", filters)
}

// BulkUnackAlarms removes acknowledgement from all alarms matching
// the given filters.
func (c *Client) BulkUnackAlarms(ctx context.Context, filters map[string]string) error {
	form := url.Values{"ack": {"false"}}
	for k, v := range filters {
		form.Set(k, v)
	}
	_, err := c.putForm(ctx, "alarms", form, nil, false)
	return err
}

// BulkClearAlarms clears all alarms matching the given filters.
func (c *Client) BulkClearAlarms(ctx context.Context, filters map[string]string) error {
	return c.bulkAlarmAction(ctx, "clear", filters)
}

// BulkEscalateAlarms escalates all alarms matching the given filters.
func (c *Client) BulkEscalateAlarms(ctx context.Context, filters map[string]string) error {
	return c.bulkAlarmAction(ctx, "escalate", filters)
}

// v2 alarms (FIQL filtering)

// GetAlarmsV2 lists alarms using the v2 API with an optional FIQL
// filter string (e.g. "alarm.severity==MAJOR"); "" lists all.
func (c *Client) GetAlarmsV2(ctx context.Context, fiql string, opts *ListOptions) (map[string]any, error) {
	params := listParams(opts)
	if fiql != "" {
		params.Set("_s", fiql)
	}
	return c.getObject(ctx, "alarms", params, true)
}

// GetAlarmV2 returns a single alarm by ID using the v2 API.
func (c *Client) GetAlarmV2(ctx context.Context, alarmID int) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("alarms/%d", alarmID), nil, true)
}
