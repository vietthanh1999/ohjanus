package postgres

import (
	"context"
	"fmt"
	"strings"
	"sync"

	pgquery "github.com/pganalyze/pg_query_go/v6"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// BlockedFunctions is the default deny list (§6.2.3 of the spec).
var BlockedFunctions = []string{
	"pg_read_file", "pg_read_binary_file", "pg_ls_dir",
	"lo_import", "lo_export", "pg_sleep",
	"dblink", "dblink_exec",
	"pg_execute_server_program",
}

// ConnRules carries the allow/deny lists for one connection alias.
type ConnRules struct {
	AllowedSchemas []string
	DeniedTables   []string
	AllowedTables  []string
}

// Validator parses PostgreSQL with libpg_query and enforces table,
// schema and function rules on the AST. Never regex.
type Validator struct {
	mu        sync.RWMutex
	conns     map[string]ConnRules
	blockedFn map[string]struct{}
}

var _ out.Validator = (*Validator)(nil)

// New builds a Postgres validator. blockedFunctions extends the default
// block list; pass nil for defaults only.
func New(conns map[string]ConnRules, blockedFunctions []string) *Validator {
	blocked := make(map[string]struct{}, len(BlockedFunctions)+len(blockedFunctions))
	for _, f := range BlockedFunctions {
		blocked[f] = struct{}{}
	}
	for _, f := range blockedFunctions {
		blocked[strings.ToLower(f)] = struct{}{}
	}
	if conns == nil {
		conns = map[string]ConnRules{}
	}
	return &Validator{conns: conns, blockedFn: blocked}
}

// SetRules registers the allow/deny lists for a connection at runtime.
func (v *Validator) SetRules(name string, rules ConnRules) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.conns[name] = rules
}

// Validate runs the §6.2.2 pipeline: parse → classify → walk → checks.
func (v *Validator) Validate(_ context.Context, connection, sql string) (*domain.ValidatedQuery, error) {
	v.mu.RLock()
	rules, ok := v.conns[connection]
	v.mu.RUnlock()
	if !ok {
		return nil, domain.ErrConnectionNotFound(connection)
	}
	tree, err := pgquery.Parse(sql)
	if err != nil {
		return nil, &domain.Error{Code: domain.CodeParseError, Message: "SQL parse error: " + err.Error(), Rule: "validator:parse"}
	}
	if len(tree.Stmts) != 1 {
		return nil, &domain.Error{Code: domain.CodeQueryDenied, Message: "exactly one statement is allowed", Rule: "validator:single-statement"}
	}
	top := tree.Stmts[0].GetStmt()
	if top == nil {
		return nil, &domain.Error{Code: domain.CodeParseError, Message: "empty statement", Rule: "validator:parse"}
	}
	stmtType, err := classifyTopLevel(top)
	if err != nil {
		return nil, err
	}

	found := collect(tree)
	if found.danger != "" {
		return nil, &domain.Error{Code: domain.CodeQueryDenied, Message: found.danger, Rule: "validator:dangerous-node"}
	}
	// Data-modifying CTEs: a top-level SELECT hiding a write becomes write.
	if stmtType.IsRead() && found.write != "" {
		stmtType = domain.StatementType(found.write)
	}

	for _, fn := range found.functions {
		if _, blocked := v.blockedFn[fn]; blocked {
			return nil, &domain.Error{
				Code:    domain.CodeQueryDenied,
				Message: fmt.Sprintf("function %q is blocked", fn),
				Rule:    "validator:blocked-function",
			}
		}
	}
	for _, sch := range found.schemas {
		if len(rules.AllowedSchemas) > 0 && !containsFold(rules.AllowedSchemas, sch) {
			return nil, &domain.Error{
				Code:    domain.CodeSchemaNotAllowed,
				Message: fmt.Sprintf("schema %q is not allowed", sch),
				Rule:    "validator:allowed-schemas",
			}
		}
	}
	for _, tbl := range found.tables {
		if isDeniedTable(rules, tbl) {
			return nil, &domain.Error{
				Code:    domain.CodeTableNotAllowed,
				Message: fmt.Sprintf("table %q is not allowed", tbl),
				Rule:    "validator:denied-tables",
			}
		}
	}

	normalized, err := pgquery.Deparse(tree)
	if err != nil {
		normalized = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(sql), ";"))
	}
	return &domain.ValidatedQuery{
		Query:         domain.Query{Connection: connection, SQL: sql},
		StatementType: stmtType,
		Tables:        found.tables,
		Schemas:       found.schemas,
		Functions:     found.functions,
		NormalizedSQL: normalized,
	}, nil
}

// classifyTopLevel maps the top-level AST node to a statement type.
// Anything outside the read/write paths is denied immediately.
func classifyTopLevel(top *pgquery.Node) (domain.StatementType, error) {
	switch {
	case top.GetSelectStmt() != nil:
		return domain.StatementSelect, nil
	case top.GetInsertStmt() != nil:
		return domain.StatementInsert, nil
	case top.GetUpdateStmt() != nil:
		return domain.StatementUpdate, nil
	case top.GetDeleteStmt() != nil:
		return domain.StatementDelete, nil
	case top.GetMergeStmt() != nil:
		return domain.StatementMerge, nil
	case top.GetVariableShowStmt() != nil:
		return domain.StatementShow, nil
	case top.GetExplainStmt() != nil:
		if explainAnalyzes(top.GetExplainStmt()) {
			return "", &domain.Error{Code: domain.CodeQueryDenied, Message: "EXPLAIN ANALYZE executes the query and is blocked", Rule: "validator:explain-analyze"}
		}
		return domain.StatementExplain, nil
	default:
		return "", &domain.Error{Code: domain.CodeQueryDenied, Message: "statement type is not allowed", Rule: "validator:statement-type"}
	}
}

