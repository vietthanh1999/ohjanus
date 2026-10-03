<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    /** Zero-based row index. Renders index + 1. */
    index?: number;
    children?: Snippet;
  }

  let { class: className = '', style = '', index, children }: Props = $props();
</script>

<td class="ohjanus-data-grid-row-num {className}" {style}>
  {#if children}{@render children()}{:else if index !== undefined}{index + 1}{/if}
</td>
