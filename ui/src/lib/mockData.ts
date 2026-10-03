import type { DatabaseNode, ServicesItem, ApprovalItem, AuditEvent, ConnectionInfo, McpToken } from './types';

export const initialTreeData: DatabaseNode[] = [
  {
    id: 'srv-dev-3',
    name: '[Dev][ReadOnly] 10.220.6.3',
    type: 'server',
    badge: '13',
    expanded: false,
    children: [
      { id: 'db-media-3', name: 'media_analytics', type: 'database', badge: '4' },
      { id: 'db-cms-3', name: 'cms_core', type: 'database', badge: '12' },
    ]
  },
  {
    id: 'srv-dev-4',
    name: '[Dev][ReadOnly] 10.220.6.4',
    type: 'server',
    badge: '15',
    expanded: true,
    connection: '[Dev][ReadOnly] 10.220.6.4',
    children: [
      { id: 'db-datahub', name: 'datahub', type: 'database', badge: '...' },
      { id: 'db-dev-auth', name: 'dev_dhb_auth', type: 'database', badge: '...' },
      { id: 'db-dev-goal', name: 'dev_goalert', type: 'database', badge: '...' },
      {
        id: 'db-dev-mh',
        name: 'dev_mh_asset',
        type: 'database',
        badge: '6',
        expanded: true,
        children: [
          { id: 'sch-es', name: 'es', type: 'schema' },
          { id: 'sch-info', name: 'information_schema', type: 'schema' },
          { id: 'sch-marts', name: 'marts', type: 'schema' },
          { id: 'sch-pg', name: 'pg_catalog', type: 'schema' },
          {
            id: 'sch-pm',
            name: 'pm',
            type: 'schema',
            expanded: false,
            children: [
              { id: 'pm-tables', name: 'tables', type: 'group', badge: '3' },
              { id: 'pm-routines', name: 'routines', type: 'group', badge: '2' },
            ]
          },
          {
            id: 'sch-public',
            name: 'public',
            type: 'schema',
            expanded: true,
            children: [
              {
                id: 'grp-tables',
                name: 'tables',
                type: 'group',
                badge: '33',
                expanded: true,
                children: [
                  { id: 'tbl-cat', name: 'category', type: 'table', targetTable: 'category' },
                  { id: 'tbl-conn-cred', name: 'connection_credential', type: 'table', targetTable: 'connection_credential' },
                  { id: 'tbl-cnt', name: 'content', type: 'table', targetTable: 'content' },
                  { id: 'tbl-cnt-dist', name: 'content_distribution', type: 'table', targetTable: 'content_distribution' },
                  { id: 'tbl-cnt-prov', name: 'content_provider_configuration', type: 'table', targetTable: 'content_provider_configuration' },
                  { id: 'tbl-trans-job', name: 'transfer_job', type: 'table', targetTable: 'transfer_job' },
                  { id: 'tbl-trans-queue', name: 'transfer_job_queue', type: 'table', targetTable: 'transfer_job_queue' },
                ]
              }
            ]
          }
        ]
      }
    ]
  },
  {
    id: 'srv-prd-23',
    name: '[PRD] 10.250.6.23',
    type: 'server',
    badge: '18',
    expanded: true,
    connection: '[PRD] 10.250.6.23',
    children: [
      { id: 'prd-db-asset', name: 'prd_mh_asset', type: 'database', badge: '8', expanded: true, children: [
        { id: 'prd-sch-pub', name: 'public', type: 'schema', expanded: true, children: [
          { id: 'prd-tbl-users', name: 'user_entitlements', type: 'table', targetTable: 'user_entitlements' },
          { id: 'prd-tbl-jobs', name: 'transfer_job', type: 'table', targetTable: 'transfer_job' },
        ]}
      ]}
    ]
  }
];

