<script lang="ts">
  import { appState } from '../../state/appState.svelte';
  import { Icon } from '@ohjanus/icons';
</script>

<div class="log-console-container">
  <div class="log-content code-text">
    {#if appState.consoleLogs.length === 0}
      <div class="log-line">
        <span class="log-ts">--</span>
        <span class="log-text">No output yet. Run a query or open a table to see live execution logs.</span>
      </div>
    {:else}
      {#each appState.consoleLogs as log (log.id)}
        <div class="log-line" class:error-line={log.type === 'error'}>
          <span class="log-ts">[{log.timestamp}]</span>
          {#if log.connection}
            <span class="log-schema">{log.connection}&gt;</span>
          {/if}
          {#if log.querySnippet}
            <span class="ident-sql">{log.querySnippet}</span>
          {:else}
            <span class="log-text">{log.summary}</span>
          {/if}
        </div>
        {#if log.querySnippet}
          <div class="log-line indent-sql">
            <span class="log-text">{log.summary}</span>
          </div>
        {/if}
      {/each}
    {/if}
  </div>

  <div class="log-action-strip">
    <button type="button" class="strip-btn" title="Reload table data" onclick={() => void appState.loadTableData()}>
      <Icon name="refresh" size={13} />
    </button>
    <button type="button" class="strip-btn" title="Clear Console" onclick={() => appState.clearLogs()}>
      <Icon name="trash" size={13} />
    </button>
  </div>
</div>

<style>
  .log-console-container {
    width: 100%;
    height: 100%;
    display: flex;
    background-color: var(--bg-canvas, #1E1F22);
    overflow: hidden;
  }

  .log-content {
    flex: 1;
    overflow-y: auto;
    padding: 10px 14px;
    font-size: var(--font-size-base, 14px);
    line-height: 20px;
    font-family: var(--font-code, 'JetBrains Mono', monospace);
  }

  .log-line {
    display: flex;
    gap: 8px;
  }

  .log-line.indent-sql {
    padding-left: 24px;
  }

  .log-line.error-line .log-text {
    color: var(--action-danger, #E55353);
  }

  .log-ts {
    color: var(--text-muted, #7A7E85);
  }

  .log-text {
    color: var(--text-primary, #DFE1E5);
  }

  .log-schema {
    color: var(--text-primary, #DFE1E5);
    font-weight: 500;
  }

  .ident-sql {
    color: var(--text-primary, #DFE1E5);
  }

  .log-action-strip {
    width: 32px;
    background-color: var(--bg-canvas, #1E1F22);
    border-left: 1px solid var(--border-subtle, #323438);
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 4px 0;
    gap: 4px;
    flex-shrink: 0;
  }

  .strip-btn {
    width: var(--icon-btn-size-sm, 26px);
    height: var(--icon-btn-size-sm, 26px);
    border-radius: var(--radius-sm, 4px);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: var(--font-size-xs, 11px);
    color: var(--text-muted, #7A7E85);
  }

  :global(.strip-btn svg) {
    width: var(--icon-size-sm, 14px);
    height: var(--icon-size-sm, 14px);
  }

  .strip-btn:hover {
    background-color: var(--bg-hover, #2B2D30);
    color: var(--text-primary, #DFE1E5);
  }
</style>
