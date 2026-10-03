<script lang="ts">
  import { getContext } from 'svelte';

  interface Props {
    class?: string;
    style?: string;
  }

  let { class: className = '', style = '' }: Props = $props();
  const ctx = getContext<{ step: number; totalSteps: number }>('ohjanus-questionnaire');
  const percent = $derived(ctx ? ((ctx.step + 1) / Math.max(1, ctx.totalSteps)) * 100 : 0);
</script>

<div
  class="ohjanus-questionnaire-progress {className}"
  {style}
  role="progressbar"
  aria-valuemin={1}
  aria-valuemax={ctx?.totalSteps ?? 1}
  aria-valuenow={(ctx?.step ?? 0) + 1}
>
  <div class="ohjanus-questionnaire-progress-indicator" style:width="{percent}%"></div>
</div>
