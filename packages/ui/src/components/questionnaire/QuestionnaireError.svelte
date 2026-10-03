<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { class: className = '', style = '', children }: Props = $props();
</script>

<p class="ohjanus-questionnaire-error {className}" {style} role="alert">
  {@render children?.()}
</p>
