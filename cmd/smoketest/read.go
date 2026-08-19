package main

import (
	"context"

	opennms "github.com/cnewkirk/go-opennms"
)

// runReadChecks exercises every read-only endpoint group.
// STUB — replaced by the read-checks port of smoke_test.py.
func runReadChecks(ctx context.Context, c *opennms.Client, r *runner) {
	_ = ctx
	_ = c
	_ = r
}
