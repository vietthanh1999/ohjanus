<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    checked?: boolean;
    disabled?: boolean;
    onCheckedChange?: (checked: boolean) => void;
    class?: string;
    style?: string;
    children?: Snippet;
  }

  let { checked = $bindable(false), disabled = false, onCheckedChange, class: className = '', style = '', children }: Props = $props();

  function toggle() {
    if (disabled) return;
    checked = !checked;
    onCheckedChange?.(checked);
  }
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
  role="menuitemcheckbox"
  aria-checked={checked}
  aria-disabled={disabled || undefined}
  class="ohjanus-dropdown-checkbox-item {className}"
  {style}
  tabindex={disabled ? -1 : 0}
  onclick={toggle}
  onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); toggle(); } }}
>
  <span class="ohjanus-dropdown-check" aria-hidden="true">{checked ? '✓' : ''}</span>
  {@render children?.()}
</div>
