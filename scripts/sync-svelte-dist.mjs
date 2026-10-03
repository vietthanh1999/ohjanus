// Mirrors every *.svelte file from src/ to dist/, preserving relative paths.
//
// Why: `tsc` compiles .ts → .js/.d.ts but never copies .svelte files, while the
// emitted dist/*.js files re-export their sibling .svelte components
// (e.g. `export { default as ThemeProvider } from './ThemeProvider.svelte'`).
// Without this step any resolver that lands on dist/ (older bundlers without
// the `svelte` export condition, SSR, published tarballs) fails with:
//   Failed to resolve import "./X.svelte" from ".../dist/index.js"
//
// Usage (run with cwd = package dir):
//   node ../../scripts/sync-svelte-dist.mjs
import { cpSync, existsSync, mkdirSync, readdirSync, statSync } from 'node:fs';
import { dirname, join, relative } from 'node:path';

const src = join(process.cwd(), 'src');
const dist = join(process.cwd(), 'dist');

function walk(dir) {
  for (const entry of readdirSync(dir)) {
    const full = join(dir, entry);
    if (statSync(full).isDirectory()) {
      walk(full);
    } else if (entry.endsWith('.svelte')) {
      const dest = join(dist, relative(src, full));
      mkdirSync(dirname(dest), { recursive: true });
      cpSync(full, dest);
      console.log('synced', relative(process.cwd(), dest));
    }
  }
}

if (!existsSync(src)) {
  console.error('no src/ directory in', process.cwd());
  process.exit(1);
}
walk(src);
