<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    label?: string;
    description?: string;
    error?: string;
    required?: boolean;
    class?: string;
    children?: Snippet;
  }

  let {
    label,
    description,
    error,
    required = false,
    class: className = '',
    children
  }: Props = $props();
</script>

<div class="ohjanus-field {className}">
  {#if label}
    <label class="ohjanus-field-label">
      {label}
      {#if required}
        <span class="ohjanus-field-required">*</span>
      {/if}
    </label>
  {/if}

  {#if description}
    <p class="ohjanus-field-description">{description}</p>
  {/if}

  <div class="ohjanus-field-control">
    {@render children?.()}
  </div>

  {#if error}
    <p class="ohjanus-field-error" role="alert">{error}</p>
  {/if}
</div>
