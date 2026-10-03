<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { class: className = '', style = '', children }: Props = $props();
  const ctx = getContext<{ open: boolean }>('ohjanus-menu-sub');
</script>

{#if ctx?.open}
  <div class="ohjanus-dropdown-sub-content {className}" {style} role="menu">
    {@render children?.()}
  </div>
{/if}
