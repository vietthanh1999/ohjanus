import { apiDownload, apiFetch } from '@/shared/api';
import type { ApiAuditEvent, Page } from '@/shared/api';

export interface AuditFilters {
  event?: string;
  connection?: string;
  token_id?: string;
  request_id?: string;
  status?: string;
  q?: string;
  from?: string;
  to?: string;
  limit?: number;
  cursor?: string;
}

export function listAudit(filters: AuditFilters = {}): Promise<Page<ApiAuditEvent>> {
  const q = new URLSearchParams();
  for (const [key, value] of Object.entries(filters)) {
    if (value !== undefined && value !== '') q.set(key, String(value));
  }
  const suffix = q.size > 0 ? `?${q}` : '';
  return apiFetch<Page<ApiAuditEvent>>(`/api/v1/audit${suffix}`);
}

export function getAuditEvent(id: string): Promise<ApiAuditEvent> {
  return apiFetch<ApiAuditEvent>(`/api/v1/audit/${encodeURIComponent(id)}`);
}

export function exportAudit(format: 'csv' | 'jsonl', filters: AuditFilters = {}): Promise<void> {
  const q = new URLSearchParams({ format });
  for (const [key, value] of Object.entries(filters)) {
    if (value !== undefined && value !== '' && key !== 'limit' && key !== 'cursor') {
      q.set(key, String(value));
    }
  }
  const stamp = new Date().toISOString().substring(0, 10);
  return apiDownload(`/api/v1/audit/export?${q}`, `ohjanus_audit_${stamp}.${format}`);
}
