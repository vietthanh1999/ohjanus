<script lang="ts">
  import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, Box, Flex, Text } from '@ohjanus/ui';
  import { Icon, type IconName } from '@ohjanus/icons';

  export interface RowActionItem {
    icon: IconName;
    label: string;
    onclick: () => void;
  }

  interface Props {
    title?: string;
    items: RowActionItem[];
  }

  let { title = 'Row actions', items }: Props = $props();
  let open = $state(false);
</script>

<!-- Custom span trigger (not DropdownMenuTrigger) so there is no nested
  <button> inside the tree-row button. The menu content is portaled to
  <body> by the core DropdownMenu, so it never clips inside the sidebar. -->
<DropdownMenu bind:open class="row-actions">
  <Box
    as="span"
    class="row-actions-trigger"
    role="button"
    tabindex={0}
    title={title}
    onclick={(e) => {
      e.stopPropagation();
      open = !open;
    }}
    onkeydown={(e: KeyboardEvent) => {
      if (e.key === 'Enter' || e.key === ' ') {
        e.preventDefault();
        e.stopPropagation();
        open = !open;
      }
    }}
  >
    <Icon name="more" size={14} />
  </Box>
  <DropdownMenuContent align="end">
    {#each items as item (item.label)}
      <DropdownMenuItem onclick={() => item.onclick()}>
        <Flex align="center" gap="8px">
          <Icon name={item.icon} size={13} />
          <Text size="sm">{item.label}</Text>
        </Flex>
      </DropdownMenuItem>
    {/each}
  </DropdownMenuContent>
</DropdownMenu>

<style>
  /* DataGrip-style: actions appear on row hover (or while the menu is open). */
  :global(.tree-node .row-actions) {
    opacity: 0;
    flex-shrink: 0;
    margin-left: 2px;
  }
  :global(.tree-node:hover .row-actions),
  :global(.tree-node:focus-within .row-actions),
  :global(.row-actions[data-open='true']) {
    opacity: 1;
  }

  :global(.row-actions-trigger) {
    width: 22px;
    height: 22px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--radius-sm, 4px);
    color: var(--text-secondary, #9DA0A8);
    cursor: pointer;
  }
  :global(.row-actions-trigger:hover) {
    background-color: var(--bg-hover, #313438);
    color: var(--text-primary, #DFE1E5);
  }
</style>
