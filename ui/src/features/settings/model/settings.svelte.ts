import { UI_DENSITY_KEY } from '@/shared/config';

export type UiDensity = 'compact' | 'standard' | 'comfortable';

class SettingsManager {
  uiDensity = $state<UiDensity>('comfortable');
  settingsModalOpen = $state<boolean>(false);

  setDensity(density: UiDensity) {
    this.uiDensity = density;
    if (typeof document !== 'undefined') {
      document.documentElement.setAttribute('data-density', density);
      try {
        localStorage.setItem(UI_DENSITY_KEY, density);
      } catch {
        // Ignore storage errors in restricted contexts
      }
    }
  }

  initDensity() {
    if (typeof document !== 'undefined') {
      try {
        const saved = localStorage.getItem(UI_DENSITY_KEY) as UiDensity | null;
        this.uiDensity = saved || 'comfortable';
      } catch {
        this.uiDensity = 'comfortable';
      }
      document.documentElement.setAttribute('data-density', this.uiDensity);
    }
  }
}

export const settingsState = new SettingsManager();
