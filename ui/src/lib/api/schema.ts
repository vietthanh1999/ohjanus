import { apiFetch } from './client';
import type { ApiSchema } from './types';

export interface SchemaFilters {
  connection: string;
  schema?: string;
  table?: string;
}

export function getSchema(filters: SchemaFilters): Promise<{ schemas: ApiSchema[] }> {
  const q = new URLSearchParams({ connection: filters.connection });
  if (filters.schema) q.set('schema', filters.schema);
  if (filters.table) q.set('table', filters.table);
  return apiFetch<{ schemas: ApiSchema[] }>(`/api/v1/schema?${q}`);
}
