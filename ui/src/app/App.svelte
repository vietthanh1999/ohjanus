<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import { bootState } from './boot.svelte';
  import { componentFor } from './tabs';
  import { handleGlobalKeydown } from './shortcuts';
  import Titlebar from './chrome/Titlebar.svelte';
  import StatusBar from './chrome/StatusBar.svelte';
  import { Sidebar, DdlModal } from '@/features/explorer';
  import { TabBar, workbenchState } from '@/features/workbench';
  import { TableDataConsole } from '@/features/table-viewer';
  import { ApprovalDecisionModal, approvalsState } from '@/features/approvals';
  import { CreateTokenModal } from '@/features/tokens';
  import { SearchPaletteModal } from '@/features/command-palette';
  import { SettingsModal } from '@/features/settings';
  import { explorerState } from '@/features/explorer';

  // Modals & UI Components
  import { Toaster, Box, Flex } from '@ohjanus/ui';

  const ActiveView = $derived(componentFor(workbenchState.activeTab));

  onMount(() => {
    void bootState.loadAll();
    approvalsState.startLive();
  });

  onDestroy(() => {
    approvalsState.stopLiveUpdates();
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
    startWidth = explorerState.sidebarWidth;

    const onMouseMove = (ev: MouseEvent) => {
      if (!isDraggingVerticalSplitter) return;
      const delta = ev.clientX - startX;
      const newWidth = Math.max(180, Math.min(600, startWidth + delta));
      explorerState.sidebarWidth = newWidth;
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
    startWorkHeight = explorerState.servicesHeight;

    const onMouseMove = (ev: MouseEvent) => {
      if (!isDraggingWorkSplitter) return;
      const delta = startWorkY - ev.clientY;
      const newHeight = Math.max(120, Math.min(450, startWorkHeight + delta));
      explorerState.servicesHeight = newHeight;
    };

    const onMouseUp = () => {
      isDraggingWorkSplitter = false;
      window.removeEventListener('mousemove', onMouseMove);
      window.removeEventListener('mouseup', onMouseUp);
    };

    window.addEventListener('mousemove', onMouseMove);
    window.addEventListener('mouseup', onMouseUp);
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

    {#if !explorerState.isSidebarCollapsed}
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
          {#if ActiveView}
            <ActiveView />
          {/if}
        </Box>
      </div>

      {#if workbenchState.activeTab?.type === 'table'}
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
        <div class="console-card-island" style="height: {explorerState.servicesHeight}px;">
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
  <Toaster />
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
