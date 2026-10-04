import type { ApiSchema } from '@/shared/api';

export interface ColumnSummary {
  name: string;
  type: string;
  nullable: boolean;
}

export interface TableSummary {
  name: string;
  columns: ColumnSummary[];
  primaryKey: string[];
}

export interface SchemaSummary {
  name: string;
  tables: TableSummary[];
}

export function mapSchema(s: ApiSchema): SchemaSummary {
  return {
    name: s.name,
    tables: (s.tables ?? []).map((t) => ({
      name: t.name,
      columns: (t.columns ?? []).map((c) => ({ name: c.name, type: c.type, nullable: c.nullable })),
      primaryKey: t.primary_key ?? []
    }))
  };
}
