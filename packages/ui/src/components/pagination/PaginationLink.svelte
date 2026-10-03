<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    href?: string;
    active?: boolean;
    disabled?: boolean;
    onclick?: (e: MouseEvent) => void;
    class?: string;
    style?: string;
    children?: Snippet;
    [key: string]: any;
  }

  let {
    href,
    active = false,
    disabled = false,
    onclick,
    class: className = '',
    style = '',
    children,
    ...restProps
  }: Props = $props();
</script>

{#if href && !disabled}
  <a
    {href}
    class="ohjanus-pagination-link {className}"
    class:active
    {style}
    aria-current={active ? 'page' : undefined}
    {onclick}
    {...restProps}
  >
    {@render children?.()}
  </a>
{:else}
  <button
    type="button"
    class="ohjanus-pagination-link {className}"
    class:active
    {style}
    {disabled}
    aria-current={active ? 'page' : undefined}
    {onclick}
    {...restProps}
  >
    {@render children?.()}
  </button>
{/if}
