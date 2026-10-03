<script lang="ts">
  import { appState } from "../../state/appState.svelte";
  import { Flex } from "@ohjanus/ui";
  import { Icon } from "@ohjanus/icons";

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
        <span class="project-name">VThanh</span>
        <Icon name="chevron-down" size={10} />
      </button>

      <!-- Version Control Branch -->
      <button class="selector-btn branch-btn" title="Git Version Control">
        <span>Version Control</span>
        <Icon name="chevron-down" size={10} />
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
        onclick={() =>
          (appState.isSidebarCollapsed = !appState.isSidebarCollapsed)}
      >
        <Icon name="database" size={15} />
      </button>

      <!-- Run circle -->
      <button
        class="quick-btn run-circle-btn"
        title="Run configurations"
        onclick={() => appState.executeQuery()}
      >
        <Icon name="play" size={14} />
      </button>

      <!-- Folder icon -->
      <button class="quick-btn" title="Project Files">
        <Icon name="folder" size={15} />
      </button>

      <!-- More ··· -->
      <button class="quick-btn" title="More IDE Actions">
        <Icon name="more" size={15} />
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
      onclick={() => (isAiOpen = !isAiOpen)}
    >
      <Icon name="cpu" size={15} />
    </button>

    <!-- Global Search -->
    <button
      class="util-btn"
      title="Search Everywhere (Cmd+K / Double Shift)"
      onclick={() => (appState.searchModalOpen = true)}
    >
      <Icon name="search" size={15} />
    </button>

    <!-- Settings Gear with Amber Dot -->
    <button
      class="util-btn settings-btn"
      title="Settings (Cmd+,)"
      onclick={() => (appState.settingsModalOpen = true)}
    >
      <Icon name="settings" size={15} />
      <span class="badge-dot"></span>
    </button>
  </Flex>
</header>

<style>
  .titlebar {
    height: 38px;
    background: linear-gradient(
      90deg,
      #1c2426 0%,
      #202628 30%,
      #24272a 70%,
      #24272a 100%
    );
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

  .light.close {
    background-color: var(--window-close);
  }
  .light.minimize {
    background-color: var(--window-minimize);
  }
  .light.maximize {
    background-color: var(--window-maximize);
  }

  /* Exact rectangular cyan avatar from design */
  :global(.avatar) {
    height: 18px;
    padding: 0 5px;
    border-radius: 3px;
    background-color: #0891b2;
    color: #ffffff;
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

  .quick-btn {
    width: var(--icon-btn-size-md, 28px);
    height: var(--icon-btn-size-md, 28px);
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--radius-sm, 4px);
    color: var(--text-muted);
    transition: all 0.1s ease;
  }

  .quick-btn:hover {
    background-color: rgba(255, 255, 255, 0.08);
    color: var(--text-primary);
  }

  .util-btn {
    width: var(--icon-btn-size-md, 28px);
    height: var(--icon-btn-size-md, 28px);
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--radius-sm, 4px);
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
    background-color: #e5a122;
    box-shadow: 0 0 0 1px #24272a;
  }

  .ai-btn:hover,
  .ai-btn.active {
    color: #dfe1e5;
  }
</style>
