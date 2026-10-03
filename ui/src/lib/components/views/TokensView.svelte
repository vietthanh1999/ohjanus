<script lang="ts">
  import { appStore } from '../../appStore.svelte';
  import type { McpToken } from '../../types';

  let showCreateModal = $state(false);
  let showSecretModal = $state(false);
  let generatedToken = $state<McpToken | null>(null);
  let copied = $state(false);

  // Form State
  let tokenName = $state('');
  let ttlDays = $state(30);
  let scopes = $state<Record<string, boolean>>({
    read: true,
    write_preview: true,
    write_execute: false,
    admin: false
  });

  function handleCreate() {
    const selectedScopes = Object.keys(scopes).filter(k => scopes[k]);
    if (!tokenName.trim()) return;

    generatedToken = appStore.createToken(tokenName, selectedScopes, ttlDays);
    showCreateModal = false;
    showSecretModal = true;
    tokenName = '';
    copied = false;
  }

  function copySecret() {
    if (generatedToken?.rawToken) {
      navigator.clipboard.writeText(generatedToken.rawToken);
      copied = true;
      setTimeout(() => copied = false, 2000);
    }
  }
</script>

<div class="tokens-view">
  <!-- Header -->
  <div class="view-header">
    <div>
      <span class="header-icon">🔑</span>
      <span class="view-title">MCP Client Tokens Management</span>
      <span class="view-desc">Issue and manage secure opaque tokens for AI agents (Cursor, Claude Desktop, custom bots)</span>
    </div>
    <button class="btn-primary" onclick={() => showCreateModal = true}>
      + Generate New Token
    </button>
  </div>

  <!-- Tokens List Table -->
  <div class="table-container font-mono">
    <table class="data-grid">
      <thead>
        <tr>
          <th>TOKEN ID</th>
          <th>AGENT / CLIENT NAME</th>
          <th>SCOPES</th>
          <th>CREATED</th>
          <th>EXPIRES</th>
          <th>LAST USED</th>
          <th>STATE</th>
          <th style="text-align: right;">ACTIONS</th>
        </tr>
      </thead>
      <tbody>
        {#each appStore.tokens as tok}
          <tr>
            <td class="id-cell">{tok.id}</td>
            <td class="name-cell font-sans">
              <strong>{tok.name}</strong>
            </td>
            <td class="scopes-cell">
              {#each tok.scopes as sc}
                <span class="scope-tag scope-{sc}">{sc}</span>
              {/each}
            </td>
            <td>{tok.createdAt}</td>
            <td>{tok.expiresAt}</td>
            <td>{tok.lastUsedAt}</td>
            <td>
              <span class="state-pill state-{tok.state}">{tok.state.toUpperCase()}</span>
            </td>
            <td style="text-align: right;">
              {#if tok.state === 'active'}
                <button 
                  class="btn-danger btn-sm"
                  onclick={() => appStore.revokeToken(tok.id)}
                >
                  Revoke
                </button>
              {:else}
                <span style="color: var(--text-muted); font-size: 11px;">Revoked</span>
              {/if}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>

<!-- Create Token Modal -->
{#if showCreateModal}
  <div class="modal-backdrop" onclick={() => showCreateModal = false}>
    <div class="modal-dialog animate-fade-in" onclick={(e) => e.stopPropagation()}>
      <div class="modal-header">
        <span>Generate New MCP Token</span>
        <button class="icon-btn" onclick={() => showCreateModal = false}>✕</button>
      </div>

      <div class="modal-body">
        <div class="form-group">
          <label for="tok-name">Agent / Client Name:</label>
          <input 
            id="tok-name"
            type="text" 
            placeholder="e.g., Cursor IDE - Production Investigator" 
            bind:value={tokenName}
            style="width: 100%;"
          />
        </div>

        <div class="form-group" style="margin-top: 12px;">
          <label>Allowed Janus Scopes (§1.2):</label>
          <div class="scopes-grid">
            <label class="chk-label">
              <input type="checkbox" bind:checked={scopes.read} />
              <span><code>read</code> (db_list_connections, db_schema, db_read, db_explain)</span>
            </label>
            <label class="chk-label">
              <input type="checkbox" bind:checked={scopes.write_preview} />
              <span><code>write_preview</code> (db_write_preview — dry-run simulation)</span>
            </label>
            <label class="chk-label">
              <input type="checkbox" bind:checked={scopes.write_execute} />
              <span><code>write_execute</code> (db_write_execute — requires human approval)</span>
            </label>
            <label class="chk-label">
              <input type="checkbox" bind:checked={scopes.admin} />
              <span><code>admin</code> (manage tokens &amp; connection configurations)</span>
            </label>
          </div>
        </div>

        <div class="form-group" style="margin-top: 12px;">
          <label for="tok-ttl">Token Validity (TTL Days):</label>
          <select id="tok-ttl" bind:value={ttlDays} style="width: 100%;">
            <option value={7}>7 Days</option>
            <option value={30}>30 Days (Recommended)</option>
            <option value={90}>90 Days</option>
            <option value={365}>1 Year</option>
          </select>
        </div>
      </div>

      <div class="modal-footer">
        <button class="btn-secondary" onclick={() => showCreateModal = false}>Cancel</button>
        <button class="btn-primary" onclick={handleCreate} disabled={!tokenName.trim()}>
          Generate Token
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- One-Time Secret Reveal Modal -->
{#if showSecretModal && generatedToken}
  <div class="modal-backdrop">
    <div class="modal-dialog animate-fade-in" style="width: 520px;">
      <div class="modal-header">
        <span>🎉 MCP Token Generated Successfully</span>
      </div>

      <div class="modal-body">
        <div class="alert-box">
          ⚠️ <strong>WARNING:</strong> This token secret is displayed only <strong>ONCE</strong> and cannot be retrieved later. Copy and paste it into your client configuration now.
        </div>

        <div style="margin: 12px 0 6px; font-weight: 600;">Opaque Bearer Token:</div>
        <div class="secret-box font-mono">
          <span>{generatedToken.rawToken}</span>
          <button class="copy-btn" onclick={copySecret}>
            {copied ? '✓ Copied' : 'Copy'}
          </button>
        </div>

        <div style="margin-top: 12px; font-size: 11px; color: var(--text-muted);">
          Example Cursor config (<code>~/.cursor/mcp.json</code>):
          <pre class="mcp-config-snippet font-mono">{`{
  "mcpServers": {
    "ohjanus": {
      "command": "janus",
      "args": ["serve", "--token", "${generatedToken.rawToken}"]
    }
  }
}`}</pre>
        </div>
      </div>

      <div class="modal-footer">
        <button class="btn-primary" onclick={() => showSecretModal = false}>
          I Have Safely Saved the Token
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .tokens-view {
    display: flex;
    flex-direction: column;
    height: 100%;
    background: var(--surface-canvas);
  }

  .view-header {
    height: 48px;
    background: var(--surface-toolbar);
    border-bottom: 1px solid var(--border-default);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 16px;
  }

  .header-icon { font-size: 16px; margin-right: 8px; }
  .view-title { font-weight: 600; font-size: 13px; margin-right: 12px; }
  .view-desc { font-size: 11px; color: var(--text-muted); }

  .table-container {
    flex: 1;
    overflow: auto;
  }

  .data-grid {
    width: 100%;
    border-collapse: collapse;
    font-size: 12px;
  }

  .data-grid th {
    background: var(--surface-table-header);
    height: var(--table-header-height);
    border-right: 1px solid var(--border-default);
    border-bottom: 1px solid var(--border-default);
    padding: 0 10px;
    text-align: left;
    color: var(--text-secondary);
    position: sticky;
    top: 0;
  }

  .data-grid td {
    height: 32px;
    border-right: 1px solid var(--border-subtle);
    border-bottom: 1px solid var(--border-subtle);
    padding: 0 10px;
  }

  .scopes-cell {
    display: flex;
    align-items: center;
    gap: 4px;
    height: 32px;
  }

  .scope-tag {
    font-size: 10px;
    padding: 1px 5px;
    border-radius: 3px;
    background: #25272A;
    border: 1px solid var(--border-default);
    color: #DFE1E5;
  }

  .scope-write_execute { border-color: var(--action-warning); color: var(--action-warning); }
  .scope-admin { border-color: var(--action-danger); color: var(--action-danger); }

  .state-pill {
    font-size: 9px;
    font-weight: 700;
    padding: 2px 6px;
    border-radius: 8px;
  }

  .state-active { background: rgba(87, 211, 140, 0.2); color: var(--action-success); }
  .state-revoked { background: rgba(229, 83, 83, 0.2); color: var(--action-danger); }

  .btn-sm {
    padding: 2px 8px;
    font-size: 11px;
  }

  .scopes-grid {
    display: flex;
    flex-direction: column;
    gap: 6px;
    margin-top: 6px;
    background: #18191B;
    padding: 8px;
    border-radius: 4px;
    border: 1px solid var(--border-default);
  }

  .chk-label {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 11px;
    cursor: pointer;
  }

  .alert-box {
    background: rgba(237, 162, 0, 0.15);
    border: 1px solid var(--action-warning);
    padding: 8px 12px;
    border-radius: 4px;
    color: #FFE082;
    font-size: 11px;
    line-height: 1.4;
  }

  .secret-box {
    background: #151618;
    border: 1px solid var(--border-accent);
    padding: 8px 12px;
    border-radius: 4px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 13px;
    color: var(--action-success);
  }

  .copy-btn {
    background: var(--surface-hover);
    border: 1px solid var(--border-default);
    color: #fff;
    padding: 3px 8px;
    border-radius: 3px;
    font-size: 11px;
  }

  .mcp-config-snippet {
    background: #151618;
    border: 1px solid var(--border-default);
    padding: 8px;
    border-radius: 4px;
    margin-top: 4px;
  }
</style>
