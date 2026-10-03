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
import type {
  ApiApproval,
  ApiAuditEvent,
  ApiConnection,
  ApiToken
} from '../api/types';

export interface TabItem {
  id: string;
  title: string;
  type: 'console' | 'table' | 'approvals' | 'audit' | 'tokens' | 'connections' | 'dashboard';
  closable: boolean;
  icon: string;
  badge?: string;
  env?: string;
}

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

export interface CredentialRow {
  id: string;
  createdDate: string;
  lastUpdatedDate: string | null;
  createdBy: string;
}

class AppStateManager {
  // Navigation & Tabs
  tabs = $state<TabItem[]>([
    {
      id: 'commands',
      title: 'commands [[Dev][ReadOnly] 10.220.6.4]',
      type: 'table',
      closable: true,
      icon: 'table',
      env: 'Dev'
    },
    {
      id: 'events',
      title: 'events [[Dev][ReadOnly] 10.220.6.4]',
      type: 'table',
      closable: true,
      icon: 'table',
      env: 'Dev'
    },
    {
      id: 'connection_credential',
      title: 'connectio...credential [[Dev][ReadOnly] 10.220.6.4]',
      type: 'table',
      closable: true,
      icon: 'table',
      env: 'Dev'
    }
  ]);

  activeTabId = $state<string>('connection_credential');

  // Sidebar Layout
  sidebarWidth = $state<number>(290);
  servicesHeight = $state<number>(230);
  isSidebarCollapsed = $state<boolean>(false);
  treeFilterQuery = $state<string>('');

  // Tree expanded states
  treeExpanded = $state<Record<string, boolean>>({
    'dev_srv': true,
    'dev_db': true,
    'dev_schema_public': true,
    'dev_tables': true,
    'prd_srv': true,
    'prd_db': true,
    'prd_schema_public': true,
    'prd_tables': true
  });

  selectedTreeNode = $state<string>('connection_credential');

  // SQL Editor State
  sqlCode = $state<string>(`join content ct  1..n <-> 1: on ct.id = ma."contentID"
LEFT JOIN transfer_job_queue q ON q."jobID" = c.process_id
WHERE c.participant = 'media-transferer'
AND c.process_name   = 'media-ingest'
AND c.name           = 'create-ingest-job'
AND c.retry_count    = 0
AND c.schedule IS NULL
AND c."timestamp"    < now() - interval '5 minutes'  -- streamer poll 10s -> >5 phút chắc chắn không ai đọc
ORDER BY c.ordinal;

select * from transfer_job where id = 'b8262180-1075-43f8-8338-2412d4734d65'`);

  isExecutingQuery = $state<boolean>(false);
  activeExecutionLine = $state<number>(74);
  activeResultTab = $state<string>('prd_mh_asset.public.transfer_job');
  lastExecutionLatency = $state<string>('630 ms');
  cursorPos = $state<{ line: number; col: number }>({ line: 71, col: 42 });

  // UI Density & Typography Scale
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

