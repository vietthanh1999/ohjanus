package domain

import (
	"strings"
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
func (c Credentials) RedactedDSN() string {
	if c.DSN == "" || c.Password == "" {
		return c.DSN
	}
	return strings.Replace(c.DSN, c.Password, "***", 1)
}

// PoolConfig controls the per-alias connection pool.
type PoolConfig struct {
	MaxOpen         int
	MaxIdle         int
	ConnMaxLifetime time.Duration
}
