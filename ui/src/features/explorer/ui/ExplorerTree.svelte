<script lang="ts">
  import { connectionsState } from '@/features/connections';
  import { consoleState } from '@/features/console';
  import { connKey, explorerState, schemaKey } from '@/features/explorer';
  import { tableViewerState } from '@/features/table-viewer';
  import { Box, Stack, Text, Button, toast } from '@ohjanus/ui';
  import ExplorerTreeItem from './ExplorerTreeItem.svelte';

  async function toggleConnection(name: string) {
    explorerState.toggleTree(connKey(name));
    if (explorerState.treeExpanded[connKey(name)]) {
      await explorerState.loadSchema(name);
    }
  }

  function openConsole(name: string) {
    explorerState.revealConsole(name);
    consoleState.openConsole(name);
  }

  async function openTable(conn: string, schema: string, table: string) {
    await tableViewerState.openTable(conn, schema, table);
  }

  function matchesFilter(text: string): boolean {
    const q = explorerState.treeFilterQuery.trim().toLowerCase();
    return q === '' || text.toLowerCase().includes(q);
  }

  /** Pill for a connection row: schema count, or `…` while unloaded (DESIGN §3.1). */
  function connectionBadge(name: string): string {
    if (explorerState.schemaLoading[name] || !explorerState.schemas[name]) return '…';
    return String(explorerState.visibleSchemas(name).length);
  }

  function copyText(text: string, what: string) {
    navigator.clipboard.writeText(text);
    toast.success(`${what} copied to clipboard`);
  }

  function copyTableDdl(conn: string, schema: string, table: string) {
    const ddl = explorerState.ddlFor(conn, schema, table);
    if (!ddl) {
      toast.info('Copy DDL', 'Schema is not loaded yet. Expand the connection and retry.');
      return;
    }
    navigator.clipboard.writeText(ddl);
    toast.success('DDL copied to clipboard');
  }
</script>

