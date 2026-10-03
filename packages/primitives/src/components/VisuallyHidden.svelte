<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    children?: Snippet;
    as?: string;
    class?: string;
  }

  let { children, as = 'span', class: className = '' }: Props = $props();
</script>

<svelte:element
  this={as}
  class={className}
  style="position: absolute; border: 0; width: 1px; height: 1px; padding: 0; margin: -1px; overflow: hidden; clip: rect(0, 0, 0, 0); white-space: nowrap; word-wrap: normal;"
>
  {@render children?.()}
</svelte:element>
