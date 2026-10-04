import { commandPaletteState } from '@/features/command-palette';
import { consoleState } from '@/features/console';
import { explorerState } from '@/features/explorer';
import { settingsState } from '@/features/settings';
import { tableViewerState } from '@/features/table-viewer';
import { workbenchState } from '@/features/workbench';

let lastShiftTime = 0;

/** Global shortcuts for the App shell (wired via <svelte:window> in App.svelte). */
export function handleGlobalKeydown(e: KeyboardEvent) {
  // Cmd+K or Ctrl+K -> Search Everywhere
  if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
    e.preventDefault();
    commandPaletteState.open();
    return;
  }

  // Double Shift -> Search Everywhere
  if (e.key === 'Shift') {
    const now = Date.now();
    if (now - lastShiftTime < 300) {
      commandPaletteState.open();
    }
    lastShiftTime = now;
    return;
  }

  // Cmd+Enter or Ctrl+Enter -> Execute Query on the active data tab
  if ((e.metaKey || e.ctrlKey) && e.key === 'Enter') {
    e.preventDefault();
    const tab = workbenchState.activeTab;
    if (tab?.type === 'console') void consoleState.executeConsoleQuery();
    else if (tab?.type === 'table') void tableViewerState.loadTableData();
    return;
  }

  // Cmd+1 -> Toggle Sidebar
  if ((e.metaKey || e.ctrlKey) && e.key === '1') {
    e.preventDefault();
    explorerState.isSidebarCollapsed = !explorerState.isSidebarCollapsed;
    return;
  }

  // Cmd+, -> Settings
  if ((e.metaKey || e.ctrlKey) && e.key === ',') {
    e.preventDefault();
    settingsState.settingsModalOpen = true;
    return;
  }
}
