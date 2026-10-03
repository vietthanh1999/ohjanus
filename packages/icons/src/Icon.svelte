<script lang="ts">
  type IconName =
    | 'database'
    | 'table'
    | 'shield'
    | 'key'
    | 'audit'
    | 'check'
    | 'x'
    | 'chevron-down'
    | 'chevron-right'
    | 'chevron-left'
    | 'chevron-up'
    | 'search'
    | 'settings'
    | 'refresh'
    | 'trash'
    | 'play'
    | 'lock'
    | 'alert-triangle'
    | 'eye'
    | 'plus'
    | 'minus'
    | 'download'
    | 'upload'
    | 'copy'
    | 'clock'
    | 'terminal'
    | 'info'
    | 'filter';

  interface Props {
    name: IconName;
    size?: number | string;
    color?: string;
    class?: string;
    [key: string]: any;
  }

  let {
    name,
    size = 16,
    color = 'currentColor',
    class: className = '',
    ...restProps
  }: Props = $props();
</script>

{#if name === 'database'}
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke={color} stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="ohjanus-icon {className}" {...restProps}>
    <ellipse cx="12" cy="5" rx="9" ry="3" />
    <path d="M3 5V19A9 3 0 0 0 21 19V5" />
    <path d="M3 12A9 3 0 0 0 21 12" />
  </svg>
{:else if name === 'table'}
  <svg width={size} height={size} viewBox="0 0 16 16" fill={color} class="ohjanus-icon {className}" {...restProps}>
    <path d="M0 2a1 1 0 011-1h14a1 1 0 011 1v12a1 1 0 01-1 1H1a1 1 0 01-1-1V2zm1 3v2h6V5H1zm7 0v2h7V5H8zm0 3v2h7V8H8zm-1 0H1v2h6V8zm0 3H1v2h6v-2zm1 0v2h7v-2H8z"/>
  </svg>
{:else if name === 'shield'}
  <svg width={size} height={size} viewBox="0 0 16 16" fill={color} class="ohjanus-icon {className}" {...restProps}>
    <path fill-rule="evenodd" d="M8 0c-.69 0-1.843.265-2.928.56-1.11.3-2.229.655-2.887.87a1.54 1.54 0 00-1.044 1.262c-.596 4.477.787 7.795 2.464 9.99 1.579 2.065 3.444 3.009 4.395 3.318.066.022.135.034.204.034s.138-.012.204-.034c.951-.309 2.816-1.253 4.395-3.318 1.677-2.195 3.06-5.513 2.464-9.99a1.54 1.54 0 00-1.044-1.263 62.467 62.467 0 00-2.887-.87C9.843.266 8.69 0 8 0zm2.146 5.146a.5.5 0 01.708.708l-3 3a.5.5 0 01-.708 0l-1.5-1.5a.5.5 0 11.708-.708L7.5 7.793l2.646-2.647z"/>
  </svg>
{:else if name === 'key'}
  <svg width={size} height={size} viewBox="0 0 16 16" fill={color} class="ohjanus-icon {className}" {...restProps}>
    <path d="M0 8a4 4 0 017.465-2H14a2 2 0 012 2v1a1 1 0 01-1 1h-1v1a1 1 0 01-1 1h-1v1a1 1 0 01-1 1H9.465A4 4 0 010 8zm4-2a2 2 0 100 4 2 2 0 000-4z"/>
  </svg>
{:else if name === 'audit'}
  <svg width={size} height={size} viewBox="0 0 16 16" fill={color} class="ohjanus-icon {className}" {...restProps}>
    <path d="M4 1.5H3a2 2 0 00-2 2V14a2 2 0 002 2h10a2 2 0 002-2V3.5a2 2 0 00-2-2h-1v1h1a1 1 0 011 1V14a1 1 0 01-1 1H3a1 1 0 01-1-1V3.5a1 1 0 011-1h1v-1z"/>
  </svg>
{:else if name === 'check'}
  <svg width={size} height={size} viewBox="0 0 16 16" fill="none" stroke={color} stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="ohjanus-icon {className}" {...restProps}>
    <polyline points="3 8 6.5 11.5 13 4"/>
  </svg>
{:else if name === 'x'}
  <svg width={size} height={size} viewBox="0 0 16 16" fill="none" stroke={color} stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="ohjanus-icon {className}" {...restProps}>
    <line x1="4" y1="4" x2="12" y2="12"/>
    <line x1="12" y1="4" x2="4" y2="12"/>
  </svg>
{:else if name === 'chevron-down'}
  <svg width={size} height={size} viewBox="0 0 16 16" fill="currentColor" class="ohjanus-icon {className}" {...restProps}>
    <path d="M4 6l4 4 4-4H4z"/>
  </svg>
{:else if name === 'chevron-right'}
  <svg width={size} height={size} viewBox="0 0 16 16" fill="currentColor" class="ohjanus-icon {className}" {...restProps}>
    <path d="M6 4l4 4-4 4V4z"/>
  </svg>
{:else if name === 'search'}
  <svg width={size} height={size} viewBox="0 0 16 16" fill="currentColor" class="ohjanus-icon {className}" {...restProps}>
    <path fill-rule="evenodd" d="M11.5 7a4.5 4.5 0 11-9 0 4.5 4.5 0 019 0zm-.82 4.74a6 6 0 111.06-1.06l3.04 3.04a.75.75 0 11-1.06 1.06l-3.04-3.04z"/>
  </svg>
{:else if name === 'settings'}
  <svg width={size} height={size} viewBox="0 0 16 16" fill="currentColor" class="ohjanus-icon {className}" {...restProps}>
    <path d="M7.07 1a1 1 0 00-.97.757l-.24 1.026a5.526 5.526 0 00-1.282.74L3.6 3.057a1 1 0 00-1.205.247l-.98 1.135a1 1 0 00-.173 1.218l.617.848a5.534 5.534 0 00-.012 1.48l-.618.847a1 1 0 00.173 1.218l.98 1.136a1 1 0 001.206.246l.978-.466c.394.3.826.55 1.282.74l.24 1.026A1 1 0 007.07 15h1.86a1 1 0 00.97-.757l.24-1.026c.456-.19.888-.44 1.282-.74l.978.466a1 1 0 001.206-.246l.98-1.136a1 1 0 00-.173-1.218l-.617-.847c.105-.486.105-.993 0-1.48l.617-.848a1 1 0 00.173-1.218l-.98-1.135a1 1 0 00-1.206-.247l-.978.466a5.527 5.527 0 00-1.282-.74l-.24-1.026A1 1 0 008.93 1H7.07zm.93 5a2 2 0 110 4 2 2 0 010-4z"/>
  </svg>
{:else if name === 'refresh'}
  <svg width={size} height={size} viewBox="0 0 16 16" fill="currentColor" class="ohjanus-icon {className}" {...restProps}>
    <path fill-rule="evenodd" d="M8 3a5 5 0 104.546 2.914.5.5 0 01.908-.417A6 6 0 118 2v1z"/>
    <path d="M8 4.466V.534a.25.25 0 01.41-.192l2.36 1.966c.12.1.12.284 0 .384L8.41 4.658A.25.25 0 018 4.466z"/>
  </svg>
{:else if name === 'trash'}
  <svg width={size} height={size} viewBox="0 0 16 16" fill="none" stroke={color} stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="ohjanus-icon {className}" {...restProps}>
    <polyline points="2 4 4 4 14 4"/>
    <path d="M5 4V2a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2"/>
    <path d="M6 7v6"/>
    <path d="M10 7v6"/>
  </svg>
{:else if name === 'play'}
  <svg width={size} height={size} viewBox="0 0 16 16" fill={color} class="ohjanus-icon {className}" {...restProps}>
    <path d="M4.5 2.5v11l9-5.5-9-5.5z"/>
  </svg>
{:else if name === 'lock'}
  <svg width={size} height={size} viewBox="0 0 16 16" fill="currentColor" class="ohjanus-icon {className}" {...restProps}>
    <path fill-rule="evenodd" d="M4 4v2h-.5A1.5 1.5 0 002 7.5v6A1.5 1.5 0 003.5 15h9a1.5 1.5 0 001.5-1.5v-6A1.5 1.5 0 0012.5 6H12V4a4 4 0 00-8 0zm6.5 2V4a2.5 2.5 0 00-5 0v2h5z"/>
  </svg>
{:else if name === 'alert-triangle'}
  <svg width={size} height={size} viewBox="0 0 16 16" fill="none" stroke={color} stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="ohjanus-icon {className}" {...restProps}>
    <path d="M8 1l7 13H1L8 1z"/>
    <line x1="8" y1="6" x2="8" y2="9"/>
    <line x1="8" y1="12" x2="8.01" y2="12"/>
  </svg>
{:else}
  <svg width={size} height={size} viewBox="0 0 16 16" fill="currentColor" class="ohjanus-icon {className}" {...restProps}>
    <circle cx="8" cy="8" r="7"/>
  </svg>
{/if}
