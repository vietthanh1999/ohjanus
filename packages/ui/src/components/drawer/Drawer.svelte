<script lang="ts">
  import type { Snippet } from 'svelte';
  import { DialogPrimitive } from '@ohjanus/primitives';

  interface Props {
    open?: boolean;
    side?: 'left' | 'right' | 'top' | 'bottom';
    onClose?: () => void;
    children?: Snippet;
  }

  let { open = $bindable(false), side = 'right', onClose, children }: Props = $props();

  function close() {
    open = false;
    onClose?.();
  }
</script>

<DialogPrimitive {open} onClose={close}>
  <div class="ohjanus-drawer-wrapper ohjanus-drawer-{side}">
    {@render children?.()}
  </div>
</DialogPrimitive>
