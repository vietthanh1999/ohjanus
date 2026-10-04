export type TabType =
  | 'console'
  | 'table'
  | 'approvals'
  | 'audit'
  | 'tokens'
  | 'connections'
  | 'dashboard';

export interface TabItem {
  id: string;
  title: string;
  type: TabType;
  closable: boolean;
  icon: string;
  badge?: string;
  /** Live data-plane coordinates. Only set for console/table tabs. */
  connection?: string;
  schema?: string;
  table?: string;
}

export const tableTabId = (connection: string, schema: string, table: string) =>
  `table:${connection}.${schema}.${table}`;

export const consoleTabId = (connection: string) => `console:${connection}`;

export function parseTableTabId(id: string): { connection: string; schema: string; table: string } | null {
  const m = /^table:(.+)\.(.+)\.(.+)$/.exec(id);
  if (!m) return null;
  return { connection: m[1], schema: m[2], table: m[3] };
}

export function parseConsoleTabId(id: string): { connection: string } | null {
  const m = /^console:(.+)$/.exec(id);
  if (!m) return null;
  return { connection: m[1] };
}
