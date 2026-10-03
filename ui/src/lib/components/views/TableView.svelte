<script lang="ts">
  import { appStore } from '../../appStore.svelte';

  let selectedCell = $state({ row: 0, col: 'id' });
  let whereInput = $state('');
  let orderByInput = $state('');
  let sortField = $state<string | null>(null);
  let sortAsc = $state(true);

  // Derived filtered & sorted rows
  let filteredRows = $derived(() => {
    let rows = [...appStore.tableRows];

    if (whereInput.trim()) {
      const q = whereInput.toLowerCase();
      rows = rows.filter(r => 
        (r.id && r.id.toLowerCase().includes(q)) ||
        (r.createdBy && r.createdBy.toLowerCase().includes(q)) ||
        (r.createdDate && r.createdDate.toLowerCase().includes(q)) ||
        (q.includes('null') && r.lastUpdatedDate === null)
      );
    }

    if (sortField) {
      rows.sort((a: any, b: any) => {
        const valA = a[sortField!];
        const valB = b[sortField!];
        if (valA === null) return 1;
        if (valB === null) return -1;
        if (valA < valB) return sortAsc ? -1 : 1;
        if (valA > valB) return sortAsc ? 1 : -1;
        return 0;
      });
    }

    return rows;
  });

  function toggleSort(field: string) {
    if (sortField === field) {
      if (sortAsc) {
        sortAsc = false;
      } else {
        sortField = null;
      }
    } else {
      sortField = field;
      sortAsc = true;
    }
    orderByInput = sortField ? `${sortField} ${sortAsc ? 'ASC' : 'DESC'}` : '';
  }

  function clearLogs() {
    appStore.consoleLogs = [];
  }
</script>

