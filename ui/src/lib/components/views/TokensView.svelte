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
  <!-- Header -->
  <Flex as="header" class="view-header" align="center" justify="between">
    <Flex align="center" gap="8px">
      <span class="header-icon"
        ><Icon name="key" size={16} color="#EDA200" /></span
      >
      <Text class="view-title" size="xl" weight="semibold">MCP Client Tokens Management</Text>
      <Text class="view-desc" color="muted"
        >Issue and manage secure opaque tokens for AI agents (Cursor, Claude
        Desktop, custom bots)</Text
      >
    </Flex>
    <Button variant="primary" onclick={() => (showCreateModal = true)}>
      <Icon name="plus" size={14} />
      <span>Generate New Token</span>
    </Button>
  </Flex>

  <!-- Tokens List Table -->
  <Box class="table-container">
    <DataGrid style="height: 100%;">
      <DataGridHead>
        <DataGridRow>
          <DataGridRowNumHead />
          <DataGridHeadCell width="150px">
            <DataGridHeaderInner>
              <Icon name="key" size={12} color="#EDA200" />
              <span>Token ID</span>
            </DataGridHeaderInner>
          </DataGridHeadCell>
          <DataGridHeadCell width="220px">
            <DataGridHeaderInner>
              <Icon name="user" size={12} color="#7A7E85" />
              <span>Agent / Client Name</span>
            </DataGridHeaderInner>
          </DataGridHeadCell>
          <DataGridHeadCell width="200px">
            <DataGridHeaderInner>
              <Icon name="shield" size={12} color="#3B82F6" />
              <span>Scopes</span>
            </DataGridHeaderInner>
          </DataGridHeadCell>
          <DataGridHeadCell width="110px">
            <DataGridHeaderInner>
              <Icon name="clock" size={12} color="#56A8F5" />
              <span>Created</span>
            </DataGridHeaderInner>
          </DataGridHeadCell>
          <DataGridHeadCell width="170px">
            <DataGridHeaderInner>
              <Icon name="clock" size={12} color="#56A8F5" />
              <span>Expires</span>
            </DataGridHeaderInner>
          </DataGridHeadCell>
          <DataGridHeadCell width="110px">
            <DataGridHeaderInner>
              <Icon name="clock" size={12} color="#56A8F5" />
              <span>Last Used</span>
            </DataGridHeaderInner>
          </DataGridHeadCell>
          <DataGridHeadCell width="100px">
            <DataGridHeaderInner>
              <Icon name="chart" size={12} color="#57D38C" />
              <span>State</span>
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
        {#each appState.tokens as tok, idx (tok.id)}
          <DataGridRow>
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

  :global(.tokens-view .view-header) {
    height: 52px;
    padding: 0 16px;
    border-bottom: 1px solid var(--border-default);
    background-color: var(--bg-toolbar);
  }

  :global(.tokens-view .header-icon) {
    font-size: 18px;
    margin-right: 6px;
  }

  :global(.tokens-view .view-title) {
    font-weight: 600;
    font-size: var(--font-size-lg, 16px);
    color: var(--text-primary);
  }

  :global(.tokens-view .view-desc) {
    font-size: var(--font-size-xs, 12px);
    color: var(--text-muted);
    margin-left: 12px;
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
