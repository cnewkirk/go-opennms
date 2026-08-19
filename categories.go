package opennms

// Categories REST API – /rest/categories.

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// categoriesAttrEscaper escapes the XML attribute-special characters
// (mirroring Python's xml.sax.saxutils.quoteattr).
var categoriesAttrEscaper = strings.NewReplacer(
	"&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;")

// Categories CRUD

// GetCategories lists all configured surveillance categories.
func (c *Client) GetCategories(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "categories", nil, false)
}

// GetCategory returns a specific category by name.
func (c *Client) GetCategory(ctx context.Context, category string) (map[string]any, error) {
	return c.getObject(ctx, "categories/"+url.PathEscape(category), nil, false)
}

// CreateCategory adds a new surveillance category.
//
// category example:
// map[string]any{"name": "Production", "authorizedGroups": []any{}}.
func (c *Client) CreateCategory(ctx context.Context, category map[string]any) (map[string]any, error) {
	return asObject(c.post(ctx, "categories", category, nil, false))
}

// UpdateCategory updates a category. data holds the category fields
// to change; the endpoint consumes form-encoded data.
func (c *Client) UpdateCategory(ctx context.Context, category string, data map[string]string) error {
	form := url.Values{}
	for k, v := range data {
		form.Set(k, v)
	}
	_, err := c.putForm(ctx, "categories/"+url.PathEscape(category), form, nil, false)
	return err
}

// DeleteCategory deletes a category.
func (c *Client) DeleteCategory(ctx context.Context, category string) error {
	_, err := c.del(ctx, "categories/"+url.PathEscape(category), nil, nil, false, "")
	return err
}

// Category ↔ Node associations

// GetNodeCategoriesList returns all categories for the node.
func (c *Client) GetNodeCategoriesList(ctx context.Context, nodeID int) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("categories/nodes/%d", nodeID), nil, false)
}

// GetCategoryForNode returns a specific category for the node.
func (c *Client) GetCategoryForNode(ctx context.Context, category string, nodeID int) (map[string]any, error) {
	return c.getObject(ctx,
		fmt.Sprintf("categories/%s/nodes/%d", url.PathEscape(category), nodeID),
		nil, false)
}

// AssociateCategoryWithNode associates the category with the node.
// The endpoint consumes only XML; the category element is built from
// the category name.
func (c *Client) AssociateCategoryWithNode(ctx context.Context, category string, nodeID int) error {
	xml := fmt.Sprintf(`<category name="%s"/>`,
		categoriesAttrEscaper.Replace(category))
	_, err := c.putText(ctx,
		fmt.Sprintf("categories/%s/nodes/%d", url.PathEscape(category), nodeID),
		xml, "application/xml", false)
	return err
}

// DissociateCategoryFromNode removes the category from the node.
func (c *Client) DissociateCategoryFromNode(ctx context.Context, category string, nodeID int) error {
	_, err := c.del(ctx,
		fmt.Sprintf("categories/%s/nodes/%d", url.PathEscape(category), nodeID),
		nil, nil, false, "")
	return err
}

// Category ↔ Group associations

// GetCategoriesForGroup returns the categories associated with the
// user group.
func (c *Client) GetCategoriesForGroup(ctx context.Context, group string) (map[string]any, error) {
	return c.getObject(ctx, "categories/groups/"+url.PathEscape(group), nil, false)
}

// AssociateCategoryWithGroup associates the category with the user
// group.
func (c *Client) AssociateCategoryWithGroup(ctx context.Context, category, group string) error {
	_, err := c.put(ctx,
		fmt.Sprintf("categories/%s/groups/%s",
			url.PathEscape(category), url.PathEscape(group)),
		nil, nil, false)
	return err
}

// DissociateCategoryFromGroup removes the category from the user
// group.
func (c *Client) DissociateCategoryFromGroup(ctx context.Context, category, group string) error {
	_, err := c.del(ctx,
		fmt.Sprintf("categories/%s/groups/%s",
			url.PathEscape(category), url.PathEscape(group)),
		nil, nil, false, "")
	return err
}
