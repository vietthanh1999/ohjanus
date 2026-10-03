<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl';
    shape?: 'circle' | 'square';
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { size = 'md', shape = 'circle', class: className = '', style = '', children }: Props = $props();
</script>

<span class="ohjanus-avatar ohjanus-avatar-{size} ohjanus-avatar-{shape} {className}" {style}>
  {@render children?.()}
</span>
