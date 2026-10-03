<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    colspan?: number;
    title?: string;
    children?: Snippet;
  }

  let { class: className = '', style = '', colspan, title, children }: Props = $props();
</script>

<td class="ohjanus-table-cell {className}" {style} colspan={colspan} {title}>
  {@render children?.()}
</td>
