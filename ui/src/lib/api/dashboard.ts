import { apiFetch } from './client';
import type { ApiSummary } from './types';

export function getSummary(): Promise<ApiSummary> {
  return apiFetch<ApiSummary>('/api/v1/dashboard/summary');
}
