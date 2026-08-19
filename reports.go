package opennms

// Reports REST API – /rest/reports.

import (
	"context"
	"net/url"
	"strconv"
)

// Templates

// GetReportTemplates returns all available report templates.
func (c *Client) GetReportTemplates(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "reports", nil, false)
}

// GetReportTemplate returns the details of a report template. Pass
// userID to resolve the template details for that user; "" = absent.
func (c *Client) GetReportTemplate(ctx context.Context, reportID, userID string) (map[string]any, error) {
	var params url.Values
	if userID != "" {
		params = url.Values{"userId": {userID}}
	}
	return c.getObject(ctx, "reports/"+url.PathEscape(reportID),
		params, false)
}

// RunReport runs a report immediately and returns the rendered
// output (e.g. PDF or CSV data). format is the output format, e.g.
// "PDF", "SVG", "CSV", or "HTML"; parameters are report parameter
// maps with "name", "type", and "value" keys.
func (c *Client) RunReport(ctx context.Context, reportID, format string, parameters []any) ([]byte, error) {
	body := map[string]any{"format": format, "parameters": parameters}
	return c.postBytes(ctx, "reports/"+url.PathEscape(reportID),
		body, nil, false)
}

// Persisted reports

// GetPersistedReports returns all persisted (stored) report
// instances.
func (c *Client) GetPersistedReports(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "reports/persisted", nil, false)
}

// DeliverReport runs a report and delivers it (email, webhook, or
// persist). format is the output format, e.g. "PDF" or "CSV";
// parameters are report parameter maps with "name", "type", and
// "value" keys; deliveryOptions is the delivery options map —
// "instanceId" is required.
func (c *Client) DeliverReport(ctx context.Context, reportID, format string, parameters []any, deliveryOptions map[string]any) error {
	body := map[string]any{
		"id":              reportID,
		"format":          format,
		"parameters":      parameters,
		"deliveryOptions": deliveryOptions,
	}
	_, err := c.post(ctx, "reports/persisted", body, nil, false)
	return err
}

// DownloadReport downloads a persisted report by its catalog
// identifier and returns the rendered output (e.g. PDF or CSV data).
// Pass format to re-render the report in that format; "" = absent.
func (c *Client) DownloadReport(ctx context.Context, locatorID int, format string) ([]byte, error) {
	params := url.Values{"locatorId": {strconv.Itoa(locatorID)}}
	if format != "" {
		params.Set("format", format)
	}
	return c.getBytes(ctx, "reports/download", params, false)
}

// DeletePersistedReports deletes all persisted reports.
func (c *Client) DeletePersistedReports(ctx context.Context) error {
	_, err := c.del(ctx, "reports/persisted", nil, nil, false, "")
	return err
}

// DeletePersistedReport deletes a persisted report by its catalog
// identifier.
func (c *Client) DeletePersistedReport(ctx context.Context, reportID int) error {
	_, err := c.del(ctx, "reports/persisted/"+strconv.Itoa(reportID),
		nil, nil, false, "")
	return err
}

// Scheduled reports

// GetScheduledReports returns all scheduled report triggers.
func (c *Client) GetScheduledReports(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "reports/scheduled", nil, false)
}

// GetScheduledReport returns the details of a scheduled report
// trigger by name.
func (c *Client) GetScheduledReport(ctx context.Context, triggerName string) (map[string]any, error) {
	return c.getObject(ctx,
		"reports/scheduled/"+url.PathEscape(triggerName), nil, false)
}

// ScheduleReport creates a scheduled report. format is the output
// format, e.g. "PDF" or "CSV"; cronExpression is a Quartz cron
// expression for the schedule; parameters are report parameter maps
// with "name", "type", and "value" keys; deliveryOptions is the
// delivery options map — "instanceId" is required.
func (c *Client) ScheduleReport(ctx context.Context, reportID, format, cronExpression string, parameters []any, deliveryOptions map[string]any) error {
	body := map[string]any{
		"id":              reportID,
		"format":          format,
		"cronExpression":  cronExpression,
		"parameters":      parameters,
		"deliveryOptions": deliveryOptions,
	}
	_, err := c.post(ctx, "reports/scheduled", body, nil, false)
	return err
}

// UpdateScheduledReport updates an existing scheduled report trigger.
// data is the updated trigger definition ("parameters",
// "deliveryOptions", "cronExpression").
func (c *Client) UpdateScheduledReport(ctx context.Context, triggerName string, data map[string]any) error {
	_, err := c.put(ctx,
		"reports/scheduled/"+url.PathEscape(triggerName),
		data, nil, false)
	return err
}

// DeleteScheduledReports deletes all scheduled report triggers.
func (c *Client) DeleteScheduledReports(ctx context.Context) error {
	_, err := c.del(ctx, "reports/scheduled", nil, nil, false, "")
	return err
}

// DeleteScheduledReport deletes a scheduled report trigger by name.
func (c *Client) DeleteScheduledReport(ctx context.Context, triggerName string) error {
	_, err := c.del(ctx,
		"reports/scheduled/"+url.PathEscape(triggerName),
		nil, nil, false, "")
	return err
}
