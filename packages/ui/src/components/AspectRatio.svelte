<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    ratio?: number;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { ratio = 16 / 9, class: className = '', style = '', children }: Props = $props();
</script>

<div class="ohjanus-aspect-ratio {className}" style="aspect-ratio: {ratio}; {style}">
  {@render children?.()}
</div>
