<script lang="ts">
  import { onMount, type Snippet } from 'svelte';
  import Portal from './Portal.svelte';

  interface Props {
    open: boolean;
    onClose: () => void;
    children?: Snippet;
    closeOnEscape?: boolean;
    closeOnClickOutside?: boolean;
    titleId?: string;
    descriptionId?: string;
  }

  let {
    open,
    onClose,
    children,
    closeOnEscape = true,
    closeOnClickOutside = true,
    titleId,
    descriptionId
  }: Props = $props();

  function handleKeydown(e: KeyboardEvent) {
    if (open && closeOnEscape && e.key === 'Escape') {
      e.stopPropagation();
      onClose();
    }
  }

  function handleBackdropClick(e: MouseEvent) {
    if (closeOnClickOutside && e.target === e.currentTarget) {
      onClose();
    }
  }
</script>

<svelte:window onkeydown={handleKeydown} />

{#if open}
  <Portal>
    <!-- Backdrop / Overlay wrapper -->
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div
      role="presentation"
      class="ohjanus-dialog-backdrop"
      onclick={handleBackdropClick}
    >
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby={titleId}
        aria-describedby={descriptionId}
        class="ohjanus-dialog-container"
      >
        {@render children?.()}
      </div>
    </div>
  </Portal>
{/if}

<style>
  .ohjanus-dialog-backdrop {
    position: fixed;
    inset: 0;
    z-index: 1000;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .ohjanus-dialog-container {
    position: relative;
    max-width: 100%;
    max-height: 100%;
    outline: none;
  }
</style>
