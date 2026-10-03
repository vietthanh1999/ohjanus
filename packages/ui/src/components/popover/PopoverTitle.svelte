<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { class: className = '', style = '', children }: Props = $props();
</script>

<h4 class="ohjanus-popover-title {className}" {style}>
  {@render children?.()}
</h4>
