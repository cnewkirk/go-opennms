package opennms

// Search REST API v2 – /api/v2/search.

import (
	"context"
	"net/url"
	"strconv"
)

// Search performs a global search across all search contexts. Results
// are grouped by context (e.g. "Node", "Action") and filtered by the
// requesting user's permissions.
//
// searchContext optionally restricts the search to one context; ""
// searches all contexts. limit is the maximum results per context
// (default 10 on the server; negative for unbounded; 0 = server
// default).
//
// Returns the list of search result groups, or nil when nothing
// matched (204 No Content).
func (c *Client) Search(ctx context.Context, query, searchContext string, limit int) ([]any, error) {
	params := url.Values{"_s": {query}}
	if searchContext != "" {
		params.Set("_c", searchContext)
	}
	if limit != 0 {
		params.Set("_l", strconv.Itoa(limit))
	}
	return c.getList(ctx, "search", params, true)
}