<div class="table-view font-mono">
  <!-- Table Toolbar -->
  <div class="table-toolbar">
    <div class="tb-left">
      <button class="icon-btn" title="Refresh Table (Cmd+R)" onclick={() => appStore.runQuery()}>
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M21.5 2v6h-6M21.34 15.57a10 10 0 1 1-.57-8.38l5.67-5.67"/></svg>
      </button>
      <button class="icon-btn" title="History">🕒</button>
      <button class="icon-btn" title="Stop">⏹</button>
      <button class="icon-btn" title="Add Row">+</button>
      <button class="icon-btn" title="Delete Row">-</button>
      <button class="icon-btn" title="Revert Changes">↩</button>
      <button class="icon-btn" title="Commit Changes">↪</button>
      <button class="icon-btn" title="Previous Page">↑</button>
      <button class="icon-btn" title="Next Page">↓</button>

      <span class="toolbar-divider"></span>
      <span class="toolbar-text">Tx: Auto ▾</span>
      <span class="toolbar-text">DDL</span>
      <button class="icon-btn" title="Pin">📌</button>
      <button class="icon-btn" title="Search">🔍</button>
      <button class="icon-btn" title="Column Filter">Y</button>
      <button class="icon-btn" title="Show Totals">📊</button>
    </div>

    <div class="tb-right">
      <span class="dropdown-pill">CSV ▾</span>
      <button class="icon-btn" title="Export">⤓</button>
      <button class="icon-btn" title="Import">⤒</button>
      <button class="icon-btn" title="Charts">📈</button>
      <button class="icon-btn" title="Column Visibility">👁️</button>
      <button class="icon-btn" title="Settings">⚙️</button>
    </div>
  </div>

  <!-- Inline WHERE / ORDER BY Filter Bar -->
  <div class="filter-clause-bar">
    <div class="filter-field where-field">
      <span class="field-icon">Y</span>
      <span class="field-label">WHERE</span>
      <input 
        type="text" 
        placeholder="id IS NOT NULL" 
        bind:value={whereInput}
      />
    </div>

    <div class="field-separator"></div>

    <div class="filter-field order-field">
      <span class="field-icon">⇅</span>
      <span class="field-label">ORDER BY</span>
      <input 
        type="text" 
        placeholder="createdDate DESC" 
        bind:value={orderByInput}
      />
    </div>
  </div>

  <!-- Main High Density Data Grid -->
  <div class="table-grid-wrapper">
    <table class="data-grid">
      <thead>
        <tr>
          <th class="row-num-col"></th>
          <th class="col-header" onclick={() => toggleSort('id')}>
            <span class="col-icon key-icon">🔑</span>
            <span class="col-name">id</span>
            <span class="filter-funnel">Y</span>
            <span class="sort-icon">{sortField === 'id' ? (sortAsc ? '▲' : '▼') : '⇅'}</span>
          </th>
          <th class="col-header" onclick={() => toggleSort('createdDate')}>
            <span class="col-icon date-icon">📅</span>
            <span class="col-name">createdDate</span>
            <span class="filter-funnel">Y</span>
            <span class="sort-icon">{sortField === 'createdDate' ? (sortAsc ? '▲' : '▼') : '⇅'}</span>
          </th>
          <th class="col-header" onclick={() => toggleSort('lastUpdatedDate')}>
            <span class="col-icon date-icon">📅</span>
            <span class="col-name">lastUpdatedDate</span>
            <span class="filter-funnel">Y</span>
            <span class="sort-icon">{sortField === 'lastUpdatedDate' ? (sortAsc ? '▲' : '▼') : '⇅'}</span>
          </th>
          <th class="col-header" onclick={() => toggleSort('createdBy')}>
            <span class="col-icon user-icon">👤</span>
            <span class="col-name">createdBy</span>
            <span class="filter-funnel">Y</span>
            <span class="sort-icon">{sortField === 'createdBy' ? (sortAsc ? '▲' : '▼') : '⇅'}</span>
          </th>
        </tr>
      </thead>
      <tbody>
        {#each filteredRows() as row, idx}
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
              class:focused={selectedCell.row === idx && selectedCell.col === 'createdDate'}
              onclick={() => selectedCell = { row: idx, col: 'createdDate' }}
            >
              {row.createdDate}
            </td>
            <td 
              class="cell" 
              class:focused={selectedCell.row === idx && selectedCell.col === 'lastUpdatedDate'}
              onclick={() => selectedCell = { row: idx, col: 'lastUpdatedDate' }}
            >
              {#if row.lastUpdatedDate === null}
                <span class="text-null">&lt;null&gt;</span>
              {:else}
                {row.lastUpdatedDate}
              {/if}
            </td>
            <td 
              class="cell" 
              class:focused={selectedCell.row === idx && selectedCell.col === 'createdBy'}
              onclick={() => selectedCell = { row: idx, col: 'createdBy' }}
            >
              {row.createdBy}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>

    <!-- Floating Row Count Pill -->
    <div class="floating-row-badge animate-fade-in">
      <span class="text">58 rows ▾</span>
      <span class="sep">|</span>
      <span class="kebab">⋮</span>
    </div>
  </div>

  <!-- Bottom Log Terminal Pane -->
  <div class="terminal-pane">
    <div class="log-content">
      {#each appStore.consoleLogs as log}
        <div class="log-line">
          {#if log.includes('SELECT')}
            <span class="ts">[2026-10-03 10:50:34]</span>
            <span class="prompt">dev_mh_asset.public&gt; </span>
            <span class="sql-kw">SELECT </span><span class="sql-id">t.*</span><br/>
            <span class="sql-indent">                                            </span>
            <span class="sql-kw">FROM </span><span class="sql-id">public.connection_credential t</span><br/>
            <span class="sql-indent">                                            </span>
            <span class="sql-kw">LIMIT </span><span class="sql-num">501</span>
          {:else if log.includes('rows retrieved')}
            <span class="ts">[2026-10-03 10:50:34]</span>
            <span class="retrieval">{log.replace(/\[.*?\]\s*/, '')}</span>
          {:else}
            <span class="ts">[2026-10-03 10:50:34]</span>
            <span class="event">{log.replace(/\[.*?\]\s*/, '')}</span>
          {/if}
        </div>
      {/each}
    </div>

    <!-- Right Action Strip -->
    <div class="terminal-action-strip">
      <button class="icon-btn" title="Toggle Output Window">📄</button>
      <button class="icon-btn" title="Clear Console Log" onclick={clearLogs}>🗑️</button>
      <button class="icon-btn" title="Toggle Soft-wrap">↩</button>
      <button class="icon-btn" title="Scroll to End">⤓</button>
      <button class="icon-btn" title="Print Log">🖨️</button>
    </div>
  </div>
</div>

<style>
  .table-view {
    display: flex;
    flex-direction: column;
    height: 100%;
    background: var(--surface-canvas);
    outline: none;
    overflow: hidden;
  }

  .table-toolbar {
    height: var(--toolbar-height);
    background: var(--surface-toolbar);
    border-bottom: 1px solid var(--border-default);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 8px;
    font-size: 11px;
  }

  .tb-left, .tb-right {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .toolbar-divider {
    width: 1px;
    height: 16px;
    background: var(--border-default);
    margin: 0 4px;
  }

  .toolbar-text {
    font-size: 11px;
    color: var(--text-muted);
    font-weight: 500;
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

  /* WHERE / ORDER BY filter bar */
  .filter-clause-bar {
    height: var(--filter-bar-height);
    background: #2B2D30;
    border-bottom: 1px solid var(--border-default);
    display: flex;
    align-items: center;
    padding: 0 8px;
    gap: 8px;
  }

  .filter-field {
    display: flex;
    align-items: center;
    gap: 6px;
    flex: 1;
  }

  .field-icon {
    font-size: 10px;
    color: var(--text-muted);
  }

  .field-label {
    font-size: 11px;
    font-weight: 600;
    color: var(--text-secondary);
  }

  .filter-field input {
    flex: 1;
    height: 22px;
    background: #1E1F22;
    border: 1px solid var(--border-default);
    border-radius: 3px;
    font-family: var(--font-code);
    font-size: 12px;
    color: var(--text-primary);
    padding: 2px 6px;
  }

  .field-separator {
    width: 1px;
    height: 18px;
    background: var(--border-default);
  }

  /* High Density Table Grid */
  .table-grid-wrapper {
    flex: 1;
    overflow: auto;
    position: relative;
    background: var(--surface-canvas);
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
    cursor: pointer;
  }

  .row-num-col {
    width: 36px;
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
  }

  .data-grid td {
    height: var(--table-row-height);
    border-right: 1px solid var(--border-subtle);
    border-bottom: 1px solid var(--border-subtle);
    padding: 0 8px;
    white-space: nowrap;
    color: var(--text-primary);
  }

  .row-index {
    background: var(--surface-table-header);
    color: var(--text-muted);
    text-align: right;
    padding-right: 6px;
    user-select: none;
    width: 36px;
  }

  .cell.focused {
    outline: 1px solid var(--action-primary);
  }

  .floating-row-badge {
    position: absolute;
    bottom: 12px;
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

  /* Bottom Terminal Pane */
  .terminal-pane {
    height: 140px;
    background: #1E1F22;
    border-top: 1px solid var(--border-default);
    display: flex;
    overflow: hidden;
  }

  .log-content {
    flex: 1;
    overflow-y: auto;
    padding: 6px 12px;
    font-size: 11px;
    line-height: 1.5;
  }

  .log-line {
    margin-bottom: 3px;
    color: var(--text-primary);
  }

  .ts {
    color: var(--text-muted);
    margin-right: 6px;
  }

  .prompt {
    color: var(--text-primary);
  }

  .sql-kw { color: #CF8E6D; font-weight: 600; }
  .sql-id { color: #DFE1E5; }
  .sql-num { color: #6897BB; }
  .retrieval { color: #DFE1E5; }

  .terminal-action-strip {
    width: 28px;
    border-left: 1px solid var(--border-default);
    background: #25272A;
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 4px 0;
    gap: 4px;
  }
</style>
