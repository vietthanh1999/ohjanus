<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext } from 'svelte';

  interface Props {
    value: string;
    disabled?: boolean;
    class?: string;
    style?: string;
    children?: Snippet;
    onSelect?: (value: string) => void;
  }

  let { value: itemValue, disabled = false, class: className = '', style = '', children, onSelect }: Props = $props();
  const ctx = getContext<{ setOpen: (v: boolean) => void; setSearch: (v: string) => void } | undefined>('ohjanus-combobox');

  function select() {
    if (disabled) return;
    onSelect?.(itemValue);
    ctx?.setSearch('');
    ctx?.setOpen(false);
  }
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<div
  role="option"
  aria-selected="false"
  aria-disabled={disabled || undefined}
  class="ohjanus-combobox-item {className}"
  class:disabled
  {style}
  tabindex={disabled ? -1 : 0}
  onclick={select}
  onkeydown={(e) => { if (e.key === 'Enter' && !disabled) select(); }}
>
  {@render children?.()}
</div>
