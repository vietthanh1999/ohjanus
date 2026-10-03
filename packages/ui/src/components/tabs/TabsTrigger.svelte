<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext } from 'svelte';

  interface Props {
    value: string;
    disabled?: boolean;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { value: triggerValue, disabled = false, class: className = '', style = '', children }: Props = $props();
  const ctx = getContext<{ value: string; select: (v: string) => void }>('ohjanus-tabs');
  const active = $derived(ctx?.value === triggerValue);
</script>

<button
  type="button"
  role="tab"
  aria-selected={active}
  aria-controls="ohjanus-tabpanel-{triggerValue}"
  id="ohjanus-tab-{triggerValue}"
  {disabled}
  class="ohjanus-tabs-trigger {className}"
  class:active
  {style}
  onclick={() => ctx?.select(triggerValue)}
>
  {@render children?.()}
</button>
