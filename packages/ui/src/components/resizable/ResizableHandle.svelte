<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext } from 'svelte';

  interface Props {
    disabled?: boolean;
    class?: string;
    style?: string;
    children?: Snippet;
    onResize?: (delta: number) => void;
  }

  let { disabled = false, class: className = '', style = '', children, onResize }: Props = $props();
  const ctx = getContext<{ direction: string }>('ohjanus-resizable');
  const vertical = $derived((ctx?.direction ?? 'horizontal') === 'vertical');

  let dragging = $state(false);
  let startPos = 0;

  function handlePointerDown(e: PointerEvent) {
    if (disabled) return;
    dragging = true;
    startPos = vertical ? e.clientY : e.clientX;
    (e.target as HTMLElement).setPointerCapture?.(e.pointerId);
  }

  function handlePointerMove(e: PointerEvent) {
    if (!dragging || disabled) return;
    const current = vertical ? e.clientY : e.clientX;
    onResize?.(current - startPos);
  }

  function handlePointerUp() {
    dragging = false;
  }
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
<div
  class="ohjanus-resizable-handle {className}"
  class:dragging
  {style}
  role="separator"
  aria-orientation={vertical ? 'vertical' : 'horizontal'}
  aria-disabled={disabled || undefined}
  tabindex={disabled ? -1 : 0}
  onpointerdown={handlePointerDown}
  onpointermove={handlePointerMove}
  onpointerup={handlePointerUp}
>
  {#if children}
    {@render children()}
  {:else}
    <span class="ohjanus-resizable-grip" aria-hidden="true"></span>
  {/if}
</div>
