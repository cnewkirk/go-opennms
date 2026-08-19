package opennms

// Scheduled Outages REST API – /rest/sched-outages.

import (
	"context"
	"net/url"
)

// Scheduled outages CRUD

// GetSchedOutages lists all configured scheduled outages.
func (c *Client) GetSchedOutages(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "sched-outages", nil, false)
}

// GetSchedOutage returns a specific scheduled outage by name.
func (c *Client) GetSchedOutage(ctx context.Context, outageName string) (map[string]any, error) {
	return c.getObject(ctx, "sched-outages/"+url.PathEscape(outageName), nil, false)
}

// CreateSchedOutage adds a new (or replaces an existing) scheduled
// outage.
//
// Valid "type" values: "weekly", "monthly", "specific", "daily". For
// "weekly" outages the "day" key in each time entry is the weekday
// name (lowercase); for "monthly" outages "day" is an integer
// day-of-month; for "specific" outages "begins" and "ends" are full
// datetime strings ("DD-Mon-YYYY HH:MM:SS"). Example:
//
//	map[string]any{
//	    "name": "Weekend-Maintenance",
//	    "type": "weekly",
//	    "time": []any{
//	        map[string]any{"day": "saturday", "begins": "00:00:00", "ends": "23:59:59"},
//	        map[string]any{"day": "sunday", "begins": "00:00:00", "ends": "23:59:59"},
//	    },
//	    "node":      []any{map[string]any{"id": 1}, map[string]any{"id": 2}},
//	    "interface": []any{map[string]any{"address": "192.168.0.1"}},
//	}
func (c *Client) CreateSchedOutage(ctx context.Context, outage map[string]any) (map[string]any, error) {
	return asObject(c.post(ctx, "sched-outages", outage, nil, false))
}

// DeleteSchedOutage deletes a scheduled outage.
func (c *Client) DeleteSchedOutage(ctx context.Context, outageName string) error {
	_, err := c.del(ctx, "sched-outages/"+url.PathEscape(outageName), nil, nil, false, "")
	return err
}

// Daemon associations

// AssociateSchedOutageCollectd associates the outage with the given
// collectd package.
func (c *Client) AssociateSchedOutageCollectd(ctx context.Context, outageName, pkg string) error {
	path := "sched-outages/" + url.PathEscape(outageName) + "/collectd/" + url.PathEscape(pkg)
	_, err := c.put(ctx, path, nil, nil, false)
	return err
}

// DissociateSchedOutageCollectd removes the collectd package
// association from the outage.
func (c *Client) DissociateSchedOutageCollectd(ctx context.Context, outageName, pkg string) error {
	path := "sched-outages/" + url.PathEscape(outageName) + "/collectd/" + url.PathEscape(pkg)
	_, err := c.del(ctx, path, nil, nil, false, "")
	return err
}

// AssociateSchedOutagePollerd associates the outage with the given
// pollerd package.
func (c *Client) AssociateSchedOutagePollerd(ctx context.Context, outageName, pkg string) error {
	path := "sched-outages/" + url.PathEscape(outageName) + "/pollerd/" + url.PathEscape(pkg)
	_, err := c.put(ctx, path, nil, nil, false)
	return err
}

// DissociateSchedOutagePollerd removes the pollerd package
// association from the outage.
func (c *Client) DissociateSchedOutagePollerd(ctx context.Context, outageName, pkg string) error {
	path := "sched-outages/" + url.PathEscape(outageName) + "/pollerd/" + url.PathEscape(pkg)
	_, err := c.del(ctx, path, nil, nil, false, "")
	return err
}

// AssociateSchedOutageThreshd associates the outage with the given
// threshd package.
func (c *Client) AssociateSchedOutageThreshd(ctx context.Context, outageName, pkg string) error {
	path := "sched-outages/" + url.PathEscape(outageName) + "/threshd/" + url.PathEscape(pkg)
	_, err := c.put(ctx, path, nil, nil, false)
	return err
}

// DissociateSchedOutageThreshd removes the threshd package
// association from the outage.
func (c *Client) DissociateSchedOutageThreshd(ctx context.Context, outageName, pkg string) error {
	path := "sched-outages/" + url.PathEscape(outageName) + "/threshd/" + url.PathEscape(pkg)
	_, err := c.del(ctx, path, nil, nil, false, "")
	return err
}

// AssociateSchedOutageNotifd associates the outage with the
// notifications daemon.
func (c *Client) AssociateSchedOutageNotifd(ctx context.Context, outageName string) error {
	path := "sched-outages/" + url.PathEscape(outageName) + "/notifd"
	_, err := c.put(ctx, path, nil, nil, false)
	return err
}

// DissociateSchedOutageNotifd removes the notifications daemon
// association from the outage.
func (c *Client) DissociateSchedOutageNotifd(ctx context.Context, outageName string) error {
	path := "sched-outages/" + url.PathEscape(outageName) + "/notifd"
	_, err := c.del(ctx, path, nil, nil, false, "")
	return err
}
