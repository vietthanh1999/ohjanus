<script lang="ts">
  import type { Snippet } from 'svelte';
  import { DialogPrimitive } from '@ohjanus/primitives';

  interface Props {
    open: boolean;
    onClose: () => void;
    modal?: boolean;
    closeOnEscape?: boolean;
    closeOnClickOutside?: boolean;
    titleId?: string;
    descriptionId?: string;
    children?: Snippet;
  }

  let {
    open,
    onClose,
    modal = true,
    closeOnEscape = true,
    closeOnClickOutside = true,
    titleId,
    descriptionId,
    children
  }: Props = $props();
</script>

<DialogPrimitive {open} {onClose} {modal} {closeOnEscape} {closeOnClickOutside} {titleId} {descriptionId}>
  {@render children?.()}
</DialogPrimitive>
