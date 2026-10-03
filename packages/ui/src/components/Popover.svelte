<script lang="ts">
  import type { Snippet } from 'svelte';
  import { setContext } from 'svelte';

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

  function toggle() {
    setOpen(!open);
  }

  setContext('ohjanus-popover', {
    get open() { return open; },
    setOpen,
    toggle
  });
</script>

<div class="ohjanus-popover {className}" {style} data-open={open}>
  {@render children?.()}
</div>
