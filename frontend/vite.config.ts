import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:80',
        changeOrigin: true,
      },
    },
  },
  build: {
    // Output to backend/web so `go run ./cmd` can serve it directly.
    outDir: '../backend/web',
    emptyOutDir: true,
  },
})
