<script lang="ts">
  import { appState } from '../../state/appState.svelte';
  import { Icon } from '@ohjanus/icons';
  import { Alert, Button, Select, Text, toast } from '@ohjanus/ui';
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
  let txMode = $state('auto');
  let exportFormat = $state('CSV');

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

  function handleExport() {
    if (viewer.rows.length === 0 || viewer.columns.length === 0) {
      toast.info('Export', 'No data to export.');
      return;
    }
    if (exportFormat === 'JSON') {
      const objects = viewer.rows.map((row) => {
        const obj: Record<string, unknown> = {};
        viewer.columns.forEach((col, idx) => {
          obj[col] = row[idx];
        });
        return obj;
      });
      const dataStr = 'data:text/json;charset=utf-8,' + encodeURIComponent(JSON.stringify(objects, null, 2));
      const downloadAnchor = document.createElement('a');
      downloadAnchor.setAttribute('href', dataStr);
      downloadAnchor.setAttribute('download', `${viewer.table || 'table'}_${Date.now()}.json`);
      downloadAnchor.click();
      toast.success('Data Exported', `Exported ${viewer.rowCount} rows to JSON`);
    } else {
      const headers = viewer.columns.map((c) => `"${c.replace(/"/g, '""')}"`).join(',');
      const rows = viewer.rows.map((row) =>
        row.map((cell) => cell === null || cell === undefined ? '' : `"${String(cell).replace(/"/g, '""')}"`).join(',')
      ).join('\n');
      const csvContent = [headers, ...rows].join('\n');
      const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.setAttribute('href', url);
      link.setAttribute('download', `${viewer.table || 'table'}_${Date.now()}.csv`);
      link.click();
      URL.revokeObjectURL(url);
      toast.success('Data Exported', `Exported ${viewer.rowCount} rows to CSV`);
    }
  }
</script>

<div class="table-data-view">
  <div class="table-toolbar">
    <div class="toolbar-left">
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Reload (Cmd+Enter)"
        onclick={handleReload}
        disabled={viewer.loading}
      >
        <Icon name="refresh" size={13} />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Audit history"
        onclick={() => (appState.activeTabId = "audit")}
      >
        <Icon name="clock" size={13} />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Stop Execution"
        disabled
      >
        <Icon name="stop" size={13} color="#E55353" />
      </Button>
      <span class="bar-separator"></span>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Add Row (+)"
        disabled
      >
        <Icon name="plus" size={13} color="#57D38C" />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Delete Row (-)"
        disabled
      >
        <Icon name="minus" size={13} />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Undo"
        disabled
      >
        <Icon name="undo" size={13} />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Redo"
        disabled
      >
        <Icon name="redo" size={13} />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Submit Changes"
        disabled
      >
        <Icon name="arrow-up" size={13} />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Rollback Changes"
        disabled
      >
        <Icon name="arrow-down" size={13} />
      </Button>
      <span class="bar-separator"></span>
      <div class="borderless-select-wrapper">
        <Select
          class="toolbar-select borderless-select"
          options={[
            { value: "auto", label: "Tx: Auto" },
            { value: "manual", label: "Tx: Manual" },
          ]}
          bind:value={txMode}
        />
      </div>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn ddl-btn"
        title="Generate Table DDL"
        onclick={() => (appState.ddlModalOpen = true)}
      >
        <Text size="xs" weight="bold" color="muted" mono>DDL</Text>
      </Button>
      <span class="bar-separator"></span>
      <span class="tx-selector" title="Target table">
        {#if viewer.connection && viewer.table}
          {viewer.connection}.{viewer.schema}.{viewer.table}
        {:else}
          No table selected
        {/if}
      </span>
      {#if viewer.loading}
        <Text size="sm" color="muted" mono style="padding: 0 6px;">Loading…</Text>
      {:else if viewer.durationMs > 0}
        <Text size="sm" color="muted" mono style="padding: 0 6px;">{viewer.rowCount} row(s){viewer.truncated ? ' (truncated)' : ''} · {viewer.durationMs} ms</Text>
      {/if}
    </div>

    <div class="toolbar-right">
      <div class="borderless-select-wrapper format-select-wrapper">
        <Select
          class="toolbar-select borderless-select"
          options={[
            { value: "CSV", label: "CSV" },
            { value: "JSON", label: "JSON" },
          ]}
          bind:value={exportFormat}
        />
      </div>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Export Data"
        onclick={handleExport}
      >
        <Icon name="download" size={13} />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Import Data"
        disabled
      >
        <Icon name="upload" size={13} />
      </Button>
      <span class="bar-separator"></span>
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
      <span class="bar-separator"></span>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Settings"
        onclick={() => (appState.settingsModalOpen = true)}
      >
        <Icon name="settings" size={13} />
      </Button>
    </div>
  </div>

  <div class="filter-bar">
    <div class="filter-group where-group">
      <span class="filter-icon"><Icon name="filter" size={11} /></span>
      <Text size="sm" weight="semibold" color="secondary" style="user-select: none;">WHERE</Text>
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
      <Text size="sm" weight="semibold" color="secondary" style="user-select: none;">ORDER BY</Text>
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
          <Text size="sm">{viewer.rowCount} rows{viewer.truncated ? '+' : ''}</Text>
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
    height: 14px;
    background-color: var(--border-default);
    margin: 0 4px;
  }

  .borderless-select-wrapper {
    display: inline-flex;
    align-items: center;
  }

  :global(.table-data-view .borderless-select) {
    min-width: unset !important;
    width: auto !important;
  }

  :global(.table-data-view .borderless-select .ohjanus-select-trigger) {
    background-color: transparent !important;
    border: none !important;
    box-shadow: none !important;
    height: 24px !important;
    padding: 0 6px !important;
    gap: 4px !important;
    font-size: var(--font-size-xs, 12px) !important;
    color: var(--text-secondary, #9DA0A8) !important;
    cursor: pointer;
  }

  :global(.table-data-view .borderless-select .ohjanus-select-trigger:hover) {
    background-color: var(--bg-hover, #313438) !important;
    color: var(--text-primary, #DFE1E5) !important;
  }

  .tx-selector {
    font-size: var(--font-size-xs, 12px);
    color: var(--text-primary);
    padding: 0 6px;
    font-family: var(--font-code);
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

  .filter-input {
    flex: 1;
    height: var(--control-height-xs, 24px);
    background-color: transparent;
    border: none;
    outline: none;
    padding: 0 8px;
    font-size: var(--font-size-sm, 13px);
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
