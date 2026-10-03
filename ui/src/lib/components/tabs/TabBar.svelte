<script lang="ts">
  import { appState } from '../../state/appState.svelte';
  import { Icon } from '@ohjanus/icons';

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
            <Icon name="lightning" size={14} color="#4A88C7" />
          {:else if tab.type === 'table'}
            <Icon name="table" size={14} color="#4A88C7" />
          {:else if tab.type === 'approvals'}
            <Icon name="shield" size={14} color="#EDA200" />
          {:else if tab.type === 'audit'}
            <Icon name="audit" size={14} color="#56A8F5" />
          {:else if tab.type === 'tokens'}
            <Icon name="key" size={14} color="#3B82F6" />
          {:else if tab.type === 'connections'}
            <Icon name="database" size={14} color="#3B82F6" />
          {:else if tab.type === 'dashboard'}
            <Icon name="chart" size={14} color="#57D38C" />
          {:else}
            <Icon name="database" size={14} color="#7A7E85" />
          {/if}
        </span>

        <!-- Tab Title -->
        <span class="tab-label truncate">{tab.title}</span>

        <!-- Close Button (x) -->
        {#if tab.closable}
          <span
            class="tab-close-icon"
            role="button"
            tabindex="0"
            title="Close Tab"
            onclick={(e) => {
              e.stopPropagation();
              appState.closeTab(tab.id);
            }}
            onkeydown={(e) => {
              if (e.key === 'Enter' || e.key === ' ') {
                e.stopPropagation();
                appState.closeTab(tab.id);
              }
            }}
          >
            <Icon name="close" size={11} />
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
      <Icon name="chevron-down" size={10} />
    </button>

    <button
      type="button"
      class="overflow-btn"
      title="Tab Actions"
    >
      <Icon name="more" size={13} />
    </button>

    {#if isDropdownOpen}
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <div class="tabs-dropdown" onclick={() => isDropdownOpen = false}>
        <div class="dropdown-header">Database Objects</div>
        <button type="button" class="dropdown-item" onclick={() => appState.openTab({ id: 'connection_credential', title: 'connectio...credential [[Dev][ReadOnly] 10.220.6.4]', type: 'table', closable: true, icon: 'table', env: 'Dev' })}>
          <span class="dropdown-item-content">
            <Icon name="table" size={14} color="#4A88C7" />
            <span>connection_credential [[Dev][ReadOnly] 10.220.6.4]</span>
          </span>
        </button>
        <button type="button" class="dropdown-item" onclick={() => appState.openTab({ id: 'console_2', title: 'console_2 [[PRD] 10.250.6.23]', type: 'console', closable: true, icon: 'lightning', env: 'PRD' })}>
          <span class="dropdown-item-content">
            <Icon name="lightning" size={14} color="#3B82F6" />
            <span>console_2 [[PRD] 10.250.6.23]</span>
          </span>
        </button>
        <div class="dropdown-separator"></div>
        <div class="dropdown-header">Gateway Modules (§ui.md)</div>
        <button type="button" class="dropdown-item" onclick={() => appState.openTab({ id: 'approvals', title: 'Approvals Queue', type: 'approvals', closable: true, icon: 'shield', badge: '2' })}>
          <span class="dropdown-item-content">
            <Icon name="shield" size={14} color="#EDA200" />
            <span>Approvals Queue</span>
          </span>
          <span class="count-tag">{appState.notificationCount}</span>
        </button>
        <button type="button" class="dropdown-item" onclick={() => appState.openTab({ id: 'audit', title: 'Audit Logs', type: 'audit', closable: true, icon: 'audit' })}>
          <span class="dropdown-item-content">
            <Icon name="audit" size={14} color="#56A8F5" />
            <span>Audit Log Trail</span>
          </span>
        </button>
        <button type="button" class="dropdown-item" onclick={() => appState.openTab({ id: 'tokens', title: 'MCP Tokens', type: 'tokens', closable: true, icon: 'key' })}>
          <span class="dropdown-item-content">
            <Icon name="key" size={14} color="#3B82F6" />
            <span>MCP Agent Tokens</span>
          </span>
        </button>
        <button type="button" class="dropdown-item" onclick={() => appState.openTab({ id: 'connections', title: 'Connection Pools', type: 'connections', closable: true, icon: 'database' })}>
          <span class="dropdown-item-content">
            <Icon name="database" size={14} color="#3B82F6" />
            <span>Connection Pools</span>
          </span>
        </button>
        <button type="button" class="dropdown-item" onclick={() => appState.openTab({ id: 'dashboard', title: 'Gateway Dashboard', type: 'dashboard', closable: true, icon: 'chart' })}>
          <span class="dropdown-item-content">
            <Icon name="chart" size={14} color="#57D38C" />
            <span>Telemetry Dashboard</span>
          </span>
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
    border-radius: 0;
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
    padding: 0 10px;
    height: var(--control-height-xs, 26px);
    max-width: 320px;
    border-radius: var(--radius-sm, 4px);
    background-color: transparent;
    border: 1px solid transparent;
    color: var(--text-secondary);
    font-size: var(--font-size-sm, 12px);
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
    width: var(--icon-btn-size-sm, 26px);
    height: var(--icon-btn-size-sm, 26px);
    border-radius: var(--radius-sm, 4px);
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

  .dropdown-item-content {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
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
