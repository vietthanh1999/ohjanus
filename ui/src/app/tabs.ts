import type { Component } from 'svelte';
import type { TabItem } from '@/entities/tab';
import { ApprovalsView } from '@/features/approvals';
import { AuditView } from '@/features/audit';
import { ConnectionsView } from '@/features/connections';
import { SqlConsoleView } from '@/features/console';
import { DashboardView } from '@/features/dashboard';
import { TableDataView } from '@/features/table-viewer';
import { TokensView } from '@/features/tokens';

/**
 * Tab registry: adding a page = adding one entry here.
 * No `{:else if}` chains in App.svelte (skill anti-pattern 5).
 */
const staticRegistry: Record<string, Component> = {
  approvals: ApprovalsView,
  audit: AuditView,
  connections: ConnectionsView,
  tokens: TokensView,
  dashboard: DashboardView
};

/** Resolve the view component for the focused tab (console/table are dynamic). */
export function componentFor(tab: TabItem | undefined): Component | undefined {
  if (!tab) return undefined;
  if (tab.type === 'console') return SqlConsoleView;
  if (tab.type === 'table') return TableDataView;
  return staticRegistry[tab.id] ?? ApprovalsView;
}
