// Command smoketest exercises the go-opennms client against a live
// OpenNMS server. It is a port of python-opennms's smoke_test.py.
//
// Read-only mode (default) is safe to run against any server,
// including production. It issues only GET requests.
//
// Write mode (--write) creates and deletes objects on the server and
// prompts for confirmation. Only use write mode against a dev or
// staging instance — never production.
//
// Environment variables:
//
//	OPENNMS_URL         Base URL, e.g. "https://onms.example.com:8443" (required)
//	OPENNMS_USER        Username with at minimum the "rest" role (required)
//	OPENNMS_PASSWORD    Password (required)
//	OPENNMS_VERIFY_SSL  Set to "false" to disable SSL verification (default: true)
//	OPENNMS_TIMEOUT     Request timeout in seconds (default: 60)
//
// Usage:
//
//	go run ./cmd/smoketest                 # read-only — safe anywhere
//	go run ./cmd/smoketest --write         # write ops — prompts
//	go run ./cmd/smoketest --write --yes   # write ops — no prompt (CI)
//	go run ./cmd/smoketest --no-color
//	go run ./cmd/smoketest --skip get_resources,get_flow
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	opennms "github.com/cnewkirk/go-opennms"
)

// failedSentinel distinguishes "call failed" from "call succeeded and
// returned nil (204)"; runner methods return ok=false alongside it.
type runner struct {
	noColor      bool
	skipPrefixes []string
	passed       int
	failed       int
	skipped      int
	warned       int
	failures     [][2]string
	warnings     [][2]string
}

func (r *runner) color(code, s string) string {
	if r.noColor {
		return s
	}
	return "\033[" + code + "m" + s + "\033[0m"
}

func (r *runner) section(title string) {
	fmt.Printf("\n[%s]\n", title)
}

func (r *runner) shouldSkip(label string) bool {
	for _, p := range r.skipPrefixes {
		if strings.HasPrefix(label, p) {
			return true
		}
	}
	return false
}

func (r *runner) skip(label, reason string) {
	r.skipped++
	suffix := ""
	if reason != "" {
		suffix = "  (" + reason + ")"
	}
	fmt.Printf("  %s  %s%s\n", r.color("33", "SKIP"), label, suffix)
}

func (r *runner) ok(label, detail string) {
	r.passed++
	suffix := ""
	if detail != "" {
		suffix = "  " + r.color("2", detail)
	}
	fmt.Printf("  %s  %s%s\n", r.color("32", "PASS"), label, suffix)
}

func (r *runner) fail(label string, err error) {
	r.failed++
	r.failures = append(r.failures, [2]string{label, err.Error()})
	fmt.Printf("  %s  %s  %s\n",
		r.color("31", "FAIL"), label, r.color("2", err.Error()))
}

func (r *runner) warnMsg(label string, err error, note string) {
	r.warned++
	msg := err.Error()
	if note != "" {
		msg += "  (" + note + ")"
	}
	r.warnings = append(r.warnings, [2]string{label, msg})
	fmt.Printf("  %s  %s  %s\n",
		r.color("33", "WARN"), label, r.color("2", msg))
}

// run calls fn and records PASS or FAIL. detail, if non-nil, renders
// the detail column from a non-nil result. ok is false when the call
// failed or was skipped.
func (r *runner) run(label string, detail func(any) string, fn func() (any, error)) (any, bool) {
	if r.shouldSkip(label) {
		r.skip(label, "--skip")
		return nil, false
	}
	result, err := fn()
	if err != nil {
		r.fail(label, err)
		return nil, false
	}
	d := ""
	if detail != nil && result != nil {
		d = detail(result)
	}
	r.ok(label, d)
	return result, true
}

// warn is like run but records WARN instead of FAIL on error — for
// endpoints that depend on optional plugins or heavy server-side
// queries. note is appended to warning output.
func (r *runner) warn(label, note string, detail func(any) string, fn func() (any, error)) (any, bool) {
	if r.shouldSkip(label) {
		r.skip(label, "--skip")
		return nil, false
	}
	result, err := fn()
	if err != nil {
		r.warnMsg(label, err, note)
		return nil, false
	}
	d := ""
	if detail != nil && result != nil {
		d = detail(result)
	}
	r.ok(label, d)
	return result, true
}

// n returns a detail func rendering a short item-count summary, using
// listKey for wrapped list responses ("" for none).
func n(listKey string) func(any) string {
	return func(result any) string {
		switch v := result.(type) {
		case int:
			return strconv.Itoa(v)
		case []any:
			return fmt.Sprintf("%d items", len(v))
		case map[string]any:
			if listKey != "" {
				if items, found := v[listKey]; found {
					if list, isList := items.([]any); isList {
						return fmt.Sprintf("%d items", len(list))
					}
					return fmt.Sprintf("%v items", items)
				}
			}
			return "ok"
		}
		return "ok"
	}
}

