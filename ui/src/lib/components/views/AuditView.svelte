<script lang="ts">
  import { appState, type AuditRecord } from '../../state/appState.svelte';

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
  }
</script>

<div class="audit-view">
  <!-- View Header & Filters -->
  <div class="audit-header">
    <div class="header-left">
      <div class="header-title">
        <svg width="15" height="15" viewBox="0 0 16 16" fill="#56A8F5">
          <path d="M4 1.5H3a2 2 0 00-2 2V14a2 2 0 002 2h10a2 2 0 002-2V3.5a2 2 0 00-2-2h-1v1h1a1 1 0 011 1V14a1 1 0 01-1 1H3a1 1 0 01-1-1V3.5a1 1 0 011-1h1v-1z"/>
        </svg>
        <span>MCP Gateway Audit Log Trail</span>
      </div>
      <span class="header-desc">End-to-end provenance records of all incoming queries and tool requests</span>
    </div>

    <div class="header-right">
      <!-- Search Input -->
      <div class="search-box">
        <svg width="12" height="12" viewBox="0 0 16 16" fill="currentColor">
          <path fill-rule="evenodd" d="M11.5 7a4.5 4.5 0 11-9 0 4.5 4.5 0 019 0zm-.82 4.74a6 6 0 111.06-1.06l3.04 3.04a.75.75 0 11-1.06 1.06l-3.04-3.04z"/>
        </svg>
        <input
          type="text"
          placeholder="Filter by client, token, SQL..."
          bind:value={searchQuery}
        />
      </div>

      <!-- Decision Filter -->
      <select class="filter-select" bind:value={filterStatus}>
        <option value="ALL">All Decisions</option>
        <option value="ALLOW">ALLOWED</option>
        <option value="REQUIRE_APPROVAL">NEEDS APPROVAL</option>
        <option value="DENY">DENIED</option>
      </select>

      <!-- Export Button -->
      <button class="export-btn" title="Export current audit dataset to CSV" onclick={exportCSV}>
        <span>⤓ Export CSV</span>
      </button>
    </div>
  </div>

  <!-- Main Content Layout -->
  <div class="content-layout">
    <!-- Grid -->
    <div class="table-container">
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
            <tr
              class:selected={selectedAudit?.id === item.id}
              onclick={() => selectedAudit = item}
            >
              <td class="cell-ts">{item.ts}</td>
              <td class="cell-decision">
                {#if item.policy_decision === 'ALLOW'}
                  <span class="decision-pill pill-allow">ALLOW</span>
                {:else if item.policy_decision === 'REQUIRE_APPROVAL'}
                  <span class="decision-pill pill-approval">APPROVAL</span>
                {:else}
                  <span class="decision-pill pill-deny">DENY</span>
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
    </div>

    <!-- Right Drawer Detail -->
    {#if selectedAudit}
      <div class="audit-drawer">
        <div class="drawer-header">
          <span style="font-weight: 600; color: var(--text-primary);">Audit Detail: {selectedAudit.id}</span>
          <button class="jb-icon-btn" onclick={() => selectedAudit = null}>✕</button>
        </div>
        <div class="drawer-body">
          <div class="drawer-row">
            <span class="d-label">Timestamp:</span>
            <span class="d-val code-text">{selectedAudit.ts}</span>
          </div>
          <div class="drawer-row">
            <span class="d-label">Request ID:</span>
            <span class="d-val code-text">{selectedAudit.request_id}</span>
          </div>
          <div class="drawer-row">
            <span class="d-label">Token ID:</span>
            <span class="d-val code-text">{selectedAudit.token_id}</span>
          </div>
          <div class="drawer-row">
            <span class="d-label">Requesting Client:</span>
            <span class="d-val">{selectedAudit.client}</span>
          </div>
          <div class="drawer-row">
            <span class="d-label">Database Target:</span>
            <span class="d-val">{selectedAudit.connection}</span>
          </div>
          <div class="drawer-row">
            <span class="d-label">Policy Decision:</span>
            <span class="d-val"><strong>{selectedAudit.policy_decision}</strong></span>
          </div>
          <div class="drawer-row">
            <span class="d-label">Policy Rule Matched:</span>
            <span class="d-val code-text">{selectedAudit.policy_rule}</span>
          </div>
          <div class="drawer-row">
            <span class="d-label">Execution Duration:</span>
            <span class="d-val code-text">{selectedAudit.duration_ms} ms</span>
          </div>

          <div style="font-size: 11px; font-weight: 600; color: var(--text-secondary); margin-top: 8px;">
            Executed Query
          </div>
          <pre class="sql-box code-text">{selectedAudit.sql_normalized}</pre>
        </div>
      </div>
    {/if}
  </div>
</div>

<style>
  .audit-view {
    display: flex;
    flex-direction: column;
    height: calc(100vh - var(--titlebar-height) - var(--tabbar-height) - var(--statusbar-height));
    background-color: var(--bg-canvas);
    overflow: hidden;
  }

  .audit-header {
    height: 48px;
    background-color: var(--bg-toolbar);
    border-bottom: 1px solid var(--border-default);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 16px;
    flex-shrink: 0;
  }

  .header-left {
    display: flex;
    flex-direction: column;
  }

  .header-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 13px;
    font-weight: 600;
    color: var(--text-primary);
  }

  .header-desc {
    font-size: 11px;
    color: var(--text-muted);
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .search-box {
    display: flex;
    align-items: center;
    gap: 6px;
    background-color: #1E1F22;
    border: 1px solid var(--border-default);
    border-radius: 4px;
    padding: 0 8px;
    height: 26px;
    color: var(--text-muted);
  }

  .search-box input {
    background: transparent;
    border: none;
    outline: none;
    font-size: 11px;
    color: var(--text-primary);
    width: 180px;
  }

  .filter-select {
    height: 26px;
    background-color: #1E1F22;
    border: 1px solid var(--border-default);
    border-radius: 4px;
    padding: 0 8px;
    font-size: 11px;
    color: var(--text-primary);
  }

  .export-btn {
    height: 26px;
    background-color: #393B40;
    border: 1px solid var(--border-strong);
    border-radius: 4px;
    padding: 0 10px;
    font-size: 11px;
    font-weight: 500;
    color: var(--text-primary);
    transition: background-color 0.15s ease;
  }

  .export-btn:hover {
    background-color: #4E5157;
  }

  .content-layout {
    flex: 1;
    display: flex;
    overflow: hidden;
  }

  .table-container {
    flex: 1;
    overflow: auto;
  }

  .audit-table {
    width: 100%;
    border-collapse: collapse;
    font-size: 12px;
  }

  .audit-table th {
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

  .audit-table td {
    height: 28px;
    border-bottom: 1px solid var(--border-subtle);
    border-right: 1px solid var(--border-subtle);
    padding: 0 8px;
    white-space: nowrap;
    color: var(--text-primary);
  }

  .audit-table tr:hover {
    background-color: var(--bg-hover);
    cursor: pointer;
  }

  .audit-table tr.selected {
    background-color: var(--bg-selected);
  }

  .cell-ts {
    color: var(--text-muted);
    font-size: 11px;
  }

  .decision-pill {
    font-size: 10px;
    font-weight: 700;
    padding: 1px 6px;
    border-radius: 3px;
  }

  .pill-allow { background-color: rgba(87, 211, 140, 0.15); color: #57D38C; border: 1px solid rgba(87, 211, 140, 0.4); }
  .pill-approval { background-color: rgba(237, 162, 0, 0.15); color: #ffc44d; border: 1px solid rgba(237, 162, 0, 0.4); }
  .pill-deny { background-color: rgba(229, 83, 83, 0.15); color: #ff8585; border: 1px solid rgba(229, 83, 83, 0.4); }

  .cell-client {
    color: #DFE1E5;
  }

  .cell-conn {
    color: var(--syntax-function);
  }

  .cell-type {
    font-weight: 600;
    font-size: 11px;
  }

  .cell-num {
    color: var(--syntax-number);
    text-align: right;
    padding-right: 12px;
  }

  .cell-sql {
    color: #9DA0A8;
    max-width: 380px;
  }

  /* Right Drawer */
  .audit-drawer {
    width: 360px;
    background-color: var(--bg-toolbar);
    border-left: 1px solid var(--border-default);
    display: flex;
    flex-direction: column;
    flex-shrink: 0;
  }

  .drawer-header {
    height: 36px;
    padding: 0 12px;
    background-color: #25272A;
    border-bottom: 1px solid var(--border-default);
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .drawer-body {
    padding: 14px;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .drawer-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 11px;
    padding: 4px 0;
    border-bottom: 1px solid var(--border-subtle);
  }

  .d-label {
    color: var(--text-muted);
  }

  .d-val {
    color: var(--text-primary);
  }

  .sql-box {
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
