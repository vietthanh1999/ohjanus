import { mapToken, type McpToken } from '@/entities/mcp-token';
import { listTokens, createToken as apiCreateToken, revokeToken as apiRevokeToken } from '../api/tokens';

class TokensManager {
  tokens = $state<McpToken[]>([]);
  createTokenModalOpen = $state<boolean>(false);
  /** Raw secret shown exactly once after creation; never persisted. */
  createdTokenSecret = $state<string | null>(null);

  async loadTokens() {
    const { items } = await listTokens();
    this.tokens = (items ?? []).map((t) => mapToken(t));
  }

  async createToken(name: string, scopes: string[], ttlDays: number) {
    const created = await apiCreateToken({ name, scopes, ttl_hours: ttlDays * 24 });
    await this.loadTokens();
    this.createdTokenSecret = created.token;
    return created.token;
  }

  async revokeToken(id: string) {
    await apiRevokeToken(id);
    await this.loadTokens();
  }
}

export const tokensState = new TokensManager();
