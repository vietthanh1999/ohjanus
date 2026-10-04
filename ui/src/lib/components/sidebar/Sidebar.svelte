<script lang="ts">
  import { appState } from '../../state/appState.svelte';
  import { Box } from '@ohjanus/ui';
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

  function handleOpenFirstConsole() {
    if (appState.connections[0]) {
      appState.selectedTreeNode = `console:${appState.connections[0].name}`;
      appState.openConsoleTab(appState.connections[0].name);
    }
  }
</script>

<Box
  class="sidebar {appState.isSidebarCollapsed ? 'collapsed' : ''}"
  style="width: {appState.isSidebarCollapsed ? '0px' : appState.sidebarWidth + 'px'};"
>
  <Box class="explorer-pane">
    <ExplorerHeader
      onreload={() => void appState.loadConnections()}
      onnewconsole={handleOpenFirstConsole}
    />

    <ExplorerFilterBar
      bind:query={appState.treeFilterQuery}
      onddl={() => (appState.ddlModalOpen = true)}
    />

    <ExplorerTree />
  </Box>

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

  <ServicesPane height={appState.servicesHeight} />
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
