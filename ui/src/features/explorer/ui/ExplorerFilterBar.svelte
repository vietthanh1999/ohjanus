<script lang="ts">
  import { Box, Flex, Button, Text, Input } from '@ohjanus/ui';

  interface Props {
    query: string;
    onddl?: () => void;
  }

  let { query = $bindable(''), onddl }: Props = $props();
</script>

<Flex align="center" gap="6px" class="explorer-subtoolbar">
  <Box class="filter-input-wrap">
    <Input
      size="sm"
      placeholder="Filter connections / tables…"
      bind:value={query}
      class="tree-filter-input"
    />
  </Box>
  <Button
    variant="secondary"
    size="xs"
    class="ddl-btn"
    title="Generate Table DDL"
    onclick={() => onddl?.()}
  >
    <Text size="xs" weight="bold" color="muted" mono>DDL</Text>
  </Button>
</Flex>

<style>
  :global(.explorer-subtoolbar) {
    height: var(--toolbar-height, 32px);
    padding: 0 8px;
    border-bottom: 1px solid var(--border-subtle);
    gap: 6px;
    flex-shrink: 0;
  }

  :global(.filter-input-wrap) {
    flex: 1;
    min-width: 0;
  }

  :global(.filter-input-wrap .ohjanus-input-wrapper) {
    height: 24px;
    background-color: var(--bg-canvas);
    border: 1px solid var(--border-default);
    border-radius: var(--radius-sm, 4px);
    padding: 0 6px;
  }

  :global(.filter-input-wrap .tree-filter-input) {
    font-size: 12px;
    color: var(--text-primary);
  }

  :global(.explorer-subtoolbar .ddl-btn) {
    height: 24px !important;
    padding: 0 6px !important;
    font-size: 10px !important;
    line-height: 1 !important;
  }
</style>
