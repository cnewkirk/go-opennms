// Package types provides optional, documentation-oriented payload
// types for the map[string]any client methods in package opennms.
//
// It is a port of the python-opennms package's types.py TypedDict
// schemas: one struct per write-payload shape, with every field
// optional at the Go level (all JSON tags carry omitempty). Fields
// marked Required in the comments must be included or the server
// will return a 400.
//
// Plain map[string]any values are accepted by every client method —
// these types exist purely for documentation and editor support. The
// Map function bridges a typed value into the map[string]any the
// client methods expect:
//
//	client.CreateNode(ctx, types.Map(types.Node{Label: "web01"}))
//
// Boolean fields are *bool and a few numeric fields are *int so that
// false and meaningful zeros survive omitempty; see the field
// comments for details.
package types
