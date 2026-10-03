<script lang="ts">
  import { appState } from '../../state/appState.svelte';
  import { Card, Badge, Button, Alert, Text, toast, Box, Flex, Grid, Stack } from '@ohjanus/ui';
  import { Icon } from '@ohjanus/icons';

  let testingConn = $state<string | null>(null);

  async function runTest(name: string) {
    if (testingConn) return;
    testingConn = name;
    try {
      const res = await appState.pingConnection(name);
      toast.success('Connection Ping Healthy', `"${name}" responded in ${res.latency_ms} ms.`);
    } catch (e) {
      toast.error('Connection Test Failed', e instanceof Error ? e.message : String(e));
    } finally {
      testingConn = null;
    }
  }
</script>

<Box class="connections-view">
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
  <!-- Header -->
  <Flex as="header" class="view-header" align="center">
    <Flex align="center" gap="8px">
      <span class="header-icon"><Icon name="database" size={16} color="#3B82F6" /></span>
      <span class="view-title">Database Connections &amp; Connection Pools</span>
      <span class="view-desc">Zero credential leakage: backend resolves credentials from OS Keychain / Vault</span>
    </Flex>
  </Flex>

  <!-- Cards Grid -->
  <Grid class="cards-grid" columns="repeat(auto-fill, minmax(340px, 1fr))" gap="16px">
      {#each appState.connections as conn}
      <Card class="conn-card-item">
        {#snippet header()}
          <Flex class="card-header-inner" justify="between" align="center">
            <Flex class="header-left" align="center" gap="8px">
              <span class="db-icon"><Icon name="database" size={15} color="#3B82F6" /></span>
              <span class="conn-title">{conn.name}</span>
              {#if conn.readonly}
                <Badge variant="warning" size="sm">
                  <Icon name="lock" size={11} />
                  <span>READ-ONLY</span>
                </Badge>
              {:else}
                <Badge variant="info" size="sm">
                  <Icon name="lightning" size={11} />
                  <span>READ-WRITE</span>
                </Badge>
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
            <span class="label">Driver:</span>
            <span class="value font-mono">{conn.driver}</span>
          </Flex>
          <Flex class="detail-row" justify="between">
            <span class="label">Last Health Ping:</span>
            <span class="value">{conn.last_ping_at}{conn.latency_ms !== undefined ? ` (${conn.latency_ms} ms)` : ''}</span>
          </Flex>
        </Stack>

        <!-- Security Guardrails -->
        <Box class="guardrails-section">
          <Box class="sec-label">ALLOWED SCHEMAS:</Box>
          <Flex class="tags-list" wrap gap="4px">
            {#each conn.allowed_schemas as sch}
              <Badge variant="default" size="sm">{sch}</Badge>
            {/each}
          </Flex>

          {#if conn.denied_tables.length > 0}
            <Box class="sec-label" style="margin-top: 6px; color: var(--action-danger);">DENIED TABLES:</Box>
            <Flex class="tags-list" wrap gap="4px">
              {#each conn.denied_tables as dt}
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
  :global(.connections-view) {
    display: flex;
    flex-direction: column;
    height: 100%;
    background: var(--bg-canvas);
    overflow-y: auto;
  }

  :global(.connections-view .view-header) {
    height: 48px;
    background: var(--bg-toolbar);
    border-bottom: 1px solid var(--border-default);
    padding: 0 16px;
  }

  :global(.connections-view .header-icon) { font-size: 16px; margin-right: 8px; }
  :global(.connections-view .view-title) { font-weight: 600; font-size: 13px; margin-right: 12px; }
  :global(.connections-view .view-desc) { font-size: 11px; color: var(--text-muted); }

  :global(.connections-view .cards-grid) {
    padding: 16px;
  }

  :global(.card-header-inner) {
    width: 100%;
  }

  :global(.connections-view .db-icon) { font-size: 16px; }
  :global(.connections-view .conn-title) { font-weight: 600; font-size: 13px; color: var(--text-primary); }

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
