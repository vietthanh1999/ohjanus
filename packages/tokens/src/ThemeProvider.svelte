<script lang="ts">
  import { setContext, type Snippet } from 'svelte';

  export type ThemeMode = 'light' | 'dark' | 'system';

  interface Props {
    theme?: ThemeMode;
    children?: Snippet;
  }

  let { theme = 'dark', children }: Props = $props();

  let resolvedTheme = $state<'light' | 'dark'>('dark');

  $effect(() => {
    if (typeof window === 'undefined') return;

    if (theme === 'system') {
      const mql = window.matchMedia('(prefers-color-scheme: dark)');
      resolvedTheme = mql.matches ? 'dark' : 'light';

      const handler = (e: MediaQueryListEvent) => {
        resolvedTheme = e.matches ? 'dark' : 'light';
        document.documentElement.dataset.theme = resolvedTheme;
      };

      mql.addEventListener('change', handler);
      document.documentElement.dataset.theme = resolvedTheme;

      return () => mql.removeEventListener('change', handler);
    } else {
      resolvedTheme = theme;
      document.documentElement.dataset.theme = theme;
    }
  });

  setContext('ohjanus-theme', {
    get theme() {
      return theme;
    },
    get resolvedTheme() {
      return resolvedTheme;
    }
  });
</script>

{@render children?.()}
