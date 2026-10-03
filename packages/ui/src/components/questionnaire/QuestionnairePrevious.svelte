<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext } from 'svelte';

  interface Props {
    disabled?: boolean;
    class?: string;
    style?: string;
    children?: Snippet;
    onclick?: (e: MouseEvent) => void;
  }

  let { disabled = false, class: className = '', style = '', children, onclick }: Props = $props();
  const ctx = getContext<{ step: number; previous: () => void }>('ohjanus-questionnaire');
  const atFirst = $derived(!ctx || ctx.step === 0);
</script>

<button
  type="button"
  class="ohjanus-btn ohjanus-btn-outline ohjanus-btn-sm {className}"
  {style}
  disabled={disabled || atFirst}
  onclick={(e) => { onclick?.(e); if (!e.defaultPrevented) ctx?.previous(); }}
>
  {#if children}{@render children()}{:else}Back{/if}
</button>
