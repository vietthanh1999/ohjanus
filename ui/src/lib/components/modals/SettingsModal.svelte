<script lang="ts">
  import { appState } from '../../state/appState.svelte';
  import { toast, Select } from '@ohjanus/ui';
  import { Icon } from '@ohjanus/icons';

  // State
  let searchQuery = $state('');
  let activeCategory = $state<'appearance' | 'database' | 'tools' | 'editor' | 'keymap' | 'plugins'>('appearance');

  // Appearance settings matching screenshot
  let theme = $state('Islands Dark');
  let syncWithOS = $state(false);
  let editorScheme = $state('Islands Dark Theme default');
  let diffToolWindowBg = $state(false);

  // Accessibility
  let zoom = $state('100%');
  let useCustomFont = $state(false);
  let customFont = $state('Inter');
  let fontSize = $state(appState.uiDensity === 'comfortable' ? '14' : appState.uiDensity === 'standard' ? '13' : '12');
  let supportScreenReaders = $state(false);
  let contrastScrollbars = $state(false);
  let adjustVision = $state(false);

  // UI Options
  let compactMode = $state(appState.uiDensity === 'compact');
  let alwaysShowFullPath = $state(false);
  let useProjectColors = $state(true);
  let dragDropAltOnly = $state(false);
  let smoothScrolling = $state(true);
  let enableMnemonicsControls = $state(true);
  let enableMnemonicsMenu = $state(false);

  // Database / Gateway API
  let apiBase = $state('http://localhost:8788/api/v1');
  let queryTimeout = $state('30');
  let autoCommit = $state(false);

  function handleCompactToggle(checked: boolean) {
    compactMode = checked;
    if (checked) {
      fontSize = '12';
      appState.setDensity('compact');
    } else {
      fontSize = '14';
      appState.setDensity('comfortable');
    }
  }

  function handleFontSizeChange(size: string) {
    fontSize = size;
    if (size === '14' || parseInt(size, 10) >= 14) {
      compactMode = false;
      appState.setDensity('comfortable');
    } else if (size === '13') {
      compactMode = false;
      appState.setDensity('standard');
    } else {
      compactMode = true;
      appState.setDensity('compact');
    }
  }

  function close() {
    appState.settingsModalOpen = false;
  }

  function handleApply() {
    toast.success('Settings applied');
  }

  function handleSave() {
    close();
    toast.success('Settings saved successfully');
  }
</script>

