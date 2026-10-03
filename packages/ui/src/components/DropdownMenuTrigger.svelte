<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { class: className = '', style = '', children }: Props = $props();
  const ctx = getContext<{ toggle: () => void; open: boolean }>('ohjanus-menu');
</script>

<button
  type="button"
  class="ohjanus-dropdown-trigger {className}"
  {style}
  aria-haspopup="menu"
  aria-expanded={ctx?.open ?? false}
  onclick={() => ctx?.toggle()}
>
  {@render children?.()}
</button>
