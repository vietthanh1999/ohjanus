<script lang="ts">
  import type { ApprovalRequest } from "../../../state/appState.svelte";
  import { Badge, Button, Alert, Box, Flex, Text } from "@ohjanus/ui";
  import { Icon } from "@ohjanus/icons";

  interface Props {
    item: ApprovalRequest;
    isSelected?: boolean;
    onselect?: () => void;
    onaction?: (mode: "approve" | "reject") => void;
  }

  let { item, isSelected = false, onselect, onaction }: Props = $props();

  function getStatementVariant(
    type: string,
  ): "default" | "success" | "warning" | "danger" | "info" {
    switch (type) {
      case "SELECT":
        return "info";
      case "INSERT":
        return "success";
      case "UPDATE":
        return "warning";
      case "DELETE":
      case "DROP":
        return "danger";
      default:
        return "default";
    }
  }
</script>

<Box
  class={`approval-card ${isSelected ? "selected" : ""}`}
  onclick={() => onselect?.()}
>
  <!-- Card Header -->
  <Flex class="card-top" align="center" justify="between">
    <Flex class="card-badges" align="center" gap="8px">
      <Badge
        variant={getStatementVariant(item.statement_type)}
        size="sm"
      >
        {item.statement_type}
      </Badge>
      <Text class="conn-pill" size="sm" color="secondary">{item.connection}</Text>
      {#if item.state === "pending"}
        <Badge variant="warning" size="sm">Pending Review</Badge>
      {:else if item.state === "approved"}
        <Badge variant="success" size="sm">Approved</Badge>
      {:else}
        <Badge variant="danger" size="sm">Rejected</Badge>
      {/if}
    </Flex>

    <Box class="time-meta">
      <Text size="xs" color="muted" class="expiry-time">Exp: {item.expires_at.substring(11, 16)}</Text>
    </Box>
  </Flex>

  <!-- Requester info -->
  <Flex class="requester-row" align="center" gap="12px">
    <Flex align="center" gap="4px" inline class="agent-client">
      <Icon name="user" size={12} />
      <Text size="sm">{item.requested_by.client}</Text>
    </Flex>
    <Text size="xs" mono class="token-tag">{item.requested_by.token_id}</Text>
    <Text class="rows-affected" size="xs" color="muted" weight="medium">
      Est. ~{item.affected_estimate.toLocaleString()} rows
    </Text>
  </Flex>

  <!-- Risk Warnings -->
  {#if item.warnings && item.warnings.length > 0}
    <Alert variant="warning">
      {#each item.warnings as warn}
        <Flex class="warn-line" align="center" gap="6px">
          <Box as="span" class="warn-icon"><Icon name="alert-triangle" size={12} /></Box>
          <Text size="xs">{warn}</Text>
        </Flex>
      {/each}
    </Alert>
  {/if}

  <!-- SQL Snippet -->
  <Box class="sql-preview code-text">
    {item.sql}
  </Box>

  <!-- Action buttons for Pending or Decision Meta -->
  {#if item.state === "pending"}
    <Flex class="card-actions" align="center" justify="end" gap="8px">
      <Button
        variant="danger"
        size="xs"
        onclick={(e: MouseEvent) => {
          e.stopPropagation();
          onaction?.("reject");
        }}
      >
        Reject...
      </Button>
      <Button
        variant="primary"
        size="xs"
        onclick={(e: MouseEvent) => {
          e.stopPropagation();
          onaction?.("approve");
        }}
      >
        Approve & Execute
      </Button>
    </Flex>
  {:else}
    <Box class="decision-meta">
      <Text size="xs" color="muted">
        Decided by <Text weight="bold">{item.decided_by}</Text> on {item.decided_at}
      </Text>
      {#if item.decision_reason}
        <Text size="xs" color="secondary" class="reason-note">Reason: "{item.decision_reason}"</Text>
      {/if}
    </Box>
  {/if}
</Box>

<style>
  :global(.approval-card) {
    background-color: var(--bg-surface);
    border: 1px solid var(--border-default);
    border-radius: 4px;
    padding: 12px 14px;
    display: flex;
    flex-direction: column;
    gap: 8px;
    cursor: pointer;
    transition: border-color 0.15s ease, background-color 0.15s ease;
  }

  :global(.approval-card:hover) {
    background-color: var(--bg-hover);
    border-color: var(--border-hover);
  }

  :global(.approval-card.selected) {
    border-color: var(--color-primary, #3574f0);
    background-color: rgba(53, 116, 240, 0.04);
  }

  :global(.approval-card .conn-pill) {
    font-family: var(--font-code);
  }

  :global(.approval-card .requester-row) {
    font-size: var(--font-size-sm, 13px);
  }

  :global(.approval-card .token-tag) {
    color: var(--syntax-number, #6897bb);
  }

  :global(.approval-card .warn-icon) {
    display: inline-flex;
    align-items: center;
    color: var(--action-warning);
  }

  :global(.approval-card .sql-preview) {
    font-size: var(--font-size-sm, 13px);
    background-color: var(--bg-canvas);
    border: 1px solid var(--border-subtle);
    border-radius: 3px;
    padding: 8px 10px;
    color: var(--text-primary);
    white-space: pre-wrap;
    word-break: break-word;
    max-height: 80px;
    overflow: hidden;
  }

  :global(.approval-card .decision-meta) {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  :global(.approval-card .reason-note) {
    font-style: italic;
  }
</style>
