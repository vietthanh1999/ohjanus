<script lang="ts">
  import { tokensState } from "@/features/tokens";
  import { settingsState } from "@/features/settings";
  import {
    Button,
    Badge,
    Alert,
    Text,
    toast,
    Box,
    Flex,
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
  } from "@/shared/ui/toolbar";

  let selectedRowIndex = $state(0);
  let whereFilter = $state("");
  let orderByFilter = $state("");
  let exportFormat = $state("CSV");
  let txMode = $state("all");
  let searchInputRef = $state<HTMLInputElement>();

  let filteredTokens = $derived.by(() => {
    let list = [...tokensState.tokens];

    // Filter by WHERE
    const w = whereFilter.trim().toLowerCase();
    if (w) {
      if (w.includes("state = 'active'") || w === "active") {
        list = list.filter((t) => t.state === "active");
      } else if (w.includes("state = 'revoked'") || w === "revoked") {
        list = list.filter((t) => t.state === "revoked");
      } else {
        list = list.filter(
          (t) =>
            t.id.toLowerCase().includes(w) ||
            t.name.toLowerCase().includes(w) ||
            t.scopes.some((s) => s.toLowerCase().includes(w)) ||
            t.state.toLowerCase().includes(w) ||
            (t.created_at && t.created_at.toLowerCase().includes(w)) ||
            (t.expires_at && t.expires_at.toLowerCase().includes(w)),
        );
      }
    }

    // Filter by Tx mode if specified
    if (txMode === "active") {
      list = list.filter((t) => t.state === "active");
    } else if (txMode === "revoked") {
      list = list.filter((t) => t.state === "revoked");
    }

    // Sort by ORDER BY
    const o = orderByFilter.trim().toLowerCase();
    if (o) {
      const isDesc = o.includes("desc");
      if (o.includes("name")) {
        list.sort((a, b) =>
          isDesc ? b.name.localeCompare(a.name) : a.name.localeCompare(b.name),
        );
      } else if (o.includes("created") || o.includes("created_at")) {
        list.sort((a, b) =>
          isDesc
            ? (b.created_at || "").localeCompare(a.created_at || "")
            : (a.created_at || "").localeCompare(b.created_at || ""),
        );
      } else if (o.includes("expires") || o.includes("expires_at")) {
        list.sort((a, b) =>
          isDesc
            ? (b.expires_at || "").localeCompare(a.expires_at || "")
            : (a.expires_at || "").localeCompare(b.expires_at || ""),
        );
      } else if (o.includes("state")) {
        list.sort((a, b) =>
          isDesc
            ? b.state.localeCompare(a.state)
            : a.state.localeCompare(b.state),
        );
      } else if (o.includes("id")) {
        list.sort((a, b) =>
          isDesc ? b.id.localeCompare(a.id) : a.id.localeCompare(b.id),
        );
      }
    }

    return list;
  });

  let selectedToken = $derived(filteredTokens[selectedRowIndex] ?? null);

  function toggleSort(col: string) {
    if (orderByFilter.startsWith(col)) {
      orderByFilter = orderByFilter.toUpperCase().includes("DESC")
        ? `${col} ASC`
        : `${col} DESC`;
    } else {
      orderByFilter = `${col} ASC`;
    }
  }

  async function handleReload() {
    try {
      await tokensState.loadTokens();
      toast.info("Tokens Refreshed", "Loaded latest MCP client tokens.");
    } catch (e) {
      toast.error("Reload Failed", e instanceof Error ? e.message : String(e));
    }
  }

  async function handleRevokeSelected() {
    if (!selectedToken || selectedToken.state !== "active") return;
    await handleRevoke(selectedToken.id);
  }

  function handleExport() {
    if (filteredTokens.length === 0) {
      toast.info("Export", "No tokens to export.");
      return;
    }
    if (exportFormat === "JSON") {
      const dataStr =
        "data:text/json;charset=utf-8," +
        encodeURIComponent(JSON.stringify(filteredTokens, null, 2));
      const downloadAnchor = document.createElement("a");
      downloadAnchor.setAttribute("href", dataStr);
      downloadAnchor.setAttribute("download", `mcp_tokens_${Date.now()}.json`);
      downloadAnchor.click();
      toast.success(
        "Tokens Exported",
        `Exported ${filteredTokens.length} tokens to JSON`,
      );
    } else {
      const headers = [
        "id",
        "name",
        "scopes",
        "created_at",
        "expires_at",
        "last_used_at",
        "state",
      ];
      const rows = filteredTokens.map((t) =>
        [
          t.id,
          `"${(t.name || "").replace(/"/g, '""')}"`,
          `"${t.scopes.join(",")}"`,
          t.created_at || "",
          t.expires_at || "",
          t.last_used_at || "",
          t.state,
        ].join(","),
      );
      const csvContent = [headers.join(","), ...rows].join("\n");
      const blob = new Blob([csvContent], { type: "text/csv;charset=utf-8;" });
      const url = URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.setAttribute("href", url);
      link.setAttribute("download", `mcp_tokens_${Date.now()}.csv`);
      link.click();
      URL.revokeObjectURL(url);
      toast.success(
        "Tokens Exported",
        `Exported ${filteredTokens.length} tokens to CSV`,
      );
    }
  }

  async function handleRevoke(id: string) {
    try {
      await tokensState.revokeToken(id);
      toast.error("Token Revoked", `Token ${id} has been revoked.`);
    } catch (e) {
      toast.error("Revoke Failed", e instanceof Error ? e.message : String(e));
    }
  }

  function handleClearFilters() {
    whereFilter = "";
    orderByFilter = "";
    txMode = "all";
  }
