import { mapApproval, type ApprovalRequest } from '@/entities/approval';
import { workbenchState } from '@/features/workbench';
import { APPROVAL_PAGE_LIMIT, APPROVAL_POLL_MS } from '@/shared/config';
import { auditState } from '@/features/audit';
import { listApprovals, approveApproval, rejectApproval, openApprovalStream } from '../api/approvals';

class ApprovalsManager {
  approvals = $state<ApprovalRequest[]>([]);
  selectedApprovalForAction = $state<ApprovalRequest | null>(null);
  approvalDecisionMode = $state<'approve' | 'reject'>('approve');
  approvalDecisionReason = $state<string>('');

  notificationCount = $derived(
    this.approvals.filter((a) => a.state === 'pending').length
  );

  async loadApprovals() {
    const page = await listApprovals({ limit: APPROVAL_PAGE_LIMIT });
    this.approvals = (page.items ?? []).map((t) => mapApproval(t));
    this.syncWorkbenchBadge();
  }

  private syncWorkbenchBadge() {
    const apprTab = workbenchState.tabs.find((t) => t.id === 'approvals');
    if (apprTab) {
      const pendingCount = this.approvals.filter((a) => a.state === 'pending').length;
      apprTab.badge = pendingCount > 0 ? String(pendingCount) : undefined;
    }
  }

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
    await auditState.loadAudit().catch(() => {});
    this.selectedApprovalForAction = null;
    this.approvalDecisionReason = '';
  }

  /** Live approval updates via SSE, with polling fallback. */
  private stopLive: (() => void) | null = null;
  private pollTimer: ReturnType<typeof setInterval> | null = null;

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
      this.pollTimer = setInterval(() => void poll(), APPROVAL_POLL_MS);
    } catch {
      this.pollTimer = setInterval(() => void poll(), APPROVAL_POLL_MS);
    }
  }

  stopLiveUpdates() {
    this.stopLive?.();
    this.stopLive = null;
  }
}

export const approvalsState = new ApprovalsManager();
