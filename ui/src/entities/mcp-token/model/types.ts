import type { ApiToken } from '@/shared/api';

/** UI-facing MCP token metadata. The raw secret is never stored — shown once at creation. */
export interface McpToken {
  id: string;
  name: string;
  scopes: string[];
  created_at: string;
  expires_at: string;
  last_used_at: string;
  state: 'active' | 'revoked' | 'expired';
  rawToken?: string;
}

export function mapToken(t: ApiToken): McpToken {
  return {
    id: t.id,
    name: t.name,
    scopes: t.scopes,
    created_at: t.created_at ?? '—',
    expires_at: t.expires_at ?? '—',
    last_used_at: t.last_used_at ?? '—',
    state: t.state
  };
}
