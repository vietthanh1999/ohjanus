<script lang="ts">
  import { appState } from '../../state/appState.svelte';

  let isAiOpen = $state(false);
</script>

<header class="titlebar">
  <!-- Left: Traffic Lights & Context Selectors -->
  <div class="left-section">
    <div class="traffic-lights">
      <span class="light close" title="Close"></span>
      <span class="light minimize" title="Minimize"></span>
      <span class="light maximize" title="Maximize"></span>
    </div>

    <div class="context-group">
      <!-- Profile Avatar -->
      <div class="avatar" title="Viet Thanh (Administrator)">
        VT
      </div>

      <!-- Project Selector -->
      <button class="selector-btn" title="Current Workspace Project">
        <span class="project-name">VTVprime</span>
        <svg class="chevron" width="10" height="10" viewBox="0 0 16 16" fill="currentColor">
          <path d="M4 6l4 4 4-4H4z" />
        </svg>
      </button>

      <span class="divider">/</span>

      <!-- VCS Branch Selector -->
      <button class="selector-btn branch-btn" title="Git Branch: main">
        <svg width="12" height="12" viewBox="0 0 16 16" fill="currentColor" class="branch-icon">
          <path fill-rule="evenodd" d="M11.75 2.5a.75.75 0 100 1.5.75.75 0 000-1.5zm-2.25.75a2.25 2.25 0 113 2.122V6A2.5 2.5 0 0110 8.5H6a1 1 0 00-1 1v1.128a2.251 2.251 0 11-1.5 0V5.372a2.25 2.25 0 111.5 0v1.836A2.492 2.492 0 016 7h4a1 1 0 001-1v-.628A2.25 2.25 0 019.5 3.25zM4.25 12a.75.75 0 100 1.5.75.75 0 000-1.5zM3.5 3.25a.75.75 0 111.5 0 .75.75 0 01-1.5 0z"/>
        </svg>
        <span>Version Control</span>
        <svg class="chevron" width="10" height="10" viewBox="0 0 16 16" fill="currentColor">
          <path d="M4 6l4 4 4-4H4z" />
        </svg>
      </button>
    </div>
  </div>

  <!-- Center: Quick Actions Bar -->
  <div class="center-section">
    <div class="quick-actions">
      <!-- Database icon -->
      <button
        class="quick-btn"
        title="Database Explorer (Cmd+1)"
        onclick={() => appState.isSidebarCollapsed = !appState.isSidebarCollapsed}
      >
        <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor">
          <path d="M8 1c3.866 0 7 1.12 7 2.5v9c0 1.38-3.134 2.5-7 2.5s-7-1.12-7-2.5v-9C1 2.12 4.134 1 8 1zm5.5 2.5c0-.44-2.126-1.2-5.5-1.2s-5.5.76-5.5 1.2 2.126 1.2 5.5 1.2 5.5-.76 5.5-1.2zM2.5 6.275V7.5c0 .44 2.126 1.2 5.5 1.2s5.5-.76 5.5-1.2V6.275c-1.39.638-3.376.925-5.5.925s-4.11-.287-5.5-.925zm0 3.5V11c0 .44 2.126 1.2 5.5 1.2s5.5-.76 5.5-1.2V9.775c-1.39.638-3.376.925-5.5.925s-4.11-.287-5.5-.925z"/>
        </svg>
      </button>

      <!-- Run configurations -->
      <button
        class="quick-btn run-btn"
        title="Execute Active Query (Cmd+Enter)"
        onclick={() => appState.executeQuery()}
      >
        <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor">
          <path d="M8 15A7 7 0 108 1a7 7 0 000 14zm-1.5-9.8v5.6L11 8l-4.5-2.8z"/>
        </svg>
      </button>

      <!-- Project Explorer -->
      <button class="quick-btn" title="Project Files">
        <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor">
          <path d="M1.5 3A1.5 1.5 0 000 4.5v7A1.5 1.5 0 001.5 13h13a1.5 1.5 0 001.5-1.5v-5A1.5 1.5 0 0014.5 5H7.707l-1.854-1.854A.5.5 0 005.5 3h-4z"/>
        </svg>
      </button>

      <!-- More -->
      <button class="quick-btn" title="More IDE Actions">
        <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor">
          <path d="M3 8a1.5 1.5 0 11-3 0 1.5 1.5 0 013 0zm6.5 0a1.5 1.5 0 11-3 0 1.5 1.5 0 013 0zm6.5 0a1.5 1.5 0 11-3 0 1.5 1.5 0 013 0z"/>
        </svg>
      </button>
    </div>
  </div>

  <!-- Right: Utilities (AI Assistant, Search, Settings) -->
  <div class="right-section">
    <!-- AI Assistant -->
    <button
      class="util-btn ai-btn"
      class:active={isAiOpen}
      title="AI Query Assistant / Explainer"
      onclick={() => isAiOpen = !isAiOpen}
    >
      <svg width="15" height="15" viewBox="0 0 16 16" fill="currentColor">
        <path d="M8 0a8 8 0 100 16A8 8 0 008 0zm0 1.5a6.5 6.5 0 015.63 9.75l-1.33-.77A5 5 0 104.5 12.19V13.8A6.5 6.5 0 018 1.5zm0 4.5a2 2 0 110 4 2 2 0 010-4z"/>
      </svg>
    </button>

    <!-- Global Search -->
    <button
      class="util-btn search-btn"
      title="Search Everywhere (Double Shift / Cmd+K)"
      onclick={() => appState.searchModalOpen = true}
    >
      <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor">
        <path fill-rule="evenodd" d="M11.5 7a4.5 4.5 0 11-9 0 4.5 4.5 0 019 0zm-.82 4.74a6 6 0 111.06-1.06l3.04 3.04a.75.75 0 11-1.06 1.06l-3.04-3.04z"/>
      </svg>
      <span class="search-hint">Double ⇧</span>
    </button>

    <!-- Settings Gear -->
    <button
      class="util-btn"
      title="Settings & Preferences (Cmd+,)"
      onclick={() => appState.settingsModalOpen = true}
    >
      <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor">
        <path d="M7.07 1a1 1 0 00-.97.757l-.24 1.026a5.526 5.526 0 00-1.282.74L3.6 3.057a1 1 0 00-1.205.247l-.98 1.135a1 1 0 00-.173 1.218l.617.848a5.534 5.534 0 00-.012 1.48l-.618.847a1 1 0 00.173 1.218l.98 1.136a1 1 0 001.206.246l.978-.466c.394.3.826.55 1.282.74l.24 1.026A1 1 0 007.07 15h1.86a1 1 0 00.97-.757l.24-1.026c.456-.19.888-.44 1.282-.74l.978.466a1 1 0 001.206-.246l.98-1.136a1 1 0 00-.173-1.218l-.617-.847c.105-.486.105-.993 0-1.48l.617-.848a1 1 0 00.173-1.218l-.98-1.135a1 1 0 00-1.206-.247l-.978.466a5.527 5.527 0 00-1.282-.74l-.24-1.026A1 1 0 008.93 1H7.07zm.93 5a2 2 0 110 4 2 2 0 010-4z"/>
      </svg>
    </button>
  </div>
