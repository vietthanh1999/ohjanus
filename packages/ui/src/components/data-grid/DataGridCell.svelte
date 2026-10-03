<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    /** Active (keyboard / click) cell highlight. */
    focused?: boolean;
    /** NULL value styling: muted + italic. */
    isNull?: boolean;
    /** Single-line ellipsis overflow. */
    truncate?: boolean;
    align?: 'left' | 'center' | 'right';
    title?: string;
    onclick?: (e: MouseEvent) => void;
    children?: Snippet;
  }

  let {
    class: className = '',
    style = '',
    focused = false,
    isNull = false,
    truncate = false,
    align = 'left',
    title,
    onclick,
    children
  }: Props = $props();
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
<td
  class="ohjanus-data-grid-cell {focused ? 'focused' : ''} {isNull ? 'null' : ''} {truncate ? 'truncate' : ''} ohjanus-data-grid-align-{align} {className}"
  {style}
  {title}
  {onclick}
>
  {@render children?.()}
</td>
