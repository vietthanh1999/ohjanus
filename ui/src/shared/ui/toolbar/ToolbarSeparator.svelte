<script lang="ts">
  import { Divider } from '@ohjanus/ui';
</script>

<Divider orientation="vertical" class="bar-separator" />

<style>
  :global(.bar-separator) {
    height: 14px !important;
    margin: 0 4px !important;
    background-color: var(--border-default) !important;
    flex-shrink: 0;
  }
</style>
