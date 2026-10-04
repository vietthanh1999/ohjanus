<script lang="ts">
  import { connectionsState } from '@/features/connections';
  import { Modal, Button, Input, Select, Checkbox, Field, Alert, toast, Stack, Text } from '@ohjanus/ui';

  let name = $state('');
  let host = $state('');
  let port = $state('5432');
  let database = $state('');
  let username = $state('');
  let password = $state('');
  let sslmode = $state('prefer');
  let readOnly = $state(true);
  let dsn = $state('');
  let creating = $state(false);
  let error = $state<string | null>(null);

  const canSubmit = $derived(
    name.trim() !== '' && (dsn.trim() !== '' || (host.trim() !== '' && database.trim() !== ''))
  );

  async function handleConnect() {
    if (!canSubmit || creating) return;
    creating = true;
    error = null;
    try {
      const created = await connectionsState.createConnection({
        name: name.trim(),
        driver: 'postgres',
        dsn: dsn.trim() || undefined,
        host: host.trim() || undefined,
        port: Number(port) || undefined,
        database: database.trim() || undefined,
        username: username.trim() || undefined,
        password: password || undefined,
        sslmode,
        read_only: readOnly
      });
      toast.success('Connection Added', `"${created.name}" connected and ready to explore.`);
      close();
    } catch (e) {
      error = e instanceof Error ? e.message : String(e);
    } finally {
      creating = false;
    }
  }

  function close() {
    connectionsState.createConnectionModalOpen = false;
    name = '';
    host = '';
    port = '5432';
    database = '';
    username = '';
    password = '';
    sslmode = 'prefer';
    readOnly = true;
    dsn = '';
    error = null;
  }
</script>

{#if connectionsState.createConnectionModalOpen}
  <Modal open={connectionsState.createConnectionModalOpen} onClose={close} title="New PostgreSQL Connection" width="480px">
    {#snippet children()}
      <Stack gap="12px">
        {#if error}
          <Alert variant="danger" title="Could not connect">
            <Text size="sm">{error}</Text>
          </Alert>
        {/if}
        <Field label="Connection Name" required>
          <Input placeholder="e.g. prod-warehouse" bind:value={name} />
        </Field>
        <Field label="Host" required>
          <Input placeholder="e.g. 127.0.0.1 or db.internal" bind:value={host} />
        </Field>
        <Stack gap="12px" style="display: grid; grid-template-columns: 1fr 1fr; gap: 12px;">
          <Field label="Port">
            <Input type="number" bind:value={port} />
          </Field>
          <Field label="SSL Mode">
            <Select
              bind:value={sslmode}
              options={[
                { value: 'prefer', label: 'prefer' },
                { value: 'require', label: 'require' },
                { value: 'disable', label: 'disable' }
              ]}
            />
          </Field>
        </Stack>
        <Field label="Database" required>
          <Input placeholder="e.g. app" bind:value={database} />
        </Field>
        <Stack gap="12px" style="display: grid; grid-template-columns: 1fr 1fr; gap: 12px;">
          <Field label="Username">
            <Input placeholder="e.g. app" bind:value={username} autocomplete="username" />
          </Field>
          <Field label="Password">
            <Input type="password" bind:value={password} autocomplete="current-password" />
          </Field>
        </Stack>
        <Field label="Or paste a full DSN (takes precedence)">
          <Input placeholder="postgres://user:pass@host:5432/db?sslmode=require" bind:value={dsn} />
        </Field>
        <label style="display: inline-flex; align-items: center; gap: 8px; cursor: pointer;">
          <Checkbox bind:checked={readOnly} ariaLabel="Read-only connection" />
          <Text size="sm">Read-only (recommended for governed access)</Text>
        </label>
        <Text size="xs" color="muted">The server tests the connection before saving — wrong credentials fail here, not later.</Text>
      </Stack>
    {/snippet}

    {#snippet footer()}
      <Button variant="secondary" onclick={close}>Cancel</Button>
      <Button variant="primary" onclick={handleConnect} disabled={!canSubmit || creating}>
        {creating ? 'Connecting…' : 'Connect'}
      </Button>
    {/snippet}
  </Modal>
{/if}
