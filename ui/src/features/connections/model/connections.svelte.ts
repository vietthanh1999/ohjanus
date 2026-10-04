import { mapConnection, type ConnectionItem } from '@/entities/connection';
import { ApiError } from '@/shared/api';
import { nowStamp } from '@/shared/lib';
import { listConnections, testConnection, createConnection as apiCreateConnection, type CreateConnectionInput } from '../api/connections';

class ConnectionsManager {
  connections = $state<ConnectionItem[]>([]);
  loading = $state<boolean>(false);
  error = $state<string | null>(null);
  createConnectionModalOpen = $state<boolean>(false);

  async loadConnections() {
    this.loading = true;
    this.error = null;
    try {
      const { items } = await listConnections();
      const prev = new Map(this.connections.map((c) => [c.name, c.latency_ms]));
      this.connections = (items ?? []).map((c) => ({ ...mapConnection(c), latency_ms: prev.get(c.name) }));
    } catch (e) {
      this.error = e instanceof ApiError ? e.message : String(e);
      throw e;
    } finally {
      this.loading = false;
    }
  }

  async pingConnection(name: string) {
    const res = await testConnection(name);
    const conn = this.connections.find((c) => c.name === name);
    if (conn) {
      conn.status = res.status;
      conn.latency_ms = res.latency_ms;
      conn.last_ping_at = nowStamp();
    }
    return res;
  }

  async createConnection(input: CreateConnectionInput) {
    const created = await apiCreateConnection(input);
    await this.loadConnections().catch(() => {});
    return created;
  }
}

export const connectionsState = new ConnectionsManager();
