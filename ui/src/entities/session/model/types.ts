/** Activity-log entry for the Services/console feed (shell-level, connection-agnostic). */
export interface ConsoleLogEntry {
  id: string;
  timestamp: string;
  connection?: string;
  querySnippet?: string;
  summary: string;
  type: 'info' | 'query' | 'success' | 'warning' | 'error';
}

/** Live session shown in the Services panel. Derived from open workbench tabs. */
export interface ServiceSession {
  id: string;
  name: string;
  type: 'table' | 'console';
  connection?: string;
  durationMs?: number;
}

export interface DashboardSummary {
  pending_approvals: number;
  requests_total: number;
  denials_total: number;
}

export interface QueryResult {
  columns: string[];
  rows: unknown[][];
  rowCount: number;
  truncated: boolean;
  durationMs: number;
}
