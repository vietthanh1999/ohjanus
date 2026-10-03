<script lang="ts">
  import { appState } from '../../state/appState.svelte';
  import { Icon } from '@ohjanus/icons';
  import { Alert } from '@ohjanus/ui';
  import {
    DataGrid,
    DataGridHead,
    DataGridBody,
    DataGridRow,
    DataGridHeadCell,
    DataGridHeaderInner,
    DataGridRowNumHead,
    DataGridRowNum,
    DataGridCell
  } from '@ohjanus/ui';

  let selectedRowIndex = $state(0);
  let selectedColumn = $state(0);

  let viewer = $derived(appState.tableViewer);
  let tableDef = $derived(
    viewer.connection && viewer.table
      ? appState.findTable(viewer.connection, viewer.schema, viewer.table)
      : null
  );

  function handleReload() {
    void appState.loadTableData();
  }

  function cellText(v: unknown): string {
    if (v === null || v === undefined) return '<null>';
    if (typeof v === 'object') return JSON.stringify(v);
    return String(v);
  }

  function isNullCell(v: unknown): boolean {
    return v === null || v === undefined;
  }

  function toggleSort(col: string) {
    if (viewer.orderBy.startsWith(col)) {
      viewer.orderBy = viewer.orderBy.toUpperCase().includes('DESC') ? `${col} ASC` : `${col} DESC`;
    } else {
      viewer.orderBy = `${col} ASC`;
    }
    handleReload();
  }
</script>

<div class="table-data-view">
  <div class="table-toolbar">
    <div class="toolbar-left">
      <button type="button" class="jb-icon-btn" title="Reload (Cmd+Enter)" onclick={handleReload} disabled={viewer.loading}>
        <Icon name="refresh" size={13} />
      </button>
      <span class="bar-separator"></span>
      <span class="tx-selector" title="Target table">
        {#if viewer.connection && viewer.table}
          {viewer.connection}.{viewer.schema}.{viewer.table}
        {:else}
          No table selected
        {/if}
      </span>
      <span class="bar-separator"></span>
      <button type="button" class="jb-icon-btn ddl-btn" title="Generate Table DDL" onclick={() => appState.ddlModalOpen = true}>
        <span style="font-size: 10px; font-weight: 700; font-family: var(--font-code); color: #7A7E85;">DDL</span>
      </button>
      {#if viewer.loading}
        <span class="loading-text">Loading…</span>
      {:else if viewer.durationMs > 0}
        <span class="loading-text">{viewer.rowCount} row(s){viewer.truncated ? ' (truncated)' : ''} · {viewer.durationMs} ms</span>
      {/if}
    </div>

    <div class="toolbar-right">
      <label class="limit-label" title="Row limit">
        Limit
        <input
          type="number"
          class="limit-input code-text"
          bind:value={viewer.limit}
          min={1}
          max={1000}
          onchange={handleReload}
        />
      </label>
    </div>
  </div>

  <div class="filter-bar">
    <div class="filter-group where-group">
      <span class="filter-icon"><Icon name="filter" size={11} /></span>
      <span class="filter-label">WHERE</span>
      <input
        type="text"
        class="filter-input code-text"
        placeholder="e.g. id IS NOT NULL"
        bind:value={viewer.where}
        onchange={handleReload}
      />
    </div>

    <div class="filter-group orderby-group">
      <span class="filter-icon"><Icon name="sort" size={11} /></span>
      <span class="filter-label">ORDER BY</span>
      <input
        type="text"
        class="filter-input code-text"
        placeholder="e.g. id DESC"
        bind:value={viewer.orderBy}
        onchange={handleReload}
      />
    </div>
  </div>

  {#if viewer.error}
    <div class="table-alert">
      <Alert variant="danger" title="Table load failed">
        <p class="code-text">{viewer.error}</p>
      </Alert>
    </div>
  {/if}

  <DataGrid style="flex: 1; min-height: 0;">
    {#snippet overlay()}
      {#if viewer.rowCount > 0}
        <div class="floating-row-badge" title="Row count">
          <span>{viewer.rowCount} rows{viewer.truncated ? '+' : ''}</span>
        </div>
      {/if}
    {/snippet}
    <DataGridHead>
      <DataGridRow>
        <DataGridRowNumHead />
        {#each viewer.columns as col, i (col + i)}
          <DataGridHeadCell width="220px">
            <DataGridHeaderInner>
              {#if tableDef?.primaryKey.includes(col)}
                <Icon name="key" size={12} color="#FACC15" />
              {/if}
              <span>{col}</span>
              <button
                type="button"
                class="ohjanus-data-grid-header-action"
                title="Sort by {col}"
                onclick={() => toggleSort(`"${col}"`)}
              >
                <Icon name="sort" size={9} />
              </button>
            </DataGridHeaderInner>
          </DataGridHeadCell>
        {/each}
      </DataGridRow>
    </DataGridHead>
    <DataGridBody>
      {#if !viewer.connection || !viewer.table}
        <DataGridRow>
          <DataGridRowNum index={0} />
          <DataGridCell tone="secondary">Select a table in the Database Explorer to browse live data.</DataGridCell>
        </DataGridRow>
      {:else if viewer.loading}
        <DataGridRow>
          <DataGridRowNum index={0} />
          <DataGridCell tone="secondary">Loading…</DataGridCell>
        </DataGridRow>
      {:else if viewer.columns.length === 0}
        <DataGridRow>
          <DataGridRowNum index={0} />
          <DataGridCell tone="secondary">No data. Check the connection and filters above.</DataGridCell>
        </DataGridRow>
      {:else}
        {#each viewer.rows as row, idx (idx)}
          <DataGridRow selected={selectedRowIndex === idx} onclick={() => selectedRowIndex = idx}>
            <DataGridRowNum index={idx} />
            {#each row as cell, c (c)}
              <DataGridCell
                isNull={isNullCell(cell)}
                focused={selectedRowIndex === idx && selectedColumn === c}
                onclick={() => { selectedRowIndex = idx; selectedColumn = c; }}
              >
                {cellText(cell)}
              </DataGridCell>
            {/each}
          </DataGridRow>
        {/each}
      {/if}
    </DataGridBody>
  </DataGrid>

</div>

<style>
  .table-data-view {
    display: flex;
    flex-direction: column;
    height: 100%;
    background-color: var(--bg-canvas);
    overflow: hidden;
  }

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

  .tx-selector {
    font-size: var(--font-size-xs, 11px);
    color: var(--text-primary);
    padding: 0 6px;
    font-family: var(--font-code);
  }

  .loading-text {
    font-size: 11px;
    color: var(--text-muted);
    font-family: var(--font-code);
    padding: 0 6px;
  }

  .limit-label {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 11px;
    color: var(--text-secondary);
  }

  .limit-input {
    width: 72px;
    height: 24px;
    background-color: var(--bg-canvas);
    border: 1px solid var(--border-default);
    border-radius: 4px;
    color: var(--text-primary);
    padding: 0 8px;
    font-size: 12px;
    outline: none;
  }

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
    background-color: transparent;
    border: none;
    outline: none;
    padding: 0 8px;
    font-size: var(--font-size-sm, 12px);
    color: var(--text-primary);
  }

  .table-alert {
    padding: 8px 12px;
    flex-shrink: 0;
  }

  .floating-row-badge {
    position: absolute;
    bottom: 12px;
    right: 16px;
    background-color: var(--bg-toolbar);
    border: 1px solid var(--border-default);
    border-radius: 12px;
    padding: 2px 10px;
    font-size: 11px;
    color: var(--text-primary);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.45);
    z-index: 5;
  }
</style>
