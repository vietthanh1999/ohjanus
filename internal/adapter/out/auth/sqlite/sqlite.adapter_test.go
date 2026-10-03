package sqlite

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
)

type mutableClock struct{ now time.Time }

func (c *mutableClock) Now() time.Time { return c.now }

func openTest(t *testing.T) (*Store, *mutableClock) {
	t.Helper()
	clk := &mutableClock{now: time.Now().Truncate(time.Second)}
	st, err := Open(filepath.Join(t.TempDir(), "tokens.db"), clk)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	return st, clk
}

func TestCreateLookupRoundTrip(t *testing.T) {
	st, _ := openTest(t)
	ctx := context.Background()
	secret, created, err := st.Create("agent-1", []domain.Scope{domain.ScopeRead}, time.Time{})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if secret == "" || created.Hash == "" {
		t.Fatal("want secret and stored hash")
	}
	got, err := st.LookupByHash(ctx, created.Hash)
	if err != nil {
		t.Fatalf("LookupByHash: %v", err)
	}
	if got.ID != created.ID || got.Name != "agent-1" || len(got.Scopes) != 1 {
		t.Fatalf("round trip mismatch: %+v", got)
	}
	if _, err := st.LookupByHash(ctx, "sha256:dead"); err == nil {
		t.Fatal("unknown hash: want error")
	}
}

func TestCreateRejectsBadScopes(t *testing.T) {
	st, _ := openTest(t)
	if _, _, err := st.Create("x", nil, time.Time{}); err == nil {
		t.Fatal("empty scopes: want error")
	}
	if _, _, err := st.Create("x", []domain.Scope{"root"}, time.Time{}); err == nil {
		t.Fatal("unknown scope: want error")
	}
}

func TestListSortedAndRevoke(t *testing.T) {
	st, _ := openTest(t)
	ctx := context.Background()
	if _, _, err := st.Create("b", []domain.Scope{domain.ScopeRead}, time.Time{}); err != nil {
		t.Fatal(err)
	}
	secret, a, err := st.Create("a", []domain.Scope{domain.ScopeAdmin}, time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	_ = secret
	list := st.List()
	if len(list) != 2 || list[0].ID > list[1].ID {
		t.Fatalf("List not sorted: %v", list)
	}
	if err := st.Revoke(a.ID); err != nil {
		t.Fatalf("Revoke: %v", err)
	}
	if _, err := st.LookupByHash(ctx, a.Hash); err == nil {
		t.Fatal("revoked token still resolves")
	}
	if err := st.Revoke("tok_missing"); err == nil {
		t.Fatal("revoke unknown: want error")
	}
	if len(st.List()) != 1 {
		t.Fatal("want 1 token after revoke")
	}
}

func TestExpiryEnforced(t *testing.T) {
	st, clk := openTest(t)
	ctx := context.Background()
	_, created, err := st.Create("short", []domain.Scope{domain.ScopeRead}, clk.now.Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.LookupByHash(ctx, created.Hash); err != nil {
		t.Fatalf("before expiry: %v", err)
	}
	clk.now = clk.now.Add(2 * time.Minute)
	if _, err := st.LookupByHash(ctx, created.Hash); err == nil {
		t.Fatal("after expiry: want error")
	}
}

func TestPersistsAcrossReopen(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tokens.db")
	clk := &mutableClock{now: time.Now().Truncate(time.Second)}
	st, err := Open(path, clk)
	if err != nil {
		t.Fatal(err)
	}
	_, created, err := st.Create("agent-1", []domain.Scope{domain.ScopeRead}, time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	st2, err := Open(path, clk)
	if err != nil {
		t.Fatal(err)
	}
	defer st2.Close()
	got, err := st2.LookupByHash(context.Background(), created.Hash)
	if err != nil {
		t.Fatalf("token lost across reopen: %v", err)
	}
	if got.Name != "agent-1" {
		t.Fatalf("wrong token after reopen: %+v", got)
	}
}

func TestSeedMergesAndIsIdempotent(t *testing.T) {
	st, _ := openTest(t)
	ctx := context.Background()
	cfg := []domain.AuthToken{{
		ID: "tok_cfg", Name: "from-config", Hash: "sha256:abc",
		Scopes: []domain.Scope{domain.ScopeRead},
	}}
	if err := st.Seed(cfg); err != nil {
		t.Fatalf("Seed: %v", err)
	}
	if err := st.Seed(cfg); err != nil {
		t.Fatalf("Seed twice: %v", err)
	}
	if len(st.List()) != 1 {
		t.Fatalf("want 1 token after double seed, got %d", len(st.List()))
	}
	if _, err := st.LookupByHash(ctx, "sha256:abc"); err != nil {
		t.Fatalf("seeded token missing: %v", err)
	}
	// Existing rows win over later seeds with the same id or hash.
	_, rt, err := st.Create("runtime", []domain.Scope{domain.ScopeAdmin}, time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Seed([]domain.AuthToken{{
		ID: rt.ID, Name: "hijack", Hash: "sha256:zzz",
		Scopes: []domain.Scope{domain.ScopeRead},
	}}); err != nil {
		t.Fatal(err)
	}
	got, err := st.LookupByHash(ctx, rt.Hash)
	if err != nil {
		t.Fatalf("runtime token clobbered by seed: %v", err)
	}
	if got.Name != "runtime" {
		t.Fatalf("seed overwrote runtime row: %+v", got)
	}
}

func TestDBFilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "tokens.db")
	clk := &mutableClock{now: time.Now()}
	st, err := Open(path, clk)
	if err != nil {
		t.Fatal(err)
	}
	st.Close()
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm() != 0o600 {
		t.Fatalf("db perms = %o, want 600", fi.Mode().Perm())
	}
}

func TestOpenRequiresPath(t *testing.T) {
	if _, err := Open("", &mutableClock{}); err == nil {
		t.Fatal("empty path: want error")
	}
}
