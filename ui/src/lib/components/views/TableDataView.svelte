<script lang="ts">
  import { appState } from '../../state/appState.svelte';
  import { Icon } from '@ohjanus/icons';
  import { Alert, Button, Text, Box, Flex, Input, toast } from '@ohjanus/ui';
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
  import {
    Toolbar,
    ToolbarSeparator,
    BorderlessSelect,
    FilterBar,
    FloatingRowCount,
  } from '../toolbar';

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

<Box class="table-data-view">
  <!-- Toolbar 1: Actions -->
  <Toolbar>
    {#snippet left()}
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
      <ToolbarSeparator />
      <BorderlessSelect
        options={[
          { value: "auto", label: "Tx: Auto" },
          { value: "manual", label: "Tx: Manual" },
        ]}
        bind:value={txMode}
      />
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn ddl-btn"
        title="Generate Table DDL"
        onclick={() => (appState.ddlModalOpen = true)}
      >
        <Text size="xs" weight="bold" color="muted" mono>DDL</Text>
      </Button>
      <ToolbarSeparator />
      <Text size="xs" color="muted" class="tx-selector" title="Target table">
        {#if viewer.connection && viewer.table}
          {viewer.connection}.{viewer.schema}.{viewer.table}
        {:else}
          No table selected
        {/if}
      </Text>
      {#if viewer.loading}
        <Text size="sm" color="muted" mono style="padding: 0 6px;">Loading…</Text>
      {:else if viewer.durationMs > 0}
        <Text size="sm" color="muted" mono style="padding: 0 6px;">{viewer.rowCount} row(s){viewer.truncated ? ' (truncated)' : ''} · {viewer.durationMs} ms</Text>
      {/if}
    {/snippet}

    {#snippet right()}
      <BorderlessSelect
        options={[
          { value: "CSV", label: "CSV" },
          { value: "JSON", label: "JSON" },
        ]}
        bind:value={exportFormat}
      />
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Export Data"
        onclick={handleExport}
      >
        <Icon name="download" size={13} />
      </Button>
      <ToolbarSeparator />
      <Flex align="center" gap="4px" class="limit-wrapper">
        <Text size="xs" color="muted" title="Row limit">Limit</Text>
        <Input
          type="number"
          size="sm"
          class="limit-input code-text"
          bind:value={viewer.limit}
          min={1}
          max={1000}
          onchange={handleReload}
        />
      </Flex>
      <ToolbarSeparator />
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Settings"
        onclick={() => (appState.settingsModalOpen = true)}
      >
        <Icon name="settings" size={13} />
      </Button>
    {/snippet}
  </Toolbar>

  <!-- Toolbar 2: WHERE & ORDER BY -->
  <FilterBar
    wherePlaceholder="e.g. id IS NOT NULL"
    orderByPlaceholder="e.g. id DESC"
    bind:whereValue={viewer.where}
    bind:orderByValue={viewer.orderBy}
    onwherechange={handleReload}
    onorderbychange={handleReload}
  />

  {#if viewer.error}
    <Box class="table-alert">
      <Alert variant="danger" title="Table load failed">
        <Text class="code-text">{viewer.error}</Text>
      </Alert>
    </Box>
  {/if}

  <DataGrid style="flex: 1; min-height: 0;">
    {#snippet overlay()}
      <FloatingRowCount count={viewer.rowCount} unit="rows{viewer.truncated ? '+' : ''}" />
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
              <Text size="xs">{col}</Text>
              <Button
                variant="ghost"
                size="icon-xs"
                class="ohjanus-data-grid-header-action"
                title="Sort by {col}"
                onclick={() => toggleSort(`"${col}"`)}
              >
                <Icon name="sort" size={9} />
              </Button>
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
</Box>

<style>
  :global(.table-data-view) {
    display: flex;
    flex-direction: column;
    height: 100%;
    background-color: var(--bg-canvas);
    color: var(--text-primary);
    overflow: hidden;
  }

  :global(.table-data-view .ddl-btn) {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 0 4px !important;
    width: auto !important;
  }

  :global(.table-data-view .tx-selector) {
    font-family: var(--font-code);
    padding: 0 4px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 220px;
  }

  :global(.table-data-view .limit-wrapper) {
    user-select: none;
  }

  :global(.table-data-view .limit-input) {
    width: 48px !important;
    height: 22px !important;
    text-align: right;
  }

  :global(.table-data-view .limit-input .ohjanus-input-field) {
    text-align: right;
    padding: 0 4px !important;
  }

  :global(.table-alert) {
    padding: 8px;
    border-bottom: 1px solid var(--border-subtle);
  }
</style>
