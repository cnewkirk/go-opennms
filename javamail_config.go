package opennms

// Javamail Configuration REST API – /rest/config/javamail.

import (
	"context"
	"net/url"
)

const javamailConfigPath = "config/javamail"

// javamailForm converts a map of key/value pairs into form data.
func javamailForm(data map[string]string) url.Values {
	form := url.Values{}
	for k, v := range data {
		form.Set(k, v)
	}
	return form
}

// Default configuration

// GetJavamailDefaultConfig returns the default Javamail
// configuration.
func (c *Client) GetJavamailDefaultConfig(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, javamailConfigPath, nil, false)
}

// SetJavamailDefaultConfig updates the default Javamail
// configuration from the given configuration data.
func (c *Client) SetJavamailDefaultConfig(ctx context.Context, data map[string]any) error {
	_, err := c.post(ctx, javamailConfigPath, data, nil, false)
	return err
}

// Read-mail configs

// GetJavamailReadmails lists all read-mail configurations.
func (c *Client) GetJavamailReadmails(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, javamailConfigPath+"/readmails", nil, false)
}

// GetJavamailReadmail returns a specific read-mail configuration by
// name.
func (c *Client) GetJavamailReadmail(ctx context.Context, name string) (map[string]any, error) {
	return c.getObject(ctx, javamailConfigPath+"/readmails/"+url.PathEscape(name), nil, false)
}

// CreateJavamailReadmail creates a new read-mail configuration.
func (c *Client) CreateJavamailReadmail(ctx context.Context, data map[string]any) error {
	_, err := c.post(ctx, javamailConfigPath+"/readmails", data, nil, false)
	return err
}

// UpdateJavamailReadmail updates the named read-mail configuration.
// data holds form-encoded key/value pairs to update.
func (c *Client) UpdateJavamailReadmail(ctx context.Context, name string, data map[string]string) error {
	_, err := c.putForm(ctx, javamailConfigPath+"/readmails/"+url.PathEscape(name),
		javamailForm(data), nil, false)
	return err
}

// DeleteJavamailReadmail deletes a read-mail configuration.
func (c *Client) DeleteJavamailReadmail(ctx context.Context, name string) error {
	_, err := c.del(ctx, javamailConfigPath+"/readmails/"+url.PathEscape(name),
		nil, nil, false, "")
	return err
}

// Send-mail configs

// GetJavamailSendmails lists all send-mail configurations.
func (c *Client) GetJavamailSendmails(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, javamailConfigPath+"/sendmails", nil, false)
}

// GetJavamailSendmail returns a specific send-mail configuration by
// name.
func (c *Client) GetJavamailSendmail(ctx context.Context, name string) (map[string]any, error) {
	return c.getObject(ctx, javamailConfigPath+"/sendmails/"+url.PathEscape(name), nil, false)
}

// CreateJavamailSendmail creates a new send-mail configuration.
func (c *Client) CreateJavamailSendmail(ctx context.Context, data map[string]any) error {
	_, err := c.post(ctx, javamailConfigPath+"/sendmails", data, nil, false)
	return err
}

// UpdateJavamailSendmail updates the named send-mail configuration.
// data holds form-encoded key/value pairs to update.
func (c *Client) UpdateJavamailSendmail(ctx context.Context, name string, data map[string]string) error {
	_, err := c.putForm(ctx, javamailConfigPath+"/sendmails/"+url.PathEscape(name),
		javamailForm(data), nil, false)
	return err
}

// DeleteJavamailSendmail deletes a send-mail configuration.
func (c *Client) DeleteJavamailSendmail(ctx context.Context, name string) error {
	_, err := c.del(ctx, javamailConfigPath+"/sendmails/"+url.PathEscape(name),
		nil, nil, false, "")
	return err
}

// End-to-end configs

// GetJavamailEnd2ends lists all end-to-end mail test configurations.
func (c *Client) GetJavamailEnd2ends(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, javamailConfigPath+"/end2ends", nil, false)
}

// GetJavamailEnd2end returns a specific end-to-end mail test
// configuration by name.
func (c *Client) GetJavamailEnd2end(ctx context.Context, name string) (map[string]any, error) {
	return c.getObject(ctx, javamailConfigPath+"/end2ends/"+url.PathEscape(name), nil, false)
}

// CreateJavamailEnd2end creates a new end-to-end mail test
// configuration.
func (c *Client) CreateJavamailEnd2end(ctx context.Context, data map[string]any) error {
	_, err := c.post(ctx, javamailConfigPath+"/end2ends", data, nil, false)
	return err
}

// UpdateJavamailEnd2end updates the named end-to-end mail test
// configuration. data holds form-encoded key/value pairs to update.
func (c *Client) UpdateJavamailEnd2end(ctx context.Context, name string, data map[string]string) error {
	_, err := c.putForm(ctx, javamailConfigPath+"/end2ends/"+url.PathEscape(name),
		javamailForm(data), nil, false)
	return err
}

// DeleteJavamailEnd2end deletes an end-to-end mail test
// configuration.
func (c *Client) DeleteJavamailEnd2end(ctx context.Context, name string) error {
	_, err := c.del(ctx, javamailConfigPath+"/end2ends/"+url.PathEscape(name),
		nil, nil, false, "")
	return err
}
