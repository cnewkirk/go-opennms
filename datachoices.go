package opennms

// Data Choices REST API – /rest/datachoices.
//
// Controls the anonymous usage-statistics collection and
// product-update enrollment settings.

import "context"

// GetUsageStatisticsReport returns the usage-statistics report for
// this system.
func (c *Client) GetUsageStatisticsReport(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "datachoices", nil, false)
}

// GetUsageStatisticsStatus returns the usage-statistics collection
// status: an object with "enabled" and "initialNoticeAcknowledged"
// keys.
func (c *Client) GetUsageStatisticsStatus(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "datachoices/status", nil, false)
}

// SetUsageStatisticsStatus updates the usage-statistics collection
// status. enabled enables or disables usage-statistics collection;
// initialNoticeAcknowledged marks the initial notice as
// acknowledged. nil leaves a field unchanged.
func (c *Client) SetUsageStatisticsStatus(ctx context.Context, enabled, initialNoticeAcknowledged *bool) error {
	body := map[string]any{}
	if enabled != nil {
		body["enabled"] = *enabled
	}
	if initialNoticeAcknowledged != nil {
		body["initialNoticeAcknowledged"] = *initialNoticeAcknowledged
	}
	_, err := c.post(ctx, "datachoices/status", body, nil, false)
	return err
}

// GetUsageStatisticsMeta returns metadata describing the
// usage-statistics fields.
func (c *Client) GetUsageStatisticsMeta(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "datachoices/meta", nil, false)
}

// GetProductUpdateStatus returns the product-update enrollment
// status: an object with "optedIn" and "noticeAcknowledged" keys.
func (c *Client) GetProductUpdateStatus(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "datachoices/productupdate/status", nil, false)
}

// SetProductUpdateStatus updates the product-update enrollment
// status. optedIn opts in to or out of product-update enrollment;
// noticeAcknowledged marks the enrollment notice as acknowledged.
// nil leaves a field unchanged.
func (c *Client) SetProductUpdateStatus(ctx context.Context, optedIn, noticeAcknowledged *bool) error {
	body := map[string]any{}
	if optedIn != nil {
		body["optedIn"] = *optedIn
	}
	if noticeAcknowledged != nil {
		body["noticeAcknowledged"] = *noticeAcknowledged
	}
	_, err := c.post(ctx, "datachoices/productupdate/status", body, nil, false)
	return err
}

// SubmitProductUpdateEnrollment submits the product-update
// enrollment form: a dict with keys such as "consent", "firstName",
// "lastName", "email", "company".
func (c *Client) SubmitProductUpdateEnrollment(ctx context.Context, formData map[string]any) error {
	_, err := c.post(ctx, "datachoices/productupdate/submit", formData, nil, false)
	return err
}
