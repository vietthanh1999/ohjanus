<script lang="ts">
  import { appState } from '../../state/appState.svelte';
  import { Icon } from '@ohjanus/icons';

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
        type="button"
        class="action-btn run-btn"
        title="Execute Entire Statement (Cmd+Enter)"
        onclick={handleExecute}
      >
        <Icon name="play" size={13} />
      </button>

      <!-- Execute Under Caret (▶_) -->
      <button type="button" class="action-btn step-btn" title="Execute Statement Under Caret">
        <Icon name="play" size={13} />
      </button>

      <!-- History (🕒) -->
      <button type="button" class="action-btn" title="Query History">
        <Icon name="clock" size={13} />
      </button>

      <!-- Parameter (P) -->
      <button type="button" class="action-btn" title="Parameters">
        <span style="font-size: 11px; font-weight: 600;">(P)</span>
      </button>

      <!-- Settings (⚙️) -->
      <button type="button" class="action-btn" title="Console Settings">
        <Icon name="settings" size={12} />
      </button>

      <!-- Split Layout 🗖 -->
      <button type="button" class="action-btn" title="Toggle Split Orientation">
        <Icon name="layout" size={12} />
      </button>

      <span class="bar-separator"></span>

      <!-- Tx Mode -->
      <button type="button" class="selector-dropdown" title="Transaction Isolation Mode">
        <span>Tx: Auto</span>
        <Icon name="chevron-down" size={8} />
      </button>

      <!-- Playground mode -->
      <label class="playground-chk" title="Sandbox execution without commit">
        <input type="checkbox" />
        <span>Playground</span>
        <Icon name="chevron-down" size={8} />
      </label>
    </div>

    <!-- Target Schema Selector (Right pinned matching design.png) -->
    <div class="toolbar-right">
      <button type="button" class="schema-btn" title="Target Database Schema Context">
        <Icon name="database" size={12} color="#7A7E85" />
        <span class="schema-name">prd_mh_asset.public</span>
        <Icon name="chevron-down" size={8} />
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
      <div class="line-num active-gutter">71</div>
      <div class="line-num">72</div>
      <div class="line-num">73</div>
      <div class="line-num success-gutter">
        <span>74</span>
        <span class="check-mark" title="Executed statement"><Icon name="check" size={10} /></span>
      </div>
    </div>

    <div class="editor-code code-text">
      <!-- Active Execution Highlight Block (lines 64 to 72) matching design.png -->
      <div class="execution-block">
        <!-- Floating Result Badge [✓ 6 ▲ ▼] -->
        <div class="floating-exec-badge" title="Query completed: 6 rows returned">
          <span class="badge-check"><Icon name="check" size={11} /></span>
          <span class="badge-count">6</span>
          <span class="badge-nav"><Icon name="chevron-up" size={10} /></span>
          <span class="badge-nav"><Icon name="chevron-down" size={10} /></span>
        </div>

        <div class="code-line">
          <span class="kw">join</span> <span class="ident">content</span> <span class="alias">ct</span>
          <span class="code-lens">&nbsp;&nbsp;1..n&lt;-&gt;1: on ct.id = ma."contentID"</span>
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
          &nbsp;&nbsp;<span class="kw">AND</span> <span class="alias">c</span>.<span class="str">"timestamp"</span> &nbsp;&nbsp;&nbsp;&lt; <span class="fn">now()</span> - <span class="fn">interval</span> <span class="str">'5 minutes'</span><span class="cursor-bar">|</span> &nbsp;
          <span class="comment">-- streamer poll 10s -&gt; &gt;5 phút <span class="spell-warn">chắc chắn không</span> ai đọc</span>
        </div>
        <div class="code-line">
          <span class="kw">ORDER BY</span> <span class="alias">c</span>.<span class="ident">ordinal</span>;
        </div>
      </div>

      <!-- Blank line 73 -->
      <div class="code-line empty">&nbsp;</div>

      <!-- Line 74: Executed select query -->
      <div class="code-line">
        <span class="kw">select</span> * <span class="kw">from</span> <span class="ident">transfer_job</span> <span class="kw">where</span> <span class="ident">id</span> = <span class="str">'b8262180-1075-43f8-8338-2412d4734d65'</span>
      </div>
    </div>
  </div>

  <!-- 3. QUERY RESULTS MULTI-TAB & ACTION BAR matching design.png -->
  <div class="results-header-tabs">
    <div class="result-tab">
      <Icon name="table" size={12} />
      <span>Result 1</span>
    </div>
    <div class="result-tab">
      <Icon name="table" size={12} />
      <span>Result 1-2</span>
    </div>
    <div class="result-tab">
      <Icon name="table" size={12} />
      <span>prd_mh_asset.public.transfer_job</span>
    </div>
    <div class="result-tab">
      <Icon name="table" size={12} />
      <span>Result 1-4</span>
    </div>
    <!-- Active Result Tab with rounded border pill -->
    <div class="result-tab active">
      <Icon name="table" size={12} />
      <span>prd_mh_asset.public.transfer_job 2</span>
      <span class="close-x">
        <Icon name="close" size={10} />
      </span>
    </div>
    <div class="result-tab-chevron">
      <Icon name="chevron-down" size={10} />
    </div>
  </div>

  <!-- Result Toolbar matching design.png -->
  <div class="results-toolbar">
    <div class="toolbar-left">
      <button type="button" class="jb-icon-btn" title="Grid View">
        <Icon name="table" size={12} />
      </button>
      <button type="button" class="jb-icon-btn" title="Text View">
        <Icon name="audit" size={12} />
      </button>
      <span class="bar-separator"></span>
      <button type="button" class="jb-icon-btn" title="Reload (Cmd+R)" onclick={handleExecute}>
        <Icon name="refresh" size={12} />
      </button>
      <button type="button" class="jb-icon-btn" title="History">
        <Icon name="clock" size={12} />
      </button>
      <button type="button" class="jb-icon-btn" title="Cancel">
        <Icon name="stop" size={11} />
      </button>
      <span class="bar-separator"></span>
      <button type="button" class="jb-icon-btn" title="Add Row">
        <Icon name="plus" size={12} />
      </button>
      <button type="button" class="jb-icon-btn" title="Delete Row">
        <Icon name="minus" size={12} />
      </button>
      <button type="button" class="jb-icon-btn" title="Revert">
        <Icon name="undo" size={12} />
      </button>
      <button type="button" class="jb-icon-btn" title="Commit">
        <Icon name="redo" size={12} />
      </button>
      <span class="bar-separator"></span>
      <button type="button" class="jb-icon-btn" title="Sort Up">
        <Icon name="arrow-up" size={12} />
      </button>
      <button type="button" class="jb-icon-btn" title="Sort Down">
        <Icon name="arrow-down" size={12} />
      </button>
      <span class="bar-separator"></span>
      <span style="color: var(--text-secondary); font-size: 11px; padding: 0 4px;">Tx: Auto</span>
      <span class="bar-separator"></span>
      <button type="button" class="jb-icon-btn" title="View DDL" onclick={() => appState.ddlModalOpen = true}>
        <span style="font-size: 10px; font-weight: 700; color: #7A7E85;">DDL</span>
      </button>
      <span class="bar-separator"></span>
      <button type="button" class="jb-icon-btn" title="Pin Tab">
        <Icon name="pin" size={12} />
      </button>
      <button type="button" class="jb-icon-btn" title="Search in Table">
        <Icon name="search" size={12} />
      </button>
      <button type="button" class="jb-icon-btn" title="Filter Funnel">
        <Icon name="filter" size={12} />
      </button>
      <button type="button" class="jb-icon-btn" title="Statistics">
        <Icon name="chart" size={12} />
      </button>
    </div>

    <div class="toolbar-right">
      <span class="export-dropdown">
        CSV <Icon name="chevron-down" size={8} />
      </span>
      <button type="button" class="jb-icon-btn" title="Export">
        <Icon name="download" size={12} />
      </button>
      <button type="button" class="jb-icon-btn" title="Import">
        <Icon name="upload" size={12} />
      </button>
      <button type="button" class="jb-icon-btn" title="Charts">
        <Icon name="chart" size={12} />
      </button>
      <button type="button" class="jb-icon-btn" title="Options">
        <Icon name="eye" size={12} />
      </button>
      <button type="button" class="jb-icon-btn" title="Settings">
        <Icon name="settings" size={12} />
      </button>
    </div>
  </div>

  <!-- 4. HIGH-DENSITY RESULTS DATA GRID matching design.png -->
  <div class="grid-container">
    <table class="jb-table code-text">
      <thead>
        <tr>
          <th class="row-num-header"></th>
          <th class="col-header" style="width: 320px;">
            <div class="header-inner">
              <Icon name="key" size={12} color="#FACC15" />
              <span>id</span>
              <span class="header-icon"><Icon name="filter" size={9} /></span>
              <span class="header-icon"><Icon name="sort" size={9} /></span>
            </div>
          </th>
          <th class="col-header" style="width: 130px;">
            <div class="header-inner">
              <span style="color: var(--type-general); font-weight: bold; font-size: 11px;">#</span>
              <span>ordinal</span>
              <span class="header-icon"><Icon name="filter" size={9} /></span>
            </div>
          </th>
          <th class="col-header" style="width: 320px;">
            <div class="header-inner">
              <span style="color: var(--type-general); font-size: 11px;">" "</span>
              <span>"connectionCredentialID"</span>
              <span class="header-icon"><Icon name="filter" size={9} /></span>
              <span class="header-icon"><Icon name="sort" size={9} /></span>
            </div>
          </th>
          <th class="col-header">
            <div class="header-inner">
              <span style="color: var(--type-general); font-size: 11px;">" "</span>
              <span>"sourcePath"</span>
              <span class="header-icon"><Icon name="filter" size={9} /></span>
            </div>
          </th>
        </tr>
      </thead>
      <tbody>
        <tr
          class="selected-row"
          onclick={() => selectedRowIndex = 0}
        >
          <td class="row-num-cell">1</td>
          <td class="cell uuid-cell" class:cell-focused={selectedColumn === 'id'} onclick={() => selectedColumn = 'id'}>
            b8262180-1075-43f8-8338-2412d4734d65
          </td>
          <td class="cell num-cell" class:cell-focused={selectedColumn === 'ordinal'} onclick={() => selectedColumn = 'ordinal'}>
            172098
          </td>
          <td class="cell uuid-cell" class:cell-focused={selectedColumn === 'connectionCredentialID'} onclick={() => selectedColumn = 'connectionCredentialID'}>
            753da7e1-523f-4834-ba74-7824e3d5aa52
          </td>
          <td class="cell path-cell truncate" class:cell-focused={selectedColumn === 'sourcePath'} onclick={() => selectedColumn = 'sourcePath'}>
            /home/vod/VOD/as...
          </td>
        </tr>
      </tbody>
    </table>

    <!-- Floating row count pill [ 1 row ⌵ | ⋮ ] -->
    <div class="floating-row-badge" title="Retrieved count">
      <span>1 row</span>
      <span style="color: var(--text-muted); opacity: 0.6;">|</span>
      <Icon name="chevron-down" size={10} />
      <span style="color: var(--text-muted); opacity: 0.6;">|</span>
      <Icon name="more" size={12} />
    </div>
  </div>
