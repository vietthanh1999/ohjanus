import type { Tab, ApprovalItem, AuditEvent, ConnectionInfo, McpToken, DatabaseNode, ServicesItem } from './types';
import {
  initialTreeData,
  initialServicesData,
  initialApprovals,
  initialAuditEvents,
  initialConnections,
  initialTokens,
  sampleTableRows,
  sampleConsoleResultRows
} from './mockData';

class AppStore {
  // Tabs
  tabs = $state<Tab[]>([
    {
      id: 'tab-console-1',
      title: 'console',
      subtitle: '[PRD] 10.250.6.23',
      icon: 'console',
      type: 'console',
      closable: true,
      metadata: { connection: '[PRD] 10.250.6.23' }
    },
    {
      id: 'tab-console-2',
      title: 'console_2',
      subtitle: '[PRD] 10.250.6.23',
      icon: 'console',
      type: 'console',
      closable: true,
      metadata: { connection: '[PRD] 10.250.6.23' }
    },
    {
      id: 'tab-table-conn',
      title: 'connectio...credential',
      subtitle: '[Dev][ReadOnly] 10.220.6.4',
      icon: 'table',
      type: 'table',
      closable: true,
      metadata: { connection: '[Dev][ReadOnly] 10.220.6.4', table: 'connection_credential', schema: 'public' }
    },
    {
      id: 'tab-approvals',
      title: 'Approvals',
      subtitle: 'Write Queue',
      icon: 'approvals',
      type: 'approvals',
      closable: true
    },
    {
      id: 'tab-audit',
      title: 'Audit Log',
      subtitle: 'Security Gateway',
      icon: 'audit',
      type: 'audit',
      closable: true
    }
  ]);

  activeTabId = $state<string>('tab-console-2');

  // Sidebar tree and services
  treeData = $state<DatabaseNode[]>(initialTreeData);
  servicesData = $state<ServicesItem[]>(initialServicesData);
  selectedTreeNodeId = $state<string>('tbl-conn-cred');

  // Gateway data
  approvals = $state<ApprovalItem[]>(initialApprovals);
  auditEvents = $state<AuditEvent[]>(initialAuditEvents);
  connections = $state<ConnectionInfo[]>(initialConnections);
  tokens = $state<McpToken[]>(initialTokens);

  // Table Viewer state
  tableRows = $state(sampleTableRows);
  whereClause = $state<string>('');
  orderByClause = $state<string>('');
  pageSize = $state<number>(58);

  // SQL Console state
  consoleSql = $state<string>(`-- Active database: prd_mh_asset.public
join content ct  1..n <-> 1: on ct.id = ma."contentID"
LEFT JOIN transfer_job_queue q ON q."jobID" = c.process_id
WHERE c.participant = 'media-transferer'
AND c.process_name   = 'media-ingest'
AND c.name           = 'create-ingest-job'
AND c.retry_count    = 0
AND c.schedule IS NULL
AND c."timestamp"    < now() - interval '5 minutes'  -- streamer poll 10s -> >5 phút chắc chắn không ai đọc
ORDER BY c.ordinal;

select * from transfer_job where id = 'b8262180-1075-43f8-8338-2412d4734d65'`);

  consoleResults = $state(sampleConsoleResultRows);
  selectedResultTab = $state<string>('prd_mh_asset.public.transfer_job 2');
  isQueryRunning = $state<boolean>(false);
  activeExecutionBlock = $state<boolean>(true);

  // Console Logs
  consoleLogs = $state<string[]>([
    '[2026-10-03 10:50:34] Connected to dev_mh_asset',
    '[2026-10-03 10:50:34] dev_mh_asset.public> SELECT t.* FROM public.connection_credential t LIMIT 501',
    '[2026-10-03 10:50:34] 58 rows retrieved starting from 1 in 586 ms (execution: 86 ms, fetching: 500 ms)'
  ]);

  // Global status
  cursorPosition = $state('71:42 (2954 chars, 73 line breaks)');
  notificationsCount = $state(3);

  // Derived state
  activeTab = $derived(this.tabs.find(t => t.id === this.activeTabId) || this.tabs[0]);
  pendingApprovalsCount = $derived(this.approvals.filter(a => a.state === 'pending').length);

  // Methods
  selectTab(id: string) {
    this.activeTabId = id;
  }

  closeTab(id: string) {
    const idx = this.tabs.findIndex(t => t.id === id);
    if (idx === -1) return;
    this.tabs.splice(idx, 1);
    if (this.activeTabId === id && this.tabs.length > 0) {
      this.activeTabId = this.tabs[Math.max(0, idx - 1)].id;
    }
  }

  openTab(tab: Tab) {
    const existing = this.tabs.find(t => t.id === tab.id);
    if (existing) {
      this.activeTabId = existing.id;
    } else {
      this.tabs.push(tab);
      this.activeTabId = tab.id;
    }
  }

