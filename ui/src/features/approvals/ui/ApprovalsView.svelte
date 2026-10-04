<script lang="ts">
  import type { ApprovalRequest } from '@/entities/approval';
  import { approvalsState } from '@/features/approvals';
  import { settingsState } from '@/features/settings';
  import {
    Button,
    Alert,
    EmptyState,
    toast,
    Box,
    Flex,
    Text,
  } from "@ohjanus/ui";
  import { Icon } from "@ohjanus/icons";
  import {
    Toolbar,
    ToolbarSeparator,
    BorderlessSelect,
    FilterBar,
  } from '@/shared/ui/toolbar';
  import ApprovalCard from './ApprovalCard.svelte';
  import ApprovalInspectionDrawer from './ApprovalInspectionDrawer.svelte';

  let selectedFilter = $state<string>("pending");
  let whereFilter = $state("");
  let orderByFilter = $state("");
  let searchRef = $state<HTMLInputElement>();
  let selectedDetail = $state<ApprovalRequest | null>(null);

  let filteredApprovals = $derived.by(() => {
    let list = approvalsState.approvals;
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

  function openActionModal(item: ApprovalRequest, mode: "approve" | "reject") {
    approvalsState.selectedApprovalForAction = item;
    approvalsState.approvalDecisionMode = mode;
    approvalsState.approvalDecisionReason =
      mode === "approve" ? "Approved for execution" : "";
  }

  async function handleReload() {
    try {
      await approvalsState.loadApprovals();
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
  {#if approvalsState.loading}
    <Box style="padding: 8px 16px;">
      <Text size="sm" color="muted">Loading live data from Admin API…</Text>
    </Box>
  {:else if approvalsState.error}
    <Box style="padding: 8px 16px;">
      <Alert variant="danger" title="Admin API unreachable">
        <Text color="danger">{approvalsState.error}</Text>
        <Button
          variant="secondary"
          size="sm"
          onclick={handleReload}
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
      <ToolbarSeparator />
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Batch Approve All Pending"
        disabled={filteredApprovals.filter(a => a.state === 'pending').length === 0}
        onclick={() => {
          const first = filteredApprovals.find(a => a.state === 'pending');
          if (first) openActionModal(first, 'approve');
        }}
      >
        <Icon name="check" size={13} color="#57D38C" />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Batch Reject"
        disabled={filteredApprovals.filter(a => a.state === 'pending').length === 0}
        onclick={() => {
          const first = filteredApprovals.find(a => a.state === 'pending');
          if (first) openActionModal(first, 'reject');
        }}
      >
        <Icon name="x" size={13} color="#E55353" />
      </Button>
      <ToolbarSeparator />
      <BorderlessSelect
        options={[
          { value: "pending", label: "Queue: Pending Review" },
          { value: "approved", label: "Queue: Approved" },
          { value: "rejected", label: "Queue: Rejected" },
          { value: "all", label: "Queue: All States" },
        ]}
        bind:value={selectedFilter}
      />
      <Text size="xs" color="muted" class="tx-selector" title="Active Filter Count">
        {filteredApprovals.length} request(s)
      </Text>
    {/snippet}

    {#snippet right()}
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
    wherePlaceholder="e.g. statement_type = 'DROP' or connection = 'main-db'"
    orderByPlaceholder="e.g. created_at DESC or connection ASC"
    bind:whereValue={whereFilter}
    bind:orderByValue={orderByFilter}
    bind:searchRef
  />

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
          <ApprovalCard
            {item}
            isSelected={selectedDetail?.id === item.id}
            onselect={() => (selectedDetail = item)}
            onaction={(mode) => openActionModal(item, mode)}
          />
        {/each}
      {/if}
    </Box>

    <!-- Right: Detail Inspection Drawer -->
    <ApprovalInspectionDrawer
      item={selectedDetail}
      onclose={() => (selectedDetail = null)}
      onaction={(mode) => {
        if (selectedDetail) {
          openActionModal(selectedDetail, mode);
        }
      }}
    />
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

  :global(.approvals-view .tx-selector) {
    font-size: var(--font-size-xs, 12px);
    color: var(--text-muted);
    font-family: var(--font-code);
    padding: 0 4px;
    white-space: nowrap;
  }

  :global(.approvals-view .content-layout) {
    flex: 1;
    overflow: hidden;
    position: relative;
  }

  :global(.approvals-view .queue-list) {
    flex: 1;
    overflow-y: auto;
    padding: 12px;
    display: flex;
    flex-direction: column;
    gap: 10px;
  }
</style>
