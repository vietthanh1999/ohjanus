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

  let search = $state('');

  function setOpen(v: boolean) {
    open = v;
    onOpenChange?.(v);
  }

  setContext('ohjanus-combobox', {
    get open() { return open; },
    setOpen,
    get search() { return search; },
    setSearch: (v: string) => (search = v)
  });
</script>

<div class="ohjanus-combobox {className}" {style} data-open={open}>
  {@render children?.()}
</div>
