<script lang="ts">
  import { appState } from '../../state/appState.svelte';
  import { Icon } from '@ohjanus/icons';
  import { Alert, Select, Text } from '@ohjanus/ui';
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

  function handleExecute() {
    void appState.executeConsoleQuery();
  }

  function handleExplain() {
    void appState.explainConsoleQuery();
  }

  function cellText(v: unknown): string {
    if (v === null || v === undefined) return '<null>';
    if (typeof v === 'object') return JSON.stringify(v);
    return String(v);
  }

  function isNullCell(v: unknown): boolean {
    return v === null || v === undefined;
  }
</script>

<div class="sql-console-view">
  <div class="sql-toolbar">
    <div class="toolbar-left">
      <button
        type="button"
        class="action-btn run-btn"
        title="Execute (Cmd+Enter)"
        onclick={handleExecute}
        disabled={appState.console.isExecuting || !appState.console.connection}
      >
        <Icon name="play" size={13} />
      </button>
      <button type="button" class="action-btn" title="EXPLAIN (never executes)" onclick={handleExplain}>
        <Text size="sm" weight="semibold" color="secondary">EX</Text>
      </button>
      <button
        type="button"
        class="action-btn"
        title="Reload last query"
        onclick={handleExecute}
      >
        <Icon name="refresh" size={12} />
      </button>
      <span class="bar-separator"></span>
      <Text size="sm" color="secondary" style="padding: 0 4px;">
        {#if appState.console.isExecuting}
          Executing…
        {:else if appState.console.durationMs > 0}
          {appState.console.durationMs} ms
        {:else}
          Tx: Auto
        {/if}
      </Text>
      <span class="bar-separator"></span>
      <button type="button" class="jb-icon-btn" title="View DDL" onclick={() => appState.ddlModalOpen = true}>
        <Text size="xs" weight="bold" color="muted">DDL</Text>
      </button>
    </div>

    <div class="toolbar-right">
      {#if appState.connections.length === 0}
        <span class="schema-name">No connection</span>
      {:else}
        <Select
          options={appState.connections.map((conn) => ({
            value: conn.name,
            label: conn.readonly ? `${conn.name} (readonly)` : conn.name
          }))}
          bind:value={appState.console.connection}
        />
      {/if}
    </div>
  </div>

  <div class="editor-container">
    {#if appState.connections.length === 0}
      <div class="editor-empty">
        {#if appState.dataLoading}
          Connecting to Admin API…
        {:else}
          No connections available. Check janus.yaml connections and Admin API status.
        {/if}
      </div>
    {:else}
      <textarea
        class="editor-textarea code-text"
        bind:value={appState.console.sql}
        spellcheck={false}
        placeholder="-- SELECT * FROM ..."
      ></textarea>
    {/if}
  </div>

  {#if appState.console.error}
    <div class="console-alert">
      <Alert variant="danger" title="Query failed">
        <p class="code-text">{appState.console.error}</p>
      </Alert>
    </div>
  {/if}

  {#if appState.console.plan}
    <div class="console-alert">
      <Alert variant="info" title="EXPLAIN plan">
        <pre class="code-text plan-pre">{appState.console.plan}</pre>
      </Alert>
    </div>
  {/if}

  <div class="results-header-tabs">
    <div class="result-tab active">
      <Icon name="table" size={12} />
      <Text size="sm" color="secondary">
        {#if appState.console.rowCount > 0}
          Result · {appState.console.rowCount} row(s){appState.console.truncated ? ' (truncated)' : ''}
        {:else}
          Result
        {/if}
      </Text>
    </div>
  </div>

  <DataGrid style="flex: 1; min-height: 0;">
    {#snippet overlay()}
      {#if appState.console.rowCount > 0}
        <div class="floating-row-badge" title="Retrieved count">
          <Text size="sm">{appState.console.rowCount} row(s)</Text>
        </div>
      {/if}
    {/snippet}
    <DataGridHead>
      <DataGridRow>
        <DataGridRowNumHead />
        {#each appState.console.columns as col, i (col + i)}
          <DataGridHeadCell width="200px">
            <DataGridHeaderInner>
              <span>{col}</span>
            </DataGridHeaderInner>
          </DataGridHeadCell>
        {/each}
      </DataGridRow>
    </DataGridHead>
    <DataGridBody>
      {#if appState.console.isExecuting}
        <DataGridRow>
          <DataGridRowNum index={0} />
          <DataGridCell>Executing…</DataGridCell>
        </DataGridRow>
      {:else if appState.console.columns.length === 0}
        <DataGridRow>
          <DataGridRowNum index={0} />
          <DataGridCell tone="secondary">No result yet — run a query above.</DataGridCell>
        </DataGridRow>
      {:else}
        {#each appState.console.rows as row, r (r)}
          <DataGridRow selected={selectedRowIndex === r} onclick={() => selectedRowIndex = r}>
            <DataGridRowNum index={r} />
            {#each row as cell, c (c)}
              <DataGridCell
                isNull={isNullCell(cell)}
                focused={selectedRowIndex === r && selectedColumn === c}
                onclick={() => { selectedRowIndex = r; selectedColumn = c; }}
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
  .sql-console-view {
    display: flex;
    flex-direction: column;
    height: 100%;
    background-color: var(--bg-canvas);
    overflow: hidden;
  }

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

  .action-btn:hover:not(:disabled) {
    background-color: var(--bg-hover);
    color: var(--text-primary);
  }

  .action-btn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }

  .action-btn.run-btn {
    color: var(--action-success);
  }

  .bar-separator {
    width: 1px;
    height: 14px;
    background-color: var(--border-default);
    margin: 0 4px;
  }

  .schema-name {
    font-size: var(--font-size-xs, 12px);
    color: var(--text-muted);
  }

  .editor-container {
    height: 220px;
    min-height: 160px;
    background-color: var(--bg-canvas);
    display: flex;
    overflow: hidden;
    position: relative;
    border-bottom: 1px solid var(--border-subtle);
  }

  .editor-textarea {
    flex: 1;
    background: transparent;
    border: none;
    outline: none;
    resize: none;
    padding: 10px 12px;
    font-size: 13px;
    line-height: 20px;
    color: var(--text-primary);
  }

  .editor-empty {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--text-muted);
    font-size: 12px;
    padding: 16px;
    text-align: center;
  }

  .console-alert {
    padding: 8px 12px;
    flex-shrink: 0;
  }

  .plan-pre {
    white-space: pre-wrap;
    font-size: 11px;
    margin: 4px 0 0;
  }

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
    user-select: none;
    white-space: nowrap;
  }

  .result-tab.active {
    background-color: #1F2E4A;
    border: 1px solid #3574F0;
    color: #DFE1E5;
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
