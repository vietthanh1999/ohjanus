<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { class: className = '', style = '', children }: Props = $props();
</script>

<h2 class="ohjanus-alert-dialog-title {className}" {style}>
  {@render children?.()}
</h2>
