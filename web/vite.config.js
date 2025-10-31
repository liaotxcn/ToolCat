import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      '/auth': {
        target: 'http://localhost:8081',
        changeOrigin: true
      },
      '/plugins': {
        target: 'http://localhost:8081',
        changeOrigin: true
      },
      '/api/v1': {
        target: 'http://localhost:8081',
        changeOrigin: true
      },
      '/health': {
        target: 'http://localhost:8081',
        changeOrigin: true
      }
    }
  }
})
