<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    htmlFor?: string;
    for?: string;
    required?: boolean;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let {
    htmlFor,
    for: forAttr,
    required = false,
    class: className = '',
    style = '',
    children
  }: Props = $props();

  const target = $derived(htmlFor ?? forAttr);
</script>

<label class="ohjanus-label {className}" for={target} {style}>
  {@render children?.()}
  {#if required}
    <span class="ohjanus-field-required" aria-hidden="true">*</span>
  {/if}
</label>
