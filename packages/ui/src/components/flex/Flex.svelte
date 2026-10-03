<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    as?: 'div' | 'section' | 'article' | 'aside' | 'main' | 'header' | 'footer' | 'nav' | 'span';
    direction?: 'row' | 'column' | 'row-reverse' | 'column-reverse';
    align?: 'start' | 'center' | 'end' | 'stretch' | 'baseline';
    justify?: 'start' | 'center' | 'end' | 'between' | 'around' | 'evenly';
    wrap?: boolean | 'wrap' | 'nowrap' | 'wrap-reverse';
    gap?: string | number;
    inline?: boolean;
    class?: string;
    style?: string;
    role?: string;
    onclick?: (e: MouseEvent) => void;
    children?: Snippet;
  }

  let {
    as = 'div',
    direction,
    align,
    justify,
    wrap,
    gap,
    inline = false,
    class: className = '',
    style = '',
    role,
    onclick,
    children
  }: Props = $props();

  const alignMap = {
    start: 'flex-start',
    center: 'center',
    end: 'flex-end',
    stretch: 'stretch',
    baseline: 'baseline'
  };

  const justifyMap = {
    start: 'flex-start',
    center: 'center',
    end: 'flex-end',
    between: 'space-between',
    around: 'space-around',
    evenly: 'space-evenly'
  };

  const gapValue = $derived(
    gap === undefined
      ? undefined
      : typeof gap === 'number'
      ? `${gap}px`
      : typeof gap === 'string' && (gap.endsWith('px') || gap.endsWith('rem') || gap.endsWith('em') || gap.endsWith('%'))
      ? gap
      : `var(--spacing-${gap}, ${gap}px)`
  );

  const wrapValue = $derived(
    wrap === undefined ? undefined : typeof wrap === 'boolean' ? (wrap ? 'wrap' : 'nowrap') : wrap
  );

  const flexStyle = $derived(
    `display: ${inline ? 'inline-flex' : 'flex'}; ` +
    (direction ? `flex-direction: ${direction}; ` : '') +
    (align ? `align-items: ${alignMap[align] || align}; ` : '') +
    (justify ? `justify-content: ${justifyMap[justify] || justify}; ` : '') +
    (wrapValue ? `flex-wrap: ${wrapValue}; ` : '') +
    (gapValue !== undefined ? `gap: ${gapValue}; ` : '') +
    `${style}`
  );
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<svelte:element
  this={as}
  class="ohjanus-flex {className}"
  style={flexStyle}
  {role}
  {onclick}
>
  {@render children?.()}
</svelte:element>
