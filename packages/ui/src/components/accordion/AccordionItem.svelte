<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    value: string;
    disabled?: boolean;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { value, disabled = false, class: className = '', style = '', children }: Props = $props();
</script>

<div class="ohjanus-accordion-item {className}" {style} data-value={value} data-disabled={disabled || undefined}>
  {@render children?.()}
</div>
