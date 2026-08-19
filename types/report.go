package types

// Reports and Data Choices APIs (keys are camelCase).

// ReportParameter is the report parameter payload for
// Client.RunReport, Client.DeliverReport, and
// Client.ScheduleReport.
type ReportParameter struct {
	// Name is the parameter name from the report template.
	// Required.
	Name string `json:"name,omitempty"`
	// Type is the parameter type: "string", "integer", "float",
	// "double", "timezone", or "date". Required.
	Type string `json:"type,omitempty"`
	// Value is the parameter value; shape depends on Type.
	Value any `json:"value,omitempty"`
}

// ReportDeliveryOptions is the report delivery options payload for
// Client.DeliverReport and Client.ScheduleReport.
type ReportDeliveryOptions struct {
	// InstanceID is the identifier for the report instance.
	// Required.
	InstanceID string `json:"instanceId,omitempty"`
	// SendMail delivers the report via email.
	SendMail *bool `json:"sendMail,omitempty"`
	// MailTo is the recipient address (required if SendMail).
	MailTo string `json:"mailTo,omitempty"`
	// Webhook delivers the report via HTTP POST.
	Webhook *bool `json:"webhook,omitempty"`
	// WebhookURL is the webhook target URL (required if Webhook).
	WebhookURL string `json:"webhookUrl,omitempty"`
	// Persist stores the report on the server.
	Persist *bool `json:"persist,omitempty"`
	// Format is the delivery format override.
	Format string `json:"format,omitempty"`
}

// ProductUpdateEnrollment is the product-update enrollment form
// payload for Client.SubmitProductUpdateEnrollment.
type ProductUpdateEnrollment struct {
	// Consent is the consent to the enrollment. Required.
	Consent *bool `json:"consent,omitempty"`
	// FirstName is the first name.
	FirstName string `json:"firstName,omitempty"`
	// LastName is the last name.
	LastName string `json:"lastName,omitempty"`
	// Email is the contact email address.
	Email string `json:"email,omitempty"`
	// Company is the company name.
	Company string `json:"company,omitempty"`
}
