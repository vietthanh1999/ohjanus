import { mapToken, type McpToken } from '@/entities/mcp-token';
import { ApiError } from '@/shared/api';
import { listTokens, createToken as apiCreateToken, revokeToken as apiRevokeToken } from '../api/tokens';

class TokensManager {
  tokens = $state<McpToken[]>([]);
  loading = $state<boolean>(false);
  error = $state<string | null>(null);
  createTokenModalOpen = $state<boolean>(false);
  /** Raw secret shown exactly once after creation; never persisted. */
  createdTokenSecret = $state<string | null>(null);

  async loadTokens() {
    this.loading = true;
    this.error = null;
    try {
      const { items } = await listTokens();
      this.tokens = (items ?? []).map((t) => mapToken(t));
    } catch (e) {
      this.error = e instanceof ApiError ? e.message : String(e);
      throw e;
    } finally {
      this.loading = false;
    }
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
