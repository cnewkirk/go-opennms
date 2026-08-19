package opennms

// Perspective Poller REST API v2 – /api/v2/perspectivepoller.

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// perspectivePollerTimeParams builds the optional start/end query
// parameters (ms since epoch; 0 = absent).
func perspectivePollerTimeParams(start, end int) url.Values {
	params := url.Values{}
	if start != 0 {
		params.Set("start", strconv.Itoa(start))
	}
	if end != 0 {
		params.Set("end", strconv.Itoa(end))
	}
	return params
}

// GetPerspectivePollerStatus returns perspective poller status for an
// application. start and end are optional times in ms since epoch;
// 0 = absent.
func (c *Client) GetPerspectivePollerStatus(ctx context.Context, appID, start, end int) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("perspectivepoller/%d", appID),
		perspectivePollerTimeParams(start, end), true)
}

// GetPerspectivePollerServiceStatus returns perspective poller status
// for a specific monitored service of an application. start and end
// are optional times in ms since epoch; 0 = absent.
func (c *Client) GetPerspectivePollerServiceStatus(ctx context.Context, appID, serviceID, start, end int) (map[string]any, error) {
	return c.getObject(ctx,
		fmt.Sprintf("perspectivepoller/%d/%d", appID, serviceID),
		perspectivePollerTimeParams(start, end), true)
}
