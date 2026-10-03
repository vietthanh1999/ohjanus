<script lang="ts">
  import type { Snippet } from 'svelte';
  import { setContext } from 'svelte';

  interface Props {
    step?: number;
    totalSteps?: number;
    class?: string;
    style?: string;
    children?: Snippet;
    onStepChange?: (step: number) => void;
    onSubmit?: () => void;
  }

  let {
    step = $bindable(0),
    totalSteps = 1,
    class: className = '',
    style = '',
    children,
    onStepChange,
    onSubmit
  }: Props = $props();

  function go(delta: number) {
    const next = Math.min(totalSteps - 1, Math.max(0, step + delta));
    step = next;
    onStepChange?.(next);
  }

  function previous() { go(-1); }
  function next() { go(1); }
  function skip() { go(1); }
  function submit() { onSubmit?.(); }

  setContext('ohjanus-questionnaire', {
    get step() { return step; },
    get totalSteps() { return totalSteps; },
    previous,
    next,
    skip,
    submit
  });
</script>

<div class="ohjanus-questionnaire {className}" {style} data-step={step}>
  {@render children?.()}
</div>
