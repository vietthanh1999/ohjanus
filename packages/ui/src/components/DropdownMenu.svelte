<script lang="ts">
  import type { Snippet } from 'svelte';
  import { setContext, onMount } from 'svelte';

  interface Props {
    open?: boolean;
    class?: string;
    style?: string;
    children?: Snippet;
    onOpenChange?: (open: boolean) => void;
  }

  let { open = $bindable(false), class: className = '', style = '', children, onOpenChange }: Props = $props();

  function setOpen(v: boolean) {
    open = v;
    onOpenChange?.(v);
  }

  setContext('ohjanus-menu', {
    get open() { return open; },
    setOpen,
    toggle: () => setOpen(!open),
    close: () => setOpen(false)
  });
</script>

<div class="ohjanus-dropdown-menu {className}" {style} data-open={open}>
  {@render children?.()}
</div>
