<script lang="ts">
  import type { Snippet } from 'svelte';
  import Flex from './Flex.svelte';

  interface Props {
    as?: 'div' | 'section' | 'article' | 'aside' | 'main' | 'span';
    inline?: boolean;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let {
    as = 'div',
    inline = false,
    class: className = '',
    style = '',
    children
  }: Props = $props();
</script>

<Flex
  {as}
  {inline}
  align="center"
  justify="center"
  class="ohjanus-center {className}"
  {style}
>
  {@render children?.()}
</Flex>
