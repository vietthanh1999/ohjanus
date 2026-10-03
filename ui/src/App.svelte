<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import { appState } from './lib/state/appState.svelte';
  import Titlebar from './lib/components/chrome/Titlebar.svelte';
  import StatusBar from './lib/components/chrome/StatusBar.svelte';
  import Sidebar from './lib/components/sidebar/Sidebar.svelte';
  import TabBar from './lib/components/tabs/TabBar.svelte';
  import SqlConsoleView from './lib/components/views/SqlConsoleView.svelte';
  import TableDataView from './lib/components/views/TableDataView.svelte';
  import TableDataConsole from './lib/components/views/TableDataConsole.svelte';
  import ApprovalsView from './lib/components/views/ApprovalsView.svelte';
  import AuditView from './lib/components/views/AuditView.svelte';
  import TokensView from './lib/components/views/TokensView.svelte';
  import ConnectionsView from './lib/components/views/ConnectionsView.svelte';
  import DashboardView from './lib/components/views/DashboardView.svelte';

  // Modals & UI Components
  import { Toaster, Box, Flex } from '@ohjanus/ui';
  import ApprovalDecisionModal from './lib/components/modals/ApprovalDecisionModal.svelte';
  import CreateTokenModal from './lib/components/modals/CreateTokenModal.svelte';
  import DdlModal from './lib/components/modals/DdlModal.svelte';
  import SearchPaletteModal from './lib/components/modals/SearchPaletteModal.svelte';
  import SettingsModal from './lib/components/modals/SettingsModal.svelte';

  let lastShiftTime = 0;

  onMount(() => {
    void appState.loadAll();
    appState.startLive();
  });

  onDestroy(() => {
    appState.stopLiveUpdates();
  });
  let isDraggingVerticalSplitter = $state(false);
  let startX = 0;
  let startWidth = 0;

  let isDraggingWorkSplitter = $state(false);
  let startWorkY = 0;
  let startWorkHeight = 0;

  function handleVerticalSplitterMouseDown(e: MouseEvent) {
    isDraggingVerticalSplitter = true;
    startX = e.clientX;
    startWidth = appState.sidebarWidth;

    const onMouseMove = (ev: MouseEvent) => {
      if (!isDraggingVerticalSplitter) return;
      const delta = ev.clientX - startX;
      const newWidth = Math.max(180, Math.min(600, startWidth + delta));
      appState.sidebarWidth = newWidth;
    };

    const onMouseUp = () => {
      isDraggingVerticalSplitter = false;
      window.removeEventListener('mousemove', onMouseMove);
      window.removeEventListener('mouseup', onMouseUp);
    };

    window.addEventListener('mousemove', onMouseMove);
    window.addEventListener('mouseup', onMouseUp);
  }

  function handleWorkSplitterMouseDown(e: MouseEvent) {
    isDraggingWorkSplitter = true;
    startWorkY = e.clientY;
    startWorkHeight = appState.servicesHeight;

    const onMouseMove = (ev: MouseEvent) => {
      if (!isDraggingWorkSplitter) return;
      const delta = startWorkY - ev.clientY;
      const newHeight = Math.max(120, Math.min(450, startWorkHeight + delta));
      appState.servicesHeight = newHeight;
    };

    const onMouseUp = () => {
      isDraggingWorkSplitter = false;
      window.removeEventListener('mousemove', onMouseMove);
      window.removeEventListener('mouseup', onMouseUp);
    };

    window.addEventListener('mousemove', onMouseMove);
    window.addEventListener('mouseup', onMouseUp);
  }

  function handleGlobalKeydown(e: KeyboardEvent) {
    // Cmd+K or Ctrl+K -> Search Everywhere
    if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
      e.preventDefault();
      appState.searchModalOpen = true;
      return;
    }

    // Double Shift -> Search Everywhere
    if (e.key === 'Shift') {
      const now = Date.now();
      if (now - lastShiftTime < 300) {
        appState.searchModalOpen = true;
      }
      lastShiftTime = now;
      return;
    }

    // Cmd+Enter or Ctrl+Enter -> Execute Query
    if ((e.metaKey || e.ctrlKey) && e.key === 'Enter') {
      e.preventDefault();
      appState.executeQuery();
      return;
    }

    // Cmd+1 -> Toggle Sidebar
    if ((e.metaKey || e.ctrlKey) && e.key === '1') {
      e.preventDefault();
      appState.isSidebarCollapsed = !appState.isSidebarCollapsed;
      return;
    }

    // Cmd+, -> Settings
    if ((e.metaKey || e.ctrlKey) && e.key === ',') {
      e.preventDefault();
      appState.settingsModalOpen = true;
      return;
    }
  }
