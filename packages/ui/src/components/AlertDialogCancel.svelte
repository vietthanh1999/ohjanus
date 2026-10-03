<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    onclick?: (e: MouseEvent) => void;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { onclick, class: className = '', style = '', children }: Props = $props();
</script>

<button type="button" class="ohjanus-btn ohjanus-btn-outline ohjanus-btn-sm {className}" {style} {onclick}>
  {@render children?.()}
</button>
