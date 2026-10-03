package admin

import (
	"context"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// Server exposes the Admin API for the UI (§4 of ui.md) on a port separate
// from MCP. Interim auth: Bearer MCP token; approve/reject require the
// admin scope so agents cannot approve their own writes. Session/OIDC lands
// in Phase 3/4.
type Server struct {
	addr      string
	authMode  string
	tokens    out.TokenStore
	resolver  out.TokenResolver
	authStore out.AuthTokenStore
	audit     out.AuditReader
	auditSink out.AuditSink
	pools     map[string]out.Pool
	metas     map[string]domain.ConnectionMeta
	clock     out.Clock
	mux       *http.ServeMux
	srv       *http.Server
}

// New wires an Admin server. authStore may be nil when token management
// is unavailable (endpoints return 503).
func New(addr, authMode string, tokens out.TokenStore, resolver out.TokenResolver, authStore out.AuthTokenStore, audit out.AuditReader, sink out.AuditSink, pools map[string]out.Pool, metas map[string]domain.ConnectionMeta, clock out.Clock) *Server {
	s := &Server{
		addr: addr, authMode: authMode, tokens: tokens, resolver: resolver,
		authStore: authStore,
		audit:     audit, auditSink: sink, pools: pools, metas: metas, clock: clock,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/readyz", s.handleReady)
	mux.HandleFunc("/api/v1/approvals/stream", s.withAuthFunc(s.handleStream))
	mux.HandleFunc("/api/v1/approvals/", s.withAuthFunc(s.handleApprovalSub))
	mux.HandleFunc("/api/v1/approvals", s.withAuth(domain.ScopeRead, s.handleApprovals))
	mux.HandleFunc("/api/v1/audit/", s.withAuth(domain.ScopeRead, s.handleAuditOne))
	mux.HandleFunc("/api/v1/audit", s.withAuth(domain.ScopeRead, s.handleAudit))
	mux.HandleFunc("/api/v1/connections/", s.withAuthFunc(s.handleConnectionSub))
	mux.HandleFunc("/api/v1/connections", s.withAuth(domain.ScopeRead, s.handleConnections))
	mux.HandleFunc("/api/v1/dashboard/summary", s.withAuth(domain.ScopeRead, s.handleSummary))
	mux.HandleFunc("/api/v1/tokens", s.withAuth(domain.ScopeAdmin, s.handleTokens))
	mux.HandleFunc("/api/v1/tokens/", s.withAuthFunc(s.handleTokenOne))
	mux.HandleFunc("/api/v1/audit/export", s.withAuth(domain.ScopeRead, s.handleAuditExport))
	s.mux = mux
	s.srv = &http.Server{Addr: addr, Handler: s.Handler(), ReadHeaderTimeout: 5 * time.Second}
	return s
}

// Handler returns the mux wrapped in CORS handling for browser clients
// (Vite dev server). Auth still uses the Authorization header; the API
// never relies on cookies so a wildcard origin is safe here.
func (s *Server) Handler() http.Handler {
	return corsMiddleware(s.mux)
}

// corsMiddleware answers preflights and stamps permissive headers on
// every response so the Admin UI can call the API cross-origin in dev.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// SetMetricsHandler mounts a public Prometheus scrape endpoint.
// Call before serving; nil disables it.
func (s *Server) SetMetricsHandler(h http.Handler) {
	if h == nil {
		return
	}
	s.mux.Handle("/metrics", h)
}

// ServeListener serves on a pre-bound listener until ctx ends.
func (s *Server) ServeListener(ctx context.Context, ln net.Listener) error {
	go func() {
		<-ctx.Done()
		shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = s.srv.Shutdown(shutCtx)
	}()
	if err := s.srv.Serve(ln); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

type ctxKey string

const adminTokenKey ctxKey = "admin_token"

// withAuth enforces a scope; /healthz and /readyz stay public.
func (s *Server) withAuth(scope domain.Scope, next func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return s.withAuthFunc(func(w http.ResponseWriter, r *http.Request, _ context.Context) {
		next(w, r)
	})
}

func (s *Server) withAuthFunc(next func(http.ResponseWriter, *http.Request, context.Context)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if s.authMode != "none" {
			tok := bearerToken(r)
			if tok == "" && r.URL.Path == "/api/v1/approvals/stream" {
				// EventSource cannot set request headers, so the live
				// stream also accepts the token as a query parameter.
				tok = r.URL.Query().Get("access_token")
			}
			if tok == "" {
				writeError(w, r, http.StatusUnauthorized, domain.CodeUnauthenticated, "missing bearer token")
				return
			}
			t, err := s.resolver.LookupByHash(ctx, hashToken(tok))
			if err != nil {
				writeError(w, r, http.StatusUnauthorized, domain.CodeTokenInvalid, "invalid token")
				return
			}
			ctx = domain.WithAuth(ctx, t.ID, t.Scopes)
			ctx = context.WithValue(ctx, adminTokenKey, t)
		} else {
			dev := domain.AuthToken{ID: "dev", Scopes: []domain.Scope{domain.ScopeAdmin}}
			ctx = domain.WithAuth(ctx, dev.ID, dev.Scopes)
			ctx = context.WithValue(ctx, adminTokenKey, dev)
		}
		if err := domain.RequireScope(ctx, requiredScope(r)); err != nil {
			if de, ok := err.(*domain.Error); ok && de.Code == domain.CodeUnauthenticated {
				writeError(w, r, http.StatusUnauthorized, de.Code, de.Message)
			} else {
				writeError(w, r, http.StatusForbidden, domain.CodeForbidden, "missing required scope")
			}
			return
		}
		next(w, r, ctx)
	}
}

// requiredScope maps calls to scopes. Mutating approval and token calls
// require admin so agents holding write_execute cannot approve their own
// writes or mint new tokens.
func requiredScope(r *http.Request) domain.Scope {
	if strings.HasSuffix(r.URL.Path, "/approve") || strings.HasSuffix(r.URL.Path, "/reject") {
		return domain.ScopeAdmin
	}
	if r.URL.Path == "/api/v1/tokens" || strings.HasPrefix(r.URL.Path, "/api/v1/tokens/") {
		return domain.ScopeAdmin
	}
	return domain.ScopeRead
}

func adminTokenID(ctx context.Context) string {
	if t, ok := ctx.Value(adminTokenKey).(domain.AuthToken); ok {
		return t.ID
	}
	return ""
}

func bearerToken(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if !strings.HasPrefix(strings.ToLower(h), "bearer ") {
		return ""
	}
	return strings.TrimSpace(h[7:])
}

func hashToken(secret string) string {
	sum := sha256.Sum256([]byte(secret))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, r *http.Request, status int, code domain.Code, msg string) {
	reqID := ""
	if r != nil {
		if id, ok := r.Context().Value(requestIDKey{}).(string); ok {
			reqID = id
		}
	}
	if reqID == "" {
		reqID = domain.NewRequestID()
	}
	writeJSON(w, status, map[string]any{
		"error": map[string]any{"code": string(code), "message": msg, "request_id": reqID},
	})
}

type requestIDKey struct{}

// ---- health ----

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	unhealthy := map[string]string{}
	for name, p := range s.pools {
		if err := p.Ping(ctx); err != nil {
			unhealthy[name] = err.Error()
		}
	}
	if len(unhealthy) > 0 {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"status": "not ready", "pools": unhealthy})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "ready"})
}

