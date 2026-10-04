package admin

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	auditmemory "github.com/vietthanh1999/ohjanus/internal/adapter/out/audit/memory"
	authconfig "github.com/vietthanh1999/ohjanus/internal/adapter/out/auth/config"
	authmemory "github.com/vietthanh1999/ohjanus/internal/adapter/out/auth/memory"
	systemclock "github.com/vietthanh1999/ohjanus/internal/adapter/out/clock/system"
	tokememory "github.com/vietthanh1999/ohjanus/internal/adapter/out/token/memory"
	pgvalidator "github.com/vietthanh1999/ohjanus/internal/adapter/out/validator/postgres"
	valrouter "github.com/vietthanh1999/ohjanus/internal/adapter/out/validator/router"
	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
	"github.com/vietthanh1999/ohjanus/internal/core/service"
)

type stubResolver struct {
	tokens map[string]domain.AuthToken
}

func (s *stubResolver) LookupByHash(_ context.Context, hash string) (domain.AuthToken, error) {
	for h, t := range s.tokens {
		if h == hash {
			return t, nil
		}
	}
	return domain.AuthToken{}, domain.NewError(domain.CodeTokenInvalid, "unknown token")
}

func testServer(authMode string) (*Server, *tokememory.Store, *auditmemory.Buffer) {
	clock := systemclock.Clock{}
	store := tokememory.New(clock, time.Minute)
	buf := auditmemory.New(100)
	metas := map[string]domain.ConnectionMeta{
		"analytics": {Connection: domain.Connection{Name: "analytics", Driver: "postgres", ReadOnly: true}},
	}
	authStore := authmemory.New(clock)
	s := New("127.0.0.1:0", authMode, store, &stubResolver{tokens: map[string]domain.AuthToken{}}, authStore, buf, buf, service.NewConnRegistry(map[string]out.Pool{}, metas), nil, nil, nil, clock)
	return s, store, buf
}

func seedToken(t *testing.T, store *tokememory.Store) {
	t.Helper()
	err := store.Create(context.Background(), &domain.PreviewToken{
		ID: "pvw_1", Connection: "analytics",
		SQLHash: "h", ParamsHash: "p",
		SQL: "UPDATE orders SET a = 1", Statement: domain.StatementUpdate,
		RequestedBy: "tok_x",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestAdminApprovals(t *testing.T) {
	s, store, _ := testServer("none")
	seedToken(t, store)
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	var list struct {
		Items []map[string]any `json:"items"`
		Total int              `json:"total"`
	}
	if err := getJSON(ts.URL+"/api/v1/approvals", &list); err != nil {
		t.Fatal(err)
	}
	if list.Total != 1 {
		t.Fatalf("total = %d, want 1", list.Total)
	}

	var detail map[string]any
	if err := getJSON(ts.URL+"/api/v1/approvals/pvw_1", &detail); err != nil {
		t.Fatal(err)
	}
	if detail["sql"] != "UPDATE orders SET a = 1" {
		t.Errorf("detail = %v", detail)
	}

	var decided map[string]any
	if err := postJSON(ts.URL+"/api/v1/approvals/pvw_1/approve", `{"reason":"ship it"}`, &decided); err != nil {
		t.Fatal(err)
	}
	if decided["state"] != "approved" {
		t.Errorf("decided = %v", decided)
	}

	var filtered struct {
		Total int `json:"total"`
	}
	if err := getJSON(ts.URL+"/api/v1/approvals?state=pending", &filtered); err != nil {
		t.Fatal(err)
	}
	if filtered.Total != 0 {
		t.Errorf("pending total = %d, want 0", filtered.Total)
	}

	resp, err := http.Get(ts.URL + "/api/v1/approvals/missing")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", resp.StatusCode)
	}
}

func TestAdminRejectRequiresReason(t *testing.T) {
	s, store, _ := testServer("none")
	seedToken(t, store)
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/api/v1/approvals/pvw_1/reject", "application/json", strings.NewReader(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400 (reason required)", resp.StatusCode)
	}
}

func TestAdminAuth(t *testing.T) {
	secret := "jn_test_secret"
	sum := sha256.Sum256([]byte(secret))
	resolver := authconfig.New([]domain.AuthToken{
		{ID: "tok_admin", Hash: "sha256:" + hex.EncodeToString(sum[:]), Scopes: []domain.Scope{domain.ScopeAdmin}},
	}, systemclock.Clock{})
	clock := systemclock.Clock{}
	store := tokememory.New(clock, time.Minute)
	buf := auditmemory.New(10)
	s := New("127.0.0.1:0", "token", store, resolver, authmemory.New(clock), buf, buf, service.NewConnRegistry(nil, nil), nil, nil, nil, clock)
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/v1/approvals")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("no token status = %d, want 401", resp.StatusCode)
	}

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/api/v1/approvals", nil)
	req.Header.Set("Authorization", "Bearer "+secret)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("admin token status = %d, want 200", resp.StatusCode)
	}
}

