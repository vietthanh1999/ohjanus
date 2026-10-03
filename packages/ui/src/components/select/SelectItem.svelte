<script lang="ts">
  import type { Snippet } from 'svelte';
  import { Icon } from '@ohjanus/icons';
  import { getContext, onMount } from 'svelte';

  interface Props {
    value: string;
    label?: string;
    disabled?: boolean;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { value: itemValue, label, disabled = false, class: className = '', style = '', children }: Props = $props();
  const ctx = getContext<{
    multiple: boolean;
    disabled: boolean;
    select: (v: string) => void;
    isSelected: (v: string) => boolean;
    registerLabel: (v: string, label: string) => void;
  }>('ohjanus-select');

  onMount(() => {
    ctx?.registerLabel(itemValue, label ?? itemValue);
  });

  const selected = $derived(ctx?.isSelected(itemValue) ?? false);
  const isDisabled = $derived(disabled || ctx?.disabled);
</script>

<button
  type="button"
  role="option"
  aria-selected={selected}
  disabled={isDisabled}
  class="ohjanus-select-option {className}"
  class:selected
  {style}
  onclick={() => ctx?.select(itemValue)}
>
  {#if ctx?.multiple}
    {#if selected}
      <span aria-hidden="true"><Icon name="check-square" size={12} /></span>
    {:else}
      <span class="ohjanus-select-check-empty" aria-hidden="true"></span>
    {/if}
  {/if}
  <span class="ohjanus-select-option-label">
    {#if children}
      {@render children()}
    {:else}
      {label ?? itemValue}
    {/if}
  </span>
  {#if !ctx?.multiple && selected}
    <span class="ohjanus-select-check" aria-hidden="true"><Icon name="check" size={12} /></span>
  {/if}
</button>
