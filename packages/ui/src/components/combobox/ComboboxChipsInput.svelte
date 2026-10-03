<script lang="ts">
  interface Props {
    value?: string;
    placeholder?: string;
    class?: string;
    style?: string;
    disabled?: boolean;
  }

  let { value = $bindable(''), placeholder = '', class: className = '', style = '', disabled = false }: Props = $props();
</script>

<input
  type="text"
  class="ohjanus-combobox-chips-input {className}"
  {style}
  {placeholder}
  {disabled}
  bind:value
/>
