<script lang="ts">
  import Refresh from 'reicon-svelte/icons/Refresh.svelte';

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

<Refresh {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
