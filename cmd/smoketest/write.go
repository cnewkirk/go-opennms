package main

import (
	"context"

	opennms "github.com/cnewkirk/go-opennms"
)

// runWriteChecks exercises write endpoints (create/update/delete),
// cleaning up everything it creates.
// STUB — replaced by the write-checks port of smoke_test.py.
func runWriteChecks(ctx context.Context, c *opennms.Client, r *runner) {
	_ = ctx
	_ = c
	_ = r
}
