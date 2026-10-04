<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext, tick } from 'svelte';
  import { portal, positionFloating } from '../../internal/floating.js';

  interface Props {
    side?: 'top' | 'bottom' | 'left' | 'right';
    class?: string;
    style?: string;
    children?: Snippet;
  }

  interface TooltipContext {
    open: boolean;
    anchor?: HTMLElement;
  }

  let { side = 'top', class: className = '', style = '', children }: Props = $props();
  const ctx = getContext<TooltipContext>('ohjanus-tooltip');
  let el: HTMLElement | undefined = $state();

  async function updatePosition() {
    await tick();
    if (el && ctx?.anchor) positionFloating(el, ctx.anchor, { side, align: 'center' });
  }

  $effect(() => {
    if (ctx?.open && el) {
      void updatePosition();
      const onReposition = () => void updatePosition();
      window.addEventListener('scroll', onReposition, true);
      window.addEventListener('resize', onReposition);
      return () => {
        window.removeEventListener('scroll', onReposition, true);
        window.removeEventListener('resize', onReposition);
      };
    }
  });
</script>

{#if ctx?.open}
  <span
    bind:this={el}
    use:portal
    class="ohjanus-tooltip-content ohjanus-tooltip-{side} {className}"
    style="position: fixed; top: 0; left: 0; {style}"
    role="tooltip"
    data-portal="true"
  >
    {@render children?.()}
  </span>
{/if}
