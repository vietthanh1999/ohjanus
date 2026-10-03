<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { class: className = '', style = '', children }: Props = $props();
</script>

<span role="link" aria-disabled="true" aria-current="page" class="ohjanus-breadcrumb-page {className}" {style}>
  {@render children?.()}
</span>
