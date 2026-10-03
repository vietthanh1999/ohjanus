<script lang="ts">
  interface Props {
    value?: string;
    placeholder?: string;
    disabled?: boolean;
    invalid?: boolean;
    class?: string;
    style?: string;
    id?: string;
    name?: string;
  }

  let {
    value = $bindable(''),
    placeholder = '',
    disabled = false,
    invalid = false,
    class: className = '',
    style = '',
    id,
    name
  }: Props = $props();
</script>

<input
  type="text"
  bind:value
  {placeholder}
  {disabled}
  {id}
  {name}
  aria-invalid={invalid ? 'true' : undefined}
  class="ohjanus-questionnaire-input {className}"
  class:invalid
  {style}
/>
