<script lang="ts">
  import { appState } from '../../state/appState.svelte';
  import { Icon } from '@ohjanus/icons';

  let isDraggingSplitter = $state(false);
  let startY = 0;
  let startHeight = 0;

  function handleSplitterMouseDown(e: MouseEvent) {
    isDraggingSplitter = true;
    startY = e.clientY;
    startHeight = appState.servicesHeight;

    const onMouseMove = (ev: MouseEvent) => {
      if (!isDraggingSplitter) return;
      const delta = startY - ev.clientY;
      const newHeight = Math.max(120, Math.min(450, startHeight + delta));
      appState.servicesHeight = newHeight;
    };

    const onMouseUp = () => {
      isDraggingSplitter = false;
      window.removeEventListener('mousemove', onMouseMove);
      window.removeEventListener('mouseup', onMouseUp);
    };

    window.addEventListener('mousemove', onMouseMove);
    window.addEventListener('mouseup', onMouseUp);
  }

  function toggleTree(id: string) {
    appState.treeExpanded[id] = !appState.treeExpanded[id];
  }

  function selectTable(tableName: string) {
    appState.selectedTreeNode = tableName;
    if (tableName === 'connection_credential') {
      appState.openTab({
        id: 'connection_credential',
        title: 'connectio...credential [[Dev][ReadOnly] 10.220.6.4]',
        type: 'table',
        closable: true,
        icon: 'table',
        env: 'Dev'
      });
    } else if (tableName === 'transfer_job' || tableName === 'console_2') {
      appState.openTab({
        id: 'console_2',
        title: 'console_2 [[PRD] 10.250.6.23]',
        type: 'console',
        closable: true,
        icon: 'lightning',
        env: 'PRD'
      });
    }
  }
</script>

<aside
  class="sidebar"
  class:collapsed={appState.isSidebarCollapsed}
  style="width: {appState.isSidebarCollapsed ? '0px' : appState.sidebarWidth + 'px'};"
