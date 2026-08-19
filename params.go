package opennms

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ListOptions carries the standard pagination and ordering query
// parameters shared by the list endpoints.
//
// A nil *ListOptions uses the server-friendly default of Limit 10 /
// Offset 0 (matching the Python wrapper). An explicit &ListOptions{}
// sends limit=0, which OpenNMS interprets as "no limit".
type ListOptions struct {
	// Limit is the maximum number of results to return. 0 means all.
	Limit int
	// Offset is the zero-based offset for pagination.
	Offset int
	// OrderBy is the field name to sort by.
	OrderBy string
	// Order is the sort direction: "asc" or "desc".
	Order string
}

// listParams converts opts into query parameters, defaulting a nil
// opts to Limit 10.
func listParams(opts *ListOptions) url.Values {
	if opts == nil {
		opts = &ListOptions{Limit: 10}
	}
	p := url.Values{}
	p.Set("limit", strconv.Itoa(opts.Limit))
	p.Set("offset", strconv.Itoa(opts.Offset))
	if opts.OrderBy != "" {
		p.Set("orderBy", opts.OrderBy)
	}
	if opts.Order != "" {
		p.Set("order", opts.Order)
	}
	return p
}

// mergeFilters adds free-form filter parameters (Hibernate query
// filters like severity=MAJOR, or a comparator override) into p.
func mergeFilters(p url.Values, filters map[string]string) url.Values {
	if p == nil {
		p = url.Values{}
	}
	for k, v := range filters {
		p.Set(k, v)
	}
	return p
}

// asObject converts a parsed response into a JSON object.
func asObject(v any, err error) (map[string]any, error) {
	if err != nil || v == nil {
		return nil, err
	}
	obj, ok := v.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("opennms: expected JSON object, got %T", v)
	}
	return obj, nil
}

// asList converts a parsed response into a JSON array.
func asList(v any, err error) ([]any, error) {
	if err != nil || v == nil {
		return nil, err
	}
	list, ok := v.([]any)
	if !ok {
		return nil, fmt.Errorf("opennms: expected JSON array, got %T", v)
	}
	return list, nil
}

// asInt converts a parsed response (plain-text count or JSON number)
// into an int.
func asInt(v any, err error) (int, error) {
	if err != nil {
		return 0, err
	}
	switch n := v.(type) {
	case int:
		return n, nil
	case float64:
		return int(n), nil
	case string:
		return strconv.Atoi(n)
	}
	return 0, fmt.Errorf("opennms: expected count, got %T", v)
}

// asString converts a parsed response into a string.
func asString(v any, err error) (string, error) {
	if err != nil || v == nil {
		return "", err
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("opennms: expected text, got %T", v)
	}
	return s, nil
}

// getObject is the common GET-and-decode-object shorthand.
func (c *Client) getObject(ctx context.Context, path string, params url.Values, v2 bool) (map[string]any, error) {
	return asObject(c.get(ctx, path, params, v2))
}

// getList is the common GET-and-decode-array shorthand.
func (c *Client) getList(ctx context.Context, path string, params url.Values, v2 bool) ([]any, error) {
	return asList(c.get(ctx, path, params, v2))
}

// getCount is the common GET-a-plain-text-count shorthand.
func (c *Client) getCount(ctx context.Context, path string, v2 bool) (int, error) {
	return asInt(c.get(ctx, path, nil, v2))
}
