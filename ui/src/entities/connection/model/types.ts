import type { ApiConnection } from '@/shared/api';

/** UI-facing connection (camelCase; snake_case mapping stays in api/). */
export interface ConnectionItem {
  name: string;
  driver: string;
  readonly: boolean;
  status: string;
  last_ping_at: string;
  latency_ms?: number;
  allowed_schemas: string[];
  denied_tables: string[];
}

export function mapConnection(c: ApiConnection): ConnectionItem {
  return {
    name: c.name,
    driver: c.driver,
    readonly: c.readonly,
    status: c.status,
    last_ping_at: c.last_ping_at,
    allowed_schemas: c.allowed_schemas ?? [],
    denied_tables: c.denied_tables ?? []
  };
}
