<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    header?: Snippet;
    footer?: Snippet;
    children?: Snippet;
  }

  let { class: className = '', style = '', header, footer, children }: Props = $props();
</script>

<div class="ohjanus-card {className}" {style}>
  {#if header}
    <div class="ohjanus-card-header">
      {@render header()}
    </div>
  {/if}
  {#if children}
    {@render children()}
  {/if}
  {#if footer}
    <div class="ohjanus-card-footer">
      {@render footer()}
    </div>
  {/if}
</div>
