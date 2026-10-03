<script lang="ts">
  import { getContext } from 'svelte';

  interface Props {
    placeholder?: string;
    class?: string;
    style?: string;
  }

  let { placeholder = 'Select...', class: className = '', style = '' }: Props = $props();
  const ctx = getContext<{
    value: string | string[];
    multiple: boolean;
    labels: Record<string, string>;
  }>('ohjanus-select');

  const text = $derived(() => {
    const v = ctx?.value;
    if (ctx?.multiple) {
      const arr = Array.isArray(v) ? v : [];
      if (!arr.length) return '';
      return arr.map((item) => ctx.labels[item] ?? item).join(', ');
    }
    if (!v || (Array.isArray(v) && !v.length)) return '';
    const key = Array.isArray(v) ? v[0] : (v as string);
    return ctx?.labels[key] ?? key ?? '';
  });
</script>

<span class="ohjanus-select-value {className}" {style} class:ohjanus-select-placeholder={!text()}>
  {text() || placeholder}
</span>