// explainAnalyzes reports whether EXPLAIN carries the ANALYZE option
// (which would execute the query).
func explainAnalyzes(stmt *pgquery.ExplainStmt) bool {
	for _, opt := range stmt.GetOptions() {
		if def := opt.GetDefElem(); def != nil {
			if strings.EqualFold(def.GetDefname(), "analyze") {
				return true
			}
		}
	}
	return false
}

type foundNodes struct {
	tables    []string
	schemas   []string
	functions []string
	write     string // nested INSERT/UPDATE/DELETE/MERGE found under a read
	danger    string // describe the first dangerous nested node
}

// collect walks the whole parse tree via protobuf reflection, so no nesting
// (subquery, CTE, function arg, UNION branch) can hide a table, a function
// call, or a dangerous utility node.
func collect(tree *pgquery.ParseResult) foundNodes {
	var f foundNodes
	seenTables := map[string]struct{}{}
	seenSchemas := map[string]struct{}{}
	seenFunctions := map[string]struct{}{}
	for _, stmt := range tree.Stmts {
		walkMessage(stmt.ProtoReflect(), func(m protoreflect.Message) {
			switch n := m.Interface().(type) {
			case *pgquery.RangeVar:
				if n.GetRelname() == "" {
					return
				}
				name := n.GetRelname()
				if sch := n.GetSchemaname(); sch != "" {
					name = sch + "." + name
					if _, ok := seenSchemas[sch]; !ok {
						seenSchemas[sch] = struct{}{}
						f.schemas = append(f.schemas, sch)
					}
				}
				if _, ok := seenTables[name]; !ok {
					seenTables[name] = struct{}{}
					f.tables = append(f.tables, name)
				}
			case *pgquery.FuncCall:
				if name := funcName(n); name != "" {
					if _, ok := seenFunctions[name]; !ok {
						seenFunctions[name] = struct{}{}
						f.functions = append(f.functions, name)
					}
				}
			case *pgquery.InsertStmt:
				f.write = firstWrite(f.write, "INSERT")
			case *pgquery.UpdateStmt:
				f.write = firstWrite(f.write, "UPDATE")
			case *pgquery.DeleteStmt:
				f.write = firstWrite(f.write, "DELETE")
			case *pgquery.MergeStmt:
				f.write = firstWrite(f.write, "MERGE")
			case *pgquery.CopyStmt:
				// Any file/program COPY touches the server filesystem.
				f.danger = "COPY with file or program is blocked"
			case *pgquery.DoStmt:
				f.danger = "DO blocks are blocked"
			case *pgquery.TransactionStmt:
				f.danger = "transaction control statements are blocked"
			case *pgquery.ListenStmt:
				f.danger = "LISTEN is blocked"
			case *pgquery.NotifyStmt:
				f.danger = "NOTIFY is blocked"
			case *pgquery.ExecuteStmt:
				f.danger = "EXECUTE (prepared statements) is blocked"
			case *pgquery.PrepareStmt:
				f.danger = "PREPARE is blocked"
			}
		})
	}
	return f
}

func firstWrite(current, next string) string {
	if current != "" {
		return current
	}
	return next
}

// funcName returns the lower-cased unqualified function name.
func funcName(call *pgquery.FuncCall) string {
	parts := call.GetFuncname()
	if len(parts) == 0 {
		return ""
	}
	last := parts[len(parts)-1].GetString_()
	if last == nil {
		return ""
	}
	return strings.ToLower(last.GetSval())
}

// walkMessage visits every populated message in the tree.
func walkMessage(m protoreflect.Message, visit func(protoreflect.Message)) {
	visit(m)
	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		switch {
		case fd.IsMap():
			// No maps carry SQL nodes.
		case fd.IsList():
			if fd.Message() == nil {
				return true
			}
			l := v.List()
			for i := 0; i < l.Len(); i++ {
				walkMessage(l.Get(i).Message(), visit)
			}
		default:
			if fd.Message() != nil {
				walkMessage(v.Message(), visit)
			}
		}
		return true
	})
}

func isDeniedTable(rules ConnRules, table string) bool {
	base := table
	if i := strings.LastIndex(table, "."); i >= 0 {
		base = table[i+1:]
	}
	for _, d := range rules.DeniedTables {
		if strings.EqualFold(d, table) || strings.EqualFold(d, base) {
			return true
		}
	}
	if len(rules.AllowedTables) == 0 {
		return false
	}
	for _, a := range rules.AllowedTables {
		if strings.EqualFold(a, table) || strings.EqualFold(a, base) {
			return false
		}
	}
	return true
}

func containsFold(list []string, s string) bool {
	for _, v := range list {
		if strings.EqualFold(v, s) {
			return true
		}
	}
	return false
}
