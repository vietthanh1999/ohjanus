<script lang="ts">
  import { appState } from '../../state/appState.svelte';
  import { Modal, Button, Input, Field, Alert, toast, Box, Flex, Stack, Text } from '@ohjanus/ui';
  import { Icon } from '@ohjanus/icons';

  let tokenName = $state('');
  let selectedScopes = $state<string[]>(['read', 'schema']);
  let hasCopied = $state(false);

  const availableScopes = [
    { id: 'read', label: 'read (db_read, SELECT)', desc: 'Allow read-only queries' },
    { id: 'write_with_approval', label: 'write_with_approval', desc: 'Allows write queries after human approval' },
    { id: 'schema', label: 'schema (db_schema, introspection)', desc: 'Inspect table definitions and columns' },
    { id: 'explain', label: 'explain (db_explain)', desc: 'Run EXPLAIN on query plans' },
    { id: 'admin', label: 'admin (full gateway config)', desc: 'Full administration scopes' }
  ];

  function toggleScope(id: string) {
    if (selectedScopes.includes(id)) {
      selectedScopes = selectedScopes.filter(s => s !== id);
    } else {
      selectedScopes = [...selectedScopes, id];
    }
  }

  let creating = $state(false);

  async function handleCreate() {
    if (!tokenName.trim() || creating) return;
    creating = true;
    try {
      await appState.createToken(tokenName.trim(), selectedScopes, 90);
      toast.success('MCP Token Generated', `Token for "${tokenName.trim()}" created successfully.`);
    } catch (e) {
      toast.error('Token Creation Failed', e instanceof Error ? e.message : String(e));
    } finally {
      creating = false;
    }
  }

  function copySecret() {
    if (appState.createdTokenSecret) {
      navigator.clipboard.writeText(appState.createdTokenSecret);
      hasCopied = true;
      toast.info('Copied to Clipboard', 'Token bearer secret copied.');
      setTimeout(() => hasCopied = false, 2000);
    }
  }

  function close() {
    appState.createTokenModalOpen = false;
    appState.createdTokenSecret = null;
    tokenName = '';
    selectedScopes = ['read', 'schema'];
  }
</script>

{#if appState.createTokenModalOpen}
  <Modal open={appState.createTokenModalOpen} onClose={close} title="Create MCP Bearer Access Token" width="500px">
    {#snippet children()}
      <Stack class="modal-body-stack" gap="14px">
        {#if !appState.createdTokenSecret}
          <Field label="Client Name / Description" required>
            <Input
              placeholder="e.g. Claude Desktop Agent - Finance Sync"
              bind:value={tokenName}
            />
          </Field>

          <Field label="Allowed Janus Scopes (§1.2)">
            <Stack class="scopes-list" gap="6px">
              {#each availableScopes as sc}
                <!-- svelte-ignore a11y_click_events_have_key_events -->
                <!-- svelte-ignore a11y_no_static_element_interactions -->
                <Flex class="scope-row" align="start" gap="8px" onclick={() => toggleScope(sc.id)}>
                  <input
                    type="checkbox"
                    checked={selectedScopes.includes(sc.id)}
                    class="scope-checkbox"
                  />
                  <Stack class="scope-text" gap="2px">
                    <span class="sc-name code-text">{sc.label}</span>
                    <Text size="sm" color="muted">{sc.desc}</Text>
                  </Stack>
                </Flex>
              {/each}
            </Stack>
          </Field>
        {:else}
          <!-- Success Token Reveal View -->
          <Stack class="token-reveal-view" gap="12px">
            <Alert variant="warning" title="Save this token key now!">
              <p style="font-size: 11px; margin-top: 2px;">
                Janus stores only cryptographic hashes in its vault. You will not be able to view this plaintext secret again.
              </p>
            </Alert>

            <Box class="secret-box code-text">
              <span>{appState.createdTokenSecret}</span>
            </Box>

            <Button variant="secondary" onclick={copySecret}>
              {#if hasCopied}
                <Icon name="check" size={14} color="#57D38C" />
                <span>Copied to Clipboard!</span>
              {:else}
                <Icon name="copy" size={14} />
                <span>Copy Token to Clipboard</span>
              {/if}
            </Button>
          </Stack>
        {/if}
      </Stack>
    {/snippet}

    {#snippet footer()}
      {#if !appState.createdTokenSecret}
        <Button variant="secondary" onclick={close}>Cancel</Button>
        <Button variant="primary" onclick={handleCreate} disabled={!tokenName.trim() || creating}>
          {creating ? 'Generating...' : 'Generate Token'}
        </Button>
      {:else}
        <Button variant="primary" onclick={close}>
          Done
        </Button>
      {/if}
    {/snippet}
  </Modal>
{/if}

<style>
  :global(.scopes-list) {
    background-color: var(--bg-canvas, #1E1F22);
    border: 1px solid var(--border-default, #393B40);
    border-radius: 4px;
    padding: 8px;
    max-height: 200px;
    overflow-y: auto;
  }

  :global(.scope-row) {
    padding: 6px 8px;
    border-radius: 3px;
    cursor: pointer;
    transition: background-color 0.1s ease;
  }

  :global(.scope-row:hover) {
    background-color: var(--bg-hover, #313438);
  }

  .scope-checkbox {
    margin-top: 2px;
    cursor: pointer;
  }

  .sc-name {
    font-size: 11.5px;
    color: var(--text-primary, #DFE1E5);
    font-weight: 500;
  }

  :global(.secret-box) {
    background-color: var(--bg-canvas, #1E1F22);
    border: 1px dashed var(--action-warning, #EDA200);
    padding: 12px;
    border-radius: 4px;
    word-break: break-all;
    font-size: 12px;
    color: #FACC15;
    user-select: all;
  }
</style>
