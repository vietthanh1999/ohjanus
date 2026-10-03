<script lang="ts">
  import { appState } from '../../state/appState.svelte';

  let apiBase = $state('http://localhost:8788/api/v1');
  let queryTimeout = $state(30);
  let autoCommit = $state(false);

  function close() {
    appState.settingsModalOpen = false;
  }
</script>

{#if appState.settingsModalOpen}
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div class="jb-modal-backdrop" onclick={close}>
    <div class="jb-modal" style="width: 500px;" onclick={(e) => e.stopPropagation()}>
      <div class="jb-modal-header">
        <span>Settings & Preferences</span>
        <button class="jb-icon-btn" onclick={close}>✕</button>
      </div>

      <div class="jb-modal-body">
        <div class="form-group">
          <label for="admin-api-base" class="form-label">Gateway Admin API Base URL</label>
          <input id="admin-api-base" type="text" class="form-input code-text" bind:value={apiBase} />
          <span class="form-hint">Port 8788 is the dedicated Janus Admin REST API port separated from MCP clients.</span>
        </div>

        <div class="form-group">
          <label for="query-timeout-sec" class="form-label">Max Query Execution Timeout (seconds)</label>
          <input id="query-timeout-sec" type="number" class="form-input code-text" bind:value={queryTimeout} />
        </div>

        <div class="form-group checkbox-group">
          <label class="chk-label">
            <input type="checkbox" bind:checked={autoCommit} />
            <span>Enable Auto-Commit on write queries (Not recommended on production)</span>
          </label>
        </div>

        <div class="form-group">
          <div class="form-label">Active Theme & Density</div>
          <div class="theme-info">
            <span>Theme: <strong>OhJanus JetBrains Dark (Official)</strong></span>
            <span>Density: <strong>High Information Density (24-26px row height)</strong></span>
            <span>Code Font: <strong>JetBrains Mono</strong></span>
            <span>UI Font: <strong>Inter</strong></span>
          </div>
        </div>
      </div>

      <div class="jb-modal-footer">
        <button class="jb-btn-primary" onclick={close}>Save & Close</button>
      </div>
    </div>
  </div>
{/if}

<style>
  .form-group {
    display: flex;
    flex-direction: column;
    gap: 6px;
    margin-bottom: 14px;
  }

  .form-label {
    font-size: 11px;
    font-weight: 500;
    color: var(--text-secondary);
  }

  .form-input {
    background-color: #1E1F22;
    border: 1px solid var(--border-default);
    border-radius: 3px;
    padding: 6px 8px;
    color: var(--text-primary);
    font-size: 12px;
  }

  .form-hint {
    font-size: 10px;
    color: var(--text-muted);
  }

  .chk-label {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 11px;
    color: var(--text-primary);
    cursor: pointer;
  }

  .theme-info {
    background-color: #1E1F22;
    border: 1px solid var(--border-default);
    padding: 10px;
    border-radius: 4px;
    display: flex;
    flex-direction: column;
    gap: 4px;
    font-size: 11px;
    color: var(--text-secondary);
  }
</style>
