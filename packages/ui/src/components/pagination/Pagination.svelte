<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    label?: string;
    children?: Snippet;
  }

  let { class: className = '', style = '', label = 'Pagination', children }: Props = $props();
</script>

<nav class="ohjanus-pagination {className}" {style} aria-label={label}>
  {@render children?.()}
</nav>
