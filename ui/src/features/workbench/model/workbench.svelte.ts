import type { ServiceSession, ConsoleLogEntry } from '@/entities/session';
import { tableTabId } from '@/entities/tab';
import type { TabItem } from '@/entities/tab';
import { nowStamp } from '@/shared/lib';
import { ACTIVITY_LOG_CAP } from '@/shared/config';

/**
 * Shell state: open tabs + activity log + per-session durations.
 * Lower layer than business features — console/table-viewer/explorer may
 * depend on workbench, never the other way around.
 */
class WorkbenchManager {
  tabs = $state<TabItem[]>([
    { id: 'approvals', title: 'Approvals Queue', type: 'approvals', closable: false, icon: 'shield' },
    { id: 'audit', title: 'Audit Log Trail', type: 'audit', closable: false, icon: 'audit' },
    { id: 'connections', title: 'Connection Pools', type: 'connections', closable: false, icon: 'database' },
    { id: 'tokens', title: 'MCP Agent Tokens', type: 'tokens', closable: false, icon: 'key' },
    { id: 'dashboard', title: 'Telemetry Dashboard', type: 'dashboard', closable: false, icon: 'chart' }
  ]);

  activeTabId = $state<string>('approvals');

  get activeTab(): TabItem | undefined {
    return this.tabs.find((t) => t.id === this.activeTabId);
  }

  openTab(tab: TabItem) {
    const existing = this.tabs.find((t) => t.id === tab.id);
    if (!existing) {
      this.tabs = [...this.tabs, tab];
    }
    this.activeTabId = tab.id;
  }

  closeTab(tabId: string) {
    const tab = this.tabs.find((t) => t.id === tabId);
    if (!tab || !tab.closable) return;
    const idx = this.tabs.findIndex((t) => t.id === tabId);
    this.tabs = this.tabs.filter((t) => t.id !== tabId);
    if (this.activeTabId === tabId && this.tabs.length > 0) {
      this.activeTabId = this.tabs[Math.max(0, idx - 1)].id;
    }
  }

  /** Open (or focus) the console tab for a connection. */
  openConsoleTab(connection: string) {
    if (!connection) return;
    this.openTab({
      id: `console:${connection}`,
      title: `console [${connection}]`,
      type: 'console',
      closable: true,
      icon: 'lightning',
      connection
    });
  }

  /** Open (or focus) the table tab for a table. */
  openTableTab(connection: string, schema: string, table: string) {
    const id = tableTabId(connection, schema, table);
    this.openTab({
      id,
      title: `${table} [${connection}.${schema}]`,
      type: 'table',
      closable: true,
      icon: 'table',
      connection,
      schema,
      table
    });
  }

  // ---- Activity log + per-session durations ----

  lastDurationBySession = $state<Record<string, number>>({});

  recordDuration(sessionId: string, durationMs: number) {
    this.lastDurationBySession[sessionId] = durationMs;
  }

  services = $derived.by((): ServiceSession[] => {
    return this.tabs
      .filter((t) => t.type === 'table' || t.type === 'console')
      .map((t) => ({
        id: t.id,
        name:
          t.type === 'table' && t.table
            ? t.table
            : (t.connection ?? t.title),
        type: t.type as 'table' | 'console',
        connection: t.connection,
        durationMs: this.lastDurationBySession[t.id]
      }));
  });

  consoleLogs = $state<ConsoleLogEntry[]>([]);

  pushLog(entry: Omit<ConsoleLogEntry, 'id' | 'timestamp'> & { timestamp?: string }): void {
    this.consoleLogs.push({
      id: `log-${Date.now()}-${Math.floor(Math.random() * 1e6)}`,
      timestamp: entry.timestamp ?? nowStamp(),
      ...entry
    });
    // Keep the log bounded so long sessions don't grow the DOM forever.
    if (this.consoleLogs.length > ACTIVITY_LOG_CAP) {
      this.consoleLogs = this.consoleLogs.slice(-ACTIVITY_LOG_CAP);
    }
  }

  clearLogs() {
    this.consoleLogs = [];
  }
}

export const workbenchState = new WorkbenchManager();
