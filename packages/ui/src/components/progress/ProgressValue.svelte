<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { class: className = '', style = '', children }: Props = $props();
  const ctx = getContext<{ value: number; max: number; percent: number }>('ohjanus-progress');
</script>

<span class="ohjanus-progress-value {className}" {style}>
  {#if children}
    {@render children()}
  {:else}
    {Math.round(ctx?.percent ?? 0)}%
  {/if}
</span>
