<script lang="ts">
  import { appState } from '../../state/appState.svelte';
  import { Card, Badge, Alert, Button, Text, Box, Flex, Grid, Stack, toast } from '@ohjanus/ui';
  import { Icon } from '@ohjanus/icons';
  import { Toolbar, ToolbarSeparator, BorderlessSelect } from '../toolbar';

  let selectedWindow = $state<string>('1h');

  async function handleReload() {
    try {
      await Promise.all([appState.loadSummary(), appState.loadApprovals(), appState.loadAudit()]);
      toast.info('Metrics Refreshed', 'Observability metrics and live KPIs updated.');
    } catch (e) {
      toast.error('Reload Failed', e instanceof Error ? e.message : String(e));
    }
  }
</script>

<Box class="dashboard-view">
  {#if appState.dataLoading}
    <Box style="padding: 8px 16px;"><Text size="sm" color="muted">Loading live data from Admin API…</Text></Box>
  {:else if appState.dataError}
    <Box style="padding: 8px 16px;">
      <Alert variant="danger" title="Admin API unreachable">
        <Text color="danger">{appState.dataError}</Text>
        <Button variant="secondary" size="sm" onclick={() => void appState.loadAll()}>Retry</Button>
      </Alert>
    </Box>
  {/if}
  <!-- Toolbar: Observability Actions (matching design2.png) -->
  <Toolbar>
    {#snippet left()}
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Reload metrics (Cmd+Enter)"
        onclick={handleReload}
      >
        <Icon name="refresh" size={13} />
      </Button>
      <ToolbarSeparator />
      <BorderlessSelect
        options={[
          { value: '1h', label: 'Window: Last 1 Hour' },
          { value: '24h', label: 'Window: Last 24 Hours' },
          { value: '7d', label: 'Window: Last 7 Days' },
        ]}
        bind:value={selectedWindow}
      />
      <ToolbarSeparator />
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Audit Trail"
        onclick={() => (appState.activeTabId = "audit")}
      >
        <Icon name="audit" size={13} />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Approval Queue"
        onclick={() => (appState.activeTabId = "approvals")}
      >
        <Icon name="shield" size={13} />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="SQL Console"
        onclick={() => (appState.activeTabId = "console")}
      >
        <Icon name="terminal" size={13} />
      </Button>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Client Tokens"
        onclick={() => (appState.activeTabId = "tokens")}
      >
        <Icon name="key" size={13} />
      </Button>
    {/snippet}

    {#snippet right()}
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Settings"
        onclick={() => (appState.settingsModalOpen = true)}
      >
        <Icon name="settings" size={13} />
      </Button>
    {/snippet}
  </Toolbar>

  <Stack class="dashboard-body" gap="16px">
    <!-- Top KPI Cards via @ohjanus/ui -->
    <Grid class="kpi-grid" columns="repeat(auto-fill, minmax(180px, 1fr))" gap="12px">
      <Card class="kpi-card-box">
        <Text class="kpi-label" size="sm" weight="semibold" color="muted">TOTAL REQUESTS (AUDIT)</Text>
        <Text class="kpi-val font-mono" weight="bold">{appState.summary ? appState.summary.requests_total.toLocaleString() : '…'}</Text>
        <Text class="kpi-trend pos" color="success">live from Admin API</Text>
      </Card>

      <Card class="kpi-card-box">
        <Text class="kpi-label" size="sm" weight="semibold" color="muted">TOTAL POLICY DENIALS</Text>
        <Text class="kpi-val font-mono warn" weight="bold" color="danger">{appState.summary ? appState.summary.denials_total.toLocaleString() : '…'}</Text>
        <Text class="kpi-sub" color="muted">Banned functions &amp; DDL blocked</Text>
      </Card>

      <Card class="kpi-card-box">
        <Text class="kpi-label" size="sm" weight="semibold" color="muted">PENDING APPROVALS</Text>
        <Text class="kpi-val font-mono pending" weight="bold" color="warning">{appState.notificationCount}</Text>
        <Text class="kpi-sub" color="muted">Human review required</Text>
      </Card>

      <Card class="kpi-card-box">
        <Text class="kpi-label" size="sm" weight="semibold" color="muted">P95 QUERY LATENCY</Text>
        <Text class="kpi-val font-mono" weight="bold">48 ms</Text>
        <Text class="kpi-sub" color="muted">Gateway overhead &lt; 2 ms</Text>
      </Card>

      <Card class="kpi-card-box">
        <Text class="kpi-label" size="sm" weight="semibold" color="muted">ACTIVE MCP TOKENS</Text>
        <Text class="kpi-val font-mono" weight="bold">{appState.tokens.filter(t => t.state === 'active').length}</Text>
        <Text class="kpi-sub" color="muted">Cursor &amp; Claude Desktop</Text>
      </Card>

      <Card class="kpi-card-box">
        <Text class="kpi-label" size="sm" weight="semibold" color="muted">HEALTHY POOLS</Text>
        <Text class="kpi-val font-mono pos" weight="bold" color="success">3 / 3</Text>
        <Text class="kpi-sub" color="muted">All connections healthy</Text>
      </Card>
    </Grid>

    <!-- Middle Split: Chart Simulation & Denials List -->
    <Grid class="mid-split" columns="3fr 2fr" gap="16px">
      <Card class="chart-panel">
        {#snippet header()}
          <Flex class="panel-header-inner" justify="between" align="center">
            <Text size="xs" weight="semibold" color="muted">REAL-TIME QUERY THROUGHPUT (REQS / SEC)</Text>
            <Badge variant="success" size="sm">
              <Box as="span" class="status-dot ok" />LIVE
            </Badge>
          </Flex>
        {/snippet}
        <Flex class="chart-bars font-mono" align="end" gap="6px">
          {#each [18, 25, 42, 38, 55, 62, 45, 80, 72, 95, 88, 64, 52, 47, 63, 78, 92, 105, 84, 76] as h, i}
            <Flex class="bar-col" direction="column" align="center" justify="end">
              <Box class="bar-fill" style="height: {h}%;" />
              <Text class="bar-label" size="sm" color="muted">{i * 3}m</Text>
            </Flex>
          {/each}
        </Flex>
      </Card>

      <Card class="security-panel">
        {#snippet header()}
          <Flex class="panel-header-inner" justify="between" align="center">
            <Text size="xs" weight="semibold" color="muted">RECENT BLOCKED OPERATIONS (AST VALIDATION)</Text>
          </Flex>
        {/snippet}
        <Stack class="denials-list font-mono" gap="8px">
          <Box class="denial-item">
            <Flex class="denial-top" justify="between" align="center">
              <Badge variant="danger" size="sm">banned-function</Badge>
              <Text class="denial-time" size="sm" color="muted">10:45:12</Text>
            </Flex>
            <Text class="denial-sql" mono size="xs">SELECT pg_read_file('/etc/passwd')</Text>
            <Text class="denial-agent" size="xs" color="muted">Agent: Claude Desktop v1.2</Text>
          </Box>

          <Box class="denial-item">
            <Flex class="denial-top" justify="between" align="center">
              <Badge variant="danger" size="sm">deny-ddl</Badge>
              <Text class="denial-time" size="sm" color="muted">09:12:33</Text>
            </Flex>
            <Text class="denial-sql" mono size="xs">DROP TABLE user_entitlements CASCADE;</Text>
            <Text class="denial-agent" size="xs" color="muted">Agent: Cursor AI Agent</Text>
          </Box>

          <Box class="denial-item">
            <Flex class="denial-top" justify="between" align="center">
              <Badge variant="danger" size="sm">denied-table</Badge>
              <Text class="denial-time" size="sm" color="muted">08:04:19</Text>
            </Flex>
            <Text class="denial-sql" mono size="xs">SELECT * FROM admin_credentials LIMIT 10;</Text>
            <Text class="denial-agent" size="xs" color="muted">Agent: Cursor AI Agent</Text>
          </Box>
        </Stack>
      </Card>
    </Grid>
  </Stack>
</Box>

<style>
  :global(.dashboard-view) {
    display: flex;
    flex-direction: column;
    height: 100%;
    background: var(--bg-canvas);
    overflow-y: auto;
  }

  :global(.dashboard-view .dashboard-body) {
    padding: 16px;
  }

  :global(.kpi-card-box) {
    display: flex;
    flex-direction: column;
    gap: 6px;
    padding: 14px 16px;
  }

  :global(.dashboard-view .kpi-label) { font-size: var(--font-size-2xs, 11px); font-weight: 600; color: var(--text-muted); letter-spacing: 0.5px; }
  :global(.dashboard-view .kpi-val) { font-size: var(--font-size-2xl, 26px); font-weight: 700; color: var(--text-primary); white-space: nowrap; }
  :global(.dashboard-view .kpi-val.pos) { color: var(--action-success); }
  :global(.dashboard-view .kpi-val.warn) { color: var(--action-danger); }
  :global(.dashboard-view .kpi-val.pending) { color: var(--action-warning); }
  :global(.dashboard-view .kpi-trend.pos) { font-size: var(--font-size-xs, 12px); color: var(--action-success); font-weight: 500; }
  :global(.dashboard-view .kpi-sub) { font-size: var(--font-size-xs, 12px); color: var(--text-muted); }

  :global(.dashboard-view .panel-header-inner) {
    width: 100%;
    font-size: var(--font-size-xs, 12px);
    font-weight: 600;
    color: var(--text-muted);
  }

  :global(.chart-bars) {
    height: 140px;
    padding-top: 16px;
  }

  :global(.bar-col) {
    flex: 1;
    height: 100%;
  }

  :global(.bar-fill) {
    width: 100%;
    background: var(--action-primary);
    border-radius: 2px 2px 0 0;
    transition: height 0.3s ease;
    opacity: 0.85;
  }

  :global(.bar-fill:hover) {
    opacity: 1;
  }

  :global(.dashboard-view .bar-label) {
    font-size: var(--font-size-2xs, 11px);
    color: var(--text-muted);
    margin-top: 4px;
  }

  :global(.denial-item) {
    background: rgba(0, 0, 0, 0.15);
    border: 1px solid var(--border-subtle);
    border-radius: 4px;
    padding: 8px 10px;
  }

  :global(.denial-top) {
    margin-bottom: 4px;
  }

  :global(.dashboard-view .denial-time) { font-size: var(--font-size-2xs, 11px); color: var(--text-muted); }
  :global(.denial-sql) { font-size: var(--font-size-xs, 12px); color: var(--action-danger); margin-bottom: 2px; }
  :global(.denial-agent) { font-size: var(--font-size-2xs, 11px); color: var(--text-muted); }
</style>
