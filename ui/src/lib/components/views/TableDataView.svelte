<script lang="ts">
  import { appState } from '../../state/appState.svelte';
  import { Icon } from '@ohjanus/icons';

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
        <Icon name="refresh" size={13} />
      </button>

      <!-- 🕒 Clock -->
      <button type="button" class="jb-icon-btn" title="Query History">
        <Icon name="clock" size={13} />
      </button>

      <!-- ⏹ Stop -->
      <button type="button" class="jb-icon-btn" title="Cancel Query">
        <Icon name="stop" size={11} />
      </button>

      <span class="bar-separator"></span>

      <!-- + Add Row -->
      <button type="button" class="jb-icon-btn" title="Add New Row">
        <Icon name="plus" size={13} />
      </button>

      <!-- - Delete Row -->
      <button type="button" class="jb-icon-btn" title="Delete Row">
        <Icon name="minus" size={13} />
      </button>

      <!-- ↩ Revert -->
      <button type="button" class="jb-icon-btn" title="Revert">
        <Icon name="undo" size={13} />
      </button>

      <!-- ↪ Submit -->
      <button type="button" class="jb-icon-btn" title="Submit">
        <Icon name="redo" size={13} />
      </button>

      <span class="bar-separator"></span>

      <!-- ↑ Reorder Up -->
      <button type="button" class="jb-icon-btn" title="Move Up">
        <Icon name="arrow-up" size={13} />
      </button>

      <!-- ↓ Reorder Down -->
      <button type="button" class="jb-icon-btn" title="Move Down">
        <Icon name="arrow-down" size={13} />
      </button>

      <span class="bar-separator"></span>

      <!-- Tx: Auto ⌵ -->
      <button type="button" class="tx-selector" title="Transaction Isolation Mode">
        <span>Tx: Auto</span>
        <Icon name="chevron-down" size={8} />
      </button>

      <span class="bar-separator"></span>

      <!-- DDL -->
      <button type="button" class="jb-icon-btn ddl-btn" title="Generate Table DDL" onclick={() => appState.ddlModalOpen = true}>
        <span style="font-size: 10px; font-weight: 700; font-family: var(--font-code); color: #7A7E85;">DDL</span>
      </button>

      <!-- Search in Table -->
      <button type="button" class="jb-icon-btn" title="Search in Table (Cmd+F)">
        <Icon name="search" size={12} />
      </button>

      <!-- Columns icon -->
      <button type="button" class="jb-icon-btn" title="Show/Hide Columns">
        <Icon name="table" size={12} />
      </button>

      <!-- Diagram icon -->
      <button type="button" class="jb-icon-btn" title="View Diagram">
        <Icon name="layout" size={12} />
      </button>
    </div>

    <div class="toolbar-right">
      <!-- CSV ⌵ -->
      <button type="button" class="export-dropdown" title="Data Format">
        <span>CSV</span>
        <Icon name="chevron-down" size={8} />
      </button>

      <!-- Download ⤓ -->
      <button type="button" class="jb-icon-btn" title="Export to File">
        <Icon name="download" size={12} />
      </button>

      <!-- Upload ⤒ -->
      <button type="button" class="jb-icon-btn" title="Import from File">
        <Icon name="upload" size={12} />
      </button>

      <!-- Maximize ↗ -->
      <button type="button" class="jb-icon-btn" title="Maximize View">
        <Icon name="maximize" size={12} />
      </button>

      <!-- Eye 👁️ -->
      <button type="button" class="jb-icon-btn" title="View Options">
        <Icon name="eye" size={12} />
      </button>

      <!-- Settings ⚙️ -->
      <button type="button" class="jb-icon-btn" title="Settings">
        <Icon name="settings" size={12} />
      </button>
    </div>
  </div>

  <!-- 2. INLINE SQL FILTER BAR (WHERE / ORDER BY) -->
  <div class="filter-bar">
    <!-- WHERE Filter Field -->
    <div class="filter-group where-group">
      <span class="filter-icon"><Icon name="filter" size={11} /></span>
      <span class="filter-label">WHERE</span>
      <input
        type="text"
        class="filter-input code-text"
        bind:value={appState.whereFilter}
      />
    </div>

    <!-- ORDER BY Filter Field -->
    <div class="filter-group orderby-group">
      <span class="filter-icon"><Icon name="sort" size={11} /></span>
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
              <Icon name="key" size={12} color="#FACC15" />
              <span class="col-name">id</span>
              <span class="header-funnel"><Icon name="filter" size={9} /></span>
              <span class="header-sort"><Icon name="sort" size={9} /></span>
            </div>
          </th>

          <!-- Column createdDate -->
          <th class="col-header" style="width: 220px;">
            <div class="header-inner">
              <Icon name="clock" size={12} color="#56A8F5" />
              <span class="col-name">createdDate</span>
              <span class="header-funnel"><Icon name="filter" size={9} /></span>
              <span class="header-sort"><Icon name="sort" size={9} /></span>
            </div>
          </th>

          <!-- Column lastUpdatedDate -->
          <th class="col-header" style="width: 220px;">
            <div class="header-inner">
              <Icon name="clock" size={12} color="#56A8F5" />
              <span class="col-name">lastUpdatedDate</span>
              <span class="header-funnel"><Icon name="filter" size={9} /></span>
              <span class="header-sort"><Icon name="sort" size={9} /></span>
            </div>
          </th>

          <!-- Column createdBy -->
          <th class="col-header">
            <div class="header-inner">
              <Icon name="user" size={12} color="#7A7E85" />
              <span class="col-name">createdBy</span>
              <span class="header-funnel"><Icon name="filter" size={9} /></span>
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
      <span>58 rows</span>
      <span style="color: var(--text-muted); opacity: 0.6;">|</span>
      <Icon name="chevron-down" size={10} />
      <span style="color: var(--text-muted); opacity: 0.6;">|</span>
      <Icon name="more" size={12} />
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
