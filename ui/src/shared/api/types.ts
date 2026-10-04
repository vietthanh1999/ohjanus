// Raw Admin API shapes (snake_case, as returned by janus serve :8788).
// UI-facing mapping (camelCase view models) lives in entities/*/model/types.ts.

export interface Page<T> {
  items: T[];
  next_cursor: string;
  total: number;
}

export interface ApiRequestedBy {
  token_id: string;
}

export interface ApiApproval {
  id: string;
  state: 'pending' | 'approved' | 'rejected' | 'expired';
  created_at: string;
  expires_at: string;
  connection: string;
  statement_type: string;
  sql: string;
  params?: unknown;
  affected_estimate: number;
  warnings: string[] | null;
  requested_by: ApiRequestedBy;
}

export interface ApiApprovalDetail extends ApiApproval {
  sql_hash: string;
  params_hash: string;
  plan?: string;
  decided_by: { id: string } | null;
  decision_reason?: string;
}

export interface ApiDecision {
  id: string;
  state: string;
  decided_by: { id: string };
  decided_at: string;
  decision_reason: string;
}

export interface ApiAuditEvent {
  id: string;
  ts: string;
  event: string;
  request_id: string;
  token_id?: string;
  connection?: string;
  tool?: string;
  sql_hash?: string;
  sql_normalized?: string;
  statement_type?: string;
  tables?: string[] | null;
  policy_decision?: string;
  policy_rule?: string;
  row_count?: number;
  truncated?: boolean;
  duration_ms?: number;
  status?: string;
  error?: string;
}

export interface ApiConnection {
  name: string;
  driver: string;
  readonly: boolean;
  status: string;
  last_ping_at: string;
  allowed_schemas?: string[] | null;
  denied_tables?: string[] | null;
}

export interface ApiConnectionTest {
  name: string;
  status: string;
  latency_ms: number;
}

export interface ApiToken {
  id: string;
  name: string;
  scopes: string[];
  created_at?: string;
  expires_at?: string | null;
  last_used_at?: string;
  state: 'active' | 'revoked' | 'expired';
}

export interface ApiTokenCreated extends ApiToken {
  token: string;
}

export interface ApiSummary {
  pending_approvals: number;
  requests_total: number;
  denials_total: number;
}

export interface ApiColumn {
  name: string;
  type: string;
  nullable: boolean;
}

export interface ApiTable {
  name: string;
  columns: ApiColumn[];
  primary_key?: string[] | null;
}

export interface ApiRoutine {
  name: string;
  kind: string;
}

export interface ApiSequence {
  name: string;
}

export interface ApiSchema {
  name: string;
  tables: ApiTable[];
  routines?: ApiRoutine[] | null;
  sequences?: ApiSequence[] | null;
}

export interface ApiQueryResult {
  columns: string[];
  rows: unknown[][];
  row_count: number;
  truncated: boolean;
  duration_ms: number;
}

export interface ApiExplainResult {
  plan: string;
  affected_estimate: number;
}

export interface ApiErrorBody {
  error: { code: string; message: string; request_id: string };
}
