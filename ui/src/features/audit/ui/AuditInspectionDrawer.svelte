<script lang="ts">
  import type { AuditRecord } from '@/entities/audit-record';
  import { Box, Flex, Stack, Badge, Text, Button } from "@ohjanus/ui";
  import { Icon } from "@ohjanus/icons";

  interface Props {
    item: AuditRecord | null;
    onclose: () => void;
  }

  let { item, onclose }: Props = $props();
</script>

{#if item}
  <Box class="audit-drawer">
    <Flex class="drawer-header" align="center" justify="between">
      <Text weight="semibold" size="md" class="drawer-title">Audit Detail: {item.id}</Text>
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
      <Flex class="drawer-row" align="center" justify="between">
        <Text size="md" color="muted">Timestamp:</Text>
        <Text size="md" mono>{item.ts}</Text>
      </Flex>
      <Flex class="drawer-row" align="center" justify="between">
        <Text size="md" color="muted">Request ID:</Text>
        <Text size="md" mono>{item.request_id}</Text>
      </Flex>
      <Flex class="drawer-row" align="center" justify="between">
        <Text size="md" color="muted">Token ID:</Text>
        <Text size="md" mono>{item.token_id}</Text>
      </Flex>
      <Flex class="drawer-row" align="center" justify="between">
        <Text size="md" color="muted">Requesting Client:</Text>
        <Text size="md">{item.client}</Text>
      </Flex>
      <Flex class="drawer-row" align="center" justify="between">
        <Text size="md" color="muted">Database Target:</Text>
        <Text size="md">{item.connection}</Text>
      </Flex>
      <Flex class="drawer-row" align="center" justify="between">
        <Text size="md" color="muted">Tables Touched:</Text>
        <Flex align="center" gap="4px" inline>
          {#if item.tables && item.tables.length > 0}
            {#each item.tables as t}
              <Badge variant="default" size="md" class="tbl-pill">{t}</Badge>
            {/each}
          {:else}
            <Text size="md" color="muted">—</Text>
          {/if}
        </Flex>
      </Flex>
      <Flex class="drawer-row" align="center" justify="between">
        <Text size="md" color="muted">Policy Decision:</Text>
        <Box>
          {#if item.policy_decision === "ALLOW"}
            <Badge variant="success" size="md">ALLOW</Badge>
          {:else if item.policy_decision === "REQUIRE_APPROVAL"}
            <Badge variant="warning" size="md">APPROVAL</Badge>
          {:else}
            <Badge variant="danger" size="md">DENY</Badge>
          {/if}
        </Box>
      </Flex>
      <Flex class="drawer-row" align="center" justify="between">
        <Text size="md" color="muted">Policy Rule Matched:</Text>
        <Text size="md" mono>{item.policy_rule}</Text>
      </Flex>
      <Flex class="drawer-row" align="center" justify="between">
        <Text size="md" color="muted">Execution Duration:</Text>
        <Text size="md" mono>{item.duration_ms} ms</Text>
      </Flex>

      <Text
        size="md"
        weight="semibold"
        color="secondary"
        style="margin-top: 8px;"
      >
        Executed Query
      </Text>
      <Box class="sql-box code-text">{item.sql_normalized}</Box>
    </Stack>
  </Box>
{/if}

<style>
  :global(.audit-drawer) {
    width: 380px;
    background-color: var(--bg-surface);
    border-left: 1px solid var(--border-default);
    display: flex;
    flex-direction: column;
    flex-shrink: 0;
  }

  :global(.audit-drawer .drawer-header) {
    height: 36px;
    padding: 0 12px;
    border-bottom: 1px solid var(--border-default);
  }

  :global(.drawer-body) {
    padding: 12px;
    overflow-y: auto;
    flex: 1;
  }

  :global(.audit-drawer .drawer-row) {
    padding: 4px 0;
    border-bottom: 1px solid var(--border-subtle);
  }

  :global(.tbl-pill) {
    font-family: var(--font-code);
  }

  :global(.sql-box) {
    background-color: var(--bg-canvas);
    border: 1px solid var(--border-subtle);
    border-radius: 4px;
    padding: 10px;
    font-size: var(--font-size-sm, 12px);
    color: var(--text-primary);
    white-space: pre-wrap;
    word-break: break-all;
    max-height: 240px;
    overflow-y: auto;
    margin: 0;
  }
</style>
