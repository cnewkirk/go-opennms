package opennms

// Alarm History REST API – /rest/alarms/history.

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// GetAlarmHistory returns the last known state of all active alarms.
//
// at is an optional millisecond epoch timestamp for historical
// lookup; 0 means now.
func (c *Client) GetAlarmHistory(ctx context.Context, at int) ([]any, error) {
	params := url.Values{}
	if at != 0 {
		params.Set("at", strconv.Itoa(at))
	}
	return c.getList(ctx, "alarms/history", params, false)
}

// GetAlarmHistoryAt returns the final known state of the alarm,
// optionally at a point in time.
//
// at is an optional millisecond epoch timestamp; 0 means now.
func (c *Client) GetAlarmHistoryAt(ctx context.Context, alarmID, at int) (map[string]any, error) {
	params := url.Values{}
	if at != 0 {
		params.Set("at", strconv.Itoa(at))
	}
	return c.getObject(ctx, fmt.Sprintf("alarms/history/%d", alarmID), params, false)
}

// GetAlarmHistoryStates returns all state transitions for the alarm.
func (c *Client) GetAlarmHistoryStates(ctx context.Context, alarmID int) ([]any, error) {
	return c.getList(ctx, fmt.Sprintf("alarms/history/%d/states", alarmID), nil, false)
}
