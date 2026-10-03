<script lang="ts">
  import { appState } from '../../state/appState.svelte';

  let isDropdownOpen = $state(false);
</script>

<nav class="tabbar-container">
  <!-- Tabs list -->
  <div class="tabs-scrollable">
    {#each appState.tabs as tab (tab.id)}
      <button
        type="button"
        class="tab-pill"
        class:active={appState.activeTabId === tab.id}
        onclick={() => appState.activeTabId = tab.id}
        title={tab.title}
      >
        <!-- Icon -->
        <span class="tab-icon">
          {#if tab.type === 'console'}
            <svg width="12" height="12" viewBox="0 0 16 16" fill="#4A88C7">
              <path d="M11.251.068a.5.5 0 01.42.58L10.07 5H14a.5.5 0 01.372.832l-8.5 9.5a.5.5 0 01-.842-.512L6.63 9.5H2.5a.5.5 0 01-.42-.772l8.5-8.5a.5.5 0 01.671-.16z"/>
            </svg>
          {:else if tab.type === 'table'}
            <svg width="12" height="12" viewBox="0 0 16 16" fill="#4A88C7">
              <path d="M0 2a1 1 0 011-1h14a1 1 0 011 1v12a1 1 0 01-1 1H1a1 1 0 01-1-1V2zm1 3v2h6V5H1zm7 0v2h7V5H8zm0 3v2h7V8H8zm-1 0H1v2h6V8zm0 3H1v2h6v-2zm1 0v2h7v-2H8z"/>
            </svg>
          {:else if tab.type === 'approvals'}
            <svg width="12" height="12" viewBox="0 0 16 16" fill="#EDA200">
              <path fill-rule="evenodd" d="M8 0c-.69 0-1.843.265-2.928.56-1.11.3-2.229.655-2.887.87a1.54 1.54 0 00-1.044 1.262c-.596 4.477.787 7.795 2.464 9.99 1.579 2.065 3.444 3.009 4.395 3.318.066.022.135.034.204.034s.138-.012.204-.034c.951-.309 2.816-1.253 4.395-3.318 1.677-2.195 3.06-5.513 2.464-9.99a1.54 1.54 0 00-1.044-1.263 62.467 62.467 0 00-2.887-.87C9.843.266 8.69 0 8 0zm2.146 5.146a.5.5 0 01.708.708l-3 3a.5.5 0 01-.708 0l-1.5-1.5a.5.5 0 11.708-.708L7.5 7.793l2.646-2.647z"/>
            </svg>
          {:else if tab.type === 'audit'}
            <svg width="12" height="12" viewBox="0 0 16 16" fill="#56A8F5">
              <path d="M4 1.5H3a2 2 0 00-2 2V14a2 2 0 002 2h10a2 2 0 002-2V3.5a2 2 0 00-2-2h-1v1h1a1 1 0 011 1V14a1 1 0 01-1 1H3a1 1 0 01-1-1V3.5a1 1 0 011-1h1v-1z"/>
            </svg>
          {:else}
            <svg width="12" height="12" viewBox="0 0 16 16" fill="#7A7E85">
              <path d="M8 1c3.866 0 7 1.12 7 2.5v9c0 1.38-3.134 2.5-7 2.5s-7-1.12-7-2.5v-9C1 2.12 4.134 1 8 1zm5.5 2.5c0-.44-2.126-1.2-5.5-1.2s-5.5.76-5.5 1.2 2.126 1.2 5.5 1.2 5.5-.76 5.5-1.2z"/>
            </svg>
          {/if}
        </span>

        <!-- Tab Title -->
        <span class="tab-label truncate">{tab.title}</span>

        <!-- Close Button (x) -->
        {#if tab.closable}
          <span
            class="tab-close-icon"
            role="button"
            tabindex="-1"
            title="Close Tab"
            onclick={(e) => {
              e.stopPropagation();
              appState.closeTab(tab.id);
            }}
          >
            ×
          </span>
        {/if}
      </button>
    {/each}
  </div>

  <!-- Right: Overflow chevron menu -->
  <div class="tabbar-actions">
    <button
      type="button"
      class="overflow-btn"
      title="All Open Tabs"
      onclick={() => isDropdownOpen = !isDropdownOpen}
    >
      <svg width="10" height="10" viewBox="0 0 16 16" fill="currentColor">
        <path d="M4 6l4 4 4-4H4z" />
      </svg>
    </button>

    <button
      type="button"
      class="overflow-btn"
      title="Tab Actions"
    >
      <svg width="12" height="12" viewBox="0 0 16 16" fill="currentColor">
        <path d="M8 3a1.5 1.5 0 110-3 1.5 1.5 0 010 3zm0 6.5a1.5 1.5 0 110-3 1.5 1.5 0 010 3zm0 6.5a1.5 1.5 0 110-3 1.5 1.5 0 010 3z"/>
      </svg>
    </button>

    {#if isDropdownOpen}
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <div class="tabs-dropdown" onclick={() => isDropdownOpen = false}>
        <div class="dropdown-header">Database Objects</div>
        <button type="button" class="dropdown-item" onclick={() => appState.activeTabId = 'connection_credential'}>
          <span>🗄️ connection_credential [[Dev][ReadOnly] 10.220.6.4]</span>
        </button>
        <button type="button" class="dropdown-item" onclick={() => appState.activeTabId = 'console_2'}>
          <span>⚡ console_2 [[PRD] 10.250.6.23]</span>
        </button>
        <div class="dropdown-separator"></div>
        <div class="dropdown-header">Gateway Modules (§ui.md)</div>
        <button type="button" class="dropdown-item" onclick={() => appState.activeTabId = 'approvals'}>
          <span>🛡️ Approvals Queue</span>
          <span class="count-tag">{appState.notificationCount}</span>
        </button>
        <button type="button" class="dropdown-item" onclick={() => appState.activeTabId = 'audit'}>
          <span>📜 Audit Log Trail</span>
        </button>
        <button type="button" class="dropdown-item" onclick={() => appState.activeTabId = 'tokens'}>
          <span>🔑 MCP Agent Tokens</span>
        </button>
        <button type="button" class="dropdown-item" onclick={() => appState.activeTabId = 'connections'}>
          <span>🔌 Connection Pools</span>
        </button>
        <button type="button" class="dropdown-item" onclick={() => appState.activeTabId = 'dashboard'}>
          <span>📊 Telemetry Dashboard</span>
        </button>
      </div>
    {/if}
  </div>
</nav>

<style>
  .tabbar-container {
    height: var(--tabbar-height);
    background-color: var(--bg-canvas);
    border-bottom: 1px solid var(--border-subtle);
    display: flex;
    align-items: center;
    padding: 0 4px;
    position: relative;
    z-index: 20;
    flex-shrink: 0;
    border-top-left-radius: 8px;
    border-top-right-radius: 8px;
  }

  .tabs-scrollable {
    display: flex;
    align-items: center;
    gap: 4px;
    height: 100%;
    overflow-x: auto;
    overflow-y: hidden;
    flex: 1;
  }

  .tabs-scrollable::-webkit-scrollbar {
    display: none;
  }

  /* Exact DataGrip Rounded Tab Pill */
  .tab-pill {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 0 8px;
    height: 24px;
    max-width: 320px;
    border-radius: 4px;
    background-color: transparent;
    border: 1px solid transparent;
    color: var(--text-secondary);
    font-size: 12px;
    cursor: pointer;
    user-select: none;
    transition: background-color 0.1s ease;
    white-space: nowrap;
  }

  .tab-pill:hover {
    background-color: #2B2D30;
    color: var(--text-primary);
  }

  /* Exact Active Tab from design2.png: Dark blue fill with 1px blue outline */
  .tab-pill.active {
    background-color: #1F2E4A;
    border: 1px solid #3574F0;
    color: #DFE1E5;
  }

  .tab-icon {
    display: flex;
    align-items: center;
    flex-shrink: 0;
  }

  .tab-label {
    flex: 1;
    font-size: 12px;
  }

  .tab-close-icon {
    font-size: 14px;
    color: var(--text-muted);
    margin-left: 2px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 14px;
    height: 14px;
    border-radius: 2px;
  }

  .tab-close-icon:hover {
    background-color: rgba(255, 255, 255, 0.15);
    color: #FFFFFF;
  }

  .tabbar-actions {
    height: 100%;
    display: flex;
    align-items: center;
    padding-left: 4px;
    position: relative;
  }

  .overflow-btn {
    width: 20px;
    height: 20px;
    border-radius: 3px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--text-muted);
  }

  .overflow-btn:hover {
    background-color: var(--bg-hover);
    color: var(--text-primary);
  }

  .tabs-dropdown {
    position: absolute;
    top: 28px;
    right: 0;
    background-color: #2B2D30;
    border: 1px solid var(--border-strong);
    border-radius: 4px;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.55);
    width: 260px;
    padding: 4px 0;
    z-index: 100;
  }

  .dropdown-header {
    font-size: 10px;
    text-transform: uppercase;
    font-weight: 600;
    color: var(--text-muted);
    padding: 6px 12px 2px;
    letter-spacing: 0.5px;
  }

  .dropdown-item {
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 6px 12px;
    font-size: 12px;
    color: var(--text-primary);
    text-align: left;
  }

  .dropdown-item:hover {
    background-color: var(--bg-hover);
  }

  .count-tag {
    background-color: #EDA200;
    color: #1E1F22;
    font-size: 9px;
    font-weight: bold;
    padding: 0 4px;
    border-radius: 4px;
  }

  .dropdown-separator {
    height: 1px;
    background-color: var(--border-default);
    margin: 4px 0;
  }
</style>
