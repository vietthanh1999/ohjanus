<script lang="ts">
  import Shield from 'reicon-svelte/icons/Shield.svelte';

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

<Shield {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
