<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { class: className = '', style = '', children }: Props = $props();
</script>

<li class="ohjanus-breadcrumb-item {className}" {style}>
  {@render children?.()}
</li>
