<script lang="ts">
  import { CheckboxPrimitive } from "@ohjanus/primitives";

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
    ariaLabel = "Checkbox",
    class: className = "",
    onchange,
    ...restProps
  }: Props = $props();

  function handleChange(next: boolean) {
    checked = next;
    onchange?.(next);
  }
</script>

<CheckboxPrimitive
  {checked}
  {indeterminate}
  {disabled}
  {name}
  {value}
  onchange={handleChange}
  class={className}
  {...restProps}
>
  {#snippet children({
    checked: isChecked,
    indeterminate: isIndeterminate,
    disabled: isDisabled,
  })}
    <span
      class="ohjanus-checkbox"
      class:checked={isChecked}
      class:indeterminate={isIndeterminate}
      class:disabled={isDisabled}
      role="presentation"
    ></span>
    <span class="ohjanus-visually-hidden">{ariaLabel}</span>
  {/snippet}
</CheckboxPrimitive>