  // SQL Result sets
  transferJobResults = $state([
    {
      id: 'b8262180-1075-43f8-8338-2412d4734d65',
      ordinal: 172098,
      connectionCredentialID: '753da7e1-523f-4834-ba74-7824e3d5aa52',
      sourcePath: '/home/vod/VOD/asset/hd_transcode_001.mp4',
      targetPath: 's3://vtv-media-cdn/ingest/hd_transcode_001.m3u8',
      status: 'COMPLETED',
      createdAt: '2024-11-14 02:15:20.124'
    },
    {
      id: '84a0d912-32b1-4f81-a901-b8471c981774',
      ordinal: 172099,
      connectionCredentialID: '753da7e1-523f-4834-ba74-7824e3d5aa52',
      sourcePath: '/home/vod/VOD/asset/4k_preview_reel.mov',
      targetPath: 's3://vtv-media-cdn/ingest/4k_preview_reel.m3u8',
      status: 'PROCESSING',
      createdAt: '2024-11-14 02:18:45.981'
    },
    {
      id: '3c92e105-1972-46bb-9372-8419db005612',
      ordinal: 172100,
      connectionCredentialID: 'ca1d2eee-6d07-48b5-af45-f551e984f3ed',
      sourcePath: '/home/vod/VOD/asset/news_live_capture.ts',
      targetPath: 's3://vtv-media-cdn/ingest/news_live_capture.m3u8',
      status: 'QUEUED',
      createdAt: '2024-11-14 02:20:00.005'
    },
    {
      id: 'd91a27cc-03f1-4328-98e1-55bb019488a0',
      ordinal: 172101,
      connectionCredentialID: 'ca1d2eee-6d07-48b5-af45-f551e984f3ed',
      sourcePath: '/home/vod/VOD/asset/promo_bumper_winter.mp4',
      targetPath: 's3://vtv-media-cdn/ingest/promo_bumper_winter.m3u8',
      status: 'COMPLETED',
      createdAt: '2024-11-14 02:21:14.332'
    },
    {
      id: 'e28401aa-7749-4bd2-a189-10fa487311c9',
      ordinal: 172102,
      connectionCredentialID: '753da7e1-523f-4834-ba74-7824e3d5aa52',
      sourcePath: '/home/vod/VOD/asset/master_dolby_atmos.mxf',
      targetPath: 's3://vtv-media-cdn/ingest/master_dolby_atmos.m3u8',
      status: 'COMPLETED',
      createdAt: '2024-11-14 02:23:40.890'
    },
    {
      id: 'f93c8801-44bb-4e63-b812-77c8e901a551',
      ordinal: 172103,
      connectionCredentialID: '9e768c19-b3e5-4027-80a9-25a83307fa14',
      sourcePath: '/home/vod/VOD/asset/sports_highlight_final.mp4',
      targetPath: 's3://vtv-media-cdn/ingest/sports_highlight_final.m3u8',
      status: 'FAILED_RETRYING',
      createdAt: '2024-11-14 02:25:01.401'
    }
  ]);

  // Table Data Viewer State (connection_credential matching design2.png)
  whereFilter = $state<string>('');
  orderByFilter = $state<string>('');
  selectedCell = $state<{ row: number; col: string } | null>({ row: 0, col: 'id' });
  tableSort = $state<{ col: string; dir: 'asc' | 'desc' | null }>({ col: '', dir: null });

  allCredentialRows = $state<CredentialRow[]>([
    {
      id: '2389a9f7-9ed7-4b7c-9b7c-7eec030b45f5',
      createdDate: '2024-10-29 10:06:03.247644',
      lastUpdatedDate: null,
      createdBy: '4ebfa335-512c-4731-b3b3-82b5faef431c'
    },
    {
      id: 'a4430e1f-3b89-4088-bed8-8a426ab293f0',
      createdDate: '2024-10-31 03:56:42.776700',
      lastUpdatedDate: null,
      createdBy: '0f379df2-47d2-47e7-b195-23c8a918f029'
    },
    {
      id: '3be52421-235b-412e-b420-23b612770315',
      createdDate: '2024-11-05 07:49:06.550088',
      lastUpdatedDate: null,
      createdBy: '0f379df2-47d2-47e7-b195-23c8a918f029'
    },
    {
      id: '6e98ccab-d08c-49af-bd10-70ca0d510879',
      createdDate: '2024-10-31 03:58:27.775837',
      lastUpdatedDate: '2024-10-31 04:05:44.633942',
      createdBy: '0f379df2-47d2-47e7-b195-23c8a918f029'
    },
    {
      id: 'ca1d2eee-6d07-48b5-af45-f551e984f3ed',
      createdDate: '2024-11-14 03:30:27.290163',
      lastUpdatedDate: null,
      createdBy: '0f379df2-47d2-47e7-b195-23c8a918f029'
    },
    {
      id: '9e768c19-b3e5-4027-80a9-25a83307fa14',
      createdDate: '2024-11-14 03:34:10.424030',
      lastUpdatedDate: null,
      createdBy: '0f379df2-47d2-47e7-b195-23c8a918f029'
    },
    {
      id: '5f1190bc-43e1-4560-a299-2810a91176b2',
      createdDate: '2024-11-15 01:12:00.129031',
      lastUpdatedDate: '2024-11-15 02:40:11.590211',
      createdBy: '1a89f902-771c-46a2-990a-55bc981240cc'
    },
    {
      id: '1c77bb32-901e-4501-83da-0012e84712aa',
      createdDate: '2024-11-16 09:20:44.810920',
      lastUpdatedDate: null,
      createdBy: '3b09cc14-411a-4d40-a190-67c822e11801'
    },
    {
      id: '88de91a0-5120-4100-bc77-90fa771801c1',
      createdDate: '2024-11-18 11:45:00.320110',
      lastUpdatedDate: '2024-11-18 12:01:29.440188',
      createdBy: '3b09cc14-411a-4d40-a190-67c822e11801'
    },
    {
      id: '7100ea21-0019-4882-a773-199bc488210e',
      createdDate: '2024-11-20 14:10:19.551020',
      lastUpdatedDate: null,
      createdBy: '4ebfa335-512c-4731-b3b3-82b5faef431c'
    },
    {
      id: 'b990145c-4472-4702-881a-66df90123841',
      createdDate: '2024-11-21 08:33:12.771900',
      lastUpdatedDate: null,
      createdBy: '0f379df2-47d2-47e7-b195-23c8a918f029'
    },
    {
      id: '448102ff-1188-4209-a110-ee1289551007',
      createdDate: '2024-11-22 16:50:41.112009',
      lastUpdatedDate: '2024-11-22 17:10:00.891223',
      createdBy: '1a89f902-771c-46a2-990a-55bc981240cc'
    }
  ]);

