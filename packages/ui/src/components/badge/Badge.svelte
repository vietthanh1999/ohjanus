<script lang="ts">
  import type { Snippet } from 'svelte';
  import { Icon } from '@ohjanus/icons';

  interface Props {
    variant?: 'default' | 'success' | 'warning' | 'danger' | 'destructive' | 'info' | 'outline' | 'secondary' | 'counter' | 'license';
    size?: 'sm' | 'md';
    dot?: boolean;
    removable?: boolean;
    onremove?: () => void;
    class?: string;
    children?: Snippet;
  }

  let {
    variant = 'default',
    size = 'sm',
    dot = false,
    removable = false,
    onremove,
    class: className = '',
    children
  }: Props = $props();
</script>

<span class="ohjanus-badge ohjanus-badge-{variant} ohjanus-badge-{size} {className}">
  {#if dot}
    <span class="ohjanus-badge-dot" aria-hidden="true"></span>
  {/if}
  {@render children?.()}
  {#if removable}
    <button
      type="button"
      class="ohjanus-badge-remove"
      aria-label="Remove"
      onclick={() => onremove?.()}
    ><Icon name="x" size={10} /></button>
  {/if}
</span>
