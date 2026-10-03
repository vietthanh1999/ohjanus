<script lang="ts">
  import type { Snippet } from 'svelte';
  import { setContext } from 'svelte';

  interface Props {
    direction?: 'horizontal' | 'vertical';
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { direction = 'horizontal', class: className = '', style = '', children }: Props = $props();

  setContext('ohjanus-resizable', {
    get direction() { return direction; }
  });
</script>

<div
  class="ohjanus-resizable-group ohjanus-resizable-{direction} {className}"
  {style}
  data-direction={direction}
>
  {@render children?.()}
</div>
