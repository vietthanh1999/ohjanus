<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    variant?: 'default' | 'success' | 'warning' | 'danger' | 'info' | 'counter' | 'license';
    size?: 'sm' | 'md';
    class?: string;
    children?: Snippet;
  }

  let {
    variant = 'default',
    size = 'sm',
    class: className = '',
    children
  }: Props = $props();
</script>

<span class="ohjanus-badge ohjanus-badge-{variant} ohjanus-badge-{size} {className}">
  {@render children?.()}
</span>
