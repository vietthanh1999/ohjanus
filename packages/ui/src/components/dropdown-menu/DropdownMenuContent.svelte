<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext } from 'svelte';

  interface Props {
    align?: 'start' | 'center' | 'end';
    side?: 'top' | 'bottom' | 'left' | 'right';
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { align = 'start', side = 'bottom', class: className = '', style = '', children }: Props = $props();
  const ctx = getContext<{ open: boolean }>('ohjanus-menu');
</script>

{#if ctx?.open}
  <div
    class="ohjanus-dropdown-content ohjanus-dropdown-align-{align} {className}"
    {style}
    role="menu"
    data-side={side}
  >
    {@render children?.()}
  </div>
{/if}
