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
      <span class="view-title">OhJanus MCP Gateway Observability &amp; Metrics</span>
    </Flex>
  </Flex>

  <Stack class="dashboard-body" gap="16px">
    <!-- Top KPI Cards via @ohjanus/ui -->
    <Grid class="kpi-grid" columns="repeat(auto-fill, minmax(180px, 1fr))" gap="12px">
      <Card class="kpi-card-box">
        <span class="kpi-label">TOTAL REQUESTS (AUDIT)</span>
        <span class="kpi-val font-mono">{appState.summary ? appState.summary.requests_total.toLocaleString() : '…'}</span>
        <span class="kpi-trend pos">live from Admin API</span>
      </Card>

      <Card class="kpi-card-box">
        <span class="kpi-label">TOTAL POLICY DENIALS</span>
        <span class="kpi-val font-mono warn">{appState.summary ? appState.summary.denials_total.toLocaleString() : '…'}</span>
        <span class="kpi-sub">Banned functions &amp; DDL blocked</span>
      </Card>

      <Card class="kpi-card-box">
        <span class="kpi-label">PENDING APPROVALS</span>
        <span class="kpi-val font-mono pending">{appState.notificationCount}</span>
        <span class="kpi-sub">Human review required</span>
      </Card>

      <Card class="kpi-card-box">
        <span class="kpi-label">P95 QUERY LATENCY</span>
        <span class="kpi-val font-mono">48 ms</span>
        <span class="kpi-sub">Gateway overhead &lt; 2 ms</span>
      </Card>

      <Card class="kpi-card-box">
        <span class="kpi-label">ACTIVE MCP TOKENS</span>
        <span class="kpi-val font-mono">{appState.tokens.filter(t => t.state === 'active').length}</span>
        <span class="kpi-sub">Cursor &amp; Claude Desktop</span>
      </Card>

      <Card class="kpi-card-box">
        <span class="kpi-label">HEALTHY POOLS</span>
        <span class="kpi-val font-mono pos">3 / 3</span>
        <span class="kpi-sub">All connections healthy</span>
      </Card>
    </Grid>

    <!-- Middle Split: Chart Simulation & Denials List -->
    <Grid class="mid-split" columns="3fr 2fr" gap="16px">
      <Card class="chart-panel">
        {#snippet header()}
          <Flex class="panel-header-inner" justify="between" align="center">
            <span>REAL-TIME QUERY THROUGHPUT (REQS / SEC)</span>
            <Badge variant="success" size="sm">● LIVE</Badge>
          </Flex>
        {/snippet}
        <Flex class="chart-bars font-mono" align="end" gap="6px">
          {#each [18, 25, 42, 38, 55, 62, 45, 80, 72, 95, 88, 64, 52, 47, 63, 78, 92, 105, 84, 76] as h, i}
            <Flex class="bar-col" direction="column" align="center" justify="end">
              <Box class="bar-fill" style="height: {h}%;"></Box>
              <span class="bar-label">{i * 3}m</span>
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
              <span class="denial-time">10:45:12</span>
            </Flex>
            <Box class="denial-sql">SELECT pg_read_file('/etc/passwd')</Box>
            <Box class="denial-agent">Agent: Claude Desktop v1.2</Box>
          </Box>

          <Box class="denial-item">
            <Flex class="denial-top" justify="between" align="center">
              <Badge variant="danger" size="sm">deny-ddl</Badge>
              <span class="denial-time">09:12:33</span>
            </Flex>
            <Box class="denial-sql">DROP TABLE user_entitlements CASCADE;</Box>
            <Box class="denial-agent">Agent: Cursor AI Agent</Box>
          </Box>

          <Box class="denial-item">
            <Flex class="denial-top" justify="between" align="center">
              <Badge variant="danger" size="sm">denied-table</Badge>
              <span class="denial-time">08:04:19</span>
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
    height: 48px;
    background: var(--bg-toolbar);
    border-bottom: 1px solid var(--border-default);
    padding: 0 16px;
  }

  :global(.dashboard-view .view-title) { font-weight: 600; font-size: 13px; }

  :global(.dashboard-view .dashboard-body) {
    padding: 16px;
  }

  :global(.kpi-card-box .ohjanus-card-body) {
    display: flex;
    flex-direction: column;
    gap: 4px;
    padding: 12px 14px;
  }

  :global(.dashboard-view .kpi-label) { font-size: 10px; font-weight: 600; color: var(--text-muted); letter-spacing: 0.5px; }
  :global(.dashboard-view .kpi-val) { font-size: 22px; font-weight: 700; color: var(--text-primary); }
  :global(.dashboard-view .kpi-val.pos) { color: var(--action-success); }
  :global(.dashboard-view .kpi-val.warn) { color: var(--action-danger); }
  :global(.dashboard-view .kpi-val.pending) { color: var(--action-warning); }
  :global(.dashboard-view .kpi-trend.pos) { font-size: 10.5px; color: var(--action-success); font-weight: 500; }
  :global(.dashboard-view .kpi-sub) { font-size: 10.5px; color: var(--text-muted); }

  :global(.dashboard-view .panel-header-inner) {
    width: 100%;
    font-size: 11px;
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
    font-size: 9px;
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

  :global(.dashboard-view .denial-time) { font-size: 10.5px; color: var(--text-muted); }
  :global(.denial-sql) { font-size: 11px; color: var(--action-danger); margin-bottom: 2px; }
  :global(.denial-agent) { font-size: 10px; color: var(--text-muted); }
</style>
