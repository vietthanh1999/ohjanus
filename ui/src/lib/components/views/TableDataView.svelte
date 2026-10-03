<script lang="ts">
  import { appState } from '../../state/appState.svelte';

  let selectedRowIndex = $state(0);
  let selectedColumn = $state('id');

  const rowsData = [
    {
      id: '2389a9f7-9ed7-4b7c-9b7c-7eec030b45f5',
      createdDate: '2024-10-29 10:06:03.247644',
      lastUpdatedDate: null,
      createdBy: '4ebfa335-5656-4052-9001-'
    },
    {
      id: 'a4430e1f-3b89-4088-bed8-8a426ab293f0',
      createdDate: '2024-10-31 03:56:42.776700',
      lastUpdatedDate: null,
      createdBy: '0f379df2-47d2-47e7-b195-'
    },
    {
      id: '3be52421-235b-412e-b420-23b612770315',
      createdDate: '2024-11-05 07:49:06.550088',
      lastUpdatedDate: null,
      createdBy: '0f379df2-47d2-47e7-b195-'
    },
    {
      id: '6e98ccab-d08c-49af-bd10-70ca0d510879',
      createdDate: '2024-10-31 03:58:27.775837',
      lastUpdatedDate: '2024-10-31 04:05:44.633942',
      createdBy: '0f379df2-47d2-47e7-b195-'
    },
    {
      id: 'ca1d2eee-6d07-48b5-af45-f551e984f3ed',
      createdDate: '2024-11-14 03:30:27.290163',
      lastUpdatedDate: null,
      createdBy: '0f379df2-47d2-47e7-b195-'
    },
    {
      id: '9e768c19-b3e5-4027-80a9-25a83307fa14',
      createdDate: '2024-11-14 03:34:10.424030',
      lastUpdatedDate: null,
      createdBy: '0f379df2-47d2-47e7-b195-'
    },
    {
      id: '09a9bcb6-cb58-47fe-bc61-e3f416acd3f6',
      createdDate: '2024-11-05 07:30:56.146845',
      lastUpdatedDate: '2024-11-05 07:31:08.435963',
      createdBy: '0f379df2-47d2-47e7-b195-'
    },
    {
      id: '29ebafe4-3484-4595-a2a5-37158b8cd0b5',
      createdDate: '2024-10-31 04:07:39.501833',
      lastUpdatedDate: '2024-11-05 07:40:31.640001',
      createdBy: '0f379df2-47d2-47e7-b195-'
    },
    {
      id: '53d2afe2-1043-4eb2-b488-1fdfa096d870',
      createdDate: '2024-11-14 0',
      lastUpdatedDate: '2024-11-18 10:30:42.374558',
      createdBy: '0f379df2-47d2-47e7-b195-'
    }
  ];

  let displayRows = $derived.by(() => {
    let list = [...rowsData];
    if (appState.whereFilter.trim()) {
      const q = appState.whereFilter.toLowerCase().trim();
      list = list.filter(r => {
        if (q.includes('is null')) return r.lastUpdatedDate === null;
        if (q.includes('is not null')) return r.lastUpdatedDate !== null;
        return r.id.toLowerCase().includes(q) || r.createdBy.toLowerCase().includes(q) || r.createdDate.toLowerCase().includes(q);
      });
    }
    return list;
  });

  function handleReload() {
    appState.consoleLogs.push({
      id: 'log-' + Date.now(),
      timestamp: '2026-10-03 10:50:34',
      type: 'query',
      summary: '58 rows retrieved starting from 1 in 586 ms (execution: 86 ms, fetching: 500 ms)'
    });
  }
</script>

