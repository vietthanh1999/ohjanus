<script lang="ts">
  import { appState } from '../../state/appState.svelte';

  let query = $state('');

  const items = [
    { title: 'connection_credential', category: 'Table', desc: 'dev_mh_asset.public', action: () => appState.activeTabId = 'connection_credential' },
    { title: 'console_2', category: 'Console', desc: 'prd_mh_asset.public query editor', action: () => appState.activeTabId = 'console_2' },
    { title: 'Approvals Queue', category: 'Gateway', desc: 'Review AI Agent write requests', action: () => appState.activeTabId = 'approvals' },
    { title: 'Audit Log Trail', category: 'Gateway', desc: 'Audit records and query logs', action: () => appState.activeTabId = 'audit' },
    { title: 'MCP Access Tokens', category: 'Gateway', desc: 'Manage agent tokens & scopes', action: () => appState.activeTabId = 'tokens' },
    { title: 'Connection Pools', category: 'Gateway', desc: 'Monitor database connections', action: () => appState.activeTabId = 'connections' },
    { title: 'Gateway Dashboard', category: 'Gateway', desc: 'Throughput and telemetry metrics', action: () => appState.activeTabId = 'dashboard' },
    { title: 'transfer_job', category: 'Table', desc: 'prd_mh_asset.public', action: () => appState.activeTabId = 'console_2' },
    { title: 'category', category: 'Table', desc: 'dev_mh_asset.public', action: () => appState.activeTabId = 'connection_credential' },
    { title: 'content', category: 'Table', desc: 'dev_mh_asset.public', action: () => appState.activeTabId = 'connection_credential' },
    { title: 'Generate DDL', category: 'Action', desc: 'Inspect current schema definition', action: () => appState.ddlModalOpen = true }
  ];

  let filtered = $derived.by(() => {
    if (!query.trim()) return items;
    const q = query.toLowerCase().trim();
    return items.filter(i =>
      i.title.toLowerCase().includes(q) ||
      i.desc.toLowerCase().includes(q) ||
      i.category.toLowerCase().includes(q)
    );
  });

  function selectItem(action: () => void) {
    action();
    appState.searchModalOpen = false;
    query = '';
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      appState.searchModalOpen = false;
    }
  }
</script>

<svelte:window onkeydown={handleKeydown} />

{#if appState.searchModalOpen}
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div class="jb-modal-backdrop" onclick={() => appState.searchModalOpen = false}>
    <div class="palette-modal" onclick={(e) => e.stopPropagation()}>
      <div class="palette-input-wrap">
        <svg width="15" height="15" viewBox="0 0 16 16" fill="var(--text-muted)">
          <path fill-rule="evenodd" d="M11.5 7a4.5 4.5 0 11-9 0 4.5 4.5 0 019 0zm-.82 4.74a6 6 0 111.06-1.06l3.04 3.04a.75.75 0 11-1.06 1.06l-3.04-3.04z"/>
        </svg>
        <!-- svelte-ignore a11y_autofocus -->
        <input
          type="text"
          class="palette-input"
          placeholder="Search Everywhere (Tables, Consoles, Gateway, Actions)..."
          bind:value={query}
          autofocus
        />
        <span class="esc-tag">ESC</span>
      </div>

      <div class="palette-results">
        {#each filtered as item}
          <!-- svelte-ignore a11y_click_events_have_key_events -->
          <!-- svelte-ignore a11y_no_static_element_interactions -->
          <div class="result-row" onclick={() => selectItem(item.action)}>
            <div class="row-left">
              <span class="category-badge">{item.category}</span>
              <span class="row-title">{item.title}</span>
            </div>
            <span class="row-desc">{item.desc}</span>
          </div>
        {/each}
      </div>
    </div>
  </div>
{/if}

<style>
  .palette-modal {
    width: 560px;
    background-color: var(--bg-sidebar);
    border: 1px solid var(--border-strong);
    border-radius: 6px;
    box-shadow: 0 12px 36px rgba(0, 0, 0, 0.75);
    overflow: hidden;
    display: flex;
    flex-direction: column;
    margin-top: -120px;
  }

  .palette-input-wrap {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px 14px;
    border-bottom: 1px solid var(--border-default);
    background-color: #25272A;
  }

  .palette-input {
    flex: 1;
    background: transparent;
    border: none;
    outline: none;
    font-size: 13px;
    color: var(--text-primary);
  }

  .esc-tag {
    font-size: 10px;
    background-color: #313438;
    color: var(--text-muted);
    padding: 2px 6px;
    border-radius: 3px;
  }

  .palette-results {
    max-height: 320px;
    overflow-y: auto;
    padding: 4px 0;
  }

  .result-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 8px 14px;
    cursor: pointer;
    transition: background-color 0.1s ease;
  }

  .result-row:hover {
    background-color: var(--bg-hover);
  }

  .row-left {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .category-badge {
    font-size: 9px;
    font-weight: 700;
    text-transform: uppercase;
    background-color: #1E1F22;
    border: 1px solid var(--border-default);
    padding: 1px 6px;
    border-radius: 3px;
    color: var(--text-secondary);
  }

  .row-title {
    font-size: 12px;
    font-weight: 500;
    color: var(--text-primary);
  }

  .row-desc {
    font-size: 11px;
    color: var(--text-muted);
  }
</style>