  // Derived filtered credential rows
  filteredCredentialRows = $derived.by(() => {
    let list = [...this.allCredentialRows];

    if (this.whereFilter.trim()) {
      const q = this.whereFilter.toLowerCase().trim();
      list = list.filter(row => {
        if (q.includes('is null')) {
          return row.lastUpdatedDate === null;
        }
        if (q.includes('is not null')) {
          return row.lastUpdatedDate !== null;
        }
        return (
          row.id.toLowerCase().includes(q) ||
          row.createdBy.toLowerCase().includes(q) ||
          row.createdDate.toLowerCase().includes(q) ||
          (row.lastUpdatedDate && row.lastUpdatedDate.toLowerCase().includes(q))
        );
      });
    }

    if (this.tableSort.col && this.tableSort.dir) {
      const { col, dir } = this.tableSort;
      list.sort((a: any, b: any) => {
        const valA = a[col] ?? '';
        const valB = b[col] ?? '';
        if (valA < valB) return dir === 'asc' ? -1 : 1;
        if (valA > valB) return dir === 'asc' ? 1 : -1;
        return 0;
      });
    }

    return list;
  });

  // Services panel sessions
  services = $state([
    { id: 's1', checked: true, name: 'connection_credential', duration: '1 s 159 ms', type: 'table' },
    { id: 's2', checked: true, name: 'events', duration: '986 ms', type: 'table' },
    { id: 's3', checked: true, name: 'commands', duration: '1 s 449 ms', type: 'table' },
    { id: 's4', checked: true, name: 'console_2', duration: '630 ms', type: 'console' }
  ]);

  // Console Logs
  consoleLogs = $state<ConsoleLogEntry[]>([
    {
      id: 'log-1',
      timestamp: '2026-10-03 10:50:34',
      type: 'info',
      summary: 'Connected to dev_mh_asset (PostgreSQL 16.2 on x86_64)'
    },
    {
      id: 'log-2',
      timestamp: '2026-10-03 10:50:34',
      connection: 'dev_mh_asset.public>',
      querySnippet: 'SELECT t.* FROM public.connection_credential t LIMIT 501',
      type: 'query',
      summary: 'SELECT query submitted to database engine'
    },
    {
      id: 'log-3',
      timestamp: '2026-10-03 10:50:34',
      type: 'success',
      summary: '58 rows retrieved starting from 1 in 586 ms (execution: 86 ms, fetching: 500 ms)'
    }
  ]);

  // MCP Gateway Approvals Queue — loaded from Admin API (see loadApprovals).
  approvals = $state<ApprovalRequest[]>([]);

  // MCP Gateway Audit Logs — loaded from Admin API (see loadAudit).
  auditLogs = $state<AuditRecord[]>([]);

  // MCP Gateway Connections — loaded from Admin API (see loadConnections).
  connections = $state<ConnectionItem[]>([]);

  // MCP Tokens — loaded from Admin API (see loadTokens).
  tokens = $state<McpToken[]>([]);

