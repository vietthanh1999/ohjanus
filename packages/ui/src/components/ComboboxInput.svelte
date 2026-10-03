<script lang="ts">
  import { getContext } from 'svelte';

  interface Props {
    value?: string;
    placeholder?: string;
    class?: string;
    style?: string;
    disabled?: boolean;
  }

  let {
    value = $bindable(''),
    placeholder = 'Search...',
    class: className = '',
    style = '',
    disabled = false
  }: Props = $props();

  const ctx = getContext<{ open: boolean; setOpen: (v: boolean) => void; setSearch: (v: string) => void; search: string }>('ohjanus-combobox');

  function handleInput(e: Event) {
    const v = (e.target as HTMLInputElement).value;
    value = v;
    ctx?.setSearch(v);
    ctx?.setOpen(true);
  }
</script>

<div class="ohjanus-combobox-input-wrapper {className}" {style}>
  <input
    type="text"
    class="ohjanus-combobox-input"
    {placeholder}
    {disabled}
    value={value}
    oninput={handleInput}
    onfocus={() => ctx?.setOpen(true)}
    role="combobox"
    aria-expanded={ctx?.open ?? false}
    aria-controls="ohjanus-combobox-listbox"
    aria-autocomplete="list"
  />
</div>