{#if appState.settingsModalOpen}
  <!-- Backdrop -->
  <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
  <div
    class="settings-backdrop"
    role="dialog"
    aria-modal="true"
    tabindex="-1"
    onkeydown={(e) => { if (e.key === 'Escape') close(); }}
  >
    <!-- Window Container -->
    <div class="settings-window">
      <!-- Titlebar / macOS Header -->
      <div class="settings-titlebar">
        <!-- Traffic lights -->
        <div class="traffic-lights">
          <button type="button" class="light close" onclick={close} aria-label="Close settings"></button>
          <button type="button" class="light minimize" onclick={close} aria-label="Minimize"></button>
          <button type="button" class="light maximize" aria-label="Maximize"></button>
        </div>

        <!-- Center Window Title -->
        <div class="window-title">Settings</div>
        <div class="titlebar-dummy"></div>
      </div>

      <!-- Main Two-Column Body -->
      <div class="settings-body">
        <!-- Left Sidebar Navigation -->
        <aside class="settings-sidebar">
          <!-- Search input box with icon -->
          <div class="sidebar-search">
            <Icon name="search" size={12} color="#7A7E85" class="search-icon" />
            <input
              type="text"
              class="search-input"
              placeholder=""
              bind:value={searchQuery}
            />
          </div>

          <!-- Tree View Navigation -->
          <nav class="sidebar-tree">
            <!-- Database -->
            <button
              type="button"
              class="tree-row"
              class:active={activeCategory === 'database'}
              onclick={() => activeCategory = 'database'}
            >
              <span class="chevron"><Icon name="chevron-right" size={11} /></span>
              <span class="tree-label">Database</span>
            </button>

            <!-- Appearance & Behavior Group -->
            <div class="tree-group">
              <button
                type="button"
                class="tree-row group-parent active-parent"
                onclick={() => activeCategory = 'appearance'}
              >
                <span class="chevron expanded"><Icon name="chevron-down" size={11} /></span>
                <span class="tree-label font-medium">Appearance &amp; Behavior</span>
              </button>

              <div class="tree-children">
                <button
                  type="button"
                  class="tree-row depth-1"
                  class:active={activeCategory === 'appearance'}
                  onclick={() => activeCategory = 'appearance'}
                >
                  <span class="tree-label">Appearance</span>
                </button>

                <div class="tree-row depth-1">
                  <span class="tree-label text-muted">Menus and Toolbars</span>
                </div>

                <div class="tree-row depth-1">
                  <span class="chevron"><Icon name="chevron-right" size={11} /></span>
                  <span class="tree-label text-muted">System Settings</span>
                </div>

                <div class="tree-row depth-1 with-icon">
                  <span class="tree-label text-muted">File Colors</span>
                  <span class="tree-icon-badge">▫</span>
                </div>

                <div class="tree-row depth-1 with-icon">
                  <span class="tree-label text-muted">Scopes</span>
                  <span class="tree-icon-badge">▫</span>
                </div>

                <div class="tree-row depth-1">
                  <span class="tree-label text-muted">Notifications</span>
                </div>

                <div class="tree-row depth-1">
                  <span class="tree-label text-muted">Data Editor and Viewer</span>
                </div>

                <div class="tree-row depth-1">
                  <span class="tree-label text-muted">Quick Lists</span>
                </div>

                <div class="tree-row depth-1 with-icon">
                  <span class="tree-label text-muted">Required Plugins</span>
                  <span class="tree-icon-badge">▫</span>
                </div>

                <div class="tree-row depth-1">
                  <span class="tree-label text-muted">Trusted Locations</span>
                </div>

                <div class="tree-row depth-1">
                  <span class="tree-label text-muted">Presentation Assistant</span>
                </div>
              </div>
            </div>

            <!-- Keymap -->
            <button
              type="button"
              class="tree-row"
              class:active={activeCategory === 'keymap'}
              onclick={() => activeCategory = 'keymap'}
            >
              <span class="tree-label">Keymap</span>
            </button>

            <!-- Editor -->
            <button
              type="button"
              class="tree-row"
              class:active={activeCategory === 'editor'}
              onclick={() => activeCategory = 'editor'}
            >
              <span class="chevron"><Icon name="chevron-right" size={11} /></span>
              <span class="tree-label">Editor</span>
            </button>

            <!-- Plugins -->
            <button
              type="button"
              class="tree-row with-badge"
              class:active={activeCategory === 'plugins'}
              onclick={() => activeCategory = 'plugins'}
            >
              <span class="tree-label">Plugins</span>
              <span class="plugin-pill">2</span>
            </button>

            <!-- Version Control -->
            <div class="tree-row with-icon">
              <span class="chevron"><Icon name="chevron-right" size={11} /></span>
              <span class="tree-label text-muted">Version Control</span>
              <span class="tree-icon-badge">▫</span>
            </div>

            <!-- Languages -->
            <div class="tree-row">
              <span class="chevron"><Icon name="chevron-right" size={11} /></span>
              <span class="tree-label text-muted">Languages</span>
            </div>

            <!-- Tools -->
            <button
              type="button"
              class="tree-row"
              class:active={activeCategory === 'tools'}
              onclick={() => activeCategory = 'tools'}
            >
              <span class="chevron"><Icon name="chevron-right" size={11} /></span>
              <span class="tree-label">Tools</span>
            </button>

            <!-- Backup and Sync -->
            <div class="tree-row">
              <span class="tree-label text-muted">Backup and Sync</span>
            </div>

            <!-- Advanced Settings -->
            <div class="tree-row">
              <span class="tree-label text-muted">Advanced Settings</span>
            </div>
          </nav>
        </aside>

        <!-- Right Content Area -->
        <main class="settings-content-area">
          <!-- Breadcrumb Row -->
          <div class="content-header">
            <div class="breadcrumb">
              <span class="crumb-parent">Appearance &amp; Behavior</span>
              <span class="crumb-separator">&rsaquo;</span>
              <span class="crumb-current">
                {#if activeCategory === 'appearance'}
                  Appearance
                {:else if activeCategory === 'database'}
                  Database
                {:else if activeCategory === 'tools'}
                  Gateway Tools &amp; API
                {:else}
                  {activeCategory.toUpperCase()}
                {/if}
              </span>
            </div>

            <!-- Navigation Arrows -->
            <div class="header-nav-arrows">
              <button type="button" class="nav-arrow-btn" title="Back"><Icon name="chevron-left" size={12} /></button>
              <button type="button" class="nav-arrow-btn" title="Forward"><Icon name="chevron-right" size={12} /></button>
            </div>
          </div>

          <!-- Main Scrollable Pane -->
          <div class="content-scrollable">
            {#if activeCategory === 'appearance'}
              <!-- 1. THEME SECTION -->
              <section class="settings-section">
                <!-- Theme Row -->
                <div class="form-row">
                  <span class="form-label">Theme:</span>
                  <div class="control-combo">
                    <Select
                      class="theme-select"
                      options={[
                        { value: 'Islands Dark', label: 'Islands Dark' },
                        { value: 'JetBrains Dark', label: 'JetBrains Dark' },
                        { value: 'Darcula', label: 'Darcula' },
                        { value: 'Light', label: 'Light' }
                      ]}
                      bind:value={theme}
                    />

                    <label class="jb-checkbox-label">
                      <input type="checkbox" bind:checked={syncWithOS} />
                      <span>Sync with OS</span>
                    </label>

                    <button type="button" class="gear-btn" title="Theme Settings">
                      <Icon name="settings" size={14} />
                    </button>
                  </div>
                </div>

                <!-- Editor Color Scheme Row -->
                <div class="form-row">
                  <span class="form-label">Editor color scheme:</span>
                  <div class="control-combo">
                    <Select
                      class="scheme-select"
                      options={[
                        { value: 'Islands Dark Theme default', label: 'Islands Dark Theme default' },
                        { value: 'Darcula', label: 'Darcula' },
                        { value: 'High Contrast', label: 'High Contrast' }
                      ]}
                      bind:value={editorScheme}
                    />

                    <button type="button" class="gear-btn" title="Scheme Settings">
                      <Icon name="settings" size={14} />
                    </button>
                  </div>
                </div>

                <!-- Different tool window background -->
                <div class="form-row single-checkbox-row">
                  <label class="jb-checkbox-label">
                    <input type="checkbox" bind:checked={diffToolWindowBg} />
                    <span class="chk-text">Different tool window background</span>
                  </label>
                  <p class="field-subtext">Use lighter color in the dark theme and darker color in the light theme as a background</p>
                </div>
              </section>

              <!-- 2. ACCESSIBILITY SECTION -->
              <section class="settings-section">
                <div class="section-title">Accessibility</div>

                <!-- Zoom Row -->
                <div class="form-row">
                  <span class="form-label">Zoom:</span>
                  <div class="control-combo">
                    <Select
                      class="short-select"
                      options={[
                        { value: '100%', label: '100%' },
                        { value: '110%', label: '110%' },
                        { value: '125%', label: '125%' },
                        { value: '150%', label: '150%' }
                      ]}
                      bind:value={zoom}
                    />
                    <span class="hint-inline">Change with ^⌥= or ^⌥-. Set to 100% with ^⌥0</span>
                  </div>
                </div>

                <!-- Use custom font & Size Row -->
                <div class="form-row">
                  <label class="jb-checkbox-label">
                    <input type="checkbox" bind:checked={useCustomFont} />
                    <span class="chk-text">Use custom font:</span>
                  </label>
                  <div class="control-combo" style="margin-left: 8px;">
                    <Select
                      class="font-select"
                      options={[
                        { value: 'Inter', label: 'Inter' },
                        { value: 'JetBrains Mono', label: 'JetBrains Mono' },
                        { value: 'System', label: 'System Default' }
                      ]}
                      bind:value={customFont}
                      disabled={!useCustomFont}
                    />

                    <span class="inline-lbl">Size:</span>
                    <Select
                      class="tiny-select"
                      options={[
                        { value: '11', label: '11' },
                        { value: '12', label: '12' },
                        { value: '13', label: '13' },
                        { value: '14', label: '14' },
                        { value: '15', label: '15' },
                        { value: '16', label: '16' },
                        { value: '18', label: '18' }
                      ]}
                      bind:value={fontSize}
                      disabled={!useCustomFont}
                      onchange={(v) => handleFontSizeChange(v as string)}
                    />
                  </div>
                </div>

                <!-- Support screen readers -->
                <div class="form-row single-checkbox-row">
                  <label class="jb-checkbox-label">
                    <input type="checkbox" bind:checked={supportScreenReaders} />
                    <span class="chk-text">Support screen readers</span>
                    <span class="restart-pill">Requires restart</span>
                  </label>
                  <p class="field-subtext">
                    ^➔| and ^⇧➔| will navigate UI controls in dialogs and will not be available for switching editor tabs or other IDE actions. Tooltips on mouse hover will be disabled.
                  </p>
                </div>

                <!-- Use contrast scrollbars -->
                <div class="form-row single-checkbox-row">
                  <label class="jb-checkbox-label">
                    <input type="checkbox" bind:checked={contrastScrollbars} />
                    <span class="chk-text">Use contrast scrollbars</span>
                  </label>
                </div>

                <!-- Adjust colors for red-green vision deficiency -->
                <div class="form-row single-checkbox-row">
                  <label class="jb-checkbox-label">
                    <input type="checkbox" bind:checked={adjustVision} />
                    <span class="chk-text">Adjust colors for red-green vision deficiency</span>
                    <!-- svelte-ignore a11y_invalid_attribute -->
                    <a href="#" class="blue-link" onclick={(e) => e.preventDefault()}>How it works</a>
                  </label>
                  <p class="field-subtext">Requires restart. For protanopia and deuteranopia.</p>
                </div>
              </section>

              <!-- 3. UI OPTIONS SECTION -->
              <section class="settings-section">
                <div class="section-title">UI Options</div>

                <div class="ui-options-grid">
                  <!-- Left column -->
                  <div class="options-col">
                    <div class="option-item">
                      <label class="jb-checkbox-label">
                        <input
                          type="checkbox"
                          checked={compactMode}
                          onchange={(e) => handleCompactToggle((e.target as HTMLInputElement).checked)}
                        />
                        <span class="chk-text">Compact mode</span>
                      </label>
                      <p class="field-subtext">UI elements take up less screen space</p>
                    </div>

                    <div class="option-item">
                      <label class="jb-checkbox-label">
                        <input type="checkbox" bind:checked={alwaysShowFullPath} />
                        <span class="chk-text">Always show full path in window header</span>
                      </label>
                    </div>

                    <div class="option-item">
                      <label class="jb-checkbox-label">
                        <input type="checkbox" bind:checked={useProjectColors} />
                        <span class="chk-text">Use project colors in main toolbar</span>
                      </label>
                      <p class="field-subtext">Distinguish projects with different toolbar</p>
                    </div>
                  </div>

                  <!-- Right column -->
                  <div class="options-col">
                    <div class="option-item">
                      <label class="jb-checkbox-label">
                        <input type="checkbox" bind:checked={dragDropAltOnly} />
                        <span class="chk-text">Drag-and-drop with Alt pressed only</span>
                      </label>
                    </div>

                    <div class="option-item">
                      <label class="jb-checkbox-label">
                        <input type="checkbox" bind:checked={smoothScrolling} />
                        <span class="chk-text">Smooth scrolling</span>
                        <span class="info-icon" title="Smooth scrolling helps readability when scrolling code">ⓘ</span>
                      </label>
                    </div>

                    <div class="option-item">
                      <label class="jb-checkbox-label">
                        <input type="checkbox" bind:checked={enableMnemonicsControls} />
                        <span class="chk-text">Enable mnemonics in controls</span>
                      </label>
                    </div>

                    <div class="option-item">
                      <label class="jb-checkbox-label">
                        <input type="checkbox" bind:checked={enableMnemonicsMenu} />
                        <span class="chk-text">Enable mnemonics in menu</span>
                      </label>
                    </div>
                  </div>
                </div>
              </section>
            {:else if activeCategory === 'database' || activeCategory === 'tools'}
              <!-- Database & Gateway API settings -->
              <section class="settings-section">
                <div class="section-title">Janus Gateway Admin API Connection</div>

                <div class="form-row vertical">
                  <label for="admin-api-base" class="form-label font-medium">Gateway Admin API Base URL:</label>
                  <input
                    id="admin-api-base"
                    type="text"
                    class="jb-input code-text"
                    bind:value={apiBase}
                    placeholder="http://localhost:8788/api/v1"
                  />
                  <p class="field-subtext">Port 8788 is the dedicated Janus Admin REST API port separated from MCP clients.</p>
                </div>

                <div class="form-row vertical">
                  <label for="query-timeout" class="form-label font-medium">Max Query Execution Timeout (seconds):</label>
                  <input
                    id="query-timeout"
                    type="number"
                    class="jb-input code-text"
                    bind:value={queryTimeout}
                    placeholder="30"
                    style="max-width: 120px;"
                  />
                  <p class="field-subtext">Default timeout before gateway terminates long-running queries.</p>
                </div>

                <div class="form-row single-checkbox-row">
                  <label class="jb-checkbox-label">
                    <input type="checkbox" bind:checked={autoCommit} />
                    <span class="chk-text">Enable Auto-Commit on write queries</span>
                  </label>
                  <p class="field-subtext" style="color: var(--action-danger, #E55353);">Not recommended for production environments.</p>
                </div>
              </section>
            {:else}
              <div class="empty-category-pane">
                <span style="font-size: 14px; font-weight: 500; color: var(--text-muted);">
                  Settings category for {activeCategory} is available in enterprise profile.
                </span>
              </div>
            {/if}
          </div>
        </main>
      </div>

      <!-- Window Footer -->
      <footer class="settings-footer">
        <button type="button" class="help-btn" title="Help">?</button>

        <div class="footer-actions">
          <button type="button" class="jb-btn jb-btn-secondary" onclick={close}>Cancel</button>
          <button type="button" class="jb-btn jb-btn-secondary" onclick={handleApply}>Apply</button>
          <button type="button" class="jb-btn jb-btn-primary" onclick={handleSave}>OK</button>
        </div>
      </footer>
    </div>
  </div>
{/if}

<style>
  .settings-backdrop {
    position: fixed;
    inset: 0;
    background-color: rgba(0, 0, 0, 0.65);
    backdrop-filter: blur(2px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }

  .settings-window {
    width: 960px;
    height: 680px;
    max-width: 96vw;
    max-height: 94vh;
    background-color: #191A1C;
    border: 1px solid var(--border-strong, #4E5157);
    border-radius: 8px;
    box-shadow: 0 16px 48px rgba(0, 0, 0, 0.75);
    display: flex;
    flex-direction: column;
    overflow: hidden;
    color: var(--text-primary, #DFE1E5);
    font-family: var(--font-ui);
    font-size: var(--font-size-base, 13px);
  }

  /* Titlebar */
  .settings-titlebar {
    height: 38px;
    background-color: #26282B;
    border-bottom: 1px solid #323438;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 12px;
    flex-shrink: 0;
    user-select: none;
  }

  .traffic-lights {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .light {
    width: 12px;
    height: 12px;
    border-radius: 50%;
    cursor: pointer;
  }
  .light.close { background-color: var(--window-close, #FF5F56); }
  .light.minimize { background-color: var(--window-minimize, #FFBD2E); }
  .light.maximize { background-color: var(--window-maximize, #27C93F); }

  .window-title {
    font-weight: 500;
    font-size: 13px;
    color: var(--text-primary, #DFE1E5);
  }

  .titlebar-dummy {
    width: 52px;
  }

  /* Body Two Columns */
  .settings-body {
    flex: 1;
    display: flex;
    overflow: hidden;
  }

  /* Left Sidebar */
  .settings-sidebar {
    width: 240px;
    background-color: #191A1C;
    border-right: 1px solid #323438;
    display: flex;
    flex-direction: column;
    flex-shrink: 0;
  }

  .sidebar-search {
    padding: 8px 10px;
    position: relative;
    border-bottom: 1px solid #323438;
  }

  :global(.search-icon) {
    position: absolute;
    left: 18px;
    top: 50%;
    transform: translateY(-50%);
  }

  .search-input {
    width: 100%;
    height: 28px;
    background-color: #26282B;
    border: 1px solid #4E5157;
    border-radius: 4px;
    padding: 0 8px 0 28px;
    color: var(--text-primary, #DFE1E5);
    font-size: 12.5px;
    outline: none;
    transition: border-color 0.15s ease;
  }

  .search-input:hover {
    border-color: #6C707E;
  }

  .search-input:focus {
    border-color: var(--action-primary, #3574F0);
    box-shadow: 0 0 0 1px var(--action-primary, #3574F0);
  }

  .sidebar-tree {
    flex: 1;
    overflow-y: auto;
    padding: 6px 4px;
  }

  .tree-row {
    width: 100%;
    height: 26px;
    display: flex;
    align-items: center;
    padding: 0 8px;
    border-radius: 4px;
    cursor: pointer;
    font-size: 12.5px;
    color: var(--text-primary, #DFE1E5);
    text-align: left;
    transition: background-color 0.1s ease;
  }

  .tree-row:hover {
    background-color: #313438;
  }

  .tree-row.active {
    background-color: #2E3A4E;
    color: #FFFFFF;
    font-weight: 500;
  }

  .tree-row.depth-1 {
    padding-left: 24px;
  }

  .tree-row.with-icon, .tree-row.with-badge {
    justify-content: space-between;
  }

  .tree-label {
    flex: 1;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .tree-label.text-muted {
    color: #9DA0A8;
  }

  .tree-icon-badge {
    font-size: 10px;
    color: #7A7E85;
    margin-left: 4px;
  }

  .plugin-pill {
    background-color: #393B40;
    color: #DFE1E5;
    font-size: 10px;
    font-weight: 600;
    padding: 0 6px;
    border-radius: 10px;
    line-height: 16px;
  }

  .chevron {
    width: 14px;
    display: inline-flex;
    align-items: center;
    color: #7A7E85;
    flex-shrink: 0;
  }

  /* Right Content Area */
  .settings-content-area {
    flex: 1;
    display: flex;
    flex-direction: column;
    background-color: #191A1C;
    overflow: hidden;
  }

  .content-header {
    height: 38px;
    padding: 0 20px;
    border-bottom: 1px solid #323438;
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-shrink: 0;
  }

  .breadcrumb {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 13px;
    font-weight: 600;
  }

  .crumb-parent {
    color: var(--text-primary);
  }

  .crumb-separator {
    color: #7A7E85;
    font-size: 14px;
  }

  .crumb-current {
    color: var(--text-primary);
  }

  .header-nav-arrows {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .nav-arrow-btn {
    width: 24px;
    height: 24px;
    border-radius: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #7A7E85;
  }

  .nav-arrow-btn:hover {
    background-color: #35373B;
    color: var(--text-primary);
  }

  .content-scrollable {
    flex: 1;
    overflow-y: auto;
    padding: 16px 24px;
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .settings-section {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .section-title {
    font-size: 13px;
    font-weight: 600;
    color: var(--text-primary);
    padding-bottom: 6px;
    border-bottom: 1px solid #323438;
    margin-top: 4px;
  }

  /* Form controls */
  .form-row {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .form-row.vertical {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }

  .form-row.single-checkbox-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 2px;
  }

  .form-label {
    min-width: 140px;
    font-size: 12.5px;
    color: var(--text-primary);
  }

  .control-combo {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
  }


  :global(.theme-select) {
    min-width: 170px !important;
  }

  :global(.scheme-select) {
    min-width: 270px !important;
  }

  :global(.short-select) {
    min-width: 90px !important;
    width: 90px !important;
  }

  :global(.font-select) {
    min-width: 180px !important;
  }

  :global(.tiny-select) {
    min-width: 68px !important;
    width: 68px !important;
  }

  .jb-input {
    width: 100%;
    height: 28px;
    background-color: #26282B;
    border: 1px solid #4E5157;
    border-radius: 4px;
    color: var(--text-primary);
    padding: 0 10px;
    font-size: 12.5px;
    outline: none;
    transition: border-color 0.15s ease;
  }

  .jb-input:hover {
    border-color: #6C707E;
  }

  .jb-input:focus {
    border-color: var(--action-primary, #3574F0);
    box-shadow: 0 0 0 1px var(--action-primary, #3574F0);
  }

  .jb-checkbox-label {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    font-size: 12.5px;
    cursor: pointer;
    user-select: none;
    color: var(--text-primary, #DFE1E5);
  }

  .jb-checkbox-label input[type="checkbox"] {
    appearance: none;
    -webkit-appearance: none;
    -moz-appearance: none;
    width: 15px;
    height: 15px;
    border: 1px solid #4E5157;
    border-radius: 3px;
    background-color: #26282B;
    cursor: pointer;
    margin: 0;
    padding: 0;
    display: inline-grid;
    place-content: center;
    position: relative;
    flex-shrink: 0;
    transition: background-color 0.12s ease, border-color 0.12s ease;
  }

  .jb-checkbox-label input[type="checkbox"]:hover {
    border-color: #6C707E;
  }

  .jb-checkbox-label input[type="checkbox"]:checked {
    background-color: #3574F0;
    border-color: #3574F0;
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 16 16' fill='none'%3E%3Cpath d='M3.5 8.5L6.5 11.5L12.5 4.5' stroke='white' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'/%3E%3C/svg%3E");
    background-repeat: no-repeat;
    background-position: center;
    background-size: 11px 11px;
  }

  .jb-checkbox-label input[type="checkbox"]:focus-visible {
    outline: none;
    box-shadow: 0 0 0 2px rgba(53, 116, 240, 0.4);
  }

  .jb-checkbox-label input[type="checkbox"]:disabled {
    background-color: #1E1F22;
    border-color: #393B40;
    opacity: 0.5;
    cursor: not-allowed;
  }

  .field-subtext {
    font-size: 11.5px;
    color: #7A7E85;
    margin-left: 22px;
    line-height: 1.4;
  }

  .hint-inline {
    font-size: 11.5px;
    color: #7A7E85;
    margin-left: 6px;
  }

  .inline-lbl {
    font-size: 12.5px;
    color: var(--text-primary);
    margin-left: 4px;
  }

  .gear-btn {
    width: 26px;
    height: 26px;
    border-radius: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #7A7E85;
  }

  .gear-btn:hover {
    background-color: #35373B;
    color: var(--text-primary);
  }

  .restart-pill {
    font-size: 10.5px;
    color: #7A7E85;
    background-color: #1E1F22;
    border: 1px solid #35373B;
    border-radius: 3px;
    padding: 1px 6px;
    margin-left: 4px;
  }

  .blue-link {
    color: #56A8F5;
    text-decoration: none;
    margin-left: 4px;
    font-size: 12px;
  }

  .blue-link:hover {
    text-decoration: underline;
  }

  .info-icon {
    font-size: 12px;
    color: #7A7E85;
    margin-left: 4px;
    cursor: help;
  }

  /* UI Options Grid */
  .ui-options-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px 24px;
  }

  .options-col {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .option-item {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .empty-category-pane {
    padding: 40px;
    text-align: center;
  }

  /* Footer */
  .settings-footer {
    height: 48px;
    background-color: #26282B;
    border-top: 1px solid #323438;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 16px;
    flex-shrink: 0;
  }

  .help-btn {
    width: 24px;
    height: 24px;
    border-radius: 50%;
    border: 1px solid var(--border-default, #393B40);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 12px;
    font-weight: 600;
    color: #7A7E85;
    transition: all 0.15s ease;
  }

  .help-btn:hover {
    background-color: #35373B;
    color: var(--text-primary);
  }

  .footer-actions {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .jb-btn {
    height: 28px;
    padding: 0 16px;
    border-radius: 4px;
    font-size: 12.5px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.15s ease;
  }

  .jb-btn-secondary {
    background-color: #393B40;
    color: var(--text-primary);
    border: 1px solid #4E5157;
  }

  .jb-btn-secondary:hover {
    background-color: #45474D;
  }

  .jb-btn-primary {
    background-color: var(--action-primary, #3574F0);
    color: #FFFFFF;
    border: 1px solid var(--action-primary, #3574F0);
  }

  .jb-btn-primary:hover {
    background-color: var(--action-primary-hover, #4682F7);
  }
</style>