  // Dialog & Modal Triggers
  searchModalOpen = $state<boolean>(false);
  settingsModalOpen = $state<boolean>(false);
  ddlModalOpen = $state<boolean>(false);
  createTokenModalOpen = $state<boolean>(false);
  createdTokenSecret = $state<string | null>(null);
  selectedApprovalForAction = $state<ApprovalRequest | null>(null);
  approvalDecisionMode = $state<'approve' | 'reject'>('approve');
  approvalDecisionReason = $state<string>('');

  // Notifications
  notificationCount = $derived(
    this.approvals.filter(a => a.state === 'pending').length
  );

  // Tab operations
  openTab(tab: TabItem) {
    const existing = this.tabs.find(t => t.id === tab.id);
    if (!existing) {
      this.tabs = [...this.tabs, tab];
    }
    this.activeTabId = tab.id;
  }

  closeTab(tabId: string) {
    const idx = this.tabs.findIndex(t => t.id === tabId);
    if (idx !== -1) {
      const isClosingActive = this.activeTabId === tabId;
      this.tabs = this.tabs.filter(t => t.id !== tabId);
      if (isClosingActive && this.tabs.length > 0) {
        this.activeTabId = this.tabs[Math.max(0, idx - 1)].id;
      }
    }
  }

  // Query Execution Simulation
  executeQuery() {
    this.isExecutingQuery = true;
    const startTime = performance.now();

    setTimeout(() => {
      this.isExecutingQuery = false;
      const duration = Math.round(performance.now() - startTime) + 32;
      this.lastExecutionLatency = `${duration} ms`;

      // Update services list
      const s = this.services.find(i => i.name === 'console_2');
      if (s) {
        s.duration = `${duration} ms`;
      }

      // Append to console log
      const now = new Date().toISOString().replace('T', ' ').substring(0, 19);
      this.consoleLogs.push({
        id: 'log-' + Date.now(),
        timestamp: now,
        connection: 'prd_mh_asset.public>',
        querySnippet: 'select * from transfer_job where id = \'b8262180-1075-43f8-8338-2412d4734d65\'',
        type: 'query',
        summary: 'Query completed successfully.'
      });

      this.consoleLogs.push({
        id: 'log-' + (Date.now() + 1),
        timestamp: now,
        type: 'success',
        summary: `6 rows retrieved starting from 1 in ${duration} ms (execution: 12 ms, fetching: ${duration - 12} ms)`
      });
    }, 400);
  }

  // ---- Admin API wiring ----
  // Per user decision: error states are shown, never silently mocked.

  dataLoading = $state<boolean>(true);
  dataError = $state<string | null>(null);

  summary = $state<DashboardSummary | null>(null);

  private stopLive: (() => void) | null = null;
  private pollTimer: ReturnType<typeof setInterval> | null = null;

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
    this.approvals = page.items.map((t) => this.mapApproval(t));
  }

  async loadAudit() {
    const page = await listAudit({ limit: 100 });
    this.auditLogs = page.items.map((e) => this.mapAudit(e));
  }

  async loadConnections() {
    const { items } = await listConnections();
    const prev = new Map(this.connections.map((c) => [c.name, c.latency_ms]));
    this.connections = items.map((c) => ({ ...this.mapConnection(c), latency_ms: prev.get(c.name) }));
  }

  async loadTokens() {
    const { items } = await listTokens();
    this.tokens = items.map((t) => this.mapToken(t));
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
    } catch (e) {
      this.dataError = e instanceof ApiError ? e.message : String(e);
    } finally {
      this.dataLoading = false;
    }
  }

  /** Live approval updates via SSE, with 10s polling fallback. */
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

  // Approve / Reject actions (async, Admin API backed).
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

  // Token actions (async, Admin API backed).
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

  // Connection ping (async, Admin API backed).
  async pingConnection(name: string) {
    const res = await testConnection(name);
    const conn = this.connections.find((c) => c.name === name);
    if (conn) {
      conn.status = res.status;
      conn.latency_ms = res.latency_ms;
      conn.last_ping_at = new Date().toISOString().replace('T', ' ').substring(0, 19);
    }
    return res;
  }
}

export const appState = new AppStateManager();
