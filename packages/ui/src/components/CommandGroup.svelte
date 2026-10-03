<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    heading?: string;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { heading, class: className = '', style = '', children }: Props = $props();
</script>

<div class="ohjanus-command-group {className}" {style}>
  {#if heading}
    <div class="ohjanus-command-group-heading">{heading}</div>
  {/if}
  {@render children?.()}
</div>
