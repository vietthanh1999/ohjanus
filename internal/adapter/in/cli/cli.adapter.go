package cli

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	adminapi "github.com/vietthanh1999/ohjanus/internal/adapter/in/admin"
	"github.com/vietthanh1999/ohjanus/internal/adapter/in/mcp"
	mcphttp "github.com/vietthanh1999/ohjanus/internal/adapter/in/transport/http"
	"github.com/vietthanh1999/ohjanus/internal/adapter/in/transport/stdio"
	cliapproval "github.com/vietthanh1999/ohjanus/internal/adapter/out/approval/cli"
	auditfile "github.com/vietthanh1999/ohjanus/internal/adapter/out/audit/file"
	auditmemory "github.com/vietthanh1999/ohjanus/internal/adapter/out/audit/memory"
	auditmulti "github.com/vietthanh1999/ohjanus/internal/adapter/out/audit/multi"
	auditstderr "github.com/vietthanh1999/ohjanus/internal/adapter/out/audit/stderr"
	auditstdout "github.com/vietthanh1999/ohjanus/internal/adapter/out/audit/stdout"
	authmemory "github.com/vietthanh1999/ohjanus/internal/adapter/out/auth/memory"
	systemclock "github.com/vietthanh1999/ohjanus/internal/adapter/out/clock/system"
	pgconnector "github.com/vietthanh1999/ohjanus/internal/adapter/out/connector/postgres"
	credchain "github.com/vietthanh1999/ohjanus/internal/adapter/out/credential/chain"
	credenv "github.com/vietthanh1999/ohjanus/internal/adapter/out/credential/env"
	credfile "github.com/vietthanh1999/ohjanus/internal/adapter/out/credential/file"
	credkeyring "github.com/vietthanh1999/ohjanus/internal/adapter/out/credential/keyring"
	vaultresolver "github.com/vietthanh1999/ohjanus/internal/adapter/out/credential/vault"
	limitmemory "github.com/vietthanh1999/ohjanus/internal/adapter/out/limits/memory"
	prommetrics "github.com/vietthanh1999/ohjanus/internal/adapter/out/metrics/prom"
	policyyaml "github.com/vietthanh1999/ohjanus/internal/adapter/out/policy/yaml"
	redactregex "github.com/vietthanh1999/ohjanus/internal/adapter/out/redact/regex"
	tokememory "github.com/vietthanh1999/ohjanus/internal/adapter/out/token/memory"
	denyall "github.com/vietthanh1999/ohjanus/internal/adapter/out/validator/denyall"
	pgvalidator "github.com/vietthanh1999/ohjanus/internal/adapter/out/validator/postgres"
	valrouter "github.com/vietthanh1999/ohjanus/internal/adapter/out/validator/router"
	"github.com/vietthanh1999/ohjanus/internal/config"
	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
	"github.com/vietthanh1999/ohjanus/internal/core/service"
)

// Version is overridden at build time: -ldflags "-X .../cli.Version=x.y.z".
var Version = "0.1.0"

// NewRootCommand builds the janus CLI (§14 of the spec).
func NewRootCommand() *cobra.Command {
	var cfgPath, logLevel, transport string
	root := &cobra.Command{
		Use:   "janus",
		Short: "Janus secure MCP gateway for databases",
	}
	root.PersistentFlags().StringVarP(&cfgPath, "config", "c", "./janus.yaml", "config file path")
	root.PersistentFlags().StringVar(&logLevel, "log-level", "", "debug|info|warn|error (overrides config)")
	root.PersistentFlags().StringVar(&transport, "transport", "", "stdio|http|sse (overrides config)")
	root.AddCommand(
		newServeCmd(&cfgPath, &logLevel, &transport),
		newValidateCmd(&cfgPath),
		newTokenCmd(&cfgPath),
		newConnectionCmd(&cfgPath),
		newVersionCmd(),
	)
	return root
}

