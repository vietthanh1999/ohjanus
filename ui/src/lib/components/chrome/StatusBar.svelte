<script lang="ts">
  import { appState } from '../../state/appState.svelte';

  function openApprovals() {
    appState.activeTabId = 'approvals';
  }
</script>

<footer class="statusbar">
  <!-- Left: Interactive Navigation Breadcrumb -->
  <div class="breadcrumb-strip">
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
        <svg width="11" height="11" viewBox="0 0 16 16" fill="#4A88C7" style="margin-right: 4px;">
          <path d="M0 2a1 1 0 011-1h14a1 1 0 011 1v12a1 1 0 01-1 1H1a1 1 0 01-1-1V2zm1 3v2h6V5H1zm7 0v2h7V5H8zm0 3v2h7V8H8zm-1 0H1v2h6V8zm0 3H1v2h6v-2zm1 0v2h7v-2H8z"/>
        </svg>
        connection_credential
      </span>
    {:else if appState.activeTabId === 'console_2'}
      <span class="crumb-item">Database Consoles</span>
      <span class="crumb-sep">&gt;</span>
      <span class="crumb-item"><strong>[PRD]</strong> 10.250.6.23</span>
      <span class="crumb-sep">&gt;</span>
      <span class="crumb-item active-crumb">
        <svg width="11" height="11" viewBox="0 0 16 16" fill="#57D38C" style="margin-right: 4px;">
          <path d="M11.251.068a.5.5 0 01.42.58L10.07 5H14a.5.5 0 01.372.832l-8.5 9.5a.5.5 0 01-.842-.512L6.63 9.5H2.5a.5.5 0 01-.42-.772l8.5-8.5a.5.5 0 01.671-.16z"/>
        </svg>
        console_2 [[PRD] 10.250.6.23]
      </span>
    {:else}
      <span class="crumb-item">MCP Gateway Security</span>
      <span class="crumb-sep">&gt;</span>
      <span class="crumb-item active-crumb">{appState.activeTabId.toUpperCase()}</span>
    {/if}
  </div>

  <!-- Right: Document, Buffer Stats & System Status -->
  <div class="status-right">
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
        <svg width="11" height="11" viewBox="0 0 16 16" fill="currentColor">
          <path fill-rule="evenodd" d="M4 4v2h-.5A1.5 1.5 0 002 7.5v6A1.5 1.5 0 003.5 15h9a1.5 1.5 0 001.5-1.5v-6A1.5 1.5 0 0012.5 6H12V4a4 4 0 00-8 0zm6.5 2V4a2.5 2.5 0 00-5 0v2h5z"/>
        </svg>
      </span>
    {:else}
      <!-- Copy / Window action icon matching design2.png -->
      <button type="button" class="status-item" title="Open in Terminal">
        <svg width="12" height="12" viewBox="0 0 16 16" fill="currentColor">
          <path d="M0 2a2 2 0 012-2h12a2 2 0 012 2v12a2 2 0 01-2 2H2a2 2 0 01-2-2V2zm2-1a1 1 0 00-1 1v12a1 1 0 001 1h12a1 1 0 001-1V2a1 1 0 00-1-1H2z"/>
          <path d="M5.5 10a.5.5 0 01-.5-.5v-3a.5.5 0 011 0v3a.5.5 0 01-.5.5zm5 0a.5.5 0 01-.5-.5v-3a.5.5 0 011 0v3a.5.5 0 01-.5.5z"/>
        </svg>
      </button>
    {/if}

    <!-- Notification Bell -->
    <button
      class="status-item bell-btn"
      class:has-notification={appState.notificationCount > 0}
      title="Pending Approvals ({appState.notificationCount})"
      onclick={openApprovals}
    >
      <svg width="12" height="12" viewBox="0 0 16 16" fill="currentColor">
        <path d="M8 1.5c-2.363 0-4 1.69-4 3.75v1.759c0 .53-.21 1.04-.586 1.416L2.354 9.485A1.5 1.5 0 003.414 12H12.586a1.5 1.5 0 001.06-2.515l-1.06-1.06A2.002 2.002 0 0112 7.009V5.25c0-2.06-1.637-3.75-4-3.75zM8 14.5a2 2 0 01-1.938-1.5h3.876A2 2 0 018 14.5z"/>
      </svg>
      {#if appState.notificationCount > 0}
        <span class="notification-badge">{appState.notificationCount}</span>
      {/if}
    </button>

    <!-- License Badge -->
    <div class="license-pill" title="Community / Evaluation License">
      Non-commercial use
    </div>
  </div>
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

  /* Left Breadcrumbs */
  .breadcrumb-strip {
    display: flex;
    align-items: center;
    gap: 4px;
    overflow: hidden;
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

  /* Right Status */
  .status-right {
    display: flex;
    align-items: center;
    gap: 12px;
    flex-shrink: 0;
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

  .license-pill {
    display: flex;
    align-items: center;
    height: 18px;
    padding: 0 8px;
    border-radius: 9px;
    background-color: var(--license-bg);
    border: 1px solid var(--license-border);
    color: var(--license-fg);
    font-size: 10px;
    font-weight: 500;
    line-height: 1;
  }
</style>