  openTableTab(tableName: string, connection: string) {
    const id = `tab-table-${tableName}`;
    this.openTab({
      id,
      title: tableName.length > 14 ? tableName.slice(0, 10) + '...' + tableName.slice(-4) : tableName,
      subtitle: connection,
      icon: 'table',
      type: 'table',
      closable: true,
      metadata: { table: tableName, connection, schema: 'public' }
    });
  }

  openConsoleTab(name: string, connection: string) {
    const id = `tab-console-${name}`;
    this.openTab({
      id,
      title: name,
      subtitle: connection,
      icon: 'console',
      type: 'console',
      closable: true,
      metadata: { connection }
    });
  }

  approveRequest(id: string, reason: string) {
    const item = this.approvals.find(a => a.id === id);
    if (!item) return;
    item.state = 'approved';
    item.decidedBy = 'thanhtran@vtvprime.vn';
    item.decidedAt = new Date().toISOString();
    item.decisionReason = reason || 'Approved manually via OhJanus Admin UI';

    // Add audit event
    this.auditEvents.unshift({
      id: `aud-${Date.now().toString().slice(-4)}`,
      ts: new Date().toISOString().replace('T', ' ').slice(0, 19),
      event: 'write.approved',
      client: item.requestedBy.client,
      connection: item.connection,
      tool: 'db_write_execute',
      statementType: item.statementType,
      sqlNormalized: item.sql.replace(/\s+/g, ' ').slice(0, 80),
      durationMs: 32,
      rowCount: item.affectedEstimate,
      status: 'SUCCESS',
      policyRule: 'manual-approval',
      policyDecision: 'ALLOW'
    });
  }

  rejectRequest(id: string, reason: string) {
    const item = this.approvals.find(a => a.id === id);
    if (!item) return;
    item.state = 'rejected';
    item.decidedBy = 'thanhtran@vtvprime.vn';
    item.decidedAt = new Date().toISOString();
    item.decisionReason = reason || 'Rejected by administrator';

    // Add audit event
    this.auditEvents.unshift({
      id: `aud-${Date.now().toString().slice(-4)}`,
      ts: new Date().toISOString().replace('T', ' ').slice(0, 19),
      event: 'write.rejected',
      client: item.requestedBy.client,
      connection: item.connection,
      tool: 'db_write_execute',
      statementType: item.statementType,
      sqlNormalized: item.sql.replace(/\s+/g, ' ').slice(0, 80),
      durationMs: 12,
      rowCount: 0,
      status: 'DENIED',
      policyRule: 'manual-approval',
      policyDecision: 'DENY',
      error: `ADMIN_REJECTED: ${reason}`
    });
  }

  runQuery(sql?: string) {
    this.isQueryRunning = true;
    const query = sql || this.consoleSql;
    const now = new Date().toISOString().replace('T', ' ').slice(0, 19);

    setTimeout(() => {
      this.isQueryRunning = false;
      const latency = Math.floor(Math.random() * 400 + 40);
      const rows = Math.floor(Math.random() * 20 + 1);

      this.consoleLogs.push(
        `[${now}] prd_mh_asset.public> ${query.slice(0, 60)}...`,
        `[${now}] ${rows} rows retrieved in ${latency} ms (execution: ${Math.floor(latency * 0.2)} ms, fetching: ${Math.floor(latency * 0.8)} ms)`
      );

      // Add to audit
      this.auditEvents.unshift({
        id: `aud-${Date.now().toString().slice(-4)}`,
        ts: now,
        event: 'query.executed',
        client: 'Console 2 (Local Operator)',
        connection: '[PRD] 10.250.6.23',
        tool: 'db_read',
        statementType: 'SELECT',
        sqlNormalized: query.replace(/\s+/g, ' ').slice(0, 80),
        durationMs: latency,
        rowCount: rows,
        status: 'SUCCESS',
        policyRule: 'allow-select',
        policyDecision: 'ALLOW'
      });
    }, 450);
  }

  createToken(name: string, scopes: string[], ttlDays: number) {
    const id = `tok-${Date.now().toString().slice(-4)}`;
    const secret = `jn_${Math.random().toString(36).substring(2, 15)}${Math.random().toString(36).substring(2, 15)}${Math.random().toString(36).substring(2, 15)}`;
    const expDate = new Date();
    expDate.setDate(expDate.getDate() + ttlDays);

    const newToken: McpToken = {
      id,
      name,
      scopes,
      createdAt: new Date().toISOString().slice(0, 10),
      expiresAt: expDate.toISOString().slice(0, 10),
      lastUsedAt: 'Never',
      state: 'active',
      rawToken: secret
    };

    this.tokens.unshift(newToken);
    return newToken;
  }

  revokeToken(id: string) {
    const t = this.tokens.find(tok => tok.id === id);
    if (t) {
      t.state = 'revoked';
    }
  }

  testConnection(name: string) {
    const conn = this.connections.find(c => c.name === name);
    if (!conn) return;
    const latency = Math.floor(Math.random() * 25 + 8);
    conn.lastPingAt = 'Just now';
    conn.latencyMs = latency;
    conn.status = 'healthy';
  }
}

export const appStore = new AppStore();
