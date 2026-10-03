package yaml

import (
	"context"
	"path"
	"strings"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// Engine evaluates ordered rules; the first match wins.
type Engine struct {
	rules []domain.PolicyRule
	def   domain.Action
}

var _ out.PolicyEngine = (*Engine)(nil)

// New builds a policy engine. An empty defaultAction means deny.
func New(rules []domain.PolicyRule, defaultAction domain.Action) *Engine {
	if defaultAction == "" {
		defaultAction = domain.ActionDeny
	}
	return &Engine{rules: rules, def: defaultAction}
}

// Evaluate returns the first matching rule, or the default action.
func (e *Engine) Evaluate(_ context.Context, q *domain.ValidatedQuery) domain.PolicyDecision {
	for _, r := range e.rules {
		if ruleMatches(r, q) {
			return domain.PolicyDecision{Action: r.Action, Rule: r.Name, Reason: r.Reason}
		}
	}
	return domain.PolicyDecision{Action: e.def, Rule: "default", Reason: "no rule matched"}
}

func ruleMatches(r domain.PolicyRule, q *domain.ValidatedQuery) bool {
	m := r.Match
	if len(m.Statements) > 0 && !anyFold(m.Statements, string(q.StatementType)) {
		return false
	}
	if len(m.Connections) > 0 && !anyFold(m.Connections, q.Query.Connection) {
		return false
	}
	if len(m.Tables) > 0 && !anyGlob(m.Tables, q.Tables) {
		return false
	}
	if len(m.Functions) > 0 && !anyGlob(m.Functions, q.Functions) {
		return false
	}
	return true
}

func anyFold(list []string, s string) bool {
	for _, v := range list {
		if strings.EqualFold(v, s) {
			return true
		}
	}
	return false
}

func anyGlob(patterns, values []string) bool {
	for _, p := range patterns {
		pl := strings.ToLower(p)
		for _, v := range values {
			if ok, _ := path.Match(pl, strings.ToLower(v)); ok {
				return true
			}
		}
	}
	return false
}
