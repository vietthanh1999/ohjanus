<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { class: className = '', style = '', children }: Props = $props();
  const ctx = getContext<{ open: boolean }>('ohjanus-select');
</script>

{#if ctx?.open}
  <div class="ohjanus-select-list {className}" {style} role="listbox">
    {@render children?.()}
  </div>
{/if}
