<script lang="ts">
  import { Text, Flex } from '@ohjanus/ui';
  import type { Snippet } from 'svelte';

  interface Props {
    count: number;
    unit?: string;
    extra?: Snippet;
  }

  let { count, unit = 'rows', extra }: Props = $props();
</script>

{#if count > 0}
  <Flex align="center" gap="8px" class="floating-row-badge">
    <Text size="sm" title="{count} {unit}">{count} {unit}</Text>
    {#if extra}
      {@render extra()}
    {/if}
  </Flex>
{/if}

<style>
  :global(.floating-row-badge) {
    position: absolute;
    bottom: 12px;
    right: 18px;
    height: var(--control-height-xs, 24px);
    padding: 0 10px;
    background-color: var(--bg-toolbar);
    border: 1px solid var(--border-default);
    border-radius: 12px;
    font-size: 11px;
    color: var(--text-primary);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.45);
    z-index: 5;
    user-select: none;
  }
</style>
