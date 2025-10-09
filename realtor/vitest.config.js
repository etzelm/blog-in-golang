import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: './src/setupTests.jsx',
    css: false,
    cors: true,
    headers: {
        'Content-Security-Policy': "default-src 'self' 'unsafe-inline' 'unsafe-eval' data: blob: https://files.mitchelletzel.com https://apis.google.com https://accounts.google.com https://www.gstatic.com https://ssl.gstatic.com; img-src 'self' data: https://files.mitchelletzel.com; connect-src 'self' https://files.mitchelletzel.com;",
    },
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
    },
    pool: 'forks',
    poolOptions: {
      forks: {
        singleFork: true
      }
    }
  },
  resolve: {
    extensions: ['.js', '.jsx', '.ts', '.tsx'],
  },
});