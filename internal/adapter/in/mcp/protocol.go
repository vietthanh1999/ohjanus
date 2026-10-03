package mcp

import "encoding/json"

// JSON-RPC 2.0 methods served by Janus (§4.2 of the spec).
const (
	MethodInitialize  = "initialize"
	MethodInitialized = "initialized"
	MethodToolsList   = "tools/list"
	MethodToolsCall   = "tools/call"
	MethodPing        = "ping"
)

// ProtocolVersion negotiated during initialize.
const ProtocolVersion = "2024-11-05"

// JSON-RPC error codes.
const (
	CodeParseError     = -32700
	CodeInvalidRequest = -32600
	CodeMethodNotFound = -32601
	CodeInvalidParams  = -32602
	CodeInternal       = -32603
	CodeJanusError     = -32000
)

// Request is a JSON-RPC 2.0 request. A missing ID means notification.
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  map[string]any  `json:"params,omitempty"`
}

// HasID reports whether the request expects a response.
func (r Request) HasID() bool {
	return len(r.ID) > 0 && string(r.ID) != "null"
}

// Response is a JSON-RPC 2.0 response.
type Response struct {
	JSONRPC string       `json:"jsonrpc"`
	ID      any          `json:"id,omitempty"`
	Result  any          `json:"result,omitempty"`
	Error   *ErrorObject `json:"error,omitempty"`
}

// ErrorObject is a JSON-RPC 2.0 error with Janus details in Data.
type ErrorObject struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    map[string]any `json:"data,omitempty"`
}

// Tool is an MCP tool definition returned by tools/list.
type Tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"inputSchema"`
}
