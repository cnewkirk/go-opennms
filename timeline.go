package opennms

// Outage Timelines REST API – /rest/timeline.

import (
	"context"
	"fmt"
	"net/url"
)

// GetTimelineHeader returns the outage timeline header image as PNG
// data. start and end are epoch timestamps (seconds); width is the
// timeline width in pixels.
func (c *Client) GetTimelineHeader(ctx context.Context, start, end, width int) ([]byte, error) {
	return c.getBytes(ctx,
		fmt.Sprintf("timeline/header/%d/%d/%d", start, end, width),
		nil, false)
}

// GetTimelineImage returns the outage timeline image for a monitored
// service as PNG data. nodeID identifies the node, ipAddress the
// interface, and serviceID the monitored service; start and end are
// epoch timestamps (seconds); width is the timeline width in pixels.
func (c *Client) GetTimelineImage(ctx context.Context, nodeID int, ipAddress string, serviceID, start, end, width int) ([]byte, error) {
	return c.getBytes(ctx,
		fmt.Sprintf("timeline/image/%d/%s/%d/%d/%d/%d",
			nodeID, url.PathEscape(ipAddress), serviceID, start, end, width),
		nil, false)
}

// GetTimelineEmpty returns an empty outage timeline image as PNG
// data, used for services that are not monitored. start and end are
// epoch timestamps (seconds); width is the timeline width in pixels.
func (c *Client) GetTimelineEmpty(ctx context.Context, start, end, width int) ([]byte, error) {
	return c.getBytes(ctx,
		fmt.Sprintf("timeline/empty/%d/%d/%d", start, end, width),
		nil, false)
}

// GetTimelineHTML returns the raw HTML embedding the timeline image.
// nodeID identifies the node, ipAddress the interface, and serviceID
// the monitored service; start and end are epoch timestamps
// (seconds); width is the timeline width in pixels.
func (c *Client) GetTimelineHTML(ctx context.Context, nodeID int, ipAddress string, serviceID, start, end, width int) (string, error) {
	return c.getText(ctx,
		fmt.Sprintf("timeline/html/%d/%s/%d/%d/%d/%d",
			nodeID, url.PathEscape(ipAddress), serviceID, start, end, width),
		nil, false, "text/html")
}
