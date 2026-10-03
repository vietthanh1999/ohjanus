<script lang="ts">
  import type { Snippet } from 'svelte';
  import { buttonVariants } from './button.variants.js';

  interface Props {
    variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger' | 'link';
    size?: 'xs' | 'sm' | 'md' | 'lg';
    disabled?: boolean;
    loading?: boolean;
    type?: 'button' | 'submit' | 'reset';
    class?: string;
    onclick?: (e: MouseEvent) => void;
    children?: Snippet;
    [key: string]: any;
  }

  let {
    variant = 'secondary',
    size = 'sm',
    disabled = false,
    loading = false,
    type = 'button',
    class: className = '',
    onclick,
    children,
    ...restProps
  }: Props = $props();

  let classes = $derived(buttonVariants({ variant, size, class: className }));
</script>

<button
  {type}
  disabled={disabled || loading}
  class={classes}
  {onclick}
  {...restProps}
>
  {#if loading}
    <span class="ohjanus-spinner" aria-hidden="true"></span>
  {/if}
  {@render children?.()}
</button>
