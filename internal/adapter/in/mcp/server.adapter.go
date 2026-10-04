package mcp

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/in"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// Server dispatches MCP methods to port/in use cases.
// It only depends on port interfaces, never on service implementations.
type Server struct {
	name     string
	version  string
	authMode string
	tokens   out.TokenResolver
	clock    out.Clock
	read     in.ReadUseCase
	write    in.WriteUseCase
	schema   in.SchemaUseCase
}

// NewServer wires an MCP server.
func NewServer(name, version, authMode string, tokens out.TokenResolver, clock out.Clock, read in.ReadUseCase, write in.WriteUseCase, schema in.SchemaUseCase) *Server {
	return &Server{
		name: name, version: version, authMode: authMode,
		tokens: tokens, clock: clock, read: read, write: write, schema: schema,
	}
}

// Handle processes one raw JSON-RPC message and returns the raw response,
// or nil for notifications (initialized, or requests without an id).
func (s *Server) Handle(ctx context.Context, raw []byte) []byte {
	reqID := domain.NewRequestID()
	ctx = domain.WithRequestID(ctx, reqID)
	var req Request
	if err := json.Unmarshal(raw, &req); err != nil {
		return marshal(Response{JSONRPC: "2.0", Error: &ErrorObject{Code: CodeParseError, Message: "parse error"}})
	}
	if req.JSONRPC != "2.0" || req.Method == "" {
		return s.respond(req, reqID, nil, &ErrorObject{Code: CodeInvalidRequest, Message: "invalid request"})
	}
	switch req.Method {
	case MethodInitialize:
		return s.respond(req, reqID, map[string]any{
			"protocolVersion": ProtocolVersion,
			"capabilities": map[string]any{
				"tools":   map[string]any{"listChanged": false},
				"logging": map[string]any{},
			},
			"serverInfo": map[string]any{"name": s.name, "version": s.version},
		}, nil)
	case MethodInitialized:
		return nil
	case MethodPing:
		return s.respond(req, reqID, map[string]any{}, nil)
	case MethodToolsList:
		authCtx, rpcErr := s.authenticate(ctx, req.Params)
		if rpcErr != nil {
			return s.respond(req, reqID, nil, rpcErr)
		}
		return s.respond(req, reqID, map[string]any{"tools": visibleTools(authCtx)}, nil)
	case MethodToolsCall:
		return s.handleCall(ctx, req, reqID)
	default:
		return s.respond(req, reqID, nil, &ErrorObject{Code: CodeMethodNotFound, Message: "method not found: " + req.Method})
	}
}

func (s *Server) respond(req Request, reqID string, result any, rpcErr *ErrorObject) []byte {
	if !req.HasID() {
		return nil
	}
	var id any
	_ = json.Unmarshal(req.ID, &id)
	if rpcErr != nil {
		if rpcErr.Data == nil {
			rpcErr.Data = map[string]any{}
		}
		rpcErr.Data["request_id"] = reqID
	}
	return marshal(Response{JSONRPC: "2.0", ID: id, Result: result, Error: rpcErr})
}

func marshal(v any) []byte {
	raw, err := json.Marshal(v)
	if err != nil {
		raw, _ = json.Marshal(Response{JSONRPC: "2.0", Error: &ErrorObject{Code: CodeInternal, Message: "INTERNAL: marshal response"}})
	}
	return raw
}

// authenticate verifies the MCP token.
// In auth mode "none" (dev only) every request runs as admin.
// Otherwise the token is read from params.token and verified by hash.
func (s *Server) authenticate(ctx context.Context, params map[string]any) (context.Context, *ErrorObject) {
	if s.authMode == "none" {
		return domain.WithAuth(ctx, "dev", []domain.Scope{domain.ScopeAdmin}), nil
	}
	tok, _ := params["token"].(string)
	if tok == "" {
		return nil, &ErrorObject{Code: CodeJanusError, Message: "UNAUTHENTICATED: missing token",
			Data: map[string]any{"code": string(domain.CodeUnauthenticated)}}
	}
	sum := sha256.Sum256([]byte(tok))
	t, err := s.tokens.LookupByHash(ctx, "sha256:"+hex.EncodeToString(sum[:]))
	if err != nil {
		code := domain.CodeTokenInvalid
		msg := err.Error()
		if de, ok := err.(*domain.Error); ok {
			code, msg = de.Code, de.Message
		}
		return nil, &ErrorObject{Code: CodeJanusError, Message: string(code) + ": " + msg,
			Data: map[string]any{"code": string(code)}}
	}
	return domain.WithAuth(ctx, t.ID, t.Scopes), nil
}

// visibleTools lists only the tools the token's scopes allow.
func visibleTools(ctx context.Context) []Tool {
	scopes := domain.ScopesFrom(ctx)
	has := func(s domain.Scope) bool {
		for _, sc := range scopes {
			if sc == s || sc == domain.ScopeAdmin {
				return true
			}
		}
		return false
	}
	all := ToolDefinitions()
	visible := make([]Tool, 0, len(all))
	for _, t := range all {
		if has(RequiredScope(t.Name)) {
			visible = append(visible, t)
		}
	}
	return visible
}

