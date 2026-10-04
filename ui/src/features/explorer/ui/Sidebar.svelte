<script lang="ts">
  import { parseTableTabId } from '@/entities/tab';
  import { connectionsState } from '@/features/connections';
  import { consoleState } from '@/features/console';
  import { explorerState } from '@/features/explorer';
  import { tableViewerState } from '@/features/table-viewer';
  import { workbenchState } from '@/features/workbench';
  import { Box, toast } from '@ohjanus/ui';
  import ExplorerHeader from './ExplorerHeader.svelte';
  import ExplorerFilterBar from './ExplorerFilterBar.svelte';
  import ExplorerTree from './ExplorerTree.svelte';
  import ServicesPane from './ServicesPane.svelte';

  let isDraggingSplitter = $state(false);
  let startY = 0;
  let startHeight = 0;

  function handleSplitterMouseDown(e: MouseEvent) {
    isDraggingSplitter = true;
    startY = e.clientY;
    startHeight = explorerState.servicesHeight;

    const onMouseMove = (ev: MouseEvent) => {
      if (!isDraggingSplitter) return;
      const delta = startY - ev.clientY;
      const newHeight = Math.max(120, Math.min(450, startHeight + delta));
      explorerState.servicesHeight = newHeight;
    };

    const onMouseUp = () => {
      isDraggingSplitter = false;
      window.removeEventListener('mousemove', onMouseMove);
      window.removeEventListener('mouseup', onMouseUp);
    };

    window.addEventListener('mousemove', onMouseMove);
    window.addEventListener('mouseup', onMouseUp);
  }

  function handleOpenFirstConsole() {
    if (connectionsState.connections[0]) {
      const name = connectionsState.connections[0].name;
      explorerState.revealConsole(name);
      consoleState.openConsole(name);
    }
  }

  function handleManageConnections() {
    workbenchState.openTab({
      id: 'connections',
      title: 'Connection Pools',
      type: 'connections',
      closable: false,
      icon: 'database'
    });
  }

  /** Refresh connections + force schema introspection for expanded nodes. */
  async function handleRefresh() {
    await connectionsState.loadConnections().catch(() => {});
    const expanded = Object.keys(explorerState.treeExpanded).filter((k) => explorerState.treeExpanded[k]);
    await Promise.all(
      connectionsState.connections
        .filter((c) => expanded.includes(`conn:${c.name}`))
        .map((c) => explorerState.loadSchema(c.name, true).catch(() => {}))
    );
  }

  /** Resolve the current table target (same rule as DdlModal). */
  function currentTarget(): { connection: string; schema: string; table: string } | null {
    const tab = workbenchState.activeTab;
    if (tab?.type === 'table' && tab.connection && tab.schema && tab.table) {
      return { connection: tab.connection, schema: tab.schema, table: tab.table };
    }
    const parsed = parseTableTabId(workbenchState.activeTabId);
    if (parsed) return parsed;
    const t = tableViewerState.tableViewer;
    if (t.connection && t.table) return { connection: t.connection, schema: t.schema, table: t.table };
    return null;
  }

  function handleCopyDdl() {
    const target = currentTarget();
    if (!target) {
      toast.info('Copy DDL', 'Open a table from the explorer first.');
      return;
    }
    const ddl = explorerState.ddlFor(target.connection, target.schema, target.table);
    if (!ddl) {
      toast.info('Copy DDL', 'Schema is not loaded yet. Expand the connection and retry.');
      return;
    }
    navigator.clipboard.writeText(ddl);
    toast.success('DDL copied to clipboard');
  }
</script>

<Box
  class="sidebar {explorerState.isSidebarCollapsed ? 'collapsed' : ''}"
  style="width: {explorerState.isSidebarCollapsed ? '0px' : explorerState.sidebarWidth + 'px'};"
>
  <Box class="explorer-pane">
    <ExplorerHeader
      onreload={() => void handleRefresh()}
      onnewconsole={handleOpenFirstConsole}
      onmanageconnections={handleManageConnections}
      oncopyddl={handleCopyDdl}
    />

    <ExplorerFilterBar bind:query={explorerState.treeFilterQuery} />

    <ExplorerTree />
  </Box>

  {#if explorerState.servicesVisible}
    <!-- Horizontal Splitter between Explorer & Services -->
    <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
    <Box
      class="horizontal-splitter"
      role="separator"
      tabindex={-1}
      aria-orientation="horizontal"
      onmousedown={handleSplitterMouseDown}
      title="Drag to resize Database Explorer / Services"
    />

    <ServicesPane height={explorerState.servicesHeight} />
  {/if}
</Box>

<style>
  :global(.sidebar) {
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

  :global(.sidebar.collapsed) {
    border: none;
    width: 0 !important;
  }

  :global(.explorer-pane) {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 120px;
    overflow: hidden;
    background-color: var(--bg-card, #191A1C);
    border-radius: var(--radius-md, 8px);
  }

  :global(.horizontal-splitter) {
    height: 6px;
    background-color: transparent;
    cursor: row-resize;
    flex-shrink: 0;
    z-index: 10;
  }

  :global(.horizontal-splitter:hover) {
    background-color: var(--border-accent, #3574F0);
  }
</style>
