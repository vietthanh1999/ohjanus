<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { class: className = '', style = '', children }: Props = $props();
</script>

<h5 class="ohjanus-alert-title {className}" {style}>
  {@render children?.()}
</h5>
