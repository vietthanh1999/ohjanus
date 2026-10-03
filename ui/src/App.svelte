<script lang="ts">
  import { appState } from './lib/state/appState.svelte';
  import Titlebar from './lib/components/chrome/Titlebar.svelte';
  import StatusBar from './lib/components/chrome/StatusBar.svelte';
  import Sidebar from './lib/components/sidebar/Sidebar.svelte';
  import TabBar from './lib/components/tabs/TabBar.svelte';
  import SqlConsoleView from './lib/components/views/SqlConsoleView.svelte';
  import TableDataView from './lib/components/views/TableDataView.svelte';
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

    <!-- Right Pane: Document Work Area -->
    <main class="work-area">
      <!-- Tabs Bar -->
      <TabBar />

      <!-- Active Content View -->
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
    background-color: var(--bg-window-frame, #24272A);
  }

  :global(.main-layout) {
    flex: 1;
    display: flex;
    overflow: hidden;
    position: relative;
    padding: 0;
    gap: 0;
  }

  :global(.work-area) {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    background-color: var(--bg-canvas);
    min-width: 0;
    border-radius: 0;
    border: none;
    border-left: 1px solid var(--border-subtle, #323438);
    position: relative;
  }

  :global(.view-content) {
    flex: 1;
    overflow: hidden;
    position: relative;
  }
</style>
