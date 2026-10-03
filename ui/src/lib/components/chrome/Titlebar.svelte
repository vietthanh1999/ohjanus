<script lang="ts">
  import { appState } from '../../state/appState.svelte';
  import { Flex } from '@ohjanus/ui';

  let isAiOpen = $state(false);
</script>

<header class="titlebar">
  <!-- Left: Traffic Lights & Context Selectors -->
  <Flex class="left-section" align="center" gap="14px">
    <Flex class="traffic-lights" align="center" gap="7px">
      <span class="light close" title="Close"></span>
      <span class="light minimize" title="Minimize"></span>
      <span class="light maximize" title="Maximize"></span>
    </Flex>

    <Flex class="context-group" align="center" gap="6px">
      <!-- Profile Avatar -->
      <Flex class="avatar" align="center" justify="center" title="Viet Thanh">
        VT
      </Flex>

      <!-- Project Selector -->
      <button class="selector-btn" title="Current Workspace Project">
        <span class="project-name">VTVprime</span>
        <svg class="chevron" width="8" height="8" viewBox="0 0 16 16" fill="currentColor">
          <path d="M4 6l4 4 4-4H4z" />
        </svg>
      </button>

      <!-- Version Control Branch -->
      <button class="selector-btn branch-btn" title="Git Version Control">
        <span>Version Control</span>
        <svg class="chevron" width="8" height="8" viewBox="0 0 16 16" fill="currentColor">
          <path d="M4 6l4 4 4-4H4z" />
        </svg>
      </button>
    </Flex>
  </Flex>

  <!-- Center: Quick Actions Bar -->
  <Flex class="center-section" align="center" justify="center">
    <Flex class="quick-actions" align="center" gap="4px">
      <!-- Database cylinder icon -->
      <button
        class="quick-btn"
        title="Database Explorer"
        onclick={() => appState.isSidebarCollapsed = !appState.isSidebarCollapsed}
      >
        <svg width="15" height="15" viewBox="0 0 16 16" fill="currentColor">
          <path d="M8 1c3.866 0 7 1.12 7 2.5v9c0 1.38-3.134 2.5-7 2.5s-7-1.12-7-2.5v-9C1 2.12 4.134 1 8 1zm5.5 2.5c0-.44-2.126-1.2-5.5-1.2s-5.5.76-5.5 1.2 2.126 1.2 5.5 1.2 5.5-.76 5.5-1.2zM2.5 6.275V7.5c0 .44 2.126 1.2 5.5 1.2s5.5-.76 5.5-1.2V6.275c-1.39.638-3.376.925-5.5.925s-4.11-.287-5.5-.925zm0 3.5V11c0 .44 2.126 1.2 5.5 1.2s5.5-.76 5.5-1.2V9.775c-1.39.638-3.376.925-5.5.925s-4.11-.287-5.5-.925z"/>
        </svg>
      </button>

      <!-- Run circle -->
      <button
        class="quick-btn run-circle-btn"
        title="Run configurations"
        onclick={() => appState.executeQuery()}
      >
        <svg width="15" height="15" viewBox="0 0 16 16" fill="currentColor">
          <path d="M8 15A7 7 0 108 1a7 7 0 000 14zm-1.5-9.8v5.6L11 8l-4.5-2.8z"/>
        </svg>
      </button>

      <!-- Folder icon -->
      <button class="quick-btn" title="Project Files">
        <svg width="15" height="15" viewBox="0 0 16 16" fill="currentColor">
          <path d="M1.5 3A1.5 1.5 0 000 4.5v7A1.5 1.5 0 001.5 13h13a1.5 1.5 0 001.5-1.5v-5A1.5 1.5 0 0014.5 5H7.707l-1.854-1.854A.5.5 0 005.5 3h-4z"/>
        </svg>
      </button>

      <!-- More ··· -->
      <button class="quick-btn" title="More IDE Actions">
        <svg width="15" height="15" viewBox="0 0 16 16" fill="currentColor">
          <path d="M3 8a1.5 1.5 0 11-3 0 1.5 1.5 0 013 0zm6.5 0a1.5 1.5 0 11-3 0 1.5 1.5 0 013 0zm6.5 0a1.5 1.5 0 11-3 0 1.5 1.5 0 013 0z"/>
        </svg>
      </button>
    </Flex>
  </Flex>

  <!-- Right: Utilities (AI Assistant, Search, Settings) -->
  <Flex class="right-section" align="center" gap="4px">
    <!-- AI Assistant Spiral -->
    <button
      class="util-btn ai-btn"
      class:active={isAiOpen}
      title="AI Assistant (Double Shift)"
      onclick={() => isAiOpen = !isAiOpen}
    >
      <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
        <path d="M8 0a8 8 0 100 16A8 8 0 008 0zm0 1.5a6.5 6.5 0 015.63 9.75l-1.33-.77A5 5 0 104.5 12.19V13.8A6.5 6.5 0 018 1.5zm0 4.5a2 2 0 110 4 2 2 0 010-4z"/>
      </svg>
    </button>

    <!-- Global Search -->
    <button
      class="util-btn"
      title="Search Everywhere (Cmd+K / Double Shift)"
      onclick={() => appState.searchModalOpen = true}
    >
      <svg width="15" height="15" viewBox="0 0 16 16" fill="currentColor">
        <path fill-rule="evenodd" d="M11.5 7a4.5 4.5 0 11-9 0 4.5 4.5 0 019 0zm-.82 4.74a6 6 0 111.06-1.06l3.04 3.04a.75.75 0 11-1.06 1.06l-3.04-3.04z"/>
      </svg>
    </button>

    <!-- Settings Gear with Amber Dot -->
    <button
      class="util-btn settings-btn"
      title="Settings (Cmd+,)"
      onclick={() => appState.settingsModalOpen = true}
    >
      <svg width="15" height="15" viewBox="0 0 16 16" fill="currentColor">
        <path d="M7.07 1a1 1 0 00-.97.757l-.24 1.026a5.526 5.526 0 00-1.282.74L3.6 3.057a1 1 0 00-1.205.247l-.98 1.135a1 1 0 00-.173 1.218l.617.848a5.534 5.534 0 00-.012 1.48l-.618.847a1 1 0 00.173 1.218l.98 1.136a1 1 0 001.206.246l.978-.466c.394.3.826.55 1.282.74l.24 1.026A1 1 0 007.07 15h1.86a1 1 0 00.97-.757l.24-1.026c.456-.19.888-.44 1.282-.74l.978.466a1 1 0 001.206-.246l.98-1.136a1 1 0 00-.173-1.218l-.617-.847c.105-.486.105-.993 0-1.48l.617-.848a1 1 0 00.173-1.218l-.98-1.135a1 1 0 00-1.206-.247l-.978.466a5.527 5.527 0 00-1.282-.74l-.24-1.026A1 1 0 008.93 1H7.07zm.93 5a2 2 0 110 4 2 2 0 010-4z"/>
      </svg>
      <span class="badge-dot"></span>
    </button>
  </Flex>
