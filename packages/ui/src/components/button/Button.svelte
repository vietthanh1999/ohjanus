<script lang="ts">
  import type { Snippet } from 'svelte';
  import { buttonVariants } from './button.variants.js';

  interface Props {
    variant?: 'default' | 'primary' | 'secondary' | 'outline' | 'ghost' | 'destructive' | 'danger' | 'link';
    size?: 'default' | 'xs' | 'sm' | 'md' | 'lg' | 'xl' | 'icon' | 'icon-xs' | 'icon-sm' | 'icon-lg';
    disabled?: boolean;
    loading?: boolean;
    fullWidth?: boolean;
    type?: 'button' | 'submit' | 'reset';
    href?: string;
    iconOnly?: boolean;
    class?: string;
    onclick?: (e: MouseEvent) => void;
    children?: Snippet;
    iconLeft?: Snippet;
    iconRight?: Snippet;
    [key: string]: any;
  }

  let {
    variant = 'default',
    size = 'default',
    disabled = false,
    loading = false,
    fullWidth = false,
    type = 'button',
    href,
    iconOnly = false,
    class: className = '',
    onclick,
    children,
    iconLeft,
    iconRight,
    ...restProps
  }: Props = $props();

  let classes = $derived(
    buttonVariants({ variant, size, class: className }) + (fullWidth ? ' ohjanus-btn-fullwidth' : '')
  );
  let isDisabled = $derived(disabled || loading);
  // `href` renders an anchor; a disabled link never navigates.
  let asLink = $derived(href != null && href !== '' && !isDisabled);
</script>

{#if asLink}
  <a
    href={href}
    class={classes}
    aria-disabled={isDisabled}
    aria-busy={loading || undefined}
    aria-label={iconOnly ? restProps['aria-label'] : undefined}
    onclick={onclick}
    {...restProps}
  >
    {#if iconLeft}{@render iconLeft()}{/if}
    {#if loading}
      <span class="ohjanus-spinner" aria-hidden="true"></span>
    {/if}
    {@render children?.()}
    {#if iconRight}{@render iconRight()}{/if}
  </a>
{:else}
  <button
    {type}
    disabled={isDisabled}
    class={classes}
    aria-busy={loading || undefined}
    {onclick}
    {...restProps}
  >
    {#if iconLeft}{@render iconLeft()}{/if}
    {#if loading}
      <span class="ohjanus-spinner" aria-hidden="true"></span>
    {/if}
    {@render children?.()}
    {#if iconRight}{@render iconRight()}{/if}
  </button>
{/if}
