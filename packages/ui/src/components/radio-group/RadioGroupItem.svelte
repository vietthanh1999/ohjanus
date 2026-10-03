<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext } from 'svelte';

  interface Props {
    value: string;
    id?: string;
    disabled?: boolean;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { value: itemValue, id, disabled = false, class: className = '', style = '', children }: Props = $props();
  const ctx = getContext<{
    value: string;
    name?: string;
    disabled: boolean;
    select: (v: string) => void;
  }>('ohjanus-radio-group');

  const checked = $derived(ctx?.value === itemValue);
  const isDisabled = $derived(disabled || ctx?.disabled);
</script>

<label class="ohjanus-radio-item {className}" class:checked class:disabled={isDisabled} {style}>
  <input
    type="radio"
    {id}
    name={ctx?.name}
    value={itemValue}
    {checked}
    disabled={isDisabled}
    class="ohjanus-radio-input"
    onchange={() => ctx?.select(itemValue)}
  />
  <span class="ohjanus-radio-dot" aria-hidden="true"></span>
  {#if children}
    <span class="ohjanus-radio-label">{@render children()}</span>
  {/if}
</label>
