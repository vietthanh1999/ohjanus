<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext } from 'svelte';

  interface Props {
    value: string;
    disabled?: boolean;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { value: itemValue, disabled = false, class: className = '', style = '', children }: Props = $props();
  const ctx = getContext<{ value: string; select: (v: string) => void }>('ohjanus-radio-group');
  const checked = $derived(ctx?.value === itemValue);
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
  role="menuitemradio"
  aria-checked={checked}
  class="ohjanus-dropdown-radio-item {className}"
  {style}
  tabindex={disabled ? -1 : 0}
  onclick={() => !disabled && ctx?.select(itemValue)}
  onkeydown={(e) => { if ((e.key === 'Enter' || e.key === ' ') && !disabled) { e.preventDefault(); ctx?.select(itemValue); } }}
>
  <span class="ohjanus-dropdown-radio-dot" class:checked aria-hidden="true"></span>
  {@render children?.()}
</div>
