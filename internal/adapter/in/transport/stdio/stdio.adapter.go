package stdio

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"log/slog"

	"github.com/vietthanh1999/ohjanus/internal/adapter/in/mcp"
)

// Transport reads JSON-RPC messages from stdin and writes responses to stdout.
// Logs go to stderr and never mix into stdout (§4.1.1 of the spec).
type Transport struct {
	server *mcp.Server
	in     io.Reader
	out    io.Writer
	log    *slog.Logger
}

// New builds a stdio transport.
func New(server *mcp.Server, in io.Reader, out io.Writer, log *slog.Logger) *Transport {
	return &Transport{server: server, in: in, out: out, log: log}
}

// Serve loops over stdin lines until EOF or context cancellation.
func (t *Transport) Serve(ctx context.Context) error {
	sc := bufio.NewScanner(t.in)
	sc.Buffer(make([]byte, 1024*1024), 10*1024*1024)
	for sc.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		line := sc.Bytes()
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		resp := t.server.Handle(ctx, line)
		if resp == nil {
			continue
		}
		if _, err := t.out.Write(append(resp, '\n')); err != nil {
			return err
		}
	}
	return sc.Err()
}
