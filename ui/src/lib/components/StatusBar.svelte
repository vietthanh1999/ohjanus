<script lang="ts">
  import { appStore } from '../appStore.svelte';

  let currentTab = $derived(appStore.activeTab);
</script>

<footer class="status-bar font-mono">
  <!-- Left Breadcrumbs -->
  <div class="breadcrumb-container">
    {#if currentTab?.type === 'console'}
      <span class="crumb">Database Consoles</span>
      <span class="sep">&gt;</span>
      <span class="crumb">{currentTab.metadata?.connection || '[PRD] 10.250.6.23'}</span>
      <span class="sep">&gt;</span>
      <span class="crumb active">⚡ {currentTab.title} [{currentTab.metadata?.connection || '[PRD] 10.250.6.23'}]</span>
    {:else if currentTab?.type === 'table'}
      <span class="crumb">Database</span>
      <span class="sep">&gt;</span>
      <span class="crumb">{currentTab.metadata?.connection || '[Dev][ReadOnly] 10.220.6.4'}</span>
      <span class="sep">&gt;</span>
      <span class="crumb">dev_mh_asset</span>
      <span class="sep">&gt;</span>
      <span class="crumb">public</span>
      <span class="sep">&gt;</span>
      <span class="crumb">tables</span>
      <span class="sep">&gt;</span>
      <span class="crumb active">📄 {currentTab.metadata?.table || 'connection_credential'}</span>
    {:else if currentTab?.type === 'approvals'}
      <span class="crumb">Janus Security Gateway</span>
      <span class="sep">&gt;</span>
      <span class="crumb">Policy Engine</span>
      <span class="sep">&gt;</span>
      <span class="crumb active">🛡️ Human Approval Queue ({appStore.pendingApprovalsCount} pending)</span>
    {:else if currentTab?.type === 'audit'}
      <span class="crumb">Janus Security Gateway</span>
      <span class="sep">&gt;</span>
      <span class="crumb">Compliance</span>
      <span class="sep">&gt;</span>
      <span class="crumb active">📋 Structured Audit Log Viewer</span>
    {:else if currentTab?.type === 'tokens'}
      <span class="crumb">Janus Security Gateway</span>
      <span class="sep">&gt;</span>
      <span class="crumb">Identity &amp; Access</span>
      <span class="sep">&gt;</span>
      <span class="crumb active">🔑 MCP Bearer Tokens</span>
    {:else if currentTab?.type === 'connections'}
      <span class="crumb">Janus Gateway</span>
      <span class="sep">&gt;</span>
      <span class="crumb">Connection Pools</span>
      <span class="sep">&gt;</span>
      <span class="crumb active">🔌 Database Endpoints &amp; Catalogs</span>
    {:else}
      <span class="crumb">Janus Gateway</span>
      <span class="sep">&gt;</span>
      <span class="crumb active">📊 Observability &amp; Metrics</span>
    {/if}
  </div>

  <!-- Right Indicators -->
  <div class="right-indicators">
    <!-- Document stats -->
    <span class="doc-stats" title="Cursor position & buffer metrics">{appStore.cursorPosition}</span>
    
    <span class="item" title="Line Feed">LF</span>
    <span class="item" title="Character Encoding">UTF-8</span>
    <span class="item" title="Indentation">4 spaces</span>

    <!-- Read-Only Lock Icon -->
    <span class="lock-icon" title="Transaction is in READ ONLY state">🔒</span>

    <!-- Notifications Bell -->
    <span class="bell-icon" class:active={appStore.notificationsCount > 0} title="{appStore.notificationsCount} Notifications">
      🔔
    </span>

    <!-- License Pill -->
    <span class="license-pill">
      Non-commercial use
    </span>
  </div>
</footer>

<style>
  .status-bar {
    height: var(--status-bar-height);
    background: var(--surface-sidebar);
    border-top: 1px solid var(--border-default);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 10px;
    font-size: 11px;
    color: var(--text-secondary);
    user-select: none;
    z-index: 100;
  }

  .breadcrumb-container {
    display: flex;
    align-items: center;
    gap: 6px;
    overflow: hidden;
    white-space: nowrap;
  }

  .crumb {
    color: var(--text-secondary);
    padding: 1px 4px;
    border-radius: 2px;
    cursor: pointer;
  }

  .crumb:hover {
    background: var(--surface-hover);
    color: var(--text-primary);
  }

  .crumb.active {
    color: var(--text-primary);
    font-weight: 500;
  }

  .sep {
    color: var(--text-muted);
    font-size: 10px;
  }

  .right-indicators {
    display: flex;
    align-items: center;
    gap: 12px;
    white-space: nowrap;
  }

  .doc-stats {
    color: var(--text-muted);
    font-size: 11px;
  }

  .item {
    color: var(--text-muted);
    cursor: pointer;
    padding: 1px 3px;
    border-radius: 2px;
  }

  .item:hover {
    background: var(--surface-hover);
    color: var(--text-primary);
  }

  .lock-icon {
    font-size: 11px;
    cursor: pointer;
    opacity: 0.8;
  }

  .bell-icon {
    font-size: 11px;
    cursor: pointer;
    position: relative;
  }

  .bell-icon.active {
    color: var(--action-primary);
  }
</style>