>
  <!-- TOP PANE: Database Explorer -->
  <div class="explorer-pane">
    <!-- Top Window-like Explorer Header -->
    <div class="explorer-header-top">
      <span class="explorer-title">Database Explorer</span>
      <div class="header-window-icons">
        <button type="button" class="jb-icon-btn" title="Locate in Tree"><Icon name="crosshairs" size={12} /></button>
        <button type="button" class="jb-icon-btn" title="Expand / Collapse"><Icon name="expand-y" size={12} /></button>
        <button type="button" class="jb-icon-btn" title="Minimize"><Icon name="minus" size={12} /></button>
        <button type="button" class="jb-icon-btn" title="More Options">
          <Icon name="more" size={12} />
        </button>
      </div>
    </div>

    <!-- Explorer Sub-toolbar Row -->
    <div class="explorer-subtoolbar">
      <!-- + New -->
      <button type="button" class="jb-icon-btn" title="New Data Source (+)">
        <Icon name="plus" size={12} />
      </button>

      <!-- DB Gear Properties -->
      <button type="button" class="jb-icon-btn" title="Data Source Properties">
        <Icon name="database" size={12} />
      </button>

      <!-- Refresh (⟳) -->
      <button type="button" class="jb-icon-btn" title="Refresh (Cmd+F5)">
        <Icon name="refresh" size={12} />
      </button>

      <!-- Properties [] -->
      <button type="button" class="jb-icon-btn" title="Options">
        <Icon name="settings" size={11} />
      </button>

      <!-- DDL -->
      <button type="button" class="jb-icon-btn ddl-btn" title="Generate Schema DDL" onclick={() => appState.ddlModalOpen = true}>
        <span style="font-size: 10px; font-weight: 700; font-family: var(--font-code); color: #7A7E85;">DDL</span>
      </button>

      <!-- Back -->
      <button type="button" class="jb-icon-btn" title="Navigate Backward">
        <Icon name="chevron-left" size={11} />
      </button>

      <!-- Eye -->
      <button type="button" class="jb-icon-btn" title="Show / Hide Schemas">
        <Icon name="eye" size={12} />
      </button>
    </div>

    <!-- Tree View -->
    <div class="tree-viewport">
      <!-- Node: Dev Server [Dev][ReadOnly] 10.220.6.4 -->
      <button type="button" class="tree-node depth-0" onclick={() => toggleTree('dev_srv')}>
        <span class="chevron" class:expanded={appState.treeExpanded['dev_srv']}><Icon name="chevron-right" size={12} /></span>
        <!-- Elephant/DB Icon in cyan/blue -->
        <Icon name="database" size={14} color="#3B82F6" class="node-icon" />
        <span class="node-label"><strong>[Dev][ReadOnly]</strong> 10.220.6.4</span>
      </button>

      {#if appState.treeExpanded['dev_srv']}
        <!-- Schema: es -->
        <button type="button" class="tree-node depth-1" onclick={() => toggleTree('dev_schema_es')}>
          <span class="chevron" class:expanded={appState.treeExpanded['dev_schema_es']}><Icon name="chevron-right" size={12} /></span>
          <Icon name="folder" size={13} color="#C29D38" class="node-icon" />
          <span class="node-label">es</span>
        </button>

        <!-- Schema: information_schema -->
        <button type="button" class="tree-node depth-1">
          <span class="chevron"><Icon name="chevron-right" size={12} /></span>
          <Icon name="folder" size={13} color="#C29D38" class="node-icon" />
          <span class="node-label">information_schema</span>
        </button>

        <!-- Schema: marts -->
        <button type="button" class="tree-node depth-1">
          <span class="chevron"><Icon name="chevron-right" size={12} /></span>
          <Icon name="folder" size={13} color="#C29D38" class="node-icon" />
          <span class="node-label">marts</span>
        </button>

        <!-- Schema: pg_catalog -->
        <button type="button" class="tree-node depth-1">
          <span class="chevron"><Icon name="chevron-right" size={12} /></span>
          <Icon name="folder" size={13} color="#C29D38" class="node-icon" />
          <span class="node-label">pg_catalog</span>
        </button>

        <!-- Schema: pm -->
        <button type="button" class="tree-node depth-1">
          <span class="chevron"><Icon name="chevron-right" size={12} /></span>
          <Icon name="folder" size={13} color="#C29D38" class="node-icon" />
          <span class="node-label">pm</span>
        </button>

        <!-- Schema: public (Expanded) -->
        <button type="button" class="tree-node depth-1" onclick={() => toggleTree('dev_schema_public')}>
          <span class="chevron" class:expanded={appState.treeExpanded['dev_schema_public']}><Icon name="chevron-right" size={12} /></span>
          <Icon name="folder" size={13} color="#C29D38" class="node-icon" />
          <span class="node-label">public</span>
        </button>

        {#if appState.treeExpanded['dev_schema_public']}
          <!-- tables 33 -->
          <button type="button" class="tree-node depth-2" onclick={() => toggleTree('dev_tables')}>
            <span class="chevron" class:expanded={appState.treeExpanded['dev_tables']}><Icon name="chevron-right" size={12} /></span>
            <Icon name="table" size={13} color="#3B82F6" class="node-icon" />
            <span class="node-label">tables <span style="color: var(--text-muted); font-size: 11px;">33</span></span>
          </button>

          {#if appState.treeExpanded['dev_tables']}
            <!-- Table: category -->
            <button
              type="button"
              class="tree-node depth-3"
              class:selected={appState.selectedTreeNode === 'category'}
              onclick={() => selectTable('category')}
            >
              <span class="chevron"><Icon name="chevron-right" size={12} /></span>
              <Icon name="table" size={12} color="#4A88C7" class="node-icon" />
              <span class="node-label">category</span>
            </button>

            <!-- Table: connection_credential -->
            <button
              type="button"
              class="tree-node depth-3 selected"
              class:selected={appState.selectedTreeNode === 'connection_credential'}
              onclick={() => selectTable('connection_credential')}
            >
              <span class="chevron"><Icon name="chevron-right" size={12} /></span>
              <Icon name="table" size={12} color="#4A88C7" class="node-icon" />
              <span class="node-label">connection_credential</span>
            </button>

            <!-- Table: content -->
            <button
              type="button"
              class="tree-node depth-3"
              class:selected={appState.selectedTreeNode === 'content'}
              onclick={() => selectTable('content')}
            >
              <span class="chevron"><Icon name="chevron-right" size={12} /></span>
              <Icon name="table" size={12} color="#4A88C7" class="node-icon" />
              <span class="node-label">content</span>
            </button>

            <!-- Table: content_distribution -->
            <button
              type="button"
              class="tree-node depth-3"
              class:selected={appState.selectedTreeNode === 'content_distribution'}
              onclick={() => selectTable('content_distribution')}
            >
              <span class="chevron"><Icon name="chevron-right" size={12} /></span>
              <Icon name="table" size={12} color="#4A88C7" class="node-icon" />
              <span class="node-label">content_distribution</span>
            </button>

            <!-- Table: content_provider_configuration -->
            <button
              type="button"
              class="tree-node depth-3"
              class:selected={appState.selectedTreeNode === 'content_provider_configuration'}
              onclick={() => selectTable('content_provider_configuration')}
            >
              <span class="chevron"><Icon name="chevron-right" size={12} /></span>
              <Icon name="table" size={12} color="#4A88C7" class="node-icon" />
              <span class="node-label">content_provider_configuration</span>
            </button>
          {/if}
        {/if}
      {/if}
    </div>
  </div>

  <!-- HORIZONTAL SPLITTER -->
  <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
  <div
    class="horizontal-splitter"
    role="separator"
    tabindex="-1"
    aria-orientation="horizontal"
    onmousedown={handleSplitterMouseDown}
    title="Drag to resize Database Explorer / Services"
  ></div>

  <!-- BOTTOM PANE: Services Panel (Execution & Session Monitor) -->
  <div class="services-pane" style="height: {appState.servicesHeight}px;">
    <!-- Services Header -->
    <div class="services-header">
      <span class="services-title">Services</span>
      <div class="services-actions">
        <button type="button" class="jb-icon-btn" title="Transaction Mode: Auto">
          <span style="font-size: 11px; font-weight: 500; color: #7A7E85;">Tx,</span>
        </button>
        <button type="button" class="jb-icon-btn" title="New Session">
          <Icon name="plus" size={12} />
        </button>
        <button type="button" class="jb-icon-btn" title="View Options">
          <Icon name="eye" size={12} />
        </button>
        <button type="button" class="jb-icon-btn" title="Pin / Duplicate">
          <Icon name="pin" size={12} />
        </button>
        <button type="button" class="jb-icon-btn" title="Expand / Collapse"><Icon name="expand-y" size={12} /></button>
        <button type="button" class="jb-icon-btn" title="Close">
          <Icon name="close" size={11} />
        </button>
      </div>
    </div>

    <!-- Services Tree matching design2.png -->
    <div class="services-viewport">
      <!-- Database parent group -->
      <div class="service-item parent-item">
        <input type="checkbox" checked class="service-chk" />
        <span class="chevron expanded"><Icon name="chevron-right" size={12} /></span>
        <Icon name="folder" size={13} color="#C29D38" class="node-icon" />
        <span class="node-label">Database</span>
      </div>

      <!-- Dev Server Group -->
      <div class="service-item depth-1">
        <span class="chevron expanded"><Icon name="chevron-right" size={12} /></span>
        <Icon name="database" size={13} color="#3B82F6" class="node-icon" />
        <span class="node-label"><strong>[Dev][ReadOnly]</strong> 10.220.6.4</span>
      </div>

      <!-- Session 1: connection_credential (Selected highlight) -->
      <button
        type="button"
        class="service-item depth-2 selected"
        onclick={() => selectTable('connection_credential')}
      >
        <Icon name="table" size={12} color="#4A88C7" class="node-icon" />
        <span class="node-label truncate">connection_credential</span>
        <span class="latency-text">1 s 159 ms</span>
      </button>

      <!-- Session 2: events -->
      <div class="service-item depth-2">
        <Icon name="table" size={12} color="#4A88C7" class="node-icon" />
        <span class="node-label truncate">events</span>
        <span class="latency-text">986 ms</span>
      </div>

      <!-- Session 3: commands -->
      <div class="service-item depth-2">
        <Icon name="table" size={12} color="#4A88C7" class="node-icon" />
        <span class="node-label truncate">commands</span>
        <span class="latency-text">1 s 449 ms</span>
      </div>

      <!-- PRD Server Group -->
      <div class="service-item depth-1">
        <span class="chevron expanded"><Icon name="chevron-right" size={12} /></span>
        <Icon name="database" size={13} color="#3B82F6" class="node-icon" />
        <span class="node-label"><strong>[PRD]</strong> 10.250.6.23</span>
      </div>

      <!-- Session 4: console -->
      <div class="service-item depth-2">
        <Icon name="lightning" size={12} color="#3B82F6" class="node-icon" />
        <span class="node-label truncate">console</span>
      </div>

      <!-- Session 5: console_2 -->
      <button
        type="button"
        class="service-item depth-2"
        onclick={() => selectTable('console_2')}
      >
        <Icon name="lightning" size={12} color="#3B82F6" class="node-icon" />
        <span class="node-label truncate">console_2</span>
        <span class="latency-text">630 ms</span>
      </button>
    </div>
  </div>
</aside>

<style>
  .sidebar {
    background-color: transparent;
    border: none;
    border-radius: 0;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    flex-shrink: 0;
    transition: width 0.15s ease-out;
    height: 100%;
    gap: 0;
  }

  .sidebar.collapsed {
    border: none;
    width: 0 !important;
  }

  /* Header Top */
  .explorer-header-top {
    height: var(--toolbar-height, 32px);
    padding: 0 8px 0 12px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-shrink: 0;
  }

  .explorer-title {
    font-size: var(--font-size-sm, 12px);
    font-weight: 600;
    color: var(--text-primary);
  }

  .header-window-icons {
    display: flex;
    align-items: center;
    gap: 2px;
  }

  /* Explorer Subtoolbar */
  .explorer-subtoolbar {
    height: var(--toolbar-height, 32px);
    padding: 0 8px;
    display: flex;
    align-items: center;
    gap: 4px;
    border-bottom: 1px solid var(--border-subtle);
    flex-shrink: 0;
  }

  /* Tree Viewport */
  .explorer-pane {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 120px;
    overflow: hidden;
    background-color: var(--bg-card, #191A1C);
    border: none;
    border-radius: 8px;
  }

  .tree-viewport, .services-viewport {
    flex: 1;
    overflow-y: auto;
    overflow-x: hidden;
    padding: 4px 0;
  }

  .tree-node, .service-item {
    height: var(--tree-row-height, 28px);
    display: flex;
    align-items: center;
    gap: 4px;
    padding-right: 8px;
    cursor: pointer;
    font-size: var(--font-size-base, 13px);
    color: var(--text-primary);
    border-radius: var(--radius-sm, 4px);
    margin: 1px 4px;
    transition: background-color 0.1s ease;
    text-align: left;
    width: calc(100% - 8px);
  }

  .tree-node:hover, .service-item:hover {
    background-color: var(--bg-hover);
  }

  /* Exact selection highlight from design2.png */
  .tree-node.selected, .service-item.selected {
    background-color: var(--bg-selected);
    color: #FFFFFF;
  }

  .depth-0 { padding-left: 6px; }
  .depth-1 { padding-left: 20px; }
  .depth-2 { padding-left: 34px; }
  .depth-3 { padding-left: 48px; }

  .chevron {
    width: 16px;
    font-size: 13px;
    color: var(--text-muted);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    transition: transform 0.1s ease;
    transform: rotate(0deg);
    flex-shrink: 0;
  }

  .chevron.expanded {
    transform: rotate(90deg);
  }

  :global(.node-icon) {
    margin-right: 2px;
    width: 15px;
    height: 15px;
    flex-shrink: 0;
  }

  .node-label {
    flex: 1;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    font-size: var(--font-size-base, 13px);
  }

  /* Horizontal Splitter */
  .horizontal-splitter {
    height: 6px;
    background-color: transparent;
    cursor: row-resize;
    flex-shrink: 0;
    z-index: 10;
  }

  .horizontal-splitter:hover {
    background-color: var(--border-accent, #3574F0);
  }

  /* Services Panel */
  .services-pane {
    display: flex;
    flex-direction: column;
    background-color: var(--bg-card, #191A1C);
    border: none;
    border-radius: 8px;
    overflow: hidden;
    flex-shrink: 0;
  }

  .services-header {
    height: var(--toolbar-height, 32px);
    padding: 0 8px 0 12px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-shrink: 0;
  }

  .services-title {
    font-size: 12px;
    font-weight: 600;
    color: var(--text-primary);
  }

  .services-actions {
    display: flex;
    align-items: center;
    gap: 2px;
  }

  .service-chk {
    margin-right: 0;
    accent-color: var(--action-primary);
  }

  .latency-text {
    font-size: 11px;
    color: var(--text-muted);
    font-family: var(--font-code);
    margin-left: 4px;
  }
</style>
