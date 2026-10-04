import { workbenchState } from '@/features/workbench';
import { ApiError, runQuery, explainQuery } from '@/shared/api';
import { CONSOLE_QUERY_LIMIT } from '@/shared/config';
import { connectionsState } from '@/features/connections';

class ConsoleManager {
  console = $state({
    connection: '',
    sql: '-- Write a read-only query (SELECT / WITH / EXPLAIN / SHOW), then press Cmd+Enter.\nSELECT 1;',
    columns: [] as string[],
    rows: [] as unknown[][],
    rowCount: 0,
    truncated: false,
    durationMs: 0,
    isExecuting: false,
    error: null as string | null,
    plan: null as string | null
  });

  get consoleCursorStats(): { chars: number; lines: number } {
    const sql = this.console.sql;
    return { chars: sql.length, lines: sql.split('\n').length };
  }

  /** Default the console to the first connection when the list first loads. */
  ensureDefaultConnection() {
    if (!this.console.connection && connectionsState.connections.length > 0) {
      this.console.connection = connectionsState.connections[0].name;
    }
  }

  /** Focus a console: open its tab and point the editor at the connection. */
  openConsole(connection: string) {
    if (!connection) return;
    workbenchState.openConsoleTab(connection);
    if (this.console.connection !== connection) {
      this.console.connection = connection;
    }
  }

  async executeConsoleQuery(): Promise<void> {
    const { connection, sql } = this.console;
    if (!connection) {
      this.console.error = 'No connection selected. Pick a connection in the explorer first.';
      return;
    }
    if (!sql.trim()) {
      this.console.error = 'Query is empty.';
      return;
    }
    this.console.isExecuting = true;
    this.console.error = null;
    this.console.plan = null;
    workbenchState.pushLog({ type: 'query', connection, querySnippet: sql.slice(0, 240), summary: 'Query submitted.' });
    try {
      const res = await runQuery({ connection, sql, limit: CONSOLE_QUERY_LIMIT });
      this.console.columns = res.columns ?? [];
      this.console.rows = res.rows ?? [];
      this.console.rowCount = res.row_count;
      this.console.truncated = res.truncated;
      this.console.durationMs = res.duration_ms;
      workbenchState.recordDuration(`console:${connection}`, res.duration_ms);
      workbenchState.pushLog({
        type: 'success',
        summary: `${res.row_count} row(s)${res.truncated ? ' (truncated)' : ''} in ${res.duration_ms} ms`
      });
    } catch (e) {
      const msg = e instanceof ApiError ? `${e.code}: ${e.message}` : String(e);
      this.console.error = msg;
      this.console.columns = [];
      this.console.rows = [];
      this.console.rowCount = 0;
      workbenchState.pushLog({ type: 'error', summary: msg });
    } finally {
      this.console.isExecuting = false;
    }
  }

  async explainConsoleQuery(): Promise<void> {
    const { connection, sql } = this.console;
    if (!connection || !sql.trim()) return;
    try {
      const res = await explainQuery({ connection, sql });
      this.console.plan = res.plan;
      workbenchState.pushLog({ type: 'info', summary: `EXPLAIN ok (estimate ${res.affected_estimate})` });
    } catch (e) {
      this.console.plan = null;
      workbenchState.pushLog({ type: 'error', summary: e instanceof Error ? e.message : String(e) });
    }
  }
}

export const consoleState = new ConsoleManager();
