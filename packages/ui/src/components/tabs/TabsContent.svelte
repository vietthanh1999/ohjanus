<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext } from 'svelte';

  interface Props {
    value: string;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { value: contentValue, class: className = '', style = '', children }: Props = $props();
  const ctx = getContext<{ value: string }>('ohjanus-tabs');
  const active = $derived(!ctx || ctx.value === contentValue);
</script>

{#if active}
  <div
    class="ohjanus-tabs-content {className}"
    {style}
    role="tabpanel"
    id="ohjanus-tabpanel-{contentValue}"
    aria-labelledby="ohjanus-tab-{contentValue}"
    tabindex="0"
  >
    {@render children?.()}
  </div>
{/if}
