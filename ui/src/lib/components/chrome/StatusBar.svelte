<script lang="ts">
  import { appState } from '../../state/appState.svelte';
  import { Flex, Badge } from '@ohjanus/ui';
  import { Icon } from '@ohjanus/icons';

  function openApprovals() {
    appState.activeTabId = 'approvals';
  }
</script>

<footer class="statusbar">
  <!-- Left: Interactive Navigation Breadcrumb -->
  <Flex class="breadcrumb-strip" align="center" gap="4px">
    {#if appState.activeTabId === 'connection_credential'}
      <span class="crumb-item">Database</span>
      <span class="crumb-sep">&gt;</span>
      <span class="crumb-item"><strong>[Dev][ReadOnly]</strong> 10.220.6.4</span>
      <span class="crumb-sep">&gt;</span>
      <span class="crumb-item">dev_mh_asset</span>
      <span class="crumb-sep">&gt;</span>
      <span class="crumb-item">public</span>
      <span class="crumb-sep">&gt;</span>
      <span class="crumb-item">tables</span>
      <span class="crumb-sep">&gt;</span>
      <span class="crumb-item active-crumb">
        <Icon name="table" size={11} color="#4A88C7" style="margin-right: 4px;" />
        connection_credential
      </span>
    {:else if appState.activeTabId === 'console_2'}
      <span class="crumb-item">Database Consoles</span>
      <span class="crumb-sep">&gt;</span>
      <span class="crumb-item"><strong>[PRD]</strong> 10.250.6.23</span>
      <span class="crumb-sep">&gt;</span>
      <span class="crumb-item active-crumb">
        <Icon name="lightning" size={11} color="#57D38C" style="margin-right: 4px;" />
        console_2 [[PRD] 10.250.6.23]
      </span>
    {:else}
      <span class="crumb-item">MCP Gateway Security</span>
      <span class="crumb-sep">&gt;</span>
      <span class="crumb-item active-crumb">{appState.activeTabId.toUpperCase()}</span>
    {/if}
  </Flex>

  <!-- Right: Document, Buffer Stats & System Status -->
  <Flex class="status-right" align="center" gap="12px">
    {#if appState.activeTabId === 'console_2'}
      <!-- Cursor Position & Buffer Analytics -->
      <span class="status-item code-text" title="Cursor line:col (character count, line breaks)">
        {appState.cursorPos.line}:{appState.cursorPos.col} (2954 chars, 73 line breaks)
      </span>

      <span class="status-item" title="Line Endings">LF</span>
      <span class="status-item" title="File Encoding">UTF-8</span>
      <span class="status-item" title="Indent Size">4 spaces</span>

      <!-- Read-Only Lock Icon -->
      <span class="status-item lock-icon" title="Database Connection is in ReadOnly Transaction Mode">
        <Icon name="lock" size={11} />
      </span>
    {:else}
      <!-- Copy / Terminal action icon matching design2.png -->
      <button type="button" class="status-item" title="Open in Terminal">
        <Icon name="terminal" size={12} />
      </button>
    {/if}

    <!-- Notification Bell -->
    <button
      class="status-item bell-btn"
      class:has-notification={appState.notificationCount > 0}
      title="Pending Approvals ({appState.notificationCount})"
      onclick={openApprovals}
    >
      <Icon name="bell" size={12} />
      {#if appState.notificationCount > 0}
        <span class="notification-badge">{appState.notificationCount}</span>
      {/if}
    </button>

    <!-- License Badge via @ohjanus/ui Badge -->
    <Badge variant="license" size="sm" class="license-pill" title="Community / Evaluation License">
      Non-commercial use
    </Badge>
  </Flex>
</footer>

<style>
  .statusbar {
    height: var(--statusbar-height);
    background-color: var(--bg-window-frame, #24272A);
    border-top: none;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 10px;
    font-size: 11px;
    color: var(--text-secondary);
    user-select: none;
    z-index: 50;
    flex-shrink: 0;
  }

  .crumb-item {
    display: flex;
    align-items: center;
    padding: 1px 4px;
    border-radius: 2px;
    cursor: pointer;
    color: var(--text-secondary);
  }

  .crumb-item:hover {
    background-color: var(--bg-hover);
    color: var(--text-primary);
  }

  .crumb-item.active-crumb {
    color: var(--text-primary);
  }

  .crumb-sep {
    color: var(--text-muted);
    font-size: 10px;
  }

  .status-item {
    padding: 1px 4px;
    border-radius: 2px;
    cursor: pointer;
    color: var(--text-muted);
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .status-item:hover {
    background-color: var(--bg-hover);
    color: var(--text-primary);
  }

  .lock-icon {
    color: var(--text-muted);
  }

  .bell-btn {
    position: relative;
    color: var(--text-muted);
  }

  .bell-btn.has-notification {
    color: #EDA200;
  }

  .notification-badge {
    position: absolute;
    top: -3px;
    right: -4px;
    background-color: #EDA200;
    color: #1E1F22;
    font-size: 8px;
    font-weight: bold;
    padding: 0 3px;
    border-radius: 5px;
    line-height: 11px;
  }

  :global(.license-pill) {
    font-size: 10px !important;
  }
</style>
