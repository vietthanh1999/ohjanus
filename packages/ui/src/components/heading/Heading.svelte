<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    level?: 1 | 2 | 3 | 4 | 5 | 6;
    size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl';
    weight?: 'normal' | 'medium' | 'semibold' | 'bold';
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let {
    level = 2,
    size,
    weight = 'semibold',
    class: className = '',
    style = '',
    children
  }: Props = $props();

  const tag = $derived(`h${level}` as 'h1' | 'h2' | 'h3' | 'h4' | 'h5' | 'h6');
  const defaultSize = $derived(
    size ?? (level === 1 ? '2xl' : level === 2 ? 'xl' : level === 3 ? 'lg' : 'md')
  );
</script>

<svelte:element
  this={tag}
  class="ohjanus-heading ohjanus-heading-{defaultSize} ohjanus-weight-{weight} {className}"
  {style}
>
  {@render children?.()}
</svelte:element>
