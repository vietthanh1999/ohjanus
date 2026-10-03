<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    width?: string;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { width = '520px', class: className = '', style = '', children }: Props = $props();
</script>

<div class="ohjanus-dialog-content {className}" style="width: {width}; {style}">
  {@render children?.()}
</div>
