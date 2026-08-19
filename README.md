# go-opennms

[![CI](https://github.com/cnewkirk/go-opennms/actions/workflows/ci.yml/badge.svg)](https://github.com/cnewkirk/go-opennms/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/cnewkirk/go-opennms.svg)](https://pkg.go.dev/github.com/cnewkirk/go-opennms)
[![Go Report Card](https://goreportcard.com/badge/github.com/cnewkirk/go-opennms)](https://goreportcard.com/report/github.com/cnewkirk/go-opennms)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

An unofficial, dependency-free Go client for the
[OpenNMS](https://www.opennms.com/) REST API (Horizon 30+ and
Meridian). A faithful port of
[python-opennms](https://github.com/cnewkirk/python-opennms), whose
request shapes are live-validated against the Meridian 2025
foundation.

> **OpenNMS resources**: [Docs](https://docs.opennms.com/) ·
> [REST API reference](https://docs.opennms.com/horizon/latest/development/rest/rest-api.html) ·
> [Community forum](https://opennms.discourse.group/)

## Features

- Covers every v1 (`/opennms/rest/`) and v2 (`/opennms/api/v2/`)
  endpoint — one `Client` method per endpoint
- Plain maps in, decoded JSON out — the few endpoints that require
  XML or form-encoded bodies are handled internally
- Standard library only — zero dependencies
- Synchronous and straightforward; every method takes a
  `context.Context`
- Typed errors — `errors.Is(err, opennms.ErrNotFound)`,
  `errors.As(err, &apiErr)` for status code and body
- Built-in retries (connection errors and 500/502/503/504, exponential
  backoff) and sensible timeouts, both configurable
- `Paginate` iterator — range over every item of any list endpoint
- Full test suite with per-method coverage (mocked HTTP — no live
  server required)

## Installation

```bash
go get github.com/cnewkirk/go-opennms
```

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"log"

	opennms "github.com/cnewkirk/go-opennms"
)

func main() {
	ctx := context.Background()
	client := opennms.NewClient(
		"https://opennms.example.com:8443", "admin", "admin",
	)

	// Server info
	info, err := client.GetInfo(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("OpenNMS", info["displayVersion"])

	// All critical alarms
	alarms, err := client.GetAlarms(ctx,
		&opennms.ListOptions{Limit: 0},
		map[string]string{"severity": "CRITICAL"})
	if err != nil {
		log.Fatal(err)
	}
	for _, a := range alarms["alarm"].([]any) {
		alarm := a.(map[string]any)
		fmt.Println(alarm["id"], alarm["logMsg"])
	}

	// Acknowledge an alarm
	if err := client.AckAlarm(ctx, 42, ""); err != nil {
		log.Fatal(err)
	}
}
```

### Options

```go
client := opennms.NewClient(url, user, pass,
	opennms.WithTimeout(60*time.Second), // read timeout (default 30s)
	opennms.WithRetries(0),              // disable retries (default 3)
	opennms.WithInsecureSkipVerify(),    // self-signed certs (dev only)
)
```

### Errors

```go
_, err := client.GetNode(ctx, 99999)
if errors.Is(err, opennms.ErrNotFound) {
	// 404 — node does not exist
}
var apiErr *opennms.APIError
if errors.As(err, &apiErr) {
	fmt.Println(apiErr.StatusCode, apiErr.Body)
}
```

Sentinels: `ErrBadRequest` (400), `ErrAuthentication` (401),
`ErrForbidden` (403), `ErrNotFound` (404), `ErrConflict` (409),
`ErrServer` (5xx).

### Pagination

```go
fetch := func(limit, offset int) (map[string]any, error) {
	return client.GetAlarms(ctx,
		&opennms.ListOptions{Limit: limit, Offset: offset},
		map[string]string{"severity": "CRITICAL"})
}
for item, err := range opennms.Paginate(fetch, "alarm", 100) {
	if err != nil {
		log.Fatal(err)
	}
	alarm := item.(map[string]any)
	fmt.Println(alarm["id"], alarm["nodeLabel"])
}
```

## Design notes

- Responses are decoded generic JSON (`map[string]any`, `[]any`;
  numbers are `float64`), mirroring the thin-wrapper philosophy of
  python-opennms rather than maintaining hundreds of typed structs.
  `/count` endpoints return `int`; 204 writes return only `error`.
- Request shapes (URLs, bodies, content types) are ported verbatim
  from python-opennms, which live-validates them against a real
  server — including the endpoints where OpenNMS forces non-JSON
  bodies (form-encoded acks and alarm/notification actions, XML-only
  group/user creation, GraphML passthrough).

## Compatibility

Same coverage baseline as python-opennms: the Meridian 2025 REST API
reference, working against Horizon 30+ and Meridian. See that repo's
COVERAGE.md for the endpoint-by-endpoint status.

## Development

```bash
go test ./...     # full suite
gofmt -l .        # formatting (must print nothing)
go vet ./...
```

See [AGENTS.md](AGENTS.md) for the porting conventions.

## License

MIT — see [LICENSE](LICENSE).
