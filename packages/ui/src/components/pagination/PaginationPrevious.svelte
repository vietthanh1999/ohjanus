<script lang="ts">
  import type { Snippet } from 'svelte';
  import PaginationLink from './PaginationLink.svelte';

  interface Props {
    href?: string;
    disabled?: boolean;
    onclick?: (e: MouseEvent) => void;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { href, disabled = false, onclick, class: className = '', style = '', children }: Props = $props();
</script>

<PaginationLink {href} {disabled} {onclick} class="ohjanus-pagination-previous {className}" {style} aria-label="Go to previous page">
  <span aria-hidden="true">‹</span>
  {@render children?.()}
</PaginationLink>
