<script lang="ts">
  import type { Snippet } from 'svelte';
  import { getContext } from 'svelte';

  interface Props {
    value: string;
    id?: string;
    disabled?: boolean;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { value: choiceValue, id, disabled = false, class: className = '', style = '', children }: Props = $props();
  const ctx = getContext<{ value: string; name?: string; select: (v: string) => void }>(
    'ohjanus-questionnaire-choices'
  );
  const checked = $derived(ctx?.value === choiceValue);
</script>

<label class="ohjanus-questionnaire-choice {className}" class:checked {style}>
  <input
    type="radio"
    {id}
    name={ctx?.name}
    value={choiceValue}
    {checked}
    {disabled}
    class="ohjanus-questionnaire-choice-input"
    onchange={() => ctx?.select(choiceValue)}
  />
  <span class="ohjanus-radio-dot" aria-hidden="true"></span>
  <span>{@render children?.()}</span>
</label>
