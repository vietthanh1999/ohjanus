<script lang="ts">
  import type { Snippet } from 'svelte';
  import { setContext, getContext } from 'svelte';

  interface Props {
    open?: boolean;
    class?: string;
    style?: string;
    children?: Snippet;
    onOpenChange?: (open: boolean) => void;
  }

  let { open = $bindable(false), class: className = '', style = '', children, onOpenChange }: Props = $props();

  function toggle() {
    open = !open;
    onOpenChange?.(open);
  }

  function setOpen(v: boolean) {
    open = v;
    onOpenChange?.(v);
  }

  setContext('ohjanus-collapsible', {
    get open() { return open; },
    toggle,
    setOpen
  });
</script>

<div class="ohjanus-collapsible {className}" {style} data-open={open}>
  {@render children?.()}
</div>
