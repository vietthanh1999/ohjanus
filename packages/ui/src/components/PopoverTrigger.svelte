<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { class: className = '', style = '', children }: Props = $props();
  const ctx = getContext<{ toggle: () => void; open: boolean }>('ohjanus-popover');
</script>

<button
  type="button"
  class="ohjanus-popover-trigger {className}"
  {style}
  aria-expanded={ctx?.open ?? false}
  aria-haspopup="dialog"
  onclick={() => ctx?.toggle()}
>
  {@render children?.()}
</button>
