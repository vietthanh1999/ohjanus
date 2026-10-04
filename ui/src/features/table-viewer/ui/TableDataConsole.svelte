<script lang="ts">
  import { tableViewerState } from '@/features/table-viewer';
  import { workbenchState } from '@/features/workbench';
  import { Icon } from '@ohjanus/icons';
  import { Text, Box, Flex, Button } from '@ohjanus/ui';
</script>

<Flex class="log-console-container">
  <Box class="log-content code-text">
    {#if workbenchState.consoleLogs.length === 0}
      <Flex gap="8px" class="log-line">
        <Text size="lg" color="muted" mono>--</Text>
        <Text size="lg" mono>No output yet. Run a query or open a table to see live execution logs.</Text>
      </Flex>
    {:else}
      {#each workbenchState.consoleLogs as log (log.id)}
        <Flex gap="8px" class="log-line">
          <Text size="lg" color="muted" mono>[{log.timestamp}]</Text>
          {#if log.connection}
            <Text size="lg" weight="medium" mono>{log.connection}&gt;</Text>
          {/if}
          {#if log.querySnippet}
            <Text size="lg" mono>{log.querySnippet}</Text>
          {:else}
            <Text
              size="lg"
              mono
              color={log.type === 'error' ? 'danger' : 'default'}
            >{log.summary}</Text>
          {/if}
        </Flex>
        {#if log.querySnippet}
          <Flex gap="8px" class="log-line indent-sql">
            <Text
              size="lg"
              mono
              color={log.type === 'error' ? 'danger' : 'default'}
            >{log.summary}</Text>
          </Flex>
        {/if}
      {/each}
    {/if}
  </Box>

  <Flex direction="column" align="center" gap="4px" class="log-action-strip">
    <Button
      variant="ghost"
      size="icon-sm"
      class="strip-btn"
      title="Reload table data"
      onclick={() => void tableViewerState.loadTableData()}
    >
      <Icon name="refresh" size={13} />
    </Button>
    <Button
      variant="ghost"
      size="icon-sm"
      class="strip-btn"
      title="Clear Console"
      onclick={() => workbenchState.clearLogs()}
    >
      <Icon name="trash" size={13} />
    </Button>
  </Flex>
</Flex>

<style>
  :global(.log-console-container) {
    width: 100%;
    height: 100%;
    background-color: var(--bg-canvas, #1E1F22);
    overflow: hidden;
  }

  :global(.log-content) {
    flex: 1;
    overflow-y: auto;
    padding: 10px 14px;
    font-size: var(--font-size-base, 14px);
    line-height: 20px;
    font-family: var(--font-code, 'JetBrains Mono', monospace);
  }

  :global(.log-line.indent-sql) {
    padding-left: 24px;
  }

  :global(.log-action-strip) {
    width: 32px;
    background-color: var(--bg-canvas, #1E1F22);
    border-left: 1px solid var(--border-subtle, #323438);
    padding: 4px 0;
    flex-shrink: 0;
  }

  :global(.strip-btn) {
    width: var(--icon-btn-size-sm, 26px) !important;
    height: var(--icon-btn-size-sm, 26px) !important;
  }
</style>
