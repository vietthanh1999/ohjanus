class CommandPaletteManager {
  searchModalOpen = $state<boolean>(false);

  open() {
    this.searchModalOpen = true;
  }

  close() {
    this.searchModalOpen = false;
  }
}

export const commandPaletteState = new CommandPaletteManager();
