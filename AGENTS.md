# AGENTS.md

This file provides guidance to AI coding agents working with code in
this repository.

## Project purpose

A thin, synchronous Go client for the OpenNMS REST API (Horizon 30+
and Meridian). This is a faithful port of the Python
[opennms-api-wrapper](https://github.com/cnewkirk/opennms-api-wrapper)
package: one `Client` type with one method per API endpoint. The
Python package is the reference implementation — its request shapes
are live-validated against a real OpenNMS server, so **the Python
source and tests are the spec**; do not invent different URLs,
bodies, or content types.

## Development environment

```bash
go build ./...    # compile
go test ./...     # run full suite
gofmt -l .        # must print nothing
go vet ./...      # must pass
```

## Repository layout

All code lives in a single `opennms` package at the repo root:

- `client.go` — `Client`, `NewClient`, functional options
- `request.go` — HTTP plumbing: `do`, `parse`, verb helpers
  (`get/post/postForm/put/putForm/del/patch/getText/getBytes/
  postBytes/postFiles/putText/postText`)
- `params.go` — `ListOptions`, `listParams`, `mergeFilters`, and the
  result converters `asObject/asList/asInt/asString` plus shorthands
  `getObject/getList/getCount`
- `errors.go` — `APIError` + sentinel errors (`ErrNotFound`, …)
- `pagination.go` — `Paginate` iterator
- `<resource>.go` — one file per API resource group (`alarms.go`,
  `nodes.go`, `alarm_stats.go`, …), mirroring the Python
  `opennms_api_wrapper/_<resource>.py` mixins
- `<resource>_test.go` — tests, mirroring `tests/test_<resource>.py`
- `testutil_test.go` — shared test helpers (`newTestClient`,
  `handleJSON/handleText/handleStatus/handleContract`, `capture`)

## Porting conventions (Python → Go)

`alarms.go` and `alarms_test.go` are the canonical template — match
their style exactly.

**Method names**: `get_alarm_count` → `GetAlarmCount`. Initialisms
follow Go style: `Snmp` → `SNMP`? **No** — keep the Python
capitalization pattern simple and mechanical: `get_snmp_config` →
`GetSnmpConfig`, `get_ksc_reports` → `GetKscReports`. Every method
takes `ctx context.Context` first.

**Signatures**:
- `limit/offset/order_by/order` keyword args → `opts *ListOptions`
  (nil = Python defaults: limit 10, offset 0).
- `**filters` → `filters map[string]string` (nil allowed).
- `Optional[str] = None` → plain `string`, `""` = absent.
- `Optional[int] = None` where 0 is not a meaningful value → plain
  `int` with 0 = absent; where 0 is meaningful → `*int`.
- `bool = False` flags → plain `bool`.
- dict payloads → `body map[string]any`.

**Return values** (mirror what the Python method actually returns per
its tests):
- JSON object → `(map[string]any, error)` via `c.getObject`/`asObject(...)`
- JSON array → `([]any, error)` via `c.getList`/`asList(...)`
- plain-text count → `(int, error)` via `c.getCount`/`asInt(...)`
- raw text (XML/GraphML) → `(string, error)` via `c.getText`/`asString(...)`
- bytes (images, rendered reports) → `([]byte, error)`
- writes whose Python return is `None` → `error` only
  (`_, err := c.put(...); return err`)

**Paths**: build with `fmt.Sprintf("alarms/%d", alarmID)`. Path
segments that may contain spaces or special characters (requisition
names, resource IDs, …) must be escaped with `url.PathEscape(...)`
exactly where the Python code relies on `requests` percent-encoding.

**v2 endpoints**: pass `v2: true` (the last bool arg of the verb
helpers) — never hardcode `/api/v2` into paths.

**Content types** (do not change — the server dispatches on them):
JSON everywhere except: `POST /rest/acks` (form), `POST /rest/groups`
and `POST /rest/users` (XML built internally from a `map[string]any`,
callers still pass plain maps), `PUT /rest/groups/{name}` and
`PUT /rest/users/{name}` (form), `/rest/graphml` (opaque XML strings
via `postText`/`getText`).

**Unexported helpers and test fixtures must be prefixed with the
resource name** (e.g. `bulkAlarmAction` is grandfathered, but write
`nodesBuildParams`, `nodeJSON`, `requisitionListJSON`) — everything
shares one package, so generic names collide.

**Doc comments**: every exported method gets a Go doc comment,
carrying over the substance of the Python docstring (including
Args/Example content, condensed to Go style).

## Test conventions

- `newTestClient(t)` returns an `*http.ServeMux` and a `Client`
  pointed at an `httptest.Server` (retries disabled).
- Register handlers with Go 1.22 method patterns:
  `"GET "+v1Path+"/alarms/{id}"`, using the `v1Path`/`v2Path`
  constants from `testutil_test.go`.
- **Request-shape assertions come from the Python tests, never from
  the Go implementation** — port each `tests/test_<resource>.py`
  case, keeping its assertions (URL params via `req.query`, bodies
  via `req.body`/`req.formBody(t)`/`json.Unmarshal`, headers via
  `req.header`).
- Where the Python test uses `add_contract` (XML-only creates,
  form-encoded updates), use `handleContract` with the same
  Content-Type — it answers 415 for anything else, exactly as the
  server does.
- Form-encoded bodies: `url.Values.Encode()` sorts keys, so compare
  parsed values (`req.formBody(t).Get("ack")`), not raw strings with
  more than one key.
- JSON numbers decode as `float64` — assert with
  `result["id"].(float64) != 42`.
- Fixture JSON lives as `const <resource>...JSON` strings in each
  test file, ported from the shapes in the Python `tests/fixtures.py`
  (they mirror real OpenNMS responses — keep the key shapes:
  singular resource name as list wrapper key, plus `totalCount`,
  `count`, `offset`).
- 204 responses: methods return nil error and no value; plain-text
  counts return `int`.

## Style

- `gofmt`-clean, `go vet`-clean. Standard library only — no external
  dependencies.
- No inline comments on code that is self-evident.
- No error handling for scenarios that cannot happen.
- No abstractions introduced for one-off operations.

## Git workflow

- All changes go through a branch + PR; branch naming
  `feature/<topic>`, `fix/<topic>`, `docs/<topic>`.
- Conventional-commit messages (`feat:`, `fix:`, `test:`, `docs:`).
- No AI attribution lines in commits or PRs.