func loadConfig(cfgPath string, logLevel, transport *string) (*config.Config, error) {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return nil, err
	}
	if *transport != "" {
		cfg.Server.Transport = *transport
	}
	if *logLevel != "" {
		cfg.Observability.LogLevel = *logLevel
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func parseLogLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func newServeCmd(cfgPath, logLevel, transport *string) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Run the gateway (default transport: stdio)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(*cfgPath, logLevel, transport)
			if err != nil {
				return err
			}
			if cfg.Server.Transport != "stdio" && cfg.Server.Transport != "http" && cfg.Server.Transport != "sse" {
				return fmt.Errorf("transport %q must be stdio|http|sse", cfg.Server.Transport)
			}
			logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: parseLogLevel(cfg.Observability.LogLevel)}))

			clock := systemclock.Clock{}
			auditSink, err := buildAuditSink(cfg, logger)
			if err != nil {
				return err
			}
			defer auditSink.Close()
			// Queryable ring buffer for the Admin UI; the configured sinks
			// remain the durable trail.
			auditBuffer := auditmemory.New(10000)
			combinedAudit := auditmulti.New(auditSink, auditBuffer)

			patterns, err := cfg.RedactPatterns()
			if err != nil {
				return err
			}
			var columns []string
			if cfg.Redaction.Enabled {
				columns = cfg.Redaction.Columns
			} else {
				patterns = nil
			}
			redactor := redactregex.New(patterns, columns)

			tokenResolver := authmemory.New(clock)
			tokenResolver.Seed(cfg.AuthTokens())
			authStore := tokenResolver
			policyEng := policyyaml.New(cfg.PolicyRules(), domain.Action(cfg.Policy.DefaultAction))
			metas := cfg.ConnectionMetas()

			validator := buildValidator(metas)

			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			creds, err := buildCredentialChain(ctx, cfg)
			if err != nil {
				return err
			}
			pools, err := openPools(ctx, cfg, creds)
			if err != nil {
				return err
			}
			defer closePools(pools)

			tokenTTL, err := time.ParseDuration(cfg.Approval.TokenTTL)
			if err != nil {
				return fmt.Errorf("approval.token_ttl: %w", err)
			}
			tokenStore := tokememory.New(clock, tokenTTL)
			approvalEng, err := buildApprovalEngine(cfg)
			if err != nil {
				return err
			}

			var limiter out.RateLimiter = out.AllowRateLimiter{}
			if cfg.Limits.RateLimit.Enabled {
				limiter = limitmemory.New(cfg.Limits.RateLimit.RequestsPerMinute, clock)
			}
			var metrics out.Metrics = out.NoopMetrics{}
			var metricsHandler http.Handler
			if cfg.Observability.Metrics.Enabled {
				pm := prommetrics.New()
				metrics = pm
				metricsHandler = pm.Handler()
			}

			gateway := service.NewGatewayWithLimits(clock, limiter, cfg.Limits.MaxConcurrentQueries, metrics)
			readSvc := service.NewReadService(gateway, validator, policyEng, pools, metas, combinedAudit, clock, redactor,
				1000, cfg.Limits.MaxQueryLength, cfg.QueryTimeout())
			schemaSvc := service.NewSchemaService(gateway, pools, metas)
			writeSvc := service.NewWriteService(gateway, validator, policyEng, pools, tokenStore, approvalEng,
				combinedAudit, clock, redactor, tokenTTL, cfg.QueryTimeout(), cfg.Limits.MaxQueryLength)
			server := mcp.NewServer(cfg.Server.Name, cfg.Server.Version, cfg.Auth.Mode, tokenResolver, clock, readSvc, writeSvc, schemaSvc)

			auditSink.Emit(context.Background(), domain.AuditEvent{
				TS:        clock.Now(),
				Event:     "session.initialized",
				RequestID: domain.NewRequestID(),
				Status:    "success",
			})

			if cfg.Admin.Enabled {
				if err := startAdmin(ctx, cfg, logger, tokenStore, tokenResolver, authStore, auditBuffer, combinedAudit, pools, metas, clock, metricsHandler); err != nil {
					return err
				}
			}
			logger.Info("janus serving", "transport", cfg.Server.Transport, "connections", len(metas))

			if cfg.Server.Transport == "stdio" {
				return stdio.New(server, os.Stdin, os.Stdout, logger).Serve(ctx)
			}
			return serveMCPHTTP(ctx, cfg, logger, server, tokenStore)
		},
	}
}

