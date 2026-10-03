<script lang="ts">
  import type { Snippet } from 'svelte';
  import { setContext } from 'svelte';

  interface Props {
    value?: string;
    class?: string;
    style?: string;
    children?: Snippet;
    onValueChange?: (value: string) => void;
  }

  let { value = $bindable(''), class: className = '', style = '', children, onValueChange }: Props = $props();

  function select(v: string) {
    value = v;
    onValueChange?.(v);
  }

  setContext('ohjanus-tabs', {
    get value() { return value; },
    select
  });
</script>

<div class="ohjanus-tabs {className}" {style} data-value={value}>
  {@render children?.()}
</div>
