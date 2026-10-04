<script lang="ts">
  import { bootState } from '@/app/boot.svelte';
  import { connectionsState } from '@/features/connections';
  import { consoleState } from '@/features/console';
  import { explorerState } from '@/features/explorer';
  import { tableViewerState } from '@/features/table-viewer';
  import { Box, Stack, Text, Button } from '@ohjanus/ui';
  import ExplorerTreeItem from './ExplorerTreeItem.svelte';

  function connKey(name: string) {
    return `conn:${name}`;
  }

  function schemaKey(conn: string, schema: string) {
    return `conn:${conn}:schema:${schema}`;
  }

  async function toggleConnection(name: string) {
    explorerState.toggleTree(connKey(name));
    if (explorerState.treeExpanded[connKey(name)]) {
      await explorerState.loadSchema(name);
    }
  }

  function openConsole(name: string) {
    explorerState.selectedTreeNode = `console:${name}`;
    consoleState.openConsole(name);
  }

  async function openTable(conn: string, schema: string, table: string) {
    await tableViewerState.openTable(conn, schema, table);
  }

  function matchesFilter(text: string): boolean {
    const q = explorerState.treeFilterQuery.trim().toLowerCase();
    return q === '' || text.toLowerCase().includes(q);
  }
</script>

<Box class="tree-viewport">
  {#if bootState.dataLoading}
    <Box class="tree-empty">
      <Text size="xs" color="muted">Loading connections from Admin API…</Text>
    </Box>
  {:else if bootState.dataError && connectionsState.connections.length === 0}
    <Stack gap="xs" class="tree-empty">
      <Text size="xs" color="muted">Admin API unreachable.</Text>
      <Text size="xs" color="danger" class="tree-error">{bootState.dataError}</Text>
      <Button
        variant="link"
        size="xs"
        class="retry-btn"
        onclick={() => void bootState.loadAll()}
      >
        Retry
      </Button>
    </Stack>
  {:else if connectionsState.connections.length === 0}
    <Box class="tree-empty">
      <Text size="xs" color="muted">
        No connections configured. Add one in janus.yaml and restart serve.
      </Text>
    </Box>
  {:else}
    {#each explorerState.explorerConnections as conn (conn.name)}
      <ExplorerTreeItem
        depth={0}
        icon="database"
        iconColor="#3B82F6"
        label={conn.name}
        badgeText={conn.readonly ? '[ReadOnly]' : undefined}
        hasChevron={true}
        isExpanded={explorerState.treeExpanded[connKey(conn.name)]}
        warnStatus={conn.status !== 'healthy' ? conn.status : undefined}
        onclick={() => void toggleConnection(conn.name)}
      />

      {#if explorerState.treeExpanded[connKey(conn.name)]}
        <ExplorerTreeItem
          depth={1}
          icon="lightning"
          iconColor="#3B82F6"
          label={`console [${conn.name}]`}
          isSelected={explorerState.selectedTreeNode === `console:${conn.name}`}
          onclick={() => openConsole(conn.name)}
        />

        {#if explorerState.schemaLoading[conn.name]}
          <Box class="tree-empty">
            <Text size="xs" color="muted">Loading schema…</Text>
          </Box>
        {:else if explorerState.schemaError[conn.name]}
          <Stack gap="xs" class="tree-empty">
            <Text size="xs" color="danger" class="tree-error">{explorerState.schemaError[conn.name]}</Text>
            <Button
              variant="link"
              size="xs"
              class="retry-btn"
              onclick={() => void explorerState.loadSchema(conn.name, true)}
            >
              Retry
            </Button>
          </Stack>
        {:else}
          {#each explorerState.schemasOf(conn.name) as schema (schema.name)}
            {#if matchesFilter(schema.name) || schema.tables.some((t) => matchesFilter(t.name))}
              <ExplorerTreeItem
                depth={1}
                icon="folder"
                iconColor="#C29D38"
                label={schema.name}
                badgeCount={schema.tables.length}
                hasChevron={true}
                isExpanded={explorerState.treeExpanded[schemaKey(conn.name, schema.name)]}
                onclick={() => explorerState.toggleTree(schemaKey(conn.name, schema.name))}
              />

              {#if explorerState.treeExpanded[schemaKey(conn.name, schema.name)]}
                {#each schema.tables.filter((t) => matchesFilter(t.name)) as table (table.name)}
                  <ExplorerTreeItem
                    depth={2}
                    icon="table"
                    iconColor="#4A88C7"
                    label={table.name}
                    isSelected={explorerState.selectedTreeNode === `table:${conn.name}.${schema.name}.${table.name}`}
                    onclick={() => void openTable(conn.name, schema.name, table.name)}
                  />
                {/each}
              {/if}
            {/if}
          {/each}
        {/if}
      {/if}
    {/each}
  {/if}
</Box>

<style>
  :global(.tree-viewport) {
    flex: 1;
    overflow-y: auto;
    overflow-x: hidden;
    padding: 4px 0;
  }

  :global(.tree-empty) {
    padding: 12px;
    line-height: 1.5;
  }

  :global(.tree-error) {
    word-break: break-word;
  }

  :global(.retry-btn) {
    padding: 0 !important;
    height: auto !important;
    font-size: 12px !important;
    justify-content: flex-start !important;
  }
</style>
