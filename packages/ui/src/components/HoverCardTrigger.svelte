<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { class: className = '', style = '', children }: Props = $props();
  const ctx = getContext<{ scheduleOpen: () => void }>('ohjanus-hovercard');
</script>

<span class="ohjanus-hovercard-trigger {className}" {style} role="button" tabindex="0" onmouseenter={() => ctx?.scheduleOpen()} onfocusin={() => ctx?.scheduleOpen()}>
  {@render children?.()}
</span>
