<script lang="ts">
  import { appState, type AuditRecord } from '../../state/appState.svelte';
  import { Button, Badge, Input, toast, Box, Flex, Stack, Text } from '@ohjanus/ui';
  import { Icon } from '@ohjanus/icons';

  let filterStatus = $state<string>('ALL');
  let searchQuery = $state<string>('');
  let selectedAudit = $state<AuditRecord | null>(null);

  let filteredAuditLogs = $derived.by(() => {
    let list = appState.auditLogs;
    if (filterStatus !== 'ALL') {
      list = list.filter(item => item.policy_decision === filterStatus);
    }
    if (searchQuery.trim()) {
      const q = searchQuery.toLowerCase().trim();
      list = list.filter(item =>
        item.client.toLowerCase().includes(q) ||
        item.sql_normalized.toLowerCase().includes(q) ||
        item.token_id.toLowerCase().includes(q) ||
        item.connection.toLowerCase().includes(q)
      );
    }
    return list;
  });

  function exportCSV() {
    const headers = ['id', 'ts', 'client', 'token_id', 'connection', 'statement_type', 'decision', 'duration_ms', 'sql'];
    const rows = filteredAuditLogs.map(r => [
      r.id,
      r.ts,
      r.client,
      r.token_id,
      r.connection,
      r.statement_type,
      r.policy_decision,
      r.duration_ms,
      `"${r.sql_normalized.replace(/"/g, '""')}"`
    ]);

    const csvContent = 'data:text/csv;charset=utf-8,' + [headers.join(','), ...rows.map(e => e.join(','))].join('\n');
    const encodedUri = encodeURI(csvContent);
    const link = document.createElement('a');
    link.setAttribute('href', encodedUri);
    link.setAttribute('download', `ohjanus_audit_${new Date().toISOString().substring(0, 10)}.csv`);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    toast.success('Audit Log Exported', `Exported ${filteredAuditLogs.length} audit records to CSV.`);
  }
</script>

