<script lang="ts">
  import type { Snippet } from 'svelte';
  import { Icon } from '@ohjanus/icons';
  import { DialogPrimitive } from '@ohjanus/primitives';

  interface Props {
    open: boolean;
    onClose: () => void;
    title?: string;
    width?: string;
    class?: string;
    children?: Snippet;
    footer?: Snippet;
  }

  let {
    open,
    onClose,
    title = '',
    width = '520px',
    class: className = '',
    children,
    footer
  }: Props = $props();
</script>

<DialogPrimitive {open} {onClose}>
  <div class="ohjanus-modal-panel {className}" style="width: {width};">
    <div class="ohjanus-modal-header">
      <h3 class="ohjanus-modal-title">{title}</h3>
      <button
        type="button"
        class="ohjanus-modal-close"
        onclick={onClose}
        aria-label="Close modal"
      >
        <Icon name="x" size={14} />
      </button>
    </div>

    <div class="ohjanus-modal-body">
      {@render children?.()}
    </div>

    {#if footer}
      <div class="ohjanus-modal-footer">
        {@render footer()}
      </div>
    {/if}
  </div>
</DialogPrimitive>
