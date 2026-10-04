<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    as?: 'div' | 'section' | 'article' | 'aside' | 'main' | 'header' | 'footer' | 'nav' | 'span';
    class?: string;
    style?: string;
    role?: string;
    onclick?: (e: MouseEvent) => void;
    children?: Snippet;
    [key: string]: any;
  }

  let {
    as = 'div',
    class: className = '',
    style = '',
    role,
    onclick,
    children,
    ...restProps
  }: Props = $props();
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<svelte:element
  this={as}
  class="ohjanus-box {className}"
  {style}
  {role}
  {onclick}
  {...restProps}
>
  {@render children?.()}
</svelte:element>