func TestAdminAuditAndHealth(t *testing.T) {
	s, _, buf := testServer("none")
	buf.Emit(context.Background(), domain.AuditEvent{Event: "query.executed", Connection: "analytics", Status: "success"})
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	var page struct {
		Items []map[string]any `json:"items"`
		Total int              `json:"total"`
	}
	if err := getJSON(ts.URL+"/api/v1/audit", &page); err != nil {
		t.Fatal(err)
	}
	if page.Total != 1 {
		t.Fatalf("total = %d, want 1", page.Total)
	}

	resp, err := http.Get(ts.URL + "/healthz")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("healthz = %d", resp.StatusCode)
	}

	resp, err = http.Get(ts.URL + "/readyz")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("readyz = %d (no pools configured)", resp.StatusCode)
	}

	var conns struct {
		Items []map[string]any `json:"items"`
	}
	if err := getJSON(ts.URL+"/api/v1/connections", &conns); err != nil {
		t.Fatal(err)
	}
	if len(conns.Items) != 1 {
		t.Errorf("connections = %v", conns.Items)
	}
}

func TestAdminSSE(t *testing.T) {
	s, store, _ := testServer("none")
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/api/v1/approvals/stream", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if ct := resp.Header.Get("Content-Type"); ct != "text/event-stream" {
		t.Fatalf("content-type = %q", ct)
	}

	lines := make(chan string, 8)
	go func() {
		sc := bufio.NewScanner(resp.Body)
		for sc.Scan() {
			lines <- sc.Text()
		}
	}()

	seedToken(t, store)
	timeout := time.After(5 * time.Second)
	var gotEvent, gotData bool
	for !(gotEvent && gotData) {
		select {
		case line := <-lines:
			if strings.HasPrefix(line, "event: approval.created") {
				gotEvent = true
			}
			if strings.HasPrefix(line, "data: ") && strings.Contains(line, "pvw_1") {
				gotData = true
			}
		case <-timeout:
			t.Fatalf("timed out waiting for SSE (event=%v data=%v)", gotEvent, gotData)
		}
	}
}

func TestAdminCORS(t *testing.T) {
	s, _, _ := testServer("none")
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodOptions, ts.URL+"/api/v1/approvals", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Headers", "authorization")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("preflight status = %d, want 204", resp.StatusCode)
	}
	if got := resp.Header.Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("allow-origin = %q, want *", got)
	}

	resp, err = http.Get(ts.URL + "/api/v1/approvals")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if got := resp.Header.Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("GET allow-origin = %q, want *", got)
	}
}

func TestAdminStreamQueryToken(t *testing.T) {
	secret := "jn_stream_secret"
	sum := sha256.Sum256([]byte(secret))
	resolver := authconfig.New([]domain.AuthToken{
		{ID: "tok_admin", Hash: "sha256:" + hex.EncodeToString(sum[:]), Scopes: []domain.Scope{domain.ScopeAdmin}},
	}, systemclock.Clock{})
	clock := systemclock.Clock{}
	store := tokememory.New(clock, time.Minute)
	buf := auditmemory.New(10)
	s := New("127.0.0.1:0", "token", store, resolver, authmemory.New(clock), buf, buf, service.NewConnRegistry(nil, nil), nil, nil, nil, clock)
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	// Bearer header still works.
	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/api/v1/approvals/stream", nil)
	req.Header.Set("Authorization", "Bearer "+secret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("header token status = %d, want 200", resp.StatusCode)
	}

	// Query fallback works for EventSource clients.
	resp, err = http.Get(ts.URL + "/api/v1/approvals/stream?access_token=" + secret)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("query token status = %d, want 200", resp.StatusCode)
	}

	// Unknown query token is rejected, and query tokens do not leak
	// into other endpoints.
	resp, err = http.Get(ts.URL + "/api/v1/approvals/stream?access_token=bogus")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("bogus token status = %d, want 401", resp.StatusCode)
	}
	resp, err = http.Get(ts.URL + "/api/v1/approvals?access_token=" + secret)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("query token on list status = %d, want 401", resp.StatusCode)
	}
}

