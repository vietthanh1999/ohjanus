<script lang="ts">
  import { appStore } from '../appStore.svelte';
  import type { DatabaseNode, ServicesItem } from '../types';

  let filterQuery = $state('');
  let isSearchActive = $state(false);
  let showHiddenSchemas = $state(false);
  let explorerHeightPercent = $state(60); // split ratio

  function toggleNode(node: DatabaseNode, e: MouseEvent) {
    e.stopPropagation();
    node.expanded = !node.expanded;
  }

  function handleNodeClick(node: DatabaseNode) {
    appStore.selectedTreeNodeId = node.id;
    if (node.type === 'table' && node.targetTable) {
      const conn = node.connection || '[Dev][ReadOnly] 10.220.6.4';
      appStore.openTableTab(node.targetTable, conn);
    }
  }

  function toggleService(item: ServicesItem, e: MouseEvent) {
    e.stopPropagation();
    item.expanded = !item.expanded;
  }

  function handleServiceClick(item: ServicesItem) {
    if (item.type === 'console') {
      appStore.openConsoleTab(item.name, '[PRD] 10.250.6.23');
    } else if (item.type === 'session') {
      appStore.openTableTab(item.name, '[Dev][ReadOnly] 10.220.6.4');
    }
  }
</script>

<aside class="sidebar">
  <!-- Top Half: Database Explorer -->
  <div class="explorer-section" style="height: {explorerHeightPercent}%;">
    <div class="panel-header">
      <span class="panel-title">Database Explorer</span>
      <div class="panel-actions">
        <button class="icon-btn" title="Add Data Source / Connection" onclick={() => appStore.openTab({ id: 'tab-connections', title: 'Connections', icon: 'connections', type: 'connections', closable: true })}>
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
        </button>
        <button class="icon-btn" title="Refresh Schemas">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M21.5 2v6h-6M21.34 15.57a10 10 0 1 1-.57-8.38l5.67-5.67"/></svg>
        </button>
        <button class="icon-btn" class:active={isSearchActive} title="Filter Tree Objects" onclick={() => isSearchActive = !isSearchActive}>
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
        </button>
        <button class="icon-btn" title="Generate DDL">
          <span style="font-size: 9px; font-weight: 700;">DDL</span>
        </button>
        <button class="icon-btn" title="Navigate Back">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="15 18 9 12 15 6"/></svg>
        </button>
        <button class="icon-btn" class:active={showHiddenSchemas} title="Toggle System Schemas &amp; Catalogs" onclick={() => showHiddenSchemas = !showHiddenSchemas}>
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>
        </button>
      </div>
    </div>

    {#if isSearchActive}
      <div class="filter-input-bar">
        <input 
          type="text" 
          placeholder="Filter tree (e.g. credential)..." 
          bind:value={filterQuery}
          autofocus
        />
        {#if filterQuery}
          <button class="clear-btn" onclick={() => filterQuery = ''}>✕</button>
        {/if}
      </div>
    {/if}

    <!-- Tree Structure -->
    <div class="tree-container">
      {#each appStore.treeData as serverNode}
        <!-- Server Node -->
        <div 
          class="tree-row depth-0" 
          class:selected={appStore.selectedTreeNodeId === serverNode.id}
          onclick={() => handleNodeClick(serverNode)}
        >
          <span class="caret" onclick={(e) => toggleNode(serverNode, e)}>
            {serverNode.expanded ? '▾' : '▸'}
          </span>
          <!-- Elephant Icon (PG) -->
          <span class="node-icon pg-icon">🐘</span>
          <span class="node-label">{serverNode.name}</span>
          {#if serverNode.badge}
            <span class="badge-pill">{serverNode.badge}</span>
          {/if}
        </div>

        {#if serverNode.expanded && serverNode.children}
          {#each serverNode.children as dbNode}
            <!-- Database Node -->
            <div 
              class="tree-row depth-1" 
              class:selected={appStore.selectedTreeNodeId === dbNode.id}
              onclick={() => handleNodeClick(dbNode)}
            >
              <span class="caret" onclick={(e) => toggleNode(dbNode, e)}>
                {dbNode.children ? (dbNode.expanded ? '▾' : '▸') : ' '}
              </span>
              <span class="node-icon folder-icon">📁</span>
              <span class="node-label">{dbNode.name}</span>
              {#if dbNode.badge}
                <span class="badge-pill">{dbNode.badge}</span>
              {/if}
            </div>

            {#if dbNode.expanded && dbNode.children}
              {#each dbNode.children as schNode}
                <!-- Schema Node -->
                <div 
                  class="tree-row depth-2" 
                  class:selected={appStore.selectedTreeNodeId === schNode.id}
                  onclick={() => handleNodeClick(schNode)}
                >
                  <span class="caret" onclick={(e) => toggleNode(schNode, e)}>
                    {schNode.children ? (schNode.expanded ? '▾' : '▸') : ' '}
                  </span>
                  <span class="node-icon schema-icon">⛁</span>
                  <span class="node-label">{schNode.name}</span>
                </div>

                {#if schNode.expanded && schNode.children}
                  {#each schNode.children as grpNode}
                    <!-- Group Node (tables / routines) -->
                    <div 
                      class="tree-row depth-3" 
                      class:selected={appStore.selectedTreeNodeId === grpNode.id}
                      onclick={() => handleNodeClick(grpNode)}
                    >
                      <span class="caret" onclick={(e) => toggleNode(grpNode, e)}>
                        {grpNode.children ? (grpNode.expanded ? '▾' : '▸') : ' '}
                      </span>
                      <span class="node-icon group-icon">🗄️</span>
                      <span class="node-label">{grpNode.name}</span>
                      {#if grpNode.badge}
                        <span class="badge-pill">{grpNode.badge}</span>
                      {/if}
                    </div>

                    {#if grpNode.expanded && grpNode.children}
                      {#each grpNode.children as tblNode}
                        <!-- Table Node -->
                        <div 
                          class="tree-row depth-4" 
                          class:selected={appStore.selectedTreeNodeId === tblNode.id}
                          ondblclick={() => handleNodeClick(tblNode)}
                          onclick={() => handleNodeClick(tblNode)}
                        >
                          <span class="caret invisible"> </span>
                          <span class="node-icon table-icon">📄</span>
                          <span class="node-label table-label">{tblNode.name}</span>
                        </div>
                      {/each}
                    {/if}
                  {/each}
                {/if}
              {/each}
            {/if}
          {/each}
        {/if}
      {/each}
    </div>
  </div>

  <!-- Horizontal Splitter -->
  <div class="horizontal-splitter" title="Drag to resize"></div>

  <!-- Bottom Half: Services Monitor -->
  <div class="services-section" style="height: {100 - explorerHeightPercent}%;">
    <div class="panel-header">
      <span class="panel-title">Services</span>
      <div class="panel-actions">
        <span class="tx-badge" title="Transaction Mode: Manual/Auto">Tx,</span>
        <button class="icon-btn" title="Add Service">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
        </button>
        <button class="icon-btn" title="Show Active Sessions">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>
        </button>
        <button class="icon-btn" title="Collapse All">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="5" y1="12" x2="19" y2="12"/></svg>
        </button>
      </div>
    </div>

    <div class="tree-container services-tree">
      {#each appStore.servicesData as rootItem}
        <div class="tree-row depth-0" onclick={(e) => toggleService(rootItem, e)}>
          <input type="checkbox" checked class="svc-chk" onclick={(e) => e.stopPropagation()} />
          <span class="caret">{rootItem.expanded ? '▾' : '▸'}</span>
          <span class="node-icon folder-icon">📁</span>
          <span class="node-label">{rootItem.name}</span>
        </div>

        {#if rootItem.expanded && rootItem.children}
          {#each rootItem.children as srvItem}
            <div class="tree-row depth-1" onclick={(e) => toggleService(srvItem, e)}>
              <span class="caret invisible"> </span>
              <span class="caret">{srvItem.expanded ? '▾' : '▸'}</span>
              <span class="node-icon pg-icon">🐘</span>
              <span class="node-label">{srvItem.name}</span>
            </div>

            {#if srvItem.expanded && srvItem.children}
              {#each srvItem.children as subItem}
                <div 
                  class="tree-row depth-2 service-item-row"
                  onclick={() => handleServiceClick(subItem)}
                >
                  <span class="caret invisible"> </span>
                  {#if subItem.type === 'console'}
                    <span class="node-icon console-icon">⚡</span>
                  {:else}
                    <span class="node-icon table-icon">📄</span>
                  {/if}
                  <span class="node-label">{subItem.name}</span>
                  {#if subItem.latency}
                    <span class="latency-pill">{subItem.latency}</span>
                  {/if}
                </div>
              {/each}
            {/if}
          {/each}
        {/if}
      {/each}
    </div>
  </div>
</aside>

<style>
  .sidebar {
    width: 280px;
    height: 100%;
    background: var(--surface-sidebar);
    border-right: 1px solid var(--border-default);
    display: flex;
    flex-direction: column;
    user-select: none;
    overflow: hidden;
  }

  .explorer-section, .services-section {
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .panel-header {
    height: 28px;
    padding: 0 8px 0 12px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    border-bottom: 1px solid var(--border-subtle);
    background: #25272A;
  }

  .panel-title {
    font-size: 11px;
    font-weight: 600;
    color: var(--text-primary);
  }

  .panel-actions {
    display: flex;
    align-items: center;
    gap: 2px;
  }

  .panel-actions .icon-btn.active {
    background: var(--surface-hover);
    color: var(--action-primary);
  }

  .tx-badge {
    font-size: 10px;
    font-weight: 700;
    color: var(--text-muted);
    padding: 1px 3px;
    cursor: pointer;
  }

  .filter-input-bar {
    padding: 4px 8px;
    background: #1E1F22;
    border-bottom: 1px solid var(--border-default);
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .filter-input-bar input {
    flex: 1;
    height: 20px;
    padding: 2px 4px;
    font-size: 11px;
    border: none;
    background: transparent;
  }

  .clear-btn {
    font-size: 10px;
    color: var(--text-muted);
  }

  .tree-container {
    flex: 1;
    overflow-y: auto;
    overflow-x: hidden;
    padding: 4px 0;
  }

  .tree-row {
    height: 22px;
    display: flex;
    align-items: center;
    padding-right: 8px;
    cursor: pointer;
    font-size: 12px;
    color: var(--text-primary);
    white-space: nowrap;
    border-radius: 4px;
    margin: 1px 4px;
  }

  .tree-row:hover {
    background: var(--surface-hover);
  }

  .tree-row.selected {
    background: var(--surface-selected);
    color: #fff;
  }

  .depth-0 { padding-left: 4px; }
  .depth-1 { padding-left: 18px; }
  .depth-2 { padding-left: 32px; }
  .depth-3 { padding-left: 46px; }
  .depth-4 { padding-left: 60px; }

  .caret {
    width: 14px;
    font-size: 10px;
    color: var(--text-muted);
    display: inline-flex;
    justify-content: center;
  }

  .caret.invisible {
    visibility: hidden;
  }

  .node-icon {
    margin-right: 6px;
    font-size: 11px;
    display: inline-flex;
    align-items: center;
  }

  .pg-icon { font-size: 13px; }
  .console-icon { color: #56A8F5; font-size: 12px; }
  .schema-icon { color: #3574F0; }
  .folder-icon { color: #EDA200; font-size: 12px; }

  .node-label {
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .horizontal-splitter {
    height: 4px;
    background: var(--border-default);
    cursor: row-resize;
    transition: background 0.15s;
  }

  .horizontal-splitter:hover {
    background: var(--border-accent);
  }

  .svc-chk {
    margin-right: 6px;
    accent-color: var(--action-primary);
  }

  .service-item-row .latency-pill {
    margin-left: auto;
    font-size: 10px;
    color: var(--badge-latency-fg);
  }
</style>
