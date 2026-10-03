<script lang="ts">
  import { appState } from '../../state/appState.svelte';
  import { Icon } from '@ohjanus/icons';
</script>

<div class="log-console-container">
  <!-- Log Output content -->
  <div class="log-content code-text">
    <div class="log-line">
      <span class="log-ts">[2026-10-03 10:50:34]</span>
      <span class="log-text">Connected to dev_mh_asset</span>
    </div>
    <div class="log-line">
      <span class="log-ts">[2026-10-03 10:50:34]</span>
      <span class="log-schema">dev_mh_asset.public&gt;</span>
      <span class="kw-sql">SELECT</span> <span class="ident-sql">t.*</span>
    </div>
    <div class="log-line indent-sql">
      <span class="kw-sql">FROM</span> <span class="ident-sql">public.connection_credential t</span>
    </div>
    <div class="log-line indent-sql">
      <span class="kw-sql">LIMIT</span> <span class="num-sql">501</span>
    </div>
    <div class="log-line">
      <span class="log-ts">[2026-10-03 10:50:34]</span>
      <span class="log-text">58 rows retrieved starting from 1 in 586 ms (execution: 86 ms, fetching: 500 ms)</span>
    </div>
  </div>

  <!-- Right Action Strip matching design2.png -->
  <div class="log-action-strip">
    <button type="button" class="strip-btn" title="Toggle Output Window">
      <Icon name="audit" size={13} />
    </button>
    <button type="button" class="strip-btn" title="Collapse">
      <Icon name="minus" size={11} />
    </button>
    <button type="button" class="strip-btn" title="Undo">
      <Icon name="undo" size={12} />
    </button>
    <button type="button" class="strip-btn" title="Print Log">
      <Icon name="printer" size={13} />
    </button>
    <button type="button" class="strip-btn" title="Clear Console" onclick={() => appState.consoleLogs = []}>
      <Icon name="trash" size={13} />
    </button>
  </div>
</div>

<style>
  .log-console-container {
    width: 100%;
    height: 100%;
    display: flex;
    background-color: var(--bg-canvas, #1E1F22);
    overflow: hidden;
  }

  .log-content {
    flex: 1;
    overflow-y: auto;
    padding: 8px 12px;
    font-size: 12px;
    line-height: 18px;
    font-family: var(--font-code, 'JetBrains Mono', monospace);
  }

  .log-line {
    display: flex;
    gap: 8px;
  }

  .log-line.indent-sql {
    padding-left: 285px;
  }

  .log-ts {
    color: var(--text-muted, #7A7E85);
  }

  .log-text {
    color: var(--text-primary, #DFE1E5);
  }

  .log-schema {
    color: var(--text-primary, #DFE1E5);
    font-weight: 500;
  }

  .kw-sql {
    color: #CF8E6D;
    font-weight: 500;
  }

  .ident-sql {
    color: var(--text-primary, #DFE1E5);
  }

  .num-sql {
    color: #6897BB;
  }

  .log-action-strip {
    width: 32px;
    background-color: var(--bg-canvas, #1E1F22);
    border-left: 1px solid var(--border-subtle, #323438);
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 4px 0;
    gap: 4px;
    flex-shrink: 0;
  }

  .strip-btn {
    width: var(--icon-btn-size-sm, 26px);
    height: var(--icon-btn-size-sm, 26px);
    border-radius: var(--radius-sm, 4px);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: var(--font-size-xs, 11px);
    color: var(--text-muted, #7A7E85);
  }

  :global(.strip-btn svg) {
    width: var(--icon-size-sm, 14px);
    height: var(--icon-size-sm, 14px);
  }

  .strip-btn:hover {
    background-color: var(--bg-hover, #2B2D30);
    color: var(--text-primary, #DFE1E5);
  }
</style>