// first fetches one item from a collection endpoint and returns the
// item and its idKey value ("" idKey → "id"). Returns (nil, nil) on
// error or empty result.
func first(fn func(limit int) (any, error), listKey, idKey string) (map[string]any, any) {
	if idKey == "" {
		idKey = "id"
	}
	result, err := fn(1)
	if err != nil {
		return nil, nil
	}
	var items []any
	switch v := result.(type) {
	case map[string]any:
		items, _ = v[listKey].([]any)
	case []any:
		items = v
	}
	if len(items) == 0 {
		return nil, nil
	}
	item, _ := items[0].(map[string]any)
	if item == nil {
		return nil, nil
	}
	return item, item[idKey]
}

// toInt converts a decoded JSON number (float64) or int to int.
func toInt(v any) int {
	switch x := v.(type) {
	case float64:
		return int(x)
	case int:
		return x
	}
	return 0
}

func main() {
	write := flag.Bool("write", false,
		"also exercise write operations (dev/staging only)")
	yes := flag.Bool("yes", false,
		"skip the write-mode confirmation prompt (for CI)")
	noColor := flag.Bool("no-color", false, "disable ANSI colour output")
	skipList := flag.String("skip", "",
		"comma-separated test label prefixes to skip")
	flag.Parse()

	r := &runner{noColor: *noColor}
	for _, s := range strings.Split(*skipList, ",") {
		if s = strings.TrimSpace(s); s != "" {
			r.skipPrefixes = append(r.skipPrefixes, s)
		}
	}

	url := os.Getenv("OPENNMS_URL")
	user := os.Getenv("OPENNMS_USER")
	password := os.Getenv("OPENNMS_PASSWORD")
	verify := !strings.EqualFold(os.Getenv("OPENNMS_VERIFY_SSL"), "false")
	timeout := 60
	if t := os.Getenv("OPENNMS_TIMEOUT"); t != "" {
		var err error
		if timeout, err = strconv.Atoi(t); err != nil {
			fmt.Fprintln(os.Stderr, "Error: OPENNMS_TIMEOUT must be an integer")
			os.Exit(2)
		}
	}
	var missing []string
	for _, pair := range [][2]string{
		{"OPENNMS_URL", url}, {"OPENNMS_USER", user},
		{"OPENNMS_PASSWORD", password},
	} {
		if pair[1] == "" {
			missing = append(missing, pair[0])
		}
	}
	if len(missing) > 0 {
		fmt.Fprintf(os.Stderr,
			"Error: missing environment variable(s): %s\n",
			strings.Join(missing, ", "))
		os.Exit(2)
	}

	if *write && !*yes {
		fmt.Println("\n  WARNING: Write mode will create and delete objects on the server.")
		fmt.Println("  All created objects are deleted at the end of the run.")
		fmt.Println("  ONLY use against a dev or staging server -- NEVER production.")
		fmt.Printf("\n  Target URL: %s\n", url)
		fmt.Print("\n  Type 'yes' to continue: ")
		confirm, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		if !strings.EqualFold(strings.TrimSpace(confirm), "yes") {
			fmt.Fprintln(os.Stderr, "Aborted.")
			os.Exit(2)
		}
		fmt.Println()
	}

	opts := []opennms.Option{
		opennms.WithTimeout(time.Duration(timeout) * time.Second),
	}
	if !verify {
		opts = append(opts, opennms.WithInsecureSkipVerify())
	}
	client := opennms.NewClient(url, user, password, opts...)
	ctx := context.Background()

	fmt.Println("OpenNMS Smoke Test (go-opennms)")
	fmt.Printf("  URL:  %s\n", url)
	fmt.Printf("  User: %s\n", user)
	mode := "read-only (getters)"
	if *write {
		mode = "read + write"
	}
	fmt.Printf("  Mode: %s\n", mode)
	if !verify {
		fmt.Println("  SSL:  verification disabled")
	}

	runReadChecks(ctx, client, r)
	if *write {
		runWriteChecks(ctx, client, r)
	}

	total := r.passed + r.failed + r.warned + r.skipped
	fmt.Printf("\n%s\n", strings.Repeat("─", 56))
	parts := []string{
		fmt.Sprintf("%d passed", r.passed),
		fmt.Sprintf("%d failed", r.failed),
	}
	if r.warned > 0 {
		parts = append(parts, fmt.Sprintf("%d warned", r.warned))
	}
	parts = append(parts, fmt.Sprintf("%d skipped", r.skipped))
	fmt.Printf("  %s  (%d total)\n", strings.Join(parts, "  ·  "), total)

	if len(r.warnings) > 0 {
		fmt.Println("\nWarnings (non-fatal):")
		for _, w := range r.warnings {
			fmt.Printf("  %s\n    %s\n", w[0], w[1])
		}
	}
	if len(r.failures) > 0 {
		fmt.Println("\nFailures:")
		for _, f := range r.failures {
			fmt.Printf("  %s\n    %s\n", f[0], f[1])
		}
	}
	if r.failed > 0 {
		os.Exit(1)
	}
}
