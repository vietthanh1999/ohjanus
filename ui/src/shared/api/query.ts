import { apiFetch } from './client';
import type { ApiExplainResult, ApiQueryResult } from './types';

export interface QueryInput {
  connection: string;
  sql: string;
  params?: unknown[];
  limit?: number;
}

/** Read-only query execution. Never executes writes; the gateway enforces policy. */
export function runQuery(input: QueryInput): Promise<ApiQueryResult> {
  return apiFetch<ApiQueryResult>('/api/v1/query', { method: 'POST', body: input });
}

/** EXPLAIN without executing the query. */
export function explainQuery(input: QueryInput): Promise<ApiExplainResult> {
  return apiFetch<ApiExplainResult>('/api/v1/explain', { method: 'POST', body: input });
}
