<script lang="ts">
  import type { AuditRecord } from '@/entities/audit-record';
  import { auditState } from '@/features/audit';
  import { settingsState } from '@/features/settings';
  import { bootState } from '@/app/boot.svelte';
  import { exportAudit } from '@/features/audit';
  import {
    Button,
    Badge,
    Alert,
    toast,
    Box,
    Flex,
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
  import {
    Toolbar,
    ToolbarSeparator,
    BorderlessSelect,
    FilterBar,
    FloatingRowCount,
  } from '@/shared/ui/toolbar';
  import AuditInspectionDrawer from './AuditInspectionDrawer.svelte';

  let filterStatus = $state<string>("ALL");
  let searchQuery = $state<string>("");
  let orderByFilter = $state<string>("");
  let exportFormat = $state<string>("CSV");
  let searchRef = $state<HTMLInputElement>();
  let selectedAudit = $state<AuditRecord | null>(null);

  let filteredAuditLogs = $derived.by(() => {
    let list = [...auditState.auditLogs];
    if (filterStatus !== "ALL") {
      list = list.filter((item) => item.policy_decision === filterStatus);
    }
    const q = searchQuery.toLowerCase().trim();
    if (q) {
      if (q.includes("decision = 'deny'") || q === "deny" || q === "denied") {
        list = list.filter((item) => item.policy_decision === "DENY");
      } else if (q.includes("decision = 'allow'") || q === "allow" || q === "allowed") {
        list = list.filter((item) => item.policy_decision === "ALLOW");
      } else if (q.includes("decision = 'approval'") || q === "approval") {
        list = list.filter((item) => item.policy_decision === "REQUIRE_APPROVAL");
      } else {
        list = list.filter(
          (item) =>
            item.client.toLowerCase().includes(q) ||
            item.sql_normalized.toLowerCase().includes(q) ||
            item.token_id.toLowerCase().includes(q) ||
            item.connection.toLowerCase().includes(q) ||
            item.policy_decision.toLowerCase().includes(q),
        );
      }
    }

    const o = orderByFilter.trim().toLowerCase();
    if (o) {
      const isDesc = o.includes("desc");
      if (o.includes("time") || o.includes("timestamp") || o.includes("ts")) {
        list.sort((a, b) =>
          isDesc
            ? (b.ts || "").localeCompare(a.ts || "")
            : (a.ts || "").localeCompare(b.ts || ""),
        );
      } else if (o.includes("client")) {
        list.sort((a, b) =>
          isDesc ? b.client.localeCompare(a.client) : a.client.localeCompare(b.client),
        );
      } else if (o.includes("conn") || o.includes("connection")) {
        list.sort((a, b) =>
          isDesc
            ? b.connection.localeCompare(a.connection)
            : a.connection.localeCompare(b.connection),
        );
      } else if (o.includes("decision")) {
        list.sort((a, b) =>
          isDesc
            ? b.policy_decision.localeCompare(a.policy_decision)
            : a.policy_decision.localeCompare(b.policy_decision),
        );
      }
    }

    return list;
  });

  function toggleSort(col: string) {
    if (orderByFilter.startsWith(col)) {
      orderByFilter = orderByFilter.toUpperCase().includes("DESC")
        ? `${col} ASC`
        : `${col} DESC`;
    } else {
      orderByFilter = `${col} ASC`;
    }
  }

  function handleClearFilters() {
    filterStatus = "ALL";
    searchQuery = "";
    orderByFilter = "";
  }

  async function handleReload() {
    try {
      await auditState.loadAudit();
      toast.info("Audit Refreshed", "Loaded latest gateway audit records.");
    } catch (e) {
      toast.error("Reload Failed", e instanceof Error ? e.message : String(e));
    }
  }

  async function handleExport() {
    if (exportFormat === "JSON") {
      const dataStr =
        "data:text/json;charset=utf-8," +
        encodeURIComponent(JSON.stringify(filteredAuditLogs, null, 2));
      const downloadAnchor = document.createElement("a");
      downloadAnchor.setAttribute("href", dataStr);
      downloadAnchor.setAttribute("download", `audit_logs_${Date.now()}.json`);
      downloadAnchor.click();
      toast.success("Audit Exported", `Exported ${filteredAuditLogs.length} audit logs to JSON`);
      return;
    }
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
  {#if bootState.dataLoading}
    <Box style="padding: 8px 16px;">
      <Text size="sm" color="muted">Loading live data from Admin API…</Text>
    </Box>
  {:else if bootState.dataError}
    <Box style="padding: 8px 16px;">
      <Alert variant="danger" title="Admin API unreachable">
        <Text color="danger">{bootState.dataError}</Text>
        <Button
          variant="secondary"
          size="sm"
          onclick={() => void bootState.loadAll()}
        >
          Retry
        </Button>
      </Alert>
    </Box>
  {/if}

  <!-- Toolbar 1: Actions -->
  <Toolbar>
    {#snippet left()}
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Reload audit records (Cmd+Enter)"
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
      <ToolbarSeparator />
      <BorderlessSelect
        options={[
          { value: "ALL", label: "Decision: All Decisions" },
          { value: "ALLOW", label: "Decision: ALLOW" },
          { value: "DENY", label: "Decision: DENY" },
          { value: "REQUIRE_APPROVAL", label: "Decision: APPROVAL" },
        ]}
        bind:value={filterStatus}
      />
      <ToolbarSeparator />
      <Text size="xs" color="muted" class="tx-selector" title="Active Filter Count">
        {filteredAuditLogs.length} audit record(s)
      </Text>
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
        title="Export audit log"
        onclick={handleExport}
      >
        <Icon name="download" size={13} />
      </Button>
      <ToolbarSeparator />
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Configure Governance Rules"
        onclick={() => (settingsState.settingsModalOpen = true)}
      >
        <Icon name="settings" size={13} />
      </Button>
    {/snippet}
  </Toolbar>

  <!-- Toolbar 2: WHERE & ORDER BY -->
  <FilterBar
    wherePlaceholder="e.g. decision = 'deny' or client name or query text"
    orderByPlaceholder="e.g. timestamp DESC or client ASC"
    bind:whereValue={searchQuery}
    bind:orderByValue={orderByFilter}
    bind:searchRef
  />

  <!-- Main Content Layout -->
  <Flex class="content-layout">
    <!-- Grid -->
    <Box class="table-container">
      <DataGrid style="height: 100%;">
        {#snippet overlay()}
          <FloatingRowCount count={filteredAuditLogs.length} />
        {/snippet}
        <DataGridHead>
          <DataGridRow>
            <DataGridRowNumHead />
            <DataGridHeadCell width="240px">
              <DataGridHeaderInner>
                <Icon name="clock" size={12} color="#56A8F5" />
                <Text size="xs">Timestamp</Text>
                <Button
                  variant="ghost"
                  size="icon-xs"
                  class="ohjanus-data-grid-header-action"
                  title="Sort by Timestamp"
                  onclick={() => toggleSort("timestamp")}
                >
                  <Icon name="sort" size={9} />
                </Button>
              </DataGridHeaderInner>
            </DataGridHeadCell>
            <DataGridHeadCell width="110px">
              <DataGridHeaderInner>
                <Icon name="shield" size={12} color="#EDA200" />
                <Text size="xs">Decision</Text>
                <Button
                  variant="ghost"
                  size="icon-xs"
                  class="ohjanus-data-grid-header-action"
                  title="Sort by Decision"
                  onclick={() => toggleSort("decision")}
                >
                  <Icon name="sort" size={9} />
                </Button>
              </DataGridHeaderInner>
            </DataGridHeadCell>
            <DataGridHeadCell width="170px">
              <DataGridHeaderInner>
                <Icon name="user" size={12} color="#7A7E85" />
                <Text size="xs">Client Agent</Text>
                <Button
                  variant="ghost"
                  size="icon-xs"
                  class="ohjanus-data-grid-header-action"
                  title="Sort by Client"
                  onclick={() => toggleSort("client")}
                >
                  <Icon name="sort" size={9} />
                </Button>
              </DataGridHeaderInner>
            </DataGridHeadCell>
            <DataGridHeadCell width="150px">
              <DataGridHeaderInner>
                <Icon name="database" size={12} color="#3B82F6" />
                <Text size="xs">Connection</Text>
                <Button
                  variant="ghost"
                  size="icon-xs"
                  class="ohjanus-data-grid-header-action"
                  title="Sort by Connection"
                  onclick={() => toggleSort("connection")}
                >
                  <Icon name="sort" size={9} />
                </Button>
              </DataGridHeaderInner>
            </DataGridHeadCell>
            <DataGridHeadCell width="90px">
              <DataGridHeaderInner>
                <Icon name="terminal" size={12} color="#9DA0A8" />
                <Text size="xs">Type</Text>
              </DataGridHeaderInner>
            </DataGridHeadCell>
            <DataGridHeadCell width="140px">
              <DataGridHeaderInner>
                <Icon name="table" size={12} color="#3B82F6" />
                <Text size="xs">Tables</Text>
              </DataGridHeaderInner>
            </DataGridHeadCell>
            <DataGridHeadCell width="70px">
              <DataGridHeaderInner>
                <Icon name="chart" size={12} color="#7A7E85" />
                <Text size="xs">Rows</Text>
              </DataGridHeaderInner>
            </DataGridHeadCell>
            <DataGridHeadCell width="90px">
              <DataGridHeaderInner>
                <Icon name="lightning" size={12} color="#FACC15" />
                <Text size="xs">Duration</Text>
              </DataGridHeaderInner>
            </DataGridHeadCell>
            <DataGridHeadCell>
              <DataGridHeaderInner>
                <Text size="xs">Normalized SQL</Text>
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
              <DataGridCell class="cell-client" truncate title={item.client}>
                <Flex align="center" gap="xs">
                  <Icon name="user" size={12} />
                  <Text size="sm" truncate>{item.client}</Text>
                </Flex>
              </DataGridCell>
              <DataGridCell class="cell-conn">{item.connection}</DataGridCell>
              <DataGridCell class="cell-type">{item.statement_type}</DataGridCell>
              <DataGridCell
                class="cell-tables"
                title={(item.tables ?? []).join(", ")}
              >
                <Flex align="center" gap="xs">
                  {#if item.tables && item.tables.length > 0}
                    {#each item.tables.slice(0, 3) as t}
                      <Badge variant="default" size="sm" class="tbl-pill">{t}</Badge>
                    {/each}
                    {#if item.tables.length > 3}
                      <Text class="more-tbl" size="sm" color="muted">+{item.tables.length - 3}</Text>
                    {/if}
                  {:else}
                    <Text size="xs" color="muted">—</Text>
                  {/if}
                </Flex>
              </DataGridCell>
              <DataGridCell tone="number" align="right">
                {item.row_count}
              </DataGridCell>
              <DataGridCell tone="number" align="right">
                {item.duration_ms} ms
              </DataGridCell>
              <DataGridCell
                tone="secondary"
                truncate
                title={item.sql_normalized}
              >
                {item.sql_normalized}
              </DataGridCell>
            </DataGridRow>
          {/each}
        </DataGridBody>
      </DataGrid>
    </Box>

    <!-- Right Drawer Detail -->
    <AuditInspectionDrawer
      item={selectedAudit}
      onclose={() => (selectedAudit = null)}
    />
  </Flex>
</Box>

<style>
  :global(.audit-view) {
    display: flex;
    flex-direction: column;
    height: 100%;
    background-color: var(--bg-canvas);
    color: var(--text-primary);
    overflow: hidden;
  }

  :global(.audit-view .tx-selector) {
    font-size: var(--font-size-xs, 12px);
    color: var(--text-muted);
    font-family: var(--font-code);
    padding: 0 4px;
    white-space: nowrap;
  }

  :global(.audit-view .content-layout) {
    flex: 1;
    overflow: hidden;
    position: relative;
    display: flex;
  }

  :global(.audit-view .table-container) {
    flex: 1;
    overflow: auto;
    position: relative;
  }

  :global(.audit-view .cell-ts) {
    font-size: var(--font-size-xs, 12px);
    color: var(--text-muted);
    font-family: var(--font-code);
  }

  :global(.audit-view .cell-client) {
    font-size: var(--font-size-sm, 13px);
  }

  :global(.audit-view .cell-conn) {
    font-family: var(--font-code);
    font-size: var(--font-size-sm, 13px);
  }

  :global(.audit-view .cell-type) {
    font-size: var(--font-size-xs, 12px);
    font-weight: 500;
  }

  :global(.audit-view .tbl-pill) {
    font-family: var(--font-code);
    font-size: 10px;
    padding: 1px 4px;
  }

  :global(.audit-view .more-tbl) {
    font-size: 11px;
    margin-left: 2px;
  }

  :global(.audit-view .no-tbl) {
    color: var(--text-muted);
  }
</style>