</header>

<style>
  .titlebar {
    height: 38px;
    background: linear-gradient(90deg, #1C2426 0%, #202628 30%, #24272A 70%, #24272A 100%);
    border-bottom: none;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 10px;
    z-index: 50;
    flex-shrink: 0;
  }

  .light {
    width: 12px;
    height: 12px;
    border-radius: 50%;
    display: inline-block;
  }

  .light.close { background-color: var(--window-close); }
  .light.minimize { background-color: var(--window-minimize); }
  .light.maximize { background-color: var(--window-maximize); }

  /* Exact rectangular cyan avatar from design */
  :global(.avatar) {
    height: 18px;
    padding: 0 5px;
    border-radius: 3px;
    background-color: #0891B2;
    color: #FFFFFF;
    font-size: 11px;
    font-weight: 700;
    line-height: 1;
    letter-spacing: -0.2px;
  }

  .selector-btn {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 3px 6px;
    border-radius: 4px;
    color: var(--text-primary);
    font-size: 12px;
    font-weight: 400;
  }

  .selector-btn:hover {
    background-color: rgba(255, 255, 255, 0.08);
  }

  .branch-btn {
    color: var(--text-secondary);
  }

  .chevron {
    color: var(--text-muted);
  }

  .quick-btn {
    width: 26px;
    height: 26px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
    color: var(--text-muted);
    transition: all 0.1s ease;
  }

  .quick-btn:hover {
    background-color: rgba(255, 255, 255, 0.08);
    color: var(--text-primary);
  }

  .util-btn {
    width: 26px;
    height: 26px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
    color: var(--text-muted);
    position: relative;
    transition: all 0.1s ease;
  }

  .util-btn:hover {
    background-color: rgba(255, 255, 255, 0.08);
    color: var(--text-primary);
  }

  .settings-btn {
    position: relative;
  }

  .badge-dot {
    position: absolute;
    top: 4px;
    right: 4px;
    width: 5px;
    height: 5px;
    border-radius: 50%;
    background-color: #E5A122;
    box-shadow: 0 0 0 1px #24272A;
  }

  .ai-btn:hover, .ai-btn.active {
    color: #DFE1E5;
  }
</style>
