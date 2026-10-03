<script lang="ts">
  import { appState } from '../../state/appState.svelte';
  import { Icon } from '@ohjanus/icons';
  import { Text } from '@ohjanus/ui';

  let isDraggingSplitter = $state(false);
  let startY = 0;
  let startHeight = 0;

  function handleSplitterMouseDown(e: MouseEvent) {
    isDraggingSplitter = true;
    startY = e.clientY;
    startHeight = appState.servicesHeight;

    const onMouseMove = (ev: MouseEvent) => {
      if (!isDraggingSplitter) return;
      const delta = startY - ev.clientY;
      const newHeight = Math.max(120, Math.min(450, startHeight + delta));
      appState.servicesHeight = newHeight;
    };

    const onMouseUp = () => {
      isDraggingSplitter = false;
      window.removeEventListener('mousemove', onMouseMove);
      window.removeEventListener('mouseup', onMouseUp);
    };

    window.addEventListener('mousemove', onMouseMove);
    window.addEventListener('mouseup', onMouseUp);
  }

  function connKey(name: string) {
    return `conn:${name}`;
  }

  function schemaKey(conn: string, schema: string) {
    return `conn:${conn}:schema:${schema}`;
  }

  async function toggleConnection(name: string) {
    appState.toggleTree(connKey(name));
    if (appState.treeExpanded[connKey(name)]) {
      await appState.loadSchema(name);
    }
  }

  function openConsole(name: string) {
    appState.selectedTreeNode = `console:${name}`;
    appState.openConsoleTab(name);
  }

  async function openTable(conn: string, schema: string, table: string) {
    await appState.selectTable(conn, schema, table);
  }

  function matchesFilter(text: string): boolean {
    const q = appState.treeFilterQuery.trim().toLowerCase();
    return q === '' || text.toLowerCase().includes(q);
  }
</script>

<aside
  class="sidebar"
  class:collapsed={appState.isSidebarCollapsed}
  style="width: {appState.isSidebarCollapsed ? '0px' : appState.sidebarWidth + 'px'};"
