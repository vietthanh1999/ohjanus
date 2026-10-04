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

export interface RoutineSummary {
  name: string;
  kind: string;
}

export interface SequenceSummary {
  name: string;
}

export interface SchemaSummary {
  name: string;
  tables: TableSummary[];
  routines: RoutineSummary[];
  sequences: SequenceSummary[];
}

export function mapSchema(s: ApiSchema): SchemaSummary {
  return {
    name: s.name,
    tables: (s.tables ?? []).map((t) => ({
      name: t.name,
      columns: (t.columns ?? []).map((c) => ({ name: c.name, type: c.type, nullable: c.nullable })),
      primaryKey: t.primary_key ?? []
    })),
    routines: (s.routines ?? []).map((r) => ({ name: r.name, kind: r.kind })),
    sequences: (s.sequences ?? []).map((q) => ({ name: q.name }))
  };
}
