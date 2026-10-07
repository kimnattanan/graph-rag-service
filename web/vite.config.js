import { defineConfig } from 'vite';

// Browser calls stay same-origin. Each service mounts its API at /api.
export default defineConfig({
  server: {
    port: 5173,
    proxy: {
      '/svc/user': {
        target: 'http://127.0.0.1:3001',
        rewrite: (path) => path.replace(/^\/svc\/user/, ''),
      },
      '/svc/knowledge': {
        target: 'http://127.0.0.1:3000',
        rewrite: (path) => path.replace(/^\/svc\/knowledge/, ''),
      },
      '/svc/conversation': {
        target: 'http://127.0.0.1:3002',
        rewrite: (path) => path.replace(/^\/svc\/conversation/, ''),
        timeout: 600_000,
        proxyTimeout: 600_000,
      },
    },
  },
});
