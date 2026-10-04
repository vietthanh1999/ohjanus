<script lang="ts">
  import { explorerState } from "@/features/explorer";
  import { Button, Flex, Text } from "@ohjanus/ui";
  import {
    DropdownMenu,
    DropdownMenuTrigger,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuCheckboxItem,
    DropdownMenuSeparator,
  } from "@ohjanus/ui";
  import { Tooltip, TooltipTrigger, TooltipContent } from "@ohjanus/ui";
  import { Icon } from "@ohjanus/icons";

  interface Props {
    onreload?: () => void;
    onnewconsole?: () => void;
    onmanageconnections?: () => void;
    oncopyddl?: () => void;
  }

  let { onreload, onnewconsole, onmanageconnections, oncopyddl }: Props =
    $props();
</script>

<Flex align="center" justify="between" class="explorer-header-top">
  <Text size="sm" weight="semibold">Database Explorer</Text>
  <Flex align="center" gap="1" class="header-window-icons">
    <Tooltip>
      <TooltipTrigger>
        <Button
          variant="ghost"
          size="icon-xs"
          class="jb-icon-btn {explorerState.scrollFromEditor ? 'toggled' : ''}"
          onclick={() =>
            (explorerState.scrollFromEditor = !explorerState.scrollFromEditor)}
        >
          <Icon name="crosshairs" size={12} />
        </Button>
      </TooltipTrigger>
      <TooltipContent side="bottom">Scroll from Editor</TooltipContent>
    </Tooltip>
    <Tooltip>
      <TooltipTrigger>
        <Button
          variant="ghost"
          size="icon-xs"
          class="jb-icon-btn"
          onclick={() => explorerState.collapseAll()}
        >
          <Icon name="chevron-up" size={12} />
        </Button>
      </TooltipTrigger>
      <TooltipContent side="bottom">Collapse All</TooltipContent>
    </Tooltip>
    <Tooltip>
      <TooltipTrigger>
        <Button
          variant="ghost"
          size="icon-xs"
          class="jb-icon-btn"
          onclick={() => explorerState.expandAll()}
        >
          <Icon name="chevron-down" size={12} />
        </Button>
      </TooltipTrigger>
      <TooltipContent side="bottom">Expand All</TooltipContent>
    </Tooltip>
    <DropdownMenu>
      <DropdownMenuTrigger>
        <span
          class="jb-icon-btn menu-trigger"
          role="button"
          tabindex="0"
          title="Explorer options"
        >
          <Icon name="more" size={12} />
        </span>
      </DropdownMenuTrigger>
      <DropdownMenuContent>
        <DropdownMenuCheckboxItem
          checked={explorerState.scrollFromEditor}
          onCheckedChange={(v) => (explorerState.scrollFromEditor = v)}
        >
          Scroll from Editor
        </DropdownMenuCheckboxItem>
        <DropdownMenuCheckboxItem
          checked={explorerState.showSystemSchemas}
          onCheckedChange={() => explorerState.toggleSystemSchemas()}
        >
          Show system schemas
        </DropdownMenuCheckboxItem>
        <DropdownMenuCheckboxItem
          checked={explorerState.servicesVisible}
          onCheckedChange={() => explorerState.toggleServices()}
        >
          Services pane
        </DropdownMenuCheckboxItem>
      </DropdownMenuContent>
    </DropdownMenu>
    <Tooltip>
      <TooltipTrigger>
        <Button
          variant="ghost"
          size="icon-xs"
          class="jb-icon-btn"
          onclick={() => (explorerState.isSidebarCollapsed = true)}
        >
          <Icon name="x" size={12} />
        </Button>
      </TooltipTrigger>
      <TooltipContent side="bottom">Hide sidebar</TooltipContent>
    </Tooltip>
    <Tooltip>
      <TooltipTrigger>
        <Button
          variant="ghost"
          size="icon-xs"
          class="jb-icon-btn"
          onclick={() => explorerState.toggleServices()}
        >
          <Icon name="minus" size={12} />
        </Button>
      </TooltipTrigger>
      <TooltipContent side="bottom">Toggle Services pane</TooltipContent>
    </Tooltip>
  </Flex>
