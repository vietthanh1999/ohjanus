<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    variant?: 'default' | 'destructive';
    onclick?: (e: MouseEvent) => void;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { variant = 'default', onclick, class: className = '', style = '', children }: Props = $props();
</script>

<button
  type="button"
  class="ohjanus-btn {variant === 'destructive' ? 'ohjanus-btn-danger' : 'ohjanus-btn-primary'} ohjanus-btn-sm {className}"
  {style}
  {onclick}
>
  {@render children?.()}
</button>
