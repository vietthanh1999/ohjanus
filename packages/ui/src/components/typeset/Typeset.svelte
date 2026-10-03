<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    preset?: 'docs' | 'chat' | 'reading' | 'compact' | 'large' | '';
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { preset = 'docs', class: className = '', style = '', children }: Props = $props();
</script>

<div class="ohjanus-typeset {preset ? `ohjanus-typeset-${preset}` : ''} {className}" {style}>
  {@render children?.()}
</div>
