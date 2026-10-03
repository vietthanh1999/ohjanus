<script lang="ts">
  import { appStore } from '../../appStore.svelte';
</script>

<div class="dashboard-view">
  <div class="view-header">
    <div style="display:flex; align-items:center; gap:8px;">
      <span style="font-size: 16px;">📊</span>
      <span class="view-title">OhJanus MCP Gateway Observability &amp; Metrics</span>
    </div>
  </div>

  <div class="dashboard-body">
    <!-- Top KPI Cards -->
    <div class="kpi-grid">
      <div class="kpi-card">
        <span class="kpi-label">24H TOTAL QUERIES</span>
        <span class="kpi-val font-mono">14,892</span>
        <span class="kpi-trend pos">↑ +12.4% vs yesterday</span>
      </div>

      <div class="kpi-card">
        <span class="kpi-label">24H POLICY DENIALS</span>
        <span class="kpi-val font-mono warn">47</span>
        <span class="kpi-sub">Banned functions &amp; DDL blocked</span>
      </div>

      <div class="kpi-card">
        <span class="kpi-label">PENDING APPROVALS</span>
        <span class="kpi-val font-mono pending">{appStore.pendingApprovalsCount}</span>
        <span class="kpi-sub">Human review required</span>
      </div>

      <div class="kpi-card">
        <span class="kpi-label">P95 QUERY LATENCY</span>
        <span class="kpi-val font-mono">48 ms</span>
        <span class="kpi-sub">Gateway overhead &lt; 2 ms</span>
      </div>

      <div class="kpi-card">
        <span class="kpi-label">ACTIVE MCP TOKENS</span>
        <span class="kpi-val font-mono">{appStore.tokens.filter(t => t.state === 'active').length}</span>
        <span class="kpi-sub">Cursor &amp; Claude Desktop</span>
      </div>

      <div class="kpi-card">
        <span class="kpi-label">HEALTHY POOLS</span>
        <span class="kpi-val font-mono pos">3 / 3</span>
        <span class="kpi-sub">All connections healthy</span>
      </div>
    </div>

    <!-- Middle Split: Chart Simulation & Denials List -->
    <div class="mid-split">
      <div class="chart-panel">
        <div class="panel-header">
          <span>REAL-TIME QUERY THROUGHPUT (REQS / SEC)</span>
          <span class="live-pill">● LIVE</span>
        </div>
        <div class="chart-bars font-mono">
          {#each [18, 25, 42, 38, 55, 62, 45, 80, 72, 95, 88, 64, 52, 47, 63, 78, 92, 105, 84, 76] as h, i}
            <div class="bar-col">
              <div class="bar-fill" style="height: {h}%;"></div>
              <span class="bar-label">{i * 3}m</span>
            </div>
          {/each}
        </div>
      </div>

      <div class="security-panel">
        <div class="panel-header">
          <span>RECENT BLOCKED OPERATIONS (AST VALIDATION)</span>
        </div>
        <div class="denials-list font-mono">
          <div class="denial-item">
            <div class="denial-top">
              <span class="rule-tag">banned-function</span>
              <span class="denial-time">10:45:12</span>
            </div>
            <div class="denial-sql">SELECT pg_read_file('/etc/passwd')</div>
            <div class="denial-agent">Agent: Claude Desktop v1.2</div>
          </div>

          <div class="denial-item">
            <div class="denial-top">
              <span class="rule-tag">deny-ddl</span>
              <span class="denial-time">09:12:33</span>
            </div>
            <div class="denial-sql">DROP TABLE user_entitlements CASCADE;</div>
            <div class="denial-agent">Agent: Cursor AI Agent</div>
          </div>

          <div class="denial-item">
            <div class="denial-top">
              <span class="rule-tag">denied-table</span>
              <span class="denial-time">08:04:19</span>
            </div>
            <div class="denial-sql">SELECT * FROM admin_credentials LIMIT 10;</div>
            <div class="denial-agent">Agent: Cursor AI Agent</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  .dashboard-view {
    display: flex;
    flex-direction: column;
    height: 100%;
    background: var(--surface-canvas);
    overflow-y: auto;
  }

  .view-header {
    height: 48px;
    background: var(--surface-toolbar);
    border-bottom: 1px solid var(--border-default);
    display: flex;
    align-items: center;
    padding: 0 16px;
  }

  .view-title { font-weight: 600; font-size: 13px; }

  .dashboard-body {
    padding: 16px;
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .kpi-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
    gap: 12px;
  }

  .kpi-card {
    background: var(--surface-sidebar);
    border: 1px solid var(--border-default);
    border-radius: 5px;
    padding: 12px;
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .kpi-label {
    font-size: 10px;
    font-weight: 700;
    color: var(--text-muted);
    letter-spacing: 0.5px;
  }

  .kpi-val {
    font-size: 22px;
    font-weight: 700;
    color: var(--text-primary);
  }

  .kpi-val.pos { color: var(--action-success); }
  .kpi-val.warn { color: var(--action-danger); }
  .kpi-val.pending { color: var(--action-warning); }

  .kpi-trend.pos {
    font-size: 10px;
    color: var(--action-success);
  }

  .kpi-sub {
    font-size: 10px;
    color: var(--text-muted);
  }

  .mid-split {
    display: grid;
    grid-template-columns: 2fr 1fr;
    gap: 16px;
  }

  .chart-panel, .security-panel {
    background: var(--surface-sidebar);
    border: 1px solid var(--border-default);
    border-radius: 6px;
    padding: 14px;
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .panel-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 11px;
    font-weight: 700;
    color: var(--text-muted);
  }

  .live-pill {
    color: var(--action-success);
    font-size: 10px;
  }

  .chart-bars {
    height: 160px;
    display: flex;
    align-items: flex-end;
    gap: 6px;
    padding-top: 10px;
  }

  .bar-col {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    height: 100%;
    justify-content: flex-end;
  }

  .bar-fill {
    width: 100%;
    background: var(--action-primary);
    border-radius: 2px 2px 0 0;
    transition: height 0.3s ease;
  }

  .bar-fill:hover {
    background: var(--action-primary-hover);
  }

  .bar-label {
    font-size: 9px;
    color: var(--text-muted);
    margin-top: 4px;
  }

  .denials-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .denial-item {
    background: #18191B;
    border: 1px solid rgba(229, 83, 83, 0.3);
    border-radius: 4px;
    padding: 8px;
    font-size: 11px;
  }

  .denial-top {
    display: flex;
    justify-content: space-between;
    margin-bottom: 4px;
  }

  .rule-tag {
    background: rgba(229, 83, 83, 0.2);
    color: var(--action-danger);
    font-size: 9px;
    font-weight: 700;
    padding: 1px 4px;
    border-radius: 2px;
  }

  .denial-time {
    color: var(--text-muted);
    font-size: 10px;
  }

  .denial-sql {
    color: var(--text-primary);
    margin-bottom: 4px;
  }

  .denial-agent {
    color: var(--text-muted);
    font-size: 10px;
  }
</style>
