import { approvalsState } from '@/features/approvals';
import { auditState } from '@/features/audit';
import { connectionsState } from '@/features/connections';
import { consoleState } from '@/features/console';
import { dashboardState } from '@/features/dashboard';
import { explorerState } from '@/features/explorer';
import { settingsState } from '@/features/settings';
import { tokensState } from '@/features/tokens';
import { workbenchState } from '@/features/workbench';
import { ApiError } from '@/shared/api';

/**
 * Shell-only boot orchestration: initial Admin API load + full reload after
 * the gateway endpoint changes (SettingsModal `onReconnect`).
 * Loading/error display is per-slice (each manager owns its flags); boot
 * only coordinates and logs the outcome.
 */
class BootManager {
  async loadAll() {
    settingsState.initDensity();
    try {
      await Promise.all([
        approvalsState.loadApprovals(),
        auditState.loadAudit(),
        connectionsState.loadConnections(),
        tokensState.loadTokens(),
        dashboardState.loadSummary()
      ]);
      consoleState.ensureDefaultConnection();
      workbenchState.pushLog({ type: 'info', summary: `Connected: ${connectionsState.connections.length} connection(s) loaded from Admin API.` });
      // Preload schema for the default connection so the explorer is live
      // on first paint; other connections load lazily on expand.
      if (consoleState.console.connection) {
        await explorerState.loadSchema(consoleState.console.connection).catch(() => {});
      }
    } catch (e) {
      const msg = e instanceof ApiError ? e.message : String(e);
      workbenchState.pushLog({ type: 'error', summary: `Admin API unreachable: ${msg}` });
    }
  }
}

export const bootState = new BootManager();
