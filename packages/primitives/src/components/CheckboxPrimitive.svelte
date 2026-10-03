<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    checked?: boolean;
    indeterminate?: boolean;
    disabled?: boolean;
    name?: string;
    value?: string;
    onchange?: (checked: boolean) => void;
    children?: Snippet<[{ checked: boolean; indeterminate: boolean; disabled: boolean }]>;
    class?: string;
    [key: string]: any;
  }

  let {
    checked = false,
    indeterminate = false,
    disabled = false,
    name,
    value,
    onchange,
    children,
    class: className = '',
    ...restProps
  }: Props = $props();

  let inputEl: HTMLInputElement | undefined = $state();

  $effect(() => {
    if (inputEl) {
      inputEl.indeterminate = indeterminate;
    }
  });

  function handleChange(e: Event) {
    const target = e.target as HTMLInputElement;
    onchange?.(target.checked);
  }
</script>

<label class="ohjanus-checkbox-primitive {className}">
  <input
    bind:this={inputEl}
    type="checkbox"
    {name}
    {value}
    {checked}
    {disabled}
    onchange={handleChange}
    style="position: absolute; opacity: 0; width: 0; height: 0;"
    {...restProps}
  />
  {#if children}
    {@render children({ checked, indeterminate, disabled })}
  {/if}
</label>
