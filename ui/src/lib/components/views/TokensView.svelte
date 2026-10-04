<script lang="ts">
  import { appState, type McpToken } from "../../state/appState.svelte";
  import {
    Button,
    Badge,
    Modal,
    Input,
    Select,
    Checkbox,
    Field,
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

  let showCreateModal = $state(false);
  let showSecretModal = $state(false);
  let generatedToken = $state<McpToken | null>(null);
  let copied = $state(false);
  let creating = $state(false);

  // Form State
  let tokenName = $state("");
  let ttlDays = $state(30);
  let scopes = $state<Record<string, boolean>>({
    read: true,
    write_preview: true,
    write_execute: false,
    admin: false,
  });

  async function handleCreate() {
    const selectedScopes = Object.keys(scopes).filter((k) => scopes[k]);
    if (!tokenName.trim() || creating) return;

    creating = true;
    try {
      const secret = await appState.createToken(
        tokenName.trim(),
        selectedScopes,
        ttlDays,
      );
      generatedToken = {
        id: "new",
        name: tokenName.trim(),
        scopes: selectedScopes,
        created_at: "",
        expires_at: "",
        last_used_at: "",
        state: "active",
        rawToken: secret,
      };
      showCreateModal = false;
      showSecretModal = true;
      tokenName = "";
      copied = false;
      toast.success("MCP Token Generated", "Opaque bearer secret created.");
    } catch (e) {
      toast.error(
        "Token Creation Failed",
        e instanceof Error ? e.message : String(e),
      );
    } finally {
      creating = false;
    }
  }

  function copySecret() {
    if (generatedToken?.rawToken) {
      navigator.clipboard.writeText(generatedToken.rawToken);
      copied = true;
      toast.info("Copied", "Token secret copied to clipboard.");
      setTimeout(() => (copied = false), 2000);
    }
  }

  let selectedRowIndex = $state(0);
  let whereFilter = $state("");
  let orderByFilter = $state("");
  let txMode = $state("auto");
  let exportFormat = $state("CSV");
  let searchInputRef = $state<HTMLInputElement>();

  let filteredTokens = $derived.by(() => {
    let list = [...appState.tokens];

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
          isDesc ? b.state.localeCompare(a.state) : a.state.localeCompare(b.state),
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
      await appState.loadTokens();
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
      toast.success("Tokens Exported", `Exported ${filteredTokens.length} tokens to JSON`);
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
      toast.success("Tokens Exported", `Exported ${filteredTokens.length} tokens to CSV`);
    }
  }

  async function handleRevoke(id: string) {
    try {
      await appState.revokeToken(id);
      toast.error("Token Revoked", `Token ${id} has been revoked.`);
    } catch (e) {
      toast.error("Revoke Failed", e instanceof Error ? e.message : String(e));
    }
  }
</script>

<Box class="tokens-view">
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
        title="Reload (Cmd+Enter)"
        onclick={handleReload}
      >
        <Icon name="refresh" size={13} />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Audit Trail"
        onclick={() => (appState.activeTabId = "audit")}
      >
        <Icon name="clock" size={13} />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Stop Execution"
        disabled
      >
        <Icon name="stop" size={13} color="#E55353" />
      </Button>
      <span class="bar-separator"></span>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Generate New Token (+)"
        onclick={() => (showCreateModal = true)}
      >
        <Icon name="plus" size={13} color="#57D38C" />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Revoke Selected Token (-)"
        disabled={!selectedToken || selectedToken.state !== "active"}
        onclick={handleRevokeSelected}
      >
        <Icon name="minus" size={13} />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Undo"
        disabled
      >
        <Icon name="undo" size={13} />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Redo"
        disabled
      >
        <Icon name="redo" size={13} />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Commit Changes"
        disabled
      >
        <Icon name="arrow-up" size={13} />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Rollback Changes"
        disabled
      >
        <Icon name="arrow-down" size={13} />
      </Button>
      <span class="bar-separator"></span>
      <div class="borderless-select-wrapper">
        <Select
          class="toolbar-select borderless-select"
          options={[
            { value: "auto", label: "Tx: Auto" },
            { value: "active", label: "Tx: Active Only" },
            { value: "revoked", label: "Tx: Revoked Only" },
          ]}
          bind:value={txMode}
        />
      </div>
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
        title="Search tokens"
        onclick={() => searchInputRef?.focus()}
      >
        <Icon name="search" size={13} />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Table layout"
      >
        <Icon name="table" size={13} />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Analytics & charts"
        onclick={() => (appState.activeTabId = "dashboard")}
      >
        <Icon name="chart" size={13} />
      </Button>
    </div>

    <div class="toolbar-right">
      <div class="borderless-select-wrapper format-select-wrapper">
        <Select
          class="toolbar-select borderless-select"
          options={[
            { value: "CSV", label: "CSV" },
            { value: "JSON", label: "JSON" },
          ]}
          bind:value={exportFormat}
        />
      </div>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Export tokens"
        onclick={handleExport}
      >
        <Icon name="download" size={13} />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Import"
        disabled
      >
        <Icon name="upload" size={13} />
      </Button>
      <span class="bar-separator"></span>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Previous"
        disabled
      >
        <Icon name="chevron-left" size={13} />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Next"
        disabled
      >
        <Icon name="chevron-right" size={13} />
      </Button>
      <span class="bar-separator"></span>
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
        placeholder="e.g. state = 'active' or client name"
        bind:this={searchInputRef}
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

  <!-- Tokens List Table -->
  <Box class="table-container">
    <DataGrid style="height: 100%;">
      {#snippet overlay()}
        {#if filteredTokens.length > 0}
          <div class="floating-row-badge" title="Token count">
            <Text size="sm">{filteredTokens.length} rows</Text>
          </div>
        {/if}
      {/snippet}
      <DataGridHead>
        <DataGridRow>
          <DataGridRowNumHead />
          <DataGridHeadCell width="150px">
            <DataGridHeaderInner>
              <Icon name="key" size={12} color="#EDA200" />
              <span>Token ID</span>
              <button
                type="button"
                class="ohjanus-data-grid-header-action"
                title="Sort by ID"
                onclick={() => toggleSort("id")}
              >
                <Icon name="sort" size={9} />
              </button>
            </DataGridHeaderInner>
          </DataGridHeadCell>
          <DataGridHeadCell width="220px">
            <DataGridHeaderInner>
              <Icon name="user" size={12} color="#7A7E85" />
              <span>Agent / Client Name</span>
              <button
                type="button"
                class="ohjanus-data-grid-header-action"
                title="Sort by Name"
                onclick={() => toggleSort("name")}
              >
                <Icon name="sort" size={9} />
              </button>
            </DataGridHeaderInner>
          </DataGridHeadCell>
          <DataGridHeadCell width="200px">
            <DataGridHeaderInner>
              <Icon name="shield" size={12} color="#3B82F6" />
              <span>Scopes</span>
            </DataGridHeaderInner>
          </DataGridHeadCell>
          <DataGridHeadCell width="120px">
            <DataGridHeaderInner>
              <Icon name="clock" size={12} color="#56A8F5" />
              <span>Created</span>
              <button
                type="button"
                class="ohjanus-data-grid-header-action"
                title="Sort by Created"
                onclick={() => toggleSort("created_at")}
              >
                <Icon name="sort" size={9} />
              </button>
            </DataGridHeaderInner>
          </DataGridHeadCell>
          <DataGridHeadCell width="170px">
            <DataGridHeaderInner>
              <Icon name="clock" size={12} color="#56A8F5" />
              <span>Expires</span>
              <button
                type="button"
                class="ohjanus-data-grid-header-action"
                title="Sort by Expires"
                onclick={() => toggleSort("expires_at")}
              >
                <Icon name="sort" size={9} />
              </button>
            </DataGridHeaderInner>
          </DataGridHeadCell>
          <DataGridHeadCell width="120px">
            <DataGridHeaderInner>
              <Icon name="clock" size={12} color="#56A8F5" />
              <span>Last Used</span>
            </DataGridHeaderInner>
          </DataGridHeadCell>
          <DataGridHeadCell width="100px">
            <DataGridHeaderInner>
              <Icon name="chart" size={12} color="#57D38C" />
              <span>State</span>
              <button
                type="button"
                class="ohjanus-data-grid-header-action"
                title="Sort by State"
                onclick={() => toggleSort("state")}
              >
                <Icon name="sort" size={9} />
              </button>
            </DataGridHeaderInner>
          </DataGridHeadCell>
          <DataGridHeadCell width="110px">
            <DataGridHeaderInner style="justify-content: flex-end;">
              <span>Actions</span>
            </DataGridHeaderInner>
          </DataGridHeadCell>
        </DataGridRow>
      </DataGridHead>
      <DataGridBody>
        {#if filteredTokens.length === 0}
          <DataGridRow>
            <DataGridRowNum index={0} />
            <DataGridCell tone="secondary">
              {#if appState.tokens.length === 0}
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
              <DataGridCell class="id-cell" truncate title={tok.id}>{tok.id}</DataGridCell>
              <DataGridCell class="name-cell" truncate title={tok.name}>
                <strong>{tok.name}</strong>
              </DataGridCell>
              <DataGridCell class="scopes-cell">
                {#each tok.scopes as sc}
                  <Badge variant="default" size="md">
                    {sc}
                  </Badge>
                {/each}
              </DataGridCell>
              <DataGridCell tone="secondary">{tok.created_at || "—"}</DataGridCell>
              <DataGridCell tone="secondary">{tok.expires_at}</DataGridCell>
              <DataGridCell tone="secondary">{tok.last_used_at || "—"}</DataGridCell>
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
                    >Revoked</Text
                  >
                {/if}
              </DataGridCell>
            </DataGridRow>
          {/each}
        {/if}
      </DataGridBody>
    </DataGrid>
  </Box>
</Box>

<!-- Create Token Modal via @ohjanus/ui -->
{#if showCreateModal}
  <Modal
    open={showCreateModal}
    onClose={() => (showCreateModal = false)}
    title="Generate New MCP Token"
    width="500px"
  >
    {#snippet children()}
      <Stack class="modal-form-stack" gap="14px">
        <Field label="Agent / Client Name:" required>
          <Input
            placeholder="e.g., Cursor IDE - Production Investigator"
            bind:value={tokenName}
          />
        </Field>

        <Field label="Allowed Janus Scopes (§1.2):">
          <Stack class="scopes-grid" gap="8px">
            <!-- svelte-ignore a11y_click_events_have_key_events -->
            <!-- svelte-ignore a11y_no_static_element_interactions -->
            <div class="chk-label" onclick={() => (scopes.read = !scopes.read)}>
              <Checkbox
                bind:checked={scopes.read}
                ariaLabel="read scope"
                onclick={(e: MouseEvent) => e.stopPropagation()}
              />
              <Text size="md"
                ><code>read</code> (db_list_connections, db_schema, db_read, db_explain)</Text
              >
            </div>
            <!-- svelte-ignore a11y_click_events_have_key_events -->
            <!-- svelte-ignore a11y_no_static_element_interactions -->
            <div
              class="chk-label"
              onclick={() => (scopes.write_preview = !scopes.write_preview)}
            >
              <Checkbox
                bind:checked={scopes.write_preview}
                ariaLabel="write_preview scope"
                onclick={(e: MouseEvent) => e.stopPropagation()}
              />
              <Text size="md"
                ><code>write_preview</code> (db_write_preview — dry-run simulation)</Text
              >
            </div>
            <!-- svelte-ignore a11y_click_events_have_key_events -->
            <!-- svelte-ignore a11y_no_static_element_interactions -->
            <div
              class="chk-label"
              onclick={() => (scopes.write_execute = !scopes.write_execute)}
            >
              <Checkbox
                bind:checked={scopes.write_execute}
                ariaLabel="write_execute scope"
                onclick={(e: MouseEvent) => e.stopPropagation()}
              />
              <Text size="md"
                ><code>write_execute</code> (db_write_execute — requires human approval)</Text
              >
            </div>
            <!-- svelte-ignore a11y_click_events_have_key_events -->
            <!-- svelte-ignore a11y_no_static_element_interactions -->
            <div class="chk-label" onclick={() => (scopes.admin = !scopes.admin)}>
              <Checkbox
                bind:checked={scopes.admin}
                ariaLabel="admin scope"
                onclick={(e: MouseEvent) => e.stopPropagation()}
              />
              <Text size="md"
                ><code>admin</code> (manage tokens &amp; connection configurations)</Text
              >
            </div>
          </Stack>
        </Field>

        <Field label="Token Validity (TTL Days):">
          <Select
            options={[
              { value: "7", label: "7 Days" },
              { value: "30", label: "30 Days (Recommended)" },
              { value: "90", label: "90 Days" },
              { value: "365", label: "1 Year" },
            ]}
            value={String(ttlDays)}
            onchange={(v) => (ttlDays = Number(v))}
          />
        </Field>
      </Stack>
    {/snippet}

    {#snippet footer()}
      <Button variant="secondary" onclick={() => (showCreateModal = false)}
        >Cancel</Button
      >
      <Button
        variant="primary"
        onclick={handleCreate}
        disabled={!tokenName.trim() || creating}
      >
        {creating ? "Generating..." : "Generate Token"}
      </Button>
    {/snippet}
  </Modal>
{/if}

<!-- One-Time Secret Reveal Modal via @ohjanus/ui -->
{#if showSecretModal && generatedToken}
  <Modal
    open={showSecretModal}
    onClose={() => (showSecretModal = false)}
    title="MCP Token Generated Successfully"
    width="520px"
  >
    {#snippet children()}
      <Stack class="secret-reveal-stack" gap="12px">
        <Alert variant="warning" title="WARNING:">
          <p style="margin-top: 2px;">
            This token secret is displayed only <strong>ONCE</strong> and cannot
            be retrieved later. Copy and paste it into your client configuration
            now.
          </p>
        </Alert>

        <Text weight="semibold" size="sm">Opaque Bearer Token:</Text>
        <Flex class="secret-box font-mono" align="center" justify="between">
          <span>{generatedToken?.rawToken}</span>
          <Button variant="secondary" size="xs" onclick={copySecret}>
            {#if copied}
              <Icon name="check" size={12} />
              <span>Copied</span>
            {:else}
              <Icon name="copy" size={12} />
              <span>Copy</span>
            {/if}
          </Button>
        </Flex>

        <Text size="xs" color="muted">
          Example Cursor config (<code>~/.cursor/mcp.json</code>):
          <pre class="mcp-config-snippet font-mono">{`{
  "mcpServers": {
    "ohjanus": {
      "command": "janus",
      "args": ["serve", "--token", "${generatedToken?.rawToken}"]
    }
  }
}`}</pre>
        </Text>
      </Stack>
    {/snippet}

    {#snippet footer()}
      <Button variant="primary" onclick={() => (showSecretModal = false)}>
        Done
      </Button>
    {/snippet}
  </Modal>
{/if}

<style>
  :global(.tokens-view) {
    display: flex;
    flex-direction: column;
    height: 100%;
    background-color: var(--bg-canvas);
    color: var(--text-primary);
    overflow: hidden;
  }

  :global(.tokens-view .table-toolbar) {
    height: var(--toolbar-height, 32px);
    background-color: var(--bg-toolbar);
    border-bottom: 1px solid var(--border-subtle);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 8px;
    flex-shrink: 0;
  }

  :global(.tokens-view .toolbar-left),
  :global(.tokens-view .toolbar-right) {
    display: flex;
    align-items: center;
    gap: 2px;
  }

  :global(.tokens-view .bar-separator) {
    width: 1px;
    height: 14px;
    background-color: var(--border-default);
    margin: 0 4px;
  }

  :global(.tokens-view .borderless-select-wrapper) {
    display: inline-flex;
    align-items: center;
  }

  :global(.tokens-view .borderless-select) {
    min-width: unset !important;
    width: auto !important;
  }

  :global(.tokens-view .borderless-select .ohjanus-select-trigger) {
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

  :global(.tokens-view .borderless-select .ohjanus-select-trigger:hover) {
    background-color: var(--bg-hover, #313438) !important;
    color: var(--text-primary, #DFE1E5) !important;
  }

  :global(.tokens-view .ddl-btn) {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 0 4px !important;
    width: auto !important;
  }

  :global(.tokens-view .filter-bar) {
    height: var(--filterbar-height, 30px);
    background-color: var(--bg-canvas);
    border-bottom: 1px solid var(--border-subtle);
    display: flex;
    align-items: center;
    padding: 0 8px;
    gap: 12px;
    flex-shrink: 0;
  }

  :global(.tokens-view .filter-group) {
    display: flex;
    align-items: center;
    gap: 6px;
  }

  :global(.tokens-view .where-group) {
    flex: 1.1;
  }

  :global(.tokens-view .orderby-group) {
    flex: 0.9;
  }

  :global(.tokens-view .filter-icon) {
    color: var(--text-muted);
    font-size: 11px;
    user-select: none;
    display: inline-flex;
    align-items: center;
  }

  :global(.tokens-view .filter-input) {
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

  :global(.tokens-view .filter-input:focus),
  :global(.tokens-view .filter-input:hover) {
    border: none !important;
    outline: none !important;
    box-shadow: none !important;
  }

  :global(.tokens-view .table-container) {
    flex: 1;
    overflow: auto;
  }

  :global(.tokens-view .id-cell) {
    color: var(--syntax-number, #6897bb);
  }

  :global(.tokens-view .scopes-cell) {
    display: flex;
    align-items: center;
    gap: 4px;
    flex-wrap: wrap;
  }

  :global(.scopes-grid) {
    background-color: var(--bg-canvas, #1e1f22);
    border: 1px solid var(--border-default, #393b40);
    border-radius: 4px;
    padding: 10px;
  }

  :global(.chk-label) {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: var(--font-size-sm, 13px);
    cursor: pointer;
  }

  :global(.secret-box) {
    background-color: var(--bg-canvas);
    border: 1px dashed var(--action-warning);
    border-radius: 4px;
    padding: 10px 12px;
    font-size: var(--font-size-sm, 13px);
    color: #facc15;
    word-break: break-all;
  }

  :global(.mcp-config-snippet) {
    background-color: var(--bg-canvas);
    border: 1px solid var(--border-default);
    border-radius: 4px;
    padding: 8px;
    margin-top: 6px;
    font-size: var(--font-size-xs, 12px);
    color: var(--syntax-string, #6aab73);
  }
</style>
