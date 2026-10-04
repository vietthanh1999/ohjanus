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

export interface CreateConnectionInput {
  name: string;
  driver?: string;
  dsn?: string;
  host?: string;
  port?: number;
  database?: string;
  username?: string;
  password?: string;
  sslmode?: string;
  read_only?: boolean;
}

export function createConnection(input: CreateConnectionInput): Promise<ApiConnection> {
  return apiFetch<ApiConnection>('/api/v1/connections', { method: 'POST', body: input });
}

export interface ConnectionProbe {
  status: string;
  latency_ms: number;
}

/** Test parameters without saving (the "Test Connection" button). */
export function probeConnection(input: CreateConnectionInput): Promise<ConnectionProbe> {
  return apiFetch<ConnectionProbe>('/api/v1/connections/test', { method: 'POST', body: input });
}
