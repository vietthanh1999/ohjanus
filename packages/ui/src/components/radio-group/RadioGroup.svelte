<script lang="ts">
  import type { Snippet } from 'svelte';
  import { setContext } from 'svelte';

  interface Props {
    value?: string;
    name?: string;
    disabled?: boolean;
    class?: string;
    style?: string;
    children?: Snippet;
    onValueChange?: (value: string) => void;
  }

  let {
    value = $bindable(''),
    name,
    disabled = false,
    class: className = '',
    style = '',
    children,
    onValueChange
  }: Props = $props();

  function select(v: string) {
    if (disabled) return;
    value = v;
    onValueChange?.(v);
  }

  setContext('ohjanus-radio-group', {
    get value() { return value; },
    get name() { return name; },
    get disabled() { return disabled; },
    select
  });
</script>

<div class="ohjanus-radio-group {className}" {style} role="radiogroup" aria-disabled={disabled || undefined}>
  {@render children?.()}
</div>
