<script lang="ts">
  import { appState } from '../../state/appState.svelte';

  let activeResultTab = $state('prd_mh_asset.public.transfer_job');
  let selectedRowIndex = $state(0);
  let selectedColumn = $state('id');

  function handleExecute() {
    appState.executeQuery();
  }
</script>

<div class="sql-console-view">
  <!-- 1. SQL SUB-TOOLBAR (Action Bar) -->
  <div class="sql-toolbar">
    <div class="toolbar-left">
      <!-- Execute Entire Statement (▶) -->
      <button
        class="action-btn run-btn"
        title="Execute Entire Statement (Cmd+Enter)"
        onclick={handleExecute}
        disabled={appState.isExecutingQuery}
      >
        <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor">
          <path d="M4 2.5l9 5.5-9 5.5V2.5z"/>
        </svg>
      </button>

      <!-- Execute Under Caret (▶_) -->
      <button class="action-btn step-btn" title="Execute Statement Under Caret">
        <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor">
          <path d="M4 3l6 4-6 4V3zM2 13h12v1.5H2V13z"/>
        </svg>
      </button>

      <!-- History (🕒) -->
      <button class="action-btn" title="Query History">
        <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor">
          <path d="M8 1a7 7 0 100 14A7 7 0 008 1zm0 1.5a5.5 5.5 0 110 11 5.5 5.5 0 010-11zM7.25 4v4.25l3.25 1.95.75-1.23-2.5-1.5V4h-1.5z"/>
        </svg>
      </button>

      <!-- Stop (⏹) -->
      <button class="action-btn stop-btn" title="Cancel / Stop Execution">
        <svg width="12" height="12" viewBox="0 0 16 16" fill="currentColor">
          <rect x="2" y="2" width="12" height="12" rx="1.5"/>
        </svg>
      </button>

      <!-- Settings (⚙️) -->
      <button class="action-btn" title="Console Settings">
        <svg width="13" height="13" viewBox="0 0 16 16" fill="currentColor">
          <path d="M7.07 1a1 1 0 00-.97.757l-.24 1.026a5.526 5.526 0 00-1.282.74L3.6 3.057a1 1 0 00-1.205.247l-.98 1.135a1 1 0 00-.173 1.218l.617.848a5.534 5.534 0 00-.012 1.48l-.618.847a1 1 0 00.173 1.218l.98 1.136a1 1 0 001.206.246l.978-.466c.394.3.826.55 1.282.74l.24 1.026A1 1 0 007.07 15h1.86a1 1 0 00.97-.757l.24-1.026c.456-.19.888-.44 1.282-.74l.978.466a1 1 0 001.206-.246l.98-1.136a1 1 0 00-.173-1.218l-.617-.847c.105-.486.105-.993 0-1.48l.617-.848a1 1 0 00.173-1.218l-.98-1.135a1 1 0 00-1.206-.247l-.978.466a5.527 5.527 0 00-1.282-.74l-.24-1.026A1 1 0 008.93 1H7.07zm.93 5a2 2 0 110 4 2 2 0 010-4z"/>
        </svg>
      </button>

      <span class="bar-separator"></span>

      <!-- Tx Mode -->
      <button class="selector-dropdown" title="Transaction Isolation Mode">
        <span>Tx: Auto</span>
        <svg width="9" height="9" viewBox="0 0 16 16" fill="currentColor">
          <path d="M4 6l4 4 4-4H4z" />
        </svg>
      </button>

      <!-- Playground mode -->
      <label class="playground-chk" title="Sandbox execution without commit">
        <input type="checkbox" />
        <span>Playground</span>
        <svg width="9" height="9" viewBox="0 0 16 16" fill="currentColor">
          <path d="M4 6l4 4 4-4H4z" />
        </svg>
      </label>
    </div>

    <!-- Target Schema Selector (Right pinned) -->
    <div class="toolbar-right">
      <button class="schema-btn" title="Current Connected Schema Context">
        <svg width="12" height="12" viewBox="0 0 16 16" fill="#56A8F5">
          <path d="M8 1c3.866 0 7 .895 7 2v10c0 1.105-3.134 2-7 2s-7-.895-7-2V3c0-1.105 3.134-2 7-2z"/>
        </svg>
        <span class="schema-name">prd_mh_asset.public</span>
        <svg width="9" height="9" viewBox="0 0 16 16" fill="currentColor">
          <path d="M4 6l4 4 4-4H4z" />
        </svg>
      </button>
    </div>
  </div>

  <!-- 2. SQL EDITOR CANVAS -->
  <div class="editor-container">
    <div class="editor-gutter code-text">
      <div class="line-num">64</div>
      <div class="line-num">65</div>
      <div class="line-num">66</div>
      <div class="line-num">67</div>
      <div class="line-num">68</div>
      <div class="line-num">69</div>
      <div class="line-num">70</div>
      <div class="line-num">71</div>
      <div class="line-num">72</div>
      <div class="line-num">73</div>
      <div class="line-num active-gutter">
        <span>74</span>
        <span class="check-mark" title="Executed statement">✓</span>
      </div>
    </div>

    <div class="editor-code code-text">
      <!-- Active Execution Highlight Block (lines 64 to 72) -->
      <div class="execution-block">
        <!-- Floating Result Badge [✓ 6 ▲ ▼] -->
        <div class="floating-exec-badge" title="Statement Result: 6 rows returned">
          <span class="badge-check">✓</span>
          <span class="badge-count">6</span>
          <span class="badge-nav">▲</span>
          <span class="badge-nav">▼</span>
        </div>

        <div class="code-line">
          <span class="kw">join</span> <span class="ident">content</span> <span class="alias">ct</span>
          <span class="code-lens">&nbsp;&nbsp;1..n &lt;-&gt; 1: on ct.id = ma."contentID"</span>
        </div>
        <div class="code-line">
          <span class="kw">LEFT JOIN</span> <span class="ident">transfer_job_queue</span> <span class="alias">q</span> <span class="kw">ON</span> <span class="alias">q</span>.<span class="str">"jobID"</span> = <span class="alias">c</span>.<span class="ident">process_id</span>
        </div>
        <div class="code-line">
          <span class="kw">WHERE</span> <span class="alias">c</span>.<span class="ident">participant</span> = <span class="str">'media-transferer'</span>
        </div>
        <div class="code-line">
          &nbsp;&nbsp;<span class="kw">AND</span> <span class="alias">c</span>.<span class="ident">process_name</span> &nbsp;&nbsp;= <span class="str">'media-ingest'</span>
        </div>
        <div class="code-line">
          &nbsp;&nbsp;<span class="kw">AND</span> <span class="alias">c</span>.<span class="ident">name</span> &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;= <span class="str">'create-ingest-job'</span>
        </div>
        <div class="code-line">
          &nbsp;&nbsp;<span class="kw">AND</span> <span class="alias">c</span>.<span class="ident">retry_count</span> &nbsp;&nbsp;&nbsp;= <span class="num">0</span>
        </div>
        <div class="code-line">
          &nbsp;&nbsp;<span class="kw">AND</span> <span class="alias">c</span>.<span class="ident">schedule</span> <span class="kw">IS NULL</span>
        </div>
        <div class="code-line">
          &nbsp;&nbsp;<span class="kw">AND</span> <span class="alias">c</span>.<span class="str">"timestamp"</span> &nbsp;&nbsp;&nbsp;&lt; <span class="fn">now()</span> - <span class="fn">interval</span> <span class="str">'5 minutes'</span> &nbsp;
          <span class="comment">-- streamer poll 10s -&gt; <span class="spell-warn">&gt;5 phút chắc chắn không ai đọc</span></span>
        </div>
        <div class="code-line">
          <span class="kw">ORDER BY</span> <span class="alias">c</span>.<span class="ident">ordinal</span>;
        </div>
      </div>

      <!-- Blank line 73 -->
      <div class="code-line empty">&nbsp;</div>

      <!-- Line 74: Executed select query -->
      <div class="code-line active-line">
        <span class="kw">select</span> * <span class="kw">from</span> <span class="ident">transfer_job</span> <span class="kw">where</span> <span class="ident">id</span> = <span class="str">'b8262180-1075-43f8-8338-2412d4734d65'</span>
      </div>
    </div>
  </div>

  <!-- 3. QUERY RESULTS MULTI-TAB & ACTION BAR -->
  <div class="results-header-tabs">
    <div class="result-tab">Result 1</div>
    <div class="result-tab">Result 1-2</div>
    <div class="result-tab active">
      <span>prd_mh_asset.public.transfer_job</span>
    </div>
    <div class="result-tab">Result 1-4</div>
    <div class="result-tab">prd_mh_asset.public.transfer_job 2 <span class="close-x">×</span></div>
    <div class="result-tab-chevron">▾</div>
  </div>

  <!-- Result Toolbar -->
  <div class="results-toolbar">
    <div class="toolbar-left">
      <button class="jb-icon-btn" title="Grid View"><span style="font-size: 11px;">⊞</span></button>
      <button class="jb-icon-btn" title="Text View"><span style="font-size: 11px;">🗎</span></button>
      <span class="bar-separator"></span>
      <button class="jb-icon-btn" title="Reload (Cmd+R)" onclick={handleExecute}><span style="font-size: 11px;">⟳</span></button>
      <button class="jb-icon-btn" title="History"><span style="font-size: 11px;">🕒</span></button>
      <button class="jb-icon-btn" title="Cancel"><span style="font-size: 11px;">⏹</span></button>
      <span class="bar-separator"></span>
      <button class="jb-icon-btn" title="Add Row"><span style="font-size: 11px;">+</span></button>
      <button class="jb-icon-btn" title="Delete Row"><span style="font-size: 11px;">-</span></button>
      <button class="jb-icon-btn" title="Revert"><span style="font-size: 11px;">↩</span></button>
      <button class="jb-icon-btn" title="Commit"><span style="font-size: 11px;">↪</span></button>
      <span class="bar-separator"></span>
      <button class="jb-icon-btn" title="Sort Up"><span style="font-size: 11px;">↑</span></button>
      <button class="jb-icon-btn" title="Sort Down"><span style="font-size: 11px;">↓</span></button>
      <span class="bar-separator"></span>
      <span style="color: var(--text-secondary); font-size: 11px; padding: 0 4px;">Tx: Auto</span>
      <span class="bar-separator"></span>
      <button class="jb-icon-btn" title="View DDL" onclick={() => appState.ddlModalOpen = true}>
        <span style="font-size: 10px; font-weight: 600; color: #4A88C7;">DDL</span>
      </button>
      <span class="bar-separator"></span>
      <button class="jb-icon-btn" title="Pin Tab"><span style="font-size: 11px;">📌</span></button>
      <button class="jb-icon-btn" title="Search in Table"><span style="font-size: 11px;">🔍</span></button>
      <button class="jb-icon-btn" title="Filter"><span style="font-size: 11px;">Y</span></button>
      <button class="jb-icon-btn" title="Statistics"><span style="font-size: 11px;">📊</span></button>
    </div>

    <div class="toolbar-right">
      <span class="export-dropdown">CSV ▾</span>
      <button class="jb-icon-btn" title="Export to File"><span>⤓</span></button>
      <button class="jb-icon-btn" title="Import from File"><span>⤒</span></button>
      <button class="jb-icon-btn" title="Charts"><span>📈</span></button>
      <button class="jb-icon-btn" title="View Options"><span>👁️</span></button>
      <button class="jb-icon-btn" title="Settings"><span>⚙️</span></button>
    </div>
  </div>

  <!-- 4. HIGH-DENSITY RESULTS DATA GRID -->
  <div class="grid-container">
    <table class="jb-table code-text">
      <thead>
        <tr>
          <th class="row-num-header"></th>
          <th class="col-header">
            <div class="header-inner">
              <svg width="12" height="12" viewBox="0 0 16 16" fill="var(--type-pk)">
                <path d="M0 8a4 4 0 017.465-2H14a2 2 0 012 2v1a1 1 0 01-1 1h-1v1a1 1 0 01-1 1h-1v1a1 1 0 01-1 1H9.465A4 4 0 010 8zm4-2a2 2 0 100 4 2 2 0 000-4z"/>
              </svg>
              <span>id</span>
              <span class="header-icon">Y</span>
              <span class="header-icon">⇅</span>
            </div>
          </th>
          <th class="col-header">
            <div class="header-inner">
              <span style="color: var(--type-general); font-weight: bold; font-size: 11px;">#</span>
              <span>ordinal</span>
            </div>
          </th>
          <th class="col-header">
            <div class="header-inner">
              <span style="color: var(--type-general); font-size: 11px;">" "</span>
              <span>connectionCredentialID</span>
              <span class="header-icon">Y</span>
              <span class="header-icon">⇅</span>
            </div>
          </th>
          <th class="col-header">
            <div class="header-inner">
              <span style="color: var(--type-general); font-size: 11px;">" "</span>
              <span>sourcePath</span>
              <span class="header-icon">Y</span>
            </div>
          </th>
          <th class="col-header">
            <div class="header-inner">
              <span style="color: var(--type-general); font-size: 11px;">" "</span>
              <span>targetPath</span>
            </div>
          </th>
          <th class="col-header">
            <div class="header-inner">
              <span>status</span>
            </div>
          </th>
        </tr>
      </thead>
      <tbody>
        {#each appState.transferJobResults as row, idx}
          <tr
            class:selected-row={selectedRowIndex === idx}
            onclick={() => selectedRowIndex = idx}
          >
            <td class="row-num-cell">{idx + 1}</td>
            <td
              class="cell uuid-cell"
              class:cell-focused={selectedRowIndex === idx && selectedColumn === 'id'}
              onclick={() => selectedColumn = 'id'}
            >
              {row.id}
            </td>
            <td
              class="cell num-cell"
              class:cell-focused={selectedRowIndex === idx && selectedColumn === 'ordinal'}
              onclick={() => selectedColumn = 'ordinal'}
            >
              {row.ordinal}
            </td>
            <td
              class="cell uuid-cell"
              class:cell-focused={selectedRowIndex === idx && selectedColumn === 'connectionCredentialID'}
              onclick={() => selectedColumn = 'connectionCredentialID'}
            >
              {row.connectionCredentialID}
            </td>
            <td
              class="cell path-cell truncate"
              class:cell-focused={selectedRowIndex === idx && selectedColumn === 'sourcePath'}
              onclick={() => selectedColumn = 'sourcePath'}
              title={row.sourcePath}
            >
              {row.sourcePath}
            </td>
            <td
              class="cell path-cell truncate"
              class:cell-focused={selectedRowIndex === idx && selectedColumn === 'targetPath'}
              onclick={() => selectedColumn = 'targetPath'}
              title={row.targetPath}
            >
              {row.targetPath}
            </td>
            <td class="cell status-cell">
              {#if row.status === 'COMPLETED'}
                <span style="color: var(--action-success);">COMPLETED</span>
              {:else if row.status === 'PROCESSING'}
                <span style="color: var(--syntax-function);">PROCESSING</span>
              {:else if row.status === 'QUEUED'}
                <span style="color: var(--action-warning);">QUEUED</span>
              {:else}
                <span style="color: var(--action-danger);">{row.status}</span>
              {/if}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>

    <!-- Floating row count pill [ 6 rows ▾ | ⋮ ] -->
    <div class="floating-row-badge" title="Page size: 50">
      <span>{appState.transferJobResults.length} rows ▾</span>
      <span style="color: var(--text-muted);">|</span>
      <span>⋮</span>
    </div>
  </div>
</div>

<style>
  .sql-console-view {
    display: flex;
    flex-direction: column;
    height: calc(100vh - var(--titlebar-height) - var(--tabbar-height) - var(--statusbar-height));
    background-color: var(--bg-canvas);
    overflow: hidden;
  }

  /* 1. Sub-toolbar */
  .sql-toolbar {
    height: var(--toolbar-height);
    background-color: var(--bg-toolbar);
    border-bottom: 1px solid var(--border-default);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 10px;
    flex-shrink: 0;
  }

  .toolbar-left, .toolbar-right {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .action-btn {
    width: 22px;
    height: 22px;
    border-radius: 3px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--text-secondary);
    transition: all 0.1s ease;
  }

  .action-btn:hover {
    background-color: var(--bg-hover);
    color: var(--text-primary);
  }

  .action-btn.run-btn {
    color: var(--action-success);
  }

  .action-btn.run-btn:hover {
    color: var(--action-success-hover);
    background-color: rgba(87, 211, 140, 0.15);
  }

  .action-btn.step-btn {
    color: var(--action-success);
  }

  .action-btn.stop-btn {
    color: var(--action-danger);
  }

  .action-btn.stop-btn:hover {
    color: var(--action-danger-hover);
    background-color: rgba(229, 83, 83, 0.15);
  }

  .bar-separator {
    width: 1px;
    height: 14px;
    background-color: var(--border-default);
    margin: 0 4px;
  }

  .selector-dropdown {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 2px 6px;
    border-radius: 3px;
    color: var(--text-secondary);
    font-size: 11px;
  }

  .selector-dropdown:hover {
    background-color: var(--bg-hover);
    color: var(--text-primary);
  }

  .playground-chk {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 11px;
    color: var(--text-secondary);
    cursor: pointer;
    padding: 2px 6px;
  }

  .playground-chk input {
    accent-color: var(--action-primary);
  }

  .schema-btn {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 2px 8px;
    border: 1px solid var(--border-default);
    border-radius: 3px;
    background-color: #1E1F22;
    color: var(--text-primary);
    font-size: 11px;
  }

  .schema-btn:hover {
    background-color: var(--bg-hover);
  }

  .schema-name {
    font-weight: 500;
  }

  /* 2. SQL Editor Canvas */
  .editor-container {
    height: 220px;
    min-height: 160px;
    background-color: var(--bg-canvas);
    display: flex;
    overflow: auto;
    position: relative;
    border-bottom: 1px solid var(--border-default);
  }

  .editor-gutter {
    width: 48px;
    background-color: var(--bg-canvas);
    border-right: 1px solid var(--border-subtle);
    padding: 8px 0;
    text-align: right;
    user-select: none;
    flex-shrink: 0;
  }

  .line-num {
    height: 20px;
    padding-right: 12px;
    color: var(--text-muted);
    font-size: 12px;
    line-height: 20px;
    display: flex;
    align-items: center;
    justify-content: flex-end;
    gap: 4px;
  }

  .line-num.active-gutter {
    color: var(--text-primary);
    font-weight: 600;
  }

  .check-mark {
    color: var(--action-success);
    font-size: 11px;
    font-weight: bold;
  }

  .editor-code {
    flex: 1;
    padding: 8px 12px;
    font-size: 13px;
    line-height: 20px;
    overflow-x: auto;
    color: var(--text-primary);
  }

  .code-line {
    height: 20px;
    white-space: pre;
    display: flex;
    align-items: center;
  }

  .code-line.empty {
    height: 20px;
  }

  .code-line.active-line {
    background-color: rgba(255, 255, 255, 0.03);
  }

  /* Active Execution Block */
  .execution-block {
    position: relative;
    background-color: var(--bg-execution-block);
    border: 1px solid var(--border-execution);
    border-radius: 2px;
    padding: 2px 4px;
    margin: -2px -4px;
  }

  .floating-exec-badge {
    position: absolute;
    top: 4px;
    right: 8px;
    background-color: #1E2B37;
    border: 1px solid #2B5B9E;
    border-radius: 3px;
    padding: 1px 6px;
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 11px;
    z-index: 5;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.4);
  }

  .badge-check {
    color: var(--action-success);
    font-weight: bold;
  }

  .badge-count {
    color: var(--text-primary);
    font-weight: 600;
  }

  .badge-nav {
    color: var(--text-muted);
    font-size: 9px;
    cursor: pointer;
  }

  /* Syntax Highlighting */
  .kw { color: var(--syntax-keyword); font-weight: 500; }
  .ident { color: var(--syntax-identifier); }
  .alias { color: var(--syntax-alias); }
  .str { color: var(--syntax-string); }
  .fn { color: var(--syntax-function); font-style: italic; }
  .num { color: var(--syntax-number); }
  .comment { color: var(--syntax-comment); font-style: italic; }
  .code-lens { color: var(--syntax-code-lens); font-style: italic; font-size: 11px; user-select: none; }
  .spell-warn { text-decoration: underline wavy var(--syntax-spell-warn); }

  /* 3. Results Tabs */
  .results-header-tabs {
    height: 26px;
    background-color: #25272A;
    border-bottom: 1px solid var(--border-default);
    display: flex;
    align-items: stretch;
    flex-shrink: 0;
    overflow-x: auto;
  }

  .result-tab {
    padding: 0 12px;
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 11px;
    color: var(--text-secondary);
    border-right: 1px solid var(--border-subtle);
    background-color: var(--bg-inactive-tab);
    cursor: pointer;
    user-select: none;
    white-space: nowrap;
  }

  .result-tab.active {
    background-color: var(--bg-canvas);
    color: var(--text-primary);
    font-weight: 500;
    border-top: 2px solid var(--border-accent);
  }

  .close-x {
    color: var(--text-muted);
    font-size: 12px;
  }

  .result-tab-chevron {
    padding: 0 8px;
    display: flex;
    align-items: center;
    color: var(--text-muted);
    cursor: pointer;
  }

  /* Result Toolbar */
  .results-toolbar {
    height: 28px;
    background-color: var(--bg-toolbar);
    border-bottom: 1px solid var(--border-default);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 8px;
    flex-shrink: 0;
  }

  .export-dropdown {
    font-size: 11px;
    color: var(--text-secondary);
    padding: 2px 6px;
    border-radius: 3px;
    cursor: pointer;
  }

  /* 4. High-Density Results Data Grid */
  .grid-container {
    flex: 1;
    overflow: auto;
    position: relative;
    background-color: var(--bg-canvas);
  }

  .jb-table {
    width: 100%;
    border-collapse: collapse;
    font-size: 12px;
  }

  .row-num-header {
    width: 32px;
    background-color: var(--bg-table-header);
    border-bottom: 1px solid var(--border-default);
    border-right: 1px solid var(--border-default);
  }

  .col-header {
    height: var(--table-header-height);
    background-color: var(--bg-table-header);
    border-bottom: 1px solid var(--border-default);
    border-right: 1px solid var(--border-default);
    color: var(--text-secondary);
    font-weight: 500;
    font-size: 11px;
    padding: 0 6px;
    text-align: left;
    user-select: none;
    white-space: nowrap;
  }

  .header-inner {
    display: flex;
    align-items: center;
    gap: 5px;
  }

  .header-icon {
    color: var(--text-muted);
    font-size: 9px;
  }

  .row-num-cell {
    width: 32px;
    height: var(--table-row-height);
    background-color: var(--bg-table-header);
    border-bottom: 1px solid var(--border-subtle);
    border-right: 1px solid var(--border-default);
    color: var(--text-muted);
    text-align: right;
    padding-right: 6px;
    font-size: 11px;
    user-select: none;
  }

  .cell {
    height: var(--table-row-height);
    border-bottom: 1px solid var(--border-subtle);
    border-right: 1px solid var(--border-subtle);
    padding: 0 8px;
    color: var(--text-primary);
    white-space: nowrap;
    max-width: 260px;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .selected-row {
    background-color: rgba(46, 58, 78, 0.4);
  }

  .cell-focused {
    outline: 1px solid var(--action-primary);
    background-color: var(--bg-selected);
  }

  .uuid-cell {
    color: #DFE1E5;
    font-size: 11px;
  }

  .num-cell {
    color: var(--syntax-number);
  }

  .path-cell {
    color: #9DA0A8;
  }

  .status-cell {
    font-weight: 500;
    font-size: 11px;
  }
</style>
