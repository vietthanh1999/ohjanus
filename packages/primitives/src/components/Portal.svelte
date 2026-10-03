<script lang="ts">
  import { onMount, type Snippet } from 'svelte';

  interface Props {
    target?: HTMLElement | string;
    children?: Snippet;
  }

  let { target = 'body', children }: Props = $props();

  let mounted = $state(false);
  let container = $state<HTMLElement | null>(null);

  onMount(() => {
    mounted = true;
    container = typeof target === 'string' ? document.querySelector(target) : target;
  });
</script>

{#if mounted && container}
  <div style="display: contents;">
    {@render children?.()}
  </div>
{/if}