>
  <div class="explorer-pane">
    <div class="explorer-header-top">
      <Text size="md" weight="semibold">Database Explorer</Text>
      <div class="header-window-icons">
        <button type="button" class="jb-icon-btn" title="Reload connections" onclick={() => void appState.loadConnections()}>
          <Icon name="refresh" size={12} />
        </button>
        <button type="button" class="jb-icon-btn" title="New console on first connection" onclick={() => appState.connections[0] && openConsole(appState.connections[0].name)}>
          <Icon name="plus" size={12} />
        </button>
      </div>
    </div>

    <div class="explorer-subtoolbar">
      <input
        type="text"
        class="tree-filter-input"
        placeholder="Filter connections / tables…"
        bind:value={appState.treeFilterQuery}
      />
      <button type="button" class="jb-icon-btn ddl-btn" title="Generate Table DDL" onclick={() => appState.ddlModalOpen = true}>
        <Text size="xs" weight="bold" color="muted" mono>DDL</Text>
      </button>
    </div>

    <div class="tree-viewport">
      {#if appState.dataLoading}
        <div class="tree-empty">Loading connections from Admin API…</div>
      {:else if appState.dataError && appState.connections.length === 0}
        <div class="tree-empty">
          <div>Admin API unreachable.</div>
          <div class="tree-error">{appState.dataError}</div>
          <button type="button" class="retry-btn" onclick={() => void appState.loadAll()}>Retry</button>
        </div>
      {:else if appState.connections.length === 0}
        <div class="tree-empty">No connections configured. Add one in janus.yaml and restart serve.</div>
      {:else}
        {#each appState.explorerConnections as conn (conn.name)}
          <button type="button" class="tree-node depth-0" onclick={() => void toggleConnection(conn.name)}>
            <span class="chevron" class:expanded={appState.treeExpanded[connKey(conn.name)]}><Icon name="chevron-right" size={12} /></span>
            <Icon name="database" size={14} color="#3B82F6" class="node-icon" />
            <Text size="lg" truncate style="flex: 1;"><strong>{conn.name}</strong>{conn.readonly ? ' [ReadOnly]' : ''}</Text>
            {#if conn.status !== 'healthy'}
              <span class="conn-warn" title={conn.status}>●</span>
            {/if}
          </button>

          {#if appState.treeExpanded[connKey(conn.name)]}
            <button type="button" class="tree-node depth-1" onclick={() => openConsole(conn.name)}>
              <span class="chevron"><Icon name="chevron-right" size={12} /></span>
              <Icon name="lightning" size={12} color="#3B82F6" class="node-icon" />
              <Text size="lg" truncate style="flex: 1;">console [{conn.name}]</Text>
            </button>

            {#if appState.schemaLoading[conn.name]}
              <div class="tree-empty">Loading schema…</div>
            {:else if appState.schemaError[conn.name]}
              <div class="tree-empty">
                <div class="tree-error">{appState.schemaError[conn.name]}</div>
                <button type="button" class="retry-btn" onclick={() => void appState.loadSchema(conn.name, true)}>Retry</button>
              </div>
            {:else}
              {#each appState.schemasOf(conn.name) as schema (schema.name)}
                {#if matchesFilter(schema.name) || schema.tables.some((t) => matchesFilter(t.name))}
                  <button type="button" class="tree-node depth-1" onclick={() => appState.toggleTree(schemaKey(conn.name, schema.name))}>
                    <span class="chevron" class:expanded={appState.treeExpanded[schemaKey(conn.name, schema.name)]}><Icon name="chevron-right" size={12} /></span>
                    <Icon name="folder" size={13} color="#C29D38" class="node-icon" />
                    <Text size="lg" truncate style="flex: 1;">{schema.name} <Text size="sm" color="muted">{schema.tables.length}</Text></Text>
                  </button>

                  {#if appState.treeExpanded[schemaKey(conn.name, schema.name)]}
                    {#each schema.tables.filter((t) => matchesFilter(t.name)) as table (table.name)}
                      <button
                        type="button"
                        class="tree-node depth-2"
                        class:selected={appState.selectedTreeNode === `table:${conn.name}.${schema.name}.${table.name}`}
                        onclick={() => void openTable(conn.name, schema.name, table.name)}
                      >
                        <span class="chevron"><Icon name="chevron-right" size={12} /></span>
                        <Icon name="table" size={12} color="#4A88C7" class="node-icon" />
                        <Text size="lg" truncate style="flex: 1;">{table.name}</Text>
                      </button>
                    {/each}
                  {/if}
                {/if}
              {/each}
            {/if}
          {/if}
        {/each}
      {/if}
    </div>
  </div>

  <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
  <div
    class="horizontal-splitter"
    role="separator"
    tabindex="-1"
    aria-orientation="horizontal"
    onmousedown={handleSplitterMouseDown}
    title="Drag to resize Database Explorer / Services"
  ></div>

  <div class="services-pane" style="height: {appState.servicesHeight}px;">
    <div class="services-header">
      <span class="services-title">Services</span>
      <Text size="xs" color="muted" style="background-color: var(--bg-hover); border-radius: 8px; padding: 0 6px;">{appState.services.length}</Text>
    </div>

    <div class="services-viewport">
      {#if appState.services.length === 0}
        <div class="tree-empty">No open sessions. Open a table or console from the explorer.</div>
      {:else}
        {#each appState.services as session (session.id)}
          <button
            type="button"
            class="service-item depth-2"
            class:selected={appState.activeTabId === session.id}
            onclick={() => appState.activeTabId = session.id}
            title={session.connection ? `${session.name} · ${session.connection}` : session.name}
          >
            {#if session.type === 'table'}
              <Icon name="table" size={12} color="#4A88C7" class="node-icon" />
            {:else}
              <Icon name="lightning" size={12} color="#3B82F6" class="node-icon" />
            {/if}
            <Text size="lg" truncate style="flex: 1;">{session.name}</Text>
            {#if session.durationMs !== undefined}
              <Text size="sm" color="muted" mono style="margin-left: 4px;">{session.durationMs} ms</Text>
            {/if}
          </button>
        {/each}
      {/if}
    </div>
  </div>
</aside>

<style>
  .sidebar {
    background-color: transparent;
    border: none;
    border-radius: 0;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    flex-shrink: 0;
    transition: width 0.15s ease-out;
    height: 100%;
    gap: 0;
  }

  .sidebar.collapsed {
    border: none;
    width: 0 !important;
  }

  .explorer-header-top {
    height: var(--toolbar-height, 32px);
    padding: 0 8px 0 12px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-shrink: 0;
  }


  .header-window-icons {
    display: flex;
    align-items: center;
    gap: 2px;
  }

  .explorer-subtoolbar {
    height: var(--toolbar-height, 32px);
    padding: 0 8px;
    display: flex;
    align-items: center;
    gap: 4px;
    border-bottom: 1px solid var(--border-subtle);
    flex-shrink: 0;
  }

  .tree-filter-input {
    flex: 1;
    height: 24px;
    background-color: var(--bg-canvas);
    border: 1px solid var(--border-default);
    border-radius: 4px;
    color: var(--text-primary);
    font-size: 12px;
    padding: 0 8px;
    outline: none;
  }

  .explorer-pane {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 120px;
    overflow: hidden;
    background-color: var(--bg-card, #191A1C);
    border: none;
    border-radius: 8px;
  }

  .tree-viewport, .services-viewport {
    flex: 1;
    overflow-y: auto;
    overflow-x: hidden;
    padding: 4px 0;
  }

  .tree-node, .service-item {
    height: var(--tree-row-height, 28px);
    display: flex;
    align-items: center;
    gap: 4px;
    padding-right: 8px;
    cursor: pointer;
    font-size: var(--font-size-base, 14px);
    color: var(--text-primary);
    border-radius: var(--radius-sm, 4px);
    margin: 1px 4px;
    transition: background-color 0.1s ease;
    text-align: left;
    width: calc(100% - 8px);
  }

  .tree-node:hover, .service-item:hover {
    background-color: var(--bg-hover);
  }

  .tree-node.selected, .service-item.selected {
    background-color: var(--bg-selected);
    color: #FFFFFF;
  }

  .depth-0 { padding-left: 6px; }
  .depth-1 { padding-left: 20px; }
  .depth-2 { padding-left: 34px; }

  .chevron {
    width: 16px;
    font-size: 13px;
    color: var(--text-muted);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    transition: transform 0.1s ease;
    transform: rotate(0deg);
    flex-shrink: 0;
  }

  .chevron.expanded {
    transform: rotate(90deg);
  }

  :global(.node-icon) {
    margin-right: 2px;
    width: 15px;
    height: 15px;
    flex-shrink: 0;
  }

  .conn-warn {
    color: var(--action-danger, #E55353);
    font-size: 10px;
  }

  .tree-empty {
    padding: 12px;
    font-size: 12px;
    color: var(--text-muted);
    line-height: 1.5;
  }

  .tree-error {
    color: var(--action-danger, #E55353);
    font-size: 11px;
    margin: 4px 0;
    word-break: break-word;
  }

  .retry-btn {
    margin-top: 6px;
    font-size: 12px;
    color: var(--action-primary, #3574F0);
    text-decoration: underline;
  }

  .horizontal-splitter {
    height: 6px;
    background-color: transparent;
    cursor: row-resize;
    flex-shrink: 0;
    z-index: 10;
  }

  .horizontal-splitter:hover {
    background-color: var(--border-accent, #3574F0);
  }

  .services-pane {
    display: flex;
    flex-direction: column;
    background-color: var(--bg-card, #191A1C);
    border: none;
    border-radius: 8px;
    overflow: hidden;
    flex-shrink: 0;
  }

  .services-header {
    height: var(--toolbar-height, 32px);
    padding: 0 8px 0 12px;
    display: flex;
    align-items: center;
    gap: 8px;
    flex-shrink: 0;
  }

  .services-title {
    font-size: 12px;
    font-weight: 600;
    color: var(--text-primary);
  }

</style>
