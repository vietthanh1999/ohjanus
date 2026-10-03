<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    onclick?: (e: MouseEvent) => void;
    children?: Snippet;
    [key: string]: any;
  }

  let { class: className = '', style = '', onclick, children, ...restProps }: Props = $props();
</script>

<button type="button" class="ohjanus-alert-dialog-trigger {className}" {style} {onclick} {...restProps}>
  {@render children?.()}
</button>
