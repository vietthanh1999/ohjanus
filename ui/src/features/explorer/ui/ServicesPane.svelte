<script lang="ts">
  import { formatLatency } from '@/shared/lib';
  import { explorerState } from '@/features/explorer';
  import { workbenchState } from '@/features/workbench';
  import { Box, Flex, Text, Badge, Button, Checkbox } from '@ohjanus/ui';
  import { Icon } from '@ohjanus/icons';

  interface Props {
    height: number;
  }

  let { height }: Props = $props();

  // Local collapse of the Database group (DESIGN §3.2 tree toggle).
  let groupCollapsed = $state(false);

  let sessions = $derived(workbenchState.services);
  let selectedCount = $derived(workbenchState.selectedSessionIds.length);
  let allChecked = $derived(
    sessions.length > 0 && sessions.every((s) => workbenchState.selectedSessionIds.includes(s.id))
  );

  function toggleAll() {
    if (allChecked) workbenchState.clearSessionSelection();
    else {
      workbenchState.selectedSessionIds = sessions.map((s) => s.id);
    }
  }

  function openSession(id: string) {
    workbenchState.activeTabId = id;
  }
</script>

<Box class="services-pane" style="height: {height}px;">
  <Flex align="center" gap="xs" class="services-header">
    <Text size="xs" weight="semibold">Services</Text>
    <Badge variant="counter" size="sm">{sessions.length}</Badge>
    <Text size="xs" color="muted" mono style="margin-left: 2px;">Tx: Auto</Text>
    <span class="services-spacer"></span>
    <Button
      variant="ghost"
      size="icon-xs"
      class="jb-icon-btn"
      title={selectedCount > 0 ? `Close ${selectedCount} selected session(s)` : 'Select sessions to close them'}
      disabled={selectedCount === 0}
      onclick={() => workbenchState.closeSessions(workbenchState.selectedSessionIds)}
    >
      <Icon name="x" size={11} />
    </Button>
    <Button
      variant="ghost"
      size="icon-xs"
      class="jb-icon-btn"
      title="Collapse all"
      onclick={() => (groupCollapsed = true)}
    >
      <Icon name="chevron-up" size={11} />
    </Button>
  </Flex>

  <Box class="services-viewport">
    {#if sessions.length === 0}
      <Box class="tree-empty">
        <Text size="xs" color="muted">
          No open sessions. Open a table or console from the explorer.
        </Text>
      </Box>
    {:else}
      <Button
        variant="ghost"
        class="service-group"
        onclick={() => (groupCollapsed = !groupCollapsed)}
      >
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <span
          class="check-wrap"
          onclick={(e) => e.stopPropagation()}
          onkeydown={(e) => e.stopPropagation()}
        >
          <Checkbox
            checked={allChecked}
            ariaLabel="Select all sessions"
            onchange={() => toggleAll()}
          />
        </span>
        <Box class="chevron {groupCollapsed ? '' : 'expanded'}">
          <Icon name="chevron-right" size={12} />
        </Box>
        <Icon name="folder" size={13} color="#C29D38" class="node-icon" />
        <Text size="sm">Database</Text>
      </Button>

      {#if !groupCollapsed}
        {#each sessions as session (session.id)}
          {@const checked = workbenchState.selectedSessionIds.includes(session.id)}
          <Button
            variant="ghost"
            class="service-item {workbenchState.activeTabId === session.id ? 'selected' : ''}"
            onclick={() => openSession(session.id)}
            title={session.connection ? `${session.name} · ${session.connection}` : session.name}
          >
            <!-- svelte-ignore a11y_no_static_element_interactions -->
            <span
              class="check-wrap"
              onclick={(e) => e.stopPropagation()}
              onkeydown={(e) => e.stopPropagation()}
            >
              <Checkbox
                checked={checked}
                ariaLabel={`Select session ${session.name}`}
                onchange={() => workbenchState.toggleSession(session.id)}
              />
            </span>
            {#if session.type === 'table'}
              <Icon name="table" size={13} color="#4A88C7" class="node-icon" />
            {:else}
              <Icon name="lightning" size={13} color="#3B82F6" class="node-icon" />
            {/if}
            <Text size="sm" truncate class="service-name">{session.name}</Text>
            {#if session.durationMs !== undefined}
              <Text size="xs" color="muted" mono style="margin-left: auto;">
                {formatLatency(session.durationMs)}
              </Text>
            {/if}
          </Button>
        {/each}
      {/if}
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

  :global(.services-header .jb-icon-btn) {
    width: 22px !important;
    height: 22px !important;
    padding: 0 !important;
    min-width: 22px !important;
    color: var(--text-secondary);
  }

  :global(.services-header .jb-icon-btn:hover:not(:disabled)) {
    color: var(--text-primary);
    background-color: var(--bg-hover);
  }

  :global(.services-spacer) {
    flex: 1;
  }

  :global(.services-viewport) {
    flex: 1;
    overflow-y: auto;
    overflow-x: hidden;
    padding: 4px 0;
  }

  :global(.service-group),
  :global(.service-item) {
    min-height: var(--tree-row-height, 28px) !important;
    display: flex !important;
    align-items: center !important;
    gap: 6px !important;
    padding: 2px 8px !important;
    cursor: pointer !important;
    border-radius: var(--radius-sm, 4px) !important;
    margin: 1px 4px !important;
    text-align: left !important;
    width: calc(100% - 8px) !important;
    justify-content: flex-start !important;
    border: none !important;
    font-weight: normal !important;
  }

  :global(.service-group:hover),
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

  :global(.service-group .chevron) {
    width: 14px;
    height: 14px;
    color: var(--text-muted);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    transition: transform 0.15s ease;
    flex-shrink: 0;
  }

  :global(.service-group .chevron.expanded) {
    transform: rotate(90deg);
  }

  :global(.check-wrap) {
    display: inline-flex;
    align-items: center;
    flex-shrink: 0;
  }
</style>
