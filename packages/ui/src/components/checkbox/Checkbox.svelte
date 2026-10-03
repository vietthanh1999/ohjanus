<script lang="ts">
  import { Icon } from '@ohjanus/icons';
  import { CheckboxPrimitive } from '@ohjanus/primitives';

  interface Props {
    checked?: boolean;
    indeterminate?: boolean;
    disabled?: boolean;
    name?: string;
    value?: string;
    ariaLabel?: string;
    class?: string;
    onchange?: (checked: boolean) => void;
    [key: string]: any;
  }

  let {
    checked = $bindable(false),
    indeterminate = false,
    disabled = false,
    name,
    value,
    ariaLabel = 'Checkbox',
    class: className = '',
    onchange,
    ...restProps
  }: Props = $props();

  function handleChange(next: boolean) {
    checked = next;
    onchange?.(next);
  }
</script>

<CheckboxPrimitive {checked} {indeterminate} {disabled} {name} {value} onchange={handleChange} class={className} {...restProps}>
  {#snippet children({ checked: isChecked, indeterminate: isIndeterminate, disabled: isDisabled })}
    <span
      class="ohjanus-checkbox"
      class:checked={isChecked}
      class:indeterminate={isIndeterminate}
      class:disabled={isDisabled}
      role="presentation"
    >
      {#if isIndeterminate}
        <span class="ohjanus-checkbox-mark" aria-hidden="true"><Icon name="minus" size={10} /></span>
      {:else if isChecked}
        <span class="ohjanus-checkbox-mark" aria-hidden="true"><Icon name="check" size={10} /></span>
      {/if}
    </span>
    <span class="ohjanus-visually-hidden">{ariaLabel}</span>
  {/snippet}
</CheckboxPrimitive>
