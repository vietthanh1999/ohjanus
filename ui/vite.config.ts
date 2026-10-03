import { svelte } from '@sveltejs/vite-plugin-svelte'
import { defineConfig } from 'vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [svelte()],
  server: {
    port: 5173,
    proxy: {
      // Same-origin /api in dev -> Admin API :8788 (no CORS, no hardcoded host).
      '/api': {
        target: process.env.VITE_API_BASE_URL ?? 'http://127.0.0.1:8788',
        changeOrigin: true
      }
    }
  }
})
