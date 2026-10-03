<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    orientation?: 'vertical' | 'horizontal' | 'both';
    maxHeight?: string;
    maxWidth?: string;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let {
    orientation = 'vertical',
    maxHeight,
    maxWidth,
    class: className = '',
    style = '',
    children
  }: Props = $props();

  const overflow = $derived(
    orientation === 'both' ? 'auto' : orientation === 'horizontal' ? 'hidden auto' : 'auto hidden'
  );
</script>

<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
<div
  class="ohjanus-scroll-area {className}"
  role="region"
  style="overflow: {overflow}; {maxHeight ? `max-height: ${maxHeight};` : ''} {maxWidth ? `max-width: ${maxWidth};` : ''} {style}"
  tabindex="0"
>
  {@render children?.()}
</div>
