<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    onclick?: (e: MouseEvent) => void;
    class?: string;
    style?: string;
    children?: Snippet;
    [key: string]: any;
  }

  let { onclick, class: className = '', style = '', children, ...restProps }: Props = $props();
</script>

<button type="button" class="ohjanus-attachment-action {className}" {style} {onclick} {...restProps}>
  {@render children?.()}
</button>
