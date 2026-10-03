<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { class: className = '', style = '', children }: Props = $props();
  const ctx = getContext<{
    open: boolean;
    disabled: boolean;
    setOpen: (v: boolean) => void;
  }>('ohjanus-select');
</script>

<button
  type="button"
  class="ohjanus-select-trigger {className}"
  {style}
  disabled={ctx?.disabled}
  onclick={() => ctx?.setOpen(!ctx.open)}
  aria-haspopup="listbox"
  aria-expanded={ctx?.open ?? false}
>
  {#if children}
    {@render children()}
  {:else}
    <span class="ohjanus-select-value-empty">Select...</span>
  {/if}
  <span class="ohjanus-select-caret" aria-hidden="true">▾</span>
</button>
