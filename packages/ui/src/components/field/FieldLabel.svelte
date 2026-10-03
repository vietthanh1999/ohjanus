<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    forId?: string;
    for?: string;
    required?: boolean;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let {
    forId,
    for: htmlfor,
    required = false,
    class: className = '',
    style = '',
    children
  }: Props = $props();

  const targetId = $derived(forId || htmlfor);
</script>

<label class="ohjanus-field-label {className}" for={targetId} {style}>
  {@render children?.()}
  {#if required}
    <span class="ohjanus-field-required">*</span>
  {/if}
</label>
