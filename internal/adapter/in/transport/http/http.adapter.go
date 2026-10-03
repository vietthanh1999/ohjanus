package http

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/vietthanh1999/ohjanus/internal/adapter/in/mcp"
	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// maxBody caps single POST bodies (SQL is capped separately by max_query_length).
const maxBody = 1 << 20

// Server exposes the MCP protocol over HTTP (§4.1.2 of the spec):
//
//	POST /mcp      JSON-RPC request (single or batch), Bearer auth
//	GET  /mcp/sse  SSE stream for server-initiated messages
//	GET  /healthz  public liveness probe
type Server struct {
	mcpServer *mcp.Server
	tokens    out.TokenStore
	authMode  string
	tls       TLSConfig
	log       *slog.Logger
	mux       *http.ServeMux
	srv       *http.Server
}

// TLSConfig mirrors server.http.tls.
type TLSConfig struct {
	Enabled  bool
	CertFile string
	KeyFile  string
}

// New wires an HTTP transport.
func New(mcpServer *mcp.Server, tokens out.TokenStore, authMode string, tls TLSConfig, log *slog.Logger) *Server {
	s := &Server{mcpServer: mcpServer, tokens: tokens, authMode: authMode, tls: tls, log: log}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/mcp", s.handleMCP)
	mux.HandleFunc("/mcp/sse", s.handleSSE)
	s.mux = mux
	s.srv = &http.Server{Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	return s
}

// ServeListener serves on a pre-bound listener until ctx ends.
func (s *Server) ServeListener(ctx context.Context, ln net.Listener) error {
	go func() {
		<-ctx.Done()
		shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = s.srv.Shutdown(shutCtx)
	}()
	var err error
	if s.tls.Enabled {
		tlsLn := tlsListener(ln, s.tls)
		err = s.srv.Serve(tlsLn)
	} else {
		err = s.srv.Serve(ln)
	}
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

// tlsListener wraps ln with the configured certificate.
func tlsListener(ln net.Listener, cfg TLSConfig) net.Listener {
	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		// Fail the next Accept with a clear error instead of panicking here.
		return &errorListener{err: err}
	}
	return tls.NewListener(ln, &tls.Config{Certificates: []tls.Certificate{cert}})
}

// errorListener fails Accept so Serve returns the TLS setup error.
type errorListener struct{ err error }

func (l *errorListener) Accept() (net.Conn, error) { return nil, l.err }
func (l *errorListener) Close() error              { return nil }
func (l *errorListener) Addr() net.Addr            { return &net.TCPAddr{} }

// handleMCP processes one POSTed JSON-RPC message or batch.
func (s *Server) handleMCP(w http.ResponseWriter, r *http.Request) {
	setCORS(w, r)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		writeRPCError(w, http.StatusMethodNotAllowed, mcp.CodeMethodNotFound, "use POST /mcp")
		return
	}
	token, ok := bearerToken(r)
	if !ok {
		writeRPCError(w, http.StatusUnauthorized, mcp.CodeJanusError, "UNAUTHENTICATED: missing bearer token")
		return
	}
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxBody))
	if err != nil {
		writeRPCError(w, http.StatusBadRequest, mcp.CodeParseError, "cannot read body")
		return
	}
	trimmed := json.RawMessage(bytesTrimSpace(body))
	var raws []json.RawMessage
	if isJSONArray(trimmed) {
		if err := json.Unmarshal(trimmed, &raws); err != nil {
			writeRPCError(w, http.StatusBadRequest, mcp.CodeParseError, "parse error")
			return
		}
	} else {
		raws = []json.RawMessage{trimmed}
	}
	ctx := r.Context()
	responses := make([]json.RawMessage, 0, len(raws))
	for _, raw := range raws {
		resp := s.mcpServer.Handle(ctx, injectToken(raw, token))
		if resp != nil {
			responses = append(responses, resp)
		}
	}
	if len(responses) == 0 {
		w.WriteHeader(http.StatusAccepted)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if len(responses) == 1 && !isJSONArray(trimmed) {
		w.WriteHeader(authStatus(responses[0]))
		_, _ = w.Write(responses[0])
		return
	}
	status := http.StatusOK
	if len(responses) > 0 {
		allAuth := true
		for _, resp := range responses {
			if authStatus(resp) == http.StatusOK {
				allAuth = false
				break
			}
		}
		if allAuth {
			status = authStatus(responses[0])
		}
	}
	out, _ := json.Marshal(responses)
	w.WriteHeader(status)
	_, _ = w.Write(out)
}

