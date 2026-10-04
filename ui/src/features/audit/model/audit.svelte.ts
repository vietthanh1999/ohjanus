import { mapAuditEvent, type AuditRecord } from '@/entities/audit-record';
import { AUDIT_PAGE_LIMIT } from '@/shared/config';
import { listAudit } from '../api/audit';

class AuditManager {
  auditLogs = $state<AuditRecord[]>([]);

  async loadAudit() {
    const page = await listAudit({ limit: AUDIT_PAGE_LIMIT });
    this.auditLogs = (page.items ?? []).map((e) => mapAuditEvent(e));
  }
}

export const auditState = new AuditManager();
