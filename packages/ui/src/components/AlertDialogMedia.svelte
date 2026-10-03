<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { class: className = '', style = '', children }: Props = $props();
</script>

<div class="ohjanus-alert-dialog-media {className}" {style} aria-hidden="true">
  {@render children?.()}
</div>
