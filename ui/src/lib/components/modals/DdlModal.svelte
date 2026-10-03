<script lang="ts">
  import { appState, parseTableTabId } from '../../state/appState.svelte';
  import { Modal, Button, toast } from '@ohjanus/ui';

  let target = $derived.by(() => {
    const tab = appState.activeTab;
    if (tab?.type === 'table' && tab.connection && tab.schema && tab.table) {
      return { connection: tab.connection, schema: tab.schema, table: tab.table };
    }
    const parsed = parseTableTabId(appState.activeTabId);
    if (parsed) return parsed;
    const t = appState.tableViewer;
    if (t.connection && t.table) return { connection: t.connection, schema: t.schema, table: t.table };
    return null;
  });

  let ddl = $derived.by(() => {
    if (!target) return null;
    return appState.ddlFor(target.connection, target.schema, target.table);
  });

  function close() {
    appState.ddlModalOpen = false;
  }

  function copyDdl() {
    if (!ddl) return;
    navigator.clipboard.writeText(ddl);
    toast.success('DDL copied to clipboard');
  }
</script>

<Modal
  open={appState.ddlModalOpen}
  title={target ? `Table DDL: ${target.schema}.${target.table}` : 'Table DDL'}
  width="580px"
  onClose={close}
>
  {#if !target}
    <p class="ddl-empty">Open a table from the explorer first — DDL is generated from the live schema.</p>
  {:else if !ddl}
    <p class="ddl-empty">
      Schema for {target.connection}.{target.schema}.{target.table} is not loaded yet.
      Expand the connection in the explorer and retry.
    </p>
  {:else}
    <pre class="ddl-code code-text">{ddl}</pre>
  {/if}

  {#snippet footer()}
    <Button variant="secondary" onclick={close}>Close</Button>
    <Button variant="primary" onclick={copyDdl} disabled={!ddl}>Copy DDL</Button>
  {/snippet}
</Modal>

<style>
  .ddl-code {
    background-color: var(--bg-canvas, #1E1F22);
    border: 1px solid var(--border-default);
    padding: 12px;
    border-radius: 4px;
    font-size: 11px;
    line-height: 18px;
    color: var(--text-primary, #DFE1E5);
    max-height: 360px;
    overflow-y: auto;
    font-family: var(--font-mono, 'JetBrains Mono', monospace);
    margin: 0;
  }

  .ddl-empty {
    font-size: 12px;
    color: var(--text-secondary, #9DA0A8);
    margin: 0;
    padding: 8px 0;
  }
</style>
