<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    class?: string;
    children?: Snippet;
    header?: Snippet;
    footer?: Snippet;
  }

  let { class: className = '', children, header, footer }: Props = $props();
</script>

<div class="ohjanus-card {className}">
  {#if header}
    <div class="ohjanus-card-header">
      {@render header()}
    </div>
  {/if}
  <div class="ohjanus-card-body">
    {@render children?.()}
  </div>
  {#if footer}
    <div class="ohjanus-card-footer">
      {@render footer()}
    </div>
  {/if}
</div>
