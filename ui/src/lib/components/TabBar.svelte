<script lang="ts">
  import { appStore } from '../appStore.svelte';
  import type { Tab, TabType } from '../types';

  let showNewTabMenu = $state(false);

  function addNewTab(type: TabType, title: string, subtitle?: string) {
    showNewTabMenu = false;
    const id = `tab-${type}-${Date.now().toString().slice(-4)}`;
    appStore.openTab({
      id,
      title,
      subtitle: subtitle || 'Active Session',
      icon: type,
      type,
      closable: true
    });
  }
</script>

<div class="tab-bar-container">
  <div class="tabs-list">
    {#each appStore.tabs as tab (tab.id)}
      <div 
        class="tab-item" 
        class:active={appStore.activeTabId === tab.id}
        onclick={() => appStore.selectTab(tab.id)}
      >
        <!-- Icon -->
        <span class="tab-icon">
          {#if tab.icon === 'console'}
            <span class="console-bolt">⚡</span>
          {:else if tab.icon === 'table'}
            <span class="table-grid">🗄️</span>
          {:else if tab.icon === 'approvals'}
            <span class="approvals-shield">🛡️</span>
          {:else if tab.icon === 'audit'}
            <span class="audit-log">📋</span>
          {:else if tab.icon === 'tokens'}
            <span class="tokens-key">🔑</span>
          {:else if tab.icon === 'connections'}
            <span class="conn-plug">🔌</span>
          {:else}
            <span class="dash-chart">📊</span>
          {/if}
        </span>

        <!-- Title & Subtitle -->
        <span class="tab-label">
          {tab.title}
          {#if tab.subtitle}
            <span class="tab-subtitle">[{tab.subtitle}]</span>
          {/if}
        </span>

        <!-- Badges (e.g. for pending approvals) -->
        {#if tab.type === 'approvals' && appStore.pendingApprovalsCount > 0}
          <span class="badge-pill tab-badge">{appStore.pendingApprovalsCount}</span>
        {/if}

        <!-- Close Button -->
        {#if tab.closable}
          <button 
            class="tab-close-btn" 
            title="Close Tab"
            onclick={(e) => { e.stopPropagation(); appStore.closeTab(tab.id); }}
          >
            ✕
          </button>
        {/if}
      </div>
    {/each}

    <!-- Add Tab Button -->
    <div class="new-tab-wrapper">
      <button class="new-tab-btn" title="Open View / New Tab" onclick={() => showNewTabMenu = !showNewTabMenu}>
        +
      </button>

      {#if showNewTabMenu}
        <div class="menu-popup new-tab-menu animate-fade-in">
          <div class="menu-header">NEW DATABASE CONSOLE</div>
          <div class="menu-item" onclick={() => addNewTab('console', 'console_3', '[PRD] 10.250.6.23')}>
            ⚡ SQL Console (PRD 10.250.6.23)
          </div>
          <div class="menu-item" onclick={() => addNewTab('console', 'console_dev', '[Dev] 10.220.6.4')}>
            ⚡ SQL Console (Dev 10.220.6.4)
          </div>

          <div class="menu-divider"></div>
          <div class="menu-header">JANUS GATEWAY MODULES</div>
          <div class="menu-item" onclick={() => addNewTab('approvals', 'Approvals', 'Write Queue')}>
            🛡️ Approval Queue ({appStore.pendingApprovalsCount} pending)
          </div>
          <div class="menu-item" onclick={() => addNewTab('audit', 'Audit Log', 'Security Gateway')}>
            📋 Audit Log Viewer
          </div>
          <div class="menu-item" onclick={() => addNewTab('tokens', 'MCP Tokens', 'Agent Access')}>
            🔑 AI Agent MCP Tokens
          </div>
          <div class="menu-item" onclick={() => addNewTab('connections', 'Connections', 'Pool Monitor')}>
            🔌 Database Connections
          </div>
          <div class="menu-item" onclick={() => addNewTab('dashboard', 'Dashboard', 'Gateway Metrics')}>
            📊 Gateway Dashboard
          </div>
        </div>
      {/if}
    </div>
  </div>

  <!-- Right Overflow Controls -->
  <div class="tab-controls">
    <button class="icon-btn" title="Tab List Dropdown">▾</button>
  </div>
</div>

<style>
  .tab-bar-container {
    height: var(--tab-bar-height);
    background: var(--surface-toolbar);
    border-bottom: 1px solid var(--border-default);
    display: flex;
    align-items: center;
    justify-content: space-between;
    overflow: hidden;
    position: relative;
    user-select: none;
  }

  .tabs-list {
    display: flex;
    align-items: center;
    height: 100%;
    overflow-x: auto;
    overflow-y: hidden;
  }

  .tabs-list::-webkit-scrollbar {
    display: none;
  }

  .tab-item {
    height: 100%;
    padding: 0 10px;
    display: flex;
    align-items: center;
    gap: 6px;
    background: var(--surface-inactive-tab);
    color: var(--text-secondary);
    font-size: 12px;
    border-right: 1px solid var(--border-default);
    cursor: pointer;
    white-space: nowrap;
    position: relative;
    transition: background 0.1s, color 0.1s;
    min-width: 130px;
    max-width: 280px;
  }

  .tab-item:hover {
    background: #35373B;
    color: var(--text-primary);
  }

  .tab-item.active {
    background: var(--surface-active-tab);
    color: var(--text-primary);
    border-top: 2px solid var(--border-accent);
  }

  .tab-icon {
    display: inline-flex;
    align-items: center;
    font-size: 12px;
  }

  .console-bolt { color: #56A8F5; }
  .approvals-shield { color: #EDA200; }

  .tab-label {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    flex: 1;
  }

  .tab-subtitle {
    color: var(--text-muted);
    font-size: 11px;
    margin-left: 2px;
  }

  .tab-badge {
    background: var(--action-warning);
    color: #14281B;
    font-weight: 700;
  }

  .tab-close-btn {
    width: 14px;
    height: 14px;
    border-radius: 3px;
    color: var(--text-muted);
    font-size: 10px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    margin-left: 2px;
    opacity: 0;
    transition: opacity 0.15s, background 0.15s;
  }

  .tab-item:hover .tab-close-btn,
  .tab-item.active .tab-close-btn {
    opacity: 1;
  }

  .tab-close-btn:hover {
    background: #4E5157;
    color: var(--text-primary);
  }

  .new-tab-wrapper {
    position: relative;
    padding: 0 4px;
  }

  .new-tab-btn {
    width: 22px;
    height: 22px;
    border-radius: 3px;
    color: var(--text-muted);
    font-size: 14px;
  }

  .new-tab-btn:hover {
    background: var(--surface-hover);
    color: var(--text-primary);
  }

  .new-tab-menu {
    position: absolute;
    top: 26px;
    left: 0;
    width: 240px;
    z-index: 1000;
  }

  .tab-controls {
    padding: 0 8px;
    display: flex;
    align-items: center;
  }
</style>
