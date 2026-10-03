<script lang="ts">
  import { appStore } from '../../appStore.svelte';

  let activeResultTab = $state('prd_mh_asset.public.transfer_job 2');
  let selectedCell = $state({ row: 0, col: 'id' });
  let isExecuting = $state(false);

  function executeQuery() {
    isExecuting = true;
    appStore.runQuery();
    setTimeout(() => {
      isExecuting = false;
    }, 450);
  }

  function handleKeydown(e: KeyboardEvent) {
    if ((e.metaKey || e.ctrlKey) && e.key === 'Enter') {
      e.preventDefault();
      executeQuery();
    }
  }
</script>

<div class="console-view" onkeydown={handleKeydown} tabindex="0">
  <!-- Top SQL Sub-Toolbar -->
  <div class="sql-action-bar">
    <div class="left-actions">
      <!-- Run Buttons -->
      <button class="action-btn run-btn" class:running={isExecuting} title="Execute Entire Script (Cmd+Enter)" onclick={executeQuery}>
        <span class="icon">▶</span>
      </button>

      <button class="action-btn run-caret-btn" title="Execute Statement Under Caret" onclick={executeQuery}>
        <span class="icon">▶_</span>
      </button>

      <button class="icon-btn" title="Query History">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
      </button>

      <button class="icon-btn stop-btn" title="Cancel Execution">
        <span style="color: var(--action-danger); font-size: 11px;">⏹</span>
      </button>

      <button class="icon-btn" title="Database Session Parameters">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>
      </button>

      <button class="icon-btn" title="Toggle Split Orientation">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2"/><line x1="12" y1="3" x2="12" y2="21"/></svg>
      </button>

      <span class="toolbar-divider"></span>

      <!-- Tx Mode & Playground -->
      <div class="dropdown-pill">
        <span>Tx: Auto</span>
        <span class="chevron">▾</span>
      </div>

      <div class="dropdown-pill checkbox-pill">
        <input type="checkbox" id="chk-playground" />
        <label for="chk-playground">Playground</label>
        <span class="chevron">▾</span>
      </div>
    </div>

    <!-- Right Schema Target -->
    <div class="right-actions">
      <div class="schema-selector-btn" title="Target Active Schema">
        <span class="schema-icon">⛁</span>
        <span class="schema-name">prd_mh_asset.public</span>
        <span class="chevron">▾</span>
      </div>
    </div>
  </div>

  <!-- Upper Work Area: SQL Query Editor -->
  <div class="editor-pane">
    <div class="gutter">
      {#each [64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74] as lineNo}
        <div class="gutter-line" class:active-line={lineNo === 71} class:success-line={lineNo === 74}>
          <span class="line-number">{lineNo}</span>
          {#if lineNo === 74}
            <span class="exec-status">✓</span>
          {/if}
        </div>
      {/each}
    </div>

    <!-- Editor Code Canvas with Active Execution Block -->
    <div class="code-canvas font-mono">
      <!-- Active Execution Block Highlight -->
      <div class="execution-block animate-fade-in">
        <div class="floating-exec-pill">
          <span class="check">✓</span>
          <span class="count">6</span>
          <span class="arrows">▲ ▼</span>
        </div>

        <div class="code-line">
          <span class="keyword">join</span> content ct 
          <span class="code-lens">1..n &lt;-&gt; 1: on ct.id = ma."contentID"</span>
        </div>
        <div class="code-line">
          <span class="keyword">LEFT JOIN</span> transfer_job_queue q <span class="keyword">ON</span> q."jobID" = c.process_id
        </div>
        <div class="code-line">
          <span class="keyword">WHERE</span> c.participant = <span class="string">'media-transferer'</span>
        </div>
        <div class="code-line">
          <span class="keyword">AND</span> c.process_name   = <span class="string">'media-ingest'</span>
        </div>
        <div class="code-line">
          <span class="keyword">AND</span> c.name           = <span class="string">'create-ingest-job'</span>
        </div>
        <div class="code-line">
          <span class="keyword">AND</span> c.retry_count    = <span class="number">0</span>
        </div>
        <div class="code-line">
          <span class="keyword">AND</span> c.schedule <span class="keyword">IS NULL</span>
        </div>
        <div class="code-line active-cursor-line">
          <span class="keyword">AND</span> c."timestamp"    &lt; <span class="func">now()</span> - <span class="keyword">interval</span> <span class="string">'5 minutes'</span> <span class="comment">-- streamer poll 10s -&gt; &gt;5 phút <span class="spell-wavy">chắc</span> <span class="spell-wavy">chắn</span> <span class="spell-wavy">không</span> ai đọc</span>
        </div>
        <div class="code-line">
          <span class="keyword">ORDER BY</span> c.ordinal;
        </div>
      </div>

      <div class="code-line empty-line"></div>

      <!-- Second Statement (Executed) -->
      <div class="code-line executed-statement">
        <span class="keyword">select</span> * <span class="keyword">from</span> transfer_job <span class="keyword">where</span> id = <span class="string">'b8262180-1075-43f8-8338-2412d4734d65'</span>
      </div>
    </div>
  </div>

  <!-- Horizontal Splitter between Editor and Results -->
  <div class="results-splitter"></div>

  <!-- Lower Work Area: Multi-Results Grid -->
  <div class="results-pane">
    <!-- Result Sub-Tabs Bar -->
    <div class="result-tabs-bar">
      <div class="res-tab" class:active={activeResultTab === 'Result 1'} onclick={() => activeResultTab = 'Result 1'}>
        Result 1
      </div>
      <div class="res-tab" class:active={activeResultTab === 'Result 1-2'} onclick={() => activeResultTab = 'Result 1-2'}>
        Result 1-2
      </div>
      <div class="res-tab" class:active={activeResultTab === 'prd_mh_asset.public.transfer_job'} onclick={() => activeResultTab = 'prd_mh_asset.public.transfer_job'}>
        prd_mh_asset.public.transfer_job
      </div>
      <div class="res-tab" class:active={activeResultTab === 'Result 1-4'} onclick={() => activeResultTab = 'Result 1-4'}>
        Result 1-4
      </div>
      <div class="res-tab" class:active={activeResultTab === 'prd_mh_asset.public.transfer_job 2'} onclick={() => activeResultTab = 'prd_mh_asset.public.transfer_job 2'}>
        prd_mh_asset.public.transfer_job 2
        <span class="close-x">✕</span>
      </div>
      <button class="icon-btn tab-overflow-btn">▾</button>
    </div>

    <!-- Result Action Toolbar -->
    <div class="result-toolbar">
      <div class="res-toolbar-left">
        <button class="icon-btn active" title="Table View">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2"/><path d="M3 9h18M3 15h18M9 3v18M15 3v18"/></svg>
        </button>
        <button class="icon-btn" title="Text View">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
        </button>
        <button class="icon-btn" title="Refresh Query" onclick={executeQuery}>
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M21.5 2v6h-6M21.34 15.57a10 10 0 1 1-.57-8.38l5.67-5.67"/></svg>
        </button>
        <button class="icon-btn" title="History">🕒</button>
        <button class="icon-btn" title="Stop">⏹</button>
        <button class="icon-btn" title="Add Row">+</button>
        <button class="icon-btn" title="Delete Row">-</button>
        <button class="icon-btn" title="Undo">↩</button>
        <button class="icon-btn" title="Redo">↪</button>
        <button class="icon-btn" title="Up">↑</button>
        <button class="icon-btn" title="Down">↓</button>

        <span class="toolbar-divider"></span>
        <span class="toolbar-text">Tx: Auto</span>
        <span class="toolbar-text">DDL</span>
        <button class="icon-btn" title="Pin">📌</button>
        <button class="icon-btn" title="Find">🔍</button>
        <button class="icon-btn" title="Filter">Y</button>
        <button class="icon-btn" title="Statistics">📊</button>
      </div>

      <div class="res-toolbar-right">
        <span class="dropdown-pill">CSV ▾</span>
        <button class="icon-btn" title="Export">⤓</button>
        <button class="icon-btn" title="Import">⤒</button>
        <button class="icon-btn" title="Visualize">📈</button>
        <button class="icon-btn" title="Show/Hide Columns">👁️</button>
        <button class="icon-btn" title="Settings">⚙️</button>
      </div>
    </div>

    <!-- Data Grid Table -->
    <div class="grid-table-container">
      <table class="data-grid font-mono">
        <thead>
          <tr>
            <th class="row-num-col"></th>
            <th class="col-header">
              <span class="col-icon key-icon">🔑</span>
              <span class="col-name">id</span>
              <span class="filter-funnel">Y</span>
              <span class="sort-icon">⇅</span>
            </th>
            <th class="col-header">
              <span class="col-icon num-icon">#</span>
              <span class="col-name">ordinal</span>
              <span class="filter-funnel">Y</span>
              <span class="sort-icon">⇅</span>
            </th>
            <th class="col-header">
              <span class="col-icon str-icon">" "</span>
              <span class="col-name">connectionCredentialID</span>
              <span class="filter-funnel">Y</span>
              <span class="sort-icon">⇅</span>
            </th>
            <th class="col-header">
              <span class="col-icon str-icon">" "</span>
              <span class="col-name">sourcePath</span>
              <span class="filter-funnel">Y</span>
              <span class="sort-icon">⇅</span>
            </th>
          </tr>
        </thead>
        <tbody>
          {#each appStore.consoleResults as row, idx}
            <tr>
              <td class="row-index">{idx + 1}</td>
              <td 
                class="cell" 
                class:focused={selectedCell.row === idx && selectedCell.col === 'id'}
                onclick={() => selectedCell = { row: idx, col: 'id' }}
              >
                {row.id}
              </td>
              <td 
                class="cell" 
                class:focused={selectedCell.row === idx && selectedCell.col === 'ordinal'}
                onclick={() => selectedCell = { row: idx, col: 'ordinal' }}
              >
                {row.ordinal}
              </td>
              <td 
                class="cell" 
                class:focused={selectedCell.row === idx && selectedCell.col === 'connectionCredentialID'}
                onclick={() => selectedCell = { row: idx, col: 'connectionCredentialID' }}
              >
                {row.connectionCredentialID}
              </td>
              <td 
                class="cell" 
                class:focused={selectedCell.row === idx && selectedCell.col === 'sourcePath'}
                onclick={() => selectedCell = { row: idx, col: 'sourcePath' }}
              >
                {row.sourcePath}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>

      <!-- Floating Row Counter Badge -->
      <div class="floating-row-badge animate-fade-in">
        <span class="text">{appStore.consoleResults.length} row ▾</span>
        <span class="sep">|</span>
        <span class="kebab">⋮</span>
      </div>
    </div>
  </div>
</div>

<style>
  .console-view {
    display: flex;
    flex-direction: column;
    height: 100%;
    background: var(--surface-canvas);
    outline: none;
    overflow: hidden;
  }

  .sql-action-bar {
    height: var(--toolbar-height);
    background: var(--surface-toolbar);
    border-bottom: 1px solid var(--border-default);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 8px;
    gap: 6px;
  }

  .left-actions, .right-actions {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .run-btn {
    color: var(--action-success);
    font-size: 13px;
    padding: 2px 6px;
    border-radius: 3px;
  }

  .run-btn:hover {
    background: rgba(87, 211, 140, 0.15);
  }

  .run-btn.running {
    animation: pulse 0.6s infinite alternate;
  }

  @keyframes pulse {
    from { opacity: 0.5; }
    to { opacity: 1; }
  }

  .run-caret-btn {
    color: var(--action-success);
    font-weight: 700;
    font-size: 12px;
    padding: 2px 5px;
  }

  .toolbar-divider {
    width: 1px;
    height: 16px;
    background: var(--border-default);
    margin: 0 4px;
  }

  .dropdown-pill {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 11px;
    color: var(--text-secondary);
    padding: 2px 6px;
    border-radius: 3px;
    cursor: pointer;
  }

  .dropdown-pill:hover {
    background: var(--surface-hover);
    color: var(--text-primary);
  }

  .checkbox-pill {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .schema-selector-btn {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 2px 8px;
    background: #1E1F22;
    border: 1px solid var(--border-default);
    border-radius: 3px;
    font-size: 11px;
    color: var(--text-primary);
    cursor: pointer;
  }

  .schema-icon {
    color: #3574F0;
  }

  /* Editor Pane */
  .editor-pane {
    height: 48%;
    background: var(--surface-canvas);
    display: flex;
    overflow-y: auto;
    position: relative;
  }

  .gutter {
    width: 48px;
    background: var(--surface-canvas);
    border-right: 1px solid #2B2D30;
    padding: 8px 0;
    user-select: none;
  }

  .gutter-line {
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: flex-end;
    padding-right: 8px;
    font-family: var(--font-code);
    font-size: 12px;
    color: var(--text-muted);
    position: relative;
  }

  .gutter-line.active-line .line-number {
    color: var(--text-primary);
    font-weight: 600;
  }

  .exec-status {
    position: absolute;
    left: 6px;
    color: var(--action-success);
    font-size: 11px;
  }

  .code-canvas {
    flex: 1;
    padding: 8px 12px;
    font-size: 13px;
    line-height: 20px;
    color: var(--syntax-identifier);
    position: relative;
    overflow-x: auto;
  }

  .execution-block {
    background: var(--surface-execution-block);
    border: 1px solid var(--border-execution);
    border-radius: 2px;
    padding: 2px 6px;
    position: relative;
    margin-bottom: 8px;
  }

  .floating-exec-pill {
    position: absolute;
    top: 4px;
    right: 8px;
    background: #1E2B37;
    border: 1px solid #2B5B9E;
    border-radius: 3px;
    padding: 1px 6px;
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 11px;
    z-index: 10;
  }

  .floating-exec-pill .check { color: var(--action-success); }
  .floating-exec-pill .count { color: #DFE1E5; font-weight: 600; }
  .floating-exec-pill .arrows { color: var(--text-muted); font-size: 9px; }

  .code-line {
    height: 20px;
    white-space: pre;
  }

  .keyword { color: var(--syntax-keyword); font-weight: 500; }
  .func { color: var(--syntax-function); font-style: italic; }
  .string { color: var(--syntax-string); }
  .number { color: var(--syntax-number); }
  .comment { color: var(--syntax-comment); font-style: italic; }

  .code-lens {
    color: var(--syntax-code-lens);
    font-style: italic;
    font-size: 11px;
    margin-left: 12px;
    opacity: 0.8;
  }

  .spell-wavy {
    text-decoration: underline wavy var(--syntax-spell-warn);
  }

  .results-splitter {
    height: 4px;
    background: var(--border-default);
    cursor: row-resize;
  }

  /* Results Pane */
  .results-pane {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    background: var(--surface-canvas);
    position: relative;
  }

  .result-tabs-bar {
    height: 28px;
    background: var(--surface-toolbar);
    border-bottom: 1px solid var(--border-default);
    display: flex;
    align-items: center;
    overflow-x: auto;
  }

  .res-tab {
    height: 100%;
    padding: 0 10px;
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 11px;
    color: var(--text-secondary);
    border-right: 1px solid var(--border-default);
    cursor: pointer;
    white-space: nowrap;
  }

  .res-tab.active {
    background: var(--surface-canvas);
    color: var(--text-primary);
    border-top: 2px solid var(--border-accent);
  }

  .close-x {
    font-size: 9px;
    color: var(--text-muted);
  }

  .result-toolbar {
    height: 26px;
    background: #25272A;
    border-bottom: 1px solid var(--border-default);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 6px;
  }

  .res-toolbar-left, .res-toolbar-right {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .toolbar-text {
    font-size: 10px;
    color: var(--text-muted);
    font-weight: 600;
    margin: 0 2px;
  }

  .grid-table-container {
    flex: 1;
    overflow: auto;
    position: relative;
  }

  .data-grid {
    width: 100%;
    border-collapse: collapse;
    font-size: 12px;
  }

  .data-grid th {
    background: var(--surface-table-header);
    height: var(--table-header-height);
    border-right: 1px solid var(--border-default);
    border-bottom: 1px solid var(--border-default);
    padding: 0 8px;
    text-align: left;
    font-weight: 500;
    color: var(--text-secondary);
    white-space: nowrap;
    position: sticky;
    top: 0;
    z-index: 5;
  }

  .row-num-col {
    width: 32px;
    background: var(--surface-table-header);
  }

  .col-icon {
    margin-right: 4px;
    font-size: 10px;
    color: #56A8F5;
  }

  .filter-funnel, .sort-icon {
    font-size: 9px;
    color: var(--text-muted);
    margin-left: 4px;
    cursor: pointer;
  }

  .data-grid td {
    height: var(--table-row-height);
    border-right: 1px solid var(--border-subtle);
    border-bottom: 1px solid var(--border-subtle);
    padding: 0 8px;
    white-space: nowrap;
  }

  .row-index {
    background: var(--surface-table-header);
    color: var(--text-muted);
    text-align: right;
    padding-right: 6px;
    user-select: none;
    width: 32px;
  }

  .cell.focused {
    outline: 1px solid var(--action-primary);
  }

  .floating-row-badge {
    position: absolute;
    bottom: 10px;
    right: 16px;
    background: var(--surface-floating);
    border: 1px solid var(--border-default);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.45);
    border-radius: 12px;
    padding: 2px 8px;
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 11px;
    color: var(--text-primary);
    cursor: pointer;
    z-index: 10;
  }

  .floating-row-badge .sep { color: var(--border-default); }
  .floating-row-badge .kebab { color: var(--text-muted); }
</style>
