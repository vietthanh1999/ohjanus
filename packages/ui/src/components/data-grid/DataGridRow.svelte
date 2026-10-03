<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
    selected?: boolean;
    onclick?: (e: MouseEvent) => void;
    children?: Snippet;
  }

  let { class: className = '', style = '', selected = false, onclick, children }: Props = $props();
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
<tr
  class="ohjanus-data-grid-row {selected ? 'selected' : ''} {className}"
  {style}
  {onclick}
>
  {@render children?.()}
</tr>
