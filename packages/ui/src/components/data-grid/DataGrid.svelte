<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    /** Pin the header row while scrolling. Defaults to true. */
    stickyHeader?: boolean;
    children?: Snippet;
  }

  let { class: className = '', style = '', stickyHeader = true, children }: Props = $props();
</script>

<div
  class="ohjanus-data-grid {stickyHeader ? 'ohjanus-data-grid-sticky' : ''} {className}"
  {style}
>
  <table class="ohjanus-data-grid-table">
    {@render children?.()}
  </table>
</div>
