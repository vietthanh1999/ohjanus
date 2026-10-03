<script lang="ts">
  interface Props {
    value?: string;
    placeholder?: string;
    rows?: number;
    disabled?: boolean;
    readonly?: boolean;
    invalid?: boolean;
    autoResize?: boolean;
    class?: string;
    style?: string;
    id?: string;
    name?: string;
    oninput?: (e: Event) => void;
    onchange?: (e: Event) => void;
    [key: string]: any;
  }

  let {
    value = $bindable(''),
    placeholder = '',
    rows = 3,
    disabled = false,
    readonly = false,
    invalid = false,
    autoResize = false,
    class: className = '',
    style = '',
    id,
    name,
    oninput,
    onchange,
    ...restProps
  }: Props = $props();

  function handleInput(e: Event) {
    if (autoResize) {
      const el = e.target as HTMLTextAreaElement;
      el.style.height = 'auto';
      el.style.height = `${el.scrollHeight}px`;
    }
    oninput?.(e);
  }
</script>

<textarea
  bind:value
  {placeholder}
  {rows}
  {disabled}
  {readonly}
  {id}
  {name}
  aria-invalid={invalid ? 'true' : undefined}
  class="ohjanus-textarea {invalid ? 'invalid' : ''} {className}"
  {style}
  oninput={handleInput}
  {onchange}
  {...restProps}
></textarea>
