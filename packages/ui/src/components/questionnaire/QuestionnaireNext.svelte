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
  const ctx = getContext<{ next: () => void }>('ohjanus-questionnaire');
</script>

<button
  type="button"
  class="ohjanus-btn ohjanus-btn-primary ohjanus-btn-sm {className}"
  {style}
  {disabled}
  onclick={(e) => { onclick?.(e); if (!e.defaultPrevented) ctx?.next(); }}
>
  {#if children}{@render children()}{:else}Next{/if}
</button>
