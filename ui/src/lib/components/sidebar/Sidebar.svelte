<script lang="ts">
  import { appState } from '../../state/appState.svelte';

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
        <button type="button" class="jb-icon-btn" title="Locate in Tree">⌖</button>
        <button type="button" class="jb-icon-btn" title="Expand / Collapse">↕</button>
        <button type="button" class="jb-icon-btn" title="Minimize">—</button>
        <button type="button" class="jb-icon-btn" title="More Options">···</button>
      </div>
    </div>

    <!-- Explorer Sub-toolbar Row -->
    <div class="explorer-subtoolbar">
      <!-- + New -->
      <button type="button" class="jb-icon-btn" title="New Data Source (+)">
        <svg width="12" height="12" viewBox="0 0 16 16" fill="currentColor">
          <path d="M8 2a.75.75 0 01.75.75v4.5h4.5a.75.75 0 010 1.5h-4.5v4.5a.75.75 0 01-1.5 0v-4.5h-4.5a.75.75 0 010-1.5h4.5v-4.5A.75.75 0 018 2z"/>
        </svg>
      </button>

      <!-- DB Gear Properties -->
      <button type="button" class="jb-icon-btn" title="Data Source Properties">
        <svg width="12" height="12" viewBox="0 0 16 16" fill="currentColor">
          <path d="M8 1c3.866 0 7 1.12 7 2.5v9c0 1.38-3.134 2.5-7 2.5s-7-1.12-7-2.5v-9C1 2.12 4.134 1 8 1zm5.5 2.5c0-.44-2.126-1.2-5.5-1.2s-5.5.76-5.5 1.2 2.126 1.2 5.5 1.2 5.5-.76 5.5-1.2z"/>
        </svg>
      </button>

      <!-- Refresh (⟳) -->
      <button type="button" class="jb-icon-btn" title="Refresh (Cmd+F5)">
        <svg width="12" height="12" viewBox="0 0 16 16" fill="currentColor">
          <path fill-rule="evenodd" d="M8 3a5 5 0 104.546 2.914.5.5 0 01.908-.417A6 6 0 118 2v1z"/>
          <path d="M8 4.466V.534a.25.25 0 01.41-.192l2.36 1.966c.12.1.12.284 0 .384L8.41 4.658A.25.25 0 018 4.466z"/>
        </svg>
      </button>

      <!-- Properties [] -->
      <button type="button" class="jb-icon-btn" title="Options">
        <svg width="11" height="11" viewBox="0 0 16 16" fill="currentColor">
          <rect x="2" y="2" width="12" height="12" rx="1" fill="none" stroke="currentColor" stroke-width="1.5"/>
        </svg>
      </button>

      <!-- DDL -->
      <button type="button" class="jb-icon-btn ddl-btn" title="Generate Schema DDL" onclick={() => appState.ddlModalOpen = true}>
        <span style="font-size: 10px; font-weight: 700; font-family: var(--font-code); color: #7A7E85;">DDL</span>
      </button>

      <!-- Back -->
      <button type="button" class="jb-icon-btn" title="Navigate Backward">
        <svg width="11" height="11" viewBox="0 0 16 16" fill="currentColor">
          <path fill-rule="evenodd" d="M9.78 12.78a.75.75 0 01-1.06 0L4.47 8.53a.75.75 0 010-1.06l4.25-4.25a.75.75 0 011.06 1.06L6.06 8l3.72 3.72a.75.75 0 010 1.06z"/>
        </svg>
      </button>

      <!-- Eye -->
      <button type="button" class="jb-icon-btn" title="Show / Hide Schemas">
        <svg width="12" height="12" viewBox="0 0 16 16" fill="currentColor">
          <path d="M16 8s-3-5.5-8-5.5S0 8 0 8s3 5.5 8 5.5S16 8 16 8zM1.173 8a13.133 13.133 0 011.66-2.043C4.12 4.668 5.88 3.5 8 3.5c2.12 0 3.879 1.168 5.168 2.457A13.133 13.133 0 0114.828 8c-.058.087-.122.183-.195.288-.335.48-.83 1.12-1.465 1.755C11.879 11.332 10.119 12.5 8 12.5c-2.12 0-3.879-1.168-5.168-2.457A13.134 13.134 0 011.172 8z"/>
          <path d="M8 5.5a2.5 2.5 0 100 5 2.5 2.5 0 000-5zM7 8a1 1 0 112 0 1 1 0 01-2 0z"/>
        </svg>
      </button>
    </div>

    <!-- Tree View -->
    <div class="tree-viewport">
      <!-- Node: Dev Server [Dev][ReadOnly] 10.220.6.4 -->
      <button type="button" class="tree-node depth-0" onclick={() => toggleTree('dev_srv')}>
        <span class="chevron" class:expanded={appState.treeExpanded['dev_srv']}>›</span>
        <!-- Exact Elephant Icon in cyan/blue -->
        <svg class="node-icon icon-pg" width="14" height="14" viewBox="0 0 24 24" fill="#3B82F6">
          <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/>
        </svg>
        <span class="node-label"><strong>[Dev][ReadOnly]</strong> 10.220.6.4</span>
      </button>

      {#if appState.treeExpanded['dev_srv']}
        <!-- Schema: es -->
        <button type="button" class="tree-node depth-1" onclick={() => toggleTree('dev_schema_es')}>
          <span class="chevron" class:expanded={appState.treeExpanded['dev_schema_es']}>›</span>
          <svg class="node-icon" width="13" height="13" viewBox="0 0 16 16" fill="#C29D38">
            <path d="M1.5 3A1.5 1.5 0 000 4.5v7A1.5 1.5 0 001.5 13h13a1.5 1.5 0 001.5-1.5v-5A1.5 1.5 0 0014.5 5H7.707l-1.854-1.854A.5.5 0 005.5 3h-4z"/>
          </svg>
          <span class="node-label">es</span>
        </button>

        <!-- Schema: information_schema -->
        <button type="button" class="tree-node depth-1">
          <span class="chevron">›</span>
          <svg class="node-icon" width="13" height="13" viewBox="0 0 16 16" fill="#C29D38">
            <path d="M1.5 3A1.5 1.5 0 000 4.5v7A1.5 1.5 0 001.5 13h13a1.5 1.5 0 001.5-1.5v-5A1.5 1.5 0 0014.5 5H7.707l-1.854-1.854A.5.5 0 005.5 3h-4z"/>
          </svg>
          <span class="node-label">information_schema</span>
        </button>

        <!-- Schema: marts -->
        <button type="button" class="tree-node depth-1">
          <span class="chevron">›</span>
          <svg class="node-icon" width="13" height="13" viewBox="0 0 16 16" fill="#C29D38">
            <path d="M1.5 3A1.5 1.5 0 000 4.5v7A1.5 1.5 0 001.5 13h13a1.5 1.5 0 001.5-1.5v-5A1.5 1.5 0 0014.5 5H7.707l-1.854-1.854A.5.5 0 005.5 3h-4z"/>
          </svg>
          <span class="node-label">marts</span>
        </button>

        <!-- Schema: pg_catalog -->
        <button type="button" class="tree-node depth-1">
          <span class="chevron">›</span>
          <svg class="node-icon" width="13" height="13" viewBox="0 0 16 16" fill="#C29D38">
            <path d="M1.5 3A1.5 1.5 0 000 4.5v7A1.5 1.5 0 001.5 13h13a1.5 1.5 0 001.5-1.5v-5A1.5 1.5 0 0014.5 5H7.707l-1.854-1.854A.5.5 0 005.5 3h-4z"/>
          </svg>
          <span class="node-label">pg_catalog</span>
        </button>

        <!-- Schema: pm -->
        <button type="button" class="tree-node depth-1">
          <span class="chevron">›</span>
          <svg class="node-icon" width="13" height="13" viewBox="0 0 16 16" fill="#C29D38">
            <path d="M1.5 3A1.5 1.5 0 000 4.5v7A1.5 1.5 0 001.5 13h13a1.5 1.5 0 001.5-1.5v-5A1.5 1.5 0 0014.5 5H7.707l-1.854-1.854A.5.5 0 005.5 3h-4z"/>
          </svg>
          <span class="node-label">pm</span>
        </button>

        <!-- Schema: public (Expanded) -->
        <button type="button" class="tree-node depth-1" onclick={() => toggleTree('dev_schema_public')}>
          <span class="chevron" class:expanded={appState.treeExpanded['dev_schema_public']}>›</span>
          <svg class="node-icon" width="13" height="13" viewBox="0 0 16 16" fill="#C29D38">
            <path d="M1.5 3A1.5 1.5 0 000 4.5v7A1.5 1.5 0 001.5 13h13a1.5 1.5 0 001.5-1.5v-5A1.5 1.5 0 0014.5 5H7.707l-1.854-1.854A.5.5 0 005.5 3h-4z"/>
          </svg>
          <span class="node-label">public</span>
        </button>

        {#if appState.treeExpanded['dev_schema_public']}
          <!-- tables 33 (Exact matching design2.png: no pill badge, just text "33") -->
          <button type="button" class="tree-node depth-2" onclick={() => toggleTree('dev_tables')}>
            <span class="chevron" class:expanded={appState.treeExpanded['dev_tables']}>›</span>
            <svg class="node-icon" width="13" height="13" viewBox="0 0 16 16" fill="#3B82F6">
              <path d="M0 2a2 2 0 012-2h12a2 2 0 012 2v12a2 2 0 01-2 2H2a2 2 0 01-2-2V2zm1 2v2h14V4H1zm0 3v3h6V7H1zm7 0v3h7V7H8zm0 4v4h7v-4H8zm-1 4v-4H1v4h6z"/>
            </svg>
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
              <span class="chevron">›</span>
              <svg class="node-icon" width="12" height="12" viewBox="0 0 16 16" fill="#4A88C7">
                <path d="M0 2a1 1 0 011-1h14a1 1 0 011 1v12a1 1 0 01-1 1H1a1 1 0 01-1-1V2zm1 3v2h6V5H1zm7 0v2h7V5H8zm0 3v2h7V8H8zm-1 0H1v2h6V8zm0 3H1v2h6v-2zm1 0v2h7v-2H8z"/>
              </svg>
              <span class="node-label">category</span>
            </button>

            <!-- Table: connection_credential (Selected matching design2.png) -->
            <button
              type="button"
              class="tree-node depth-3 selected"
              class:selected={appState.selectedTreeNode === 'connection_credential'}
              onclick={() => selectTable('connection_credential')}
            >
              <span class="chevron">›</span>
              <svg class="node-icon" width="12" height="12" viewBox="0 0 16 16" fill="#4A88C7">
                <path d="M0 2a1 1 0 011-1h14a1 1 0 011 1v12a1 1 0 01-1 1H1a1 1 0 01-1-1V2zm1 3v2h6V5H1zm7 0v2h7V5H8zm0 3v2h7V8H8zm-1 0H1v2h6V8zm0 3H1v2h6v-2zm1 0v2h7v-2H8z"/>
              </svg>
              <span class="node-label">connection_credential</span>
            </button>

            <!-- Table: content -->
            <button
              type="button"
              class="tree-node depth-3"
              class:selected={appState.selectedTreeNode === 'content'}
              onclick={() => selectTable('content')}
            >
              <span class="chevron">›</span>
              <svg class="node-icon" width="12" height="12" viewBox="0 0 16 16" fill="#4A88C7">
                <path d="M0 2a1 1 0 011-1h14a1 1 0 011 1v12a1 1 0 01-1 1H1a1 1 0 01-1-1V2zm1 3v2h6V5H1zm7 0v2h7V5H8zm0 3v2h7V8H8zm-1 0H1v2h6V8zm0 3H1v2h6v-2zm1 0v2h7v-2H8z"/>
              </svg>
              <span class="node-label">content</span>
            </button>

            <!-- Table: content_distribution -->
            <button
              type="button"
              class="tree-node depth-3"
              class:selected={appState.selectedTreeNode === 'content_distribution'}
              onclick={() => selectTable('content_distribution')}
            >
              <span class="chevron">›</span>
              <svg class="node-icon" width="12" height="12" viewBox="0 0 16 16" fill="#4A88C7">
                <path d="M0 2a1 1 0 011-1h14a1 1 0 011 1v12a1 1 0 01-1 1H1a1 1 0 01-1-1V2zm1 3v2h6V5H1zm7 0v2h7V5H8zm0 3v2h7V8H8zm-1 0H1v2h6V8zm0 3H1v2h6v-2zm1 0v2h7v-2H8z"/>
              </svg>
              <span class="node-label">content_distribution</span>
            </button>

            <!-- Table: content_provider_configuration -->
            <button
              type="button"
              class="tree-node depth-3"
              class:selected={appState.selectedTreeNode === 'content_provider_configuration'}
              onclick={() => selectTable('content_provider_configuration')}
            >
              <span class="chevron">›</span>
              <svg class="node-icon" width="12" height="12" viewBox="0 0 16 16" fill="#4A88C7">
                <path d="M0 2a1 1 0 011-1h14a1 1 0 011 1v12a1 1 0 01-1 1H1a1 1 0 01-1-1V2zm1 3v2h6V5H1zm7 0v2h7V5H8zm0 3v2h7V8H8zm-1 0H1v2h6V8zm0 3H1v2h6v-2zm1 0v2h7v-2H8z"/>
              </svg>
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
        <button type="button" class="jb-icon-btn" title="New Session">+</button>
        <button type="button" class="jb-icon-btn" title="View Options">👁️</button>
        <button type="button" class="jb-icon-btn" title="Pin / Duplicate">⧉</button>
        <button type="button" class="jb-icon-btn" title="Expand / Collapse">↕</button>
        <button type="button" class="jb-icon-btn" title="Close">✕</button>
      </div>
    </div>

    <!-- Services Tree matching design2.png -->
    <div class="services-viewport">
      <!-- Database parent group -->
      <div class="service-item parent-item">
        <input type="checkbox" checked class="service-chk" />
        <span class="chevron expanded">›</span>
        <svg class="node-icon" width="13" height="13" viewBox="0 0 16 16" fill="#C29D38">
          <path d="M1.5 3A1.5 1.5 0 000 4.5v7A1.5 1.5 0 001.5 13h13a1.5 1.5 0 001.5-1.5v-5A1.5 1.5 0 0014.5 5H7.707l-1.854-1.854A.5.5 0 005.5 3h-4z"/>
        </svg>
        <span class="node-label">Database</span>
      </div>

      <!-- Dev Server Group -->
      <div class="service-item depth-1">
        <span class="chevron expanded">›</span>
        <svg class="node-icon" width="13" height="13" viewBox="0 0 24 24" fill="#3B82F6">
          <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/>
        </svg>
        <span class="node-label"><strong>[Dev][ReadOnly]</strong> 10.220.6.4</span>
      </div>

      <!-- Session 1: connection_credential (Selected highlight) -->
      <button
        type="button"
        class="service-item depth-2 selected"
        onclick={() => selectTable('connection_credential')}
      >
        <svg class="node-icon" width="12" height="12" viewBox="0 0 16 16" fill="#4A88C7">
          <path d="M0 2a1 1 0 011-1h14a1 1 0 011 1v12a1 1 0 01-1 1H1a1 1 0 01-1-1V2zm1 3v2h6V5H1zm7 0v2h7V5H8zm0 3v2h7V8H8zm-1 0H1v2h6V8zm0 3H1v2h6v-2zm1 0v2h7v-2H8z"/>
        </svg>
        <span class="node-label truncate">connection_credential</span>
        <span class="latency-text">1 s 159 ms</span>
      </button>

      <!-- Session 2: events -->
      <div class="service-item depth-2">
        <svg class="node-icon" width="12" height="12" viewBox="0 0 16 16" fill="#4A88C7">
          <path d="M0 2a1 1 0 011-1h14a1 1 0 011 1v12a1 1 0 01-1 1H1a1 1 0 01-1-1V2zm1 3v2h6V5H1zm7 0v2h7V5H8zm0 3v2h7V8H8zm-1 0H1v2h6V8zm0 3H1v2h6v-2zm1 0v2h7v-2H8z"/>
        </svg>
        <span class="node-label truncate">events</span>
        <span class="latency-text">986 ms</span>
      </div>

      <!-- Session 3: commands -->
      <div class="service-item depth-2">
        <svg class="node-icon" width="12" height="12" viewBox="0 0 16 16" fill="#4A88C7">
          <path d="M0 2a1 1 0 011-1h14a1 1 0 011 1v12a1 1 0 01-1 1H1a1 1 0 01-1-1V2zm1 3v2h6V5H1zm7 0v2h7V5H8zm0 3v2h7V8H8zm-1 0H1v2h6V8zm0 3H1v2h6v-2zm1 0v2h7v-2H8z"/>
        </svg>
        <span class="node-label truncate">commands</span>
        <span class="latency-text">1 s 449 ms</span>
      </div>

      <!-- PRD Server Group -->
      <div class="service-item depth-1">
        <span class="chevron expanded">›</span>
        <svg class="node-icon" width="13" height="13" viewBox="0 0 24 24" fill="#3B82F6">
          <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/>
        </svg>
        <span class="node-label"><strong>[PRD]</strong> 10.250.6.23</span>
      </div>

      <!-- Session 4: console -->
      <div class="service-item depth-2">
        <svg class="node-icon" width="12" height="12" viewBox="0 0 16 16" fill="#3B82F6">
          <path d="M11.251.068a.5.5 0 01.42.58L10.07 5H14a.5.5 0 01.372.832l-8.5 9.5a.5.5 0 01-.842-.512L6.63 9.5H2.5a.5.5 0 01-.42-.772l8.5-8.5a.5.5 0 01.671-.16z"/>
        </svg>
        <span class="node-label truncate">console</span>
      </div>

      <!-- Session 5: console_2 -->
      <button
        type="button"
        class="service-item depth-2"
        onclick={() => selectTable('console_2')}
      >
        <svg class="node-icon" width="12" height="12" viewBox="0 0 16 16" fill="#3B82F6">
          <path d="M11.251.068a.5.5 0 01.42.58L10.07 5H14a.5.5 0 01.372.832l-8.5 9.5a.5.5 0 01-.842-.512L6.63 9.5H2.5a.5.5 0 01-.42-.772l8.5-8.5a.5.5 0 01.671-.16z"/>
        </svg>
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
    background-color: var(--bg-canvas, #1E1F22);
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

  .node-icon {
    margin-right: 6px;
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
    background-color: var(--bg-canvas, #1E1F22);
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
    margin-right: 4px;
    accent-color: var(--action-primary);
  }

  .latency-text {
    font-size: 11px;
    color: var(--text-muted);
    font-family: var(--font-code);
    margin-left: 8px;
  }
</style>