func buildValidator(metas map[string]domain.ConnectionMeta) out.Validator {
	rules := make(map[string]pgvalidator.ConnRules, len(metas))
	byConn := map[string]out.Validator{}
	for name, m := range metas {
		rules[name] = pgvalidator.ConnRules{
			AllowedSchemas: m.AllowedSchemas,
			DeniedTables:   m.DeniedTables,
			AllowedTables:  m.AllowedTables,
		}
	}
	pgVal := pgvalidator.New(rules, nil)
	for name, m := range metas {
		if m.Connection.Driver == "postgres" {
			byConn[name] = pgVal
		}
	}
	return valrouter.New(byConn, denyall.New())
}

func openPools(ctx context.Context, cfg *config.Config, creds *credchain.Chain) (map[string]out.Pool, error) {
	pools := map[string]out.Pool{}
	for _, conn := range cfg.Connections {
		resolved, err := creds.Resolve(ctx, conn.DSNRef)
		if err != nil {
			closePools(pools)
			return nil, fmt.Errorf("connection %q: %w", conn.Name, err)
		}
		pc := domain.PoolConfig{MaxOpen: conn.Pool.MaxOpen, MaxIdle: conn.Pool.MaxIdle}
		if conn.Pool.ConnMaxLifetime != "" {
			d, err := time.ParseDuration(conn.Pool.ConnMaxLifetime)
			if err != nil {
				closePools(pools)
				return nil, fmt.Errorf("connection %q: %w", conn.Name, err)
			}
			pc.ConnMaxLifetime = d
		}
		var pool out.Pool
		switch conn.Driver {
		case "postgres":
			pool, err = pgconnector.New().Open(ctx, resolved, pc)
		default:
			err = fmt.Errorf("driver %q not implemented in v0.1", conn.Driver)
		}
		if err != nil {
			closePools(pools)
			return nil, fmt.Errorf("connection %q: %w", conn.Name, err)
		}
		pools[conn.Name] = pool
	}
	return pools, nil
}

func closePools(pools map[string]out.Pool) {
	for _, p := range pools {
		_ = p.Close()
	}
}

// buildCredentialChain wires env + file always, keyring always (used only
// when a ref matches), and Vault when configured. Vault AppRole credentials
// resolve through the env/file chain to avoid recursion.
func buildCredentialChain(ctx context.Context, cfg *config.Config) (*credchain.Chain, error) {
	envR := credenv.New()
	fileR := credfile.New(cfg.FileModeForScheme())
	base := credchain.New(envR, fileR)
	resolvers := []out.CredentialResolver{envR, fileR, credkeyring.New(cfg.CredentialKeyringService())}
	for _, rc := range cfg.Credential.Resolvers {
		if !strings.EqualFold(rc.Scheme, "vault") {
			continue
		}
		roleID, err := resolveValue(ctx, base, rc.RoleIDRef)
		if err != nil {
			return nil, fmt.Errorf("vault role_id: %w", err)
		}
		secretID, err := resolveValue(ctx, base, rc.SecretIDRef)
		if err != nil {
			return nil, fmt.Errorf("vault secret_id: %w", err)
		}
		v, err := vaultresolver.New(rc.Address, roleID, secretID)
		if err != nil {
			return nil, fmt.Errorf("vault resolver: %w", err)
		}
		resolvers = append(resolvers, v)
	}
	return credchain.New(resolvers...), nil
}

// resolveValue resolves an env:/file: ref, or returns the literal.
func resolveValue(ctx context.Context, chain *credchain.Chain, ref string) (string, error) {
	if ref == "" {
		return "", fmt.Errorf("empty ref")
	}
	if strings.Contains(ref, "://") || strings.HasPrefix(ref, "env:") || strings.HasPrefix(ref, "file:") {
		creds, err := chain.Resolve(ctx, ref)
		if err != nil {
			return "", err
		}
		return creds.DSN, nil
	}
	return ref, nil
}

