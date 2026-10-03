<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext } from 'svelte';

  interface Props {
    align?: 'start' | 'center' | 'end';
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { align = 'start', class: className = '', style = '', children }: Props = $props();
  const ctx = getContext<{ open: boolean; setOpen: (v: boolean) => void }>('ohjanus-popover');
</script>

{#if ctx?.open}
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div
    class="ohjanus-popover-content ohjanus-popover-align-{align} {className}"
    {style}
    role="dialog"
    tabindex="-1"
    onclick={(e) => e.stopPropagation()}
  >
    {@render children?.()}
  </div>
{/if}
