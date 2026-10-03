<script lang="ts">
  import type { Snippet } from 'svelte';
  import { setContext } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { class: className = '', style = '', children }: Props = $props();
  let open = $state(false);

  setContext('ohjanus-menu-sub', {
    get open() { return open; },
    setOpen: (v: boolean) => (open = v),
    toggle: () => (open = !open)
  });
</script>

<div class="ohjanus-dropdown-sub {className}" {style} data-open={open}>
  {@render children?.()}
</div>
