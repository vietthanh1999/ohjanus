import { apiFetch } from './client';
import type { ApiToken, ApiTokenCreated } from './types';

export function listTokens(): Promise<{ items: ApiToken[] }> {
  return apiFetch<{ items: ApiToken[] }>('/api/v1/tokens');
}

export interface CreateTokenInput {
  name: string;
  scopes: string[];
  ttl_hours: number;
}

export function createToken(input: CreateTokenInput): Promise<ApiTokenCreated> {
  return apiFetch<ApiTokenCreated>('/api/v1/tokens', { method: 'POST', body: input });
}

export async function revokeToken(id: string): Promise<void> {
  await apiFetch<undefined>(`/api/v1/tokens/${encodeURIComponent(id)}`, { method: 'DELETE' });
}
