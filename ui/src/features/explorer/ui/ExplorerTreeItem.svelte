<script lang="ts">
  import { Button, Text, Box } from '@ohjanus/ui';
  import { Icon, type IconName } from '@ohjanus/icons';

  interface Props {
    depth?: 0 | 1 | 2;
    icon: IconName;
    iconColor?: string;
    label: string;
    badgeText?: string;
    badgeCount?: number;
    hasChevron?: boolean;
    isExpanded?: boolean;
    isSelected?: boolean;
    warnStatus?: string;
    onclick?: () => void;
  }

  let {
    depth = 0,
    icon,
    iconColor = '#3B82F6',
    label,
    badgeText,
    badgeCount,
    hasChevron = false,
    isExpanded = false,
    isSelected = false,
    warnStatus,
    onclick
  }: Props = $props();
</script>

<Button
  variant="ghost"
  class="tree-node depth-{depth} {isSelected ? 'selected' : ''}"
  onclick={() => onclick?.()}
>
  {#if hasChevron}
    <Box class="chevron {isExpanded ? 'expanded' : ''}">
      <Icon name="chevron-right" size={12} />
    </Box>
  {:else}
    <Box class="chevron-placeholder" />
  {/if}

  <Icon name={icon} size={13} color={iconColor} class="node-icon" />

  <Text size="sm" truncate class="node-label">
    {label}
    {#if badgeText}
      <Text size="xs" color="muted" style="margin-left: 4px;">{badgeText}</Text>
    {/if}
    {#if badgeCount !== undefined}
      <Text size="xs" color="muted" style="margin-left: 4px;">{badgeCount}</Text>
    {/if}
  </Text>

  {#if warnStatus}
    <Text size="xs" color="danger" class="conn-warn" title={warnStatus}>●</Text>
  {/if}
</Button>

<style>
  :global(.tree-node) {
    height: var(--tree-row-height, 28px) !important;
    display: flex !important;
    align-items: center !important;
    gap: 4px !important;
    padding-right: 8px !important;
    cursor: pointer !important;
    font-size: var(--font-size-sm, 12px) !important;
    color: var(--text-primary) !important;
    border-radius: var(--radius-sm, 4px) !important;
    margin: 1px 4px !important;
    transition: background-color 0.1s ease !important;
    text-align: left !important;
    width: calc(100% - 8px) !important;
    justify-content: flex-start !important;
    border: none !important;
    font-weight: normal !important;
  }

  :global(.tree-node:hover) {
    background-color: var(--bg-hover) !important;
  }

  :global(.tree-node.selected) {
    background-color: var(--bg-selected) !important;
    color: #FFFFFF !important;
  }

  :global(.tree-node.depth-0) {
    padding-left: 6px !important;
    font-weight: 600 !important;
  }
  :global(.tree-node.depth-1) {
    padding-left: 20px !important;
  }
  :global(.tree-node.depth-2) {
    padding-left: 34px !important;
  }

  :global(.chevron) {
    width: 14px;
    height: 14px;
    font-size: 12px;
    color: var(--text-muted);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    transition: transform 0.15s ease;
    transform: rotate(0deg);
    flex-shrink: 0;
  }

  :global(.chevron.expanded) {
    transform: rotate(90deg);
  }

  :global(.chevron-placeholder) {
    width: 14px;
    height: 14px;
    flex-shrink: 0;
  }

  :global(.node-icon) {
    margin-right: 2px;
    width: 14px;
    height: 14px;
    flex-shrink: 0;
  }

  :global(.node-label) {
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  :global(.conn-warn) {
    margin-left: auto;
    font-size: 9px;
  }
</style>
