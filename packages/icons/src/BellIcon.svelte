<script lang="ts">
  import Bell from 'reicon-svelte/icons/Bell.svelte';

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

<Bell {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
