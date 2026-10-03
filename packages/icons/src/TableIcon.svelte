<script lang="ts">
  import Grid from 'reicon-svelte/icons/Grid.svelte';

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

<Grid {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
