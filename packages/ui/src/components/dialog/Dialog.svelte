<script lang="ts">
  import type { Snippet } from 'svelte';
  import { DialogPrimitive } from '@ohjanus/primitives';

  interface Props {
    open: boolean;
    onClose: () => void;
    children?: Snippet;
  }

  let { open, onClose, children }: Props = $props();
</script>

<DialogPrimitive {open} {onClose}>
  {@render children?.()}
</DialogPrimitive>
