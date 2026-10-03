<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    as?: 'div' | 'section' | 'article' | 'aside' | 'main';
    maxWidth?: 'sm' | 'md' | 'lg' | 'xl' | '2xl' | 'full' | string;
    centered?: boolean;
    padding?: string | number;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let {
    as = 'div',
    maxWidth = 'xl',
    centered = true,
    padding = '16px',
    class: className = '',
    style = '',
    children
  }: Props = $props();

  const widthMap: Record<string, string> = {
    sm: '640px',
    md: '768px',
    lg: '1024px',
    xl: '1280px',
    '2xl': '1536px',
    full: '100%'
  };

  const padValue = $derived(
    typeof padding === 'number' ? `${padding}px` : padding
  );

  const containerStyle = $derived(
    `width: 100%; ` +
    `max-width: ${widthMap[maxWidth] || maxWidth}; ` +
    (centered ? `margin-left: auto; margin-right: auto; ` : '') +
    `padding-left: ${padValue}; padding-right: ${padValue}; ` +
    `${style}`
  );
</script>

<svelte:element
  this={as}
  class="ohjanus-container {className}"
  style={containerStyle}
>
  {@render children?.()}
</svelte:element>
