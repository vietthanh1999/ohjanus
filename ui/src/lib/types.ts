export type TabType = 
  | 'console' 
  | 'table' 
  | 'approvals' 
  | 'audit' 
  | 'tokens' 
  | 'connections' 
  | 'dashboard';

export interface Tab {
  id: string;
  title: string;
  subtitle?: string;
  icon: TabType;
  type: TabType;
  closable: boolean;
  metadata?: {
    table?: string;
    connection?: string;
    schema?: string;
  };
}

export interface DatabaseNode {
  id: string;
  name: string;
  type: 'server' | 'database' | 'schema' | 'group' | 'table' | 'routine';
  badge?: string;
  expanded?: boolean;
  iconType?: string;
  children?: DatabaseNode[];
  targetTable?: string;
  connection?: string;
}

export interface ServicesItem {
  id: string;
  name: string;
  type: 'server' | 'session' | 'console';
  latency?: string;
  checked?: boolean;
  expanded?: boolean;
  children?: ServicesItem[];
}

export type StatementType = 'SELECT' | 'INSERT' | 'UPDATE' | 'DELETE' | 'DROP' | 'ALTER' | 'MERGE';

export interface ApprovalItem {
  id: string;
  connection: string;
  statementType: StatementType;
  sql: string;
  params: Record<string, any> | string;
  affectedEstimate: number;
  requestedBy: {
    tokenId: string;
    client: string;
  };
  state: 'pending' | 'approved' | 'rejected' | 'expired';
  createdAt: string;
  expiresAt: string;
  warnings: string[];
  plan?: string;
  decisionReason?: string;
  decidedBy?: string;
  decidedAt?: string;
}

export interface AuditEvent {
  id: string;
  ts: string;
  event: string;
  client: string;
  connection: string;
  tool: string;
  sqlNormalized: string;
  statementType: StatementType;
  durationMs: number;
  rowCount: number;
  status: 'SUCCESS' | 'DENIED' | 'FAILED';
  policyRule?: string;
  policyDecision?: 'ALLOW' | 'DENY' | 'REQUIRE_APPROVAL';
  paramsRedacted?: string;
  error?: string;
}

export interface ConnectionInfo {
  name: string;
  host: string;
  driver: 'postgres' | 'mysql' | 'sqlite';
  readonly: boolean;
  status: 'healthy' | 'degraded' | 'unreachable';
  lastPingAt: string;
  latencyMs: number;
  pool: {
    open: number;
    idle: number;
    inUse: number;
    max: number;
  };
  version: string;
  allowedSchemas: string[];
  deniedTables: string[];
}

export interface McpToken {
  id: string;
  name: string;
  scopes: string[];
  createdAt: string;
  expiresAt: string;
  lastUsedAt: string;
  state: 'active' | 'revoked' | 'expired';
  rawToken?: string;
}
