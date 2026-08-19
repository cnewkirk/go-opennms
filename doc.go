// Package opennms is a thin, synchronous Go client for the OpenNMS
// REST API (Horizon 30+ and Meridian).
//
// It is a port of the Python opennms-api-wrapper package
// (https://github.com/cnewkirk/opennms-api-wrapper) and mirrors its
// design: one Client type with one method per API endpoint, JSON-first
// request bodies, and responses returned as decoded generic values
// (map[string]any, []any, int, string) rather than typed structs.
//
// Usage:
//
//	client := opennms.NewClient(
//		"https://onms.example.com", "admin", "admin",
//	)
//
//	alarms, err := client.GetAlarms(ctx, nil, map[string]string{
//		"severity": "CRITICAL",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, a := range alarms["alarm"].([]any) {
//		alarm := a.(map[string]any)
//		fmt.Println(alarm["id"], alarm["logMsg"])
//	}
//
// Failed requests return an *APIError carrying the status code and
// response body; match error classes with errors.Is against the
// package sentinels (ErrNotFound, ErrAuthentication, ...).
package opennms
