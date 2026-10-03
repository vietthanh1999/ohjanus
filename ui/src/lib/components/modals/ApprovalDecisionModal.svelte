<script lang="ts">
  import { appState } from '../../state/appState.svelte';

  function close() {
    appState.selectedApprovalForAction = null;
  }

  function confirm() {
    appState.confirmApprovalAction();
  }
</script>

{#if appState.selectedApprovalForAction}
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div class="jb-modal-backdrop" onclick={close}>
    <div class="jb-modal" style="width: 520px;" onclick={(e) => e.stopPropagation()}>
      <div class="jb-modal-header">
        <span>
          {#if appState.approvalDecisionMode === 'approve'}
            <span style="color: var(--action-success);">Approve Operation</span>: {appState.selectedApprovalForAction.id}
          {:else}
            <span style="color: var(--action-danger);">Reject Operation</span>: {appState.selectedApprovalForAction.id}
          {/if}
        </span>
        <button class="jb-icon-btn" onclick={close}>✕</button>
      </div>

      <div class="jb-modal-body">
        <div class="summary-banner" class:danger={appState.approvalDecisionMode === 'reject'}>
          {#if appState.approvalDecisionMode === 'approve'}
            <p><strong>Warning:</strong> You are about to authorize an agent write query to execute on <strong>{appState.selectedApprovalForAction.connection}</strong>.</p>
          {:else}
            <p><strong>Note:</strong> Rejecting will return an access denial error with your explanation back to the requesting agent.</p>
          {/if}
        </div>

        <div class="form-group">
          <label for="decision-reason" class="form-label">
            {#if appState.approvalDecisionMode === 'approve'}
              Approval Note (Optional):
            {:else}
              Rejection Reason <span style="color: var(--action-danger);">*</span>:
            {/if}
          </label>
          <textarea
            id="decision-reason"
            rows="3"
            class="form-textarea code-text"
            placeholder={appState.approvalDecisionMode === 'approve' ? 'e.g. Verified change with DevOps lead' : 'e.g. Unbounded DELETE without date filter is prohibited'}
            bind:value={appState.approvalDecisionReason}
          ></textarea>
        </div>

        <div class="sql-box-preview code-text">
          {appState.selectedApprovalForAction.sql}
        </div>
      </div>

      <div class="jb-modal-footer">
        <button class="jb-btn-secondary" onclick={close}>Cancel</button>
        {#if appState.approvalDecisionMode === 'approve'}
          <button class="jb-btn-primary" style="background-color: var(--action-success); color: #14281B; font-weight: 600;" onclick={confirm}>
            Confirm & Execute Statement
          </button>
        {:else}
          <button class="jb-btn-danger" onclick={confirm}>
            Confirm Rejection
          </button>
        {/if}
      </div>
    </div>
  </div>
{/if}

<style>
  .summary-banner {
    background-color: rgba(87, 211, 140, 0.1);
    border: 1px solid rgba(87, 211, 140, 0.3);
    padding: 10px;
    border-radius: 4px;
    font-size: 11px;
    color: var(--text-primary);
    margin-bottom: 12px;
  }

  .summary-banner.danger {
    background-color: rgba(229, 83, 83, 0.1);
    border-color: rgba(229, 83, 83, 0.3);
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 6px;
    margin-bottom: 12px;
  }

  .form-label {
    font-size: 11px;
    font-weight: 500;
    color: var(--text-secondary);
  }

  .form-textarea {
    background-color: #1E1F22;
    border: 1px solid var(--border-default);
    border-radius: 3px;
    padding: 6px 8px;
    color: var(--text-primary);
    font-size: 11px;
    resize: vertical;
  }

  .form-textarea:focus {
    border-color: var(--border-accent);
  }

  .sql-box-preview {
    background-color: #1E1F22;
    border: 1px solid var(--border-default);
    padding: 8px 10px;
    border-radius: 4px;
    font-size: 10px;
    color: #DFE1E5;
    white-space: pre-wrap;
    max-height: 100px;
    overflow-y: auto;
  }
</style>
