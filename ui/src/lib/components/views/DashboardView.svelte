<script lang="ts">
  import { appState } from '../../state/appState.svelte';
  import { Card, Badge, Alert, Button, Text, Box, Flex, Grid, Stack } from '@ohjanus/ui';
  import { Icon } from '@ohjanus/icons';
</script>

<Box class="dashboard-view">
  {#if appState.dataLoading}
    <Box style="padding: 8px 16px;"><Text size="sm" color="muted">Loading live data from Admin API…</Text></Box>
  {:else if appState.dataError}
    <Box style="padding: 8px 16px;">
      <Alert variant="danger" title="Admin API unreachable">
        <p>{appState.dataError}</p>
        <Button variant="secondary" size="sm" onclick={() => void appState.loadAll()}>Retry</Button>
      </Alert>
    </Box>
  {/if}
  <Flex as="header" class="view-header" align="center">
    <Flex align="center" gap="8px">
      <Icon name="chart" size={16} color="#3574F0" />
      <Text class="view-title" size="xl" weight="semibold">OhJanus MCP Gateway Observability &amp; Metrics</Text>
    </Flex>
  </Flex>

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
            <span>REAL-TIME QUERY THROUGHPUT (REQS / SEC)</span>
            <Badge variant="success" size="sm"><span class="status-dot ok" aria-hidden="true"></span>LIVE</Badge>
          </Flex>
        {/snippet}
        <Flex class="chart-bars font-mono" align="end" gap="6px">
          {#each [18, 25, 42, 38, 55, 62, 45, 80, 72, 95, 88, 64, 52, 47, 63, 78, 92, 105, 84, 76] as h, i}
            <Flex class="bar-col" direction="column" align="center" justify="end">
              <Box class="bar-fill" style="height: {h}%;"></Box>
              <Text class="bar-label" size="sm" color="muted">{i * 3}m</Text>
            </Flex>
          {/each}
        </Flex>
      </Card>

      <Card class="security-panel">
        {#snippet header()}
          <Flex class="panel-header-inner" justify="between" align="center">
            <span>RECENT BLOCKED OPERATIONS (AST VALIDATION)</span>
          </Flex>
        {/snippet}
        <Stack class="denials-list font-mono" gap="8px">
          <Box class="denial-item">
            <Flex class="denial-top" justify="between" align="center">
              <Badge variant="danger" size="sm">banned-function</Badge>
              <Text class="denial-time" size="sm" color="muted">10:45:12</Text>
            </Flex>
            <Box class="denial-sql">SELECT pg_read_file('/etc/passwd')</Box>
            <Box class="denial-agent">Agent: Claude Desktop v1.2</Box>
          </Box>

          <Box class="denial-item">
            <Flex class="denial-top" justify="between" align="center">
              <Badge variant="danger" size="sm">deny-ddl</Badge>
              <Text class="denial-time" size="sm" color="muted">09:12:33</Text>
            </Flex>
            <Box class="denial-sql">DROP TABLE user_entitlements CASCADE;</Box>
            <Box class="denial-agent">Agent: Cursor AI Agent</Box>
          </Box>

          <Box class="denial-item">
            <Flex class="denial-top" justify="between" align="center">
              <Badge variant="danger" size="sm">denied-table</Badge>
              <Text class="denial-time" size="sm" color="muted">08:04:19</Text>
            </Flex>
            <Box class="denial-sql">SELECT * FROM admin_credentials LIMIT 10;</Box>
            <Box class="denial-agent">Agent: Cursor AI Agent</Box>
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

  :global(.dashboard-view .view-header) {
    height: 52px;
    background: var(--bg-toolbar);
    border-bottom: 1px solid var(--border-default);
    padding: 0 16px;
  }

  :global(.dashboard-view .view-title) { font-weight: 600; font-size: var(--font-size-lg, 16px); }

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
