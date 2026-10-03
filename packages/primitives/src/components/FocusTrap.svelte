<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    disabled?: boolean;
    children: Snippet;
  }

  let { disabled = false, children }: Props = $props();

  let container: HTMLElement | undefined = $state();
  let previousFocused: Element | null = null;

  const FOCUSABLE =
    'a[href], button:not([disabled]), input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])';

  function focusables(): HTMLElement[] {
    if (!container) return [];
    return [...container.querySelectorAll<HTMLElement>(FOCUSABLE)].filter(
      (el) => el.offsetParent !== null || el === document.activeElement
    );
  }

  $effect(() => {
    if (disabled || !container) return;
    previousFocused = document.activeElement;
    const first = focusables()[0] ?? container;
    first.focus({ preventScroll: true });
    return () => {
      if (previousFocused instanceof HTMLElement) {
        previousFocused.focus({ preventScroll: true });
      }
    };
  });

  function handleKeydown(e: KeyboardEvent) {
    if (disabled || e.key !== 'Tab' || !container) return;
    const items = focusables();
    if (items.length === 0) {
      e.preventDefault();
      return;
    }
    const first = items[0];
    const last = items[items.length - 1];
    if (e.shiftKey && document.activeElement === first) {
      e.preventDefault();
      last.focus();
    } else if (!e.shiftKey && document.activeElement === last) {
      e.preventDefault();
      first.focus();
    }
  }
</script>

<!-- Focus trap container: intentionally intercepts Tab to cycle focus. -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div bind:this={container} onkeydown={handleKeydown} tabindex="-1" style="display: contents;">
  {@render children()}
</div>
