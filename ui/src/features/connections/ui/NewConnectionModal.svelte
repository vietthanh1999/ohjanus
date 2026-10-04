<script lang="ts">
  import { connectionsState, probeConnection } from "@/features/connections";
  import {
    Modal,
    Button,
    Input,
    Select,
    Checkbox,
    Alert,
    toast,
    Text,
    Box,
    Flex,
    Stack,
    Grid,
    Label,
  } from "@ohjanus/ui";

  let name = $state("");
  let host = $state("");
  let port = $state("5432");
  let database = $state("");
  let username = $state("");
  let password = $state("");
  let sslmode = $state("prefer");
  let readOnly = $state(true);
  let dsn = $state("");
  let creating = $state(false);
  let testing = $state(false);
  let error = $state<string | null>(null);
  let probe = $state<{
    ok: boolean;
    latency?: number;
    message?: string;
  } | null>(null);

  const canSubmit = $derived(
    name.trim() !== "" && (dsn.trim() !== "" || host.trim() !== ""),
  );
  const canTest = $derived(
    dsn.trim() !== "" || host.trim() !== "",
  );

  function payload() {
    return {
      name: name.trim(),
      driver: "postgres",
      dsn: dsn.trim() || undefined,
      host: host.trim() || undefined,
      port: Number(port) || undefined,
      database: database.trim() || undefined,
      username: username.trim() || undefined,
      password: password || undefined,
      sslmode,
      read_only: readOnly,
    };
  }

  // Live URL preview, DataGrip-style (password masked).
  const previewUrl = $derived.by(() => {
    if (dsn.trim() !== "") return dsn.trim();
    if (host.trim() === "") return "";
    const auth =
      username.trim() === ""
        ? ""
        : `${username.trim()}${password !== "" ? ":***" : ""}@`;
    const portPart =
      port.trim() === "" || port.trim() === "5432" ? "" : `:${port.trim()}`;
    const dbPart = database.trim() === "" ? "postgres" : database.trim();
    return `postgres://${auth}${host.trim()}${portPart}/${dbPart}?sslmode=${sslmode}`;
  });

  async function handleTest() {
    if (!canTest || testing) return;
    testing = true;
    probe = null;
    try {
      const res = await probeConnection(payload());
      probe = { ok: true, latency: res.latency_ms };
    } catch (e) {
      probe = {
        ok: false,
        message: e instanceof Error ? e.message : String(e),
      };
    } finally {
      testing = false;
    }
  }

  async function handleConnect() {
    if (!canSubmit || creating) return;
    creating = true;
    error = null;
    try {
      const created = await connectionsState.createConnection(payload());
      toast.success(
        "Connection Added",
        `"${created.name}" connected and ready to explore.`,
      );
      close();
    } catch (e) {
      error = e instanceof Error ? e.message : String(e);
    } finally {
      creating = false;
    }
  }

  function close() {
    connectionsState.createConnectionModalOpen = false;
    name = "";
    host = "";
    port = "5432";
    database = "";
    username = "";
    password = "";
    sslmode = "prefer";
    readOnly = true;
    dsn = "";
    error = null;
    probe = null;
  }
</script>

