<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { class: className = '', style = '', children }: Props = $props();
</script>

<li role="presentation" aria-hidden="true" class="ohjanus-breadcrumb-separator {className}" {style}>
  {#if children}
    {@render children()}
  {:else}
    &gt;
  {/if}
</li>
