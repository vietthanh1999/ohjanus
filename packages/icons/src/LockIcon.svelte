<script lang="ts">
  import Lock from 'reicon-svelte/icons/Lock.svelte';

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

<Lock {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
