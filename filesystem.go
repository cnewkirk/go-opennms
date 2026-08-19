package opennms

// Filesystem REST API – /rest/filesystem.
//
// Read and write configuration files in the OpenNMS etc directory.
// All calls require the FILESYSTEM EDITOR role.

import (
	"context"
	"net/url"
)

// GetFilesystemFiles returns the names of the configuration files
// accessible via the API. When changedOnly is true, only files that
// differ from the shipped defaults are listed.
func (c *Client) GetFilesystemFiles(ctx context.Context, changedOnly bool) ([]any, error) {
	var params url.Values
	if changedOnly {
		params = url.Values{"changedFilesOnly": {"true"}}
	}
	return c.getList(ctx, "filesystem", params, false)
}

// GetFilesystemExtensions returns the file extensions supported by
// the filesystem API.
func (c *Client) GetFilesystemExtensions(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "filesystem/extensions", nil, false)
}

// GetFilesystemHelp returns the Markdown help text for the named
// configuration file.
func (c *Client) GetFilesystemHelp(ctx context.Context, filename string) (string, error) {
	return c.getText(ctx, "filesystem/help", url.Values{"f": {filename}}, false, "")
}

// GetFilesystemContents returns the contents of the named
// configuration file (e.g. "discovery-configuration.xml").
func (c *Client) GetFilesystemContents(ctx context.Context, filename string) (string, error) {
	return c.getText(ctx, "filesystem/contents", url.Values{"f": {filename}}, false, "")
}

// UploadFilesystemContents creates or overwrites the named
// configuration file with content.
func (c *Client) UploadFilesystemContents(ctx context.Context, filename string, content []byte) error {
	files := map[string]fileUpload{
		"upload": {filename: filename, content: content},
	}
	_, err := c.postFiles(ctx, "filesystem/contents", files,
		url.Values{"f": {filename}}, false)
	return err
}

// DeleteFilesystemFile deletes the named configuration file.
func (c *Client) DeleteFilesystemFile(ctx context.Context, filename string) error {
	_, err := c.del(ctx, "filesystem/contents",
		url.Values{"f": {filename}}, nil, false, "")
	return err
}
