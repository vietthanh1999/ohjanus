import { tableTabId } from '@/entities/tab';
import { explorerState } from '@/features/explorer';
import { workbenchState } from '@/features/workbench';
import { ApiError, runQuery } from '@/shared/api';
import { buildTableSelect } from '@/shared/lib';

class TableViewerManager {
  tableViewer = $state({
    connection: '',
    schema: 'public',
    table: '',
    where: '',
    orderBy: '',
    limit: 500,
    columns: [] as string[],
    rows: [] as unknown[][],
    rowCount: 0,
    truncated: false,
    durationMs: 0,
    loading: false,
    error: null as string | null
  });

  get tableViewerSql(): string {
    const t = this.tableViewer;
    if (!t.connection || !t.table) return '';
    return buildTableSelect(t.schema, t.table, t.where, t.orderBy, t.limit);
  }

  /** Focus a table: open its tab, point the viewer at it, and load live data. */
  async openTable(connection: string, schema: string, table: string): Promise<void> {
    await explorerState.loadSchema(connection);
    if (explorerState.scrollFromEditor) {
      explorerState.revealTable(connection, schema, table);
    } else {
      explorerState.selectedTreeNode = `table:${connection}.${schema}.${table}`;
    }
    workbenchState.openTableTab(connection, schema, table);
    this.tableViewer.connection = connection;
    this.tableViewer.schema = schema;
    this.tableViewer.table = table;
    await this.loadTableData();
  }

  /** Point the viewer at the workbench's focused table tab (tab switches). */
  syncToActiveTab() {
    const tab = workbenchState.activeTab;
    if (tab?.type === 'table' && tab.connection && tab.table) {
      this.tableViewer.connection = tab.connection;
      this.tableViewer.schema = tab.schema ?? 'public';
      this.tableViewer.table = tab.table;
      void this.loadTableData();
    }
  }

  async loadTableData(): Promise<void> {
    const t = this.tableViewer;
    if (!t.connection || !t.table) return;
    const sql = this.tableViewerSql;
    t.loading = true;
    t.error = null;
    try {
      const res = await runQuery({ connection: t.connection, sql, limit: t.limit });
      t.columns = res.columns ?? [];
      t.rows = res.rows ?? [];
      t.rowCount = res.row_count;
      t.truncated = res.truncated;
      t.durationMs = res.duration_ms;
      workbenchState.recordDuration(tableTabId(t.connection, t.schema, t.table), res.duration_ms);
      workbenchState.pushLog({
        type: 'query',
        connection: `${t.connection}.${t.schema}`,
        querySnippet: sql.slice(0, 240),
        summary: `${t.table}: ${res.row_count} row(s) in ${res.duration_ms} ms`
      });
    } catch (e) {
      const msg = e instanceof ApiError ? `${e.code}: ${e.message}` : String(e);
      t.error = msg;
      t.columns = [];
      t.rows = [];
      t.rowCount = 0;
      workbenchState.pushLog({ type: 'error', summary: msg });
    } finally {
      t.loading = false;
    }
  }
}

export const tableViewerState = new TableViewerManager();
