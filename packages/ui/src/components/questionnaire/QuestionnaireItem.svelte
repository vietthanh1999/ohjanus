<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext } from 'svelte';

  interface Props {
    index?: number;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { index, class: className = '', style = '', children }: Props = $props();
  const ctx = getContext<{ step: number } | undefined>('ohjanus-questionnaire');
  const visible = $derived(!ctx || index === undefined || ctx.step === index);
</script>

{#if visible}
  <div class="ohjanus-questionnaire-item {className}" {style}>
    {@render children?.()}
  </div>
{/if}
