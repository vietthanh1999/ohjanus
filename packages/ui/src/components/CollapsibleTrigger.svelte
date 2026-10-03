<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { class: className = '', style = '', children }: Props = $props();
  const ctx = getContext<{ open: boolean; toggle: () => void }>('ohjanus-collapsible');
</script>

<button
  type="button"
  class="ohjanus-collapsible-trigger {className}"
  {style}
  aria-expanded={ctx?.open ?? false}
  onclick={() => ctx?.toggle()}
>
  {@render children?.()}
</button>