<Box class="audit-view">
  <!-- View Header & Filters -->
  <Flex as="header" class="audit-header" align="center" justify="between">
    <Stack class="header-left" gap="2px">
      <Flex class="header-title" align="center" gap="8px">
        <Icon name="audit" size={16} color="#56A8F5" />
        <span>MCP Gateway Audit Log Trail</span>
      </Flex>
      <span class="header-desc">End-to-end provenance records of all incoming queries and tool requests</span>
    </Stack>

    <Flex class="header-right" align="center" gap="8px">
      <!-- Search Input via UI Kit -->
      <Box style="width: 240px;">
        <Input
          placeholder="Filter by client, token, SQL..."
          bind:value={searchQuery}
        />
      </Box>

      <!-- Decision Filter -->
      <select class="filter-select" bind:value={filterStatus}>
        <option value="ALL">All Decisions</option>
        <option value="ALLOW">ALLOWED</option>
        <option value="REQUIRE_APPROVAL">NEEDS APPROVAL</option>
        <option value="DENY">DENIED</option>
      </select>

      <!-- Export Button via UI Kit -->
      <Button variant="secondary" size="sm" onclick={exportCSV}>
        <Icon name="download" size={14} />
        <span>Export CSV</span>
      </Button>
    </Flex>
  </Flex>

  <!-- Main Content Layout -->
  <Flex class="content-layout">
    <!-- Grid -->
    <Box class="table-container">
      <table class="audit-table code-text">
        <thead>
          <tr>
            <th style="width: 140px;">Timestamp</th>
            <th style="width: 90px;">Decision</th>
            <th style="width: 160px;">Client Agent</th>
            <th style="width: 140px;">Connection</th>
            <th style="width: 70px;">Type</th>
            <th style="width: 60px;">Rows</th>
            <th style="width: 70px;">Duration</th>
            <th>Normalized SQL</th>
          </tr>
        </thead>
        <tbody>
          {#each filteredAuditLogs as item (item.id)}
            <!-- svelte-ignore a11y_click_events_have_key_events -->
            <!-- svelte-ignore a11y_no_static_element_interactions -->
            <tr
              class:selected={selectedAudit?.id === item.id}
              onclick={() => selectedAudit = item}
            >
              <td class="cell-ts">{item.ts}</td>
              <td class="cell-decision">
                {#if item.policy_decision === 'ALLOW'}
                  <Badge variant="success" size="sm">ALLOW</Badge>
                {:else if item.policy_decision === 'REQUIRE_APPROVAL'}
                  <Badge variant="warning" size="sm">APPROVAL</Badge>
                {:else}
                  <Badge variant="danger" size="sm">DENY</Badge>
                {/if}
              </td>
              <td class="cell-client truncate" title={item.client}>
                🤖 {item.client}
              </td>
              <td class="cell-conn">{item.connection}</td>
              <td class="cell-type">{item.statement_type}</td>
              <td class="cell-num">{item.row_count}</td>
              <td class="cell-num">{item.duration_ms} ms</td>
              <td class="cell-sql truncate" title={item.sql_normalized}>
                {item.sql_normalized}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </Box>

    <!-- Right Drawer Detail -->
    {#if selectedAudit}
      <Box class="audit-drawer">
        <Flex class="drawer-header" align="center" justify="between">
          <span style="font-weight: 600; color: var(--text-primary);">Audit Detail: {selectedAudit.id}</span>
          <button class="jb-icon-btn" onclick={() => selectedAudit = null}>✕</button>
        </Flex>
        <Stack class="drawer-body" gap="8px">
          <Flex class="drawer-row" align="center" justify="between">
            <span class="d-label">Timestamp:</span>
            <span class="d-val code-text">{selectedAudit.ts}</span>
          </Flex>
          <Flex class="drawer-row" align="center" justify="between">
            <span class="d-label">Request ID:</span>
            <span class="d-val code-text">{selectedAudit.request_id}</span>
          </Flex>
          <Flex class="drawer-row" align="center" justify="between">
            <span class="d-label">Token ID:</span>
            <span class="d-val code-text">{selectedAudit.token_id}</span>
          </Flex>
          <Flex class="drawer-row" align="center" justify="between">
            <span class="d-label">Requesting Client:</span>
            <span class="d-val">{selectedAudit.client}</span>
          </Flex>
          <Flex class="drawer-row" align="center" justify="between">
            <span class="d-label">Database Target:</span>
            <span class="d-val">{selectedAudit.connection}</span>
          </Flex>
          <Flex class="drawer-row" align="center" justify="between">
            <span class="d-label">Policy Decision:</span>
            <span class="d-val"><strong>{selectedAudit.policy_decision}</strong></span>
          </Flex>
          <Flex class="drawer-row" align="center" justify="between">
            <span class="d-label">Policy Rule Matched:</span>
            <span class="d-val code-text">{selectedAudit.policy_rule}</span>
          </Flex>
          <Flex class="drawer-row" align="center" justify="between">
            <span class="d-label">Execution Duration:</span>
            <span class="d-val code-text">{selectedAudit.duration_ms} ms</span>
          </Flex>

          <Text size="xs" weight="semibold" color="secondary" style="margin-top: 8px;">
            Executed Query
          </Text>
          <pre class="sql-box code-text">{selectedAudit.sql_normalized}</pre>
        </Stack>
      </Box>
    {/if}
  </Flex>
</Box>

<style>
  :global(.audit-view) {
    display: flex;
    flex-direction: column;
    height: 100%;
    background-color: var(--bg-canvas);
    overflow: hidden;
  }

  :global(.audit-view .audit-header) {
    height: 48px;
    background-color: var(--bg-toolbar);
    border-bottom: 1px solid var(--border-default);
    padding: 0 16px;
    flex-shrink: 0;
  }

  :global(.audit-view .header-left) {
    display: flex;
    flex-direction: column;
  }

  :global(.audit-view .header-title) {
    font-size: 13px;
    font-weight: 600;
    color: var(--text-primary);
  }

  :global(.audit-view .header-desc) {
    font-size: 11px;
    color: var(--text-muted);
  }

  :global(.audit-view .filter-select) {
    height: 26px;
    background-color: #1E1F22;
    border: 1px solid var(--border-default);
    border-radius: 4px;
    padding: 0 8px;
    font-size: 11px;
    color: var(--text-primary);
  }

  :global(.audit-view .content-layout) {
    flex: 1;
    display: flex;
    flex-direction: row;
    overflow: hidden;
  }

  :global(.audit-view .table-container) {
    flex: 1;
    overflow: auto;
  }

  :global(.audit-view .audit-table) {
    width: 100%;
    border-collapse: collapse;
    font-size: 12px;
  }

  :global(.audit-view .audit-table th) {
    height: var(--table-header-height);
    background-color: var(--bg-table-header);
    border-bottom: 1px solid var(--border-default);
    border-right: 1px solid var(--border-default);
    color: var(--text-secondary);
    font-weight: 500;
    font-size: 11px;
    padding: 0 8px;
    text-align: left;
    white-space: nowrap;
    position: sticky;
    top: 0;
    z-index: 10;
  }

  :global(.audit-view .audit-table td) {
    height: 28px;
    border-bottom: 1px solid var(--border-subtle);
    border-right: 1px solid var(--border-subtle);
    padding: 0 8px;
    white-space: nowrap;
    color: var(--text-primary);
  }

  :global(.audit-view .audit-table tr:hover) {
    background-color: var(--bg-hover);
    cursor: pointer;
  }

  :global(.audit-view .audit-table tr.selected) {
    background-color: var(--bg-selected);
  }

  :global(.audit-view .cell-ts) {
    color: var(--text-muted);
    font-size: 11px;
  }

  :global(.audit-view .cell-client) {
    color: #DFE1E5;
  }

  :global(.audit-view .cell-conn) {
    color: var(--syntax-function);
  }

  :global(.audit-view .cell-type) {
    font-weight: 600;
    font-size: 11px;
  }

  :global(.audit-view .cell-num) {
    color: var(--syntax-number);
    text-align: right;
    padding-right: 12px;
  }

  :global(.audit-view .cell-sql) {
    color: #9DA0A8;
    max-width: 380px;
  }

  /* Right Drawer */
  :global(.audit-view .audit-drawer) {
    width: 360px;
    background-color: var(--bg-toolbar);
    border-left: 1px solid var(--border-default);
    display: flex;
    flex-direction: column;
    flex-shrink: 0;
  }

  :global(.audit-view .drawer-header) {
    height: 36px;
    padding: 0 12px;
    background-color: #25272A;
    border-bottom: 1px solid var(--border-default);
  }

  :global(.audit-view .drawer-body) {
    padding: 14px;
    overflow-y: auto;
  }

  :global(.audit-view .drawer-row) {
    font-size: 11px;
    padding: 4px 0;
    border-bottom: 1px solid var(--border-subtle);
  }

  :global(.audit-view .d-label) {
    color: var(--text-muted);
  }

  :global(.audit-view .d-val) {
    color: var(--text-primary);
  }

  :global(.audit-view .sql-box) {
    background-color: #1E1F22;
    border: 1px solid var(--border-default);
    padding: 8px 10px;
    border-radius: 4px;
    font-size: 11px;
    line-height: 16px;
    color: #DFE1E5;
    white-space: pre-wrap;
    max-height: 200px;
    overflow-y: auto;
  }
</style>