</div>

<style>
  .sql-console-view {
    display: flex;
    flex-direction: column;
    height: 100%;
    background-color: var(--bg-canvas);
    overflow: hidden;
  }

  /* 1. Sub-toolbar */
  .sql-toolbar {
    height: var(--toolbar-height);
    background-color: var(--bg-toolbar);
    border-bottom: 1px solid var(--border-subtle);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 8px;
    flex-shrink: 0;
  }

  .toolbar-left, .toolbar-right {
    display: flex;
    align-items: center;
    gap: 3px;
  }

  .action-btn {
    width: var(--icon-btn-size-sm, 26px);
    height: var(--icon-btn-size-sm, 26px);
    border-radius: var(--radius-sm, 4px);
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--text-secondary);
    transition: all 0.1s ease;
  }

  :global(.action-btn svg) {
    width: var(--icon-size-sm, 14px);
    height: var(--icon-size-sm, 14px);
  }

  .action-btn:hover {
    background-color: var(--bg-hover);
    color: var(--text-primary);
  }

  .action-btn.run-btn {
    color: var(--action-success);
  }

  .action-btn.step-btn {
    color: var(--action-success);
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
    height: var(--control-height-xs, 24px);
    padding: 0 6px;
    border-radius: var(--radius-sm, 4px);
    color: var(--text-secondary);
    font-size: var(--font-size-xs, 11px);
  }

  .playground-chk {
    display: flex;
    align-items: center;
    gap: 4px;
    height: var(--control-height-xs, 24px);
    font-size: var(--font-size-xs, 11px);
    color: var(--text-secondary);
    cursor: pointer;
    padding: 0 6px;
  }

  .playground-chk input {
    accent-color: var(--action-primary);
  }

  .schema-btn {
    display: flex;
    align-items: center;
    gap: 6px;
    height: var(--control-height-xs, 24px);
    padding: 0 8px;
    border-radius: var(--radius-sm, 4px);
    color: var(--text-primary);
    font-size: var(--font-size-xs, 11px);
  }

  .schema-btn:hover {
    background-color: var(--bg-hover);
  }

  .schema-name {
    font-weight: 400;
  }

  /* 2. SQL Editor Canvas */
  .editor-container {
    height: 220px;
    min-height: 160px;
    background-color: var(--bg-canvas);
    display: flex;
    overflow: auto;
    position: relative;
    border-bottom: 1px solid var(--border-subtle);
  }

  .editor-gutter {
    width: 44px;
    background-color: var(--bg-canvas);
    border-right: 1px solid var(--border-subtle);
    padding: 6px 0;
    text-align: right;
    user-select: none;
    flex-shrink: 0;
  }

  .line-num {
    height: 20px;
    padding-right: 10px;
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
  }

  .line-num.success-gutter {
    color: var(--text-primary);
  }

  .check-mark {
    color: var(--action-success);
    font-size: 11px;
    font-weight: bold;
    display: inline-flex;
    align-items: center;
  }

  .editor-code {
    flex: 1;
    padding: 6px 12px;
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

  /* Active Execution Block from design.png */
  .execution-block {
    position: relative;
    background-color: var(--bg-execution-block);
    border: 1px solid var(--border-execution);
    border-radius: 2px;
    padding: 1px 4px;
    margin: -1px -4px;
  }

  .floating-exec-badge {
    position: absolute;
    top: 4px;
    right: 8px;
    background-color: #1A3125;
    border: 1px solid #2B5B3E;
    border-radius: 3px;
    padding: 1px 6px;
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 11px;
    z-index: 5;
  }

  .badge-check {
    color: var(--action-success);
    font-weight: bold;
    display: inline-flex;
    align-items: center;
  }

  .badge-count {
    color: #FFFFFF;
    font-weight: 500;
  }

  .badge-nav {
    color: var(--text-muted);
    font-size: 8px;
    cursor: pointer;
    display: inline-flex;
    align-items: center;
  }

  /* Syntax Highlighting */
  .kw { color: var(--syntax-keyword); font-weight: 500; }
  .ident { color: var(--syntax-identifier); }
  .alias { color: var(--syntax-alias); }
  .str { color: var(--syntax-string); }
  .fn { color: var(--syntax-function); }
  .num { color: var(--syntax-number); }
  .comment { color: var(--syntax-comment); }
  .code-lens { color: var(--syntax-code-lens); font-size: 11px; }
  .spell-warn { text-decoration: underline wavy var(--syntax-spell-warn); }
  .cursor-bar { color: #FFFFFF; font-weight: bold; animation: blink 1s step-start infinite; }

  @keyframes blink {
    50% { opacity: 0; }
  }

  /* 3. Results Tabs */
  .results-header-tabs {
    height: 26px;
    background-color: var(--bg-canvas);
    border-bottom: 1px solid var(--border-subtle);
    display: flex;
    align-items: center;
    padding: 0 4px;
    flex-shrink: 0;
    overflow-x: auto;
    gap: 4px;
  }

  .result-tab {
    padding: 0 8px;
    height: 22px;
    display: flex;
    align-items: center;
    gap: 5px;
    font-size: 11px;
    color: var(--text-secondary);
    border-radius: 3px;
    background-color: transparent;
    cursor: pointer;
    user-select: none;
    white-space: nowrap;
  }

  .result-tab.active {
    background-color: #1F2E4A;
    border: 1px solid #3574F0;
    color: #DFE1E5;
  }

  .close-x {
    color: var(--text-muted);
    font-size: 13px;
    margin-left: 2px;
  }

  .result-tab-chevron {
    padding: 0 4px;
    display: flex;
    align-items: center;
    color: var(--text-muted);
    cursor: pointer;
  }

  /* Result Toolbar */
  .results-toolbar {
    height: 26px;
    background-color: var(--bg-toolbar);
    border-bottom: 1px solid var(--border-subtle);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 8px;
    flex-shrink: 0;
  }

  .export-dropdown {
    font-size: 11px;
    color: var(--text-secondary);
    padding: 1px 4px;
    cursor: pointer;
  }

  /* 4. Results Data Grid */
  .grid-container {
    flex: 1;
    overflow: auto;
    position: relative;
    background-color: var(--bg-canvas);
  }

  .jb-table {
    width: 100%;
    border-collapse: collapse;
    font-size: var(--font-size-sm, 13px);
  }

  .row-num-header {
    width: 32px;
    background-color: var(--bg-canvas);
    border-bottom: 1px solid var(--border-subtle);
    border-right: 1px solid var(--border-subtle);
  }

  .col-header {
    height: var(--table-header-height);
    background-color: var(--bg-canvas);
    border-bottom: 1px solid var(--border-subtle);
    border-right: 1px solid var(--border-subtle);
    color: var(--text-secondary);
    font-weight: 400;
    font-size: var(--font-size-xs, 12px);
    padding: 0 8px;
    text-align: left;
    user-select: none;
    white-space: nowrap;
  }

  .header-inner {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .header-icon {
    color: var(--text-muted);
    font-size: var(--font-size-2xs, 11px);
    display: inline-flex;
    align-items: center;
  }

  .row-num-cell {
    width: 32px;
    height: var(--table-row-height);
    background-color: var(--bg-canvas);
    border-bottom: 1px solid #25272A;
    border-right: 1px solid var(--border-subtle);
    color: #7A7E85;
    text-align: right;
    padding-right: 6px;
    font-size: var(--font-size-xs, 12px);
    user-select: none;
  }

  .cell {
    height: var(--table-row-height);
    border-bottom: 1px solid #25272A;
    border-right: 1px solid var(--border-subtle);
    padding: 0 8px;
    color: var(--text-primary);
    white-space: nowrap;
  }

  .selected-row {
    background-color: rgba(46, 58, 78, 0.35);
  }

  .cell-focused {
    outline: 1px solid var(--action-primary);
    background-color: var(--bg-selected);
  }

  .uuid-cell {
    color: #DFE1E5;
    font-size: var(--font-size-base, 14px);
  }

  .num-cell {
    color: var(--syntax-number);
  }

  .path-cell {
    color: #9DA0A8;
  }
</style>
