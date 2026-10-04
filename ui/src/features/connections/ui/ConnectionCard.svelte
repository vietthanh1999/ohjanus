<script lang="ts">
  import type { ConnectionItem } from '@/entities/connection';
  import { Card, Badge, Button, Text, Box, Flex, Stack } from "@ohjanus/ui";
  import { Icon } from "@ohjanus/icons";

  interface Props {
    conn: ConnectionItem;
    isTesting: boolean;
    ontest: (name: string) => void;
  }

  let { conn, isTesting, ontest }: Props = $props();
</script>

<Card class="conn-card-item">
  {#snippet header()}
    <Flex class="card-header-inner" justify="between" align="center">
      <Flex class="header-left" align="center" gap="8px">
        <Box class="db-icon">
          <Icon name="database" size={15} color="#3B82F6" />
        </Box>
        <Text class="conn-title" weight="semibold">{conn.name}</Text>
        {#if conn.readonly}
          <Badge variant="warning" size="sm">
            <Icon name="lock" size={11} />
            <Text size="xs" color="warning">READ-ONLY</Text>
          </Badge>
        {:else}
          <Badge variant="info" size="sm">
            <Icon name="lightning" size={11} />
            <Text size="xs" style="color: #56A8F5;">READ-WRITE</Text>
          </Badge>
        {/if}
      </Flex>
      {#if conn.status === "healthy"}
        <Badge variant="success" size="sm">
          <Box as="span" class="status-dot ok" />HEALTHY
        </Badge>
      {:else}
        <Badge variant="danger" size="sm">
          <Box as="span" class="status-dot bad" />DEGRADED
        </Badge>
      {/if}
    </Flex>
  {/snippet}

  <Stack class="conn-details" gap="6px">
    <Flex class="detail-row" justify="between">
      <Text class="label" size="md" color="muted">Driver:</Text>
      <Text class="value" size="md" mono>{conn.driver}</Text>
    </Flex>
    <Flex class="detail-row" justify="between">
      <Text class="label" size="md" color="muted">Last Health Ping:</Text>
      <Text class="value" size="md">
        {conn.last_ping_at}{conn.latency_ms !== undefined ? ` (${conn.latency_ms} ms)` : ""}
      </Text>
    </Flex>
  </Stack>

  <!-- Security Guardrails -->
  <Box class="guardrails-section">
    <Text class="sec-label" size="xs" weight="bold" color="muted">ALLOWED SCHEMAS:</Text>
    <Flex class="tags-list" wrap gap="4px">
      {#each conn.allowed_schemas as sch}
        <Badge variant="default" size="sm">{sch}</Badge>
      {/each}
    </Flex>

    {#if conn.denied_tables.length > 0}
      <Text
        class="sec-label"
        size="xs"
        weight="bold"
        color="danger"
        style="margin-top: 6px;"
      >
        DENIED TABLES:
      </Text>
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
      onclick={() => ontest(conn.name)}
      loading={isTesting}
    >
      {#if isTesting}
        Pinging...
      {:else}
        <Icon name="refresh" size={12} />
        <Text size="md">Test Connection</Text>
      {/if}
    </Button>
  {/snippet}
</Card>

<style>
  :global(.conn-card-item) {
    background-color: var(--bg-surface);
    border: 1px solid var(--border-default);
    border-radius: 4px;
    display: flex;
    flex-direction: column;
  }

  :global(.conn-card-item .card-header-inner) {
    width: 100%;
  }

  :global(.conn-card-item .header-left) {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  :global(.conn-card-item .status-dot) {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    margin-right: 4px;
    display: inline-block;
  }

  :global(.conn-card-item .status-dot.ok) {
    background-color: var(--action-success, #22c55e);
  }

  :global(.conn-card-item .status-dot.bad) {
    background-color: var(--action-danger, #ef4444);
  }

  :global(.conn-card-item .guardrails-section) {
    margin-top: 10px;
    padding-top: 10px;
    border-top: 1px solid var(--border-subtle);
  }

  :global(.conn-card-item .sec-label) {
    font-size: 10px;
    letter-spacing: 0.5px;
    margin-bottom: 4px;
    display: block;
  }
</style>
