<script lang="ts">
  import type { Snippet } from 'svelte';
  import Flex from '../flex/Flex.svelte';

  interface Props {
    as?: 'div' | 'section' | 'article' | 'aside' | 'main' | 'header' | 'footer' | 'nav' | 'span';
    direction?: 'vertical' | 'horizontal';
    align?: 'start' | 'center' | 'end' | 'stretch' | 'baseline';
    justify?: 'start' | 'center' | 'end' | 'between' | 'around' | 'evenly';
    gap?: string | number;
    wrap?: boolean | 'wrap' | 'nowrap';
    class?: string;
    style?: string;
    role?: string;
    onclick?: (e: MouseEvent) => void;
    children?: Snippet;
  }

  let {
    as = 'div',
    direction = 'vertical',
    align = 'stretch',
    justify = 'start',
    gap = '12px',
    wrap = false,
    class: className = '',
    style = '',
    role,
    onclick,
    children
  }: Props = $props();
</script>

<Flex
  {as}
  direction={direction === 'vertical' ? 'column' : 'row'}
  {align}
  {justify}
  {gap}
  {wrap}
  class="ohjanus-stack {className}"
  {style}
  {role}
  {onclick}
>
  {@render children?.()}
</Flex>