// ---- approvals ----

func approvalItem(t *domain.PreviewToken) map[string]any {
	return map[string]any{
		"id": t.ID, "state": string(t.State),
		"created_at": t.CreatedAt.UTC().Format(time.RFC3339),
		"expires_at": t.ExpiresAt.UTC().Format(time.RFC3339),
		"connection": t.Connection, "statement_type": string(t.Statement),
		"sql": t.SQL, "params": t.Params, "affected_estimate": t.Estimate,
		"warnings":     t.Warnings,
		"requested_by": map[string]any{"token_id": t.RequestedBy},
	}
}

func approvalDetail(t *domain.PreviewToken) map[string]any {
	m := approvalItem(t)
	m["sql_hash"] = t.SQLHash
	m["params_hash"] = t.ParamsHash
	m["plan"] = t.PlanText
	m["decided_by"] = nil
	if t.ApprovedBy != "" {
		m["decided_by"] = map[string]any{"id": t.ApprovedBy}
	}
	m["decision_reason"] = t.DecidedReason
	return m
}

func (s *Server) handleApprovals(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, r, http.StatusMethodNotAllowed, domain.CodeInternal, "method not allowed")
		return
	}
	all, err := s.tokens.List(r.Context())
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, domain.CodeInternal, err.Error())
		return
	}
	state := r.URL.Query().Get("state")
	conn := r.URL.Query().Get("connection")
	limit := queryInt(r, "limit", 20, 100)
	offset := queryCursor(r)
	items := []map[string]any{}
	for _, t := range all {
		if state != "" && string(t.State) != state {
			continue
		}
		if conn != "" && t.Connection != conn {
			continue
		}
		items = append(items, approvalItem(t))
	}
	total := len(items)
	if offset > total {
		offset = total
	}
	end := offset + limit
	if end > total {
		end = total
	}
	next := ""
	if end < total {
		next = strconv.Itoa(end)
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"items": items[offset:end], "next_cursor": next, "total": total,
	})
}

