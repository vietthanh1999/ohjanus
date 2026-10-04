import { apiFetch, type ApiSummary } from '@/shared/api';
import type { DashboardSummary } from '@/entities/session';

export function getSummary(): Promise<DashboardSummary> {
  return apiFetch<ApiSummary>('/api/v1/dashboard/summary');
}