</Flex>

<!-- Action toolbar (DESIGN §3.1): new / refresh / copy DDL / DDL / eye -->
<Flex align="center" gap="1" class="explorer-toolbar">
  <DropdownMenu>
    <DropdownMenuTrigger>
      <span
        class="jb-icon-btn menu-trigger"
        role="button"
        tabindex="0"
        title="New data source / console"
      >
        <Icon name="plus" size={13} />
      </span>
    </DropdownMenuTrigger>
    <DropdownMenuContent>
      <DropdownMenuItem onclick={() => onnewconsole?.()}>
        New console on first connection
      </DropdownMenuItem>
      <DropdownMenuSeparator />
      <DropdownMenuItem onclick={() => onmanageconnections?.()}>
        Manage connections…
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
  <span class="toolbar-sep"></span>
  <Tooltip>
    <TooltipTrigger>
      <Button
        variant="ghost"
        size="icon-xs"
        class="jb-icon-btn"
        onclick={() => onreload?.()}
      >
        <Icon name="refresh" size={12} />
      </Button>
    </TooltipTrigger>
    <TooltipContent side="bottom">Refresh schema introspection</TooltipContent>
  </Tooltip>
  <Tooltip>
    <TooltipTrigger>
      <Button
        variant="ghost"
        size="icon-xs"
        class="jb-icon-btn"
        onclick={() => oncopyddl?.()}
      >
        <Icon name="copy" size={12} />
      </Button>
    </TooltipTrigger>
    <TooltipContent side="bottom">Copy DDL of current table</TooltipContent>
  </Tooltip>
  <span class="toolbar-sep"></span>
  <Button
    variant="ghost"
    size="icon-xs"
    class="jb-icon-btn ddl-label"
    title="View DDL of current table"
    onclick={() => (explorerState.ddlModalOpen = true)}
  >
    <Text size="xs" weight="bold" color="muted">DDL</Text>
  </Button>
  <span class="toolbar-sep"></span>
  <Tooltip>
    <TooltipTrigger>
      <Button
        variant="ghost"
        size="icon-xs"
        class="jb-icon-btn {explorerState.showSystemSchemas ? 'toggled' : ''}"
        onclick={() => explorerState.toggleSystemSchemas()}
      >
        <Icon
          name={explorerState.showSystemSchemas ? "eye" : "eye-off"}
          size={12}
        />
      </Button>
    </TooltipTrigger>
    <TooltipContent side="bottom">Show/hide system catalogs</TooltipContent>
  </Tooltip>
</Flex>

<style>
  :global(.explorer-header-top) {
    height: var(--toolbar-height, 32px);
    padding: 0 8px 0 12px;
    flex-shrink: 0;
  }

  :global(.explorer-toolbar) {
    min-height: 30px;
    padding: 2px 8px;
    border-bottom: 1px solid var(--border-subtle);
    flex-shrink: 0;
  }

  :global(.explorer-header-top .jb-icon-btn),
  :global(.explorer-toolbar .jb-icon-btn) {
    width: 22px !important;
    height: 22px !important;
    padding: 0 !important;
    min-width: 22px !important;
    color: var(--text-secondary);
    display: inline-flex;
    align-items: center;
    justify-content: center;
  }

  :global(.explorer-header-top .jb-icon-btn:hover),
  :global(.explorer-toolbar .jb-icon-btn:hover) {
    color: var(--text-primary);
    background-color: var(--bg-hover);
  }

  :global(.jb-icon-btn.toggled) {
    color: var(--text-primary) !important;
    background-color: var(--bg-hover) !important;
  }

  :global(.menu-trigger) {
    cursor: pointer;
    border-radius: var(--radius-sm, 4px);
  }

  :global(.toolbar-sep) {
    width: 1px;
    height: 16px;
    background-color: var(--border-default);
    margin: 0 3px;
    flex-shrink: 0;
  }

  :global(.explorer-toolbar .ddl-label) {
    width: auto !important;
    padding: 0 5px !important;
  }
</style>