export const initialServicesData: ServicesItem[] = [
  {
    id: 'svc-db',
    name: 'Database',
    type: 'server',
    expanded: true,
    children: [
      {
        id: 'svc-dev4',
        name: '[Dev][ReadOnly] 10.220.6.4',
        type: 'server',
        expanded: true,
        children: [
          { id: 'svc-conn-cred', name: 'connection_credential', type: 'session', latency: '1 s 159 ms' },
          { id: 'svc-events', name: 'events', type: 'session', latency: '986 ms' },
          { id: 'svc-commands', name: 'commands', type: 'session', latency: '1 s 449 ms' },
        ]
      },
      {
        id: 'svc-prd23',
        name: '[PRD] 10.250.6.23',
        type: 'server',
        expanded: true,
        children: [
          { id: 'svc-console', name: 'console', type: 'console' },
          { id: 'svc-console-2', name: 'console_2', type: 'console', latency: '630 ms' },
        ]
      }
    ]
  }
];

export const sampleTableRows = [
  {
    id: '2389a9f7-9ed7-4b7c-9b7c-7eec030b45f5',
    createdDate: '2024-10-29 10:06:03.247644',
    lastUpdatedDate: null,
    createdBy: '4ebfa335-5656-4052-9001-f112da909c21'
  },
  {
    id: 'a4430e1f-3b89-4088-bed8-8a426ab293f0',
    createdDate: '2024-10-31 03:56:42.776700',
    lastUpdatedDate: null,
    createdBy: '0f379df2-47d2-47e7-b195-2acfa39a8bc4'
  },
  {
    id: '3be52421-235b-412e-b420-23b612770315',
    createdDate: '2024-11-05 07:49:06.550088',
    lastUpdatedDate: null,
    createdBy: '0f379df2-47d2-47e7-b195-2acfa39a8bc4'
  },
  {
    id: '6e98ccab-d08c-49af-bd10-70ca0d510879',
    createdDate: '2024-10-31 03:58:27.775837',
    lastUpdatedDate: '2024-10-31 04:05:44.633942',
    createdBy: '0f379df2-47d2-47e7-b195-2acfa39a8bc4'
  },
  {
    id: 'ca1d2eee-6d07-48b5-af45-f551e984f3ed',
    createdDate: '2024-11-14 03:30:27.290163',
    lastUpdatedDate: null,
    createdBy: '0f379df2-47d2-47e7-b195-2acfa39a8bc4'
  },
  {
    id: '9e768c19-b3e5-4027-80a9-25a83307fa14',
    createdDate: '2024-11-14 03:34:10.424030',
    lastUpdatedDate: null,
    createdBy: '0f379df2-47d2-47e7-b195-2acfa39a8bc4'
  },
  {
    id: '09a9bcb6-cb58-47fe-bc61-e3f416acd3f6',
    createdDate: '2024-11-05 07:30:56.146845',
    lastUpdatedDate: '2024-11-05 07:31:08.435963',
    createdBy: '0f379df2-47d2-47e7-b195-2acfa39a8bc4'
  },
  {
    id: '29ebafe4-3484-4595-a2a5-37158b8cd0b5',
    createdDate: '2024-10-31 04:07:39.501833',
    lastUpdatedDate: '2024-11-05 07:40:31.640001',
    createdBy: '0f379df2-47d2-47e7-b195-2acfa39a8bc4'
  },
  {
    id: '53d2afe2-1043-4eb2-b488-1fdfa096d870',
    createdDate: '2024-11-14 03:36:12.190281',
    lastUpdatedDate: '2024-11-18 10:30:42.374558',
    createdBy: '0f379df2-47d2-47e7-b195-2acfa39a8bc4'
  }
];

export const sampleConsoleResultRows = [
  {
    id: 'b8262180-1075-43f8-8338-2412d4734d65',
    ordinal: 172098,
    connectionCredentialID: '753da7e1-523f-4834-ba74-7824e3d5aa52',
    sourcePath: '/home/vod/VOD/assets/2024/11/video_hd_098.mp4',
    status: 'COMPLETED',
    duration_ms: 1450
  },
  {
    id: 'c9371291-2186-59a9-9449-3523e5845e76',
    ordinal: 172099,
    connectionCredentialID: '753da7e1-523f-4834-ba74-7824e3d5aa52',
    sourcePath: '/home/vod/VOD/assets/2024/11/stream_manifest_4k.m3u8',
    status: 'IN_PROGRESS',
    duration_ms: 320
  }
];

