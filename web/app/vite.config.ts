import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// In dev the SPA runs at :5173 (base '/') and proxies API calls to the Go
// server on :8080. The production build uses base '/app/' so the Go server can
// serve it under /app/ alongside the marketing landing at /.
export default defineConfig(({ mode }) => ({
  base: mode === 'production' ? '/app/' : '/',
  plugins: [svelte()],
  server: {
    port: 5173,
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
  build: {
    outDir: 'dist',
    emptyOutDir: true,
  },
}))
