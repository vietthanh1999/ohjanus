<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext } from 'svelte';

  interface Props {
    align?: 'start' | 'center' | 'end';
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { align = 'center', class: className = '', style = '', children }: Props = $props();
  const ctx = getContext<{ open: boolean }>('ohjanus-hovercard');
</script>

{#if ctx?.open}
  <div class="ohjanus-hovercard-content ohjanus-hovercard-align-{align} {className}" {style} role="dialog">
    {@render children?.()}
  </div>
{/if}
