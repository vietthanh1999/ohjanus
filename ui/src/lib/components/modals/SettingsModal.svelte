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

    <Box class="density-section">
      <Text size="sm" weight="semibold" style="margin-bottom: 8px; display: block; color: var(--text-primary);">
        UI Density & Typography Scale
      </Text>
      <div class="density-buttons">
        <button
          type="button"
          class="density-btn"
          class:active={appState.uiDensity === 'comfortable'}
          onclick={() => appState.setDensity('comfortable')}
        >
          <span class="d-title">Comfortable (14px)</span>
          <span class="d-desc">Base 14px — Dễ đọc, thoáng đãng, tối ưu màn hình lớn</span>
        </button>

        <button
          type="button"
          class="density-btn"
          class:active={appState.uiDensity === 'standard'}
          onclick={() => appState.setDensity('standard')}
        >
          <span class="d-title">Standard (13px)</span>
          <span class="d-desc">Base 13px — Cân bằng giữa mật độ thông tin và độ rõ nét</span>
        </button>

        <button
          type="button"
          class="density-btn"
          class:active={appState.uiDensity === 'compact'}
          onclick={() => appState.setDensity('compact')}
        >
          <span class="d-title">Compact (12px)</span>
          <span class="d-desc">Base 12px — Hiển thị tối đa dòng dữ liệu (DataGrip legacy)</span>
        </button>
      </div>
    </Box>

    <Box class="theme-info-box">
      <Text size="xs" weight="semibold" class="theme-title">Active Theme & Details</Text>
      <Stack class="theme-details" gap="4px">
        <span>Theme: <strong>OhJanus JetBrains Dark (Official)</strong></span>
        <span>Selected Density: <strong style="text-transform: capitalize;">{appState.uiDensity}</strong></span>
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

  .density-buttons {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .density-btn {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    padding: 10px 12px;
    border-radius: 6px;
    border: 1px solid var(--border-default, #393B40);
    background-color: var(--bg-canvas, #1E1F22);
    cursor: pointer;
    text-align: left;
    transition: all 0.15s ease;
  }

  .density-btn:hover {
    border-color: var(--border-strong, #4E5157);
    background-color: var(--bg-hover, #2B2D30);
  }

  .density-btn.active {
    border-color: var(--border-accent, #3574F0);
    background-color: rgba(53, 116, 240, 0.12);
  }

  .density-btn .d-title {
    font-size: var(--font-size-sm, 13px);
    font-weight: 600;
    color: var(--text-primary, #DFE1E5);
  }

  .density-btn.active .d-title {
    color: #56A8F5;
  }

  .density-btn .d-desc {
    font-size: var(--font-size-xs, 12px);
    color: var(--text-muted, #7A7E85);
    margin-top: 2px;
  }
</style>
