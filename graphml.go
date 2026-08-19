package opennms

// GraphML REST API – /rest/graphml.
//
// The GraphML endpoint is XML-only by design: GraphML is itself an
// XML format, so request and response bodies are raw GraphML
// documents rather than JSON.

import (
	"context"
	"net/url"
)

// GetGraphml returns a stored GraphML graph definition as an XML
// string.
func (c *Client) GetGraphml(ctx context.Context, graphName string) (string, error) {
	return c.getText(ctx, "graphml/"+url.PathEscape(graphName), nil,
		false, "application/xml")
}

// CreateGraphml creates a new GraphML graph stored under graphName
// from a GraphML document given as an XML string.
//
// The server returns 500 (not 400) when the graph already exists or
// the document fails validation — a valid document needs a "label"
// <key> and per-graph <data> entries (see the OpenNMS GraphML docs).
func (c *Client) CreateGraphml(ctx context.Context, graphName, graphmlXML string) error {
	_, err := c.postText(ctx, "graphml/"+url.PathEscape(graphName),
		graphmlXML, "application/xml", false, "*/*", nil)
	return err
}

// DeleteGraphml deletes a stored GraphML graph.
func (c *Client) DeleteGraphml(ctx context.Context, graphName string) error {
	_, err := c.del(ctx, "graphml/"+url.PathEscape(graphName), nil,
		nil, false, "*/*")
	return err
}
