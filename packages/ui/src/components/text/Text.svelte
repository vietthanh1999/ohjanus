<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    as?: 'span' | 'p' | 'div' | 'label' | 'code';
    size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl';
    weight?: 'normal' | 'medium' | 'semibold' | 'bold';
    color?: 'default' | 'muted' | 'secondary' | 'accent' | 'success' | 'warning' | 'danger';
    mono?: boolean;
    truncate?: boolean;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let {
    as = 'span',
    size = 'sm',
    weight = 'normal',
    color = 'default',
    mono = false,
    truncate = false,
    class: className = '',
    style = '',
    children
  }: Props = $props();

  const sizeClass = $derived(`ohjanus-text-${size}`);
  const weightClass = $derived(`ohjanus-weight-${weight}`);
  const colorClass = $derived(`ohjanus-color-${color}`);
  const monoClass = $derived(mono ? 'ohjanus-mono' : '');
  const truncateClass = $derived(truncate ? 'ohjanus-truncate' : '');
</script>

<svelte:element
  this={as}
  class="ohjanus-text {sizeClass} {weightClass} {colorClass} {monoClass} {truncateClass} {className}"
  {style}
>
  {@render children?.()}
</svelte:element>
