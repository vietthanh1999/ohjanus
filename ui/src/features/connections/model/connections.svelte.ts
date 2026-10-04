import { mapConnection, type ConnectionItem } from '@/entities/connection';
import { nowStamp } from '@/shared/lib';
import { listConnections, testConnection } from '../api/connections';

class ConnectionsManager {
  connections = $state<ConnectionItem[]>([]);

  async loadConnections() {
    const { items } = await listConnections();
    const prev = new Map(this.connections.map((c) => [c.name, c.latency_ms]));
    this.connections = (items ?? []).map((c) => ({ ...mapConnection(c), latency_ms: prev.get(c.name) }));
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
}

export const connectionsState = new ConnectionsManager();
