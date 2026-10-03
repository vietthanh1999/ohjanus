package mcp

import "github.com/vietthanh1999/ohjanus/internal/core/domain"

// Tool names served by Janus (§5.1 of the spec).
const (
	ToolListConnections = "db_list_connections"
	ToolSchema          = "db_schema"
	ToolRead            = "db_read"
	ToolWritePreview    = "db_write_preview"
	ToolWriteExecute    = "db_write_execute"
	ToolExplain         = "db_explain"
)

func schema(props map[string]any, required ...string) map[string]any {
	if required == nil {
		required = []string{}
	}
	return map[string]any{
		"type":                 "object",
		"properties":           props,
		"required":             required,
		"additionalProperties": false,
	}
}

// ToolDefinitions returns every tool with its JSON input schema.
func ToolDefinitions() []Tool {
	return []Tool{
		{
			Name:        ToolListConnections,
			Description: "List the database connection aliases the agent may see.",
			InputSchema: schema(map[string]any{}),
		},
		{
			Name:        ToolSchema,
			Description: "List schemas, tables and columns of a connection.",
			InputSchema: schema(map[string]any{
				"connection": map[string]any{"type": "string"},
				"schema":     map[string]any{"type": "string"},
				"table":      map[string]any{"type": "string"},
			}, "connection"),
		},
		{
			Name:        ToolRead,
			Description: "Run a read-only query (SELECT/WITH/EXPLAIN/SHOW).",
			InputSchema: schema(map[string]any{
				"connection": map[string]any{"type": "string"},
				"sql":        map[string]any{"type": "string"},
				"params":     map[string]any{"type": "array"},
				"limit":      map[string]any{"type": "integer", "minimum": 1, "maximum": 1000},
			}, "connection", "sql"),
		},
		{
			Name:        ToolWritePreview,
			Description: "Preview a write statement and get a single-use token. Never executes.",
			InputSchema: schema(map[string]any{
				"connection": map[string]any{"type": "string"},
				"sql":        map[string]any{"type": "string"},
				"params":     map[string]any{"type": "array"},
			}, "connection", "sql"),
		},
		{
			Name:        ToolWriteExecute,
			Description: "Execute a write statement with a single-use preview token.",
			InputSchema: schema(map[string]any{
				"connection":    map[string]any{"type": "string"},
				"sql":           map[string]any{"type": "string"},
				"params":        map[string]any{"type": "array"},
				"preview_token": map[string]any{"type": "string"},
			}, "connection", "sql", "params", "preview_token"),
		},
		{
			Name:        ToolExplain,
			Description: "Run EXPLAIN on a query without executing it.",
			InputSchema: schema(map[string]any{
				"connection": map[string]any{"type": "string"},
				"sql":        map[string]any{"type": "string"},
				"params":     map[string]any{"type": "array"},
			}, "connection", "sql"),
		},
	}
}

// RequiredScope returns the scope needed to call a tool.
func RequiredScope(tool string) domain.Scope {
	switch tool {
	case ToolWritePreview:
		return domain.ScopeWritePreview
	case ToolWriteExecute:
		return domain.ScopeWriteExecute
	default:
		return domain.ScopeRead
	}
}
