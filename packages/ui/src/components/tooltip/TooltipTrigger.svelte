<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { class: className = '', style = '', children }: Props = $props();
  const ctx = getContext<{ scheduleOpen: () => void; scheduleClose: () => void }>('ohjanus-tooltip');
</script>

<span
  class="ohjanus-tooltip-trigger {className}"
  {style}
  role="button"
  tabindex="0"
  onmouseenter={() => ctx?.scheduleOpen()}
  onmouseleave={() => ctx?.scheduleClose()}
  onfocusin={() => ctx?.scheduleOpen()}
  onfocusout={() => ctx?.scheduleClose()}
>
  {@render children?.()}
</span>
