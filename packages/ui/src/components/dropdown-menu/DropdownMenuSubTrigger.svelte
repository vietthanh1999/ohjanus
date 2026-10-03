<script lang="ts">
  import type { Snippet } from 'svelte';
  import { Icon } from '@ohjanus/icons';
  import { getContext } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { class: className = '', style = '', children }: Props = $props();
  const ctx = getContext<{ toggle: () => void }>('ohjanus-menu-sub');
</script>

<button type="button" class="ohjanus-dropdown-sub-trigger {className}" {style} onclick={() => ctx?.toggle()} aria-haspopup="menu">
  {@render children?.()}
  <span aria-hidden="true"><Icon name="chevron-right" size={12} /></span>
</button>
