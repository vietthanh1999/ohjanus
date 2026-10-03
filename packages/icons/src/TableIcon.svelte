<script lang="ts">
  interface Props {
    size?: number | string;
    color?: string;
    class?: string;
    [key: string]: any;
  }

  let {
    size = 16,
    color = 'currentColor',
    class: className = '',
    ...restProps
  }: Props = $props();
</script>

<svg
  xmlns="http://www.w3.org/2000/svg"
  width={size}
  height={size}
  viewBox="0 0 16 16"
  fill={color}
  class="ohjanus-icon {className}"
  {...restProps}
>
  <path d="M0 2a1 1 0 011-1h14a1 1 0 011 1v12a1 1 0 01-1 1H1a1 1 0 01-1-1V2zm1 3v2h6V5H1zm7 0v2h7V5H8zm0 3v2h7V8H8zm-1 0H1v2h6V8zm0 3H1v2h6v-2zm1 0v2h7v-2H8z"/>
</svg>
