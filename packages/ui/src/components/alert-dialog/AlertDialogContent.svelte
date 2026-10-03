<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    width?: string;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { width = '480px', class: className = '', style = '', children }: Props = $props();
</script>

<div class="ohjanus-alert-dialog-content {className}" style="width: {width}; {style}" role="alertdialog" aria-modal="true">
  {@render children?.()}
</div>
