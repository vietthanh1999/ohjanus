import { mapAuditEvent, type AuditRecord } from '@/entities/audit-record';
import { ApiError } from '@/shared/api';
import { AUDIT_PAGE_LIMIT } from '@/shared/config';
import { listAudit } from '../api/audit';

class AuditManager {
  auditLogs = $state<AuditRecord[]>([]);
  loading = $state<boolean>(false);
  error = $state<string | null>(null);

  async loadAudit() {
    this.loading = true;
    this.error = null;
    try {
      const page = await listAudit({ limit: AUDIT_PAGE_LIMIT });
      this.auditLogs = (page.items ?? []).map((e) => mapAuditEvent(e));
    } catch (e) {
      this.error = e instanceof ApiError ? e.message : String(e);
      throw e;
    } finally {
      this.loading = false;
    }
  }
}

export const auditState = new AuditManager();
