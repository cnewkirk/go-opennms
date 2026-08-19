package opennms

// Groups REST API – /rest/groups.

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// groupsXMLReplacer escapes &, < and > in XML text content, matching
// Python's xml.sax.saxutils.escape.
var groupsXMLReplacer = strings.NewReplacer(
	"&", "&amp;", "<", "&lt;", ">", "&gt;")

// groupsMemberList normalizes a group's "user" entry into a string
// slice.
func groupsMemberList(v any) []string {
	switch members := v.(type) {
	case []string:
		return members
	case []any:
		out := make([]string, len(members))
		for i, m := range members {
			out[i] = fmt.Sprint(m)
		}
		return out
	}
	return nil
}

// GetGroups lists all user groups.
func (c *Client) GetGroups(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "groups", nil, false)
}

// GetGroup returns a specific group by name.
func (c *Client) GetGroup(ctx context.Context, groupName string) (map[string]any, error) {
	return c.getObject(ctx, "groups/"+url.PathEscape(groupName), nil, false)
}

// CreateGroup creates a new user group.
//
// The body is sent as XML — POST /rest/groups does not accept JSON on
// any OpenNMS version. The XML document is built internally; callers
// pass a plain map with keys "name" (required), "comments", and
// "user" (list of member usernames). Example:
// map[string]any{"name": "network-ops", "comments": "Network operations team"}.
func (c *Client) CreateGroup(ctx context.Context, group map[string]any) (map[string]any, error) {
	var b strings.Builder
	b.WriteString("<group>")
	fmt.Fprintf(&b, "<name>%s</name>",
		groupsXMLReplacer.Replace(fmt.Sprint(group["name"])))
	if comments, ok := group["comments"]; ok {
		fmt.Fprintf(&b, "<comments>%s</comments>",
			groupsXMLReplacer.Replace(fmt.Sprint(comments)))
	}
	for _, member := range groupsMemberList(group["user"]) {
		fmt.Fprintf(&b, "<user>%s</user>",
			groupsXMLReplacer.Replace(member))
	}
	b.WriteString("</group>")
	return asObject(c.postText(ctx, "groups", b.String(),
		"application/xml", false, "", nil))
}

// UpdateGroup updates group metadata (e.g. the comments field).
//
// The body is sent form-encoded — PUT /rest/groups/{name} does not
// accept JSON on any OpenNMS version. group is the map of group
// fields to change, e.g. map[string]string{"comments": "..."}.
func (c *Client) UpdateGroup(ctx context.Context, groupName string, group map[string]string) error {
	form := url.Values{}
	for k, v := range group {
		form.Set(k, v)
	}
	_, err := c.putForm(ctx, "groups/"+url.PathEscape(groupName),
		form, nil, false)
	return err
}

// DeleteGroup deletes a user group.
func (c *Client) DeleteGroup(ctx context.Context, groupName string) error {
	_, err := c.del(ctx, "groups/"+url.PathEscape(groupName),
		nil, nil, false, "")
	return err
}

// GetGroupUsers lists users in the group.
func (c *Client) GetGroupUsers(ctx context.Context, groupName string) (map[string]any, error) {
	return c.getObject(ctx,
		fmt.Sprintf("groups/%s/users", url.PathEscape(groupName)),
		nil, false)
}

// AddUserToGroup adds the user to the group.
func (c *Client) AddUserToGroup(ctx context.Context, groupName, username string) error {
	_, err := c.put(ctx,
		fmt.Sprintf("groups/%s/users/%s",
			url.PathEscape(groupName), url.PathEscape(username)),
		nil, nil, false)
	return err
}

// RemoveUserFromGroup removes the user from the group.
func (c *Client) RemoveUserFromGroup(ctx context.Context, groupName, username string) error {
	_, err := c.del(ctx,
		fmt.Sprintf("groups/%s/users/%s",
			url.PathEscape(groupName), url.PathEscape(username)),
		nil, nil, false, "")
	return err
}

// GetGroupCategories lists surveillance categories associated with
// the group.
func (c *Client) GetGroupCategories(ctx context.Context, groupName string) (map[string]any, error) {
	return c.getObject(ctx,
		fmt.Sprintf("groups/%s/categories", url.PathEscape(groupName)),
		nil, false)
}

// AddCategoryToGroup associates the category with the group.
func (c *Client) AddCategoryToGroup(ctx context.Context, groupName, categoryName string) error {
	_, err := c.put(ctx,
		fmt.Sprintf("groups/%s/categories/%s",
			url.PathEscape(groupName), url.PathEscape(categoryName)),
		nil, nil, false)
	return err
}

// RemoveCategoryFromGroup removes the category from the group.
func (c *Client) RemoveCategoryFromGroup(ctx context.Context, groupName, categoryName string) error {
	_, err := c.del(ctx,
		fmt.Sprintf("groups/%s/categories/%s",
			url.PathEscape(groupName), url.PathEscape(categoryName)),
		nil, nil, false, "")
	return err
}
