<script lang="ts">
  import Key from 'reicon-svelte/icons/Key.svelte';

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

<Key {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
