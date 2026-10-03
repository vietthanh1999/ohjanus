<script lang="ts">
  import type { Snippet } from 'svelte';
  import { Icon } from '@ohjanus/icons';
  import { setContext } from 'svelte';

  interface Option {
    value: string;
    label: string;
    disabled?: boolean;
  }

  interface Props {
    value?: string | string[];
    options?: Option[];
    placeholder?: string;
    multiple?: boolean;
    searchable?: boolean;
    clearable?: boolean;
    disabled?: boolean;
    invalid?: boolean;
    class?: string;
    style?: string;
    id?: string;
    name?: string;
    onchange?: (value: string | string[]) => void;
    children?: Snippet;
    [key: string]: any;
  }

  let {
    value = $bindable(''),
    options = [],
    placeholder = 'Select...',
    multiple = false,
    searchable = false,
    clearable = false,
    disabled = false,
    invalid = false,
    class: className = '',
    style = '',
    id,
    name,
    onchange,
    children,
    ...restProps
  }: Props = $props();

  let search = $state('');
  let open = $state(false);
  let rootEl: HTMLDivElement | undefined = $state();

  let filtered = $derived(
    searchable && search
      ? options.filter((o) => o.label.toLowerCase().includes(search.toLowerCase()))
      : options
  );

  function handleDocumentClick(e: MouseEvent) {
    if (open && rootEl && !rootEl.contains(e.target as Node)) {
      open = false;
    }
  }

  function handleDocumentKeydown(e: KeyboardEvent) {
    if (open && e.key === 'Escape') {
      open = false;
    }
  }

  $effect(() => {
    if (open) {
      window.addEventListener('click', handleDocumentClick, true);
      window.addEventListener('keydown', handleDocumentKeydown, true);
      return () => {
        window.removeEventListener('click', handleDocumentClick, true);
        window.removeEventListener('keydown', handleDocumentKeydown, true);
      };
    }
  });

  function isSelected(optValue: string): boolean {
    return Array.isArray(value) ? value.includes(optValue) : value === optValue;
  }

  function select(optValue: string) {
    if (multiple) {
      const arr = Array.isArray(value) ? value : [];
      const next = arr.includes(optValue) ? arr.filter((v) => v !== optValue) : [...arr, optValue];
      value = next;
      onchange?.(next);
    } else {
      value = optValue;
      onchange?.(optValue);
      open = false;
    }
  }

  function clear(e: MouseEvent) {
    e.stopPropagation();
    value = multiple ? [] : '';
    onchange?.(value);
  }

  // Labels registered by compound <SelectItem /> parts.
  const labels = $state<Record<string, string>>({});

  function registerLabel(itemValue: string, label: string) {
    labels[itemValue] = label;
  }

  const displayLabel = $derived(() => {
    if (multiple) {
      const arr = Array.isArray(value) ? value : [];
      if (!arr.length) return '';
      return options
        .filter((o) => arr.includes(o.value))
        .map((o) => o.label)
        .join(', ');
    }
    return options.find((o) => o.value === value)?.label ?? labels[value as string] ?? '';
  });

  function setOpen(v: boolean) {
    open = v;
  }

  setContext('ohjanus-select', {
    get value() { return value; },
    get open() { return open; },
    get multiple() { return multiple; },
    get disabled() { return disabled; },
    get labels() { return labels; },
    select,
    setOpen,
    isSelected,
    registerLabel
  });
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div bind:this={rootEl} class="ohjanus-select {invalid ? 'invalid' : ''} {disabled ? 'disabled' : ''} {className}" {style}>
  {#if searchable}
    <input
      type="text"
      class="ohjanus-select-search"
      placeholder={displayLabel() || placeholder}
      bind:value={search}
      {disabled}
      onfocus={() => (open = true)}
    />
  {:else}
    <button
      type="button"
      class="ohjanus-select-trigger"
      {disabled}
      onclick={() => (open = !open)}
      aria-haspopup="listbox"
      aria-expanded={open}
    >
      <span class:ohjanus-select-placeholder={!displayLabel()}>{displayLabel() || placeholder}</span>
      <span class="ohjanus-select-caret" class:open aria-hidden="true"><Icon name="chevron-down" size={11} /></span>
    </button>
  {/if}
  {#if clearable && (Array.isArray(value) ? value.length : value)}
    <button type="button" class="ohjanus-select-clear" onclick={clear} aria-label="Clear selection"><Icon name="x" size={10} /></button>
  {/if}
  {#if open}
    <ul class="ohjanus-select-list" role="listbox" aria-multiselectable={multiple}>
      {#each filtered as opt (opt.value)}
        <li>
          <button
            type="button"
            role="option"
            aria-selected={isSelected(opt.value)}
            disabled={opt.disabled}
            class="ohjanus-select-option"
            class:selected={isSelected(opt.value)}
            onclick={() => select(opt.value)}
          >
            {#if multiple}
              {#if isSelected(opt.value)}
                <span aria-hidden="true"><Icon name="check-square" size={12} /></span>
              {:else}
                <span class="ohjanus-select-check-empty" aria-hidden="true"></span>
              {/if}
            {/if}
            <span class="ohjanus-select-option-label">{opt.label}</span>
            {#if !multiple && isSelected(opt.value)}
              <span class="ohjanus-select-check" aria-hidden="true"><Icon name="check" size={12} /></span>
            {/if}
          </button>
        </li>
      {:else}
        <li class="ohjanus-select-empty">No options</li>
      {/each}
    </ul>
  {/if}
  {#if children}
    {@render children()}
  {/if}
  {#if name}
    <select {name} {id} {multiple} {disabled} style="display:none" aria-hidden="true" tabindex="-1">
      {#each options as opt (opt.value)}
        <option value={opt.value} selected={isSelected(opt.value)}>{opt.label}</option>
      {/each}
    </select>
  {/if}
</div>
