<script lang="ts">
  import { connectionsState } from '@/features/connections';
  import { consoleState } from '@/features/console';
  import { explorerState } from '@/features/explorer';
  import { bootState } from '@/app/boot.svelte';
  import { Icon } from '@ohjanus/icons';
  import { Alert, Button, Text, Box, Flex, Textarea } from '@ohjanus/ui';
  import { Toolbar, ToolbarSeparator, BorderlessSelect, FloatingRowCount } from '@/shared/ui/toolbar';
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
    void consoleState.executeConsoleQuery();
  }

  function handleExplain() {
    void consoleState.explainConsoleQuery();
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

<Box class="sql-console-view">
  <Toolbar>
    {#snippet left()}
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn run-btn"
        title="Execute (Cmd+Enter)"
        onclick={handleExecute}
        disabled={consoleState.console.isExecuting || !consoleState.console.connection}
      >
        <Icon name="play" size={13} color="#57D38C" />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="EXPLAIN (never executes)"
        onclick={handleExplain}
      >
        <Text size="xs" weight="semibold" color="secondary">EX</Text>
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Reload last query"
        onclick={handleExecute}
      >
        <Icon name="refresh" size={12} />
      </Button>
      <ToolbarSeparator />
      <Text size="sm" color="secondary" style="padding: 0 4px;">
        {#if consoleState.console.isExecuting}
          Executing…
        {:else if consoleState.console.durationMs > 0}
          {consoleState.console.durationMs} ms
        {:else}
          Tx: Auto
        {/if}
      </Text>
      <ToolbarSeparator />
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="View DDL"
        onclick={() => explorerState.ddlModalOpen = true}
      >
        <Text size="xs" weight="bold" color="muted">DDL</Text>
      </Button>
    {/snippet}

    {#snippet right()}
      {#if connectionsState.connections.length === 0}
        <Text size="xs" color="muted" class="schema-name">No connection</Text>
      {:else}
        <BorderlessSelect
          options={connectionsState.connections.map((conn) => ({
            value: conn.name,
            label: conn.readonly ? `${conn.name} (readonly)` : conn.name
          }))}
          bind:value={consoleState.console.connection}
        />
      {/if}
    {/snippet}
  </Toolbar>

  <Box class="editor-container">
    {#if connectionsState.connections.length === 0}
      <Flex align="center" justify="center" class="editor-empty">
        <Text size="sm" color="muted">
          {#if bootState.dataLoading}
            Connecting to Admin API…
          {:else}
            No connections available. Check janus.yaml connections and Admin API status.
          {/if}
        </Text>
      </Flex>
    {:else}
      <Textarea
        class="editor-textarea code-text"
        bind:value={consoleState.console.sql}
        spellcheck={false}
        placeholder="-- SELECT * FROM ..."
      />
    {/if}
  </Box>

  {#if consoleState.console.error}
    <Box class="console-alert">
      <Alert variant="danger" title="Query failed">
        <Text class="code-text">{consoleState.console.error}</Text>
      </Alert>
    </Box>
  {/if}

  {#if consoleState.console.plan}
    <Box class="console-alert">
      <Alert variant="info" title="EXPLAIN plan">
        <Box class="code-text plan-pre">{consoleState.console.plan}</Box>
      </Alert>
    </Box>
  {/if}

  <Flex align="center" class="results-header-tabs">
    <Flex align="center" gap="6px" class="result-tab active">
      <Icon name="table" size={12} />
      <Text size="sm" color="secondary">
        {#if consoleState.console.rowCount > 0}
          Result · {consoleState.console.rowCount} row(s){consoleState.console.truncated ? ' (truncated)' : ''}
        {:else}
          Result
        {/if}
      </Text>
    </Flex>
  </Flex>

  <DataGrid style="flex: 1; min-height: 0;">
    {#snippet overlay()}
      <FloatingRowCount count={consoleState.console.rowCount} />
    {/snippet}
    <DataGridHead>
      <DataGridRow>
        <DataGridRowNumHead />
        {#each consoleState.console.columns as col, i (col + i)}
          <DataGridHeadCell width="200px">
            <DataGridHeaderInner>
              <Text size="xs">{col}</Text>
            </DataGridHeaderInner>
          </DataGridHeadCell>
        {/each}
      </DataGridRow>
    </DataGridHead>
    <DataGridBody>
      {#if consoleState.console.isExecuting}
        <DataGridRow>
          <DataGridRowNum index={0} />
          <DataGridCell>Executing…</DataGridCell>
        </DataGridRow>
      {:else if consoleState.console.columns.length === 0}
        <DataGridRow>
          <DataGridRowNum index={0} />
          <DataGridCell tone="secondary">No result yet — run a query above.</DataGridCell>
        </DataGridRow>
      {:else}
        {#each consoleState.console.rows as row, r (r)}
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
</Box>

<style>
  :global(.sql-console-view) {
    display: flex;
    flex-direction: column;
    height: 100%;
    background-color: var(--bg-canvas);
    overflow: hidden;
  }

  :global(.sql-console-view .schema-name) {
    font-size: var(--font-size-xs, 12px);
    color: var(--text-muted);
    font-family: var(--font-code);
  }

  :global(.sql-console-view .editor-container) {
    height: 180px;
    background-color: var(--bg-canvas);
    border-bottom: 1px solid var(--border-subtle);
    position: relative;
    flex-shrink: 0;
  }

  :global(.sql-console-view .editor-empty) {
    width: 100%;
    height: 100%;
    padding: 16px;
    font-size: var(--font-size-sm);
  }

  :global(.sql-console-view .editor-textarea) {
    width: 100% !important;
    height: 100% !important;
    padding: 10px !important;
    background-color: transparent !important;
    border: none !important;
    color: var(--text-primary) !important;
    font-family: var(--font-code) !important;
    font-size: var(--font-size-sm) !important;
    resize: none !important;
    outline: none !important;
  }

  :global(.sql-console-view .console-alert) {
    padding: 6px 8px;
    border-bottom: 1px solid var(--border-subtle);
  }

  :global(.sql-console-view .plan-pre) {
    margin: 0;
    font-family: var(--font-code);
    font-size: var(--font-size-xs, 12px);
    white-space: pre-wrap;
    max-height: 140px;
    overflow-y: auto;
  }

  :global(.sql-console-view .results-header-tabs) {
    height: 28px;
    background-color: var(--bg-toolbar);
    border-bottom: 1px solid var(--border-subtle);
    padding: 0 8px;
    gap: 4px;
    flex-shrink: 0;
  }

  :global(.sql-console-view .result-tab) {
    height: 24px;
    padding: 0 8px;
    border-radius: var(--radius-sm, 4px) var(--radius-sm, 4px) 0 0;
    user-select: none;
  }

  :global(.sql-console-view .result-tab.active) {
    background-color: var(--bg-card);
    border-bottom: 2px solid var(--border-accent, #3574F0);
  }
</style>
