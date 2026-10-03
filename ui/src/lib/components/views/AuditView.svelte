<script lang="ts">
  import { appState, type AuditRecord } from "../../state/appState.svelte";
  import { exportAudit } from "../../api/audit";
  import {
    Button,
    Badge,
    Input,
    Alert,
    toast,
    Box,
    Flex,
    Stack,
    Text,
    DataGrid,
    DataGridHead,
    DataGridBody,
    DataGridRow,
    DataGridHeadCell,
    DataGridHeaderInner,
    DataGridRowNumHead,
    DataGridRowNum,
    DataGridCell,
  } from "@ohjanus/ui";
  import { Icon } from "@ohjanus/icons";

  let filterStatus = $state<string>("ALL");
  let searchQuery = $state<string>("");
  let selectedAudit = $state<AuditRecord | null>(null);

  let filteredAuditLogs = $derived.by(() => {
    let list = appState.auditLogs;
    if (filterStatus !== "ALL") {
      list = list.filter((item) => item.policy_decision === filterStatus);
    }
    if (searchQuery.trim()) {
      const q = searchQuery.toLowerCase().trim();
      list = list.filter(
        (item) =>
          item.client.toLowerCase().includes(q) ||
          item.sql_normalized.toLowerCase().includes(q) ||
          item.token_id.toLowerCase().includes(q) ||
          item.connection.toLowerCase().includes(q),
      );
    }
    return list;
  });

  async function exportCSV() {
    // Note: the server export supports text search (q) but has no
    // policy-decision filter, so the on-screen decision filter is not forwarded.
    try {
      await exportAudit("csv", { q: searchQuery.trim() || undefined });
      toast.success(
        "Audit Log Exported",
        "Server-side CSV download started (up to 10,000 rows).",
      );
    } catch (e) {
      toast.error("Export Failed", e instanceof Error ? e.message : String(e));
    }
  }
</script>

