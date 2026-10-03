package config

import (
	"os"
	"path/filepath"
	"testing"
)

const minimalYAML = `
server:
  name: janus
  transport: stdio
auth:
  mode: none
connections:
  - name: analytics
    driver: postgres
    dsn_ref: "env:JANUS_ANALYTICS_DSN"
    readonly: true
    allowed_schemas: ["public"]
policy:
  default_action: deny
  rules:
    - name: allow-select
      match:
        statement: ["SELECT"]
      action: allow
limits:
  max_query_length: 100
  query_timeout: "5s"
`

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "janus.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadMinimal(t *testing.T) {
	cfg, err := Load(writeTemp(t, minimalYAML))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Server.Transport != "stdio" {
		t.Errorf("transport = %q, want stdio", cfg.Server.Transport)
	}
	if len(cfg.Connections) != 1 || cfg.Connections[0].Name != "analytics" {
		t.Errorf("connections = %+v, want [analytics]", cfg.Connections)
	}
	if cfg.Connections[0].RowLimit != 1000 {
		t.Errorf("row_limit = %d, want default 1000", cfg.Connections[0].RowLimit)
	}
	if cfg.Limits.MaxQueryLength != 100 {
		t.Errorf("max_query_length = %d, want 100", cfg.Limits.MaxQueryLength)
	}
}

func TestEnvOverride(t *testing.T) {
	t.Setenv("JANUS_SERVER_TRANSPORT", "http")
	t.Setenv("JANUS_AUTH_MODE", "none")
	t.Setenv("JANUS_OBSERVABILITY_LOG_LEVEL", "debug")
	cfg, err := Load(writeTemp(t, minimalYAML))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Server.Transport != "http" {
		t.Errorf("transport = %q, want http (env)", cfg.Server.Transport)
	}
	if cfg.Auth.Mode != "none" {
		t.Errorf("auth.mode = %q, want none (env)", cfg.Auth.Mode)
	}
	if cfg.Observability.LogLevel != "debug" {
		t.Errorf("log_level = %q, want debug (env)", cfg.Observability.LogLevel)
	}
}

func TestValidateBadDriver(t *testing.T) {
	bad := `
server:
  transport: stdio
auth:
  mode: none
connections:
  - name: x
    driver: oracle
    dsn_ref: "env:X"
policy:
  default_action: deny
`
	if _, err := Load(writeTemp(t, bad)); err == nil {
		t.Error("expected error for bad driver, got nil")
	}
}

func TestValidateBadDuration(t *testing.T) {
	bad := `
server:
  transport: stdio
auth:
  mode: none
connections: []
policy:
  default_action: deny
limits:
  query_timeout: "not-a-duration"
`
	if _, err := Load(writeTemp(t, bad)); err == nil {
		t.Error("expected error for bad duration, got nil")
	}
}

func TestConnectionMetas(t *testing.T) {
	cfg, err := Load(writeTemp(t, minimalYAML))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	metas := cfg.ConnectionMetas()
	m, ok := metas["analytics"]
	if !ok {
		t.Fatal("missing analytics meta")
	}
	if !m.Connection.ReadOnly {
		t.Error("expected readonly connection")
	}
	if len(m.AllowedSchemas) != 1 || m.AllowedSchemas[0] != "public" {
		t.Errorf("allowed_schemas = %v", m.AllowedSchemas)
	}
}
