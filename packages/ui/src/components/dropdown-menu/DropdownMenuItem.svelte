<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext } from 'svelte';

  interface Props {
    disabled?: boolean;
    onclick?: (e: MouseEvent) => void;
    class?: string;
    style?: string;
    children?: Snippet;
    [key: string]: any;
  }

  let { disabled = false, onclick, class: className = '', style = '', children, ...restProps }: Props = $props();
  const ctx = getContext<{ close: () => void }>('ohjanus-menu');

  function handleClick(e: MouseEvent) {
    if (disabled) return;
    onclick?.(e);
    if (!e.defaultPrevented) ctx?.close();
  }
</script>

<button
  type="button"
  role="menuitem"
  {disabled}
  class="ohjanus-dropdown-item {className}"
  {style}
  onclick={handleClick}
  {...restProps}
>
  {@render children?.()}
</button>
