<script lang="ts">
  import { appState } from '../../state/appState.svelte';

  let tokenName = $state('');
  let selectedScopes = $state<string[]>(['read', 'schema']);
  let hasCopied = $state(false);

  const availableScopes = [
    { id: 'read', label: 'read (db_read, SELECT)', desc: 'Allow read-only queries' },
    { id: 'write_with_approval', label: 'write_with_approval', desc: 'Allows write queries after human approval' },
    { id: 'schema', label: 'schema (db_schema, introspection)', desc: 'Inspect table definitions and columns' },
    { id: 'explain', label: 'explain (db_explain)', desc: 'Run EXPLAIN on query plans' },
    { id: 'admin', label: 'admin (full gateway config)', desc: 'Full administration scopes' }
  ];

  function toggleScope(id: string) {
    if (selectedScopes.includes(id)) {
      selectedScopes = selectedScopes.filter(s => s !== id);
    } else {
      selectedScopes = [...selectedScopes, id];
    }
  }

  function handleCreate() {
    if (!tokenName.trim()) return;
    appState.createToken(tokenName.trim(), selectedScopes);
  }

  function copySecret() {
    if (appState.createdTokenSecret) {
      navigator.clipboard.writeText(appState.createdTokenSecret);
      hasCopied = true;
      setTimeout(() => hasCopied = false, 2000);
    }
  }

  function close() {
    appState.createTokenModalOpen = false;
    appState.createdTokenSecret = null;
    tokenName = '';
    selectedScopes = ['read', 'schema'];
  }
</script>

{#if appState.createTokenModalOpen}
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div class="jb-modal-backdrop" onclick={close}>
    <div class="jb-modal" style="width: 480px;" onclick={(e) => e.stopPropagation()}>
      <div class="jb-modal-header">
        <span>Create MCP Bearer Access Token</span>
        <button class="jb-icon-btn" onclick={close}>✕</button>
      </div>

      <div class="jb-modal-body">
        {#if !appState.createdTokenSecret}
          <div class="form-group">
            <label for="tok-name" class="form-label">Client Name / Description <span style="color: var(--action-danger);">*</span></label>
            <input
              id="tok-name"
              type="text"
              class="form-input"
              placeholder="e.g. Claude Desktop Agent - Finance Sync"
              bind:value={tokenName}
            />
          </div>

          <div class="form-group">
            <div class="form-label">Token Scopes</div>
            <div class="scopes-list">
              {#each availableScopes as sc}
                <!-- svelte-ignore a11y_click_events_have_key_events -->
                <!-- svelte-ignore a11y_no_static_element_interactions -->
                <div class="scope-row" onclick={() => toggleScope(sc.id)}>
                  <input
                    type="checkbox"
                    checked={selectedScopes.includes(sc.id)}
                    class="scope-checkbox"
                  />
                  <div class="scope-text">
                    <span class="sc-name code-text">{sc.label}</span>
                    <span class="sc-desc">{sc.desc}</span>
                  </div>
                </div>
              {/each}
            </div>
          </div>
        {:else}
          <!-- Success Token Reveal View -->
          <div class="token-reveal-view">
            <div class="alert-box">
              <span style="font-size: 18px;">⚠️</span>
              <div>
                <strong>Save this token key now!</strong>
                <p style="font-size: 11px; margin-top: 2px;">Janus stores only cryptographic hashes in its vault. You will not be able to view this plaintext secret again.</p>
              </div>
            </div>

            <div class="secret-box code-text">
              <span>{appState.createdTokenSecret}</span>
            </div>

            <button class="copy-btn" onclick={copySecret}>
              {hasCopied ? '✓ Copied to Clipboard!' : '📋 Copy Token to Clipboard'}
            </button>
          </div>
        {/if}
      </div>

      <div class="jb-modal-footer">
        {#if !appState.createdTokenSecret}
          <button class="jb-btn-secondary" onclick={close}>Cancel</button>
          <button class="jb-btn-primary" onclick={handleCreate} disabled={!tokenName.trim()}>
            Generate Token
          </button>
        {:else}
          <button class="jb-btn-primary" onclick={close}>
            Done
          </button>
        {/if}
      </div>
    </div>
  </div>
{/if}

<style>
  .form-group {
    display: flex;
    flex-direction: column;
    gap: 6px;
    margin-bottom: 14px;
  }

  .form-label {
    font-size: 11px;
    font-weight: 500;
    color: var(--text-secondary);
  }

  .form-input {
    background-color: #1E1F22;
    border: 1px solid var(--border-default);
    border-radius: 3px;
    padding: 6px 8px;
    color: var(--text-primary);
    font-size: 12px;
  }

  .form-input:focus {
    border-color: var(--border-accent);
  }

  .scopes-list {
    background-color: #1E1F22;
    border: 1px solid var(--border-default);
    border-radius: 4px;
    overflow: hidden;
  }

  .scope-row {
    display: flex;
    align-items: flex-start;
    gap: 8px;
    padding: 8px 10px;
    border-bottom: 1px solid var(--border-subtle);
    cursor: pointer;
    transition: background-color 0.1s ease;
  }

  .scope-row:last-child {
    border-bottom: none;
  }

  .scope-row:hover {
    background-color: var(--bg-hover);
  }

  .scope-checkbox {
    margin-top: 2px;
    accent-color: var(--action-primary);
  }

  .scope-text {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .sc-name {
    font-size: 11px;
    color: var(--text-primary);
    font-weight: 500;
  }

  .sc-desc {
    font-size: 10px;
    color: var(--text-muted);
  }

  /* Token Reveal */
  .token-reveal-view {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .alert-box {
    background-color: rgba(237, 162, 0, 0.1);
    border: 1px solid rgba(237, 162, 0, 0.4);
    padding: 10px;
    border-radius: 4px;
    display: flex;
    gap: 10px;
    align-items: center;
    color: #ffc44d;
  }

  .secret-box {
    background-color: #1E1F22;
    border: 1px solid #3574F0;
    border-radius: 4px;
    padding: 10px;
    font-size: 13px;
    color: #57D38C;
    word-break: break-all;
    user-select: all;
  }

  .copy-btn {
    padding: 8px 14px;
    background-color: var(--action-primary);
    color: #FFFFFF;
    border-radius: 4px;
    font-weight: 500;
    font-size: 12px;
    text-align: center;
    transition: background-color 0.15s ease;
  }

  .copy-btn:hover {
    background-color: var(--action-primary-hover);
  }
</style>
