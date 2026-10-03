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
  <path fill-rule="evenodd" d="M8 0c-.69 0-1.843.265-2.928.56-1.11.3-2.229.655-2.887.87a1.54 1.54 0 00-1.044 1.262c-.596 4.477.787 7.795 2.464 9.99 1.579 2.065 3.444 3.009 4.395 3.318.066.022.135.034.204.034s.138-.012.204-.034c.951-.309 2.816-1.253 4.395-3.318 1.677-2.195 3.06-5.513 2.464-9.99a1.54 1.54 0 00-1.044-1.263 62.467 62.467 0 00-2.887-.87C9.843.266 8.69 0 8 0zm2.146 5.146a.5.5 0 01.708.708l-3 3a.5.5 0 01-.708 0l-1.5-1.5a.5.5 0 11.708-.708L7.5 7.793l2.646-2.647z"/>
</svg>
