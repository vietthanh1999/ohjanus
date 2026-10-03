<script lang="ts">
  import { appState, type ApprovalRequest } from '../../state/appState.svelte';

  let selectedFilter = $state<'pending' | 'approved' | 'rejected' | 'all'>('pending');
  let selectedDetail = $state<ApprovalRequest | null>(null);

  let filteredApprovals = $derived.by(() => {
    if (selectedFilter === 'all') return appState.approvals;
    return appState.approvals.filter(a => a.state === selectedFilter);
  });

  function getStatementBadgeClass(type: string) {
    switch (type) {
      case 'SELECT': return 'badge-select';
      case 'INSERT': return 'badge-insert';
      case 'UPDATE': return 'badge-update';
      case 'DELETE': return 'badge-delete';
      case 'DROP': return 'badge-drop';
      default: return '';
    }
  }

  function openActionModal(item: ApprovalRequest, mode: 'approve' | 'reject') {
    appState.selectedApprovalForAction = item;
    appState.approvalDecisionMode = mode;
    appState.approvalDecisionReason = mode === 'approve' ? 'Approved for execution' : '';
  }
</script>

<div class="approvals-view">
  <!-- Header Bar -->
  <div class="view-header">
    <div class="header-left">
      <div class="header-title">
        <svg width="16" height="16" viewBox="0 0 16 16" fill="#EDA200">
          <path fill-rule="evenodd" d="M8 0c-.69 0-1.843.265-2.928.56-1.11.3-2.229.655-2.887.87a1.54 1.54 0 00-1.044 1.262c-.596 4.477.787 7.795 2.464 9.99 1.579 2.065 3.444 3.009 4.395 3.318.066.022.135.034.204.034s.138-.012.204-.034c.951-.309 2.816-1.253 4.395-3.318 1.677-2.195 3.06-5.513 2.464-9.99a1.54 1.54 0 00-1.044-1.263 62.467 62.467 0 00-2.887-.87C9.843.266 8.69 0 8 0zm2.146 5.146a.5.5 0 01.708.708l-3 3a.5.5 0 01-.708 0l-1.5-1.5a.5.5 0 11.708-.708L7.5 7.793l2.646-2.647z"/>
        </svg>
        <span>MCP Gateway Approval Queue</span>
      </div>
      <span class="header-desc">Review and authorize AI Agent write operations before database execution</span>
    </div>

    <!-- Filter tabs -->
    <div class="state-tabs">
      <button
        class="state-tab-btn"
        class:active={selectedFilter === 'pending'}
        onclick={() => selectedFilter = 'pending'}
      >
        Pending
        <span class="count-pill">{appState.approvals.filter(a => a.state === 'pending').length}</span>
      </button>
      <button
        class="state-tab-btn"
        class:active={selectedFilter === 'approved'}
        onclick={() => selectedFilter = 'approved'}
      >
        Approved
      </button>
      <button
        class="state-tab-btn"
        class:active={selectedFilter === 'rejected'}
        onclick={() => selectedFilter = 'rejected'}
      >
        Rejected
      </button>
      <button
        class="state-tab-btn"
        class:active={selectedFilter === 'all'}
        onclick={() => selectedFilter = 'all'}
      >
        All
      </button>
    </div>
  </div>

  <!-- Content Split: Queue List & Inspection Drawer -->
  <div class="content-layout">
    <div class="queue-list">
      {#if filteredApprovals.length === 0}
        <div class="empty-state">
          <span style="font-size: 28px;">🎉</span>
          <p style="color: var(--text-primary); font-weight: 500; margin-top: 8px;">No requests in this queue</p>
          <span style="color: var(--text-muted); font-size: 11px;">All agent queries are up to date and policies satisfied.</span>
        </div>
      {:else}
        {#each filteredApprovals as item (item.id)}
          <div
            class="approval-card"
            class:selected={selectedDetail?.id === item.id}
            onclick={() => selectedDetail = item}
          >
            <!-- Card Header -->
            <div class="card-top">
              <div class="card-badges">
                <span class="stmt-badge {getStatementBadgeClass(item.statement_type)}">
                  {item.statement_type}
                </span>
                <span class="conn-pill">{item.connection}</span>
                {#if item.state === 'pending'}
                  <span class="status-pill status-pending">Pending Review</span>
                {:else if item.state === 'approved'}
                  <span class="status-pill status-approved">Approved</span>
                {:else}
                  <span class="status-pill status-rejected">Rejected</span>
                {/if}
              </div>

              <div class="time-meta">
                <span class="expiry-time">Exp: {item.expires_at.substring(11, 16)}</span>
              </div>
            </div>

            <!-- Requester info -->
            <div class="requester-row">
              <span class="agent-client">🤖 {item.requested_by.client}</span>
              <span class="token-tag code-text">{item.requested_by.token_id}</span>
              <span class="rows-affected">Est. ~{item.affected_estimate.toLocaleString()} rows</span>
            </div>

            <!-- Risk Warnings -->
            {#if item.warnings && item.warnings.length > 0}
              <div class="warning-box">
                {#each item.warnings as warn}
                  <div class="warn-line">
                    <span class="warn-icon">⚠️</span>
                    <span>{warn}</span>
                  </div>
                {/each}
              </div>
            {/if}

            <!-- SQL Snippet -->
            <div class="sql-preview code-text">
              {item.sql}
            </div>

            <!-- Action buttons for Pending -->
            {#if item.state === 'pending'}
              <div class="card-actions">
                <button
                  class="action-btn reject-btn"
                  onclick={(e) => {
                    e.stopPropagation();
                    openActionModal(item, 'reject');
                  }}
                >
                  Reject...
                </button>
                <button
                  class="action-btn approve-btn"
                  onclick={(e) => {
                    e.stopPropagation();
                    openActionModal(item, 'approve');
                  }}
                >
                  Approve & Execute
                </button>
              </div>
            {:else}
              <div class="decision-meta">
                <span>Decided by <strong>{item.decided_by}</strong> on {item.decided_at}</span>
                {#if item.decision_reason}
                  <span class="reason-note">Reason: "{item.decision_reason}"</span>
                {/if}
              </div>
            {/if}
          </div>
        {/each}
      {/if}
    </div>

    <!-- Right: Detail Inspection Pane -->
    <div class="detail-drawer">
      {#if selectedDetail}
        <div class="drawer-header">
          <span style="font-weight: 600; color: var(--text-primary);">Request Inspection: {selectedDetail.id}</span>
          <button class="jb-icon-btn" onclick={() => selectedDetail = null}>✕</button>
        </div>
        <div class="drawer-body">
          <div class="meta-section">
            <div class="meta-row">
              <span class="meta-label">Connection:</span>
              <span class="meta-val">{selectedDetail.connection} (PostgreSQL 16)</span>
            </div>
            <div class="meta-row">
              <span class="meta-label">Requesting Client:</span>
              <span class="meta-val">{selectedDetail.requested_by.client}</span>
            </div>
            <div class="meta-row">
              <span class="meta-label">Agent Token ID:</span>
              <span class="meta-val code-text">{selectedDetail.requested_by.token_id}</span>
            </div>
            <div class="meta-row">
              <span class="meta-label">Submitted At:</span>
              <span class="meta-val code-text">{selectedDetail.created_at}</span>
            </div>
            <div class="meta-row">
              <span class="meta-label">Expires At:</span>
              <span class="meta-val code-text">{selectedDetail.expires_at}</span>
            </div>
          </div>

          <div class="section-title">Full SQL Statement</div>
          <pre class="full-sql code-text">{selectedDetail.sql}</pre>

          <div class="section-title">Execution Safety Assessment</div>
          <div class="safety-box">
            <div class="safety-item">
              <span class="safe-dot" class:risk={selectedDetail.affected_estimate > 100}></span>
              <span>Blast Radius: {selectedDetail.affected_estimate} row(s) estimated</span>
            </div>
            <div class="safety-item">
              <span class="safe-dot" class:risk={selectedDetail.statement_type === 'DROP'}></span>
              <span>Statement Type: {selectedDetail.statement_type}</span>
            </div>
            <div class="safety-item">
              <span class="safe-dot"></span>
              <span>AST Validator: Parsed successfully via pg_query_go</span>
            </div>
          </div>

          {#if selectedDetail.state === 'pending'}
            <div class="drawer-actions">
              <button
                class="jb-btn-danger"
                style="flex: 1;"
                onclick={() => openActionModal(selectedDetail!, 'reject')}
              >
                Reject Request
              </button>
              <button
                class="jb-btn-primary"
                style="flex: 1; background-color: var(--action-success);"
                onclick={() => openActionModal(selectedDetail!, 'approve')}
              >
                Approve Request
              </button>
            </div>
          {/if}
        </div>
      {:else}
        <div class="drawer-empty">
          <p>Select an approval request to inspect full AST validation and EXPLAIN execution plan</p>
        </div>
      {/if}
    </div>
  </div>
</div>

<style>
  .approvals-view {
    display: flex;
    flex-direction: column;
    height: calc(100vh - var(--titlebar-height) - var(--tabbar-height) - var(--statusbar-height));
    background-color: var(--bg-canvas);
    overflow: hidden;
  }

  .view-header {
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

  .state-tabs {
    display: flex;
    align-items: center;
    background-color: #1E1F22;
    padding: 2px;
    border-radius: 4px;
    border: 1px solid var(--border-default);
  }

  .state-tab-btn {
    padding: 3px 10px;
    font-size: 11px;
    color: var(--text-secondary);
    border-radius: 3px;
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .state-tab-btn.active {
    background-color: #313438;
    color: var(--text-primary);
    font-weight: 500;
  }

  .count-pill {
    background-color: #EDA200;
    color: #1E1F22;
    font-size: 10px;
    font-weight: 700;
    padding: 0 4px;
    border-radius: 6px;
  }

  .content-layout {
    flex: 1;
    display: flex;
    overflow: hidden;
  }

  .queue-list {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 60px 0;
  }

  .approval-card {
    background-color: var(--bg-sidebar);
    border: 1px solid var(--border-default);
    border-radius: 6px;
    padding: 12px 14px;
    display: flex;
    flex-direction: column;
    gap: 8px;
    cursor: pointer;
    transition: all 0.15s ease;
  }

  .approval-card:hover {
    border-color: #4E5157;
    background-color: #2E3035;
  }

  .approval-card.selected {
    border-color: var(--border-accent);
    box-shadow: 0 0 0 1px var(--border-accent);
  }

  .card-top {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .card-badges {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .stmt-badge {
    font-size: 10px;
    font-weight: 700;
    padding: 2px 6px;
    border-radius: 3px;
    letter-spacing: 0.5px;
  }

  .badge-select { background-color: rgba(53, 116, 240, 0.2); color: #79a8ff; border: 1px solid rgba(53, 116, 240, 0.4); }
  .badge-insert { background-color: rgba(87, 211, 140, 0.2); color: #57D38C; border: 1px solid rgba(87, 211, 140, 0.4); }
  .badge-update { background-color: rgba(237, 162, 0, 0.2); color: #ffc44d; border: 1px solid rgba(237, 162, 0, 0.4); }
  .badge-delete { background-color: rgba(249, 115, 22, 0.2); color: #fb923c; border: 1px solid rgba(249, 115, 22, 0.4); }
  .badge-drop   { background-color: rgba(229, 83, 83, 0.2); color: #ff8585; border: 1px solid rgba(229, 83, 83, 0.4); }

  .conn-pill {
    background-color: #1E1F22;
    border: 1px solid var(--border-default);
    color: var(--text-secondary);
    font-size: 11px;
    padding: 1px 6px;
    border-radius: 3px;
  }

  .status-pill {
    font-size: 10px;
    font-weight: 500;
    padding: 1px 6px;
    border-radius: 3px;
  }

  .status-pending { background-color: rgba(237, 162, 0, 0.15); color: #ffc44d; }
  .status-approved { background-color: rgba(87, 211, 140, 0.15); color: #57D38C; }
  .status-rejected { background-color: rgba(229, 83, 83, 0.15); color: #ff8585; }

  .time-meta {
    font-size: 11px;
    color: var(--text-muted);
  }

  .requester-row {
    display: flex;
    align-items: center;
    gap: 12px;
    font-size: 11px;
    color: var(--text-secondary);
  }

  .agent-client {
    font-weight: 500;
    color: var(--text-primary);
  }

  .token-tag {
    color: var(--text-muted);
    background-color: #1E1F22;
    padding: 1px 4px;
    border-radius: 2px;
  }

  .rows-affected {
    margin-left: auto;
    color: var(--syntax-number);
    font-weight: 500;
  }

  .warning-box {
    background-color: rgba(237, 162, 0, 0.08);
    border-left: 3px solid #EDA200;
    padding: 6px 10px;
    border-radius: 2px;
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .warn-line {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 11px;
    color: #ffc44d;
  }

  .sql-preview {
    background-color: #1E1F22;
    border: 1px solid var(--border-subtle);
    border-radius: 4px;
    padding: 8px 10px;
    font-size: 11px;
    color: #DFE1E5;
    line-height: 16px;
    white-space: pre-wrap;
    max-height: 90px;
    overflow: hidden;
  }

  .card-actions {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    gap: 8px;
    margin-top: 4px;
  }

  .action-btn {
    padding: 4px 12px;
    border-radius: 3px;
    font-size: 11px;
    font-weight: 500;
    transition: all 0.15s ease;
  }

  .reject-btn {
    background-color: #393B40;
    color: #DFE1E5;
    border: 1px solid #4E5157;
  }

  .reject-btn:hover {
    background-color: var(--action-danger);
    color: #FFFFFF;
    border-color: var(--action-danger);
  }

  .approve-btn {
    background-color: var(--action-success);
    color: #14281B;
    font-weight: 600;
  }

  .approve-btn:hover {
    background-color: var(--action-success-hover);
  }

  .decision-meta {
    font-size: 11px;
    color: var(--text-muted);
    border-top: 1px solid var(--border-subtle);
    padding-top: 6px;
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .reason-note {
    font-style: italic;
    color: var(--text-secondary);
  }

  /* Right Drawer */
  .detail-drawer {
    width: 380px;
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
    gap: 12px;
  }

  .drawer-empty {
    padding: 40px 20px;
    color: var(--text-muted);
    text-align: center;
    font-size: 12px;
  }

  .meta-section {
    display: flex;
    flex-direction: column;
    gap: 6px;
    background-color: #1E1F22;
    padding: 8px 10px;
    border-radius: 4px;
    border: 1px solid var(--border-subtle);
  }

  .meta-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 11px;
  }

  .meta-label {
    color: var(--text-muted);
  }

  .meta-val {
    color: var(--text-primary);
  }

  .section-title {
    font-size: 11px;
    font-weight: 600;
    color: var(--text-secondary);
    text-transform: uppercase;
  }

  .full-sql {
    background-color: #1E1F22;
    border: 1px solid var(--border-default);
    padding: 10px;
    border-radius: 4px;
    font-size: 11px;
    line-height: 18px;
    color: #DFE1E5;
    white-space: pre-wrap;
    max-height: 220px;
    overflow-y: auto;
  }

  .safety-box {
    background-color: #1E1F22;
    border: 1px solid var(--border-default);
    border-radius: 4px;
    padding: 8px 10px;
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .safety-item {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 11px;
    color: var(--text-secondary);
  }

  .safe-dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background-color: var(--action-success);
  }

  .safe-dot.risk {
    background-color: var(--action-danger);
  }

  .drawer-actions {
    display: flex;
    gap: 8px;
    margin-top: 8px;
  }
</style>
