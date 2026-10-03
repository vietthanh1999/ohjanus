<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    onRemove?: () => void;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { onRemove, class: className = '', style = '', children }: Props = $props();
</script>

<span class="ohjanus-combobox-chip {className}" {style}>
  {@render children?.()}
  {#if onRemove}
    <button type="button" class="ohjanus-combobox-chip-remove" onclick={onRemove} aria-label="Remove">×</button>
  {/if}
</span>
