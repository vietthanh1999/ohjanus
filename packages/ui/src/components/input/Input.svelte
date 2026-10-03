<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    value?: string | number;
    placeholder?: string;
    type?: string;
    disabled?: boolean;
    readonly?: boolean;
    invalid?: boolean;
    size?: 'sm' | 'md' | 'lg';
    class?: string;
    prefix?: Snippet;
    suffix?: Snippet;
    oninput?: (e: Event) => void;
    onchange?: (e: Event) => void;
    [key: string]: any;
  }

  let {
    value = $bindable(''),
    placeholder = '',
    type = 'text',
    disabled = false,
    readonly = false,
    invalid = false,
    size = 'sm',
    class: className = '',
    prefix,
    suffix,
    oninput,
    onchange,
    ...restProps
  }: Props = $props();
</script>

<div class="ohjanus-input-wrapper ohjanus-input-{size}" class:invalid class:disabled>
  {#if prefix}
    <span class="ohjanus-input-prefix">
      {@render prefix()}
    </span>
  {/if}
  <input
    {type}
    bind:value
    {placeholder}
    {disabled}
    {readonly}
    aria-invalid={invalid ? 'true' : undefined}
    class="ohjanus-input-field {className}"
    {oninput}
    {onchange}
    {...restProps}
  />
  {#if suffix}
    <span class="ohjanus-input-suffix">
      {@render suffix()}
    </span>
  {/if}
</div>
