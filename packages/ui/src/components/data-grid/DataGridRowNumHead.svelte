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

<th class="ohjanus-data-grid-row-num-head {className}" {style} scope="col">
  {#if children}{@render children()}{:else if index !== undefined}{index + 1}{/if}
</th>
