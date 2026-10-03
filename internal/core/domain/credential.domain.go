package domain

import (
	"regexp"
	"time"
)

// Credentials holds the resolved database credentials.
// It must never leave the process: never log it, never return it
// in an MCP response, never include it in error messages.
type Credentials struct {
	DSN      string
	Username string
	Password string
	Host     string
	Port     int
	Database string
}

// RedactedDSN returns the DSN with the password replaced by ***.
// Resolvers usually fill only DSN, so the password is parsed out of the
// DSN itself. Substring replacement is deliberately avoided: it mangles
// hosts when the password is a substring of them.
func (c Credentials) RedactedDSN() string {
	if c.DSN == "" {
		return ""
	}
	if out := urlDSNPasswordRe.ReplaceAllString(c.DSN, "$1***@"); out != c.DSN {
		return out
	}
	return bareDSNPasswordRe.ReplaceAllString(c.DSN, "$1***@")
}

// postgres://user:pass@host/db (pass itself may contain @, but not /).
var urlDSNPasswordRe = regexp.MustCompile(`(://[^/:@?]+:)[^/]+@`)

// Non-URL DSNs, e.g. MySQL user:pass@tcp(host)/db.
var bareDSNPasswordRe = regexp.MustCompile(`^([^:@/]+:)[^/]+@`)

// PoolConfig controls the per-alias connection pool.
type PoolConfig struct {
	MaxOpen         int
	MaxIdle         int
	ConnMaxLifetime time.Duration
}
