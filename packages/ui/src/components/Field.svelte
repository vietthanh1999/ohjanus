<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    id?: string;
    label?: string;
    hint?: string;
    description?: string;
    error?: string;
    required?: boolean;
    class?: string;
    children?: Snippet;
  }

  let {
    id,
    label,
    hint,
    description,
    error,
    required = false,
    class: className = '',
    children
  }: Props = $props();

  const descText = $derived(hint ?? description);
</script>

<div class="ohjanus-field {className}">
  {#if label}
    <label class="ohjanus-field-label" for={id}>
      {label}
      {#if required}
        <span class="ohjanus-field-required">*</span>
      {/if}
    </label>
  {/if}

  {#if descText}
    <p class="ohjanus-field-description">{descText}</p>
  {/if}

  <div class="ohjanus-field-control">
    {@render children?.()}
  </div>

  {#if error}
    <p class="ohjanus-field-error" role="alert">{error}</p>
  {/if}
</div>
