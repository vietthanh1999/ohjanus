package domain

// Action is the outcome of policy evaluation.
type Action string

const (
	ActionAllow           Action = "allow"
	ActionDeny            Action = "deny"
	ActionRequireApproval Action = "require_approval"
)

// PolicyDecision is the result of evaluating one validated query.
type PolicyDecision struct {
	Action Action
	Rule   string
	Reason string
}

// RuleMatch describes what a policy rule matches on.
type RuleMatch struct {
	Statements  []string
	Connections []string
	Tables      []string
	Functions   []string
}

// PolicyRule is a single allow/deny/require_approval rule.
// Rules are evaluated in order; evaluation stops at the first match.
type PolicyRule struct {
	Name   string
	Match  RuleMatch
	Action Action
	Reason string
}
