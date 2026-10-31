import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api/v1/ws': {
        target: 'ws://localhost:80',
        ws: true,
        rewriteWsOrigin: true,
      },
      '/api': 'http://localhost:80',
    },
  },
})
