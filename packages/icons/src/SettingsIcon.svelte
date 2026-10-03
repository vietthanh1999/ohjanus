<script lang="ts">
  import Settings from 'reicon-svelte/icons/Settings.svelte';

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

<Settings {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
