package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// Engine prompts a human on the terminal for write approval.
// It prints to err (never stdout) and reads the decision from in.
type Engine struct {
	in  io.Reader
	err io.Writer
}

var _ out.ApprovalEngine = (*Engine)(nil)

// NewApprovalEngine builds a CLI approval engine reading from stdin.
func NewApprovalEngine() *Engine {
	return &Engine{in: os.Stdin, err: os.Stderr}
}

// Approve prints the write request and waits for y/n.
func (e *Engine) Approve(_ context.Context, req out.ApprovalRequest) (bool, error) {
	fmt.Fprintf(e.err, `
+-------------------------------+
| JANUS WRITE REQUEST           |
| Connection: %-17s |
| Statement:  %-17s |
| SQL: %.60s
| Params: %v
| Affected estimate: ~%d rows
+-------------------------------+
Approve? [y/N] `,
		req.Connection, req.Statement, singleLine(req.SQL), req.Params, req.Estimate)
	line, err := bufio.NewReader(e.in).ReadString('\n')
	if err != nil && len(line) == 0 {
		return false, nil
	}
	switch strings.ToLower(strings.TrimSpace(line)) {
	case "y", "yes":
		return true, nil
	default:
		return false, nil
	}
}

func singleLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
