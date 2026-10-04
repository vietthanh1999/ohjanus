import { apiFetch } from '@/shared/api';
import type { ApiConnection, ApiConnectionTest } from '@/shared/api';

export function listConnections(): Promise<{ items: ApiConnection[] }> {
  return apiFetch<{ items: ApiConnection[] }>('/api/v1/connections');
}

export function testConnection(name: string): Promise<ApiConnectionTest> {
  return apiFetch<ApiConnectionTest>(`/api/v1/connections/${encodeURIComponent(name)}/test`, {
    method: 'POST'
  });
}
