import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/files': {
        target: 'http://localhost:9000',
        changeOrigin: true,
      },
      '/users': {
        target: 'http://localhost:9000',
        changeOrigin: true,
      },
      '/share': {
        target: 'http://localhost:9000',
        changeOrigin: true,
        bypass(req) {
          if (req.url && req.url.includes('download=true')) {
            return null;
          }
          if (req.headers.accept && req.headers.accept.includes('text/html')) {
            return '/index.html';
          }
        },
      },
      '/health': {
        target: 'http://localhost:9000',
        changeOrigin: true,
      },
      '/auth-test': {
        target: 'http://localhost:9000',
        changeOrigin: true,
      },
    },
  },
});