</script>

<Box class="tokens-view">
  {#if tokensState.loading}
    <Box style="padding: 8px 16px;">
      <Text size="sm" color="muted">Loading live data from Admin API…</Text>
    </Box>
  {:else if tokensState.error}
    <Box style="padding: 8px 16px;">
      <Alert variant="danger" title="Admin API unreachable">
        <Text color="danger">{tokensState.error}</Text>
        <Button variant="secondary" size="sm" onclick={handleReload}>
          Retry
        </Button>
      </Alert>
    </Box>
  {/if}

  <!-- Toolbar 1: Actions (matching design2.png) -->
  <Toolbar>
    {#snippet left()}
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Reload tokens (Cmd+Enter)"
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
        title="Generate New MCP Token (+)"
        onclick={() => (tokensState.createTokenModalOpen = true)}
      >
        <Icon name="plus" size={13} color="#57D38C" />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Revoke selected token (-)"
        disabled={!selectedToken || selectedToken.state !== "active"}
        onclick={handleRevokeSelected}
      >
        <Icon name="minus" size={13} />
      </Button>
      <ToolbarSeparator />
      <BorderlessSelect
        options={[
          { value: "all", label: "State: All Tokens" },
          { value: "active", label: "State: Active Only" },
          { value: "revoked", label: "State: Revoked Only" },
        ]}
        bind:value={txMode}
      />
      <ToolbarSeparator />
      <Text
        size="xs"
        color="muted"
        class="tx-selector"
        title="Active Filter Count"
      >
        {filteredTokens.length} token(s)
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
        title="Export tokens"
        onclick={handleExport}
      >
        <Icon name="download" size={13} />
      </Button>
      <ToolbarSeparator />
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Settings"
        onclick={() => (settingsState.settingsModalOpen = true)}
      >
        <Icon name="settings" size={13} />
      </Button>
    {/snippet}
  </Toolbar>

  <!-- Toolbar 2: WHERE & ORDER BY (matching design2.png) -->
  <FilterBar
    wherePlaceholder="e.g. state = 'active' or client name"
    orderByPlaceholder="e.g. created_at DESC"
    bind:whereValue={whereFilter}
    bind:orderByValue={orderByFilter}
    bind:searchRef={searchInputRef}
  />

  <!-- Tokens List Table -->
  <Box class="table-container">
    <DataGrid style="height: 100%;">
      {#snippet overlay()}
        <FloatingRowCount count={filteredTokens.length} />
      {/snippet}
      <DataGridHead>
        <DataGridRow>
          <DataGridRowNumHead />
          <DataGridHeadCell width="150px">
            <DataGridHeaderInner>
              <Icon name="key" size={12} color="#EDA200" />
              <Text size="md">Token ID</Text>
              <Button
                variant="ghost"
                size="icon-xs"
                class="ohjanus-data-grid-header-action"
                title="Sort by ID"
                onclick={() => toggleSort("id")}
              >
                <Icon name="sort" size={9} />
              </Button>
            </DataGridHeaderInner>
          </DataGridHeadCell>
          <DataGridHeadCell width="220px">
            <DataGridHeaderInner>
              <Icon name="user" size={12} color="#7A7E85" />
              <Text size="md">Agent / Client Name</Text>
              <Button
                variant="ghost"
                size="icon-xs"
                class="ohjanus-data-grid-header-action"
                title="Sort by Name"
                onclick={() => toggleSort("name")}
              >
                <Icon name="sort" size={9} />
              </Button>
            </DataGridHeaderInner>
          </DataGridHeadCell>
          <DataGridHeadCell>
            <DataGridHeaderInner>
              <Icon name="shield" size={12} color="#3B82F6" />
              <Text size="md">Scopes</Text>
            </DataGridHeaderInner>
          </DataGridHeadCell>
          <DataGridHeadCell width="120px">
            <DataGridHeaderInner>
              <Icon name="clock" size={12} color="#56A8F5" />
              <Text size="md">Created</Text>
              <Button
                variant="ghost"
                size="icon-xs"
                class="ohjanus-data-grid-header-action"
                title="Sort by Created"
                onclick={() => toggleSort("created_at")}
              >
                <Icon name="sort" size={9} />
              </Button>
            </DataGridHeaderInner>
          </DataGridHeadCell>
          <DataGridHeadCell width="170px">
            <DataGridHeaderInner>
              <Icon name="clock" size={12} color="#56A8F5" />
              <Text size="md">Expires</Text>
              <Button
                variant="ghost"
                size="icon-xs"
                class="ohjanus-data-grid-header-action"
                title="Sort by Expires"
                onclick={() => toggleSort("expires_at")}
              >
                <Icon name="sort" size={9} />
              </Button>
            </DataGridHeaderInner>
          </DataGridHeadCell>
          <DataGridHeadCell width="120px">
            <DataGridHeaderInner>
              <Icon name="clock" size={12} color="#56A8F5" />
              <Text size="md">Last Used</Text>
            </DataGridHeaderInner>
          </DataGridHeadCell>
          <DataGridHeadCell width="100px">
            <DataGridHeaderInner>
              <Icon name="chart" size={12} color="#57D38C" />
              <Text size="md">State</Text>
              <Button
                variant="ghost"
                size="icon-xs"
                class="ohjanus-data-grid-header-action"
                title="Sort by State"
                onclick={() => toggleSort("state")}
              >
                <Icon name="sort" size={9} />
              </Button>
            </DataGridHeaderInner>
          </DataGridHeadCell>
          <DataGridHeadCell width="90px">
            <DataGridHeaderInner style="justify-content: flex-end;">
              <Text size="md">Actions</Text>
            </DataGridHeaderInner>
          </DataGridHeadCell>
        </DataGridRow>
      </DataGridHead>
      <DataGridBody>
        {#if filteredTokens.length === 0}
          <DataGridRow>
            <DataGridRowNum index={0} />
            <DataGridCell tone="secondary">
              {#if tokensState.tokens.length === 0}
                No tokens generated yet. Click '+' in toolbar to create one.
              {:else}
                No tokens match the filter. Clear WHERE filter to show all.
              {/if}
            </DataGridCell>
          </DataGridRow>
        {:else}
          {#each filteredTokens as tok, idx (tok.id)}
            <DataGridRow
              selected={selectedRowIndex === idx}
              onclick={() => (selectedRowIndex = idx)}
            >
              <DataGridRowNum index={idx} />
              <DataGridCell class="id-cell" truncate title={tok.id}
                >{tok.id}</DataGridCell
              >
              <DataGridCell class="name-cell" truncate title={tok.name}>
                <Text weight="bold">{tok.name}</Text>
              </DataGridCell>
              <DataGridCell class="scopes-cell">
                <Flex align="center" gap="xs">
                  {#each tok.scopes as sc}
                    <Badge variant="default" size="md">
                      {sc}
                    </Badge>
                  {/each}
                </Flex>
              </DataGridCell>
              <DataGridCell tone="secondary"
                >{tok.created_at || "—"}</DataGridCell
              >
              <DataGridCell tone="secondary">{tok.expires_at}</DataGridCell>
              <DataGridCell tone="secondary"
                >{tok.last_used_at || "—"}</DataGridCell
              >
              <DataGridCell>
                {#if tok.state === "active"}
                  <Badge variant="success" size="md">ACTIVE</Badge>
                {:else}
                  <Badge variant="danger" size="md">REVOKED</Badge>
                {/if}
              </DataGridCell>
              <DataGridCell align="right">
                {#if tok.state === "active"}
                  <Button
                    variant="danger"
                    size="xs"
                    onclick={() => handleRevoke(tok.id)}
                  >
                    Revoke
                  </Button>
                {:else}
                  <Text
                    color="muted"
                    style="color: var(--text-muted); font-size: var(--font-size-xs, 12px);"
                  >
                    Revoked
                  </Text>
                {/if}
              </DataGridCell>
            </DataGridRow>
          {/each}
        {/if}
      </DataGridBody>
    </DataGrid>
  </Box>
</Box>

<style>
  :global(.tokens-view) {
    display: flex;
    flex-direction: column;
    height: 100%;
    background-color: var(--bg-canvas);
    color: var(--text-primary);
    overflow: hidden;
  }

  :global(.tokens-view .tx-selector) {
    font-size: var(--font-size-xs, 12px);
    color: var(--text-muted);
    font-family: var(--font-code);
    padding: 0 4px;
    white-space: nowrap;
  }

  :global(.tokens-view .table-container) {
    flex: 1;
    overflow: auto;
    position: relative;
  }

  :global(.tokens-view .id-cell) {
    color: var(--syntax-number, #6897bb);
  }
</style>
