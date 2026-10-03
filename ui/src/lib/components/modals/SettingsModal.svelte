<script lang="ts">
  import { appState } from '../../state/appState.svelte';
  import { Modal, Field, Input, Switch, Button, toast, Stack, Box, Text } from '@ohjanus/ui';

  let apiBase = $state('http://localhost:8788/api/v1');
  let queryTimeout = $state('30');
  let autoCommit = $state(false);

  function close() {
    appState.settingsModalOpen = false;
  }

  function handleSave() {
    close();
    toast.success('Settings updated successfully');
  }
</script>

<Modal
  open={appState.settingsModalOpen}
  title="Settings & Preferences"
  width="520px"
  onClose={close}
>
  <Stack class="settings-content" gap="16px">
    <Field
      label="Gateway Admin API Base URL"
      hint="Port 8788 is the dedicated Janus Admin REST API port separated from MCP clients."
    >
      <Input
        bind:value={apiBase}
        placeholder="http://localhost:8788/api/v1"
        class="code-text"
      />
    </Field>

    <Field
      label="Max Query Execution Timeout (seconds)"
      hint="Default timeout before gateway terminates long-running queries."
    >
      <Input
        type="number"
        bind:value={queryTimeout}
        placeholder="30"
        class="code-text"
      />
    </Field>

    <Box class="toggle-section">
      <Switch
        checked={autoCommit}
        label="Enable Auto-Commit on write queries (Not recommended on production)"
        onchange={(val) => { autoCommit = val; }}
      />
    </Box>

    <Box class="theme-info-box">
      <Text size="xs" weight="semibold" class="theme-title">Active Theme & Density</Text>
      <Stack class="theme-details" gap="4px">
        <span>Theme: <strong>OhJanus JetBrains Dark (Official)</strong></span>
        <span>Density: <strong>High Information Density (24-26px row height)</strong></span>
        <span>Code Font: <strong>JetBrains Mono</strong></span>
        <span>UI Font: <strong>Inter</strong></span>
      </Stack>
    </Box>
  </Stack>

  {#snippet footer()}
    <Button variant="secondary" onclick={close}>Cancel</Button>
    <Button variant="primary" onclick={handleSave}>Save & Close</Button>
  {/snippet}
</Modal>

<style>
  :global(.toggle-section) {
    padding: 6px 0;
  }

  :global(.theme-info-box) {
    background-color: var(--bg-canvas, #1E1F22);
    border: 1px solid var(--border-default, #393B40);
    padding: 12px;
    border-radius: 4px;
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  :global(.theme-title) {
    color: var(--text-secondary, #9DA0A8);
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  :global(.theme-details) {
    font-size: 11px;
    color: var(--text-secondary, #9DA0A8);
  }

  :global(.theme-details strong) {
    color: var(--text-primary, #DFE1E5);
  }
</style>
