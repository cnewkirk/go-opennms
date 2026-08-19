package types

import "encoding/json"

// Map converts a typed payload value into the map[string]any accepted
// by the client's write methods, by marshalling v to JSON and
// unmarshalling the result into a map. Fields left at their zero
// value are dropped by their omitempty tags, so only the fields you
// set appear in the map.
//
// Map returns nil for a nil input and panics only if v cannot be
// marshalled to a JSON object (e.g. an unsupported type such as a
// channel, or a value that is not a JSON object at the top level).
//
// Usage:
//
//	client.CreateNode(ctx, types.Map(types.Node{Label: "web01"}))
func Map(v any) map[string]any {
	if v == nil {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		panic(err)
	}
	return m
}
