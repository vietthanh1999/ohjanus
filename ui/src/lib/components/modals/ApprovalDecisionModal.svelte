<script lang="ts">
  import { appState } from '../../state/appState.svelte';
  import { Modal, Button, Alert, Field, toast, Box, Stack } from '@ohjanus/ui';

  function close() {
    appState.selectedApprovalForAction = null;
  }

  async function confirm() {
    const isApprove = appState.approvalDecisionMode === 'approve';
    const id = appState.selectedApprovalForAction?.id ?? '';
    try {
      await appState.confirmApprovalAction();
    } catch (e) {
      toast.error('Decision Failed', e instanceof Error ? e.message : String(e));
      return;
    }

    if (isApprove) {
      toast.success('Query Approved & Executed', `Statement ${id} successfully authorized.`);
    } else {
      toast.error('Query Rejected', `Access denied for statement ${id}.`);
    }
  }

  let isOpen = $derived(Boolean(appState.selectedApprovalForAction));
  let isApprove = $derived(appState.approvalDecisionMode === 'approve');
  let title = $derived(
    isApprove
      ? `Approve Operation: ${appState.selectedApprovalForAction?.id ?? ''}`
      : `Reject Operation: ${appState.selectedApprovalForAction?.id ?? ''}`
  );
</script>

{#if appState.selectedApprovalForAction}
  <Modal open={isOpen} onClose={close} {title} width="540px">
    {#snippet children()}
      <Stack class="modal-content-stack" gap="14px">
        <Alert variant={isApprove ? 'warning' : 'danger'}>
          {#if isApprove}
            <p><strong>Warning:</strong> You are about to authorize an agent write query to execute on <strong>{appState.selectedApprovalForAction?.connection}</strong>.</p>
          {:else}
            <p><strong>Note:</strong> Rejecting will return an access denial error with your explanation back to the requesting agent.</p>
          {/if}
        </Alert>

        <Field
          label={isApprove ? 'Approval Note (Optional):' : 'Rejection Reason:'}
          required={!isApprove}
        >
          <textarea
            id="decision-reason"
            rows="3"
            class="form-textarea code-text"
            placeholder={isApprove ? 'e.g. Verified change with DevOps lead' : 'e.g. Unbounded DELETE without date filter is prohibited'}
            bind:value={appState.approvalDecisionReason}
          ></textarea>
        </Field>

        <Box class="sql-box-preview code-text">
          {appState.selectedApprovalForAction?.sql}
        </Box>
      </Stack>
    {/snippet}

    {#snippet footer()}
      <Button variant="secondary" onclick={close}>Cancel</Button>
      {#if isApprove}
        <Button variant="primary" onclick={confirm}>
          Confirm & Execute Statement
        </Button>
      {:else}
        <Button variant="danger" onclick={confirm}>
          Confirm Rejection
        </Button>
      {/if}
    {/snippet}
  </Modal>
{/if}

<style>
  .form-textarea {
    width: 100%;
    background-color: var(--bg-canvas, #1E1F22);
    border: 1px solid var(--border-default, #393B40);
    border-radius: 4px;
    padding: 8px 10px;
    color: var(--text-primary, #DFE1E5);
    font-size: 12px;
    resize: vertical;
    outline: none;
  }

  .form-textarea:focus {
    border-color: var(--border-accent, #3574F0);
  }

  :global(.sql-box-preview) {
    background-color: var(--bg-canvas, #1E1F22);
    border: 1px solid var(--border-default, #393B40);
    padding: 10px 12px;
    border-radius: 4px;
    font-size: 11px;
    color: #DFE1E5;
    white-space: pre-wrap;
    max-height: 120px;
    overflow-y: auto;
  }
</style>
