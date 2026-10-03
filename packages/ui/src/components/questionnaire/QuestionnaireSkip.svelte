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
  const ctx = getContext<{ skip: () => void }>('ohjanus-questionnaire');
</script>

<button
  type="button"
  class="ohjanus-btn ohjanus-btn-ghost ohjanus-btn-sm {className}"
  {style}
  {disabled}
  onclick={(e) => { onclick?.(e); if (!e.defaultPrevented) ctx?.skip(); }}
>
  {#if children}{@render children()}{:else}Skip{/if}
</button>
