<script lang="ts">
  import Search from 'reicon-svelte/icons/Search.svelte';

  interface Props {
    size?: number | string;
    color?: string;
    strokeWidth?: number | string;
    class?: string;
    [key: string]: any;
  }

  let {
    size = 16,
    color = 'currentColor',
    strokeWidth,
    class: className = '',
    ...restProps
  }: Props = $props();
</script>

<Search {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
