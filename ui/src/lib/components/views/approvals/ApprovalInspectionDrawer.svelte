<script lang="ts">
  import type { ApprovalRequest } from "../../../state/appState.svelte";
  import { Box, Flex, Stack, Button, Text } from "@ohjanus/ui";
  import { Icon } from "@ohjanus/icons";

  interface Props {
    item: ApprovalRequest | null;
    onclose: () => void;
    onaction?: (mode: "approve" | "reject") => void;
  }

  let { item, onclose, onaction }: Props = $props();
</script>

<Box class="detail-drawer">
  {#if item}
    <Flex class="drawer-header" align="center" justify="between">
      <Text weight="semibold" size="sm" class="drawer-title">Request Inspection: {item.id}</Text>
      <Button
        variant="ghost"
        size="icon-sm"
        class="jb-icon-btn"
        title="Close Inspector"
        onclick={onclose}
      >
        <Icon name="x" size={12} />
      </Button>
    </Flex>
    <Stack class="drawer-body" gap="8px">
      <Box class="meta-section">
        <Flex class="meta-row" align="center" justify="between">
          <Text size="xs" color="muted">Connection:</Text>
          <Text size="xs">{item.connection} (PostgreSQL 16)</Text>
        </Flex>
        <Flex class="meta-row" align="center" justify="between">
          <Text size="xs" color="muted">Requesting Client:</Text>
          <Text size="xs">{item.requested_by.client}</Text>
        </Flex>
        <Flex class="meta-row" align="center" justify="between">
          <Text size="xs" color="muted">Agent Token ID:</Text>
          <Text size="xs" mono>{item.requested_by.token_id}</Text>
        </Flex>
        <Flex class="meta-row" align="center" justify="between">
          <Text size="xs" color="muted">Submitted At:</Text>
          <Text size="xs" mono>{item.created_at}</Text>
        </Flex>
        <Flex class="meta-row" align="center" justify="between">
          <Text size="xs" color="muted">Expires At:</Text>
          <Text size="xs" mono>{item.expires_at}</Text>
        </Flex>
      </Box>

      <Text size="xs" weight="semibold" color="secondary" class="section-title">Full SQL Statement</Text>
      <Box class="full-sql code-text">{item.sql}</Box>

      <Text size="xs" weight="semibold" color="secondary" class="section-title">Execution Safety Assessment</Text>
      <Stack class="safety-box" gap="6px">
        <Flex class="safety-item" align="center" gap="8px">
          <Box
            as="span"
            class={`safe-dot ${item.affected_estimate > 100 ? "risk" : ""}`}
          />
          <Text size="xs">Blast Radius: {item.affected_estimate} row(s) estimated</Text>
        </Flex>
        <Flex class="safety-item" align="center" gap="8px">
          <Box
            as="span"
            class={`safe-dot ${item.statement_type === "DROP" ? "risk" : ""}`}
          />
          <Text size="xs">Statement Type: {item.statement_type}</Text>
        </Flex>
        <Flex class="safety-item" align="center" gap="8px">
          <Box as="span" class="safe-dot" />
          <Text size="xs">AST Validator: Parsed successfully via pg_query_go</Text>
        </Flex>
      </Stack>

      {#if item.state === "pending"}
        <Flex class="drawer-actions" gap="8px">
          <Button
            variant="danger"
            style="flex: 1;"
            onclick={() => onaction?.("reject")}
          >
            Reject Request
          </Button>
          <Button
            variant="primary"
            style="flex: 1;"
            onclick={() => onaction?.("approve")}
          >
            Approve Request
          </Button>
        </Flex>
      {/if}
    </Stack>
  {:else}
    <Flex class="drawer-empty" align="center" justify="center">
      <Text size="sm" color="muted">
        Select an approval request to inspect full AST validation and
        EXPLAIN execution plan
      </Text>
    </Flex>
  {/if}
</Box>

<style>
  :global(.detail-drawer) {
    width: 360px;
    background-color: var(--bg-surface);
    border-left: 1px solid var(--border-default);
    display: flex;
    flex-direction: column;
    flex-shrink: 0;
  }

  :global(.detail-drawer .drawer-header) {
    height: 36px;
    padding: 0 12px;
    border-bottom: 1px solid var(--border-default);
  }

  :global(.drawer-body) {
    padding: 12px;
    overflow-y: auto;
    flex: 1;
  }

  :global(.detail-drawer .meta-section) {
    display: flex;
    flex-direction: column;
    gap: 6px;
    background-color: var(--bg-canvas);
    border: 1px solid var(--border-subtle);
    border-radius: 4px;
    padding: 8px 10px;
  }

  :global(.detail-drawer .section-title) {
    text-transform: uppercase;
    letter-spacing: 0.5px;
    margin-top: 8px;
  }

  :global(.detail-drawer .full-sql) {
    background-color: var(--bg-canvas);
    border: 1px solid var(--border-subtle);
    border-radius: 4px;
    padding: 10px;
    font-size: var(--font-size-xs, 12px);
    color: var(--text-primary);
    white-space: pre-wrap;
    word-break: break-all;
    max-height: 180px;
    overflow-y: auto;
    margin: 0;
  }

  :global(.detail-drawer .safety-box) {
    background-color: var(--bg-canvas);
    border: 1px solid var(--border-subtle);
    border-radius: 4px;
    padding: 8px 10px;
  }

  :global(.detail-drawer .safe-dot) {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background-color: var(--action-success, #22c55e);
    flex-shrink: 0;
    display: inline-block;
  }

  :global(.detail-drawer .safe-dot.risk) {
    background-color: var(--action-danger, #ef4444);
  }

  :global(.detail-drawer .drawer-actions) {
    margin-top: 12px;
  }

  :global(.detail-drawer .drawer-empty) {
    height: 100%;
    padding: 24px;
    text-align: center;
  }
</style>
