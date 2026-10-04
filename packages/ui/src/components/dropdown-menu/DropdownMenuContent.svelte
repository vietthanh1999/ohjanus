<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext, tick, onDestroy } from 'svelte';
  import { portal, positionFloating } from '../../internal/floating.js';

  interface Props {
    align?: 'start' | 'center' | 'end';
    side?: 'top' | 'bottom' | 'left' | 'right';
    class?: string;
    style?: string;
    children?: Snippet;
  }

  interface MenuContext {
    open: boolean;
    anchor?: HTMLElement;
    setContentEl?: (el: HTMLElement | undefined) => void;
  }

  let { align = 'start', side = 'bottom', class: className = '', style = '', children }: Props = $props();
  const ctx = getContext<MenuContext>('ohjanus-menu');
  let el: HTMLElement | undefined = $state();

  async function updatePosition() {
    await tick();
    if (el && ctx?.anchor) positionFloating(el, ctx.anchor, { side, align });
  }

  $effect(() => {
    if (el) {
      ctx?.setContentEl?.(el);
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

  onDestroy(() => ctx?.setContentEl?.(undefined));
</script>

{#if ctx?.open}
  <div
    bind:this={el}
    use:portal
    class="ohjanus-dropdown-content ohjanus-dropdown-align-{align} {className}"
    style="position: fixed; top: 0; left: 0; {style}"
    role="menu"
    data-side={side}
    data-portal="true"
  >
    {@render children?.()}
  </div>
{/if}
