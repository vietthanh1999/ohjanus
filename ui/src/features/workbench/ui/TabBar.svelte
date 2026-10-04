<script lang="ts">
  import { workbenchState } from '@/features/workbench';
  import { Icon } from '@ohjanus/icons';
  import { Text } from '@ohjanus/ui';

  let isDropdownOpen = $state(false);

  // Pending-approvals badge, synced by features/approvals (no import: avoids a
  // workbench <-> approvals module cycle; TabBar only reads shell tab state).
  let approvalsBadge = $derived(workbenchState.tabs.find((t) => t.id === 'approvals')?.badge);
</script>

<nav class="tabbar-container">
  <!-- Tabs list -->
  <div class="tabs-scrollable">
    {#each workbenchState.tabs as tab (tab.id)}
      <button
        type="button"
        class="tab-pill"
        class:active={workbenchState.activeTabId === tab.id}
        onclick={() => workbenchState.activeTabId = tab.id}
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
        <Text size="md" truncate style="flex: 1;">{tab.title}</Text>

        <!-- Close Button (x) -->
        {#if tab.closable}
          <span
            class="tab-close-icon"
            role="button"
            tabindex="0"
            title="Close Tab"
            onclick={(e) => {
              e.stopPropagation();
              workbenchState.closeTab(tab.id);
            }}
            onkeydown={(e) => {
              if (e.key === 'Enter' || e.key === ' ') {
                e.stopPropagation();
                workbenchState.closeTab(tab.id);
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
        {#if workbenchState.tabs.some((t) => t.type === 'table' || t.type === 'console')}
          <div class="dropdown-header">Database Objects</div>
          {#each workbenchState.tabs.filter((t) => t.type === 'table' || t.type === 'console') as tab (tab.id)}
            <button type="button" class="dropdown-item" onclick={() => workbenchState.openTab(tab)}>
              <span class="dropdown-item-content">
                {#if tab.type === 'table'}
                  <Icon name="table" size={14} color="#4A88C7" />
                {:else}
                  <Icon name="lightning" size={14} color="#3B82F6" />
                {/if}
                <Text size="md">{tab.title}</Text>
              </span>
            </button>
          {/each}
          <div class="dropdown-separator"></div>
        {/if}
        <div class="dropdown-header">Gateway Modules</div>
        <button type="button" class="dropdown-item" onclick={() => workbenchState.openTab({ id: 'approvals', title: 'Approvals Queue', type: 'approvals', closable: false, icon: 'shield' })}>
          <span class="dropdown-item-content">
            <Icon name="shield" size={14} color="#EDA200" />
            <Text size="md">Approvals Queue</Text>
          </span>
          {#if approvalsBadge}
            <span class="count-tag">{approvalsBadge}</span>
          {/if}
        </button>
        <button type="button" class="dropdown-item" onclick={() => workbenchState.openTab({ id: 'audit', title: 'Audit Logs', type: 'audit', closable: false, icon: 'audit' })}>
          <span class="dropdown-item-content">
            <Icon name="audit" size={14} color="#56A8F5" />
            <Text size="md">Audit Log Trail</Text>
          </span>
        </button>
        <button type="button" class="dropdown-item" onclick={() => workbenchState.openTab({ id: 'tokens', title: 'MCP Tokens', type: 'tokens', closable: false, icon: 'key' })}>
          <span class="dropdown-item-content">
            <Icon name="key" size={14} color="#3B82F6" />
            <Text size="md">MCP Agent Tokens</Text>
          </span>
        </button>
        <button type="button" class="dropdown-item" onclick={() => workbenchState.openTab({ id: 'connections', title: 'Connection Pools', type: 'connections', closable: false, icon: 'database' })}>
          <span class="dropdown-item-content">
            <Icon name="database" size={14} color="#3B82F6" />
            <Text size="md">Connection Pools</Text>
          </span>
        </button>
        <button type="button" class="dropdown-item" onclick={() => workbenchState.openTab({ id: 'dashboard', title: 'Gateway Dashboard', type: 'dashboard', closable: false, icon: 'chart' })}>
          <span class="dropdown-item-content">
            <Icon name="chart" size={14} color="#57D38C" />
            <Text size="md">Telemetry Dashboard</Text>
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
    font-size: var(--font-size-sm, 13px);
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
    width: var(--icon-btn-size-sm, 28px);
    height: var(--icon-btn-size-sm, 28px);
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
    font-size: var(--font-size-2xs, 11px);
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
    font-size: var(--font-size-sm, 13px);
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
