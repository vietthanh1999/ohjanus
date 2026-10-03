package yaml

import (
	"context"
	"testing"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
)

func vq(stmt domain.StatementType, conn string, tables ...string) *domain.ValidatedQuery {
	return &domain.ValidatedQuery{
		Query:         domain.Query{Connection: conn},
		StatementType: stmt,
		Tables:        tables,
	}
}

func TestFirstMatchWins(t *testing.T) {
	e := New([]domain.PolicyRule{
		{Name: "allow-select", Match: domain.RuleMatch{Statements: []string{"SELECT"}}, Action: domain.ActionAllow},
		{Name: "deny-all", Action: domain.ActionDeny, Reason: "blocked"},
	}, domain.ActionDeny)
	d := e.Evaluate(context.Background(), vq(domain.StatementSelect, "a", "orders"))
	if d.Action != domain.ActionAllow || d.Rule != "allow-select" {
		t.Errorf("got %+v, want allow via allow-select", d)
	}
	d = e.Evaluate(context.Background(), vq(domain.StatementDrop, "a", "orders"))
	if d.Action != domain.ActionDeny || d.Rule != "deny-all" {
		t.Errorf("got %+v, want deny via deny-all", d)
	}
}

func TestDefaultAction(t *testing.T) {
	e := New(nil, domain.ActionDeny)
	d := e.Evaluate(context.Background(), vq(domain.StatementSelect, "a"))
	if d.Action != domain.ActionDeny || d.Rule != "default" {
		t.Errorf("got %+v, want default deny", d)
	}
}

func TestTableGlobAndConnection(t *testing.T) {
	e := New([]domain.PolicyRule{
		{
			Name:   "deny-secret",
			Match:  domain.RuleMatch{Tables: []string{"secret*", "*.keys"}},
			Action: domain.ActionDeny,
		},
		{
			Name:   "analytics-only",
			Match:  domain.RuleMatch{Connections: []string{"analytics"}},
			Action: domain.ActionAllow,
		},
	}, domain.ActionDeny)

	if d := e.Evaluate(context.Background(), vq(domain.StatementSelect, "a", "secret_tokens")); d.Action != domain.ActionDeny {
		t.Errorf("glob secret* should deny, got %+v", d)
	}
	if d := e.Evaluate(context.Background(), vq(domain.StatementSelect, "a", "public.keys")); d.Action != domain.ActionDeny {
		t.Errorf("glob *.keys should deny, got %+v", d)
	}
	if d := e.Evaluate(context.Background(), vq(domain.StatementSelect, "analytics", "orders")); d.Action != domain.ActionAllow {
		t.Errorf("analytics should allow, got %+v", d)
	}
	if d := e.Evaluate(context.Background(), vq(domain.StatementSelect, "other", "orders")); d.Action != domain.ActionDeny {
		t.Errorf("other connection should fall to default deny, got %+v", d)
	}
}

func TestRequireApproval(t *testing.T) {
	e := New([]domain.PolicyRule{
		{Name: "warn-write", Match: domain.RuleMatch{Statements: []string{"UPDATE"}}, Action: domain.ActionRequireApproval},
	}, domain.ActionDeny)
	d := e.Evaluate(context.Background(), vq(domain.StatementUpdate, "a", "orders"))
	if d.Action != domain.ActionRequireApproval {
		t.Errorf("got %+v, want require_approval", d)
	}
}
