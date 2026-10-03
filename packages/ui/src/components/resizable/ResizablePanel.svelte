<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    defaultSize?: number;
    minSize?: number;
    maxSize?: number;
    collapsible?: boolean;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let {
    defaultSize,
    minSize,
    maxSize,
    collapsible = false,
    class: className = '',
    style = '',
    children
  }: Props = $props();
</script>

<div
  class="ohjanus-resizable-panel {className}"
  style="{defaultSize !== undefined ? `flex: ${defaultSize} ${defaultSize} 0%;` : ''} {minSize !== undefined ? `min-width: ${minSize}%;` : ''} {maxSize !== undefined ? `max-width: ${maxSize}%;` : ''} {style}"
  data-collapsible={collapsible || undefined}
>
  {@render children?.()}
</div>
