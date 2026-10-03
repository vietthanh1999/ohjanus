<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    href?: string;
    class?: string;
    style?: string;
    onclick?: (e: MouseEvent) => void;
    children?: Snippet;
  }

  let { href = '#', class: className = '', style = '', onclick, children }: Props = $props();
</script>

<a class="ohjanus-breadcrumb-link {className}" {href} {style} {onclick}>
  {@render children?.()}
</a>
