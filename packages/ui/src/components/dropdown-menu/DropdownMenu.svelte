<script lang="ts">
  import type { Snippet } from 'svelte';
  import { setContext } from 'svelte';

  interface Props {
    open?: boolean;
    class?: string;
    style?: string;
    children?: Snippet;
    onOpenChange?: (open: boolean) => void;
  }

  let { open = $bindable(false), class: className = '', style = '', children, onOpenChange }: Props = $props();

  let root: HTMLElement | undefined = $state();
  let contentEl: HTMLElement | undefined = $state();

  function setOpen(v: boolean) {
    open = v;
    onOpenChange?.(v);
  }

  setContext('ohjanus-menu', {
    get open() { return open; },
    get anchor() { return root; },
    setContentEl: (el: HTMLElement | undefined) => (contentEl = el),
    setOpen,
    toggle: () => setOpen(!open),
    close: () => setOpen(false)
  });

  // Close on outside click / Escape — the content is portaled to <body>,
  // so both the trigger root and the portaled content count as "inside".
  $effect(() => {
    if (!open) return;
    const onPointerDown = (e: PointerEvent) => {
      const t = e.target as Node | null;
      if (t && (root?.contains(t) || contentEl?.contains(t))) return;
      setOpen(false);
    };
    const onKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') setOpen(false);
    };
    document.addEventListener('pointerdown', onPointerDown, true);
    document.addEventListener('keydown', onKeyDown);
    return () => {
      document.removeEventListener('pointerdown', onPointerDown, true);
      document.removeEventListener('keydown', onKeyDown);
    };
  });
</script>

<div bind:this={root} class="ohjanus-dropdown-menu {className}" {style} data-open={open}>
  {@render children?.()}
</div>