<Box class="audit-view">
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
  <!-- View Header & Filters -->
  <Flex as="header" class="audit-header" align="center" justify="between">
    <Stack class="header-left" gap="2px">
      <Flex class="header-title" align="center" gap="8px">
        <Icon name="audit" size={16} color="#56A8F5" />
        <span>MCP Gateway Audit Log Trail</span>
      </Flex>
      <span class="header-desc"
        >End-to-end provenance records of all incoming queries and tool requests</span
      >
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
      <DataGrid style="height: 100%;">
        {#snippet overlay()}
          <!-- Floating row count pill -->
          <div class="floating-row-badge" title="Filter count">
            <span>{filteredAuditLogs.length} rows</span>
            <span style="color: var(--text-muted); opacity: 0.6;">|</span>
            <Icon name="more" size={12} />
          </div>
        {/snippet}
        <DataGridHead>
          <DataGridRow>
            <DataGridRowNumHead />
            <DataGridHeadCell width="150px">
              <DataGridHeaderInner>
                <Icon name="clock" size={12} color="#56A8F5" />
                <span>Timestamp</span>
                <span class="ohjanus-data-grid-header-action"
                  ><Icon name="filter" size={9} /></span
                >
                <span class="ohjanus-data-grid-header-action"><Icon name="sort" size={9} /></span>
              </DataGridHeaderInner>
            </DataGridHeadCell>
            <DataGridHeadCell width="110px">
              <DataGridHeaderInner>
                <Icon name="shield" size={12} color="#EDA200" />
                <span>Decision</span>
                <span class="ohjanus-data-grid-header-action"
                  ><Icon name="filter" size={9} /></span
                >
                <span class="ohjanus-data-grid-header-action"><Icon name="sort" size={9} /></span>
              </DataGridHeaderInner>
            </DataGridHeadCell>
            <DataGridHeadCell width="170px">
              <DataGridHeaderInner>
                <Icon name="user" size={12} color="#7A7E85" />
                <span>Client Agent</span>
                <span class="ohjanus-data-grid-header-action"
                  ><Icon name="filter" size={9} /></span
                >
                <span class="ohjanus-data-grid-header-action"><Icon name="sort" size={9} /></span>
              </DataGridHeaderInner>
            </DataGridHeadCell>
            <DataGridHeadCell width="150px">
              <DataGridHeaderInner>
                <Icon name="database" size={12} color="#3B82F6" />
                <span>Connection</span>
                <span class="ohjanus-data-grid-header-action"
                  ><Icon name="filter" size={9} /></span
                >
                <span class="ohjanus-data-grid-header-action"><Icon name="sort" size={9} /></span>
              </DataGridHeaderInner>
            </DataGridHeadCell>
            <DataGridHeadCell width="90px">
              <DataGridHeaderInner>
                <Icon name="terminal" size={12} color="#9DA0A8" />
                <span>Type</span>
                <span class="ohjanus-data-grid-header-action"
                  ><Icon name="filter" size={9} /></span
                >
                <span class="ohjanus-data-grid-header-action"><Icon name="sort" size={9} /></span>
              </DataGridHeaderInner>
            </DataGridHeadCell>
            <DataGridHeadCell width="140px">
              <DataGridHeaderInner>
                <Icon name="table" size={12} color="#3B82F6" />
                <span>Tables</span>
                <span class="ohjanus-data-grid-header-action"
                  ><Icon name="filter" size={9} /></span
                >
                <span class="ohjanus-data-grid-header-action"><Icon name="sort" size={9} /></span>
              </DataGridHeaderInner>
            </DataGridHeadCell>
            <DataGridHeadCell width="70px">
              <DataGridHeaderInner>
                <Icon name="chart" size={12} color="#7A7E85" />
                <span>Rows</span>
                <span class="ohjanus-data-grid-header-action"
                  ><Icon name="filter" size={9} /></span
                >
                <span class="ohjanus-data-grid-header-action"><Icon name="sort" size={9} /></span>
              </DataGridHeaderInner>
            </DataGridHeadCell>
            <DataGridHeadCell width="90px">
              <DataGridHeaderInner>
                <Icon name="lightning" size={12} color="#FACC15" />
                <span>Duration</span>
                <span class="ohjanus-data-grid-header-action"
                  ><Icon name="filter" size={9} /></span
                >
                <span class="ohjanus-data-grid-header-action"><Icon name="sort" size={9} /></span>
              </DataGridHeaderInner>
            </DataGridHeadCell>
            <DataGridHeadCell>
              <DataGridHeaderInner>
                <span>Normalized SQL</span>
                <span class="ohjanus-data-grid-header-action"
                  ><Icon name="filter" size={9} /></span
                >
                <span class="ohjanus-data-grid-header-action"><Icon name="sort" size={9} /></span>
              </DataGridHeaderInner>
            </DataGridHeadCell>
          </DataGridRow>
        </DataGridHead>
        <DataGridBody>
          {#each filteredAuditLogs as item, idx (item.id)}
            <DataGridRow
              selected={selectedAudit?.id === item.id}
              onclick={() => (selectedAudit = item)}
            >
              <DataGridRowNum index={idx} />
              <DataGridCell class="cell-ts">{item.ts}</DataGridCell>
              <DataGridCell class="cell-decision">
                {#if item.policy_decision === "ALLOW"}
                  <Badge variant="success" size="sm">ALLOW</Badge>
                {:else if item.policy_decision === "REQUIRE_APPROVAL"}
                  <Badge variant="warning" size="sm">APPROVAL</Badge>
                {:else}
                  <Badge variant="danger" size="sm">DENY</Badge>
                {/if}
              </DataGridCell>
              <DataGridCell class="cell-client truncate" title={item.client}>
                <Icon name="user" size={12} />
                {item.client}
              </DataGridCell>
              <DataGridCell class="cell-conn">{item.connection}</DataGridCell>
              <DataGridCell class="cell-type">{item.statement_type}</DataGridCell>
              <DataGridCell
                class="cell-tables truncate"
                title={(item.tables ?? []).join(", ")}
              >
                {#if item.tables && item.tables.length > 0}
                  {#each item.tables.slice(0, 3) as t}
                    <Badge variant="default" size="sm" class="tbl-pill"
                      >{t}</Badge
                    >
                  {/each}
                  {#if item.tables.length > 3}
                    <span class="more-tbl">+{item.tables.length - 3}</span>
                  {/if}
                {:else}
                  <span class="no-tbl">—</span>
                {/if}
              </DataGridCell>
              <DataGridCell tone="number" align="right">{item.row_count}</DataGridCell>
              <DataGridCell tone="number" align="right">{item.duration_ms} ms</DataGridCell>
              <DataGridCell tone="secondary" truncate title={item.sql_normalized}>
                {item.sql_normalized}
              </DataGridCell>
            </DataGridRow>
          {/each}
        </DataGridBody>
      </DataGrid>
    </Box>

    <!-- Right Drawer Detail -->
    {#if selectedAudit}
      <Box class="audit-drawer">
        <Flex class="drawer-header" align="center" justify="between">
          <span style="font-weight: 600; color: var(--text-primary);"
            >Audit Detail: {selectedAudit.id}</span
          >
          <button class="jb-icon-btn" onclick={() => (selectedAudit = null)}
            ><Icon name="x" size={12} /></button
          >
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
            <span class="d-label">Tables Touched:</span>
            <span class="d-val">
              {#if selectedAudit.tables && selectedAudit.tables.length > 0}
                {#each selectedAudit.tables as t}
                  <Badge variant="default" size="sm" class="tbl-pill">{t}</Badge
                  >
                {/each}
              {:else}
                <span class="no-tbl">—</span>
              {/if}
            </span>
          </Flex>
          <Flex class="drawer-row" align="center" justify="between">
            <span class="d-label">Policy Decision:</span>
            <span class="d-val"
              ><strong>{selectedAudit.policy_decision}</strong></span
            >
          </Flex>
          <Flex class="drawer-row" align="center" justify="between">
            <span class="d-label">Policy Rule Matched:</span>
            <span class="d-val code-text">{selectedAudit.policy_rule}</span>
          </Flex>
          <Flex class="drawer-row" align="center" justify="between">
            <span class="d-label">Execution Duration:</span>
            <span class="d-val code-text">{selectedAudit.duration_ms} ms</span>
          </Flex>

          <Text
            size="xs"
            weight="semibold"
            color="secondary"
            style="margin-top: 8px;"
          >
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
    font-size: var(--font-size-lg, 16px);
    font-weight: 600;
    color: var(--text-primary);
  }

  :global(.audit-view .header-desc) {
    font-size: var(--font-size-xs, 12px);
    color: var(--text-muted);
  }

  :global(.audit-view .filter-select) {
    height: var(--control-height-sm, 28px);
    background-color: #1e1f22;
    border: 1px solid var(--border-default);
    border-radius: 4px;
    padding: 0 8px;
    font-size: var(--font-size-xs, 12px);
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
    position: relative;
  }

  /* Data grid structure comes from the DataGrid compound (@ohjanus/ui) */
  :global(.audit-view .ohjanus-data-grid-row) {
    cursor: pointer;
  }

  :global(.audit-view .floating-row-badge) {
    position: absolute;
    bottom: 12px;
    right: 16px;
    display: flex;
    align-items: center;
    gap: 6px;
    background-color: var(--bg-toolbar);
    border: 1px solid var(--border-default);
    border-radius: 4px;
    padding: 4px 10px;
    font-size: var(--font-size-xs, 12px);
    color: var(--text-secondary);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.4);
    user-select: none;
  }

  :global(.audit-view .cell-ts) {
    color: var(--text-muted);
    font-size: var(--font-size-xs, 12px);
  }

  :global(.audit-view .cell-client) {
    color: #dfe1e5;
  }

  :global(.audit-view .cell-conn) {
    color: var(--syntax-function);
  }

  :global(.audit-view .cell-type) {
    font-weight: 600;
    font-size: var(--font-size-xs, 12px);
  }

  :global(.audit-view .cell-tables) {
    max-width: 240px;
  }

  :global(.audit-view .tbl-pill) {
    font-size: var(--font-size-2xs, 11px) !important;
    margin-right: 3px;
    font-family: var(--font-code);
  }

  :global(.audit-view .more-tbl) {
    font-size: var(--font-size-2xs, 11px);
    color: var(--text-muted);
  }

  :global(.audit-view .no-tbl) {
    color: var(--text-null, var(--text-muted));
    font-style: italic;
    font-size: var(--font-size-xs, 12px);
  }

  :global(.audit-view .cell-num) {
    color: var(--syntax-number);
    text-align: right;
    padding-right: 12px;
  }

  :global(.audit-view .cell-sql) {
    color: #9da0a8;
    max-width: 380px;
  }

  /* Right Drawer */
  :global(.audit-view .audit-drawer) {
    width: 380px;
    background-color: var(--bg-toolbar);
    border-left: 1px solid var(--border-default);
    display: flex;
    flex-direction: column;
    flex-shrink: 0;
  }

  :global(.audit-view .drawer-header) {
    height: 40px;
    padding: 0 14px;
    background-color: #25272a;
    border-bottom: 1px solid var(--border-default);
    font-size: var(--font-size-md, 15px);
  }

  :global(.audit-view .drawer-body) {
    padding: 16px;
    overflow-y: auto;
  }

  :global(.audit-view .drawer-row) {
    font-size: var(--font-size-xs, 12px);
    padding: 6px 0;
    border-bottom: 1px solid var(--border-subtle);
  }

  :global(.audit-view .d-label) {
    color: var(--text-muted);
  }

  :global(.audit-view .d-val) {
    color: var(--text-primary);
  }

  :global(.audit-view .sql-box) {
    background-color: #1e1f22;
    border: 1px solid var(--border-default);
    padding: 10px 12px;
    border-radius: 4px;
    font-size: var(--font-size-xs, 12px);
    line-height: var(--line-height-normal, 1.45);
    color: #dfe1e5;
    white-space: pre-wrap;
    max-height: 200px;
    overflow-y: auto;
  }
</style>
