<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    open: boolean;
    class?: string;
    children: Snippet;
  }

  let { open, class: className = '', children }: Props = $props();
</script>

{#if open}
  <div class={className}>
    {@render children()}
  </div>
{/if}
