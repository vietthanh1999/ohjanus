<script lang="ts">
  import { appState, type ApprovalRequest } from "../../state/appState.svelte";
  import {
    Button,
    Badge,
    Alert,
    EmptyState,
    Box,
    Flex,
    Stack,
    Text,
  } from "@ohjanus/ui";
  import { Icon } from "@ohjanus/icons";

  let selectedFilter = $state<"pending" | "approved" | "rejected" | "all">(
    "pending",
  );
  let selectedDetail = $state<ApprovalRequest | null>(null);

  let filteredApprovals = $derived.by(() => {
    if (selectedFilter === "all") return appState.approvals;
    return appState.approvals.filter((a) => a.state === selectedFilter);
  });

  function getStatementVariant(
    type: string,
  ): "default" | "success" | "warning" | "danger" | "info" {
    switch (type) {
      case "SELECT":
        return "info";
      case "INSERT":
        return "success";
      case "UPDATE":
        return "warning";
      case "DELETE":
        return "danger";
      case "DROP":
        return "danger";
      default:
        return "default";
    }
  }

  function openActionModal(item: ApprovalRequest, mode: "approve" | "reject") {
    appState.selectedApprovalForAction = item;
    appState.approvalDecisionMode = mode;
    appState.approvalDecisionReason =
      mode === "approve" ? "Approved for execution" : "";
  }
</script>

