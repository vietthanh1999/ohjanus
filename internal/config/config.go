package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"gopkg.in/yaml.v3"
)

// Config mirrors configs/janus.example.yaml (§3.1 of the spec).
type Config struct {
	Server        ServerConfig        `yaml:"server"`
	Admin         AdminConfig         `yaml:"admin"`
	Auth          AuthConfig          `yaml:"auth"`
	Connections   []ConnectionConfig  `yaml:"connections"`
	Policy        PolicyConfig        `yaml:"policy"`
	Approval      ApprovalConfig      `yaml:"approval"`
	Credential    CredentialConfig    `yaml:"credential"`
	Audit         AuditConfig         `yaml:"audit"`
	Redaction     RedactionConfig     `yaml:"redaction"`
	Observability ObservabilityConfig `yaml:"observability"`
	Limits        LimitsConfig        `yaml:"limits"`
}

type ServerConfig struct {
	Name      string     `yaml:"name"`
	Version   string     `yaml:"version"`
	Transport string     `yaml:"transport"`
	HTTP      HTTPConfig `yaml:"http"`
}

type HTTPConfig struct {
	Listen string    `yaml:"listen"`
	TLS    TLSConfig `yaml:"tls"`
}

type TLSConfig struct {
	Enabled  bool   `yaml:"enabled"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

// AdminConfig controls the Admin API for the UI (port 8788, separate from MCP).
type AdminConfig struct {
	Enabled bool   `yaml:"enabled"`
	Listen  string `yaml:"listen"`
}

type AuthConfig struct {
	Mode   string        `yaml:"mode"`
	Tokens []TokenConfig `yaml:"tokens"`
}

type TokenConfig struct {
	ID        string    `yaml:"id"`
	Hash      string    `yaml:"hash"`
	Scopes    []string  `yaml:"scopes"`
	ExpiresAt time.Time `yaml:"expires_at"`
}

type ConnectionConfig struct {
	Name           string     `yaml:"name"`
	Driver         string     `yaml:"driver"`
	DSNRef         string     `yaml:"dsn_ref"`
	Pool           PoolConfig `yaml:"pool"`
	ReadOnly       bool       `yaml:"readonly"`
	AllowedSchemas []string   `yaml:"allowed_schemas"`
	DeniedTables   []string   `yaml:"denied_tables"`
	AllowedTables  []string   `yaml:"allowed_tables"`
	RowLimit       int        `yaml:"row_limit"`
}

type PoolConfig struct {
	MaxOpen         int    `yaml:"max_open"`
	MaxIdle         int    `yaml:"max_idle"`
	ConnMaxLifetime string `yaml:"conn_max_lifetime"`
}

type PolicyConfig struct {
	DefaultAction string       `yaml:"default_action"`
	Rules         []RuleConfig `yaml:"rules"`
}

type RuleConfig struct {
	Name   string          `yaml:"name"`
	Match  RuleMatchConfig `yaml:"match"`
	Action string          `yaml:"action"`
	Reason string          `yaml:"reason"`
}

type RuleMatchConfig struct {
	Statement  []string `yaml:"statement"`
	Connection []string `yaml:"connection"`
	Tables     []string `yaml:"tables"`
	Functions  []string `yaml:"functions"`
}

type ApprovalConfig struct {
	Enabled  bool      `yaml:"enabled"`
	Method   string    `yaml:"method"`
	CLI      CLIConfig `yaml:"cli"`
	TokenTTL string    `yaml:"token_ttl"`
}

type CLIConfig struct {
	Prompt bool `yaml:"prompt"`
}

type CredentialConfig struct {
	Resolvers []ResolverConfig `yaml:"resolvers"`
}

type ResolverConfig struct {
	Scheme        string `yaml:"scheme"`
	ServicePrefix string `yaml:"service_prefix"`
	BasePath      string `yaml:"base_path"`
	FileMode      uint32 `yaml:"file_mode"`
	Prefix        string `yaml:"prefix"`
	Address       string `yaml:"address"`
	AuthMethod    string `yaml:"auth_method"`
	RoleIDRef     string `yaml:"role_id_ref"`
	SecretIDRef   string `yaml:"secret_id_ref"`
}

type AuditConfig struct {
	Enabled bool         `yaml:"enabled"`
	Sinks   []SinkConfig `yaml:"sinks"`
}

type SinkConfig struct {
	Type   string       `yaml:"type"`
	Path   string       `yaml:"path"`
	Format string       `yaml:"format"`
	Rotate RotateConfig `yaml:"rotate"`
}

type RotateConfig struct {
	MaxSizeMB  int `yaml:"max_size_mb"`
	MaxBackups int `yaml:"max_backups"`
}

type RedactionConfig struct {
	Enabled  bool            `yaml:"enabled"`
	Patterns []PatternConfig `yaml:"patterns"`
	Columns  []string        `yaml:"columns"`
}

type PatternConfig struct {
	Name        string `yaml:"name"`
	Regex       string `yaml:"regex"`
	Replacement string `yaml:"replacement"`
}

type ObservabilityConfig struct {
	LogLevel string        `yaml:"log_level"`
	Metrics  MetricsConfig `yaml:"metrics"`
	Tracing  TracingConfig `yaml:"tracing"`
}

type MetricsConfig struct {
	Enabled bool   `yaml:"enabled"`
	Listen  string `yaml:"listen"`
	Path    string `yaml:"path"`
}

type TracingConfig struct {
	Enabled      bool   `yaml:"enabled"`
	OTLPEndpoint string `yaml:"otlp_endpoint"`
}

type LimitsConfig struct {
	MaxQueryLength       int             `yaml:"max_query_length"`
	QueryTimeout         string          `yaml:"query_timeout"`
	MaxConcurrentQueries int             `yaml:"max_concurrent_queries"`
	RateLimit            RateLimitConfig `yaml:"rate_limit"`
}

type RateLimitConfig struct {
	Enabled           bool `yaml:"enabled"`
	RequestsPerMinute int  `yaml:"requests_per_minute"`
}

// Load reads, parses, overrides (ENV) and validates a config file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parse config %q: %w", path, err)
	}
	c.SetDefaults()
	applyEnvOverrides(&c)
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return &c, nil
}

// SetDefaults fills zero values with spec defaults.
func (c *Config) SetDefaults() {
	if c.Server.Transport == "" {
		c.Server.Transport = "stdio"
	}
	if c.Admin.Listen == "" {
		c.Admin.Listen = "127.0.0.1:8788"
	}
	if c.Auth.Mode == "" {
		c.Auth.Mode = "token"
	}
	if c.Policy.DefaultAction == "" {
		c.Policy.DefaultAction = "deny"
	}
	if c.Observability.LogLevel == "" {
		c.Observability.LogLevel = "info"
	}
	if c.Limits.MaxQueryLength == 0 {
		c.Limits.MaxQueryLength = 10000
	}
	if c.Limits.QueryTimeout == "" {
		c.Limits.QueryTimeout = "30s"
	}
	if c.Limits.MaxConcurrentQueries == 0 {
		c.Limits.MaxConcurrentQueries = 20
	}
	if c.Approval.TokenTTL == "" {
		c.Approval.TokenTTL = "5m"
	}
	for i := range c.Connections {
		if c.Connections[i].RowLimit == 0 {
			c.Connections[i].RowLimit = 1000
		}
	}
}

// envOverride maps JANUS_* variables onto config fields.
var envOverride = map[string]func(*Config, string){
	"JANUS_SERVER_TRANSPORT":        func(c *Config, v string) { c.Server.Transport = v },
	"JANUS_SERVER_HTTP_LISTEN":      func(c *Config, v string) { c.Server.HTTP.Listen = v },
	"JANUS_AUTH_MODE":               func(c *Config, v string) { c.Auth.Mode = v },
	"JANUS_OBSERVABILITY_LOG_LEVEL": func(c *Config, v string) { c.Observability.LogLevel = v },
	"JANUS_LIMITS_MAX_QUERY_LENGTH": func(c *Config, v string) { fmt.Sscan(v, &c.Limits.MaxQueryLength) },
	"JANUS_LIMITS_QUERY_TIMEOUT":    func(c *Config, v string) { c.Limits.QueryTimeout = v },
	"JANUS_APPROVAL_TOKEN_TTL":      func(c *Config, v string) { c.Approval.TokenTTL = v },
	"JANUS_CREDENTIAL_FILE_BASE_PATH": func(c *Config, v string) {
		for i := range c.Credential.Resolvers {
			if c.Credential.Resolvers[i].Scheme == "file" {
				c.Credential.Resolvers[i].BasePath = v
			}
		}
	},
}

func applyEnvOverrides(c *Config) {
	for key, apply := range envOverride {
		if v, ok := os.LookupEnv(key); ok && v != "" {
			apply(c, v)
		}
	}
}

// Validate checks the config for spec violations (fail fast at startup).
func (c *Config) Validate() error {
	switch c.Server.Transport {
	case "stdio", "http", "sse":
	default:
		return fmt.Errorf("server.transport must be stdio|http|sse, got %q", c.Server.Transport)
	}
	switch c.Auth.Mode {
	case "token", "none":
	default:
		return fmt.Errorf("auth.mode must be token|none, got %q", c.Auth.Mode)
	}
	switch c.Policy.DefaultAction {
	case "allow", "deny":
	default:
		return fmt.Errorf("policy.default_action must be allow|deny, got %q", c.Policy.DefaultAction)
	}
	seen := map[string]struct{}{}
	for i, conn := range c.Connections {
		if conn.Name == "" {
			return fmt.Errorf("connections[%d].name is required", i)
		}
		if _, dup := seen[conn.Name]; dup {
			return fmt.Errorf("duplicate connection name %q", conn.Name)
		}
		seen[conn.Name] = struct{}{}
		switch conn.Driver {
		case "postgres", "mysql", "sqlite":
		default:
			return fmt.Errorf("connections[%q].driver must be postgres|mysql|sqlite, got %q", conn.Name, conn.Driver)
		}
		if conn.DSNRef == "" {
			return fmt.Errorf("connections[%q].dsn_ref is required", conn.Name)
		}
		if conn.Pool.ConnMaxLifetime != "" {
			if _, err := time.ParseDuration(conn.Pool.ConnMaxLifetime); err != nil {
				return fmt.Errorf("connections[%q].pool.conn_max_lifetime: %w", conn.Name, err)
			}
		}
	}
	for i, r := range c.Policy.Rules {
		if r.Name == "" {
			return fmt.Errorf("policy.rules[%d].name is required", i)
		}
		switch r.Action {
		case "allow", "deny", "require_approval":
		default:
			return fmt.Errorf("policy.rules[%q].action must be allow|deny|require_approval", r.Name)
		}
	}
	if c.Limits.QueryTimeout != "" {
		if _, err := time.ParseDuration(c.Limits.QueryTimeout); err != nil {
			return fmt.Errorf("limits.query_timeout: %w", err)
		}
	}
	if c.Approval.TokenTTL != "" {
		if _, err := time.ParseDuration(c.Approval.TokenTTL); err != nil {
			return fmt.Errorf("approval.token_ttl: %w", err)
		}
	}
	return nil
}

// QueryTimeout returns the parsed query timeout.
func (c *Config) QueryTimeout() time.Duration {
	d, _ := time.ParseDuration(c.Limits.QueryTimeout)
	return d
}

// ConnectionMetas converts connection configs to domain metadata.
func (c *Config) ConnectionMetas() map[string]domain.ConnectionMeta {
	metas := make(map[string]domain.ConnectionMeta, len(c.Connections))
	for _, conn := range c.Connections {
		metas[conn.Name] = domain.ConnectionMeta{
			Connection:     domain.Connection{Name: conn.Name, Driver: conn.Driver, ReadOnly: conn.ReadOnly},
			AllowedSchemas: conn.AllowedSchemas,
			DeniedTables:   conn.DeniedTables,
			AllowedTables:  conn.AllowedTables,
			RowLimit:       conn.RowLimit,
		}
	}
	return metas
}

// PolicyRules converts policy configs to domain rules.
func (c *Config) PolicyRules() []domain.PolicyRule {
	rules := make([]domain.PolicyRule, 0, len(c.Policy.Rules))
	for _, r := range c.Policy.Rules {
		rules = append(rules, domain.PolicyRule{
			Name:   r.Name,
			Match:  domain.RuleMatch{Statements: r.Match.Statement, Connections: r.Match.Connection, Tables: r.Match.Tables, Functions: r.Match.Functions},
			Action: domain.Action(r.Action),
			Reason: r.Reason,
		})
	}
	return rules
}

// AuthTokens converts token configs to domain tokens.
func (c *Config) AuthTokens() []domain.AuthToken {
	tokens := make([]domain.AuthToken, 0, len(c.Auth.Tokens))
	for _, t := range c.Auth.Tokens {
		scopes := make([]domain.Scope, 0, len(t.Scopes))
		for _, s := range t.Scopes {
			scopes = append(scopes, domain.Scope(s))
		}
		tokens = append(tokens, domain.AuthToken{ID: t.ID, Hash: t.Hash, Scopes: scopes, ExpiresAt: t.ExpiresAt})
	}
	return tokens
}

// RedactPatterns compiles redaction regexes.
func (c *Config) RedactPatterns() ([]domain.RedactPattern, error) {
	patterns := make([]domain.RedactPattern, 0, len(c.Redaction.Patterns))
	for _, p := range c.Redaction.Patterns {
		re, err := regexp.Compile(p.Regex)
		if err != nil {
			return nil, fmt.Errorf("redaction.patterns[%q].regex: %w", p.Name, err)
		}
		patterns = append(patterns, domain.RedactPattern{Name: p.Name, Regex: re, Replacement: p.Replacement})
	}
	return patterns, nil
}

// FileModeForScheme returns the expected file mode for file resolver refs,
// or 0 when no file resolver is configured.
func (c *Config) FileModeForScheme() uint32 {
	for _, r := range c.Credential.Resolvers {
		if r.Scheme == "file" && r.FileMode != 0 {
			return r.FileMode
		}
	}
	return 0
}

// FileBasePath returns the base path of the file resolver, if configured.
func (c *Config) FileBasePath() string {
	for _, r := range c.Credential.Resolvers {
		if r.Scheme == "file" {
			return r.BasePath
		}
	}
	return ""
}

// EnvPrefix returns the prefix of the env resolver, if configured.
func (c *Config) EnvPrefix() string {
	for _, r := range c.Credential.Resolvers {
		if r.Scheme == "env" {
			return r.Prefix
		}
	}
	return ""
}

// HasScheme reports whether a credential resolver scheme is configured.
func (c *Config) HasScheme(scheme string) bool {
	for _, r := range c.Credential.Resolvers {
		if strings.EqualFold(r.Scheme, scheme) {
			return true
		}
	}
	return len(c.Credential.Resolvers) == 0
}
