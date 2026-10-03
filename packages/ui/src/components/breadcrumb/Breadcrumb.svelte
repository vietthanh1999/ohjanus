<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { class: className = '', style = '', children }: Props = $props();
</script>

<nav aria-label="breadcrumb" class="ohjanus-breadcrumb {className}" {style}>
  {@render children?.()}
</nav>
