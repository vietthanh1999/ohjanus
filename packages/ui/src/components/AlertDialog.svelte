<script lang="ts">
  import type { Snippet } from 'svelte';
  import { DialogPrimitive } from '@ohjanus/primitives';

  interface Props {
    open?: boolean;
    onClose?: () => void;
    onCancel?: () => void;
    onAction?: () => void;
    children?: Snippet;
  }

  // AlertDialog: confirm-only — no click-outside dismiss, Escape only via Cancel/Action.
  let { open = $bindable(false), onClose, onCancel, onAction, children }: Props = $props();

  function close() {
    open = false;
    onClose?.();
  }
</script>

<DialogPrimitive {open} onClose={close} closeOnClickOutside={false}>
  {@render children?.()}
</DialogPrimitive>
