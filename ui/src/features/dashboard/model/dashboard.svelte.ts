import type { DashboardSummary } from '@/entities/session';
import { ApiError } from '@/shared/api';
import { getSummary } from '../api/dashboard';

class DashboardManager {
  summary = $state<DashboardSummary | null>(null);
  loading = $state<boolean>(false);
  error = $state<string | null>(null);

  async loadSummary() {
    this.loading = true;
    this.error = null;
    try {
      this.summary = await getSummary();
    } catch (e) {
      this.error = e instanceof ApiError ? e.message : String(e);
      throw e;
    } finally {
      this.loading = false;
    }
  }
}

export const dashboardState = new DashboardManager();
