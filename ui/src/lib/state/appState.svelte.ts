import { ApiError } from '../api/client';
import {
  listApprovals,
  approveApproval,
  rejectApproval,
  openApprovalStream
} from '../api/approvals';
import { listAudit } from '../api/audit';
import { listConnections, testConnection } from '../api/connections';
import { listTokens, createToken as apiCreateToken, revokeToken as apiRevokeToken } from '../api/tokens';
import { getSummary } from '../api/dashboard';
import { getSchema } from '../api/schema';
import { runQuery, explainQuery, buildTableSelect } from '../api/query';
import type {
  ApiApproval,
  ApiAuditEvent,
  ApiConnection,
  ApiSchema,
  ApiToken
} from '../api/types';

// ---- Tab registry (single convention, no hardcoded hosts/tables) ----

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

// ---- View models (camelCase; snake_case mapping stays in api/) ----

export interface ApprovalRequest {
  id: string;
  state: 'pending' | 'approved' | 'rejected' | 'expired';
  connection: string;
  statement_type: 'SELECT' | 'INSERT' | 'UPDATE' | 'DELETE' | 'DROP';
  sql: string;
  params?: Record<string, any>;
  affected_estimate: number;
  requested_by: {
    client: string;
    token_id: string;
  };
  warnings: string[];
  created_at: string;
  expires_at: string;
  decided_by?: string;
  decided_at?: string;
  decision_reason?: string;
}

export interface AuditRecord {
  id: string;
  ts: string;
  event: string;
  request_id: string;
  token_id: string;
  client: string;
  connection: string;
  statement_type: string;
  tables: string[];
  policy_decision: 'ALLOW' | 'DENY' | 'REQUIRE_APPROVAL';
  policy_rule: string;
  row_count: number;
  duration_ms: number;
  status: 'OK' | 'DENIED' | 'ERROR';
  sql_normalized: string;
}

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

export interface QueryResult {
  columns: string[];
  rows: unknown[][];
  rowCount: number;
  truncated: boolean;
  durationMs: number;
}

export interface DashboardSummary {
  pending_approvals: number;
  requests_total: number;
  denials_total: number;
}

export interface McpToken {
  id: string;
  name: string;
  scopes: string[];
  created_at: string;
  expires_at: string;
  last_used_at: string;
  state: 'active' | 'revoked' | 'expired';
  rawToken?: string;
}

export interface ConsoleLogEntry {
  id: string;
  timestamp: string;
  connection?: string;
  querySnippet?: string;
  summary: string;
  type: 'info' | 'query' | 'success' | 'warning' | 'error';
}

export interface ServiceSession {
  id: string;
  name: string;
  type: 'table' | 'console';
  connection?: string;
  durationMs?: number;
}

function nowStamp(): string {
  return new Date().toISOString().replace('T', ' ').substring(0, 19);
}

class AppStateManager {
  // ---- Shell: tabs ----
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

  // ---- Shell: sidebar layout ----
  sidebarWidth = $state<number>(290);
  servicesHeight = $state<number>(230);
  isSidebarCollapsed = $state<boolean>(false);
  treeFilterQuery = $state<string>('');
  treeExpanded = $state<Record<string, boolean>>({});
  selectedTreeNode = $state<string | null>(null);

  toggleTree(id: string) {
    this.treeExpanded[id] = !this.treeExpanded[id];
  }

  // ---- Explorer: live connections + schema cache (no mocks) ----
  connections = $state<ConnectionItem[]>([]);
  schemas = $state<Record<string, SchemaSummary[]>>({});
  schemaLoading = $state<Record<string, boolean>>({});
  schemaError = $state<Record<string, string | null>>({});

  get explorerConnections(): ConnectionItem[] {
    const q = this.treeFilterQuery.trim().toLowerCase();
    if (!q) return this.connections;
    return this.connections.filter((c) => c.name.toLowerCase().includes(q));
  }

  schemasOf(connection: string): SchemaSummary[] {
    return this.schemas[connection] ?? [];
  }

  findTable(connection: string, schema: string, table: string): TableSummary | null {
    for (const s of this.schemasOf(connection)) {
      if (s.name !== schema) continue;
      for (const t of s.tables) {
        if (t.name === table) return t;
      }
    }
    return null;
  }

