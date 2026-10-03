<script lang="ts">
  import { getContext } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
  }

  let { class: className = '', style = '' }: Props = $props();
  const ctx = getContext<{ percent: number }>('ohjanus-progress');
</script>

<div
  class="ohjanus-progress-indicator {className}"
  {style}
  style:width="{ctx?.percent ?? 0}%; {style}"
></div>