</header>

<style>
  .titlebar {
    height: var(--titlebar-height);
    background-color: var(--bg-toolbar);
    border-bottom: 1px solid var(--border-default);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 12px;
    z-index: 50;
    flex-shrink: 0;
  }

  .left-section {
    display: flex;
    align-items: center;
    gap: 14px;
  }

  /* Traffic lights */
  .traffic-lights {
    display: flex;
    align-items: center;
    gap: 8px;
    padding-right: 4px;
  }

  .light {
    width: 12px;
    height: 12px;
    border-radius: 50%;
    display: inline-block;
    cursor: pointer;
    transition: transform 0.1s ease;
  }

  .light.close {
    background-color: var(--window-close);
  }
  .light.close:hover {
    background-color: #E0443E;
  }

  .light.minimize {
    background-color: var(--window-minimize);
  }
  .light.minimize:hover {
    background-color: #DEA123;
  }

  .light.maximize {
    background-color: var(--window-maximize);
  }
  .light.maximize:hover {
    background-color: #1AAB29;
  }

  /* Context selectors */
  .context-group {
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .avatar {
    width: 20px;
    height: 20px;
    border-radius: 50%;
    background-color: var(--profile-bg);
    color: var(--profile-fg);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 10px;
    font-weight: 700;
    margin-right: 2px;
  }

  .selector-btn {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 2px 6px;
    border-radius: 4px;
    color: var(--text-primary);
    font-size: 12px;
    font-weight: 500;
  }

  .selector-btn:hover {
    background-color: var(--bg-hover);
  }

  .branch-btn {
    color: var(--text-secondary);
    font-weight: 400;
  }

  .branch-icon {
    color: var(--text-muted);
  }

  .chevron {
    color: var(--text-muted);
  }

  .divider {
    color: var(--border-default);
    font-size: 11px;
  }

  /* Center quick actions */
  .center-section {
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .quick-actions {
    display: flex;
    align-items: center;
    gap: 4px;
    background-color: #222427;
    padding: 2px 4px;
    border-radius: 4px;
    border: 1px solid var(--border-subtle);
  }

  .quick-btn {
    width: 26px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 3px;
    color: var(--text-muted);
    transition: all 0.15s ease;
  }

  .quick-btn:hover {
    background-color: var(--bg-hover);
    color: var(--text-primary);
  }

  .quick-btn.run-btn {
    color: var(--action-success);
  }

  .quick-btn.run-btn:hover {
    color: var(--action-success-hover);
    background-color: rgba(87, 211, 140, 0.12);
  }

  /* Right utility section */
  .right-section {
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .util-btn {
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 6px;
    padding: 0 6px;
    border-radius: 4px;
    color: var(--text-muted);
    transition: all 0.15s ease;
  }

  .util-btn:hover {
    background-color: var(--bg-hover);
    color: var(--text-primary);
  }

  .ai-btn:hover, .ai-btn.active {
    color: #9DA0A8;
  }

  .search-btn {
    background-color: #25272A;
    border: 1px solid var(--border-default);
    padding: 0 8px;
    height: 24px;
    border-radius: 4px;
  }

  .search-hint {
    font-size: 10px;
    color: var(--text-muted);
    background-color: #313438;
    padding: 1px 4px;
    border-radius: 2px;
  }
</style>
