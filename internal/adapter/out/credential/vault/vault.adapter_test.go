package vault

import (
	"context"
	"encoding/binary"
	"io"
	"strings"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/vault"
)

func TestSplitRef(t *testing.T) {
	cases := []struct {
		ref                string
		mount, path, field string
		expectError        bool
	}{
		{"vault://secret/data/janus/analytics#dsn", "secret", "janus/analytics", "dsn", false},
		{"vault://secret/janus/analytics#password", "secret", "janus/analytics", "password", false},
		{"vault://secret/janus/analytics", "secret", "janus/analytics", "dsn", false},
		{"vault://secret", "", "", "", true},
		{"vault://secret/", "", "", "", true},
		{"vault://", "", "", "", true},
		{"env:X", "", "", "", true},
	}
	for _, tc := range cases {
		mount, path, field, err := splitRef(tc.ref)
		if tc.expectError {
			if err == nil {
				t.Errorf("splitRef(%q) should fail", tc.ref)
			}
			continue
		}
		if err != nil {
			t.Errorf("splitRef(%q): %v", tc.ref, err)
			continue
		}
		if mount != tc.mount || path != tc.path || field != tc.field {
			t.Errorf("splitRef(%q) = %q/%q#%q, want %q/%q#%q",
				tc.ref, mount, path, field, tc.mount, tc.path, tc.field)
		}
	}
}

func ctrExec(t *testing.T, ctx context.Context, ctr *vault.VaultContainer, cmd ...string) string {
	t.Helper()
	code, reader, err := ctr.Exec(ctx, cmd)
	if err != nil {
		t.Fatalf("exec %v: %v", cmd, err)
	}
	raw, _ := io.ReadAll(reader)
	out := demux(raw)
	if code != 0 {
		t.Fatalf("exec %v exit %d: %s", cmd, code, out)
	}
	return strings.TrimSpace(out)
}

// demux strips Docker multiplexed-stream headers (8 bytes per frame).
func demux(raw []byte) string {
	var sb strings.Builder
	for len(raw) >= 8 {
		n := int(binary.BigEndian.Uint32(raw[4:8]))
		if len(raw) < 8+n {
			break
		}
		sb.Write(raw[8 : 8+n])
		raw = raw[8+n:]
	}
	return sb.String()
}

func TestVaultApproleIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	ctx := context.Background()
	ctr, err := vault.Run(ctx, "hashicorp/vault:1.15",
		vault.WithToken("root-token"),
	)
	if err != nil {
		t.Fatalf("run vault: %v", err)
	}
	t.Cleanup(func() { _ = ctr.Terminate(context.Background()) })

	ctrExec(t, ctx, ctr, "vault", "kv", "put", "secret/janus/analytics", "dsn=postgres://vault-user:vault-pass@db:5432/app")
	ctrExec(t, ctx, ctr, "vault", "auth", "enable", "approle")
	ctrExec(t, ctx, ctr, "sh", "-c", `echo 'path "secret/data/janus/*" { capabilities = ["read"] }' | vault policy write janus-ro -`)
	ctrExec(t, ctx, ctr, "vault", "write", "auth/approle/role/janus", "token_policies=janus-ro")
	roleID := ctrExec(t, ctx, ctr, "vault", "read", "-field=role_id", "auth/approle/role/janus/role-id")
	secretID := ctrExec(t, ctx, ctr, "vault", "write", "-f", "-field=secret_id", "auth/approle/role/janus/secret-id")

	addr, err := ctr.HttpHostAddress(ctx)
	if err != nil {
		t.Fatalf("host address: %v", err)
	}
	r, err := New(addr, roleID, secretID)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	creds, err := r.Resolve(ctx, "vault://secret/data/janus/analytics#dsn")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if creds.DSN != "postgres://vault-user:vault-pass@db:5432/app" {
		t.Errorf("wrong DSN resolved")
	}

	if _, err := r.Resolve(ctx, "vault://secret/data/janus/analytics#missing"); err == nil {
		t.Error("missing field should fail")
	}
	if _, err := r.Resolve(ctx, "vault://secret/data/janus/nope#dsn"); err == nil {
		t.Error("missing secret should fail")
	}
}
