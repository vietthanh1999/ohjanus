<script lang="ts">
  import { connectionsState } from '@/features/connections';
  import { settingsState } from '@/features/settings';
  import { bootState } from '@/app/boot.svelte';
  import {
    Button,
    Alert,
    Text,
    toast,
    Box,
    Grid,
  } from "@ohjanus/ui";
  import { Icon } from "@ohjanus/icons";
  import {
    Toolbar,
    ToolbarSeparator,
    BorderlessSelect,
    FilterBar,
  } from '@/shared/ui/toolbar';
  import ConnectionCard from './ConnectionCard.svelte';

  let testingConn = $state<string | null>(null);
  let selectedMode = $state<string>("ALL");
  let whereFilter = $state("");
  let orderByFilter = $state("");
  let searchRef = $state<HTMLInputElement>();

  let filteredConnections = $derived.by(() => {
    let list = connectionsState.connections;
    if (selectedMode === "READONLY") {
      list = list.filter((c) => c.readonly);
    } else if (selectedMode === "READWRITE") {
      list = list.filter((c) => !c.readonly);
    }

    const w = whereFilter.trim().toLowerCase();
    if (w) {
      list = list.filter(
        (c) =>
          c.name.toLowerCase().includes(w) ||
          c.driver.toLowerCase().includes(w) ||
          c.status.toLowerCase().includes(w) ||
          c.allowed_schemas.some((s) => s.toLowerCase().includes(w)),
      );
    }

    const o = orderByFilter.trim().toLowerCase();
    if (o) {
      const isDesc = o.includes("desc");
      if (o.includes("latency")) {
        list = [...list].sort((a, b) =>
          isDesc
            ? (b.latency_ms ?? 0) - (a.latency_ms ?? 0)
            : (a.latency_ms ?? 0) - (b.latency_ms ?? 0),
        );
      } else if (o.includes("name")) {
        list = [...list].sort((a, b) =>
          isDesc ? b.name.localeCompare(a.name) : a.name.localeCompare(b.name),
        );
      } else if (o.includes("status")) {
        list = [...list].sort((a, b) =>
          isDesc
            ? b.status.localeCompare(a.status)
            : a.status.localeCompare(b.status),
        );
      }
    }
    return list;
  });

  async function runTest(name: string) {
    if (testingConn) return;
    testingConn = name;
    try {
      const res = await connectionsState.pingConnection(name);
      toast.success(
        "Connection Ping Healthy",
        `"${name}" responded in ${res.latency_ms} ms.`,
      );
    } catch (e) {
      toast.error(
        "Connection Test Failed",
        e instanceof Error ? e.message : String(e),
      );
    } finally {
      testingConn = null;
    }
  }

  async function runTestAll() {
    if (testingConn || connectionsState.connections.length === 0) return;
    for (const c of connectionsState.connections) {
      await runTest(c.name);
    }
  }

  async function handleReload() {
    try {
      await connectionsState.loadConnections();
      toast.info("Connections Refreshed", "Reloaded active connection pools.");
    } catch (e) {
      toast.error("Reload Failed", e instanceof Error ? e.message : String(e));
    }
  }

  function handleClearFilters() {
    selectedMode = "ALL";
    whereFilter = "";
    orderByFilter = "";
  }
</script>

<Box class="connections-view">
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
        title="Reload connections (Cmd+Enter)"
        onclick={handleReload}
      >
        <Icon name="refresh" size={13} />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Ping all data sources"
        disabled={testingConn !== null || connectionsState.connections.length === 0}
        onclick={runTestAll}
      >
        <Icon name="lightning" size={13} />
      </Button>
      <ToolbarSeparator />
      <BorderlessSelect
        options={[
          { value: "ALL", label: "Mode: All Connections" },
          { value: "READONLY", label: "Mode: Read-Only" },
          { value: "READWRITE", label: "Mode: Read-Write" },
        ]}
        bind:value={selectedMode}
      />
      <ToolbarSeparator />
      <Text size="xs" color="muted" class="tx-selector" title="Active Filter Count">
        {filteredConnections.length} pool(s)
      </Text>
    {/snippet}

    {#snippet right()}
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Configure data sources"
        onclick={() => (settingsState.settingsModalOpen = true)}
      >
        <Icon name="settings" size={13} />
      </Button>
    {/snippet}
  </Toolbar>

  <!-- Toolbar 2: WHERE & ORDER BY -->
  <FilterBar
    wherePlaceholder="e.g. driver = 'postgres' or name = 'prod-db'"
    orderByPlaceholder="e.g. latency_ms ASC or name ASC"
    bind:whereValue={whereFilter}
    bind:orderByValue={orderByFilter}
    bind:searchRef
  />

  <!-- Cards Grid -->
  <Box class="cards-container">
    {#if filteredConnections.length === 0}
      <Box style="padding: 24px;">
        <Text color="muted" size="sm">No database data sources match the current filter.</Text>
      </Box>
    {:else}
      <Grid
        class="cards-grid"
        columns="repeat(auto-fill, minmax(340px, 1fr))"
        gap="16px"
      >
        {#each filteredConnections as conn (conn.name)}
          <ConnectionCard
            {conn}
            isTesting={testingConn === conn.name}
            ontest={runTest}
          />
        {/each}
      </Grid>
    {/if}
  </Box>
</Box>

<style>
  :global(.connections-view) {
    display: flex;
    flex-direction: column;
    height: 100%;
    background: var(--bg-canvas);
    overflow: hidden;
  }

  :global(.connections-view .tx-selector) {
    font-size: var(--font-size-xs, 12px);
    color: var(--text-muted);
    font-family: var(--font-code);
    padding: 0 4px;
    white-space: nowrap;
  }

  :global(.connections-view .cards-container) {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
  }
</style>
