<script lang="ts">
  import { TogglePrimitive } from '@ohjanus/primitives';

  interface Props {
    checked?: boolean;
    disabled?: boolean;
    label?: string;
    class?: string;
    onchange?: (checked: boolean) => void;
  }

  let {
    checked = $bindable(false),
    disabled = false,
    label,
    class: className = '',
    onchange
  }: Props = $props();
</script>

<label class="ohjanus-switch-wrapper {className}">
  <TogglePrimitive {checked} {disabled} {onchange}>
    {#snippet children({ checked: isChecked, disabled: isDisabled })}
      <span class="ohjanus-switch-track" class:checked={isChecked} class:disabled={isDisabled}>
        <span class="ohjanus-switch-thumb"></span>
      </span>
    {/snippet}
  </TogglePrimitive>
  {#if label}
    <span class="ohjanus-switch-label">{label}</span>
  {/if}
</label>
