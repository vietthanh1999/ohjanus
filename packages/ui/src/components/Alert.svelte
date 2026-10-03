<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    variant?: 'info' | 'success' | 'warning' | 'danger';
    title?: string;
    class?: string;
    children?: Snippet;
  }

  let {
    variant = 'info',
    title,
    class: className = '',
    children
  }: Props = $props();
</script>

<div class="ohjanus-alert ohjanus-alert-{variant} {className}" role="alert">
  {#if title}
    <div class="ohjanus-alert-title">{title}</div>
  {/if}
  <div class="ohjanus-alert-content">
    {@render children?.()}
  </div>
</div>
