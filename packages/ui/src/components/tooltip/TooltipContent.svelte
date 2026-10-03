<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext } from 'svelte';

  interface Props {
    side?: 'top' | 'bottom' | 'left' | 'right';
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { side = 'top', class: className = '', style = '', children }: Props = $props();
  const ctx = getContext<{ open: boolean }>('ohjanus-tooltip');
</script>

{#if ctx?.open}
  <span class="ohjanus-tooltip-content ohjanus-tooltip-{side} {className}" {style} role="tooltip">
    {@render children?.()}
  </span>
{/if}
