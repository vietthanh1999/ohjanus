<script lang="ts">
  import type { Snippet } from 'svelte';
  import { setContext } from 'svelte';

  interface Props {
    type?: 'single' | 'multiple';
    collapsible?: boolean;
    value?: string | string[];
    class?: string;
    style?: string;
    children?: Snippet;
    onValueChange?: (value: string | string[]) => void;
  }

  let {
    type = 'single',
    collapsible = true,
    value = $bindable(type === 'multiple' ? [] : ''),
    class: className = '',
    style = '',
    children,
    onValueChange
  }: Props = $props();

  function toggle(itemValue: string) {
    if (type === 'multiple') {
      const arr = Array.isArray(value) ? value : [];
      const next = arr.includes(itemValue) ? arr.filter((v) => v !== itemValue) : [...arr, itemValue];
      value = next;
      onValueChange?.(next);
    } else {
      const next = value === itemValue ? (collapsible ? '' : itemValue) : itemValue;
      value = next;
      onValueChange?.(next);
    }
  }

  function isOpen(itemValue: string): boolean {
    return Array.isArray(value) ? value.includes(itemValue) : value === itemValue;
  }

  setContext('ohjanus-accordion', {
    toggle,
    isOpen,
    get type() { return type; }
  });
</script>

<div class="ohjanus-accordion {className}" {style} data-type={type}>
  {@render children?.()}
</div>
