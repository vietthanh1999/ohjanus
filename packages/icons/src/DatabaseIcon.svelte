<script lang="ts">
  interface Props {
    size?: number | string;
    color?: string;
    strokeWidth?: number;
    class?: string;
    [key: string]: any;
  }

  let {
    size = 16,
    color = 'currentColor',
    strokeWidth = 2,
    class: className = '',
    ...restProps
  }: Props = $props();
</script>

<svg
  xmlns="http://www.w3.org/2000/svg"
  width={size}
  height={size}
  viewBox="0 0 24 24"
  fill="none"
  stroke={color}
  stroke-width={strokeWidth}
  stroke-linecap="round"
  stroke-linejoin="round"
  class="ohjanus-icon {className}"
  {...restProps}
>
  <ellipse cx="12" cy="5" rx="9" ry="3" />
  <path d="M3 5V19A9 3 0 0 0 21 19V5" />
  <path d="M3 12A9 3 0 0 0 21 12" />
</svg>