func TestAdminTokens(t *testing.T) {
	s, _, _ := testServer("none")
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	var created map[string]any
	resp, err := http.Post(ts.URL+"/api/v1/tokens", "application/json", strings.NewReader(`{"name":"ci","scopes":["read"],"ttl_hours":1}`))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		resp.Body.Close()
		t.Fatalf("create status = %d, want 201", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		resp.Body.Close()
		t.Fatal(err)
	}
	resp.Body.Close()
	secret, _ := created["token"].(string)
	if !strings.HasPrefix(secret, "jn_") {
		t.Errorf("token = %q, want jn_ prefix", secret)
	}
	id, _ := created["id"].(string)

	var list struct {
		Items []map[string]any `json:"items"`
	}
	if err := getJSON(ts.URL+"/api/v1/tokens", &list); err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 1 || list.Items[0]["name"] != "ci" {
		t.Errorf("list = %v", list.Items)
	}

	var bad map[string]any
	if err := postJSON(ts.URL+"/api/v1/tokens", `{"name":"x","scopes":["bogus"]}`, &bad); err == nil {
		t.Error("unknown scope should fail")
	}

	req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/v1/tokens/"+id, nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("delete status = %d, want 204", resp.StatusCode)
	}
	if err := getJSON(ts.URL+"/api/v1/tokens", &list); err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 0 {
		t.Errorf("list after revoke = %v", list.Items)
	}
}

func TestAdminAuditExport(t *testing.T) {
	s, _, buf := testServer("none")
	buf.Emit(context.Background(), domain.AuditEvent{Event: "query.executed", Connection: "analytics", Status: "success"})
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/v1/audit/export?format=csv")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if ct := resp.Header.Get("Content-Type"); ct != "text/csv" {
		t.Errorf("content-type = %q", ct)
	}
	sc := bufio.NewScanner(resp.Body)
	lines := 0
	for sc.Scan() {
		lines++
	}
	if lines != 2 { // header + 1 event
		t.Errorf("csv lines = %d, want 2", lines)
	}

	resp2, err := http.Get(ts.URL + "/api/v1/audit/export?format=bogus")
	if err != nil {
		t.Fatal(err)
	}
	resp2.Body.Close()
	if resp2.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", resp2.StatusCode)
	}
}

func getJSON(url string, v any) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return &testHTTPError{status: resp.StatusCode, url: url}
	}
	return json.NewDecoder(resp.Body).Decode(v)
}

func postJSON(url, body string, v any) error {
	resp, err := http.Post(url, "application/json", strings.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return &testHTTPError{status: resp.StatusCode, url: url}
	}
	return json.NewDecoder(resp.Body).Decode(v)
}

type testHTTPError struct {
	status int
	url    string
}

func (e *testHTTPError) Error() string {
	return http.StatusText(e.status) + " for " + e.url
}

type stubPool struct{}

func (stubPool) Query(_ context.Context, _ domain.Query, _ domain.QueryOpts) (*domain.ResultSet, error) {
	return &domain.ResultSet{}, nil
}
func (stubPool) Exec(_ context.Context, _ domain.Query) (*domain.ExecResult, error) {
	return &domain.ExecResult{}, nil
}
func (stubPool) Explain(_ context.Context, _ domain.Query) (*domain.Plan, error) {
	return &domain.Plan{}, nil
}
func (stubPool) Schema(_ context.Context, _, _ string) ([]domain.Schema, error) { return nil, nil }
func (stubPool) Ping(_ context.Context) error                                   { return nil }
func (stubPool) Close() error                                                   { return nil }

type stubConnector struct {
	pool out.Pool
	err  error
}

func (s *stubConnector) Driver() string { return "postgres" }
func (s *stubConnector) Open(_ context.Context, _ domain.Credentials, _ domain.PoolConfig) (out.Pool, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.pool, nil
}

func testConnServer(t *testing.T) (*httptest.Server, *Server, string, string) {
	t.Helper()
	adminSecret := "jn_admin_secret"
	readSecret := "jn_read_secret"
	adminSum := sha256.Sum256([]byte(adminSecret))
	readSum := sha256.Sum256([]byte(readSecret))
	resolver := authconfig.New([]domain.AuthToken{
		{ID: "tok_admin", Hash: "sha256:" + hex.EncodeToString(adminSum[:]), Scopes: []domain.Scope{domain.ScopeAdmin}},
		{ID: "tok_read", Hash: "sha256:" + hex.EncodeToString(readSum[:]), Scopes: []domain.Scope{domain.ScopeRead}},
	}, systemclock.Clock{})
	clock := systemclock.Clock{}
	store := tokememory.New(clock, time.Minute)
	buf := auditmemory.New(10)
	pgVal := pgvalidator.New(nil, nil)
	router := valrouter.New(map[string]out.Validator{}, nil)
	s := New("127.0.0.1:0", "token", store, resolver, authmemory.New(clock), buf, buf,
		service.NewConnRegistry(nil, nil), router, pgVal, &stubConnector{pool: stubPool{}}, clock)
	ts := httptest.NewServer(s.mux)
	t.Cleanup(ts.Close)
	return ts, s, adminSecret, readSecret
}

