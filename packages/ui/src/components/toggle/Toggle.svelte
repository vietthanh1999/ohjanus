<script lang="ts">
  import type { Snippet } from 'svelte';
  import { TogglePrimitive } from '@ohjanus/primitives';

  interface Props {
    pressed?: boolean;
    disabled?: boolean;
    variant?: 'default' | 'outline';
    size?: 'sm' | 'md' | 'lg';
    class?: string;
    style?: string;
    ariaLabel?: string;
    onchange?: (pressed: boolean) => void;
    children?: Snippet;
  }

  let {
    pressed = $bindable(false),
    disabled = false,
    variant = 'default',
    size = 'md',
    class: className = '',
    style = '',
    ariaLabel,
    onchange,
    children: content
  }: Props = $props();

  function handleChange(next: boolean) {
    pressed = next;
    onchange?.(next);
  }
</script>

<TogglePrimitive
  checked={pressed}
  {disabled}
  onchange={handleChange}
  class="ohjanus-toggle ohjanus-toggle-{variant} ohjanus-toggle-{size} {className}"
  style={style}
  aria-label={ariaLabel}
  aria-pressed={pressed}
>
  {#snippet children({ checked, disabled: isDisabled })}
    <span class="ohjanus-toggle-content" data-pressed={checked} data-disabled={isDisabled}>
      {@render content?.()}
    </span>
  {/snippet}
</TogglePrimitive>
