<script lang="ts">
  import { Flex, Text, Button } from '@ohjanus/ui';
  import { Icon } from '@ohjanus/icons';

  interface Props {
    onreload?: () => void;
    onnewconsole?: () => void;
  }

  let { onreload, onnewconsole }: Props = $props();
</script>

<Flex align="center" justify="between" class="explorer-header-top">
  <Text size="sm" weight="semibold">Database Explorer</Text>
  <Flex align="center" gap="xs" class="header-window-icons">
    <Button
      variant="ghost"
      size="icon-xs"
      class="jb-icon-btn"
      title="Reload connections"
      onclick={() => onreload?.()}
    >
      <Icon name="refresh" size={12} />
    </Button>
    <Button
      variant="ghost"
      size="icon-xs"
      class="jb-icon-btn"
      title="New console on first connection"
      onclick={() => onnewconsole?.()}
    >
      <Icon name="plus" size={12} />
    </Button>
  </Flex>
</Flex>

<style>
  :global(.explorer-header-top) {
    height: var(--toolbar-height, 32px);
    padding: 0 8px 0 12px;
    flex-shrink: 0;
  }

  :global(.explorer-header-top .jb-icon-btn) {
    width: 22px !important;
    height: 22px !important;
    padding: 0 !important;
    min-width: 22px !important;
    color: var(--text-secondary);
  }

  :global(.explorer-header-top .jb-icon-btn:hover) {
    color: var(--text-primary);
    background-color: var(--bg-hover);
  }
</style>
