package cli

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/vietthanh1999/ohjanus/internal/adapter/in/mcp"
	"github.com/vietthanh1999/ohjanus/internal/adapter/in/transport/stdio"
	auditfile "github.com/vietthanh1999/ohjanus/internal/adapter/out/audit/file"
	auditmulti "github.com/vietthanh1999/ohjanus/internal/adapter/out/audit/multi"
	auditstderr "github.com/vietthanh1999/ohjanus/internal/adapter/out/audit/stderr"
	auditstdout "github.com/vietthanh1999/ohjanus/internal/adapter/out/audit/stdout"
	authconfig "github.com/vietthanh1999/ohjanus/internal/adapter/out/auth/config"
	systemclock "github.com/vietthanh1999/ohjanus/internal/adapter/out/clock/system"
	pgconnector "github.com/vietthanh1999/ohjanus/internal/adapter/out/connector/postgres"
	credchain "github.com/vietthanh1999/ohjanus/internal/adapter/out/credential/chain"
	credenv "github.com/vietthanh1999/ohjanus/internal/adapter/out/credential/env"
	credfile "github.com/vietthanh1999/ohjanus/internal/adapter/out/credential/file"
	policyyaml "github.com/vietthanh1999/ohjanus/internal/adapter/out/policy/yaml"
	redactregex "github.com/vietthanh1999/ohjanus/internal/adapter/out/redact/regex"
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
			if cfg.Server.Transport != "stdio" {
				return fmt.Errorf("transport %q lands in v0.2, only stdio is wired in v0.1", cfg.Server.Transport)
			}
			logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: parseLogLevel(cfg.Observability.LogLevel)}))

			clock := systemclock.Clock{}
			auditSink, err := buildAuditSink(cfg, logger)
			if err != nil {
				return err
			}
			defer auditSink.Close()

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

			tokenResolver := authconfig.New(cfg.AuthTokens(), clock)
			policyEng := policyyaml.New(cfg.PolicyRules(), domain.Action(cfg.Policy.DefaultAction))
			metas := cfg.ConnectionMetas()

			validator := buildValidator(metas)

			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			creds := credchain.New(credenv.New(), credfile.New(cfg.FileModeForScheme()))
			pools, err := openPools(ctx, cfg, creds)
			if err != nil {
				return err
			}
			defer closePools(pools)

			gateway := service.NewGateway(clock)
			readSvc := service.NewReadService(gateway, validator, policyEng, pools, metas, auditSink, clock, redactor,
				1000, cfg.Limits.MaxQueryLength, cfg.QueryTimeout())
			schemaSvc := service.NewSchemaService(gateway, pools, metas)
			writeSvc := service.NewWriteService(gateway)
			server := mcp.NewServer(cfg.Server.Name, cfg.Server.Version, cfg.Auth.Mode, tokenResolver, clock, readSvc, writeSvc, schemaSvc)

			auditSink.Emit(context.Background(), domain.AuditEvent{
				TS:        clock.Now(),
				Event:     "session.initialized",
				RequestID: domain.NewRequestID(),
				Status:    "success",
			})
			logger.Info("janus serving", "transport", "stdio", "connections", len(metas))

			return stdio.New(server, os.Stdin, os.Stdout, logger).Serve(ctx)
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
			secret, err := randomString(32)
			if err != nil {
				return err
			}
			secret = "jn_" + secret
			sum := sha256.Sum256([]byte(secret))
			hash := "sha256:" + hex.EncodeToString(sum[:])
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
			chain := credchain.New(credenv.New(), credfile.New(cfg.FileModeForScheme()))
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
