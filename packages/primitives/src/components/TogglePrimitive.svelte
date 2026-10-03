<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    checked?: boolean;
    disabled?: boolean;
    onchange?: (checked: boolean) => void;
    children?: Snippet<[{ checked: boolean; disabled: boolean }]>;
    class?: string;
    [key: string]: any;
  }

  let {
    checked = false,
    disabled = false,
    onchange,
    children,
    class: className = '',
    ...restProps
  }: Props = $props();

  function toggle() {
    if (disabled) return;
    const next = !checked;
    onchange?.(next);
  }

  function handleKeydown(e: KeyboardEvent) {
    if (disabled) return;
    if (e.key === ' ' || e.key === 'Enter') {
      e.preventDefault();
      toggle();
    }
  }
</script>

<button
  type="button"
  role="switch"
  aria-checked={checked}
  {disabled}
  class="ohjanus-toggle-primitive {className}"
  onclick={toggle}
  onkeydown={handleKeydown}
  {...restProps}
>
  {#if children}
    {@render children({ checked, disabled })}
  {/if}
</button>
