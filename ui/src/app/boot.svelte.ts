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
 * Global boot orchestration: initial Admin API load for the App shell.
 * Views read `dataLoading`/`dataError` for their loading/error states and
 * call `loadAll()` to retry. Only the shell (and the deprecated facade)
 * trigger it.
 */
class BootManager {
  dataLoading = $state<boolean>(true);
  dataError = $state<string | null>(null);

  async loadAll() {
    settingsState.initDensity();
    this.dataLoading = true;
    this.dataError = null;
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
      this.dataError = e instanceof ApiError ? e.message : String(e);
      workbenchState.pushLog({ type: 'error', summary: `Admin API unreachable: ${this.dataError}` });
    } finally {
      this.dataLoading = false;
    }
  }
}

export const bootState = new BootManager();
