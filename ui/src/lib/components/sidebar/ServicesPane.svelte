<script lang="ts">
  import { appState } from '../../state/appState.svelte';
  import { Box, Flex, Text, Badge, Button } from '@ohjanus/ui';
  import { Icon } from '@ohjanus/icons';

  interface Props {
    height: number;
  }

  let { height }: Props = $props();
</script>

<Box class="services-pane" style="height: {height}px;">
  <Flex align="center" gap="xs" class="services-header">
    <Text size="xs" weight="semibold">Services</Text>
    <Badge variant="counter" size="sm">{appState.services.length}</Badge>
  </Flex>

  <Box class="services-viewport">
    {#if appState.services.length === 0}
      <Box class="tree-empty">
        <Text size="xs" color="muted">
          No open sessions. Open a table or console from the explorer.
        </Text>
      </Box>
    {:else}
      {#each appState.services as session (session.id)}
        <Button
          variant="ghost"
          class="service-item {appState.activeTabId === session.id ? 'selected' : ''}"
          onclick={() => appState.activeTabId = session.id}
          title={session.connection ? `${session.name} · ${session.connection}` : session.name}
        >
          {#if session.type === 'table'}
            <Icon name="table" size={13} color="#4A88C7" class="node-icon" />
          {:else}
            <Icon name="lightning" size={13} color="#3B82F6" class="node-icon" />
          {/if}
          <Text size="sm" truncate class="service-name">{session.name}</Text>
          {#if session.durationMs !== undefined}
            <Text size="xs" color="muted" mono style="margin-left: auto;">
              {session.durationMs} ms
            </Text>
          {/if}
        </Button>
      {/each}
    {/if}
  </Box>
</Box>

<style>
  :global(.services-pane) {
    display: flex;
    flex-direction: column;
    background-color: var(--bg-card, #191A1C);
    border-radius: var(--radius-md, 8px);
    overflow: hidden;
    flex-shrink: 0;
  }

  :global(.services-header) {
    height: var(--toolbar-height, 32px);
    padding: 0 8px 0 12px;
    border-bottom: 1px solid var(--border-subtle);
    flex-shrink: 0;
  }

  :global(.services-viewport) {
    flex: 1;
    overflow-y: auto;
    overflow-x: hidden;
    padding: 4px 0;
  }

  :global(.service-item) {
    height: var(--tree-row-height, 28px) !important;
    display: flex !important;
    align-items: center !important;
    gap: 6px !important;
    padding: 0 8px !important;
    cursor: pointer !important;
    border-radius: var(--radius-sm, 4px) !important;
    margin: 1px 4px !important;
    text-align: left !important;
    width: calc(100% - 8px) !important;
    justify-content: flex-start !important;
    border: none !important;
    font-weight: normal !important;
  }

  :global(.service-item:hover) {
    background-color: var(--bg-hover) !important;
  }

  :global(.service-item.selected) {
    background-color: var(--bg-selected) !important;
    color: #FFFFFF !important;
  }

  :global(.service-name) {
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
</style>
