package http

import (
	"bufio"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/vietthanh1999/ohjanus/internal/adapter/in/mcp"
	authmemory "github.com/vietthanh1999/ohjanus/internal/adapter/out/auth/memory"
	systemclock "github.com/vietthanh1999/ohjanus/internal/adapter/out/clock/system"
	tokememory "github.com/vietthanh1999/ohjanus/internal/adapter/out/token/memory"
	"github.com/vietthanh1999/ohjanus/internal/core/domain"
)

func testTransport(authMode string) (*Server, string) {
	clock := systemclock.Clock{}
	authStore := authmemory.New(clock)
	secret, token, err := authStore.Create("test", []domain.Scope{domain.ScopeRead}, time.Now().Add(time.Hour))
	if err != nil {
		panic(err)
	}
	_ = token
	srv := mcp.NewServer("janus", "0.1.0", authMode, authStore, clock, nil, nil, nil)
	return New(srv, tokememory.New(clock, time.Minute), authMode, TLSConfig{}, nil), secret
}

func post(t *testing.T, ts *httptest.Server, token, body string) (int, map[string]any) {
	t.Helper()
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/mcp", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var out map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&out)
	return resp.StatusCode, out
}

func TestInitPing(t *testing.T) {
	s, secret := testTransport("token")
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	status, out := post(t, ts, secret, `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`)
	if status != http.StatusOK {
		t.Fatalf("status = %d", status)
	}
	if out["result"] == nil {
		t.Errorf("missing result: %v", out)
	}

	status, _ = post(t, ts, secret, `{"jsonrpc":"2.0","id":2,"method":"ping"}`)
	if status != http.StatusOK {
		t.Errorf("ping status = %d", status)
	}
}

func TestAuthStatuses(t *testing.T) {
	s, _ := testTransport("token")
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	status, _ := post(t, ts, "", `{"jsonrpc":"2.0","id":1,"method":"ping"}`)
	if status != http.StatusUnauthorized {
		t.Errorf("missing token status = %d, want 401", status)
	}

	status, out := post(t, ts, "jn_wrong", `{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}`)
	if status != http.StatusUnauthorized {
		t.Errorf("wrong token status = %d, want 401", status)
	}
	if out["error"] == nil {
		t.Errorf("want JSON-RPC error body, got %v", out)
	}
}

func TestBatchAndMethod(t *testing.T) {
	s, secret := testTransport("token")
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/mcp", strings.NewReader(
		`[{"jsonrpc":"2.0","id":1,"method":"ping"},{"jsonrpc":"2.0","id":2,"method":"ping"}]`))
	req.Header.Set("Authorization", "Bearer "+secret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var batch []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&batch); err != nil {
		t.Fatal(err)
	}
	if len(batch) != 2 {
		t.Errorf("batch responses = %d, want 2", len(batch))
	}

	resp2, err := http.Get(ts.URL + "/mcp")
	if err != nil {
		t.Fatal(err)
	}
	resp2.Body.Close()
	if resp2.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("GET /mcp status = %d, want 405", resp2.StatusCode)
	}
}

func TestSSEHeadersAndEvent(t *testing.T) {
	s, secret := testTransport("token")
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/mcp/sse", nil)
	req.Header.Set("Authorization", "Bearer "+secret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if ct := resp.Header.Get("Content-Type"); ct != "text/event-stream" {
		t.Fatalf("content-type = %q", ct)
	}

	store := s.tokens
	if err := store.Create(context.Background(), &domain.PreviewToken{ID: "pvw_1", Connection: "a"}); err != nil {
		t.Fatal(err)
	}
	sc := bufio.NewScanner(resp.Body)
	deadline := time.After(5 * time.Second)
	var gotEvent, gotData bool
	for !(gotEvent && gotData) {
		select {
		case <-time.After(0):
			if !sc.Scan() {
				t.Fatal("stream closed")
			}
			line := sc.Text()
			if strings.HasPrefix(line, "event: message") {
				gotEvent = true
			}
			if strings.Contains(line, "janus-approval") && strings.Contains(line, "pvw_1") {
				gotData = true
			}
		case <-deadline:
			t.Fatalf("timed out (event=%v data=%v)", gotEvent, gotData)
		}
	}
}