func (s *Server) handleCall(ctx context.Context, req Request, reqID string) []byte {
	authCtx, rpcErr := s.authenticate(ctx, req.Params)
	if rpcErr != nil {
		return s.respond(req, reqID, nil, rpcErr)
	}
	name, _ := req.Params["name"].(string)
	args, _ := req.Params["arguments"].(map[string]any)
	if args == nil {
		args = map[string]any{}
	}
	if err := domain.RequireScope(authCtx, RequiredScope(name)); err != nil {
		return s.respond(req, reqID, nil, janusError(reqID, err))
	}
	var result any
	var err error
	switch name {
	case ToolListConnections:
		result, err = s.listConnections(authCtx)
	case ToolSchema:
		result, err = s.getSchema(authCtx, args)
	case ToolRead:
		result, err = s.readQuery(authCtx, args)
	case ToolExplain:
		result, err = s.explain(authCtx, args)
	case ToolWritePreview:
		result, err = s.writePreview(authCtx, args)
	case ToolWriteExecute:
		result, err = s.writeExecute(authCtx, args)
	case "":
		return s.respond(req, reqID, nil, &ErrorObject{Code: CodeInvalidParams, Message: "missing tool name"})
	default:
		return s.respond(req, reqID, nil, &ErrorObject{Code: CodeInvalidParams, Message: "unknown tool: " + name})
	}
	if err != nil {
		return s.respond(req, reqID, nil, janusError(reqID, err))
	}
	return s.respond(req, reqID, map[string]any{
		"content": []any{map[string]any{"type": "text", "text": mustJSON(result)}},
	}, nil)
}

// janusError maps a domain error to the MCP error format (§12.2).
func janusError(reqID string, err error) *ErrorObject {
	if de, ok := err.(*domain.Error); ok {
		return &ErrorObject{
			Code:    CodeJanusError,
			Message: string(de.Code) + ": " + de.Message,
			Data: map[string]any{
				"code":       string(de.Code),
				"rule":       de.Rule,
				"request_id": reqID,
			},
		}
	}
	return &ErrorObject{Code: CodeInternal, Message: "INTERNAL: " + err.Error(),
		Data: map[string]any{"code": string(domain.CodeInternal), "request_id": reqID}}
}

func mustJSON(v any) string {
	raw, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(raw)
}

func strParam(args map[string]any, key string) string {
	s, _ := args[key].(string)
	return s
}

func paramsFrom(args map[string]any) []any {
	p, _ := args["params"].([]any)
	return p
}

func limitFrom(args map[string]any) int {
	f, _ := args["limit"].(float64)
	return int(f)
}

func (s *Server) listConnections(ctx context.Context) (any, error) {
	conns, err := s.schema.ListConnections(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]map[string]any, 0, len(conns))
	for _, c := range conns {
		items = append(items, map[string]any{"name": c.Name, "driver": c.Driver, "readonly": c.ReadOnly})
	}
	return map[string]any{"connections": items}, nil
}

func (s *Server) getSchema(ctx context.Context, args map[string]any) (any, error) {
	schemas, err := s.schema.GetSchema(ctx, strParam(args, "connection"), strParam(args, "schema"), strParam(args, "table"))
	if err != nil {
		return nil, err
	}
	items := make([]map[string]any, 0, len(schemas))
	for _, sc := range schemas {
		tables := make([]map[string]any, 0, len(sc.Tables))
		for _, t := range sc.Tables {
			cols := make([]map[string]any, 0, len(t.Columns))
			for _, c := range t.Columns {
				cols = append(cols, map[string]any{"name": c.Name, "type": c.Type, "nullable": c.Nullable})
			}
			tables = append(tables, map[string]any{
				"name": t.Name, "columns": cols, "primary_key": t.PrimaryKey,
			})
		}
		routines := make([]map[string]any, 0, len(sc.Routines))
		for _, r := range sc.Routines {
			routines = append(routines, map[string]any{"name": r.Name, "kind": r.Kind})
		}
		sequences := make([]map[string]any, 0, len(sc.Sequences))
		for _, sq := range sc.Sequences {
			sequences = append(sequences, map[string]any{"name": sq.Name})
		}
		items = append(items, map[string]any{"name": sc.Name, "tables": tables, "routines": routines, "sequences": sequences})
	}
	return map[string]any{"schemas": items}, nil
}

func readRequestFrom(args map[string]any) domain.ReadRequest {
	return domain.ReadRequest{
		Connection: strParam(args, "connection"),
		SQL:        strParam(args, "sql"),
		Params:     paramsFrom(args),
		Limit:      limitFrom(args),
	}
}

func (s *Server) readQuery(ctx context.Context, args map[string]any) (any, error) {
	res, err := s.read.Read(ctx, readRequestFrom(args))
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"columns": res.Columns, "rows": res.Rows, "row_count": res.RowCount,
		"truncated": res.Truncated, "duration_ms": res.DurationMs,
	}, nil
}

func (s *Server) explain(ctx context.Context, args map[string]any) (any, error) {
	plan, err := s.read.Explain(ctx, readRequestFrom(args))
	if err != nil {
		return nil, err
	}
	return map[string]any{"plan": plan.Text, "affected_estimate": plan.AffectedEstimate}, nil
}

func (s *Server) writePreview(ctx context.Context, args map[string]any) (any, error) {
	p, err := s.write.Preview(ctx, readRequestFrom(args))
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"preview_token": p.PreviewTokenID, "expires_at": p.ExpiresAt,
		"statement_type": string(p.StatementType), "affected_estimate": p.AffectedEstimate,
		"plan": p.Plan, "warnings": p.Warnings, "sql_normalized": p.SQLNormalized,
	}, nil
}

func (s *Server) writeExecute(ctx context.Context, args map[string]any) (any, error) {
	res, err := s.write.Execute(ctx, domain.ExecuteRequest{
		Connection:     strParam(args, "connection"),
		SQL:            strParam(args, "sql"),
		Params:         paramsFrom(args),
		PreviewTokenID: strParam(args, "preview_token"),
	})
	if err != nil {
		return nil, err
	}
	return map[string]any{"rows_affected": res.RowsAffected, "duration_ms": res.DurationMs}, nil
}
