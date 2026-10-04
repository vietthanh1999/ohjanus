<script lang="ts">
  import { Select, Box } from '@ohjanus/ui';

  interface Option {
    value: string;
    label: string;
    disabled?: boolean;
  }

  interface Props {
    options: Option[];
    value?: string;
    class?: string;
    onchange?: (val: string | string[]) => void;
  }

  let { options, value = $bindable(''), class: className = '', onchange }: Props = $props();
</script>

<Box class="borderless-select-wrapper {className}">
  <Select
    class="borderless-select"
    {options}
    bind:value
    {onchange}
  />
</Box>

<style>
  :global(.borderless-select-wrapper) {
    display: inline-flex;
    align-items: center;
  }

  :global(.borderless-select) {
    min-width: unset !important;
    width: auto !important;
  }

  :global(.borderless-select .ohjanus-select-trigger) {
    background-color: transparent !important;
    border: none !important;
    box-shadow: none !important;
    height: 24px !important;
    padding: 0 6px !important;
    gap: 4px !important;
    font-size: var(--font-size-xs, 12px) !important;
    color: var(--text-secondary, #9DA0A8) !important;
    cursor: pointer;
    transition: background-color 0.1s ease, color 0.1s ease;
  }

  :global(.borderless-select .ohjanus-select-trigger:hover) {
    background-color: var(--bg-hover, #313438) !important;
    color: var(--text-primary, #DFE1E5) !important;
  }
</style>