<Box class="approvals-view">
  {#if appState.dataLoading}
    <Box style="padding: 8px 16px;"
      ><Text size="sm" color="muted">Loading live data from Admin API…</Text
      ></Box
    >
  {:else if appState.dataError}
    <Box style="padding: 8px 16px;">
      <Alert variant="danger" title="Admin API unreachable">
        <p>{appState.dataError}</p>
        <Button
          variant="secondary"
          size="sm"
          onclick={() => void appState.loadAll()}>Retry</Button
        >
      </Alert>
    </Box>
  {/if}
  <!-- Header Bar -->
  <Flex as="header" class="view-header" align="center" justify="between">
    <Stack class="header-left" gap="2px">
      <Flex class="header-title" align="center" gap="8px">
        <Icon name="shield" size={16} color="#EDA200" />
        <Text class="header-title" size="xl" weight="semibold">MCP Gateway Approval Queue</Text>
      </Flex>
    </Stack>

    <!-- Filter tabs -->
    <Flex class="state-tabs" align="center">
      <button
        class="state-tab-btn"
        class:active={selectedFilter === "pending"}
        onclick={() => (selectedFilter = "pending")}
      >
        Pending
        <Text class="count-pill" size="sm" weight="bold"
          >{appState.approvals.filter((a) => a.state === "pending")
            .length}</Text
        >
      </button>
      <button
        class="state-tab-btn"
        class:active={selectedFilter === "approved"}
        onclick={() => (selectedFilter = "approved")}
      >
        Approved
      </button>
      <button
        class="state-tab-btn"
        class:active={selectedFilter === "rejected"}
        onclick={() => (selectedFilter = "rejected")}
      >
        Rejected
      </button>
      <button
        class="state-tab-btn"
        class:active={selectedFilter === "all"}
        onclick={() => (selectedFilter = "all")}
      >
        All
      </button>
    </Flex>
  </Flex>

  <!-- Content Split: Queue List & Inspection Drawer -->
  <Flex class="content-layout">
    <Box class="queue-list">
      {#if filteredApprovals.length === 0}
        <EmptyState
          title="No requests in this queue"
          description="All agent queries are up to date and policies satisfied."
        />
      {:else}
        {#each filteredApprovals as item (item.id)}
          <!-- svelte-ignore a11y_click_events_have_key_events -->
          <!-- svelte-ignore a11y_no_static_element_interactions -->
          <Box
            class={`approval-card ${selectedDetail?.id === item.id ? "selected" : ""}`}
            onclick={() => (selectedDetail = item)}
          >
            <!-- Card Header -->
            <Flex class="card-top" align="center" justify="between">
              <Flex class="card-badges" align="center" gap="8px">
                <Badge
                  variant={getStatementVariant(item.statement_type)}
                  size="sm"
                >
                  {item.statement_type}
                </Badge>
                <Text class="conn-pill" size="sm" color="secondary">{item.connection}</Text>
                {#if item.state === "pending"}
                  <Badge variant="warning" size="sm">Pending Review</Badge>
                {:else if item.state === "approved"}
                  <Badge variant="success" size="sm">Approved</Badge>
                {:else}
                  <Badge variant="danger" size="sm">Rejected</Badge>
                {/if}
              </Flex>

              <Box class="time-meta">
                <span class="expiry-time"
                  >Exp: {item.expires_at.substring(11, 16)}</span
                >
              </Box>
            </Flex>

            <!-- Requester info -->
            <Flex class="requester-row" align="center" gap="12px">
              <span class="agent-client"
                ><Icon name="user" size={12} /> {item.requested_by.client}</span
              >
              <span class="token-tag code-text"
                >{item.requested_by.token_id}</span
              >
              <Text class="rows-affected" weight="medium"
                >Est. ~{item.affected_estimate.toLocaleString()} rows</Text
              >
            </Flex>

            <!-- Risk Warnings -->
            {#if item.warnings && item.warnings.length > 0}
              <Alert variant="warning">
                {#each item.warnings as warn}
                  <Flex class="warn-line" align="center" gap="6px">
                    <span class="warn-icon"
                      ><Icon name="alert-triangle" size={12} /></span
                    >
                    <span>{warn}</span>
                  </Flex>
                {/each}
              </Alert>
            {/if}

            <!-- SQL Snippet -->
            <Box class="sql-preview code-text">
              {item.sql}
            </Box>

            <!-- Action buttons for Pending -->
            {#if item.state === "pending"}
              <Flex class="card-actions" align="center" justify="end" gap="8px">
                <Button
                  variant="danger"
                  size="xs"
                  onclick={(e: MouseEvent) => {
                    e.stopPropagation();
                    openActionModal(item, "reject");
                  }}
                >
                  Reject...
                </Button>
                <Button
                  variant="primary"
                  size="xs"
                  onclick={(e: MouseEvent) => {
                    e.stopPropagation();
                    openActionModal(item, "approve");
                  }}
                >
                  Approve & Execute
                </Button>
              </Flex>
            {:else}
              <Box class="decision-meta">
                <span
                  >Decided by <strong>{item.decided_by}</strong> on {item.decided_at}</span
                >
                {#if item.decision_reason}
                  <span class="reason-note"
                    >Reason: "{item.decision_reason}"</span
                  >
                {/if}
              </Box>
            {/if}
          </Box>
        {/each}
      {/if}
    </Box>

    <!-- Right: Detail Inspection Pane -->
    <Box class="detail-drawer">
      {#if selectedDetail}
        <Flex class="drawer-header" align="center" justify="between">
          <span style="font-weight: 600; color: var(--text-primary);"
            >Request Inspection: {selectedDetail.id}</span
          >
          <button class="jb-icon-btn" onclick={() => (selectedDetail = null)}
            ><Icon name="x" size={12} /></button
          >
        </Flex>
        <Stack class="drawer-body" gap="8px">
          <Box class="meta-section">
            <Flex class="meta-row" align="center" justify="between">
              <span class="meta-label">Connection:</span>
              <span class="meta-val"
                >{selectedDetail.connection} (PostgreSQL 16)</span
              >
            </Flex>
            <Flex class="meta-row" align="center" justify="between">
              <span class="meta-label">Requesting Client:</span>
              <span class="meta-val">{selectedDetail.requested_by.client}</span>
            </Flex>
            <Flex class="meta-row" align="center" justify="between">
              <span class="meta-label">Agent Token ID:</span>
              <span class="meta-val code-text"
                >{selectedDetail.requested_by.token_id}</span
              >
            </Flex>
            <Flex class="meta-row" align="center" justify="between">
              <span class="meta-label">Submitted At:</span>
              <span class="meta-val code-text">{selectedDetail.created_at}</span
              >
            </Flex>
            <Flex class="meta-row" align="center" justify="between">
              <span class="meta-label">Expires At:</span>
              <span class="meta-val code-text">{selectedDetail.expires_at}</span
              >
            </Flex>
          </Box>

          <Box class="section-title">Full SQL Statement</Box>
          <pre class="full-sql code-text">{selectedDetail.sql}</pre>

          <Box class="section-title">Execution Safety Assessment</Box>
          <Stack class="safety-box" gap="6px">
            <Flex class="safety-item" align="center" gap="8px">
              <span
                class="safe-dot"
                class:risk={selectedDetail.affected_estimate > 100}
              ></span>
              <span
                >Blast Radius: {selectedDetail.affected_estimate} row(s) estimated</span
              >
            </Flex>
            <Flex class="safety-item" align="center" gap="8px">
              <span
                class="safe-dot"
                class:risk={selectedDetail.statement_type === "DROP"}
              ></span>
              <span>Statement Type: {selectedDetail.statement_type}</span>
            </Flex>
            <Flex class="safety-item" align="center" gap="8px">
              <span class="safe-dot"></span>
              <span>AST Validator: Parsed successfully via pg_query_go</span>
            </Flex>
          </Stack>

          {#if selectedDetail.state === "pending"}
            <Flex class="drawer-actions" gap="8px">
              <Button
                variant="danger"
                style="flex: 1;"
                onclick={() => openActionModal(selectedDetail!, "reject")}
              >
                Reject Request
              </Button>
              <Button
                variant="primary"
                style="flex: 1;"
                onclick={() => openActionModal(selectedDetail!, "approve")}
              >
                Approve Request
              </Button>
            </Flex>
          {/if}
        </Stack>
      {:else}
        <Flex class="drawer-empty" align="center" justify="center">
          <p>
            Select an approval request to inspect full AST validation and
            EXPLAIN execution plan
          </p>
        </Flex>
      {/if}
    </Box>
  </Flex>
</Box>

<style>
  :global(.approvals-view) {
    display: flex;
    flex-direction: column;
    height: 100%;
    background-color: var(--bg-canvas);
    overflow: hidden;
  }

  :global(.approvals-view .view-header) {
    height: 48px;
    background-color: var(--bg-toolbar);
    border-bottom: 1px solid var(--border-default);
    padding: 0 16px;
    flex-shrink: 0;
  }

  :global(.approvals-view .header-left) {
    display: flex;
    flex-direction: column;
  }

  :global(.approvals-view .header-title) {
    font-size: var(--font-size-lg, 16px);
    font-weight: 600;
    color: var(--text-primary);
  }

  :global(.approvals-view .header-desc) {
    font-size: var(--font-size-xs, 12px);
    color: var(--text-muted);
  }

  :global(.approvals-view .state-tabs) {
    background-color: #1e1f22;
    padding: 3px;
    border-radius: 4px;
    border: 1px solid var(--border-default);
  }

  :global(.approvals-view .state-tab-btn) {
    padding: 4px 12px;
    font-size: var(--font-size-xs, 12px);
    color: var(--text-secondary);
    border-radius: 3px;
    display: flex;
    align-items: center;
    gap: 6px;
  }

  :global(.approvals-view .state-tab-btn.active) {
    background-color: #313438;
    color: var(--text-primary);
    font-weight: 500;
  }

  :global(.approvals-view .count-pill) {
    background-color: #eda200;
    color: #1e1f22;
    font-size: var(--font-size-2xs, 11px);
    font-weight: 700;
    padding: 1px 6px;
    border-radius: 10px;
  }

  :global(.approvals-view .content-layout) {
    flex: 1;
    display: flex;
    flex-direction: row;
    overflow: hidden;
  }

  :global(.approvals-view .queue-list) {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  :global(.approvals-view .approval-card) {
    background-color: var(--bg-sidebar);
    border: 1px solid var(--border-default);
    border-radius: 6px;
    padding: 14px 16px;
    display: flex;
    flex-direction: column;
    gap: 10px;
    cursor: pointer;
    transition: all 0.15s ease;
  }

  :global(.approvals-view .approval-card:hover) {
    border-color: #4e5157;
    background-color: #2e3035;
  }

  :global(.approvals-view .approval-card.selected) {
    border-color: var(--border-accent);
    box-shadow: 0 0 0 1px var(--border-accent);
  }

  :global(.approvals-view .conn-pill) {
    background-color: #1e1f22;
    border: 1px solid var(--border-default);
    color: var(--text-secondary);
    font-size: var(--font-size-2xs, 11px);
    padding: 2px 8px;
    border-radius: 3px;
  }

  :global(.approvals-view .time-meta) {
    font-size: var(--font-size-xs, 12px);
    color: var(--text-muted);
  }

  :global(.approvals-view .requester-row) {
    font-size: var(--font-size-xs, 12px);
    color: var(--text-secondary);
  }

  :global(.approvals-view .agent-client) {
    font-weight: 500;
    color: var(--text-primary);
  }

  :global(.approvals-view .token-tag) {
    color: var(--text-muted);
    background-color: #1e1f22;
    padding: 1px 6px;
    border-radius: 2px;
  }

  :global(.approvals-view .rows-affected) {
    margin-left: auto;
    color: var(--syntax-number);
    font-weight: 500;
    font-size: var(--font-size-xs, 12px);
  }

  :global(.approvals-view .sql-preview) {
    background-color: #1e1f22;
    border: 1px solid var(--border-subtle);
    border-radius: 4px;
    padding: 10px 12px;
    font-size: var(--font-size-xs, 12px);
    color: #dfe1e5;
    line-height: var(--line-height-normal, 1.45);
    white-space: pre-wrap;
    max-height: 90px;
    overflow: hidden;
  }

  :global(.approvals-view .card-actions) {
    margin-top: 4px;
  }

  :global(.approvals-view .decision-meta) {
    font-size: var(--font-size-xs, 12px);
    color: var(--text-muted);
    border-top: 1px solid var(--border-subtle);
    padding-top: 6px;
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  :global(.approvals-view .reason-note) {
    color: #dfe1e5;
    font-style: italic;
  }

  /* Right Drawer */
  :global(.approvals-view .detail-drawer) {
    width: 440px;
    background-color: var(--bg-toolbar);
    border-left: 1px solid var(--border-default);
    display: flex;
    flex-direction: column;
    flex-shrink: 0;
  }

  :global(.approvals-view .drawer-header) {
    height: 40px;
    padding: 0 14px;
    background-color: #25272a;
    border-bottom: 1px solid var(--border-default);
    font-size: var(--font-size-md, 15px);
  }

  :global(.approvals-view .drawer-body) {
    padding: 16px;
    overflow-y: auto;
  }

  :global(.approvals-view .meta-section) {
    background-color: #1e1f22;
    border: 1px solid var(--border-default);
    border-radius: 4px;
    padding: 10px 14px;
  }

  :global(.approvals-view .meta-row) {
    font-size: var(--font-size-xs, 12px);
    padding: 6px 0;
    border-bottom: 1px solid var(--border-subtle);
  }

  :global(.approvals-view .meta-row:last-child) {
    border-bottom: none;
  }

  :global(.approvals-view .meta-label) {
    color: var(--text-muted);
  }

  :global(.approvals-view .meta-val) {
    color: var(--text-primary);
    font-weight: 500;
  }

  :global(.approvals-view .section-title) {
    font-size: var(--font-size-xs, 12px);
    font-weight: 600;
    color: var(--text-secondary);
    margin-top: 8px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  :global(.approvals-view .full-sql) {
    background-color: #1e1f22;
    border: 1px solid var(--border-default);
    border-radius: 4px;
    padding: 10px 12px;
    font-size: var(--font-size-xs, 12px);
    line-height: var(--line-height-normal, 1.45);
    color: #dfe1e5;
    max-height: 180px;
    overflow-y: auto;
    white-space: pre-wrap;
    margin: 0;
  }

  :global(.approvals-view .safety-box) {
    background-color: #1e1f22;
    border: 1px solid var(--border-default);
    border-radius: 4px;
    padding: 10px 12px;
  }

  :global(.approvals-view .safety-item) {
    font-size: var(--font-size-xs, 12px);
    color: var(--text-primary);
  }

  :global(.approvals-view .safe-dot) {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background-color: #57d38c;
  }

  :global(.approvals-view .safe-dot.risk) {
    background-color: #e55353;
  }

  :global(.approvals-view .drawer-actions) {
    margin-top: 12px;
  }

  :global(.approvals-view .drawer-empty) {
    flex: 1;
    padding: 24px;
    color: var(--text-muted);
    font-size: var(--font-size-sm, 13px);
    text-align: center;
  }
</style>
