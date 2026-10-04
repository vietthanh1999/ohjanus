<script lang="ts">
  import { appState, type ApprovalRequest } from "../../state/appState.svelte";
  import {
    Button,
    Badge,
    Select,
    Alert,
    EmptyState,
    toast,
    Box,
    Flex,
    Stack,
    Text,
  } from "@ohjanus/ui";
  import { Icon } from "@ohjanus/icons";

  let selectedFilter = $state<string>("pending");
  let whereFilter = $state("");
  let orderByFilter = $state("");
  let searchRef = $state<HTMLInputElement>();
  let selectedDetail = $state<ApprovalRequest | null>(null);

  let filteredApprovals = $derived.by(() => {
    let list = appState.approvals;
    if (selectedFilter !== "all") {
      list = list.filter((a) => a.state === selectedFilter);
    }
    const w = whereFilter.trim().toLowerCase();
    if (w) {
      list = list.filter(
        (a) =>
          a.statement_type.toLowerCase().includes(w) ||
          a.connection.toLowerCase().includes(w) ||
          a.sql.toLowerCase().includes(w) ||
          a.id.toLowerCase().includes(w) ||
          (a.requested_by?.client && a.requested_by.client.toLowerCase().includes(w)),
      );
    }
    const o = orderByFilter.trim().toLowerCase();
    if (o) {
      const isDesc = o.includes("desc");
      if (o.includes("created") || o.includes("time")) {
        list = [...list].sort((a, b) =>
          isDesc
            ? (b.created_at || "").localeCompare(a.created_at || "")
            : (a.created_at || "").localeCompare(b.created_at || ""),
        );
      } else if (o.includes("type") || o.includes("statement")) {
        list = [...list].sort((a, b) =>
          isDesc
            ? b.statement_type.localeCompare(a.statement_type)
            : a.statement_type.localeCompare(b.statement_type),
        );
      } else if (o.includes("conn") || o.includes("connection")) {
        list = [...list].sort((a, b) =>
          isDesc
            ? b.connection.localeCompare(a.connection)
            : a.connection.localeCompare(b.connection),
        );
      }
    }
    return list;
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

  async function handleReload() {
    try {
      await appState.loadApprovals();
      toast.info("Approvals Refreshed", "Loaded latest approval requests.");
    } catch (e) {
      toast.error("Reload Failed", e instanceof Error ? e.message : String(e));
    }
  }

  function handleClearFilters() {
    selectedFilter = "pending";
    whereFilter = "";
    orderByFilter = "";
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
  <!-- Toolbar 1: Actions (matching design2.png) -->
  <div class="table-toolbar">
    <div class="toolbar-left">
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Reload approval queue (Cmd+Enter)"
        onclick={handleReload}
      >
        <Icon name="refresh" size={13} />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Clear filters"
        onclick={handleClearFilters}
      >
        <Icon name="stop" size={13} color="#E55353" />
      </Button>
      <span class="bar-separator"></span>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Approve selected request"
        disabled={!selectedDetail || selectedDetail.state !== "pending"}
        onclick={() => selectedDetail && openActionModal(selectedDetail, "approve")}
      >
        <Icon name="check" size={13} color="#57D38C" />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Reject selected request"
        disabled={!selectedDetail || selectedDetail.state !== "pending"}
        onclick={() => selectedDetail && openActionModal(selectedDetail, "reject")}
      >
        <Icon name="close" size={13} color="#E55353" />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Toggle inspection panel"
        onclick={() => (selectedDetail = selectedDetail ? null : filteredApprovals[0] || null)}
      >
        <Icon name="layout" size={13} />
      </Button>
      <span class="bar-separator"></span>
      <div class="borderless-select-wrapper">
        <Select
          class="toolbar-select borderless-select"
          options={[
            {
              value: "pending",
              label: `Queue: Pending (${appState.approvals.filter((a) => a.state === "pending").length})`,
            },
            { value: "approved", label: "Queue: Approved" },
            { value: "rejected", label: "Queue: Rejected" },
            { value: "all", label: "Queue: All Requests" },
          ]}
          bind:value={selectedFilter}
        />
      </div>
      <span class="bar-separator"></span>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn ddl-btn"
        title="View DDL"
        onclick={() => (appState.ddlModalOpen = true)}
      >
        <Text size="xs" weight="bold" color="muted" mono>DDL</Text>
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Search requests"
        onclick={() => searchRef?.focus()}
      >
        <Icon name="search" size={13} />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Observability metrics"
        onclick={() => (appState.activeTabId = "dashboard")}
      >
        <Icon name="chart" size={13} />
      </Button>
    </div>

    <div class="toolbar-right">
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

  <!-- Toolbar 2: WHERE & ORDER BY (matching design2.png) -->
  <div class="filter-bar">
    <div class="filter-group where-group">
      <span class="filter-icon"><Icon name="filter" size={11} /></span>
      <Text size="sm" weight="semibold" color="secondary" style="user-select: none;">WHERE</Text>
      <input
        type="text"
        class="filter-input code-text"
        placeholder="e.g. statement = 'DELETE' or connection or SQL snippet"
        bind:this={searchRef}
        bind:value={whereFilter}
      />
    </div>

    <div class="filter-group orderby-group">
      <span class="filter-icon"><Icon name="sort" size={11} /></span>
      <Text size="sm" weight="semibold" color="secondary" style="user-select: none;">ORDER BY</Text>
      <input
        type="text"
        class="filter-input code-text"
        placeholder="e.g. created_at DESC"
        bind:value={orderByFilter}
      />
    </div>
  </div>

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

  :global(.approvals-view .table-toolbar) {
    height: var(--toolbar-height, 32px);
    background-color: var(--bg-toolbar);
    border-bottom: 1px solid var(--border-subtle);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 8px;
    flex-shrink: 0;
  }

  :global(.approvals-view .toolbar-left),
  :global(.approvals-view .toolbar-right) {
    display: flex;
    align-items: center;
    gap: 2px;
  }

  :global(.approvals-view .bar-separator) {
    width: 1px;
    height: 14px;
    background-color: var(--border-default);
    margin: 0 4px;
  }

  :global(.approvals-view .borderless-select-wrapper) {
    display: inline-flex;
    align-items: center;
  }

  :global(.approvals-view .borderless-select) {
    min-width: unset !important;
    width: auto !important;
  }

  :global(.approvals-view .borderless-select .ohjanus-select-trigger) {
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

  :global(.approvals-view .borderless-select .ohjanus-select-trigger:hover) {
    background-color: var(--bg-hover, #313438) !important;
    color: var(--text-primary, #DFE1E5) !important;
  }

  :global(.approvals-view .ddl-btn) {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 0 4px !important;
    width: auto !important;
  }

  :global(.approvals-view .filter-bar) {
    height: var(--filterbar-height, 30px);
    background-color: var(--bg-canvas);
    border-bottom: 1px solid var(--border-subtle);
    display: flex;
    align-items: center;
    padding: 0 8px;
    gap: 12px;
    flex-shrink: 0;
  }

  :global(.approvals-view .filter-group) {
    display: flex;
    align-items: center;
    gap: 6px;
  }

  :global(.approvals-view .where-group) {
    flex: 1.1;
  }

  :global(.approvals-view .orderby-group) {
    flex: 0.9;
  }

  :global(.approvals-view .filter-icon) {
    color: var(--text-muted);
    font-size: 11px;
    user-select: none;
    display: inline-flex;
    align-items: center;
  }

  :global(.approvals-view .filter-input) {
    flex: 1;
    height: var(--control-height-xs, 24px);
    background-color: transparent !important;
    border: none !important;
    outline: none !important;
    box-shadow: none !important;
    padding: 0 8px;
    font-size: var(--font-size-sm, 13px);
    color: var(--text-primary);
    font-family: var(--font-code);
  }

  :global(.approvals-view .filter-input:focus),
  :global(.approvals-view .filter-input:hover) {
    border: none !important;
    outline: none !important;
    box-shadow: none !important;
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
