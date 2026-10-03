<script lang="ts">
  import type { Snippet } from 'svelte';
  import Portal from './Portal.svelte';
  import FocusTrap from './FocusTrap.svelte';

  interface Props {
    open: boolean;
    onClose: () => void;
    children?: Snippet;
    modal?: boolean;
    closeOnEscape?: boolean;
    closeOnClickOutside?: boolean;
    titleId?: string;
    descriptionId?: string;
  }

  let {
    open,
    onClose,
    children,
    modal = true,
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
    if (modal && closeOnClickOutside && e.target === e.currentTarget) {
      onClose();
    }
  }

  // Scroll-lock the page while the dialog is open.
  $effect(() => {
    if (!open || typeof document === 'undefined') return;
    const previous = document.body.style.overflow;
    document.body.style.overflow = 'hidden';
    return () => {
      document.body.style.overflow = previous;
    };
  });
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
      class:ohjanus-dialog-backdrop-transparent={!modal}
      onclick={handleBackdropClick}
    >
      <FocusTrap>
        <div
          role="dialog"
          aria-modal={modal}
          aria-labelledby={titleId}
          aria-describedby={descriptionId}
          class="ohjanus-dialog-container"
        >
          {@render children?.()}
        </div>
      </FocusTrap>
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
  .ohjanus-dialog-backdrop-transparent {
    pointer-events: none;
  }
  .ohjanus-dialog-backdrop-transparent > * {
    pointer-events: auto;
  }
  .ohjanus-dialog-container {
    position: relative;
    max-width: 100%;
    max-height: 100%;
    outline: none;
  }
</style>