<Box class="tree-viewport">
  {#if connectionsState.loading}
    <Box class="tree-empty">
      <Text size="xs" color="muted">Loading connections from Admin API…</Text>
    </Box>
  {:else if connectionsState.error && connectionsState.connections.length === 0}
    <Stack gap="xs" class="tree-empty">
      <Text size="xs" color="muted">Admin API unreachable.</Text>
      <Text size="xs" color="danger" class="tree-error">{connectionsState.error}</Text>
      <Button
        variant="link"
        size="xs"
        class="retry-btn"
        onclick={() => void connectionsState.loadConnections()}
      >
        Retry
      </Button>
    </Stack>
  {:else if connectionsState.connections.length === 0}
    <Box class="tree-empty">
      <Text size="xs" color="muted">
        No connections yet. Press + in the toolbar to add a PostgreSQL connection.
      </Text>
    </Box>
  {:else}
    {#each explorerState.explorerConnections as conn (conn.name)}
      <ExplorerTreeItem
        depth={0}
        icon="database"
        iconColor="#4A88C7"
        label={conn.name}
        locked={conn.readonly}
        badge={connectionBadge(conn.name)}
        hasChevron={true}
        isExpanded={explorerState.treeExpanded[connKey(conn.name)]}
        warnStatus={conn.status !== 'healthy' ? conn.status : undefined}
        title={conn.readonly ? `${conn.name} (read-only)` : conn.name}
        actions={[
          { icon: 'lightning', label: 'New console', onclick: () => openConsole(conn.name) },
          { icon: 'refresh', label: 'Refresh schemas', onclick: () => void explorerState.loadSchema(conn.name, true) },
          { icon: 'copy', label: 'Copy name', onclick: () => copyText(conn.name, 'Connection name') }
        ]}
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
          {#each explorerState.visibleSchemas(conn.name) as schema (schema.name)}
            {@const tables = schema.tables.filter((t) => matchesFilter(t.name))}
            {@const routines = schema.routines.filter((r) => matchesFilter(r.name))}
            {@const sequences = schema.sequences.filter((q) => matchesFilter(q.name))}
            {#if matchesFilter(schema.name) || tables.length > 0 || routines.length > 0 || sequences.length > 0}
              <ExplorerTreeItem
                depth={1}
                icon="folder"
                iconColor="#C29D38"
                label={schema.name}
                badge={String(schema.tables.length)}
                hasChevron={true}
                isExpanded={explorerState.treeExpanded[schemaKey(conn.name, schema.name)]}
                title={`${conn.name}.${schema.name} (${schema.tables.length} tables)`}
                actions={[
                  { icon: 'refresh', label: 'Refresh schemas', onclick: () => void explorerState.loadSchema(conn.name, true) },
                  { icon: 'copy', label: 'Copy qualified name', onclick: () => copyText(`${conn.name}.${schema.name}`, 'Schema name') }
                ]}
                onclick={() => explorerState.toggleTree(schemaKey(conn.name, schema.name))}
              />

              {#if explorerState.treeExpanded[schemaKey(conn.name, schema.name)]}
                {@const baseKey = schemaKey(conn.name, schema.name)}
                {#if tables.length > 0}
                  <ExplorerTreeItem
                    depth={2}
                    icon="table"
                    iconColor="#4A88C7"
                    label="tables"
                    badge={String(tables.length)}
                    hasChevron={true}
                    isExpanded={explorerState.treeExpanded[`${baseKey}:tables`]}
                    title={`${tables.length} tables in ${conn.name}.${schema.name}`}
                    onclick={() => explorerState.toggleTree(`${baseKey}:tables`)}
                  />
                  {#if explorerState.treeExpanded[`${baseKey}:tables`]}
                    {#each tables as table (table.name)}
                      <ExplorerTreeItem
                        depth={3}
                        icon="table"
                        iconColor="#4A88C7"
                        label={table.name}
                        isSelected={explorerState.selectedTreeNode === `table:${conn.name}.${schema.name}.${table.name}`}
                        title={`${conn.name}.${schema.name}.${table.name}`}
                        actions={[
                          { icon: 'table', label: 'Open table', onclick: () => void openTable(conn.name, schema.name, table.name) },
                          { icon: 'copy', label: 'Copy DDL', onclick: () => copyTableDdl(conn.name, schema.name, table.name) },
                          { icon: 'copy', label: 'Copy qualified name', onclick: () => copyText(`${conn.name}.${schema.name}.${table.name}`, 'Table name') }
                        ]}
                        onclick={() => void openTable(conn.name, schema.name, table.name)}
                      />
                    {/each}
                  {/if}
                {/if}

                {#if routines.length > 0}
                  <ExplorerTreeItem
                    depth={2}
                    icon="terminal"
                    iconColor="#57D38C"
                    label="routines"
                    badge={String(routines.length)}
                    hasChevron={true}
                    isExpanded={explorerState.treeExpanded[`${baseKey}:routines`]}
                    title={`${routines.length} routines in ${conn.name}.${schema.name}`}
                    onclick={() => explorerState.toggleTree(`${baseKey}:routines`)}
                  />
                  {#if explorerState.treeExpanded[`${baseKey}:routines`]}
                    {#each routines as routine (routine.name)}
                      <ExplorerTreeItem
                        depth={3}
                        icon="terminal"
                        iconColor="#57D38C"
                        label={routine.name}
                        isSelected={explorerState.selectedTreeNode === `routine:${conn.name}.${schema.name}.${routine.name}`}
                        title={`${routine.kind} ${conn.name}.${schema.name}.${routine.name}`}
                        actions={[
                          { icon: 'copy', label: 'Copy qualified name', onclick: () => copyText(`${conn.name}.${schema.name}.${routine.name}`, 'Routine name') }
                        ]}
                        onclick={() => (explorerState.selectedTreeNode = `routine:${conn.name}.${schema.name}.${routine.name}`)}
                      />
                    {/each}
                  {/if}
                {/if}

                {#if sequences.length > 0}
                  <ExplorerTreeItem
                    depth={2}
                    icon="sort"
                    iconColor="#9DA0A8"
                    label="sequences"
                    badge={String(sequences.length)}
                    hasChevron={true}
                    isExpanded={explorerState.treeExpanded[`${baseKey}:sequences`]}
                    title={`${sequences.length} sequences in ${conn.name}.${schema.name}`}
                    onclick={() => explorerState.toggleTree(`${baseKey}:sequences`)}
                  />
                  {#if explorerState.treeExpanded[`${baseKey}:sequences`]}
                    {#each sequences as sequence (sequence.name)}
                      <ExplorerTreeItem
                        depth={3}
                        icon="sort"
                        iconColor="#9DA0A8"
                        label={sequence.name}
                        isSelected={explorerState.selectedTreeNode === `sequence:${conn.name}.${schema.name}.${sequence.name}`}
                        title={`${conn.name}.${schema.name}.${sequence.name}`}
                        actions={[
                          { icon: 'copy', label: 'Copy qualified name', onclick: () => copyText(`${conn.name}.${schema.name}.${sequence.name}`, 'Sequence name') }
                        ]}
                        onclick={() => (explorerState.selectedTreeNode = `sequence:${conn.name}.${schema.name}.${sequence.name}`)}
                      />
                    {/each}
                  {/if}
                {/if}
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
