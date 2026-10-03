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
  <path d="M0 8a4 4 0 017.465-2H14a2 2 0 012 2v1a1 1 0 01-1 1h-1v1a1 1 0 01-1 1h-1v1a1 1 0 01-1 1H9.465A4 4 0 010 8zm4-2a2 2 0 100 4 2 2 0 000-4z"/>
</svg>
