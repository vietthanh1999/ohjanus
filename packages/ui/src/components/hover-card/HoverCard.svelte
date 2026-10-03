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

  let { openDelay = 200, closeDelay = 150, class: className = '', style = '', children }: Props = $props();

  let open = $state(false);
  let timer: ReturnType<typeof setTimeout> | undefined;

  function scheduleOpen() {
    clearTimeout(timer);
    timer = setTimeout(() => (open = true), openDelay);
  }

  function scheduleClose() {
    clearTimeout(timer);
    timer = setTimeout(() => (open = false), closeDelay);
  }

  setContext('ohjanus-hovercard', {
    get open() { return open; },
    scheduleOpen,
    scheduleClose
  });
</script>

<div
  class="ohjanus-hovercard {className}"
  {style}
  data-open={open}
  role="group"
  onmouseenter={scheduleOpen}
  onmouseleave={scheduleClose}
  onfocusin={scheduleOpen}
  onfocusout={scheduleClose}
>
  {@render children?.()}
</div>
