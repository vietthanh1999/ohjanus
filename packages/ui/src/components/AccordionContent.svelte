<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext } from 'svelte';

  interface AccordionContext {
    toggle: (value: string) => void;
    isOpen: (value: string) => boolean;
  }

  interface Props {
    itemValue: string;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { itemValue, class: className = '', style = '', children }: Props = $props();

  const ctx = getContext<AccordionContext>('ohjanus-accordion');
  const open = $derived(ctx ? ctx.isOpen(itemValue) : false);
</script>

{#if open}
  <div class="ohjanus-accordion-content {className}" {style} role="region">
    {@render children?.()}
  </div>
{/if}
