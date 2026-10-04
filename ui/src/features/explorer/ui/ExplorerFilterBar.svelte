<script lang="ts">
  import { Box, Flex, Input } from '@ohjanus/ui';

  interface Props {
    query: string;
  }

  let { query = $bindable('') }: Props = $props();
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
</style>
