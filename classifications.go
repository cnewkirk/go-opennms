package opennms

// Flow Classification REST API – /rest/classifications.

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// GetClassificationRules lists classification rules.
func (c *Client) GetClassificationRules(ctx context.Context, opts *ListOptions) (map[string]any, error) {
	return c.getObject(ctx, "classifications", listParams(opts), false)
}

// GetClassificationRule returns a specific classification rule by ID.
func (c *Client) GetClassificationRule(ctx context.Context, ruleID int) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("classifications/%d", ruleID), nil, false)
}

// CreateClassificationRule creates a new classification rule.
//
// rule is the rule definition. Common keys: "name", "dstAddress",
// "dstPort", "srcAddress", "srcPort", "protocol", "exporterFilter",
// "groupId".
func (c *Client) CreateClassificationRule(ctx context.Context, rule map[string]any) error {
	_, err := c.post(ctx, "classifications", rule, nil, false)
	return err
}

// UpdateClassificationRule updates the classification rule with the
// given database ID using the updated rule definition.
func (c *Client) UpdateClassificationRule(ctx context.Context, ruleID int, rule map[string]any) error {
	_, err := c.put(ctx, fmt.Sprintf("classifications/%d", ruleID), rule, nil, false)
	return err
}

// DeleteClassificationRules deletes classification rules, optionally
// filtered by group. When groupID is non-zero, only rules in that
// group are deleted.
func (c *Client) DeleteClassificationRules(ctx context.Context, groupID int) error {
	var params url.Values
	if groupID != 0 {
		params = url.Values{"groupId": {strconv.Itoa(groupID)}}
	}
	_, err := c.del(ctx, "classifications", params, nil, false, "")
	return err
}

// DeleteClassificationRule deletes a specific classification rule.
func (c *Client) DeleteClassificationRule(ctx context.Context, ruleID int) error {
	_, err := c.del(ctx, fmt.Sprintf("classifications/%d", ruleID), nil, nil, false, "")
	return err
}

// Classify classifies a flow record against the configured rules.
//
// request is the classification request with keys such as
// "srcAddress", "srcPort", "dstAddress", "dstPort", "protocol",
// "exporterAddress".
func (c *Client) Classify(ctx context.Context, request map[string]any) (map[string]any, error) {
	return asObject(c.post(ctx, "classifications/classify", request, nil, false))
}

// GetClassificationGroups lists classification groups.
func (c *Client) GetClassificationGroups(ctx context.Context, opts *ListOptions) (map[string]any, error) {
	return c.getObject(ctx, "classifications/groups", listParams(opts), false)
}

// GetClassificationGroup returns a specific classification group by ID.
func (c *Client) GetClassificationGroup(ctx context.Context, groupID int) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("classifications/groups/%d", groupID), nil, false)
}

// CreateClassificationGroup creates a new classification group.
// group is the group definition with key "name".
func (c *Client) CreateClassificationGroup(ctx context.Context, group map[string]any) error {
	_, err := c.post(ctx, "classifications/groups", group, nil, false)
	return err
}

// UpdateClassificationGroup updates the classification group with the
// given database ID using the updated group definition.
func (c *Client) UpdateClassificationGroup(ctx context.Context, groupID int, group map[string]any) error {
	_, err := c.put(ctx, fmt.Sprintf("classifications/groups/%d", groupID), group, nil, false)
	return err
}

// DeleteClassificationGroup deletes a classification group.
func (c *Client) DeleteClassificationGroup(ctx context.Context, groupID int) error {
	_, err := c.del(ctx, fmt.Sprintf("classifications/groups/%d", groupID), nil, nil, false, "")
	return err
}

// ImportClassificationRules imports classification rules from CSV
// text into the group with the given database ID.
func (c *Client) ImportClassificationRules(ctx context.Context, groupID int, csvText string) error {
	_, err := c.postText(ctx, fmt.Sprintf("classifications/groups/%d", groupID),
		csvText, "text/comma-separated-values", false, "", nil)
	return err
}

// GetClassificationProtocols lists all known protocols used by the
// classification engine.
func (c *Client) GetClassificationProtocols(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "classifications/protocols", nil, false)
}
