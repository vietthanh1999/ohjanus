<script lang="ts">
  import type { Snippet } from 'svelte';
  import { setContext } from 'svelte';

  interface Props {
    value?: number;
    max?: number;
    class?: string;
    style?: string;
    label?: string;
    children?: Snippet;
  }

  let {
    value = $bindable(0),
    max = 100,
    class: className = '',
    style = '',
    label = 'Progress',
    children
  }: Props = $props();

  const percent = $derived(Math.min(100, Math.max(0, ((value ?? 0) / max) * 100)));

  setContext('ohjanus-progress', {
    get value() { return value; },
    get max() { return max; },
    get percent() { return percent; }
  });
</script>

<div
  class="ohjanus-progress {className}"
  {style}
  role="progressbar"
  aria-label={label}
  aria-valuemin={0}
  aria-valuemax={max}
  aria-valuenow={Math.round(value ?? 0)}
>
  {@render children?.()}
</div>
