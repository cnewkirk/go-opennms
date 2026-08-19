package opennms

// Users REST API – /rest/users.

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// GetUsers lists all users.
func (c *Client) GetUsers(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "users", nil, false)
}

// GetUser returns a specific user by username.
func (c *Client) GetUser(ctx context.Context, username string) (map[string]any, error) {
	return c.getObject(ctx, "users/"+url.PathEscape(username), nil, false)
}

// usersXMLEscape escapes &, < and > in XML character data, matching
// Python's xml.sax.saxutils.escape.
var usersXMLEscape = strings.NewReplacer(
	"&", "&amp;", "<", "&lt;", ">", "&gt;").Replace

// usersXMLValue renders a user attribute value for the XML body:
// bools become lowercase "true"/"false", everything else its
// default string form.
func usersXMLValue(v any) string {
	if b, ok := v.(bool); ok {
		return strconv.FormatBool(b)
	}
	return fmt.Sprint(v)
}

// usersStringList normalizes a list-valued user attribute
// (duty-schedule, role) given as []string or []any.
func usersStringList(v any) []string {
	switch list := v.(type) {
	case []string:
		return list
	case []any:
		out := make([]string, len(list))
		for i, e := range list {
			out[i] = fmt.Sprint(e)
		}
		return out
	}
	return nil
}

// usersBuildXML builds the <user> document for POST /rest/users,
// keeping the element order of the reference wrapper.
func usersBuildXML(user map[string]any) string {
	var b strings.Builder
	b.WriteString("<user>")
	for _, key := range [...]string{"user-id", "full-name", "user-comments",
		"email", "password", "passwordSalt"} {
		if value, ok := user[key]; ok {
			fmt.Fprintf(&b, "<%s>%s</%s>",
				key, usersXMLEscape(usersXMLValue(value)), key)
		}
	}
	for _, schedule := range usersStringList(user["duty-schedule"]) {
		fmt.Fprintf(&b, "<duty-schedule>%s</duty-schedule>",
			usersXMLEscape(schedule))
	}
	for _, role := range usersStringList(user["role"]) {
		fmt.Fprintf(&b, "<role>%s</role>", usersXMLEscape(role))
	}
	b.WriteString("</user>")
	return b.String()
}

// CreateUser creates a new user.
//
// user is the user attribute map. Required key: "user-id". Optional
// keys: "full-name", "user-comments", "password", "email",
// "duty-schedule" (list of schedule strings), "role" (list of role
// names), "passwordSalt". Example:
//
//	map[string]any{
//		"user-id":   "jsmith",
//		"full-name": "Jane Smith",
//		"password":  "secret",
//		"email":     "jsmith@example.com",
//	}
//
// When hashPassword is true OpenNMS hashes the plain-text password.
//
// The body is sent as XML — POST /rest/users does not accept JSON on
// any OpenNMS version. The XML document is built internally; callers
// still pass a plain map.
func (c *Client) CreateUser(ctx context.Context, user map[string]any, hashPassword bool) (map[string]any, error) {
	var params url.Values
	if hashPassword {
		params = url.Values{"hashPassword": {"true"}}
	}
	return asObject(c.postText(ctx, "users", usersBuildXML(user),
		"application/xml", false, "", params))
}

// UpdateUser updates user properties. Pass only the fields to
// change.
//
// The body is sent form-encoded — PUT /rest/users/{name} does not
// accept JSON on any OpenNMS version. Keys are bean property names,
// e.g. "fullName", "email", "password".
func (c *Client) UpdateUser(ctx context.Context, username string, user map[string]string) error {
	form := url.Values{}
	for k, v := range user {
		form.Set(k, v)
	}
	_, err := c.putForm(ctx, "users/"+url.PathEscape(username), form, nil, false)
	return err
}

// DeleteUser deletes a user.
func (c *Client) DeleteUser(ctx context.Context, username string) error {
	_, err := c.del(ctx, "users/"+url.PathEscape(username), nil, nil, false, "")
	return err
}

// AssignRoleToUser assigns roleName to username.
func (c *Client) AssignRoleToUser(ctx context.Context, username, roleName string) error {
	_, err := c.put(ctx, fmt.Sprintf("users/%s/roles/%s",
		url.PathEscape(username), url.PathEscape(roleName)), nil, nil, false)
	return err
}

// RevokeRoleFromUser revokes roleName from username.
func (c *Client) RevokeRoleFromUser(ctx context.Context, username, roleName string) error {
	_, err := c.del(ctx, fmt.Sprintf("users/%s/roles/%s",
		url.PathEscape(username), url.PathEscape(roleName)), nil, nil, false, "")
	return err
}
