<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext } from 'svelte';

  interface AccordionContext {
    toggle: (value: string) => void;
    isOpen: (value: string) => boolean;
    type: string;
  }

  interface Props {
    itemValue: string;
    disabled?: boolean;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { itemValue, disabled = false, class: className = '', style = '', children }: Props = $props();

  const ctx = getContext<AccordionContext>('ohjanus-accordion');
  const open = $derived(ctx ? ctx.isOpen(itemValue) : false);
</script>

<button
  type="button"
  class="ohjanus-accordion-trigger {className}"
  {style}
  {disabled}
  aria-expanded={open}
  onclick={() => ctx?.toggle(itemValue)}
>
  <span class="ohjanus-accordion-trigger-label">{@render children?.()}</span>
  <span class="ohjanus-accordion-chevron" aria-hidden="true" data-open={open}>▾</span>
</button>