// serveMCPHTTP pre-binds the MCP HTTP port (fail fast) then serves
// POST /mcp + GET /mcp/sse until ctx ends. Both "http" and "sse"
// transports land here; sse is the stream-oriented alias.
func serveMCPHTTP(ctx context.Context, cfg *config.Config, logger *slog.Logger, server *mcp.Server, tokenStore out.TokenStore) error {
	if cfg.Server.HTTP.Listen == "" {
		return fmt.Errorf("server.http.listen is required for http/sse transport")
	}
	ln, err := net.Listen("tcp", cfg.Server.HTTP.Listen)
	if err != nil {
		return fmt.Errorf("mcp http listen %s: %w", cfg.Server.HTTP.Listen, err)
	}
	transport := mcphttp.New(server, tokenStore, cfg.Auth.Mode, mcphttp.TLSConfig{
		Enabled: cfg.Server.HTTP.TLS.Enabled, CertFile: cfg.Server.HTTP.TLS.CertFile, KeyFile: cfg.Server.HTTP.TLS.KeyFile,
	}, logger)
	logger.Info("mcp http serving", "listen", cfg.Server.HTTP.Listen, "tls", cfg.Server.HTTP.TLS.Enabled)
	return transport.ServeListener(ctx, ln)
}

// startAdmin pre-binds the Admin port (fail fast) then serves it in the
// background next to the MCP transport.
func startAdmin(ctx context.Context, cfg *config.Config, logger *slog.Logger, tokenStore out.TokenStore, tokenResolver out.TokenResolver, authStore out.AuthTokenStore, auditReader out.AuditReader, auditSink out.AuditSink, pools map[string]out.Pool, metas map[string]domain.ConnectionMeta, clock out.Clock, metricsHandler http.Handler) error {
	ln, err := net.Listen("tcp", cfg.Admin.Listen)
	if err != nil {
		return fmt.Errorf("admin listen %s: %w", cfg.Admin.Listen, err)
	}
	admin := adminapi.New(cfg.Admin.Listen, cfg.Auth.Mode, tokenStore, tokenResolver, authStore, auditReader, auditSink, pools, metas, clock)
	admin.SetMetricsHandler(metricsHandler)
	go func() {
		if err := admin.ServeListener(ctx, ln); err != nil {
			logger.Error("admin server stopped", "err", err)
		}
	}()
	logger.Info("admin api serving", "listen", cfg.Admin.Listen)
	return nil
}

func buildApprovalEngine(cfg *config.Config) (out.ApprovalEngine, error) {
	switch cfg.Approval.Method {
	case "", "cli":
		return cliapproval.NewApprovalEngine(), nil
	default:
		return nil, fmt.Errorf("approval.method %q not implemented in v0.2 (use cli)", cfg.Approval.Method)
	}
}

func buildAuditSink(cfg *config.Config, logger *slog.Logger) (out.AuditSink, error) {
	if !cfg.Audit.Enabled {
		return auditmulti.New(), nil
	}
	var sinks []out.AuditSink
	for _, s := range cfg.Audit.Sinks {
		switch s.Type {
		case "stdout":
			if cfg.Server.Transport == "stdio" {
				logger.Warn("audit stdout sink redirected to stderr to protect the stdio stream")
				sinks = append(sinks, auditstderr.New())
				continue
			}
			sinks = append(sinks, auditstdout.New())
		case "stderr":
			sinks = append(sinks, auditstderr.New())
		case "file":
			if s.Path == "" {
				return nil, fmt.Errorf("audit file sink requires path")
			}
			f, err := auditfile.New(s.Path)
			if err != nil {
				return nil, err
			}
			sinks = append(sinks, f)
		default:
			return nil, fmt.Errorf("unknown audit sink type %q", s.Type)
		}
	}
	if len(sinks) == 0 {
		sinks = append(sinks, auditstdout.New())
	}
	return auditmulti.New(sinks...), nil
}

func newValidateCmd(cfgPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate the config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := config.Load(*cfgPath); err != nil {
				return err
			}
			fmt.Println("config valid")
			return nil
		},
	}
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("janus", Version)
		},
	}
}

const tokenAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func randomString(n int) (string, error) {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	out := make([]byte, n)
	for i := range out {
		out[i] = tokenAlphabet[int(b[i%len(b)])%len(tokenAlphabet)]
	}
	return string(out), nil
}

