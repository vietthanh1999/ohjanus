<script lang="ts">
  import { Icon } from '@ohjanus/icons';
  import { toast } from '@ohjanus/primitives';
  import Portal from '@ohjanus/primitives/components/Portal.svelte';
</script>

<Portal>
  <div class="ohjanus-toaster-container">
    {#each toast.toasts as item (item.id)}
      <div
        class="ohjanus-toast-item ohjanus-toast-{item.type}"
        role="alert"
      >
        <div class="ohjanus-toast-content">
          {#if item.title}
            <div class="ohjanus-toast-title">{item.title}</div>
          {/if}
          {#if item.description}
            <div class="ohjanus-toast-description">{item.description}</div>
          {/if}
        </div>
        <button
          type="button"
          class="ohjanus-toast-close"
          onclick={() => toast.dismiss(item.id)}
          aria-label="Dismiss toast"
        >
          <Icon name="x" size={12} />
        </button>
      </div>
    {/each}
  </div>
</Portal>
