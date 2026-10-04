<script lang="ts">
  import { appState } from "../../state/appState.svelte";
  import {
    Card,
    Badge,
    Button,
    Select,
    Alert,
    Text,
    toast,
    Box,
    Flex,
    Grid,
    Stack,
  } from "@ohjanus/ui";
  import { Icon } from "@ohjanus/icons";

  let testingConn = $state<string | null>(null);
  let selectedMode = $state<string>("ALL");
  let whereFilter = $state("");
  let orderByFilter = $state("");
  let searchRef = $state<HTMLInputElement>();

  let filteredConnections = $derived.by(() => {
    let list = appState.connections;
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
      const res = await appState.pingConnection(name);
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
    if (testingConn || appState.connections.length === 0) return;
    for (const c of appState.connections) {
      await runTest(c.name);
    }
  }

  async function handleReload() {
    try {
      await appState.loadConnections();
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
        disabled={testingConn !== null || appState.connections.length === 0}
        onclick={runTestAll}
      >
        <Icon name="lightning" size={13} />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Configure data sources in Settings"
        onclick={() => (appState.settingsModalOpen = true)}
      >
        <Icon name="plus" size={13} color="#57D38C" />
      </Button>
      <span class="bar-separator"></span>
      <div class="borderless-select-wrapper">
        <Select
          class="toolbar-select borderless-select"
          options={[
            { value: "ALL", label: "Mode: All Data Sources" },
            { value: "READWRITE", label: "Mode: Read-Write" },
            { value: "READONLY", label: "Mode: Read-Only" },
          ]}
          bind:value={selectedMode}
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
        title="Search data sources"
        onclick={() => searchRef?.focus()}
      >
        <Icon name="search" size={13} />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Open SQL Console"
        onclick={() => (appState.activeTabId = "console")}
      >
        <Icon name="terminal" size={13} />
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
        placeholder="e.g. name = 'PRD' or driver = 'postgres'"
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
        placeholder="e.g. latency_ms ASC or name ASC"
        bind:value={orderByFilter}
      />
    </div>
  </div>

  <!-- Cards Grid -->
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
      {#each filteredConnections as conn}
      <Card class="conn-card-item">
        {#snippet header()}
          <Flex class="card-header-inner" justify="between" align="center">
            <Flex class="header-left" align="center" gap="8px">
              <span class="db-icon"
                ><Icon name="database" size={15} color="#3B82F6" /></span
              >
              <Text class="conn-title" weight="semibold">{conn.name}</Text>
              {#if conn.readonly}
                <Badge variant="warning" size="sm">
                  <Icon name="lock" size={11} />
                  <Text size="xs" color="warning">READ-ONLY</Text>
                </Badge>
              {:else}
                <Badge variant="info" size="sm">
                  <Icon name="lightning" size={11} />
                  <Text size="xs" style="color: #56A8F5;">READ-WRITE</Text>
                </Badge>
              {/if}
            </Flex>
            {#if conn.status === "healthy"}
              <Badge variant="success" size="sm"
                ><span class="status-dot ok" aria-hidden="true"
                ></span>HEALTHY</Badge
              >
            {:else}
              <Badge variant="danger" size="sm"
                ><span class="status-dot bad" aria-hidden="true"
                ></span>DEGRADED</Badge
              >
            {/if}
          </Flex>
        {/snippet}

        <Stack class="conn-details" gap="6px">
          <Flex class="detail-row" justify="between">
            <Text class="label" size="md" color="muted">Driver:</Text>
            <Text class="value" size="md" mono>{conn.driver}</Text>
          </Flex>
          <Flex class="detail-row" justify="between">
            <Text class="label" size="md" color="muted">Last Health Ping:</Text>
            <Text class="value" size="md"
              >{conn.last_ping_at}{conn.latency_ms !== undefined
                ? ` (${conn.latency_ms} ms)`
                : ""}</Text
            >
          </Flex>
        </Stack>

        <!-- Security Guardrails -->
        <Box class="guardrails-section">
          <Box class="sec-label">ALLOWED SCHEMAS:</Box>
          <Flex class="tags-list" wrap gap="4px">
            {#each conn.allowed_schemas as sch}
              <Badge variant="default" size="sm">{sch}</Badge>
            {/each}
          </Flex>

          {#if conn.denied_tables.length > 0}
            <Box
              class="sec-label"
              style="margin-top: 6px; color: var(--action-danger);"
              >DENIED TABLES:</Box
            >
            <Flex class="tags-list" wrap gap="4px">
              {#each conn.denied_tables as dt}
                <Badge variant="danger" size="sm">{dt}</Badge>
              {/each}
            </Flex>
          {/if}
        </Box>

        {#snippet footer()}
          <Button
            variant="secondary"
            size="sm"
            onclick={() => runTest(conn.name)}
            loading={testingConn === conn.name}
          >
            {#if testingConn === conn.name}
              Pinging...
            {:else}
              <Icon name="refresh" size={12} />
              <Text size="md">Test Connection</Text>
            {/if}
          </Button>
        {/snippet}
      </Card>
    {/each}
  </Grid>
  {/if}
</Box>

<style>
  :global(.connections-view) {
    display: flex;
    flex-direction: column;
    height: 100%;
    background: var(--bg-canvas);
    overflow-y: auto;
  }

  :global(.connections-view .table-toolbar) {
    height: var(--toolbar-height, 32px);
    background-color: var(--bg-toolbar);
    border-bottom: 1px solid var(--border-subtle);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 8px;
    flex-shrink: 0;
  }

  :global(.connections-view .toolbar-left),
  :global(.connections-view .toolbar-right) {
    display: flex;
    align-items: center;
    gap: 2px;
  }

  :global(.connections-view .bar-separator) {
    width: 1px;
    height: 14px;
    background-color: var(--border-default);
    margin: 0 4px;
  }

  :global(.connections-view .borderless-select-wrapper) {
    display: inline-flex;
    align-items: center;
  }

  :global(.connections-view .borderless-select) {
    min-width: unset !important;
    width: auto !important;
  }

  :global(.connections-view .borderless-select .ohjanus-select-trigger) {
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

  :global(.connections-view .borderless-select .ohjanus-select-trigger:hover) {
    background-color: var(--bg-hover, #313438) !important;
    color: var(--text-primary, #DFE1E5) !important;
  }

  :global(.connections-view .ddl-btn) {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 0 4px !important;
    width: auto !important;
  }

  :global(.connections-view .filter-bar) {
    height: var(--filterbar-height, 30px);
    background-color: var(--bg-canvas);
    border-bottom: 1px solid var(--border-subtle);
    display: flex;
    align-items: center;
    padding: 0 8px;
    gap: 12px;
    flex-shrink: 0;
  }

  :global(.connections-view .filter-group) {
    display: flex;
    align-items: center;
    gap: 6px;
  }

  :global(.connections-view .where-group) {
    flex: 1.1;
  }

  :global(.connections-view .orderby-group) {
    flex: 0.9;
  }

  :global(.connections-view .filter-icon) {
    color: var(--text-muted);
    font-size: 11px;
    user-select: none;
    display: inline-flex;
    align-items: center;
  }

  :global(.connections-view .filter-input) {
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

  :global(.connections-view .filter-input:focus),
  :global(.connections-view .filter-input:hover) {
    border: none !important;
    outline: none !important;
    box-shadow: none !important;
  }

  :global(.connections-view .cards-grid) {
    padding: 16px;
  }

  :global(.card-header-inner) {
    width: 100%;
  }

  :global(.connections-view .db-icon) {
    font-size: 18px;
  }
  :global(.connections-view .conn-title) {
    font-weight: 600;
    font-size: var(--font-size-md, 15px);
    color: var(--text-primary);
  }

  :global(.conn-details) {
    margin-bottom: 14px;
    font-size: var(--font-size-sm, 13px);
  }

  :global(.detail-row .label) {
    color: var(--text-muted);
  }
  :global(.detail-row .value) {
    color: var(--text-primary);
  }

  :global(.pool-section) {
    background: rgba(0, 0, 0, 0.2);
    border: 1px solid var(--border-subtle);
    border-radius: 4px;
    padding: 10px;
    margin-bottom: 12px;
  }

  :global(.pool-title) {
    font-size: var(--font-size-2xs, 11px);
    font-weight: 600;
    color: var(--text-muted);
    letter-spacing: 0.5px;
    margin-bottom: 8px;
  }

  :global(.pool-stats) {
    text-align: center;
  }

  :global(.stat-box .num) {
    display: block;
    font-size: var(--font-size-lg, 16px);
    font-weight: 700;
    color: var(--action-primary);
  }

  :global(.stat-box .lbl) {
    font-size: var(--font-size-2xs, 11px);
    color: var(--text-muted);
  }

  :global(.guardrails-section) {
    margin-bottom: 8px;
  }

  :global(.sec-label) {
    font-size: var(--font-size-2xs, 11px);
    font-weight: 600;
    color: var(--text-muted);
    margin-bottom: 4px;
  }
</style>
