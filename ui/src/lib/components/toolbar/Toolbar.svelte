<script lang="ts">
  import type { Snippet } from 'svelte';
  import { Flex } from '@ohjanus/ui';

  interface Props {
    left?: Snippet;
    right?: Snippet;
    children?: Snippet;
    class?: string;
  }

  let { left, right, children, class: className = '' }: Props = $props();
</script>

<Flex align="center" justify="between" class="table-toolbar {className}">
  {#if left}
    <Flex align="center" gap="2px" class="toolbar-left">
      {@render left()}
    </Flex>
  {/if}
  {#if children}
    {@render children()}
  {/if}
  {#if right}
    <Flex align="center" gap="2px" class="toolbar-right">
      {@render right()}
    </Flex>
  {/if}
</Flex>

<style>
  :global(.table-toolbar) {
    height: var(--toolbar-height, 32px);
    background-color: var(--bg-toolbar);
    border-bottom: 1px solid var(--border-subtle);
    padding: 0 8px;
    flex-shrink: 0;
  }
</style>