func (s *Server) handleApprovalSub(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	rest := strings.TrimPrefix(r.URL.Path, "/api/v1/approvals/")
	parts := strings.SplitN(rest, "/", 2)
	id := parts[0]
	action := ""
	if len(parts) == 2 {
		action = parts[1]
	}
	switch {
	case action == "" && r.Method == http.MethodGet:
		t, err := s.tokens.Get(ctx, id)
		if err != nil {
			writeError(w, r, statusForTokenErr(err), errorCode(err), errMsg(err))
			return
		}
		writeJSON(w, http.StatusOK, approvalDetail(t))
	case (action == "approve" || action == "reject") && r.Method == http.MethodPost:
		var body struct {
			Reason string `json:"reason"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		var t *domain.PreviewToken
		var err error
		if action == "approve" {
			t, err = s.tokens.Approve(ctx, id, adminTokenID(ctx), body.Reason)
		} else {
			t, err = s.tokens.Reject(ctx, id, adminTokenID(ctx), body.Reason)
		}
		if err != nil {
			writeError(w, r, statusForTokenErr(err), errorCode(err), errMsg(err))
			return
		}
		s.emitAdmin(ctx, "write."+action+"d", t)
		writeJSON(w, http.StatusOK, map[string]any{
			"id": t.ID, "state": string(t.State),
			"decided_by":      map[string]any{"id": t.ApprovedBy},
			"decided_at":      t.DecidedAt.UTC().Format(time.RFC3339),
			"decision_reason": t.DecidedReason,
		})
	default:
		writeError(w, r, http.StatusNotFound, domain.CodeInternal, "not found")
	}
}

func (s *Server) emitAdmin(ctx context.Context, event string, t *domain.PreviewToken) {
	if s.auditSink == nil {
		return
	}
	s.auditSink.Emit(ctx, domain.AuditEvent{
		TS: s.clock.Now(), Event: event,
		RequestID: domain.NewRequestID(), TokenID: adminTokenID(ctx),
		Connection: t.Connection, Tool: "admin",
		SQLHash: t.SQLHash, StatementType: t.Statement,
		Status: "success",
	})
}

// ---- SSE ----

// handleStream pushes token lifecycle events to the UI.
func (s *Server) handleStream(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	if s.tokens == nil {
		writeError(w, r, http.StatusServiceUnavailable, domain.CodeInternal, "approvals are not enabled")
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, r, http.StatusInternalServerError, domain.CodeInternal, "streaming unsupported")
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	// Flush headers immediately so clients (and proxies) know the stream
	// is live before the first event.
	flusher.Flush()
	ch, cancel := s.tokens.Watch()
	defer cancel()
	heartbeat := time.NewTicker(15 * time.Second)
	defer heartbeat.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-ch:
			writeSSE(w, "approval."+string(ev.Type), approvalDetail(ev.Token))
			flusher.Flush()
		case <-heartbeat.C:
			_, _ = fmt.Fprintln(w, ":ping")
			flusher.Flush()
		}
	}
}

func writeSSE(w http.ResponseWriter, event string, v any) {
	raw, _ := json.Marshal(v)
	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, raw)
}

// ---- audit ----

func (s *Server) handleAudit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, r, http.StatusMethodNotAllowed, domain.CodeInternal, "method not allowed")
		return
	}
	if s.audit == nil {
		writeError(w, r, http.StatusServiceUnavailable, domain.CodeInternal, "audit query is not enabled")
		return
	}
	q := r.URL.Query()
	filter := out.AuditFilter{
		Event: q.Get("event"), Connection: q.Get("connection"),
		TokenID: q.Get("token_id"), RequestID: q.Get("request_id"),
		Status: q.Get("status"), Search: q.Get("q"),
		Limit: queryInt(r, "limit", 20, 1000), Offset: queryCursor(r),
	}
	if from := q.Get("from"); from != "" {
		if ts, err := time.Parse(time.RFC3339, from); err == nil {
			filter.From = ts
		}
	}
	if to := q.Get("to"); to != "" {
		if ts, err := time.Parse(time.RFC3339, to); err == nil {
			filter.To = ts
		}
	}
	page, err := s.audit.Query(r.Context(), filter)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, domain.CodeInternal, err.Error())
		return
	}
	next := ""
	if page.NextOffset >= 0 {
		next = strconv.Itoa(page.NextOffset)
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"items": page.Events, "next_cursor": next, "total": page.Total,
	})
}

func (s *Server) handleAuditOne(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, r, http.StatusMethodNotAllowed, domain.CodeInternal, "method not allowed")
		return
	}
	if s.audit == nil {
		writeError(w, r, http.StatusServiceUnavailable, domain.CodeInternal, "audit query is not enabled")
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/audit/")
	e, err := s.audit.Get(r.Context(), id)
	if err != nil {
		writeError(w, r, http.StatusNotFound, domain.CodeInternal, "audit event not found")
		return
	}
	writeJSON(w, http.StatusOK, e)
}

// ---- connections ----

func (s *Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, r, http.StatusMethodNotAllowed, domain.CodeInternal, "method not allowed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": s.connectionItems(r.Context())})
}

func (s *Server) connectionItems(ctx context.Context) []map[string]any {
	names := make([]string, 0, len(s.metas))
	for n := range s.metas {
		names = append(names, n)
	}
	sort.Strings(names)
	items := make([]map[string]any, 0, len(names))
	for _, n := range names {
		m := s.metas[n]
		status := "healthy"
		if p, ok := s.pools[n]; ok {
			if err := p.Ping(ctx); err != nil {
				status = "unhealthy: " + err.Error()
			}
		} else {
			status = "unhealthy: no pool"
		}
		items = append(items, map[string]any{
			"name": m.Connection.Name, "driver": m.Connection.Driver,
			"readonly": m.Connection.ReadOnly, "status": status,
			"last_ping_at":    s.clock.Now().UTC().Format(time.RFC3339),
			"allowed_schemas": m.AllowedSchemas, "denied_tables": m.DeniedTables,
		})
	}
	return items
}

func (s *Server) handleConnectionSub(w http.ResponseWriter, r *http.Request, _ context.Context) {
	rest := strings.TrimPrefix(r.URL.Path, "/api/v1/connections/")
	parts := strings.SplitN(rest, "/", 2)
	if len(parts) != 2 || parts[1] != "test" || r.Method != http.MethodPost {
		writeError(w, r, http.StatusNotFound, domain.CodeInternal, "not found")
		return
	}
	p, ok := s.pools[parts[0]]
	if !ok {
		writeError(w, r, http.StatusNotFound, domain.CodeConnectionNotFound, "unknown connection")
		return
	}
	start := time.Now()
	if err := p.Ping(r.Context()); err != nil {
		writeError(w, r, http.StatusServiceUnavailable, domain.CodeDBError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"name": parts[0], "status": "healthy",
		"latency_ms": time.Since(start).Milliseconds(),
	})
}

// ---- dashboard ----

func (s *Server) handleSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, r, http.StatusMethodNotAllowed, domain.CodeInternal, "method not allowed")
		return
	}
	pending := 0
	if all, err := s.tokens.List(r.Context()); err == nil {
		for _, t := range all {
			if t.State == domain.TokenPending {
				pending++
			}
		}
	}
	total, denied := 0, 0
	if s.audit != nil {
		if page, err := s.audit.Query(r.Context(), out.AuditFilter{Limit: 1}); err == nil {
			total = page.Total
		}
		if page, err := s.audit.Query(r.Context(), out.AuditFilter{Status: "denied", Limit: 1}); err == nil {
			denied = page.Total
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"pending_approvals": pending, "requests_total": total, "denials_total": denied,
	})
}

// ---- tokens ----

func tokenItem(t domain.AuthToken, now time.Time) map[string]any {
	state := "active"
	if t.Expired(now) {
		state = "expired"
	}
	var expires any
	if !t.ExpiresAt.IsZero() {
		expires = t.ExpiresAt.UTC().Format(time.RFC3339)
	}
	scopes := make([]string, 0, len(t.Scopes))
	for _, s := range t.Scopes {
		scopes = append(scopes, string(s))
	}
	return map[string]any{
		"id": t.ID, "name": t.Name, "scopes": scopes,
		"expires_at": expires, "state": state,
	}
}

func (s *Server) handleTokens(w http.ResponseWriter, r *http.Request) {
	if s.authStore == nil {
		writeError(w, r, http.StatusServiceUnavailable, domain.CodeInternal, "token management is not enabled")
		return
	}
	switch r.Method {
	case http.MethodGet:
		items := []map[string]any{}
		for _, t := range s.authStore.List() {
			items = append(items, tokenItem(t, s.clock.Now()))
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": items})
	case http.MethodPost:
		var body struct {
			Name     string   `json:"name"`
			Scopes   []string `json:"scopes"`
			TTLHours float64  `json:"ttl_hours"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeError(w, r, http.StatusBadRequest, domain.CodeInternal, "invalid body")
			return
		}
		scopes := make([]domain.Scope, 0, len(body.Scopes))
		for _, sc := range body.Scopes {
			scopes = append(scopes, domain.Scope(sc))
		}
		var expires time.Time
		if body.TTLHours > 0 {
			expires = s.clock.Now().Add(time.Duration(body.TTLHours * float64(time.Hour)))
		}
		secret, t, err := s.authStore.Create(body.Name, scopes, expires)
		if err != nil {
			writeError(w, r, http.StatusBadRequest, domain.CodeInternal, err.Error())
			return
		}
		s.emitTokenEvent(r.Context(), "token.created", t.ID)
		item := tokenItem(t, s.clock.Now())
		item["token"] = secret // shown exactly once
		writeJSON(w, http.StatusCreated, item)
	default:
		writeError(w, r, http.StatusMethodNotAllowed, domain.CodeInternal, "method not allowed")
	}
}

