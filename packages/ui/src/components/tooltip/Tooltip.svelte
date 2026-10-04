<script lang="ts">
  import type { Snippet } from 'svelte';
  import { setContext } from 'svelte';

  interface Props {
    openDelay?: number;
    closeDelay?: number;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { openDelay = 300, closeDelay = 150, class: className = '', style = '', children }: Props = $props();

  let open = $state(false);
  let timer: ReturnType<typeof setTimeout> | undefined;
  let anchor: HTMLElement | undefined = $state();

  function scheduleOpen() {
    clearTimeout(timer);
    timer = setTimeout(() => (open = true), openDelay);
  }

  function scheduleClose() {
    clearTimeout(timer);
    timer = setTimeout(() => (open = false), closeDelay);
  }

  setContext('ohjanus-tooltip', {
    get open() { return open; },
    get anchor() { return anchor; },
    scheduleOpen,
    scheduleClose
  });
</script>

<span
  bind:this={anchor}
  class="ohjanus-tooltip {className}"
  {style}
  data-open={open}
  role="group"
  onmouseenter={scheduleOpen}
  onmouseleave={scheduleClose}
  onfocusin={scheduleOpen}
  onfocusout={scheduleClose}
>
  {@render children?.()}
</span>
