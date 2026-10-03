<script lang="ts">
  import { appStore } from '../../appStore.svelte';
  import type { ConnectionInfo } from '../../types';

  let testingConn = $state<string | null>(null);

  function runTest(name: string) {
    testingConn = name;
    setTimeout(() => {
      appStore.testConnection(name);
      testingConn = null;
    }, 400);
  }
</script>

<div class="connections-view">
  <!-- Header -->
  <div class="view-header">
    <div>
      <span class="header-icon">🔌</span>
      <span class="view-title">Database Connections &amp; Connection Pools</span>
      <span class="view-desc">Zero credential leakage: backend resolves credentials from OS Keychain / Vault</span>
    </div>
  </div>

  <!-- Cards Grid -->
  <div class="cards-grid">
    {#each appStore.connections as conn}
      <div class="conn-card">
        <div class="card-header">
          <div class="header-left">
            <span class="db-icon">🐘</span>
            <span class="conn-title">{conn.name}</span>
            {#if conn.readonly}
              <span class="ro-badge">🔒 READ-ONLY</span>
            {:else}
              <span class="rw-badge">⚡ READ-WRITE</span>
            {/if}
          </div>
          <span class="status-badge {conn.status}">● {conn.status.toUpperCase()}</span>
        </div>

        <div class="conn-details">
          <div class="detail-row">
            <span class="label">Host &amp; Port:</span>
            <span class="value font-mono">{conn.host}</span>
          </div>
          <div class="detail-row">
            <span class="label">Engine Version:</span>
            <span class="value">{conn.version}</span>
          </div>
          <div class="detail-row">
            <span class="label">Last Health Ping:</span>
            <span class="value">{conn.lastPingAt} ({conn.latencyMs} ms)</span>
          </div>
        </div>

        <!-- Pool Stats -->
        <div class="pool-section">
          <div class="pool-title">CONNECTION POOL METRICS</div>
          <div class="pool-stats font-mono">
            <div class="stat-box">
              <span class="num">{conn.pool.inUse}</span>
              <span class="lbl">In-Use</span>
            </div>
            <div class="stat-box">
              <span class="num">{conn.pool.idle}</span>
              <span class="lbl">Idle</span>
            </div>
            <div class="stat-box">
              <span class="num">{conn.pool.open}</span>
              <span class="lbl">Open</span>
            </div>
            <div class="stat-box">
              <span class="num">{conn.pool.max}</span>
              <span class="lbl">Max Pool</span>
            </div>
          </div>
        </div>

        <!-- Security Guardrails -->
        <div class="guardrails-section">
          <div class="sec-label">ALLOWED SCHEMAS:</div>
          <div class="tags-list">
            {#each conn.allowedSchemas as sch}
              <span class="tag-pill">{sch}</span>
            {/each}
          </div>

          {#if conn.deniedTables.length > 0}
            <div class="sec-label" style="margin-top: 6px; color: var(--action-danger);">DENIED TABLES:</div>
            <div class="tags-list">
              {#each conn.deniedTables as dt}
                <span class="tag-pill denied">{dt}</span>
              {/each}
            </div>
          {/if}
        </div>

        <!-- Actions -->
        <div class="card-footer">
          <button 
            class="btn-secondary test-btn" 
            class:testing={testingConn === conn.name}
            onclick={() => runTest(conn.name)}
          >
            {testingConn === conn.name ? 'Pinging...' : '⟳ Test Connection'}
          </button>
        </div>
      </div>
    {/each}
  </div>
</div>

<style>
  .connections-view {
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

  .header-icon { font-size: 16px; margin-right: 8px; }
  .view-title { font-weight: 600; font-size: 13px; margin-right: 12px; }
  .view-desc { font-size: 11px; color: var(--text-muted); }

  .cards-grid {
    padding: 16px;
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
    gap: 16px;
  }

  .conn-card {
    background: var(--surface-sidebar);
    border: 1px solid var(--border-default);
    border-radius: 6px;
    padding: 14px;
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding-bottom: 8px;
    border-bottom: 1px solid var(--border-default);
  }

  .header-left {
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .conn-title {
    font-weight: 600;
    font-size: 13px;
  }

  .ro-badge {
    background: rgba(86, 168, 245, 0.15);
    color: #56A8F5;
    font-size: 9px;
    font-weight: 700;
    padding: 1px 4px;
    border-radius: 3px;
  }

  .rw-badge {
    background: rgba(237, 162, 0, 0.15);
    color: var(--action-warning);
    font-size: 9px;
    font-weight: 700;
    padding: 1px 4px;
    border-radius: 3px;
  }

  .status-badge.healthy {
    color: var(--action-success);
    font-size: 10px;
    font-weight: 700;
  }

  .conn-details {
    display: flex;
    flex-direction: column;
    gap: 4px;
    font-size: 11px;
  }

  .detail-row {
    display: flex;
    justify-content: space-between;
  }

  .label { color: var(--text-muted); }
  .value { color: var(--text-primary); }

  .pool-section {
    background: #18191B;
    border: 1px solid var(--border-default);
    border-radius: 4px;
    padding: 8px 10px;
  }

  .pool-title {
    font-size: 9px;
    font-weight: 700;
    color: var(--text-muted);
    letter-spacing: 0.5px;
    margin-bottom: 6px;
  }

  .pool-stats {
    display: flex;
    justify-content: space-between;
  }

  .stat-box {
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .stat-box .num {
    font-size: 14px;
    font-weight: 600;
    color: var(--text-primary);
  }

  .stat-box .lbl {
    font-size: 9px;
    color: var(--text-muted);
  }

  .guardrails-section {
    font-size: 10px;
  }

  .sec-label {
    font-weight: 700;
    color: var(--text-muted);
    margin-bottom: 4px;
  }

  .tags-list {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
  }

  .tag-pill {
    background: #25272A;
    border: 1px solid var(--border-default);
    padding: 1px 6px;
    border-radius: 3px;
    font-size: 10px;
  }

  .tag-pill.denied {
    border-color: rgba(229, 83, 83, 0.4);
    color: #FFB3B3;
  }

  .card-footer {
    display: flex;
    justify-content: flex-end;
    padding-top: 8px;
    border-top: 1px solid var(--border-default);
  }

  .test-btn {
    font-size: 11px;
  }

  .test-btn.testing {
    opacity: 0.7;
  }
</style>
