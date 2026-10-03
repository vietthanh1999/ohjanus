<script lang="ts">
  import { appStore } from '../appStore.svelte';

  let showProjectDropdown = $state(false);
  let showBranchDropdown = $state(false);
  let showSearchModal = $state(false);
  let showAiModal = $state(false);
  let aiPrompt = $state('');
  let searchQuery = $state('');

  function openTab(type: any, title: string, subtitle?: string) {
    appStore.openTab({
      id: `tab-${type}`,
      title,
      subtitle,
      icon: type,
      type,
      closable: true
    });
  }
</script>

<header class="titlebar">
  <!-- Left: Traffic Lights & Context -->
  <div class="left-section">
    <div class="traffic-lights">
      <div class="traffic-dot close" title="Close"></div>
      <div class="traffic-dot minimize" title="Minimize"></div>
      <div class="traffic-dot maximize" title="Zoom"></div>
    </div>

    <div class="context-group">
      <!-- Avatar Badge -->
      <div class="avatar-badge" title="User Profile: Thanh Tran">VT</div>

      <!-- Project Selector -->
      <div class="dropdown-trigger" onclick={() => showProjectDropdown = !showProjectDropdown}>
        <span class="project-name">VTVprime</span>
        <span class="chevron">▾</span>
      </div>

      {#if showProjectDropdown}
        <div class="menu-popup project-menu animate-fade-in">
          <div class="menu-item active">● VTVprime (Active)</div>
          <div class="menu-item">○ OpenJanus Core Gateway</div>
          <div class="menu-item">○ Data Analytics Platform</div>
          <div class="menu-divider"></div>
          <div class="menu-item action">+ Open Project...</div>
        </div>
      {/if}

      <!-- Version Control -->
      <div class="dropdown-trigger vcs-trigger" onclick={() => showBranchDropdown = !showBranchDropdown}>
        <span>Version Control</span>
        <span class="chevron">▾</span>
      </div>

      {#if showBranchDropdown}
        <div class="menu-popup vcs-menu animate-fade-in">
          <div class="menu-header">Git Branch: main</div>
          <div class="menu-item active">✓ main (origin/main)</div>
          <div class="menu-item">feat/approval-engine</div>
          <div class="menu-item">fix/redact-tokens</div>
          <div class="menu-divider"></div>
          <div class="menu-item action">⟳ Fetch &amp; Pull</div>
        </div>
      {/if}
    </div>
  </div>

  <!-- Center Quick Actions -->
  <div class="center-section">
    <div class="quick-actions">
      <button class="icon-btn" title="Database Explorer" onclick={() => appStore.selectTab('tab-console-2')}>
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <ellipse cx="12" cy="5" rx="9" ry="3"/>
          <path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/>
          <path d="M3 12c0 1.66 4 3 9 3s9-1.34 9-3"/>
        </svg>
      </button>

      <button class="icon-btn run-btn" title="Run Selected Query (Cmd+Enter)" onclick={() => appStore.runQuery()}>
        <svg width="15" height="15" viewBox="0 0 24 24" fill="var(--action-success)" stroke="var(--action-success)" stroke-width="1">
          <polygon points="5 3 19 12 5 21 5 3"/>
        </svg>
      </button>

      <button class="icon-btn" title="Open Approvals Queue" onclick={() => openTab('approvals', 'Approvals', 'Write Queue')}>
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
          <path d="m9 12 2 2 4-4"/>
        </svg>
      </button>

      <button class="icon-btn" title="Audit Log Viewer" onclick={() => openTab('audit', 'Audit Log', 'Security Gateway')}>
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
          <polyline points="14 2 14 8 20 8"/>
          <line x1="16" y1="13" x2="8" y2="13"/>
          <line x1="16" y1="17" x2="8" y2="17"/>
        </svg>
      </button>

      <button class="icon-btn" title="MCP Tokens" onclick={() => openTab('tokens', 'MCP Tokens', 'Agent Access')}>
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="7.5" cy="15.5" r="5.5"/>
          <path d="m21 2-9.6 9.6"/>
          <path d="m15.5 7.5 3 3L22 7l-3-3"/>
        </svg>
      </button>

      <button class="icon-btn" title="More Tools">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
          <circle cx="12" cy="12" r="2"/>
          <circle cx="19" cy="12" r="2"/>
          <circle cx="5" cy="12" r="2"/>
        </svg>
      </button>
    </div>
  </div>

  <!-- Right Utilities: AI Assistant, Search, Settings -->
  <div class="right-section">
    <!-- AI Assistant -->
    <button class="icon-btn ai-btn" title="Janus AI Assistant" onclick={() => showAiModal = true}>
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#56A8F5" stroke-width="2">
        <path d="M12 2a10 10 0 1 0 10 10A10 10 0 0 0 12 2zm0 14a4 4 0 1 1 4-4 4 4 0 0 1-4 4z"/>
        <path d="M12 6v2m0 8v2M6 12h2m8 0h2"/>
      </svg>
    </button>

    <!-- Global Search -->
    <button class="icon-btn" title="Search Everything (Cmd+K)" onclick={() => showSearchModal = true}>
      <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="11" cy="11" r="8"/>
        <line x1="21" y1="21" x2="16.65" y2="16.65"/>
      </svg>
    </button>

    <!-- Settings Gear -->
    <button class="icon-btn" title="Database &amp; Gateway Settings" onclick={() => openTab('connections', 'Connections', 'Pool Monitor')}>
      <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="3"/>
        <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/>
      </svg>
    </button>
  </div>
</header>

<!-- Global Search Modal -->
{#if showSearchModal}
  <div class="modal-backdrop" onclick={() => showSearchModal = false}>
    <div class="modal-dialog search-dialog animate-fade-in" onclick={(e) => e.stopPropagation()}>
      <div class="search-input-wrapper">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8"/>
          <line x1="21" y1="21" x2="16.65" y2="16.65"/>
        </svg>
        <input 
          type="text" 
          placeholder="Search tables, schemas, queries, or approvals (e.g. transfer_job)..."
          bind:value={searchQuery}
          autofocus
        />
        <span class="esc-badge">ESC</span>
      </div>

      <div class="search-results">
        <div class="search-group-title">TABLES &amp; ENTITIES</div>
        <div class="search-item" onclick={() => { appStore.openTableTab('connection_credential', '[Dev][ReadOnly] 10.220.6.4'); showSearchModal = false; }}>
          <span class="type-tag table">TABLE</span>
          <span class="item-name">public.connection_credential</span>
          <span class="item-sub">10.220.6.4 (Dev)</span>
        </div>
        <div class="search-item" onclick={() => { appStore.openTableTab('transfer_job', '[PRD] 10.250.6.23'); showSearchModal = false; }}>
          <span class="type-tag table">TABLE</span>
          <span class="item-name">public.transfer_job</span>
          <span class="item-sub">10.250.6.23 (PRD)</span>
        </div>

        <div class="search-group-title">ACTIONS &amp; TOOLS</div>
        <div class="search-item" onclick={() => { openTab('approvals', 'Approvals', 'Write Queue'); showSearchModal = false; }}>
          <span class="type-tag action">GATEWAY</span>
          <span class="item-name">Pending Approvals Queue</span>
          <span class="badge-pill">{appStore.pendingApprovalsCount} pending</span>
        </div>
        <div class="search-item" onclick={() => { openTab('tokens', 'MCP Tokens', 'Agent Access'); showSearchModal = false; }}>
          <span class="type-tag action">GATEWAY</span>
          <span class="item-name">Manage AI Agent Tokens</span>
        </div>
      </div>
    </div>
  </div>
{/if}

<!-- AI Assistant Modal -->
{#if showAiModal}
  <div class="modal-backdrop" onclick={() => showAiModal = false}>
    <div class="modal-dialog ai-dialog animate-fade-in" onclick={(e) => e.stopPropagation()}>
      <div class="modal-header">
        <div style="display:flex; align-items:center; gap:8px;">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#56A8F5" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <path d="M12 6v2m0 8v2M6 12h2m8 0h2"/>
          </svg>
          <span>Janus SQL AI Copilot</span>
        </div>
        <button class="icon-btn" onclick={() => showAiModal = false}>✕</button>
      </div>
      <div class="modal-body">
        <p style="color:var(--text-secondary); margin-bottom:12px;">
          Ask Janus to construct safe, AST-validated queries, explain execution plans, or check schema constraints.
        </p>
        <textarea 
          rows="3" 
          placeholder="e.g., Lấy danh sách job media-ingest đã retry trên 3 lần trong 24h qua..."
          bind:value={aiPrompt}
          style="width: 100%;"
        ></textarea>
        <div style="margin-top:12px; display:flex; gap:6px;">
          <button class="btn-primary" onclick={() => {
            appStore.consoleSql += `\n-- AI Generated Query:\nSELECT * FROM transfer_job WHERE retry_count >= 3 AND created_at > NOW() - INTERVAL '24 hours';`;
            showAiModal = false;
          }}>Generate SQL &amp; Insert to Editor</button>
          <button class="btn-secondary" onclick={() => showAiModal = false}>Cancel</button>
        </div>
      </div>
    </div>
  </div>
{/if}

<style>
  .titlebar {
    height: var(--window-header-height);
    background: var(--surface-sidebar);
    border-bottom: 1px solid var(--border-default);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 12px;
    z-index: 100;
  }

  .left-section, .center-section, .right-section {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .traffic-lights {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-right: 12px;
  }

  .traffic-dot {
    width: 12px;
    height: 12px;
    border-radius: 50%;
    cursor: pointer;
    transition: filter 0.15s;
  }

  .traffic-dot.close { background: var(--window-close); }
  .traffic-dot.minimize { background: var(--window-minimize); }
  .traffic-dot.maximize { background: var(--window-maximize); }
  .traffic-dot:hover { filter: brightness(1.2); }

  .context-group {
    display: flex;
    align-items: center;
    gap: 8px;
    position: relative;
  }

  .avatar-badge {
    width: 20px;
    height: 20px;
    background: var(--badge-profile-bg);
    color: var(--badge-profile-fg);
    border-radius: 50%;
    font-size: 10px;
    font-weight: 700;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .dropdown-trigger {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 2px 6px;
    border-radius: 4px;
    cursor: pointer;
    font-size: 12px;
    color: var(--text-primary);
  }

  .dropdown-trigger:hover {
    background: var(--surface-hover);
  }

  .project-name {
    font-weight: 500;
  }

  .vcs-trigger {
    color: var(--text-secondary);
  }

  .chevron {
    font-size: 10px;
    color: var(--text-muted);
  }

  .quick-actions {
    display: flex;
    align-items: center;
    gap: 4px;
    background: #25272A;
    padding: 2px 4px;
    border-radius: 4px;
    border: 1px solid var(--border-default);
  }

  .run-btn:hover svg {
    filter: drop-shadow(0 0 3px rgba(87, 211, 140, 0.6));
  }

  .menu-popup {
    position: absolute;
    top: 28px;
    left: 20px;
    background: var(--surface-sidebar);
    border: 1px solid var(--border-default);
    border-radius: 4px;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.55);
    min-width: 180px;
    padding: 4px 0;
    z-index: 1000;
  }

  .menu-header {
    padding: 4px 10px;
    font-size: 11px;
    color: var(--text-muted);
  }

  .menu-item {
    padding: 5px 12px;
    font-size: 12px;
    color: var(--text-primary);
    cursor: pointer;
  }

  .menu-item:hover {
    background: var(--surface-hover);
  }

  .menu-item.active {
    color: var(--action-primary);
    font-weight: 500;
  }

  .menu-item.action {
    color: var(--text-secondary);
  }

  .menu-divider {
    height: 1px;
    background: var(--border-default);
    margin: 4px 0;
  }

  /* Search Modal */
  .search-dialog {
    width: 540px;
    padding: 0;
  }

  .search-input-wrapper {
    display: flex;
    align-items: center;
    padding: 10px 14px;
    border-bottom: 1px solid var(--border-default);
    gap: 10px;
  }

  .search-input-wrapper input {
    flex: 1;
    border: none;
    font-size: 13px;
    background: transparent;
    padding: 0;
  }

  .esc-badge {
    font-size: 10px;
    color: var(--text-muted);
    background: #393B40;
    padding: 2px 4px;
    border-radius: 3px;
  }

  .search-results {
    max-height: 320px;
    overflow-y: auto;
    padding: 6px 0;
  }

  .search-group-title {
    padding: 6px 14px 2px;
    font-size: 10px;
    font-weight: 600;
    color: var(--text-muted);
  }

  .search-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 7px 14px;
    cursor: pointer;
  }

  .search-item:hover {
    background: var(--surface-hover);
  }

  .type-tag {
    font-size: 9px;
    font-weight: 700;
    padding: 1px 4px;
    border-radius: 2px;
  }

  .type-tag.table {
    background: rgba(86, 168, 245, 0.15);
    color: #56A8F5;
  }

  .type-tag.action {
    background: rgba(87, 211, 140, 0.15);
    color: var(--action-success);
  }

  .item-name {
    flex: 1;
    color: var(--text-primary);
  }

  .item-sub {
    color: var(--text-muted);
    font-size: 11px;
  }
</style>
