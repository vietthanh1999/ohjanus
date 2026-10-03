<script lang="ts">
  import type { Snippet } from 'svelte';
  import { setContext } from 'svelte';

  interface Props {
    type?: 'single' | 'multiple';
    value?: string | string[];
    disabled?: boolean;
    variant?: 'default' | 'outline';
    size?: 'sm' | 'md' | 'lg';
    class?: string;
    style?: string;
    children?: Snippet;
    onValueChange?: (value: string | string[]) => void;
  }

  let {
    type = 'single',
    value = $bindable(type === 'multiple' ? [] : ''),
    disabled = false,
    variant = 'default',
    size = 'md',
    class: className = '',
    style = '',
    children,
    onValueChange
  }: Props = $props();

  function toggle(itemValue: string) {
    if (disabled) return;
    if (type === 'multiple') {
      const arr = Array.isArray(value) ? value : [];
      const next = arr.includes(itemValue) ? arr.filter((v) => v !== itemValue) : [...arr, itemValue];
      value = next;
      onValueChange?.(next);
    } else {
      const next = value === itemValue ? '' : itemValue;
      value = next;
      onValueChange?.(next);
    }
  }

  function isPressed(itemValue: string): boolean {
    return Array.isArray(value) ? value.includes(itemValue) : value === itemValue;
  }

  setContext('ohjanus-toggle-group', {
    toggle,
    isPressed,
    get disabled() { return disabled; },
    get variant() { return variant; },
    get size() { return size; }
  });
</script>

<div class="ohjanus-toggle-group ohjanus-toggle-group-{variant} {className}" {style} role="group" data-type={type}>
  {@render children?.()}
</div>