{#if connectionsState.createConnectionModalOpen}
  <Modal
    open={connectionsState.createConnectionModalOpen}
    onClose={close}
    title="New PostgreSQL Connection"
    width="560px"
  >
    {#snippet children()}
      <Stack class="conn-form" gap="10px">
        {#if error}
          <Alert variant="danger" title="Could not connect">
            <Text size="sm">{error}</Text>
          </Alert>
        {/if}

        <Grid class="conn-row" columns="88px 1fr" gap="10px" align="center">
          <Text class="conn-label" size="sm">Name:</Text>
          <Input placeholder="e.g. prod-warehouse" bind:value={name} />
        </Grid>

        <Grid class="conn-row" columns="88px 1fr" gap="10px" align="center">
          <Text class="conn-label" size="sm">Host:</Text>
          <Flex class="conn-split" align="center" gap="8px">
            <Input
              placeholder="e.g. 127.0.0.1 or db.internal"
              bind:value={host}
            />
            <Text class="conn-sublabel" size="sm">Port:</Text>
            <Box as="span" class="conn-port">
              <Input type="number" bind:value={port} />
            </Box>
          </Flex>
        </Grid>

        <Grid class="conn-row" columns="88px 1fr" gap="10px" align="center">
          <Text class="conn-label" size="sm">Database:</Text>
          <Input placeholder="e.g. app (empty = postgres)" bind:value={database} />
        </Grid>

        <Grid class="conn-row" columns="88px 1fr" gap="10px" align="center">
          <Text class="conn-label" size="sm">User:</Text>
          <Input
            placeholder="e.g. app"
            bind:value={username}
            autocomplete="username"
          />
        </Grid>

        <Grid class="conn-row" columns="88px 1fr" gap="10px" align="center">
          <Text class="conn-label" size="sm">Password:</Text>
          <Input
            type="password"
            bind:value={password}
            autocomplete="current-password"
          />
        </Grid>

        <Grid class="conn-row" columns="88px 1fr" gap="10px" align="center">
          <Text class="conn-label" size="sm">SSL:</Text>
          <Flex class="conn-split" align="center" gap="8px">
            <Select
              bind:value={sslmode}
              options={[
                { value: "prefer", label: "prefer" },
                { value: "require", label: "require" },
                { value: "disable", label: "disable" },
              ]}
            />
            <Label class="conn-check">
              <Checkbox
                bind:checked={readOnly}
                ariaLabel="Read-only connection"
              />
              <Text size="sm">Read-only</Text>
            </Label>
          </Flex>
        </Grid>

        <Grid class="conn-row" columns="88px 1fr" gap="10px" align="center">
          <Text class="conn-label" size="sm">DSN:</Text>
          <Input
            placeholder="postgres://user:pass@host:5432/db (overrides fields above)"
            bind:value={dsn}
            class="conn-mono"
          />
        </Grid>

        {#if previewUrl !== ""}
          <Grid class="conn-row" columns="88px 1fr" gap="10px" align="center">
            <Text class="conn-label" size="sm">URL:</Text>
            <Box class="conn-url">{previewUrl}</Box>
          </Grid>
        {/if}

        {#if probe}
          <Grid class="conn-row" columns="88px 1fr" gap="10px" align="center">
            <Text class="conn-label" size="sm"></Text>
            {#if probe.ok}
              <Text size="sm" color="success"
                >Connected{probe.latency !== undefined
                  ? ` in ${probe.latency} ms`
                  : ""}.</Text
              >
            {:else}
              <Text size="sm" color="danger"
                >{probe.message ?? "Connection failed."}</Text
              >
            {/if}
          </Grid>
        {/if}
      </Stack>
    {/snippet}

    {#snippet footer()}
      <Button
        size="sm"
        variant="secondary"
        onclick={handleTest}
        disabled={!canTest || testing}
      >
        {testing ? "Testing…" : "Test Connection"}
      </Button>
      <Button size="sm" variant="secondary" onclick={close}>Cancel</Button>
      <Button
        size="sm"
        variant="primary"
        onclick={handleConnect}
        disabled={!canSubmit || creating}
      >
        {creating ? "Connecting…" : "Connect"}
      </Button>
    {/snippet}
  </Modal>
{/if}

<!-- <style>
  :global(.conn-label.conn-label) {
    font-size: 12.5px;
    color: var(--text-secondary, #9da0a8);
    text-align: right;
    white-space: nowrap;
  }

  :global(.conn-split) {
    min-width: 0;
  }

  :global(.conn-split > .ohjanus-input-wrapper) {
    flex: 1;
    min-width: 0;
  }

  :global(.conn-sublabel.conn-sublabel) {
    font-size: 12.5px;
    color: var(--text-secondary, #9da0a8);
    flex-shrink: 0;
  }

  :global(.conn-port.conn-port) {
    width: 84px;
    flex: none;
  }

  :global(.conn-port .ohjanus-input-wrapper) {
    width: 100%;
  }

  :global(.conn-check.conn-check) {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
    flex-shrink: 0;
    margin-left: auto;
  }

  :global(.conn-mono.conn-mono) {
    font-family: var(--font-code, "JetBrains Mono", monospace);
    font-size: 11.5px;
  }

  :global(.conn-url.conn-url) {
    font-family: var(--font-code, "JetBrains Mono", monospace);
    font-size: 11.5px;
    color: var(--text-secondary, #9da0a8);
    background-color: var(--bg-canvas, #1e1f22);
    border: 1px solid var(--border-subtle, #2b2d30);
    border-radius: 4px;
    padding: 7px 10px;
    overflow-x: auto;
    white-space: nowrap;
  }
</style> -->
