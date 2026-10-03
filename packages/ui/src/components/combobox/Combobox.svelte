<script lang="ts">
  import type { Snippet } from 'svelte';
  import { Icon } from '@ohjanus/icons';

  interface Option {
    value: string;
    label: string;
    disabled?: boolean;
    group?: string;
  }

  interface Props {
    value?: string | string[];
    options?: Option[];
    placeholder?: string;
    multiple?: boolean;
    disabled?: boolean;
    clearable?: boolean;
    creatable?: boolean;
    invalid?: boolean;
    size?: 'sm' | 'md' | 'lg';
    class?: string;
    style?: string;
    id?: string;
    filterFn?: (option: Option, search: string) => boolean;
    onchange?: (value: string | string[]) => void;
    onSelect?: (value: string) => void;
    option?: Snippet<[Option]>;
    empty?: Snippet;
    children?: Snippet;
    [key: string]: any;
  }

  let {
    value = $bindable(''),
    options = [],
    placeholder = 'Search...',
    multiple = false,
    disabled = false,
    clearable = false,
    creatable = false,
    invalid = false,
    size = 'md',
    class: className = '',
    style = '',
    id,
    filterFn,
    onchange,
    onSelect,
    option,
    empty,
    children,
    ...restProps
  }: Props = $props();

  let search = $state('');
  let open = $state(false);
  let highlight = $state(0);
  let listId = $derived(`${id ?? 'combobox'}-listbox`);

  const isMulti = $derived(Array.isArray(value));
  const selectedValues: string[] = $derived(isMulti ? (value as string[]) : value ? [value as string] : []);

  function defaultFilter(opt: Option, q: string): boolean {
    return opt.label.toLowerCase().includes(q.toLowerCase());
  }

  let filtered = $derived.by(() => {
    const q = search.trim();
    const fn = filterFn ?? defaultFilter;
    let list = q ? options.filter((o) => fn(o, q)) : options;
    if (creatable && q && !options.some((o) => o.label === q || o.value === q)) {
      list = [...list, { value: q, label: `Create "${q}"` }];
    }
    return list;
  });

  function labelFor(v: string): string {
    return options.find((o) => o.value === v)?.label ?? v;
  }

  function commit(next: string | string[]) {
    value = next;
    onchange?.(next);
  }

  function choose(optValue: string, created = false) {
    onSelect?.(optValue);
    if (isMulti) {
      const current = value as string[];
      commit(current.includes(optValue) ? current.filter((v) => v !== optValue) : [...current, optValue]);
    } else {
      commit(optValue);
      search = '';
      open = false;
    }
    if (created) search = '';
  }

  function clear() {
    commit(isMulti ? [] : '');
    search = '';
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'ArrowDown') {
      e.preventDefault();
      open = true;
      highlight = (highlight + 1) % Math.max(filtered.length, 1);
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      highlight = (highlight - 1 + Math.max(filtered.length, 1)) % Math.max(filtered.length, 1);
    } else if (e.key === 'Enter' && open && filtered[highlight]) {
      e.preventDefault();
      const opt = filtered[highlight];
      choose(opt.value, creatable && !options.some((o) => o.value === opt.value));
    } else if (e.key === 'Escape') {
      open = false;
    }
  }

  function handleBlur(e: FocusEvent) {
    if (!(e.currentTarget as HTMLElement).contains(e.relatedTarget as Node)) {
      open = false;
    }
  }
</script>

<div
  class="ohjanus-combobox ohjanus-combobox-{size} {className}"
  class:ohjanus-combobox-invalid={invalid}
  class:ohjanus-combobox-disabled={disabled}
  style={style}
  onfocusout={handleBlur}
  {...restProps}
>
  <div class="ohjanus-combobox-tags">
    {#if isMulti}
      {#each selectedValues as v (v)}
        <span class="ohjanus-combobox-chip">
          {labelFor(v)}
          {#if !disabled}
            <button type="button" aria-label="Remove {labelFor(v)}" onclick={() => choose(v)}><Icon name="x" size={10} /></button>
          {/if}
        </span>
      {/each}
    {/if}
    <input
      {id}
      type="text"
      role="combobox"
      aria-expanded={open}
      aria-controls={listId}
      aria-autocomplete="list"
      {placeholder}
      {disabled}
      aria-invalid={invalid || undefined}
      bind:value={search}
      onfocus={() => { if (!disabled) open = true; }}
      oninput={() => { open = true; highlight = 0; }}
      onkeydown={handleKeydown}
    />
    {#if clearable && (search || selectedValues.length > 0) && !disabled}
      <button type="button" class="ohjanus-combobox-clear" aria-label="Clear" onclick={clear}><Icon name="x" size={10} /></button>
    {/if}
  </div>
  {#if open && !disabled}
    <ul id={listId} role="listbox" class="ohjanus-combobox-list">
      {#each filtered as opt, i (opt.value)}
        {@const selected = selectedValues.includes(opt.value)}
        <li
          role="option"
          aria-selected={selected}
          aria-disabled={opt.disabled || undefined}
          class:ohjanus-combobox-highlight={i === highlight}
          class:ohjanus-combobox-selected={selected}
        >
          <button
            type="button"
            disabled={opt.disabled}
            onclick={() => choose(opt.value, creatable && !options.some((o) => o.value === opt.value))}
            onmouseenter={() => (highlight = i)}
          >
            {#if option}{@render option(opt)}{:else}{opt.label}{/if}
          </button>
        </li>
      {:else}
        <li class="ohjanus-combobox-empty">
          {#if empty}{@render empty()}{:else}No results{/if}
        </li>
      {/each}
    </ul>
  {/if}
  {#if children}{@render children()}{/if}
</div>

<style>
  .ohjanus-combobox { position: relative; }
  .ohjanus-combobox-list {
    position: absolute; z-index: 50; inset-inline: 0; top: 100%;
    margin: 4px 0 0; padding: 4px; list-style: none;
    background: var(--ohjanus-surface, #fff);
    border: 1px solid var(--ohjanus-border, #ddd); border-radius: 6px;
    max-height: 240px; overflow: auto;
  }
  .ohjanus-combobox-highlight { background: var(--ohjanus-hover, #f0f0f0); }
  .ohjanus-combobox-tags { display: flex; flex-wrap: wrap; gap: 4px; align-items: center; }
  .ohjanus-combobox-tags input { flex: 1; min-width: 80px; border: none; outline: none; background: transparent; }
</style>