<div class="table-data-view">
  <!-- 1. TOP TOOLBAR -->
  <div class="table-toolbar">
    <div class="toolbar-left">
      <!-- ⟳ Reload -->
      <button type="button" class="jb-icon-btn" title="Reload (Cmd+R)" onclick={handleReload}>
        <svg width="13" height="13" viewBox="0 0 16 16" fill="currentColor">
          <path fill-rule="evenodd" d="M8 3a5 5 0 104.546 2.914.5.5 0 01.908-.417A6 6 0 118 2v1z"/>
          <path d="M8 4.466V.534a.25.25 0 01.41-.192l2.36 1.966c.12.1.12.284 0 .384L8.41 4.658A.25.25 0 018 4.466z"/>
        </svg>
      </button>

      <!-- 🕒 Clock -->
      <button type="button" class="jb-icon-btn" title="Query History">
        <svg width="13" height="13" viewBox="0 0 16 16" fill="currentColor">
          <path d="M8 1a7 7 0 100 14A7 7 0 008 1zm0 1.5a5.5 5.5 0 110 11 5.5 5.5 0 010-11zM7.25 4v4.25l3.25 1.95.75-1.23-2.5-1.5V4h-1.5z"/>
        </svg>
      </button>

      <!-- ⏹ Stop -->
      <button type="button" class="jb-icon-btn" title="Cancel Query">
        <svg width="11" height="11" viewBox="0 0 16 16" fill="currentColor">
          <rect x="2" y="2" width="12" height="12" rx="1.5"/>
        </svg>
      </button>

      <span class="bar-separator"></span>

      <!-- + Add Row -->
      <button type="button" class="jb-icon-btn" title="Add New Row">+</button>

      <!-- - Delete Row -->
      <button type="button" class="jb-icon-btn" title="Delete Row">—</button>

      <!-- ↩ Revert -->
      <button type="button" class="jb-icon-btn" title="Revert">↩</button>

      <!-- ↪ Submit -->
      <button type="button" class="jb-icon-btn" title="Submit">↪</button>

      <span class="bar-separator"></span>

      <!-- ↑ Reorder Up -->
      <button type="button" class="jb-icon-btn" title="Move Up">↑</button>

      <!-- ↓ Reorder Down -->
      <button type="button" class="jb-icon-btn" title="Move Down">↓</button>

      <span class="bar-separator"></span>

      <!-- Tx: Auto ⌵ -->
      <button type="button" class="tx-selector" title="Transaction Isolation Mode">
        <span>Tx: Auto</span>
        <svg width="8" height="8" viewBox="0 0 16 16" fill="currentColor"><path d="M4 6l4 4 4-4H4z" /></svg>
      </button>

      <span class="bar-separator"></span>

      <!-- DDL -->
      <button type="button" class="jb-icon-btn ddl-btn" title="Generate Table DDL" onclick={() => appState.ddlModalOpen = true}>
        <span style="font-size: 10px; font-weight: 700; font-family: var(--font-code); color: #7A7E85;">DDL</span>
      </button>

      <!-- Search in Table -->
      <button type="button" class="jb-icon-btn" title="Search in Table (Cmd+F)">
        <svg width="12" height="12" viewBox="0 0 16 16" fill="currentColor">
          <path fill-rule="evenodd" d="M11.5 7a4.5 4.5 0 11-9 0 4.5 4.5 0 019 0zm-.82 4.74a6 6 0 111.06-1.06l3.04 3.04a.75.75 0 11-1.06 1.06l-3.04-3.04z"/>
        </svg>
      </button>

      <!-- Columns icon -->
      <button type="button" class="jb-icon-btn" title="Show/Hide Columns">
        <svg width="12" height="12" viewBox="0 0 16 16" fill="currentColor">
          <path d="M0 2a1 1 0 011-1h14a1 1 0 011 1v12a1 1 0 01-1 1H1a1 1 0 01-1-1V2zm1 3v2h6V5H1zm7 0v2h7V5H8zm0 3v2h7V8H8zm-1 0H1v2h6V8zm0 3H1v2h6v-2zm1 0v2h7v-2H8z"/>
        </svg>
      </button>

      <!-- Diagram icon -->
      <button type="button" class="jb-icon-btn" title="View Diagram">
        <svg width="12" height="12" viewBox="0 0 16 16" fill="currentColor">
          <path d="M2 3h4v4H2V3zm0 6h4v4H2V9zm8-6h4v4h-4V3zm0 6h4v4h-4V9z"/>
        </svg>
      </button>
    </div>

    <div class="toolbar-right">
      <!-- CSV ⌵ -->
      <button type="button" class="export-dropdown" title="Data Format">
        <span>CSV</span>
        <svg width="8" height="8" viewBox="0 0 16 16" fill="currentColor"><path d="M4 6l4 4 4-4H4z" /></svg>
      </button>

      <!-- Download ⤓ -->
      <button type="button" class="jb-icon-btn" title="Export to File">⤓</button>

      <!-- Upload ⤒ -->
      <button type="button" class="jb-icon-btn" title="Import from File">⤒</button>

      <!-- Maximize ↗ -->
      <button type="button" class="jb-icon-btn" title="Maximize View">↗</button>

      <!-- Eye 👁️ -->
      <button type="button" class="jb-icon-btn" title="View Options">👁️</button>

      <!-- Settings ⚙️ -->
      <button type="button" class="jb-icon-btn" title="Settings">⚙️</button>
    </div>
  </div>

  <!-- 2. INLINE SQL FILTER BAR (WHERE / ORDER BY) -->
  <div class="filter-bar">
    <!-- WHERE Filter Field -->
    <div class="filter-group where-group">
      <span class="filter-icon">⌵</span>
      <span class="filter-label">WHERE</span>
      <input
        type="text"
        class="filter-input code-text"
        bind:value={appState.whereFilter}
      />
    </div>

    <!-- ORDER BY Filter Field -->
    <div class="filter-group orderby-group">
      <span class="filter-icon">≡</span>
      <span class="filter-label">ORDER BY</span>
      <input
        type="text"
        class="filter-input code-text"
        bind:value={appState.orderByFilter}
      />
    </div>
  </div>

  <!-- 3. DATA GRID -->
  <div class="grid-section">
    <table class="jb-table code-text">
      <thead>
        <tr>
          <th class="row-num-header"></th>
          <!-- Column id -->
          <th class="col-header" style="width: 280px;">
            <div class="header-inner">
              <svg width="12" height="12" viewBox="0 0 16 16" fill="#FACC15">
                <path d="M0 8a4 4 0 017.465-2H14a2 2 0 012 2v1a1 1 0 01-1 1h-1v1a1 1 0 01-1 1h-1v1a1 1 0 01-1 1H9.465A4 4 0 010 8zm4-2a2 2 0 100 4 2 2 0 000-4z"/>
              </svg>
              <span class="col-name">id</span>
              <span class="header-funnel">▽</span>
              <span class="header-sort">⇅</span>
            </div>
          </th>

          <!-- Column createdDate -->
          <th class="col-header" style="width: 220px;">
            <div class="header-inner">
              <svg width="12" height="12" viewBox="0 0 16 16" fill="#56A8F5">
                <path d="M3.5 0a.5.5 0 01.5.5V1h8V.5a.5.5 0 011 0V1h1a2 2 0 012 2v11a2 2 0 01-2 2H2a2 2 0 01-2-2V3a2 2 0 012-2h1V.5a.5.5 0 01.5-.5zM1 4v10a1 1 0 001 1h12a1 1 0 001-1V4H1z"/>
              </svg>
              <span class="col-name">createdDate</span>
              <span class="header-funnel">▽</span>
              <span class="header-sort">⇅</span>
            </div>
          </th>

          <!-- Column lastUpdatedDate -->
          <th class="col-header" style="width: 220px;">
            <div class="header-inner">
              <svg width="12" height="12" viewBox="0 0 16 16" fill="#56A8F5">
                <path d="M3.5 0a.5.5 0 01.5.5V1h8V.5a.5.5 0 011 0V1h1a2 2 0 012 2v11a2 2 0 01-2 2H2a2 2 0 01-2-2V3a2 2 0 012-2h1V.5a.5.5 0 01.5-.5zM1 4v10a1 1 0 001 1h12a1 1 0 001-1V4H1z"/>
              </svg>
              <span class="col-name">lastUpdatedDate</span>
              <span class="header-funnel">▽</span>
              <span class="header-sort">⇅</span>
            </div>
          </th>

          <!-- Column createdBy -->
          <th class="col-header">
            <div class="header-inner">
              <svg width="12" height="12" viewBox="0 0 16 16" fill="#7A7E85">
                <path d="M8 8a3 3 0 100-6 3 3 0 000 6zm2-3a2 2 0 11-4 0 2 2 0 014 0zm4 8c0 1-1 1-1 1H3s-1 0-1-1 1-4 6-4 6 3 6 4zm-1-.004c-.001-.246-.154-.986-.832-1.664C11.516 10.68 10.289 10 8 10c-2.29 0-3.516.68-4.168 1.332-.678.678-.83 1.418-.832 1.664h10z"/>
              </svg>
              <span class="col-name">createdBy</span>
              <span class="header-funnel">▽</span>
            </div>
          </th>
        </tr>
      </thead>
      <tbody>
        {#each displayRows as row, idx}
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
              class="cell date-cell"
              class:cell-focused={selectedRowIndex === idx && selectedColumn === 'createdDate'}
              onclick={() => selectedColumn = 'createdDate'}
            >
              {row.createdDate}
            </td>
            <td
              class="cell date-cell"
              class:null-cell={row.lastUpdatedDate === null}
              class:cell-focused={selectedRowIndex === idx && selectedColumn === 'lastUpdatedDate'}
              onclick={() => selectedColumn = 'lastUpdatedDate'}
            >
              {row.lastUpdatedDate ?? '<null>'}
            </td>
            <td
              class="cell uuid-cell"
              class:cell-focused={selectedRowIndex === idx && selectedColumn === 'createdBy'}
              onclick={() => selectedColumn = 'createdBy'}
            >
              {row.createdBy}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>

    <!-- Floating row count pill [ 58 rows ⌵ | ⋮ ] -->
    <div class="floating-row-badge" title="Filter count">
      <span>58 rows ⌵</span>
      <span style="color: var(--text-muted); opacity: 0.6;">|</span>
      <span>⋮</span>
    </div>
  </div>
</div>

<style>
  .table-data-view {
    display: flex;
    flex-direction: column;
    height: 100%;
    background-color: var(--bg-canvas);
    overflow: hidden;
  }

  /* 1. Top toolbar */
  .table-toolbar {
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

  .bar-separator {
    width: 1px;
    height: 12px;
    background-color: var(--border-default);
    margin: 0 3px;
  }

  .tx-selector, .export-dropdown {
    display: flex;
    align-items: center;
    gap: 4px;
    height: var(--control-height-xs, 24px);
    font-size: var(--font-size-xs, 11px);
    color: var(--text-secondary);
    padding: 0 6px;
    border-radius: var(--radius-sm, 4px);
  }

  .tx-selector:hover, .export-dropdown:hover {
    background-color: var(--bg-hover);
    color: var(--text-primary);
  }

  /* 2. Inline Filter Bar */
  .filter-bar {
    height: var(--filterbar-height);
    background-color: var(--bg-canvas);
    border-bottom: 1px solid var(--border-subtle);
    display: flex;
    align-items: center;
    padding: 0 8px;
    gap: 12px;
    flex-shrink: 0;
  }

  .filter-group {
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .where-group {
    flex: 1.1;
  }

  .orderby-group {
    flex: 0.9;
  }

  .filter-icon {
    color: var(--text-muted);
    font-size: 11px;
    user-select: none;
  }

  .filter-label {
    font-size: 11px;
    font-weight: 600;
    color: var(--text-secondary);
    user-select: none;
  }

  .filter-input {
    flex: 1;
    height: var(--control-height-xs, 24px);
    background-color: #1E1F22;
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-sm, 4px);
    padding: 0 8px;
    font-size: var(--font-size-sm, 12px);
    color: var(--text-primary);
  }

  .filter-input:focus {
    border-color: var(--border-accent);
  }

  /* 3. Data Grid */
  .grid-section {
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
    font-size: 11px;
    padding: 0 6px;
    text-align: left;
    user-select: none;
    white-space: nowrap;
  }

  .header-inner {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .col-name {
    color: var(--text-primary);
    font-weight: 500;
  }

  .header-funnel, .header-sort {
    color: var(--text-muted);
    font-size: 9px;
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
    font-size: 11px;
    user-select: none;
  }

  .cell {
    height: var(--table-row-height);
    border-bottom: 1px solid #25272A;
    border-right: 1px solid var(--border-subtle);
    padding: 0 6px;
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

  .uuid-cell, .date-cell {
    font-size: 12px;
  }

  .null-cell {
    color: var(--text-null);
    font-style: italic;
  }
</style>
