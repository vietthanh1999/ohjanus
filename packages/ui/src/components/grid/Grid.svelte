<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    as?: 'div' | 'section' | 'article' | 'aside' | 'main' | 'header' | 'footer' | 'nav';
    columns?: number | string;
    rows?: number | string;
    gap?: string | number;
    rowGap?: string | number;
    columnGap?: string | number;
    align?: 'start' | 'center' | 'end' | 'stretch';
    justify?: 'start' | 'center' | 'end' | 'stretch' | 'space-between';
    class?: string;
    style?: string;
    role?: string;
    onclick?: (e: MouseEvent) => void;
    children?: Snippet;
  }

  let {
    as = 'div',
    columns = 1,
    rows,
    gap = '12px',
    rowGap,
    columnGap,
    align = 'stretch',
    justify = 'stretch',
    class: className = '',
    style = '',
    role,
    onclick,
    children
  }: Props = $props();

  const colsValue = $derived(
    typeof columns === 'number' ? `repeat(${columns}, 1fr)` : columns
  );

  const rowsValue = $derived(
    rows !== undefined
      ? typeof rows === 'number'
        ? `repeat(${rows}, 1fr)`
        : rows
      : 'none'
  );

  const formatGap = (g: string | number) =>
    typeof g === 'number'
      ? `${g}px`
      : g.endsWith('px') || g.endsWith('rem') || g.endsWith('em') || g.endsWith('%')
      ? g
      : `var(--spacing-${g}, ${g}px)`;

  const gridStyle = $derived(
    `display: grid; ` +
    `grid-template-columns: ${colsValue}; ` +
    (rows ? `grid-template-rows: ${rowsValue}; ` : '') +
    `gap: ${formatGap(gap)}; ` +
    (rowGap ? `row-gap: ${formatGap(rowGap)}; ` : '') +
    (columnGap ? `column-gap: ${formatGap(columnGap)}; ` : '') +
    `align-items: ${align}; ` +
    `justify-items: ${justify}; ` +
    `${style}`
  );
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<svelte:element
  this={as}
  class="ohjanus-grid {className}"
  style={gridStyle}
  {role}
  {onclick}
>
  {@render children?.()}
</svelte:element>
