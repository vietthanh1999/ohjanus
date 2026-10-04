import type { ApiApproval } from '@/shared/api';

/** UI-facing approval request (camelCase; snake_case mapping stays in api/). */
export interface ApprovalRequest {
  id: string;
  state: 'pending' | 'approved' | 'rejected' | 'expired';
  connection: string;
  statement_type: 'SELECT' | 'INSERT' | 'UPDATE' | 'DELETE' | 'DROP';
  sql: string;
  params?: Record<string, any>;
  affected_estimate: number;
  requested_by: {
    client: string;
    token_id: string;
  };
  warnings: string[];
  created_at: string;
  expires_at: string;
  decided_by?: string;
  decided_at?: string;
  decision_reason?: string;
}

export function mapApproval(t: ApiApproval): ApprovalRequest {
  return {
    id: t.id,
    state: t.state,
    connection: t.connection,
    statement_type: t.statement_type as ApprovalRequest['statement_type'],
    sql: t.sql,
    params: (t.params as Record<string, any> | undefined) ?? {},
    affected_estimate: t.affected_estimate,
    requested_by: { client: t.requested_by.token_id, token_id: t.requested_by.token_id },
    warnings: t.warnings ?? [],
    created_at: t.created_at,
    expires_at: t.expires_at
  };
}