</script>

<svelte:window onkeydown={handleGlobalKeydown} />

<Box class="app-root">
  <!-- Top: Window Titlebar -->
  <Titlebar />

  <!-- Center: Main Layout (Sidebar + Work Area) -->
  <Flex class="main-layout">
    <!-- Left Pane: Sidebar (Database Explorer + Services) -->
    <Sidebar />

    {#if !appState.isSidebarCollapsed}
      <!-- Vertical Splitter / 6px gap between Sidebar & Work Area -->
      <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
      <div
        class="vertical-splitter"
        role="separator"
        tabindex="-1"
        aria-orientation="vertical"
        onmousedown={handleVerticalSplitterMouseDown}
        title="Drag to resize Database Explorer"
      ></div>
    {/if}

    <!-- Right Pane: Document Work Area -->
    <main class="work-area">
      <!-- Upper Card Island: TabBar + Active Content View -->
      <div class="editor-card-island">
        <TabBar />

        <Box class="view-content">
          {#if appState.activeTabId === 'console_2'}
            <SqlConsoleView />
          {:else if appState.activeTabId === 'connection_credential'}
            <TableDataView />
          {:else if appState.activeTabId === 'approvals'}
            <ApprovalsView />
          {:else if appState.activeTabId === 'audit'}
            <AuditView />
          {:else if appState.activeTabId === 'tokens'}
            <TokensView />
          {:else if appState.activeTabId === 'connections'}
            <ConnectionsView />
          {:else if appState.activeTabId === 'dashboard'}
            <DashboardView />
          {/if}
        </Box>
      </div>

      {#if ['connection_credential', 'commands', 'events'].includes(appState.activeTabId)}
        <!-- Horizontal Splitter / 6px gap between Editor and Console -->
        <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
        <div
          class="horizontal-work-splitter"
          role="separator"
          tabindex="-1"
          aria-orientation="horizontal"
          onmousedown={handleWorkSplitterMouseDown}
          title="Drag to resize Editor / Console"
        ></div>

        <!-- Lower Card Island: Table Log Console & Action Strip -->
        <div class="console-card-island" style="height: {appState.servicesHeight}px;">
          <TableDataConsole />
        </div>
      {/if}
    </main>
  </Flex>

  <!-- Bottom: Global Status Bar -->
  <StatusBar />

  <!-- Global Modals & Notifications -->
  <ApprovalDecisionModal />
  <CreateTokenModal />
  <DdlModal />
  <SearchPaletteModal />
  <SettingsModal />
  <Toaster position="bottom-right" />
</Box>

<style>
  :global(.app-root) {
    width: 100vw;
    height: 100vh;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    background: linear-gradient(
      90deg,
      #2a3032 0%,
      #292f32 30%,
      #26282b 70%,
      #26282b 100%
    );
  }

  :global(.main-layout) {
    flex: 1;
    display: flex;
    overflow: hidden;
    position: relative;
    padding: 0 6px 6px 6px;
    gap: 0;
    background: transparent;
  }

  .vertical-splitter {
    width: 6px;
    height: 100%;
    cursor: col-resize;
    flex-shrink: 0;
    background: transparent;
    z-index: 10;
  }

  .vertical-splitter:hover {
    background-color: var(--border-accent, #3574F0);
  }

  :global(.work-area) {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    background-color: transparent;
    min-width: 0;
    border-radius: 0;
    border: none;
    position: relative;
    gap: 0;
  }

  .editor-card-island {
    flex: 1;
    min-height: 0;
    display: flex;
    flex-direction: column;
    background-color: var(--bg-card, #191A1C);
    border: none;
    border-radius: 8px;
    overflow: hidden;
  }

  .horizontal-work-splitter {
    height: 6px;
    width: 100%;
    cursor: row-resize;
    flex-shrink: 0;
    background: transparent;
    z-index: 10;
  }

  .horizontal-work-splitter:hover {
    background-color: var(--border-accent, #3574F0);
  }

  .console-card-island {
    display: flex;
    background-color: var(--bg-card, #191A1C);
    border: none;
    border-radius: 8px;
    overflow: hidden;
    flex-shrink: 0;
  }

  :global(.view-content) {
    flex: 1;
    overflow: hidden;
    position: relative;
  }
</style>
