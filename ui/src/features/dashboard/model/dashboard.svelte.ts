import type { DashboardSummary } from '@/entities/session';
import { getSummary } from '../api/dashboard';

class DashboardManager {
  summary = $state<DashboardSummary | null>(null);

  async loadSummary() {
    this.summary = await getSummary();
  }
}

export const dashboardState = new DashboardManager();
