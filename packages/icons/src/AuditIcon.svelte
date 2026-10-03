<script lang="ts">
  import Document from 'reicon-svelte/icons/Document.svelte';

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

<Document {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
