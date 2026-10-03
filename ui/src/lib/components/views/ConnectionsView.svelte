<script lang="ts">
  import { appStore } from '../../appStore.svelte';
  import type { ConnectionInfo } from '../../types';
  import { Card, Badge, Button, toast, Box, Flex, Grid, Stack } from '@ohjanus/ui';

  let testingConn = $state<string | null>(null);

  function runTest(name: string) {
    testingConn = name;
    setTimeout(() => {
      appStore.testConnection(name);
      testingConn = null;
      toast.success('Connection Ping Healthy', `Successfully pinged database "${name}".`);
    }, 400);
  }
</script>

<Box class="connections-view">
  <!-- Header -->
  <Flex as="header" class="view-header" align="center">
    <Flex align="center">
      <span class="header-icon">🔌</span>
      <span class="view-title">Database Connections &amp; Connection Pools</span>
      <span class="view-desc">Zero credential leakage: backend resolves credentials from OS Keychain / Vault</span>
    </Flex>
  </Flex>

  <!-- Cards Grid -->
  <Grid class="cards-grid" columns="repeat(auto-fill, minmax(340px, 1fr))" gap="16px">
    {#each appStore.connections as conn}
      <Card class="conn-card-item">
        {#snippet header()}
          <Flex class="card-header-inner" justify="between" align="center">
            <Flex class="header-left" align="center" gap="8px">
              <span class="db-icon">🐘</span>
              <span class="conn-title">{conn.name}</span>
              {#if conn.readonly}
                <Badge variant="warning" size="sm">🔒 READ-ONLY</Badge>
              {:else}
                <Badge variant="info" size="sm">⚡ READ-WRITE</Badge>
              {/if}
            </Flex>
            {#if conn.status === 'healthy'}
              <Badge variant="success" size="sm">● HEALTHY</Badge>
            {:else}
              <Badge variant="danger" size="sm">● DEGRADED</Badge>
            {/if}
          </Flex>
        {/snippet}

        <Stack class="conn-details" gap="6px">
          <Flex class="detail-row" justify="between">
            <span class="label">Host &amp; Port:</span>
            <span class="value font-mono">{conn.host}</span>
          </Flex>
          <Flex class="detail-row" justify="between">
            <span class="label">Engine Version:</span>
            <span class="value">{conn.version}</span>
          </Flex>
          <Flex class="detail-row" justify="between">
            <span class="label">Last Health Ping:</span>
            <span class="value">{conn.lastPingAt} ({conn.latencyMs} ms)</span>
          </Flex>
        </Stack>

        <!-- Pool Stats -->
        <Box class="pool-section">
          <Box class="pool-title">CONNECTION POOL METRICS</Box>
          <Flex class="pool-stats font-mono" justify="between">
            <Box class="stat-box">
              <span class="num">{conn.pool.inUse}</span>
              <span class="lbl">In-Use</span>
            </Box>
            <Box class="stat-box">
              <span class="num">{conn.pool.idle}</span>
              <span class="lbl">Idle</span>
            </Box>
            <Box class="stat-box">
              <span class="num">{conn.pool.open}</span>
              <span class="lbl">Open</span>
            </Box>
            <Box class="stat-box">
              <span class="num">{conn.pool.max}</span>
              <span class="lbl">Max Pool</span>
            </Box>
          </Flex>
        </Box>

        <!-- Security Guardrails -->
        <Box class="guardrails-section">
          <Box class="sec-label">ALLOWED SCHEMAS:</Box>
          <Flex class="tags-list" wrap gap="4px">
            {#each conn.allowedSchemas as sch}
              <Badge variant="default" size="sm">{sch}</Badge>
            {/each}
          </Flex>

          {#if conn.deniedTables.length > 0}
            <Box class="sec-label" style="margin-top: 6px; color: var(--action-danger);">DENIED TABLES:</Box>
            <Flex class="tags-list" wrap gap="4px">
              {#each conn.deniedTables as dt}
                <Badge variant="danger" size="sm">{dt}</Badge>
              {/each}
            </Flex>
          {/if}
        </Box>

        {#snippet footer()}
          <Button
            variant="secondary"
            size="sm"
            onclick={() => runTest(conn.name)}
            loading={testingConn === conn.name}
          >
            {testingConn === conn.name ? 'Pinging...' : '⟳ Test Connection'}
          </Button>
        {/snippet}
      </Card>
    {/each}
  </Grid>
</Box>

<style>
  .connections-view {
    display: flex;
    flex-direction: column;
    height: 100%;
    background: var(--bg-canvas);
    overflow-y: auto;
  }

  .view-header {
    height: 48px;
    background: var(--bg-toolbar);
    border-bottom: 1px solid var(--border-default);
    padding: 0 16px;
  }

  .header-icon { font-size: 16px; margin-right: 8px; }
  .view-title { font-weight: 600; font-size: 13px; margin-right: 12px; }
  .view-desc { font-size: 11px; color: var(--text-muted); }

  .cards-grid {
    padding: 16px;
  }

  :global(.card-header-inner) {
    width: 100%;
  }

  .db-icon { font-size: 16px; }
  .conn-title { font-weight: 600; font-size: 13px; color: var(--text-primary); }

  :global(.conn-details) {
    margin-bottom: 14px;
    font-size: 11.5px;
  }

  :global(.detail-row .label) { color: var(--text-muted); }
  :global(.detail-row .value) { color: var(--text-primary); }

  :global(.pool-section) {
    background: rgba(0, 0, 0, 0.2);
    border: 1px solid var(--border-subtle);
    border-radius: 4px;
    padding: 10px;
    margin-bottom: 12px;
  }

  :global(.pool-title) {
    font-size: 10px;
    font-weight: 600;
    color: var(--text-muted);
    letter-spacing: 0.5px;
    margin-bottom: 8px;
  }

  :global(.pool-stats) {
    text-align: center;
  }

  :global(.stat-box .num) {
    display: block;
    font-size: 14px;
    font-weight: 700;
    color: var(--action-primary);
  }

  :global(.stat-box .lbl) {
    font-size: 9.5px;
    color: var(--text-muted);
  }

  :global(.guardrails-section) {
    margin-bottom: 8px;
  }

  :global(.sec-label) {
    font-size: 10px;
    font-weight: 600;
    color: var(--text-muted);
    margin-bottom: 4px;
  }
</style>
