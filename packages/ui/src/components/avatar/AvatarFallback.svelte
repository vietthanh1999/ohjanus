<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { class: className = '', style = '', children }: Props = $props();
</script>

<span class="ohjanus-avatar-fallback {className}" {style} aria-hidden="true">
  {@render children?.()}
</span>
