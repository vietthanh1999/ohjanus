<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    disabled?: boolean;
    onclick?: (e: MouseEvent) => void;
    children?: Snippet;
  }

  let { class: className = '', style = '', disabled = false, onclick, children }: Props = $props();
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<div
  class="ohjanus-command-item {disabled ? 'disabled' : ''} {className}"
  {style}
  role="option"
  aria-selected="false"
  tabindex={disabled ? -1 : 0}
  {onclick}
>
  {@render children?.()}
</div>