func postConn(t *testing.T, ts *httptest.Server, secret, body string) (int, map[string]any) {
	t.Helper()
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/v1/connections", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if secret != "" {
		req.Header.Set("Authorization", "Bearer "+secret)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var item map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&item)
	return resp.StatusCode, item
}

func TestAdminConnectionCreate(t *testing.T) {
	ts, s, adminSecret, readSecret := testConnServer(t)

	// Valid field-based create → 201, tested-on-create (stub pings OK).
	status, item := postConn(t, ts, adminSecret, `{"name":"prod","host":"db.internal","port":5432,"database":"app","username":"app","password":"s3cret","read_only":true}`)
	if status != http.StatusCreated {
		t.Fatalf("create status = %d, item = %v", status, item)
	}
	if item["name"] != "prod" || item["driver"] != "postgres" || item["status"] != "healthy" {
		t.Errorf("item = %v", item)
	}

	// Duplicate name → 409.
	if status, _ := postConn(t, ts, adminSecret, `{"name":"prod","host":"db.internal","database":"app"}`); status != http.StatusConflict {
		t.Errorf("duplicate status = %d, want 409", status)
	}

	// Bad name → 400.
	if status, _ := postConn(t, ts, adminSecret, `{"name":"bad name!","host":"h","database":"d"}`); status != http.StatusBadRequest {
		t.Errorf("bad name status = %d, want 400", status)
	}

	// Unsupported driver → 400.
	if status, _ := postConn(t, ts, adminSecret, `{"name":"m","driver":"mysql","host":"h","database":"d"}`); status != http.StatusBadRequest {
		t.Errorf("bad driver status = %d, want 400", status)
	}

	// Read-scoped token may list but not create → 403.
	if status, _ := postConn(t, ts, readSecret, `{"name":"other","host":"h","database":"d"}`); status != http.StatusForbidden {
		t.Errorf("read token create status = %d, want 403", status)
	}

	// Unauthenticated → 401.
	if status, _ := postConn(t, ts, "", `{"name":"anon","host":"h","database":"d"}`); status != http.StatusUnauthorized {
		t.Errorf("anon create status = %d, want 401", status)
	}

	// Created connection shows up in the list.
	var list struct {
		Items []map[string]any `json:"items"`
	}
	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/api/v1/connections", nil)
	req.Header.Set("Authorization", "Bearer "+adminSecret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 1 || list.Items[0]["name"] != "prod" {
		t.Errorf("list = %v", list.Items)
	}

	// Probe tests parameters without saving: valid → 200, nothing added.
	probeReq, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/v1/connections/test", strings.NewReader(`{"host":"db.internal","database":"app"}`))
	probeReq.Header.Set("Content-Type", "application/json")
	probeReq.Header.Set("Authorization", "Bearer "+adminSecret)
	probeResp, err := http.DefaultClient.Do(probeReq)
	if err != nil {
		t.Fatal(err)
	}
	defer probeResp.Body.Close()
	if probeResp.StatusCode != http.StatusOK {
		t.Errorf("probe status = %d, want 200", probeResp.StatusCode)
	}
	if s.conns.Has("db.internal") || len(list.Items) != 1 {
		t.Errorf("probe must not register a connection")
	}

	// Probe with missing fields → 400; probe with read token → 403.
	probeReq, _ = http.NewRequest(http.MethodPost, ts.URL+"/api/v1/connections/test", strings.NewReader(`{"host":"","database":""}`))
	probeReq.Header.Set("Content-Type", "application/json")
	probeReq.Header.Set("Authorization", "Bearer "+adminSecret)
	if resp, err := http.DefaultClient.Do(probeReq); err != nil {
		t.Fatal(err)
	} else {
		resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("bad probe status = %d, want 400", resp.StatusCode)
		}
	}
	probeReq, _ = http.NewRequest(http.MethodPost, ts.URL+"/api/v1/connections/test", strings.NewReader(`{"host":"h","database":"d"}`))
	probeReq.Header.Set("Content-Type", "application/json")
	probeReq.Header.Set("Authorization", "Bearer "+readSecret)
	if resp, err := http.DefaultClient.Do(probeReq); err != nil {
		t.Fatal(err)
	} else {
		resp.Body.Close()
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("read probe status = %d, want 403", resp.StatusCode)
		}
	}
}
