<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    /** Fixed column width, e.g. "280px". */
    width?: string;
    children?: Snippet;
  }

  let { class: className = '', style = '', width, children }: Props = $props();
</script>

<th class="ohjanus-data-grid-head-cell {className}" style="{width ? `width: ${width};` : ''}{style ?? ''}">
  {@render children?.()}
</th>
