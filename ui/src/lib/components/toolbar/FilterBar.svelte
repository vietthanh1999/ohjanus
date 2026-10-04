<script lang="ts">
  import { Icon } from '@ohjanus/icons';
  import { Flex, Box, Text, Input } from '@ohjanus/ui';

  interface Props {
    whereValue?: string;
    wherePlaceholder?: string;
    orderByValue?: string;
    orderByPlaceholder?: string;
    searchRef?: HTMLInputElement;
    onwherechange?: (val: string) => void;
    onorderbychange?: (val: string) => void;
    class?: string;
  }

  let {
    whereValue = $bindable(''),
    wherePlaceholder = "e.g. id IS NOT NULL",
    orderByValue = $bindable(''),
    orderByPlaceholder = "e.g. created_at DESC",
    searchRef = $bindable(),
    onwherechange,
    onorderbychange,
    class: className = ''
  }: Props = $props();
</script>

<Flex align="center" gap="12px" class="filter-bar {className}">
  <Flex align="center" gap="6px" class="filter-group where-group">
    <Box as="span" class="filter-icon">
      <Icon name="filter" size={11} />
    </Box>
    <Text size="sm" weight="semibold" color="secondary" style="user-select: none;">WHERE</Text>
    <Input
      size="sm"
      class="filter-input code-text"
      placeholder={wherePlaceholder}
      bind:value={whereValue}
      onchange={() => onwherechange?.(whereValue)}
    />
  </Flex>

  <Flex align="center" gap="6px" class="filter-group orderby-group">
    <Box as="span" class="filter-icon">
      <Icon name="sort" size={11} />
    </Box>
    <Text size="sm" weight="semibold" color="secondary" style="user-select: none;">ORDER BY</Text>
    <Input
      size="sm"
      class="filter-input code-text"
      placeholder={orderByPlaceholder}
      bind:value={orderByValue}
      onchange={() => onorderbychange?.(orderByValue)}
    />
  </Flex>
</Flex>

<style>
  :global(.filter-bar) {
    height: var(--filterbar-height, 30px);
    background-color: var(--bg-canvas);
    border-bottom: 1px solid var(--border-subtle);
    padding: 0 8px;
    flex-shrink: 0;
  }

  :global(.filter-group) {
    display: flex;
    align-items: center;
  }

  :global(.where-group) {
    flex: 1.1;
  }

  :global(.orderby-group) {
    flex: 0.9;
  }

  :global(.filter-icon) {
    color: var(--text-muted);
    font-size: 11px;
    user-select: none;
    display: inline-flex;
    align-items: center;
  }

  :global(.filter-bar .ohjanus-input-wrapper) {
    flex: 1;
    background-color: transparent !important;
    border: none !important;
    box-shadow: none !important;
    height: var(--control-height-xs, 24px) !important;
    padding: 0 !important;
  }

  :global(.filter-bar .ohjanus-input-field) {
    background-color: transparent !important;
    border: none !important;
    outline: none !important;
    box-shadow: none !important;
    padding: 0 6px !important;
    font-size: var(--font-size-sm, 13px) !important;
    color: var(--text-primary) !important;
    font-family: var(--font-code) !important;
    height: 100% !important;
  }

  :global(.filter-bar .ohjanus-input-field:focus),
  :global(.filter-bar .ohjanus-input-field:hover) {
    border: none !important;
    outline: none !important;
    box-shadow: none !important;
  }
</style>