  async loadSchema(connection: string, force = false): Promise<void> {
    if (!connection) return;
    if (this.schemaLoading[connection]) return;
    if (!force && this.schemas[connection]) return;
    this.schemaLoading[connection] = true;
    this.schemaError[connection] = null;
    try {
      const { schemas } = await getSchema({ connection });
      this.schemas[connection] = schemas.map((s: ApiSchema) => ({
        name: s.name,
        tables: (s.tables ?? []).map((t) => ({
          name: t.name,
          columns: (t.columns ?? []).map((c) => ({ name: c.name, type: c.type, nullable: c.nullable })),
          primaryKey: t.primary_key ?? []
        }))
      }));
    } catch (e) {
      this.schemaError[connection] = e instanceof ApiError ? e.message : String(e);
    } finally {
      this.schemaLoading[connection] = false;
    }
  }

  /** Explorer selection: table -> opens table tab; connection -> opens console. */
  async selectTable(connection: string, schema: string, table: string): Promise<void> {
    this.selectedTreeNode = `table:${connection}.${schema}.${table}`;
    await this.loadSchema(connection);
    this.openTableTab(connection, schema, table);
  }

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
    // Sync the table viewer to the newly focused tab.
    this.tableViewer.connection = connection;
    this.tableViewer.schema = schema;
    this.tableViewer.table = table;
    void this.loadTableData();
  }

  openConsoleTab(connection: string) {
    if (!connection) return;
    const id = consoleTabId(connection);
    this.openTab({
      id,
      title: `console [${connection}]`,
      type: 'console',
      closable: true,
      icon: 'lightning',
      connection
    });
    if (this.console.connection !== connection) {
      this.console.connection = connection;
    }
  }

  // ---- SQL console: live query execution (no simulation) ----
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
    const started = performance.now();
    this.pushLog({ type: 'query', connection, querySnippet: sql.slice(0, 240), summary: 'Query submitted.' });
    try {
      const res = await runQuery({ connection, sql, limit: 500 });
      this.console.columns = res.columns ?? [];
      this.console.rows = res.rows ?? [];
      this.console.rowCount = res.row_count;
      this.console.truncated = res.truncated;
      this.console.durationMs = res.duration_ms;
      this.lastDurationBySession[`console:${connection}`] = res.duration_ms;
      void started;
      this.pushLog({
        type: 'success',
        summary: `${res.row_count} row(s)${res.truncated ? ' (truncated)' : ''} in ${res.duration_ms} ms`
      });
    } catch (e) {
      const msg = e instanceof ApiError ? `${e.code}: ${e.message}` : String(e);
      this.console.error = msg;
      this.console.columns = [];
      this.console.rows = [];
      this.console.rowCount = 0;
      this.pushLog({ type: 'error', summary: msg });
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
      this.pushLog({ type: 'info', summary: `EXPLAIN ok (estimate ${res.affected_estimate})` });
    } catch (e) {
      this.console.plan = null;
      this.pushLog({ type: 'error', summary: e instanceof Error ? e.message : String(e) });
    }
  }

  // ---- Table viewer: live SELECT over the focused table ----
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
      this.lastDurationBySession[tableTabId(t.connection, t.schema, t.table)] = res.duration_ms;
      this.pushLog({
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
      this.pushLog({ type: 'error', summary: msg });
    } finally {
      t.loading = false;
    }
  }

  /** DDL generated from the live schema cache (no hardcoded strings). */
  ddlFor(connection: string, schema: string, table: string): string | null {
    const t = this.findTable(connection, schema, table);
    if (!t) return null;
    const lines = t.columns.map((c) => `    "${c.name}" ${c.type}${c.nullable ? '' : ' NOT NULL'}`);
    if (t.primaryKey.length > 0) {
      lines.push(`    CONSTRAINT ${table}_pkey PRIMARY KEY (${t.primaryKey.map((k) => `"${k}"`).join(', ')})`);
    }
    const qualified = schema ? `"${schema}"."${table}"` : `"${table}"`;
    return `-- Generated from live schema: ${connection}.${schema}.${table}\nCREATE TABLE ${qualified} (\n${lines.join(',\n')}\n);`;
  }

  // ---- Services panel + console log: derived from live sessions ----
  lastDurationBySession = $state<Record<string, number>>({});

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
    if (this.consoleLogs.length > 300) {
      this.consoleLogs = this.consoleLogs.slice(-300);
    }
  }

  clearLogs() {
    this.consoleLogs = [];
  }

  // ---- Gateway entities (approvals / audit / tokens / dashboard) ----
  approvals = $state<ApprovalRequest[]>([]);
  auditLogs = $state<AuditRecord[]>([]);
  tokens = $state<McpToken[]>([]);
  summary = $state<DashboardSummary | null>(null);

  dataLoading = $state<boolean>(true);
  dataError = $state<string | null>(null);

  searchModalOpen = $state<boolean>(false);
  settingsModalOpen = $state<boolean>(false);
  ddlModalOpen = $state<boolean>(false);
  createTokenModalOpen = $state<boolean>(false);
  createdTokenSecret = $state<string | null>(null);
  selectedApprovalForAction = $state<ApprovalRequest | null>(null);
  approvalDecisionMode = $state<'approve' | 'reject'>('approve');
  approvalDecisionReason = $state<string>('');

  notificationCount = $derived(
    this.approvals.filter((a) => a.state === 'pending').length
  );

  // ---- UI density ----
  uiDensity = $state<'compact' | 'standard' | 'comfortable'>('comfortable');

  setDensity(density: 'compact' | 'standard' | 'comfortable') {
    this.uiDensity = density;
    if (typeof document !== 'undefined') {
      document.documentElement.setAttribute('data-density', density);
      try {
        localStorage.setItem('ohjanus_ui_density', density);
      } catch {
        // Ignore storage errors in restricted contexts
      }
    }
  }

  initDensity() {
    if (typeof document !== 'undefined') {
      try {
        const saved = localStorage.getItem('ohjanus_ui_density') as 'compact' | 'standard' | 'comfortable' | null;
        this.uiDensity = saved || 'comfortable';
      } catch {
        this.uiDensity = 'comfortable';
      }
      document.documentElement.setAttribute('data-density', this.uiDensity);
    }
  }

  private mapApproval(t: ApiApproval): ApprovalRequest {
    return {
      id: t.id,
      state: t.state,
      connection: t.connection,
      statement_type: t.statement_type as ApprovalRequest['statement_type'],
      sql: t.sql,
      params: (t.params as Record<string, any> | undefined) ?? {},
      affected_estimate: t.affected_estimate,
      requested_by: { client: t.requested_by.token_id, token_id: t.requested_by.token_id },
      warnings: t.warnings ?? [],
      created_at: t.created_at,
      expires_at: t.expires_at
    };
  }

  private mapAudit(e: ApiAuditEvent): AuditRecord {
    const status = e.status === 'success' ? 'OK' : e.status === 'denied' ? 'DENIED' : 'ERROR';
    return {
      id: e.id,
      ts: e.ts,
      event: e.event,
      request_id: e.request_id,
      token_id: e.token_id ?? '',
      client: e.token_id ?? '',
      connection: e.connection ?? '',
      statement_type: e.statement_type ?? '',
      tables: e.tables ?? [],
      policy_decision: (e.policy_decision ?? '') as AuditRecord['policy_decision'],
      policy_rule: e.policy_rule ?? '',
      row_count: e.row_count ?? 0,
      duration_ms: e.duration_ms ?? 0,
      status,
      sql_normalized: e.sql_normalized ?? ''
    };
  }

  private mapConnection(c: ApiConnection): ConnectionItem {
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

  private mapToken(t: ApiToken): McpToken {
    return {
      id: t.id,
      name: t.name,
      scopes: t.scopes,
      created_at: t.created_at ?? '—',
      expires_at: t.expires_at ?? '—',
      last_used_at: t.last_used_at ?? '—',
      state: t.state
    };
  }

  async loadApprovals() {
    const page = await listApprovals({ limit: 100 });
    this.approvals = (page.items ?? []).map((t) => this.mapApproval(t));
  }

  async loadAudit() {
    const page = await listAudit({ limit: 100 });
    this.auditLogs = (page.items ?? []).map((e) => this.mapAudit(e));
  }

  async loadConnections() {
    const { items } = await listConnections();
    const prev = new Map(this.connections.map((c) => [c.name, c.latency_ms]));
    this.connections = (items ?? []).map((c) => ({ ...this.mapConnection(c), latency_ms: prev.get(c.name) }));
    // Default the console to the first connection so empty states resolve
    // to a real target instead of a hardcoded sample host.
    if (!this.console.connection && this.connections.length > 0) {
      this.console.connection = this.connections[0].name;
    }
  }

  async loadTokens() {
    const { items } = await listTokens();
    this.tokens = (items ?? []).map((t) => this.mapToken(t));
  }

  async loadSummary() {
    this.summary = await getSummary();
  }

  async loadAll() {
    this.initDensity();
    this.dataLoading = true;
    this.dataError = null;
    try {
      await Promise.all([
        this.loadApprovals(),
        this.loadAudit(),
        this.loadConnections(),
        this.loadTokens(),
        this.loadSummary()
      ]);
      this.pushLog({ type: 'info', summary: `Connected: ${this.connections.length} connection(s) loaded from Admin API.` });
      // Preload schema for the default connection so the explorer is live
      // on first paint; other connections load lazily on expand.
      if (this.console.connection) {
        await this.loadSchema(this.console.connection).catch(() => {});
      }
    } catch (e) {
      this.dataError = e instanceof ApiError ? e.message : String(e);
      this.pushLog({ type: 'error', summary: `Admin API unreachable: ${this.dataError}` });
    } finally {
      this.dataLoading = false;
    }
  }

  /** Live approval updates via SSE, with 10s polling fallback. */
  private stopLive: (() => void) | null = null;
  private pollTimer: ReturnType<typeof setInterval> | null = null;

  startLive() {
    this.stopLive?.();
    let stopped = false;
    const stop = () => {
      stopped = true;
    };
    this.stopLive = stop;

    const poll = async () => {
      if (stopped) return;
      try {
        await this.loadApprovals();
      } catch {
        // Error state is surfaced by explicit reloads; polling stays quiet.
      }
    };

    try {
      const closeStream = openApprovalStream(() => void poll());
      const prevStop = stop;
      this.stopLive = () => {
        prevStop();
        closeStream();
        if (this.pollTimer) clearInterval(this.pollTimer);
        this.pollTimer = null;
      };
      // Fallback polling in case the stream silently drops.
      this.pollTimer = setInterval(() => void poll(), 10000);
    } catch {
      this.pollTimer = setInterval(() => void poll(), 10000);
    }
  }

  stopLiveUpdates() {
    this.stopLive?.();
    this.stopLive = null;
  }

  async confirmApprovalAction() {
    if (!this.selectedApprovalForAction) return;
    const id = this.selectedApprovalForAction.id;
    const mode = this.approvalDecisionMode;
    const reason = this.approvalDecisionReason || (mode === 'approve' ? 'Approved via Admin UI' : '');
    if (mode === 'approve') {
      await approveApproval(id, reason);
    } else {
      await rejectApproval(id, reason);
    }
    await this.loadApprovals();
    await this.loadAudit().catch(() => {});
    const apprTab = this.tabs.find((t) => t.id === 'approvals');
    if (apprTab) {
      const pendingCount = this.approvals.filter((a) => a.state === 'pending').length;
      apprTab.badge = pendingCount > 0 ? String(pendingCount) : undefined;
    }
    this.selectedApprovalForAction = null;
    this.approvalDecisionReason = '';
  }

  async createToken(name: string, scopes: string[], ttlDays: number) {
    const created = await apiCreateToken({ name, scopes, ttl_hours: ttlDays * 24 });
    await this.loadTokens();
    this.createdTokenSecret = created.token;
    return created.token;
  }

  async revokeToken(id: string) {
    await apiRevokeToken(id);
    await this.loadTokens();
  }

  async pingConnection(name: string) {
    const res = await testConnection(name);
    const conn = this.connections.find((c) => c.name === name);
    if (conn) {
      conn.status = res.status;
      conn.latency_ms = res.latency_ms;
      conn.last_ping_at = nowStamp();
    }
    return res;
  }
}

export const appState = new AppStateManager();
