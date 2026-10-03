<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    max?: number;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { max, class: className = '', style = '', children }: Props = $props();
</script>

<div class="ohjanus-avatar-group {className}" {style} data-max={max}>
  {@render children?.()}
</div>
