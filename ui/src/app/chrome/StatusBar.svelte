<script lang="ts">
  import { parseConsoleTabId, parseTableTabId } from '@/entities/tab';
  import { approvalsState } from '@/features/approvals';
  import { connectionsState } from '@/features/connections';
  import { consoleState } from '@/features/console';
  import { tableViewerState } from '@/features/table-viewer';
  import { workbenchState } from '@/features/workbench';
  import { Flex, Badge, Text } from '@ohjanus/ui';
  import { Icon } from '@ohjanus/icons';

  let activeTab = $derived(workbenchState.activeTab);
  let tableCoords = $derived(
    activeTab?.type === 'table' && activeTab.connection
      ? { connection: activeTab.connection, schema: activeTab.schema ?? '', table: activeTab.table ?? '' }
      : parseTableTabId(workbenchState.activeTabId)
  );
  let consoleCoords = $derived(
    activeTab?.type === 'console' ? parseConsoleTabId(activeTab.id) : null
  );
  let activeConnection = $derived(
    tableCoords?.connection ?? consoleCoords?.connection ?? consoleState.console.connection ?? ''
  );
  let connectionItem = $derived(connectionsState.connections.find((c) => c.name === activeConnection));
  let cursorStats = $derived(consoleState.consoleCursorStats);

  function openApprovals() {
    workbenchState.openTab({ id: 'approvals', title: 'Approvals Queue', type: 'approvals', closable: false, icon: 'shield' });
  }
</script>

<footer class="statusbar">
  <Flex class="breadcrumb-strip" align="center" gap="4px">
    {#if tableCoords}
      <span class="crumb-item">Database</span>
      <Text size="xs" color="muted">&gt;</Text>
      <span class="crumb-item">{tableCoords.connection}</span>
      <Text size="xs" color="muted">&gt;</Text>
      <span class="crumb-item">{tableCoords.schema}</span>
      <Text size="xs" color="muted">&gt;</Text>
      <span class="crumb-item">tables</span>
      <Text size="xs" color="muted">&gt;</Text>
      <span class="crumb-item active-crumb">
        <Icon name="table" size={11} color="#4A88C7" style="margin-right: 4px;" />
        {tableCoords.table}
      </span>
    {:else if consoleCoords}
      <span class="crumb-item">Database Consoles</span>
      <Text size="xs" color="muted">&gt;</Text>
      <span class="crumb-item">{consoleCoords.connection}</span>
      <Text size="xs" color="muted">&gt;</Text>
      <span class="crumb-item active-crumb">
        <Icon name="lightning" size={11} color="#57D38C" style="margin-right: 4px;" />
        console [{consoleCoords.connection}]
      </span>
    {:else}
      <span class="crumb-item">MCP Gateway Security</span>
      <Text size="xs" color="muted">&gt;</Text>
      <span class="crumb-item active-crumb">{workbenchState.activeTabId.toUpperCase()}</span>
    {/if}
  </Flex>

  <Flex class="status-right" align="center" gap="12px">
    {#if activeTab?.type === 'console'}
      <span class="status-item code-text" title="SQL buffer stats">
        {cursorStats.lines} lines ({cursorStats.chars} chars)
      </span>
      <span class="status-item" title="Line Endings">LF</span>
      <span class="status-item" title="File Encoding">UTF-8</span>
      {#if connectionItem?.readonly}
        <span class="status-item lock-icon" title="Connection is read-only">
          <Icon name="lock" size={11} />
        </span>
      {/if}
    {:else if tableCoords}
      <span class="status-item code-text" title="Last execution">
        {tableViewerState.tableViewer.rowCount} row(s){tableViewerState.tableViewer.truncated ? ' (truncated)' : ''} · {tableViewerState.tableViewer.durationMs} ms
      </span>
      {#if connectionItem?.readonly}
        <span class="status-item lock-icon" title="Connection is read-only">
          <Icon name="lock" size={11} />
        </span>
      {/if}
    {:else}
      <span class="status-item code-text" title="Connections">
        {connectionsState.connections.length} connection(s)
      </span>
    {/if}

    <button
      class="status-item bell-btn"
      class:has-notification={approvalsState.notificationCount > 0}
      title="Pending Approvals ({approvalsState.notificationCount})"
      onclick={openApprovals}
    >
      <Icon name="bell" size={12} />
      {#if approvalsState.notificationCount > 0}
        <span class="notification-badge">{approvalsState.notificationCount}</span>
      {/if}
    </button>

    <Badge variant="license" size="sm" class="license-pill">
      Non-commercial use
    </Badge>
  </Flex>
</footer>

<style>
  .statusbar {
    height: var(--statusbar-height);
    background-color: var(--bg-window-frame, #24272A);
    border-top: none;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 10px;
    font-size: 11px;
    color: var(--text-secondary);
    user-select: none;
    z-index: 50;
    flex-shrink: 0;
  }

  .crumb-item {
    display: flex;
    align-items: center;
    padding: 1px 4px;
    border-radius: 2px;
    cursor: pointer;
    color: var(--text-secondary);
  }

  .crumb-item:hover {
    background-color: var(--bg-hover);
    color: var(--text-primary);
  }

  .crumb-item.active-crumb {
    color: var(--text-primary);
  }

  .status-item {
    padding: 1px 4px;
    border-radius: 2px;
    cursor: pointer;
    color: var(--text-muted);
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .status-item:hover {
    background-color: var(--bg-hover);
    color: var(--text-primary);
  }

  .lock-icon {
    color: var(--text-muted);
  }

  .bell-btn {
    position: relative;
    color: var(--text-muted);
  }

  .bell-btn.has-notification {
    color: #EDA200;
  }

  .notification-badge {
    position: absolute;
    top: -3px;
    right: -4px;
    background-color: #EDA200;
    color: #1E1F22;
    font-size: 8px;
    font-weight: bold;
    padding: 0 3px;
    border-radius: 5px;
    line-height: 11px;
  }

  :global(.license-pill) {
    font-size: 10px !important;
  }
</style>
