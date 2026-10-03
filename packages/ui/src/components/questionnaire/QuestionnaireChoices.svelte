<script lang="ts">
  import type { Snippet } from 'svelte';
  import { setContext } from 'svelte';

  interface Props {
    value?: string;
    name?: string;
    class?: string;
    style?: string;
    children?: Snippet;
    onValueChange?: (value: string) => void;
  }

  let {
    value = $bindable(''),
    name,
    class: className = '',
    style = '',
    children,
    onValueChange
  }: Props = $props();

  function select(v: string) {
    value = v;
    onValueChange?.(v);
  }

  setContext('ohjanus-questionnaire-choices', {
    get value() { return value; },
    get name() { return name; },
    select
  });
</script>

<div class="ohjanus-questionnaire-choices {className}" {style} role="radiogroup">
  {@render children?.()}
</div>
