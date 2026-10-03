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
  <path d="M4 1.5H3a2 2 0 00-2 2V14a2 2 0 002 2h10a2 2 0 002-2V3.5a2 2 0 00-2-2h-1v1h1a1 1 0 011 1V14a1 1 0 01-1 1H3a1 1 0 01-1-1V3.5a1 1 0 011-1h1v-1z"/>
</svg>
