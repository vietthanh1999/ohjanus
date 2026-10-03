<script lang="ts">
  import { appState } from '../../state/appState.svelte';

  let isDropdownOpen = $state(false);

  function getTabIconColor(type: string) {
    switch (type) {
      case 'console': return '#57D38C';
      case 'table': return '#4A88C7';
      case 'approvals': return '#EDA200';
      case 'audit': return '#56A8F5';
      case 'tokens': return '#C792EA';
      case 'connections': return '#4A88C7';
      case 'dashboard': return '#57D38C';
      default: return '#7A7E85';
    }
  }
</script>

<nav class="tabbar-container">
  <!-- Tabs list -->
  <div class="tabs-scrollable">
    {#each appState.tabs as tab (tab.id)}
      <div
        class="tab-item"
        class:active={appState.activeTabId === tab.id}
        onclick={() => appState.activeTabId = tab.id}
        title={tab.title}
      >
        <!-- Icon -->
        <span class="tab-icon" style="color: {getTabIconColor(tab.type)};">
          {#if tab.type === 'console'}
            <svg width="13" height="13" viewBox="0 0 16 16" fill="currentColor">
              <path d="M11.251.068a.5.5 0 01.42.58L10.07 5H14a.5.5 0 01.372.832l-8.5 9.5a.5.5 0 01-.842-.512L6.63 9.5H2.5a.5.5 0 01-.42-.772l8.5-8.5a.5.5 0 01.671-.16z"/>
            </svg>
          {:else if tab.type === 'table'}
            <svg width="13" height="13" viewBox="0 0 16 16" fill="currentColor">
              <path d="M0 2a1 1 0 011-1h14a1 1 0 011 1v12a1 1 0 01-1 1H1a1 1 0 01-1-1V2zm1 3v2h6V5H1zm7 0v2h7V5H8zm0 3v2h7V8H8zm-1 0H1v2h6V8zm0 3H1v2h6v-2zm1 0v2h7v-2H8z"/>
            </svg>
          {:else if tab.type === 'approvals'}
            <svg width="13" height="13" viewBox="0 0 16 16" fill="currentColor">
              <path fill-rule="evenodd" d="M8 0c-.69 0-1.843.265-2.928.56-1.11.3-2.229.655-2.887.87a1.54 1.54 0 00-1.044 1.262c-.596 4.477.787 7.795 2.464 9.99 1.579 2.065 3.444 3.009 4.395 3.318.066.022.135.034.204.034s.138-.012.204-.034c.951-.309 2.816-1.253 4.395-3.318 1.677-2.195 3.06-5.513 2.464-9.99a1.54 1.54 0 00-1.044-1.263 62.467 62.467 0 00-2.887-.87C9.843.266 8.69 0 8 0zm2.146 5.146a.5.5 0 01.708.708l-3 3a.5.5 0 01-.708 0l-1.5-1.5a.5.5 0 11.708-.708L7.5 7.793l2.646-2.647z"/>
            </svg>
          {:else if tab.type === 'audit'}
            <svg width="13" height="13" viewBox="0 0 16 16" fill="currentColor">
              <path d="M4 1.5H3a2 2 0 00-2 2V14a2 2 0 002 2h10a2 2 0 002-2V3.5a2 2 0 00-2-2h-1v1h1a1 1 0 011 1V14a1 1 0 01-1 1H3a1 1 0 01-1-1V3.5a1 1 0 011-1h1v-1z"/>
              <path d="M9.5 1a.5.5 0 01.5.5v1a.5.5 0 01-.5.5h-3a.5.5 0 01-.5-.5v-1a.5.5 0 01.5-.5h3zm-3-1A1.5 1.5 0 005 1.5v1A1.5 1.5 0 006.5 4h3A1.5 1.5 0 0011 2.5v-1A1.5 1.5 0 009.5 0h-3z"/>
            </svg>
          {:else if tab.type === 'tokens'}
            <svg width="13" height="13" viewBox="0 0 16 16" fill="currentColor">
              <path fill-rule="evenodd" d="M6.5 0a6.5 6.5 0 00-5.348 10.185L.146 11.146A.5.5 0 000 11.5V15a1 1 0 001 1h3.5a.5.5 0 00.354-.146l1.353-1.354H7.5a.5.5 0 00.354-.146l1.353-1.354.908.908a6.5 6.5 0 10-3.615-13.908zM11 6a2 2 0 100-4 2 2 0 000 4z"/>
            </svg>
          {:else if tab.type === 'connections'}
            <svg width="13" height="13" viewBox="0 0 16 16" fill="currentColor">
              <path d="M8 1c3.866 0 7 1.12 7 2.5v9c0 1.38-3.134 2.5-7 2.5s-7-1.12-7-2.5v-9C1 2.12 4.134 1 8 1zm5.5 2.5c0-.44-2.126-1.2-5.5-1.2s-5.5.76-5.5 1.2 2.126 1.2 5.5 1.2 5.5-.76 5.5-1.2z"/>
            </svg>
          {:else if tab.type === 'dashboard'}
            <svg width="13" height="13" viewBox="0 0 16 16" fill="currentColor">
              <path d="M0 0h1v15h15v1H0V0zm10 3.5a.5.5 0 01.5-.5h4a.5.5 0 01.5.5v10a.5.5 0 01-.5.5h-4a.5.5 0 01-.5-.5v-10zm-4 4a.5.5 0 01.5-.5h2a.5.5 0 01.5.5v6a.5.5 0 01-.5.5h-2a.5.5 0 01-.5-.5v-6zm-4 3a.5.5 0 01.5-.5h2a.5.5 0 01.5.5v3a.5.5 0 01-.5.5h-2a.5.5 0 01-.5-.5v-3z"/>
            </svg>
          {/if}
        </span>

        <!-- Tab Title -->
        <span class="tab-label truncate">{tab.title}</span>

        <!-- Optional Badge -->
        {#if tab.badge}
          <span class="tab-badge">{tab.badge}</span>
        {/if}

        <!-- Close Button (x) -->
        {#if tab.closable}
          <button
            class="tab-close-btn"
            title="Close Tab (Cmd+W)"
            onclick={(e) => {
              e.stopPropagation();
              appState.closeTab(tab.id);
            }}
          >
            <svg width="10" height="10" viewBox="0 0 16 16" fill="currentColor">
              <path fill-rule="evenodd" d="M3.72 3.72a.75.75 0 011.06 0L8 6.94l3.22-3.22a.75.75 0 111.06 1.06L9.06 8l3.22 3.22a.75.75 0 11-1.06 1.06L8 9.06l-3.22 3.22a.75.75 0 01-1.06-1.06L6.94 8 3.72 4.78a.75.75 0 010-1.06z"/>
            </svg>
          </button>
        {/if}
      </div>
    {/each}
  </div>

  <!-- Right: Overflow chevron menu & quick switcher -->
  <div class="tabbar-actions">
    <button
      class="overflow-btn"
      title="All Open Tabs & Gateway Modules"
      onclick={() => isDropdownOpen = !isDropdownOpen}
    >
      <svg width="12" height="12" viewBox="0 0 16 16" fill="currentColor">
        <path d="M4 6l4 4 4-4H4z" />
      </svg>
    </button>

    {#if isDropdownOpen}
      <div class="tabs-dropdown" onclick={() => isDropdownOpen = false}>
        <div class="dropdown-header">Switch View / Module</div>
        <button class="dropdown-item" onclick={() => appState.activeTabId = 'console_2'}>
          <span>⚡ console_2 (SQL Editor)</span>
        </button>
        <button class="dropdown-item" onclick={() => appState.activeTabId = 'connection_credential'}>
          <span>🗄️ connection_credential (Data Grid)</span>
        </button>
        <div class="dropdown-separator"></div>
        <div class="dropdown-header">MCP Gateway (Admin)</div>
        <button class="dropdown-item" onclick={() => appState.activeTabId = 'approvals'}>
          <span>🛡️ Approvals Queue</span>
          <span class="jb-badge jb-badge-orange">{appState.notificationCount}</span>
        </button>
        <button class="dropdown-item" onclick={() => appState.activeTabId = 'audit'}>
          <span>📜 Audit Log Trail</span>
        </button>
        <button class="dropdown-item" onclick={() => appState.activeTabId = 'tokens'}>
          <span>🔑 MCP Agent Tokens</span>
        </button>
        <button class="dropdown-item" onclick={() => appState.activeTabId = 'connections'}>
          <span>🔌 Connection Pools & Latency</span>
        </button>
        <button class="dropdown-item" onclick={() => appState.activeTabId = 'dashboard'}>
          <span>📊 Gateway Dashboard Metrics</span>
        </button>
      </div>
    {/if}
  </div>
</nav>

<style>
  .tabbar-container {
    height: var(--tabbar-height);
    background-color: var(--bg-toolbar);
    border-bottom: 1px solid var(--border-default);
    display: flex;
    align-items: center;
    position: relative;
    z-index: 20;
    flex-shrink: 0;
  }

  .tabs-scrollable {
    display: flex;
    align-items: stretch;
    height: 100%;
    overflow-x: auto;
    overflow-y: hidden;
    flex: 1;
  }

  .tabs-scrollable::-webkit-scrollbar {
    display: none;
  }

  .tab-item {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 0 10px;
    height: 100%;
    min-width: 140px;
    max-width: 280px;
    background-color: var(--bg-inactive-tab);
    border-right: 1px solid var(--border-default);
    color: var(--text-secondary);
    font-size: 12px;
    cursor: pointer;
    position: relative;
    user-select: none;
    transition: background-color 0.1s ease;
  }

  .tab-item:hover {
    background-color: #35373B;
    color: var(--text-primary);
  }

  .tab-item.active {
    background-color: var(--bg-active-tab);
    color: var(--text-primary);
    border-bottom: 1px solid var(--bg-active-tab);
  }

  .tab-item.active::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 2px;
    background-color: var(--border-accent);
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

  .tab-badge {
    background-color: rgba(237, 162, 0, 0.2);
    border: 1px solid rgba(237, 162, 0, 0.5);
    color: #ffc44d;
    padding: 0 4px;
    border-radius: 6px;
    font-size: 10px;
    font-weight: 600;
  }

  .tab-close-btn {
    width: 16px;
    height: 16px;
    border-radius: 3px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--text-muted);
    opacity: 0.6;
    transition: all 0.15s ease;
  }

  .tab-item:hover .tab-close-btn {
    opacity: 1;
  }

  .tab-close-btn:hover {
    background-color: #393B40;
    color: var(--text-primary);
  }

  .tabbar-actions {
    height: 100%;
    display: flex;
    align-items: center;
    padding: 0 6px;
    background-color: var(--bg-toolbar);
    border-left: 1px solid var(--border-default);
    position: relative;
  }

  .overflow-btn {
    width: 22px;
    height: 22px;
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
    top: 32px;
    right: 4px;
    background-color: var(--bg-sidebar);
    border: 1px solid var(--border-strong);
    border-radius: 4px;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.55);
    width: 240px;
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

  .dropdown-separator {
    height: 1px;
    background-color: var(--border-default);
    margin: 4px 0;
  }
</style>