func (s *Server) handleTokenOne(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	if s.authStore == nil {
		writeError(w, r, http.StatusServiceUnavailable, domain.CodeInternal, "token management is not enabled")
		return
	}
	if r.Method != http.MethodDelete {
		writeError(w, r, http.StatusNotFound, domain.CodeInternal, "not found")
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/tokens/")
	if id == "" {
		writeError(w, r, http.StatusBadRequest, domain.CodeInternal, "missing token id")
		return
	}
	// Never revoke your own token through the API.
	if id == adminTokenID(ctx) {
		writeError(w, r, http.StatusBadRequest, domain.CodeInternal, "cannot revoke the token in use")
		return
	}
	if err := s.authStore.Revoke(id); err != nil {
		writeError(w, r, http.StatusNotFound, errorCode(err), errMsg(err))
		return
	}
	s.emitTokenEvent(ctx, "token.revoked", id)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) emitTokenEvent(ctx context.Context, event, id string) {
	if s.auditSink == nil {
		return
	}
	s.auditSink.Emit(ctx, domain.AuditEvent{
		TS: s.clock.Now(), Event: event,
		RequestID: domain.NewRequestID(), TokenID: adminTokenID(ctx),
		Tool: "admin", Status: "success", Error: id,
	})
}

// ---- audit export ----

func (s *Server) handleAuditExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, r, http.StatusMethodNotAllowed, domain.CodeInternal, "method not allowed")
		return
	}
	if s.audit == nil {
		writeError(w, r, http.StatusServiceUnavailable, domain.CodeInternal, "audit query is not enabled")
		return
	}
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "jsonl"
	}
	q := r.URL.Query()
	filter := out.AuditFilter{
		Event: q.Get("event"), Connection: q.Get("connection"),
		TokenID: q.Get("token_id"), RequestID: q.Get("request_id"),
		Status: q.Get("status"), Search: q.Get("q"),
		Limit: 10000,
	}
	page, err := s.audit.Query(r.Context(), filter)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, domain.CodeInternal, err.Error())
		return
	}
	switch format {
	case "jsonl":
		w.Header().Set("Content-Type", "application/jsonl")
		w.Header().Set("Content-Disposition", `attachment; filename="janus-audit.jsonl"`)
		enc := json.NewEncoder(w)
		for _, e := range page.Events {
			_ = enc.Encode(e)
		}
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", `attachment; filename="janus-audit.csv"`)
		cw := csv.NewWriter(w)
		_ = cw.Write([]string{"id", "ts", "event", "request_id", "token_id", "connection", "tool", "statement_type", "status", "duration_ms", "error", "sql_hash"})
		for _, e := range page.Events {
			_ = cw.Write([]string{
				e.ID, e.TS.UTC().Format(time.RFC3339), e.Event, e.RequestID, e.TokenID,
				e.Connection, e.Tool, string(e.StatementType), e.Status,
				strconv.FormatInt(e.DurationMs, 10), e.Error, e.SQLHash,
			})
		}
		cw.Flush()
	default:
		writeError(w, r, http.StatusBadRequest, domain.CodeInternal, "format must be jsonl|csv")
	}
}

// ---- helpers ----

func queryInt(r *http.Request, key string, def, max int) int {
	n, err := strconv.Atoi(r.URL.Query().Get(key))
	if err != nil || n <= 0 {
		return def
	}
	if n > max {
		return max
	}
	return n
}

func queryCursor(r *http.Request) int {
	n, err := strconv.Atoi(r.URL.Query().Get("cursor"))
	if err != nil || n < 0 {
		return 0
	}
	return n
}

func statusForTokenErr(err error) int {
	switch errorCode(err) {
	case domain.CodeTokenExpired:
		return http.StatusGone
	case domain.CodeTokenMismatch:
		return http.StatusNotFound
	case domain.CodeTokenAlreadyUsed:
		return http.StatusConflict
	default:
		return http.StatusBadRequest
	}
}

func errorCode(err error) domain.Code {
	if de, ok := err.(*domain.Error); ok {
		return de.Code
	}
	return domain.CodeInternal
}

func errMsg(err error) string {
	if de, ok := err.(*domain.Error); ok {
		return de.Message
	}
	return err.Error()
}
