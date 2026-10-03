<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext } from 'svelte';

  interface Props {
    value: string;
    disabled?: boolean;
    class?: string;
    style?: string;
    ariaLabel?: string;
    children?: Snippet;
  }

  let { value: itemValue, disabled = false, class: className = '', style = '', ariaLabel, children }: Props = $props();
  const ctx = getContext<{
    toggle: (v: string) => void;
    isPressed: (v: string) => boolean;
    disabled: boolean;
    variant: string;
    size: string;
  }>('ohjanus-toggle-group');

  const pressed = $derived(ctx?.isPressed(itemValue) ?? false);
  const isDisabled = $derived(disabled || ctx?.disabled);
</script>

<button
  type="button"
  class="ohjanus-toggle ohjanus-toggle-{ctx?.variant ?? 'default'} ohjanus-toggle-{ctx?.size ?? 'md'} {className}"
  class:pressed
  {style}
  disabled={isDisabled}
  aria-pressed={pressed}
  aria-label={ariaLabel}
  onclick={() => ctx?.toggle(itemValue)}
>
  {@render children?.()}
</button>
