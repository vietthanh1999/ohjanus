<script lang="ts">
  import { appState } from '../../state/appState.svelte';
  import { DialogPrimitive, Badge, Kbd, Box, Flex } from '@ohjanus/ui';
  import { Icon } from '@ohjanus/icons';

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

  function close() {
    appState.searchModalOpen = false;
    query = '';
  }
</script>

<DialogPrimitive open={appState.searchModalOpen} onClose={close}>
  <Box class="palette-modal">
    <Flex class="palette-input-wrap" align="center" gap="10px">
      <Icon name="search" size={15} color="var(--text-muted)" />
      <!-- svelte-ignore a11y_autofocus -->
      <input
        type="text"
        class="palette-input"
        placeholder="Search Everywhere (Tables, Consoles, Gateway, Actions)..."
        bind:value={query}
        autofocus
      />
      <Kbd>ESC</Kbd>
    </Flex>

    <Box class="palette-results">
      {#each filtered as item}
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <Flex class="result-row" align="center" justify="between" onclick={() => selectItem(item.action)}>
          <Flex class="row-left" align="center" gap="8px">
            <Badge variant="outline" size="sm">{item.category}</Badge>
            <span class="row-title">{item.title}</span>
          </Flex>
          <span class="row-desc">{item.desc}</span>
        </Flex>
      {/each}
    </Box>
  </Box>
</DialogPrimitive>

<style>
  :global(.palette-modal) {
    width: 560px;
    background-color: var(--bg-sidebar, #2B2D30);
    border: 1px solid var(--border-strong, #43454A);
    border-radius: 6px;
    box-shadow: 0 12px 36px rgba(0, 0, 0, 0.75);
    overflow: hidden;
    display: flex;
    flex-direction: column;
    margin-top: -100px;
  }

  :global(.palette-input-wrap) {
    padding: 10px 14px;
    border-bottom: 1px solid var(--border-default, #393B40);
    background-color: var(--bg-canvas, #1E1F22);
  }

  .palette-input {
    flex: 1;
    background: transparent;
    border: none;
    outline: none;
    font-size: 13px;
    color: var(--text-primary, #DFE1E5);
  }

  :global(.palette-results) {
    max-height: 320px;
    overflow-y: auto;
    padding: 4px 0;
  }

  :global(.result-row) {
    padding: 8px 14px;
    cursor: pointer;
    transition: background-color 0.1s ease;
  }

  :global(.result-row:hover) {
    background-color: var(--bg-hover, #2E3136);
  }

  .row-title {
    font-size: 12px;
    font-weight: 500;
    color: var(--text-primary, #DFE1E5);
  }

  .row-desc {
    font-size: 11px;
    color: var(--text-muted, #767980);
  }
</style>
