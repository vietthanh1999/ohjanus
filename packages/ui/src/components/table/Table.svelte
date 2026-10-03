<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { class: className = '', style = '', children }: Props = $props();
</script>

<div class="ohjanus-table-container">
  <table class="ohjanus-table {className}" {style}>
    {@render children?.()}
  </table>
</div>
