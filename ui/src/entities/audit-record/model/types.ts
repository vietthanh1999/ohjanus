import type { ApiAuditEvent } from '@/shared/api';

/** UI-facing audit record (camelCase; snake_case mapping stays in api/). */
export interface AuditRecord {
  id: string;
  ts: string;
  event: string;
  request_id: string;
  token_id: string;
  client: string;
  connection: string;
  statement_type: string;
  tables: string[];
  policy_decision: 'ALLOW' | 'DENY' | 'REQUIRE_APPROVAL';
  policy_rule: string;
  row_count: number;
  duration_ms: number;
  status: 'OK' | 'DENIED' | 'ERROR';
  sql_normalized: string;
}

export function mapAuditEvent(e: ApiAuditEvent): AuditRecord {
  const status = e.status === 'success' ? 'OK' : e.status === 'denied' ? 'DENIED' : 'ERROR';
  return {
    id: e.id,
    ts: e.ts,
    event: e.event,
    request_id: e.request_id,
    token_id: e.token_id ?? '',
    client: e.token_id ?? '',
    connection: e.connection ?? '',
    statement_type: e.statement_type ?? '',
    tables: e.tables ?? [],
    policy_decision: (e.policy_decision ?? '') as AuditRecord['policy_decision'],
    policy_rule: e.policy_rule ?? '',
    row_count: e.row_count ?? 0,
    duration_ms: e.duration_ms ?? 0,
    status,
    sql_normalized: e.sql_normalized ?? ''
  };
}
