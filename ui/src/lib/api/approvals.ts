import { adminToken, apiBaseUrl, apiFetch } from './client';
import type { ApiApproval, ApiApprovalDetail, ApiDecision, Page } from './types';

export interface ApprovalFilters {
  state?: 'pending' | 'approved' | 'rejected' | 'expired';
  connection?: string;
  limit?: number;
  cursor?: string;
}

export function listApprovals(filters: ApprovalFilters = {}): Promise<Page<ApiApproval>> {
  const q = new URLSearchParams();
  if (filters.state) q.set('state', filters.state);
  if (filters.connection) q.set('connection', filters.connection);
  if (filters.limit) q.set('limit', String(filters.limit));
  if (filters.cursor) q.set('cursor', filters.cursor);
  const suffix = q.size > 0 ? `?${q}` : '';
  return apiFetch<Page<ApiApproval>>(`/api/v1/approvals${suffix}`);
}

export function getApproval(id: string): Promise<ApiApprovalDetail> {
  return apiFetch<ApiApprovalDetail>(`/api/v1/approvals/${encodeURIComponent(id)}`);
}

export function approveApproval(id: string, reason: string): Promise<ApiDecision> {
  return apiFetch<ApiDecision>(`/api/v1/approvals/${encodeURIComponent(id)}/approve`, {
    method: 'POST',
    body: { reason }
  });
}

export function rejectApproval(id: string, reason: string): Promise<ApiDecision> {
  return apiFetch<ApiDecision>(`/api/v1/approvals/${encodeURIComponent(id)}/reject`, {
    method: 'POST',
    body: { reason }
  });
}

export type StreamHandler = () => void;

const STREAM_EVENTS = [
  'approval.created',
  'approval.approved',
  'approval.rejected',
  'approval.used',
  'approval.expired',
  'write.approved',
  'write.rejected'
];

/**
 * Live approval events. EventSource cannot set headers, so the token (when
 * configured) travels as ?access_token= — accepted only by /approvals/stream.
 * Returns an unsubscribe function. Caller must provide a polling fallback.
 */
export function openApprovalStream(onEvent: StreamHandler): () => void {
  const token = adminToken();
  const url = `${apiBaseUrl()}/api/v1/approvals/stream${token ? `?access_token=${encodeURIComponent(token)}` : ''}`;
  const source = new EventSource(url);
  for (const type of STREAM_EVENTS) {
    source.addEventListener(type, onEvent);
  }
  return () => source.close();
}
