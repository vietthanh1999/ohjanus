<script lang="ts">
  import { appState, type McpToken } from "../../state/appState.svelte";
  import {
    Button,
    Badge,
    Modal,
    Input,
    Field,
    Alert,
    toast,
    Box,
    Flex,
    Stack,
    Text,
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
      <span class="view-title">MCP Client Tokens Management</span>
      <span class="view-desc"
        >Issue and manage secure opaque tokens for AI agents (Cursor, Claude
        Desktop, custom bots)</span
      >
    </Flex>
    <Button variant="primary" onclick={() => (showCreateModal = true)}>
      <Icon name="plus" size={14} />
      <span>Generate New Token</span>
    </Button>
  </Flex>

  <!-- Tokens List Table -->
  <Box class="table-container font-mono">
    <table class="data-grid">
      <thead>
        <tr>
          <th>TOKEN ID</th>
          <th>AGENT / CLIENT NAME</th>
          <th>SCOPES</th>
          <th>CREATED</th>
          <th>EXPIRES</th>
          <th>LAST USED</th>
          <th>STATE</th>
          <th style="text-align: right;">ACTIONS</th>
        </tr>
      </thead>
      <tbody>
        {#each appState.tokens as tok}
          <tr>
            <td class="id-cell">{tok.id}</td>
            <td class="name-cell font-sans">
              <strong>{tok.name}</strong>
            </td>
            <td class="scopes-cell">
              {#each tok.scopes as sc}
                <Badge variant="default" size="sm" class="scope-pill">
                  {sc}
                </Badge>
              {/each}
            </td>
            <td>{tok.created_at}</td>
            <td>{tok.expires_at}</td>
            <td>{tok.last_used_at}</td>
            <td>
              {#if tok.state === "active"}
                <Badge variant="success" size="sm">ACTIVE</Badge>
              {:else}
                <Badge variant="danger" size="sm">REVOKED</Badge>
              {/if}
            </td>
            <td style="text-align: right;">
              {#if tok.state === "active"}
                <Button
                  variant="danger"
                  size="xs"
                  onclick={() => handleRevoke(tok.id)}
                >
                  Revoke
                </Button>
              {:else}
                <span
                  style="color: var(--text-muted); font-size: var(--font-size-xs, 12px);"
                  >Revoked</span
                >
              {/if}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
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
            <label class="chk-label">
              <input type="checkbox" bind:checked={scopes.read} />
              <span
                ><code>read</code> (db_list_connections, db_schema, db_read, db_explain)</span
              >
            </label>
            <label class="chk-label">
              <input type="checkbox" bind:checked={scopes.write_preview} />
              <span
                ><code>write_preview</code> (db_write_preview — dry-run simulation)</span
              >
            </label>
            <label class="chk-label">
              <input type="checkbox" bind:checked={scopes.write_execute} />
              <span
                ><code>write_execute</code> (db_write_execute — requires human approval)</span
              >
            </label>
            <label class="chk-label">
              <input type="checkbox" bind:checked={scopes.admin} />
              <span
                ><code>admin</code> (manage tokens &amp; connection configurations)</span
              >
            </label>
          </Stack>
        </Field>

        <Field label="Token Validity (TTL Days):">
          <select id="tok-ttl" bind:value={ttlDays} class="ohjanus-select">
            <option value={7}>7 Days</option>
            <option value={30}>30 Days (Recommended)</option>
            <option value={90}>90 Days</option>
            <option value={365}>1 Year</option>
          </select>
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

  :global(.tokens-view .data-grid) {
    width: 100%;
    border-collapse: collapse;
    font-size: var(--font-size-base, 14px);
  }

  :global(.tokens-view .data-grid th) {
    background-color: var(--bg-table-header);
    border-bottom: 1px solid var(--border-default);
    padding: 8px 12px;
    text-align: left;
    color: var(--text-muted);
    font-weight: 600;
    font-size: var(--font-size-xs, 12px);
  }

  :global(.tokens-view .data-grid td) {
    border-bottom: 1px solid var(--border-subtle);
    padding: 9px 12px;
    color: var(--text-primary);
  }

  :global(.tokens-view .data-grid tr:hover) {
    background-color: var(--bg-hover);
  }

  :global(.tokens-view .id-cell) {
    color: var(--syntax-number, #6897bb);
  }

  :global(.tokens-view .scopes-cell) {
    display: flex;
    gap: 4px;
    flex-wrap: wrap;
  }

  :global(.scope-pill) {
    font-size: var(--font-size-2xs, 11px) !important;
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

  :global(.ohjanus-select) {
    width: 100%;
    height: var(--control-height-md, 34px);
    background-color: var(--bg-canvas, #1e1f22);
    border: 1px solid var(--border-default, #393b40);
    border-radius: 4px;
    color: var(--text-primary, #dfe1e5);
    padding: 0 8px;
    font-size: var(--font-size-base, 14px);
    outline: none;
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
