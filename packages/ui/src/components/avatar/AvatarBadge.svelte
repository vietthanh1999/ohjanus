<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    position?: 'top-right' | 'bottom-right' | 'top-left' | 'bottom-left';
    status?: 'online' | 'offline' | 'busy' | 'away';
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { position = 'bottom-right', status, class: className = '', style = '', children }: Props = $props();
</script>

<span
  class="ohjanus-avatar-badge ohjanus-avatar-badge-{position} {status ? `ohjanus-avatar-badge-${status}` : ''} {className}"
  {style}
>
  {@render children?.()}
</span>
