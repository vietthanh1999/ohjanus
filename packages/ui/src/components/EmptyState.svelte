<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    title?: string;
    description?: string;
    class?: string;
    icon?: Snippet;
    action?: Snippet;
  }

  let {
    title = 'No data available',
    description = '',
    class: className = '',
    icon,
    action
  }: Props = $props();
</script>

<div class="ohjanus-empty-state {className}">
  {#if icon}
    <div class="ohjanus-empty-icon">
      {@render icon()}
    </div>
  {/if}
  <h4 class="ohjanus-empty-title">{title}</h4>
  {#if description}
    <p class="ohjanus-empty-description">{description}</p>
  {/if}
  {#if action}
    <div class="ohjanus-empty-action">
      {@render action()}
    </div>
  {/if}
</div>
