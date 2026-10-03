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
  status: 'healthy' | 'degraded' | 'error';
  last_ping_at: string;
  latency_ms: number;
  pool: {
    open: number;
    idle: number;
    in_use: number;
    max: number;
  };
}

export interface McpToken {
  id: string;
  name: string;
  scopes: string[];
  created_at: string;
  expires_at: string;
  last_used_at: string;
  state: 'active' | 'revoked' | 'expired';
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
    },
    {
      id: 'console_2',
      title: 'console_2 [[PRD] 10.250.6.23]',
      type: 'console',
      closable: true,
      icon: 'lightning',
      env: 'PRD'
    },
    {
      id: 'approvals',
      title: 'Approvals Queue',
      type: 'approvals',
      closable: true,
      icon: 'shield',
      badge: '2'
    },
    {
      id: 'audit',
      title: 'Audit Logs',
      type: 'audit',
      closable: true,
      icon: 'audit'
    },
    {
      id: 'tokens',
      title: 'MCP Tokens',
      type: 'tokens',
      closable: true,
      icon: 'key'
    },
    {
      id: 'connections',
      title: 'Connection Pools',
      type: 'connections',
      closable: true,
      icon: 'database'
    },
    {
      id: 'dashboard',
      title: 'Gateway Dashboard',
      type: 'dashboard',
      closable: true,
      icon: 'chart'
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

  // MCP Gateway Approvals Queue
  approvals = $state<ApprovalRequest[]>([
    {
      id: 'appr-9812',
      state: 'pending',
      connection: 'prd_mh_asset',
      statement_type: 'UPDATE',
      sql: `UPDATE user_accounts
SET active = false,
    deactivation_reason = 'Dormant account auto-sweep'
WHERE last_login < NOW() - INTERVAL '365 days'
  AND active = true;`,
      affected_estimate: 1420,
      requested_by: {
        client: 'Cursor IDE Agent (DevOps)',
        token_id: 'jn_agent_cursor_992'
      },
      warnings: [
        'Batch UPDATE on production database',
        'Impact estimation exceeds threshold (> 1,000 rows)'
      ],
      created_at: '2026-10-03 11:20:00',
      expires_at: '2026-10-03 12:20:00'
    },
    {
      id: 'appr-9813',
      state: 'pending',
      connection: 'dev_mh_asset',
      statement_type: 'DROP',
      sql: `DROP TABLE deprecated_ingest_temp_2024;`,
      affected_estimate: 0,
      requested_by: {
        client: 'Claude Desktop Agent (DBA)',
        token_id: 'jn_agent_claude_104'
      },
      warnings: [
        'DDL operation: irreversible drop of database object',
        'Cannot be rolled back automatically'
      ],
      created_at: '2026-10-03 11:32:10',
      expires_at: '2026-10-03 12:02:10'
    },
    {
      id: 'appr-9810',
      state: 'approved',
      connection: 'prd_mh_asset',
      statement_type: 'INSERT',
      sql: `INSERT INTO category (id, code, name, created_at)
VALUES (gen_random_uuid(), 'TECH_NEWS', 'Technology & AI News', NOW());`,
      affected_estimate: 1,
      requested_by: {
        client: 'Content Sync Service',
        token_id: 'jn_agent_sync_552'
      },
      warnings: [],
      created_at: '2026-10-03 10:14:00',
      expires_at: '2026-10-03 11:14:00',
      decided_by: 'thanhtran',
      decided_at: '2026-10-03 10:18:22',
      decision_reason: 'Verified new category code is authorized by content team.'
    },
    {
      id: 'appr-9807',
      state: 'rejected',
      connection: 'prd_mh_asset',
      statement_type: 'DELETE',
      sql: `DELETE FROM transfer_job WHERE status = 'FAILED';`,
      affected_estimate: 890,
      requested_by: {
        client: 'Cleanup Cron Agent',
        token_id: 'jn_agent_cleanup_001'
      },
      warnings: ['Destructive DELETE without date filter'],
      created_at: '2026-10-03 09:05:00',
      expires_at: '2026-10-03 10:05:00',
      decided_by: 'thanhtran',
      decided_at: '2026-10-03 09:12:00',
      decision_reason: 'Rejected. Failed jobs must be preserved for 30 days audit retention.'
    }
  ]);

  // MCP Gateway Audit Logs
  auditLogs = $state<AuditRecord[]>([
    {
      id: 'aud-001',
      ts: '2026-10-03 11:42:15',
      event: 'query.executed',
      request_id: 'req_88192a',
      token_id: 'jn_agent_cursor_992',
      client: 'Cursor Agent (DevOps)',
      connection: 'prd_mh_asset',
      statement_type: 'SELECT',
      tables: ['transfer_job'],
      policy_decision: 'ALLOW',
      policy_rule: 'allow-select-all',
      row_count: 6,
      duration_ms: 86,
      status: 'OK',
      sql_normalized: 'SELECT * FROM transfer_job WHERE id = ?'
    },
    {
      id: 'aud-002',
      ts: '2026-10-03 11:32:10',
      event: 'approval.requested',
      request_id: 'req_88190c',
      token_id: 'jn_agent_claude_104',
      client: 'Claude Desktop Agent (DBA)',
      connection: 'dev_mh_asset',
      statement_type: 'DROP',
      tables: ['deprecated_ingest_temp_2024'],
      policy_decision: 'REQUIRE_APPROVAL',
      policy_rule: 'require-approval-ddl',
      row_count: 0,
      duration_ms: 12,
      status: 'OK',
      sql_normalized: 'DROP TABLE deprecated_ingest_temp_2024'
    },
    {
      id: 'aud-003',
      ts: '2026-10-03 11:20:00',
      event: 'approval.requested',
      request_id: 'req_88184f',
      token_id: 'jn_agent_cursor_992',
      client: 'Cursor Agent (DevOps)',
      connection: 'prd_mh_asset',
      statement_type: 'UPDATE',
      tables: ['user_accounts'],
      policy_decision: 'REQUIRE_APPROVAL',
      policy_rule: 'require-approval-production-write',
      row_count: 0,
      duration_ms: 18,
      status: 'OK',
      sql_normalized: 'UPDATE user_accounts SET active = ? WHERE last_login < ?'
    },
    {
      id: 'aud-004',
      ts: '2026-10-03 10:55:01',
      event: 'query.denied',
      request_id: 'req_88172d',
      token_id: 'jn_agent_guest_301',
      client: 'Anonymous Agent',
      connection: 'prd_mh_asset',
      statement_type: 'SELECT',
      tables: ['pg_shadow'],
      policy_decision: 'DENY',
      policy_rule: 'deny-system-catalogs',
      row_count: 0,
      duration_ms: 4,
      status: 'DENIED',
      sql_normalized: 'SELECT * FROM pg_shadow'
    },
    {
      id: 'aud-005',
      ts: '2026-10-03 10:50:34',
      event: 'query.executed',
      request_id: 'req_88165b',
      token_id: 'jn_agent_local_admin',
      client: 'OhJanus UI Console',
      connection: 'dev_mh_asset',
      statement_type: 'SELECT',
      tables: ['connection_credential'],
      policy_decision: 'ALLOW',
      policy_rule: 'allow-select-all',
      row_count: 58,
      duration_ms: 86,
      status: 'OK',
      sql_normalized: 'SELECT t.* FROM public.connection_credential t LIMIT 501'
    }
  ]);

  // MCP Gateway Connections
  connections = $state<ConnectionItem[]>([
    {
      name: 'dev_mh_asset',
      driver: 'PostgreSQL 16.2',
      readonly: true,
      status: 'healthy',
      last_ping_at: '2026-10-03 11:43:00',
      latency_ms: 12,
      pool: { open: 12, idle: 8, in_use: 4, max: 25 }
    },
    {
      name: 'prd_mh_asset',
      driver: 'PostgreSQL 16.2',
      readonly: false,
      status: 'healthy',
      last_ping_at: '2026-10-03 11:43:10',
      latency_ms: 18,
      pool: { open: 35, idle: 22, in_use: 13, max: 50 }
    },
    {
      name: 'analytics_clickhouse',
      driver: 'ClickHouse 24.3',
      readonly: true,
      status: 'healthy',
      last_ping_at: '2026-10-03 11:42:45',
      latency_ms: 34,
      pool: { open: 8, idle: 6, in_use: 2, max: 20 }
    }
  ]);

  // MCP Tokens
  tokens = $state<McpToken[]>([
    {
      id: 'jn_agent_cursor_992',
      name: 'Cursor IDE Dev Agent',
      scopes: ['read', 'write_with_approval', 'explain'],
      created_at: '2026-09-15 08:00:00',
      expires_at: '2026-12-15 08:00:00',
      last_used_at: '2026-10-03 11:42:15',
      state: 'active'
    },
    {
      id: 'jn_agent_claude_104',
      name: 'Claude Desktop DBA Assistant',
      scopes: ['read', 'write_with_approval', 'schema'],
      created_at: '2026-09-20 10:30:00',
      expires_at: '2026-12-20 10:30:00',
      last_used_at: '2026-10-03 11:32:10',
      state: 'active'
    },
    {
      id: 'jn_agent_sync_552',
      name: 'Automated Content Sync Worker',
      scopes: ['read', 'write_with_approval'],
      created_at: '2026-10-01 00:00:00',
      expires_at: '2027-01-01 00:00:00',
      last_used_at: '2026-10-03 10:14:00',
      state: 'active'
    },
    {
      id: 'jn_agent_temp_eval',
      name: 'Security Audit Scanner (Temporary)',
      scopes: ['read'],
      created_at: '2026-10-01 14:00:00',
      expires_at: '2026-10-02 14:00:00',
      last_used_at: '2026-10-02 13:58:00',
      state: 'expired'
    }
  ]);

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

  // Approve / Reject actions
  confirmApprovalAction() {
    if (!this.selectedApprovalForAction) return;

    const item = this.approvals.find(a => a.id === this.selectedApprovalForAction!.id);
    if (item) {
      item.state = this.approvalDecisionMode === 'approve' ? 'approved' : 'rejected';
      item.decided_by = 'thanhtran';
      item.decided_at = new Date().toISOString().replace('T', ' ').substring(0, 19);
      item.decision_reason = this.approvalDecisionReason || (this.approvalDecisionMode === 'approve' ? 'Approved via Admin UI' : 'Rejected via Admin UI');

      // Update approval tab badge
      const apprTab = this.tabs.find(t => t.id === 'approvals');
      if (apprTab) {
        const pendingCount = this.approvals.filter(a => a.state === 'pending').length;
        apprTab.badge = pendingCount > 0 ? String(pendingCount) : undefined;
      }

      // Add to audit log
      this.auditLogs.unshift({
        id: 'aud-' + Date.now(),
        ts: item.decided_at,
        event: item.state === 'approved' ? 'approval.granted' : 'approval.rejected',
        request_id: 'req_' + item.id,
        token_id: item.requested_by.token_id,
        client: item.requested_by.client,
        connection: item.connection,
        statement_type: item.statement_type,
        tables: [item.statement_type],
        policy_decision: item.state === 'approved' ? 'ALLOW' : 'DENY',
        policy_rule: 'human-approval-decision',
        row_count: item.affected_estimate,
        duration_ms: 25,
        status: item.state === 'approved' ? 'OK' : 'DENIED',
        sql_normalized: item.sql.split('\n')[0]
      });
    }

    this.selectedApprovalForAction = null;
    this.approvalDecisionReason = '';
  }

  // Token actions
  createToken(name: string, scopes: string[]) {
    const randomHex = Array.from({ length: 32 }, () => Math.floor(Math.random() * 16).toString(16)).join('');
    const rawSecret = `jn_${randomHex}`;
    const id = `jn_tok_${Date.now().toString(36)}`;
    const now = new Date().toISOString().replace('T', ' ').substring(0, 19);

    const expires = new Date();
    expires.setDate(expires.getDate() + 90);
    const expiresStr = expires.toISOString().replace('T', ' ').substring(0, 19);

    this.tokens.unshift({
      id,
      name,
      scopes,
      created_at: now,
      expires_at: expiresStr,
      last_used_at: 'Never',
      state: 'active'
    });

    this.createdTokenSecret = rawSecret;
  }

  revokeToken(id: string) {
    const tok = this.tokens.find(t => t.id === id);
    if (tok) {
      tok.state = 'revoked';
    }
  }

  // Connection ping
  pingConnection(name: string) {
    const conn = this.connections.find(c => c.name === name);
    if (conn) {
      const start = performance.now();
      setTimeout(() => {
        conn.latency_ms = Math.floor(Math.random() * 15) + 8;
        conn.last_ping_at = new Date().toISOString().replace('T', ' ').substring(0, 19);
      }, 250);
    }
  }
}

export const appState = new AppStateManager();