// authStatus maps Janus auth errors to HTTP status (§4.1.2: 401 missing or
// invalid token, 403 expired scope or forbidden). Non-auth payloads are 200.
func authStatus(resp json.RawMessage) int {
	var parsed struct {
		Error *struct {
			Data map[string]any `json:"data"`
		} `json:"error"`
	}
	if err := json.Unmarshal(resp, &parsed); err != nil || parsed.Error == nil {
		return http.StatusOK
	}
	code, _ := parsed.Error.Data["code"].(string)
	switch code {
	case "UNAUTHENTICATED", "TOKEN_INVALID", "TOKEN_EXPIRED":
		return http.StatusUnauthorized
	case "FORBIDDEN":
		return http.StatusForbidden
	default:
		return http.StatusOK
	}
}

// injectToken sets params.token so the MCP server authenticates the Bearer
// credential through its single token path (params.token / auth mode none).
func injectToken(raw json.RawMessage, token string) []byte {
	var msg map[string]any
	if err := json.Unmarshal(raw, &msg); err != nil {
		return raw
	}
	params, _ := msg["params"].(map[string]any)
	if params == nil {
		params = map[string]any{}
		msg["params"] = params
	}
	params["token"] = token
	out, err := json.Marshal(msg)
	if err != nil {
		return raw
	}
	return out
}

// handleSSE streams server-initiated approval notifications.
func (s *Server) handleSSE(w http.ResponseWriter, r *http.Request) {
	setCORS(w, r)
	if r.Method != http.MethodGet {
		writeRPCError(w, http.StatusMethodNotAllowed, mcp.CodeMethodNotFound, "use GET /mcp/sse")
		return
	}
	if _, ok := s.requireBearer(r); !ok {
		writeRPCError(w, http.StatusUnauthorized, mcp.CodeJanusError, "UNAUTHENTICATED: missing bearer token")
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeRPCError(w, http.StatusInternalServerError, mcp.CodeInternal, "streaming unsupported")
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	flusher.Flush()

	var ch <-chan domain.TokenEvent
	var cancel func()
	if s.tokens != nil {
		ch, cancel = s.tokens.Watch()
		defer cancel()
	}
	heartbeat := time.NewTicker(15 * time.Second)
	defer heartbeat.Stop()
	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-ch:
			writeSSEEvent(w, "message", mcpApprovalNotification(ev))
			flusher.Flush()
		case <-heartbeat.C:
			_, _ = fmt.Fprintln(w, ":ping")
			flusher.Flush()
		}
	}
}

// mcpApprovalNotification wraps a token event as an MCP notification.
// Method "notifications/janus-approval" is a Janus extension clients
// may ignore; the payload mirrors GET /api/v1/approvals/{id}.
func mcpApprovalNotification(ev domain.TokenEvent) map[string]any {
	t := ev.Token
	return map[string]any{
		"jsonrpc": "2.0",
		"method":  "notifications/janus-approval",
		"params": map[string]any{
			"type": string(ev.Type), "id": t.ID, "state": string(t.State),
			"connection": t.Connection, "statement_type": string(t.Statement),
		},
	}
}

func writeSSEEvent(w http.ResponseWriter, event string, v any) {
	raw, _ := json.Marshal(v)
	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, raw)
}

func writeRPCError(w http.ResponseWriter, status, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(mcp.Response{
		JSONRPC: "2.0",
		Error:   &mcp.ErrorObject{Code: code, Message: msg},
	})
}

func setCORS(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Vary", "Origin")
	}
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
}

func bearerToken(r *http.Request) (string, bool) {
	h := r.Header.Get("Authorization")
	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") || strings.TrimSpace(parts[1]) == "" {
		return "", false
	}
	return strings.TrimSpace(parts[1]), true
}

func (s *Server) requireBearer(r *http.Request) (string, bool) {
	if s.authMode == "none" {
		return "", true
	}
	return bearerToken(r)
}

func isJSONArray(raw json.RawMessage) bool {
	for _, b := range []byte(raw) {
		if b == ' ' || b == '\n' || b == '\r' || b == '\t' {
			continue
		}
		return b == '['
	}
	return false
}

func bytesTrimSpace(b []byte) []byte {
	return []byte(strings.TrimSpace(string(b)))
}
