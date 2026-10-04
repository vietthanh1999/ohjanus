<script lang="ts">
  import { appState } from '../../state/appState.svelte';
  import { DialogPrimitive, Badge, Kbd, Box, Flex, Text, Input } from '@ohjanus/ui';
  import { Icon } from '@ohjanus/icons';

  let query = $state('');

  interface PaletteItem {
    title: string;
    category: string;
    desc: string;
    action: () => void;
  }

  let items = $derived.by((): PaletteItem[] => {
    const list: PaletteItem[] = [
      { title: 'Approvals Queue', category: 'Gateway', desc: 'Review write requests', action: () => appState.openTab({ id: 'approvals', title: 'Approvals Queue', type: 'approvals', closable: false, icon: 'shield' }) },
      { title: 'Audit Log Trail', category: 'Gateway', desc: 'Audit records and query logs', action: () => appState.openTab({ id: 'audit', title: 'Audit Logs', type: 'audit', closable: false, icon: 'audit' }) },
      { title: 'MCP Access Tokens', category: 'Gateway', desc: 'Manage agent tokens & scopes', action: () => appState.openTab({ id: 'tokens', title: 'MCP Tokens', type: 'tokens', closable: false, icon: 'key' }) },
      { title: 'Connection Pools', category: 'Gateway', desc: 'Monitor database connections', action: () => appState.openTab({ id: 'connections', title: 'Connections', type: 'connections', closable: false, icon: 'database' }) },
      { title: 'Gateway Dashboard', category: 'Gateway', desc: 'Throughput and telemetry metrics', action: () => appState.openTab({ id: 'dashboard', title: 'Dashboard', type: 'dashboard', closable: false, icon: 'chart' }) },
      { title: 'Generate DDL', category: 'Action', desc: 'Inspect current schema definition', action: () => appState.ddlModalOpen = true }
    ];
    for (const conn of appState.connections) {
      list.push({
        title: `console [${conn.name}]`,
        category: 'Console',
        desc: `${conn.name} query editor`,
        action: () => appState.openConsoleTab(conn.name)
      });
      for (const schema of appState.schemasOf(conn.name)) {
        for (const table of schema.tables) {
          list.push({
            title: table.name,
            category: 'Table',
            desc: `${conn.name}.${schema.name}`,
            action: () => void appState.selectTable(conn.name, schema.name, table.name)
          });
        }
      }
    }
    return list;
  });

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
      <Input
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
            <Text size="xs" weight="medium" class="row-title">{item.title}</Text>
          </Flex>
          <Text class="row-desc" size="sm" color="muted">{item.desc}</Text>
        </Flex>
      {/each}
      {#if filtered.length === 0}
        <Text size="xs" color="muted" class="palette-empty">
          No matches. Tables appear here after the explorer loads their schema.
        </Text>
      {/if}
    </Box>
  </Box>
</DialogPrimitive>

<style>
  :global(.palette-modal) {
    width: 560px;
    background-color: #191A1C;
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
    background-color: #191A1C;
  }

  :global(.palette-input-wrap .ohjanus-input-wrapper) {
    flex: 1;
    background: transparent !important;
    border: none !important;
    box-shadow: none !important;
    padding: 0 !important;
    height: auto !important;
  }

  :global(.palette-input-wrap .ohjanus-input-field) {
    background: transparent !important;
    border: none !important;
    outline: none !important;
    box-shadow: none !important;
    font-size: 13px !important;
    color: var(--text-primary, #DFE1E5) !important;
    padding: 0 !important;
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

  :global(.row-title) {
    color: var(--text-primary, #DFE1E5);
  }

  :global(.palette-empty) {
    display: block;
    padding: 16px;
    text-align: center;
  }
</style>