export const initialApprovals: ApprovalItem[] = [
  {
    id: 'appr-9824f1',
    connection: '[PRD] 10.250.6.23',
    statementType: 'UPDATE',
    sql: `UPDATE user_entitlements\nSET status = 'ACTIVE', tier = 'VIP_PREMIUM'\nWHERE user_id = $1 AND renewal_date > NOW();`,
    params: { "$1": "usr_99812480" },
    affectedEstimate: 1,
    requestedBy: {
      tokenId: 'jn_c9f28a9b182e01',
      client: 'Cursor AI Agent v0.42.3'
    },
    state: 'pending',
    createdAt: '2026-10-03T11:40:12Z',
    expiresAt: '2026-10-03T11:45:12Z',
    warnings: [
      'Modifies user permissions on Production environment',
      'Requires explicit sign-off per security policy §4.3'
    ],
    plan: 'Index Scan using idx_user_entitlements_uid on user_entitlements (cost=0.29..8.31 rows=1 width=72)',
  },
  {
    id: 'appr-77a3d9',
    connection: '[Dev][ReadOnly] 10.220.6.4',
    statementType: 'DELETE',
    sql: `DELETE FROM transfer_job_queue\nWHERE retry_count >= 5\n  AND created_at < NOW() - INTERVAL '7 days';`,
    params: {},
    affectedEstimate: 42,
    requestedBy: {
      tokenId: 'jn_a10b98f244109c',
      client: 'Claude Desktop v1.2'
    },
    state: 'pending',
    createdAt: '2026-10-03T11:36:00Z',
    expiresAt: '2026-10-03T11:51:00Z',
    warnings: [
      'Bulk deletion: 42 rows will be removed from transfer_job_queue',
      'Queue consumers might miss failed task history'
    ],
    plan: 'Seq Scan on transfer_job_queue (cost=0.00..12.50 rows=42 width=0)'
  },
  {
    id: 'appr-33c91e',
    connection: '[PRD] 10.250.6.23',
    statementType: 'DROP',
    sql: `DROP TABLE legacy_content_backup;`,
    params: {},
    affectedEstimate: 0,
    requestedBy: {
      tokenId: 'jn_b441ca1209774a',
      client: 'Migration Script (CI/CD)'
    },
    state: 'pending',
    createdAt: '2026-10-03T11:43:00Z',
    expiresAt: '2026-10-03T11:46:00Z',
    warnings: [
      'CRITICAL DANGER: DDL operation DROP TABLE in production!',
      'Policy deny-ddl flagged this query as high risk'
    ],
    plan: 'DROP TABLE legacy_content_backup;'
  },
  {
    id: 'appr-55b201',
    connection: '[PRD] 10.250.6.23',
    statementType: 'INSERT',
    sql: `INSERT INTO content_distribution (content_id, target_cdn, is_active)\nVALUES ('cnt_99014', 'cdn-akamai-vn', true);`,
    params: {},
    affectedEstimate: 1,
    requestedBy: {
      tokenId: 'jn_c9f28a9b182e01',
      client: 'Cursor AI Agent v0.42.3'
    },
    state: 'approved',
    createdAt: '2026-10-03T10:15:00Z',
    expiresAt: '2026-10-03T10:20:00Z',
    decidedBy: 'thanhtran@vtvprime.vn',
    decidedAt: '2026-10-03T10:17:22Z',
    decisionReason: 'Verified CDN configuration parameters with network team.',
    warnings: []
  }
];

