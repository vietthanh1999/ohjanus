<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    value?: string;
    placeholder?: string;
    class?: string;
    style?: string;
    autofocus?: boolean;
    children?: Snippet;
  }

  let {
    value = $bindable(''),
    placeholder = 'Type a command or search...',
    class: className = '',
    style = '',
    autofocus = true,
    children
  }: Props = $props();
</script>

<div class="ohjanus-command-input-wrapper {className}" {style}>
  <svg class="ohjanus-command-search-icon" width="15" height="15" viewBox="0 0 16 16" fill="currentColor">
    <path fill-rule="evenodd" d="M11.5 7a4.5 4.5 0 11-9 0 4.5 4.5 0 019 0zm-.82 4.74a6 6 0 111.06-1.06l3.04 3.04a.75.75 0 11-1.06 1.06l-3.04-3.04z"/>
  </svg>
  <!-- svelte-ignore a11y_autofocus -->
  <input
    type="text"
    class="ohjanus-command-input"
    {placeholder}
    bind:value
    {autofocus}
  />
  {@render children?.()}
</div>