func newTokenCmd(cfgPath *string) *cobra.Command {
	cmd := &cobra.Command{Use: "token", Short: "Manage MCP tokens"}
	var name, scopes, ttl string

	create := &cobra.Command{
		Use:   "create",
		Short: "Create a token (prints the secret once)",
		RunE: func(cmd *cobra.Command, args []string) error {
			secret, err := authmemory.GenerateSecret()
			if err != nil {
				return err
			}
			hash := authmemory.HashSecret(secret)
			id := name
			if id == "" {
				r, err := randomString(8)
				if err != nil {
					return err
				}
				id = "tok_" + r
			}
			expires := "never"
			expiresYAML := ""
			if ttl != "" {
				d, err := time.ParseDuration(ttl)
				if err != nil {
					return fmt.Errorf("bad --ttl: %w", err)
				}
				at := time.Now().UTC().Add(d)
				expires = at.Format(time.RFC3339)
				expiresYAML = fmt.Sprintf("\n      expires_at: %q", expires)
			}
			fmt.Printf("id: %s\ntoken: %s\n\nAdd to config (tokens are stored hashed):\n", id, secret)
			fmt.Printf("    - id: %q\n      hash: %q\n      scopes: [%s]%s\n",
				id, hash, quotedScopes(scopes), expiresYAML)
			return nil
		},
	}
	create.Flags().StringVar(&name, "name", "", "token id (default: random tok_*)")
	create.Flags().StringVar(&scopes, "scopes", "read", "comma-separated scopes")
	create.Flags().StringVar(&ttl, "ttl", "720h", "time to live (empty = never expires)")

	list := &cobra.Command{
		Use:   "list",
		Short: "List configured tokens (metadata only)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(*cfgPath)
			if err != nil {
				return err
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tSCOPES\tEXPIRES")
			for _, t := range cfg.Auth.Tokens {
				exp := "never"
				if !t.ExpiresAt.IsZero() {
					exp = t.ExpiresAt.Format("2006-01-02")
				}
				fmt.Fprintf(w, "%s\t%s\t%s\n", t.ID, strings.Join(t.Scopes, ","), exp)
			}
			return w.Flush()
		},
	}
	cmd.AddCommand(create, list)
	return cmd
}

func quotedScopes(scopes string) string {
	parts := strings.Split(scopes, ",")
	for i := range parts {
		parts[i] = fmt.Sprintf("%q", strings.TrimSpace(parts[i]))
	}
	return strings.Join(parts, ", ")
}

func newConnectionCmd(cfgPath *string) *cobra.Command {
	cmd := &cobra.Command{Use: "connection", Short: "Manage connections"}

	list := &cobra.Command{
		Use:   "list",
		Short: "List configured connections",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(*cfgPath)
			if err != nil {
				return err
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tDRIVER\tREADONLY\tDSN_REF")
			for _, c := range cfg.Connections {
				fmt.Fprintf(w, "%s\t%s\t%v\t%s\n", c.Name, c.Driver, c.ReadOnly, c.DSNRef)
			}
			return w.Flush()
		},
	}
	test := &cobra.Command{
		Use:   "test <name>",
		Short: "Resolve credentials and ping the database",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(*cfgPath)
			if err != nil {
				return err
			}
			var found *config.ConnectionConfig
			for i := range cfg.Connections {
				if cfg.Connections[i].Name == args[0] {
					found = &cfg.Connections[i]
					break
				}
			}
			if found == nil {
				return fmt.Errorf("unknown connection %q", args[0])
			}
			ctx := context.Background()
			chain, err := buildCredentialChain(ctx, cfg)
			if err != nil {
				return err
			}
			creds, err := chain.Resolve(ctx, found.DSNRef)
			if err != nil {
				return fmt.Errorf("resolve %s: %w", found.DSNRef, err)
			}
			if found.Driver != "postgres" {
				fmt.Printf("credentials resolve OK (dsn=%s); driver %q pool is not implemented yet\n",
					creds.RedactedDSN(), found.Driver)
				return nil
			}
			pc := domain.PoolConfig{MaxOpen: 1, MaxIdle: 1}
			start := time.Now()
			pool, err := pgconnector.New().Open(ctx, creds, pc)
			if err != nil {
				return err
			}
			defer pool.Close()
			fmt.Printf("connection %q healthy (driver=postgres, dsn=%s, latency=%s)\n",
				found.Name, creds.RedactedDSN(), time.Since(start).Round(time.Millisecond))
			return nil
		},
	}
	cmd.AddCommand(list, test)
	return cmd
}