export const initialAuditEvents: AuditEvent[] = [
  {
    id: 'aud-1092',
    ts: '2026-10-03 11:42:08',
    event: 'query.executed',
    client: 'Cursor Agent (jn_c9f28a)',
    connection: '[PRD] 10.250.6.23',
    tool: 'db_read',
    statementType: 'SELECT',
    sqlNormalized: 'SELECT id, ordinal, status FROM transfer_job WHERE id = $1 LIMIT 1',
    durationMs: 42,
    rowCount: 1,
    status: 'SUCCESS',
    policyRule: 'allow-select',
    policyDecision: 'ALLOW'
  },
  {
    id: 'aud-1091',
    ts: '2026-10-03 11:40:12',
    event: 'write.preview_created',
    client: 'Cursor Agent (jn_c9f28a)',
    connection: '[PRD] 10.250.6.23',
    tool: 'db_write_preview',
    statementType: 'UPDATE',
    sqlNormalized: 'UPDATE user_entitlements SET status = $1 WHERE user_id = $2',
    durationMs: 18,
    rowCount: 0,
    status: 'SUCCESS',
    policyRule: 'warn-write',
    policyDecision: 'REQUIRE_APPROVAL'
  },
  {
    id: 'aud-1090',
    ts: '2026-10-03 11:38:45',
    event: 'query.denied',
    client: 'Claude Desktop (jn_a10b98)',
    connection: '[PRD] 10.250.6.23',
    tool: 'db_read',
    statementType: 'SELECT',
    sqlNormalized: 'SELECT pg_read_file($1)',
    durationMs: 4,
    rowCount: 0,
    status: 'DENIED',
    policyRule: 'banned-functions',
    policyDecision: 'DENY',
    error: 'ACCESS_DENIED: Execution of banned system function pg_read_file is prohibited'
  },
  {
    id: 'aud-1089',
    ts: '2026-10-03 11:35:10',
    event: 'query.executed',
    client: 'DataGrip UI (Console 2)',
    connection: '[Dev][ReadOnly] 10.220.6.4',
    tool: 'db_read',
    statementType: 'SELECT',
    sqlNormalized: 'SELECT t.* FROM public.connection_credential t LIMIT 501',
    durationMs: 86,
    rowCount: 58,
    status: 'SUCCESS',
    policyRule: 'allow-select',
    policyDecision: 'ALLOW'
  }
];

export const initialConnections: ConnectionInfo[] = [
  {
    name: '[Dev][ReadOnly] 10.220.6.4',
    host: '10.220.6.4:5432',
    driver: 'postgres',
    readonly: true,
    status: 'healthy',
    lastPingAt: 'Just now',
    latencyMs: 12,
    pool: {
      open: 10,
      idle: 7,
      inUse: 3,
      max: 20
    },
    version: 'PostgreSQL 16.4 on aarch64-unknown-linux-gnu',
    allowedSchemas: ['public', 'pm', 'marts', 'es'],
    deniedTables: ['secrets', 'admin_credentials', 'payment_keys']
  },
  {
    name: '[PRD] 10.250.6.23',
    host: '10.250.6.23:5432',
    driver: 'postgres',
    readonly: false,
    status: 'healthy',
    lastPingAt: '12s ago',
    latencyMs: 28,
    pool: {
      open: 16,
      idle: 11,
      inUse: 5,
      max: 30
    },
    version: 'PostgreSQL 16.2 (Debian 16.2-1.pgdg120+1)',
    allowedSchemas: ['public'],
    deniedTables: ['user_salts', 'audit_vault']
  },
  {
    name: '[Dev][ReadOnly] 10.220.6.3',
    host: '10.220.6.3:5432',
    driver: 'postgres',
    readonly: true,
    status: 'healthy',
    lastPingAt: '1m ago',
    latencyMs: 15,
    pool: {
      open: 8,
      idle: 6,
      inUse: 2,
      max: 15
    },
    version: 'PostgreSQL 15.6',
    allowedSchemas: ['media_analytics', 'cms_core'],
    deniedTables: []
  }
];

export const initialTokens: McpToken[] = [
  {
    id: 'tok-01',
    name: 'Cursor AI - Thanh Tran Laptop',
    scopes: ['read', 'write_preview', 'write_execute'],
    createdAt: '2026-09-20',
    expiresAt: '2026-11-20',
    lastUsedAt: '2 minutes ago',
    state: 'active'
  },
  {
    id: 'tok-02',
    name: 'Claude Desktop - Data Analyst Team',
    scopes: ['read'],
    createdAt: '2026-09-25',
    expiresAt: '2026-10-25',
    lastUsedAt: '15 minutes ago',
    state: 'active'
  },
  {
    id: 'tok-03',
    name: 'VTVPrime CI Deployment Bot',
    scopes: ['read', 'write_preview', 'admin'],
    createdAt: '2026-08-01',
    expiresAt: '2026-12-31',
    lastUsedAt: '1 hour ago',
    state: 'active'
  }
];
