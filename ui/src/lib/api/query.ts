import { apiFetch } from './client';
import type { ApiExplainResult, ApiQueryResult } from './types';

export interface QueryInput {
  connection: string;
  sql: string;
  params?: unknown[];
  limit?: number;
}

export function runQuery(input: QueryInput): Promise<ApiQueryResult> {
  return apiFetch<ApiQueryResult>('/api/v1/query', { method: 'POST', body: input });
}

export function explainQuery(input: QueryInput): Promise<ApiExplainResult> {
  return apiFetch<ApiExplainResult>('/api/v1/explain', { method: 'POST', body: input });
}

/** Quote an identifier for generated SELECTs (table viewer, DDL). */
export function quoteIdent(name: string): string {
  return `"${name.replace(/"/g, '""')}"`;
}

/** Build the SELECT used by the table viewer from WHERE/ORDER BY fragments. */
export function buildTableSelect(
  schema: string,
  table: string,
  where: string,
  orderBy: string,
  limit: number
): string {
  const from = schema ? `${quoteIdent(schema)}.${quoteIdent(table)}` : quoteIdent(table);
  let sql = `SELECT * FROM ${from}`;
  if (where.trim()) sql += ` WHERE ${where.trim()}`;
  if (orderBy.trim()) sql += ` ORDER BY ${orderBy.trim()}`;
  sql += ` LIMIT ${limit}`;
  return sql;
}
