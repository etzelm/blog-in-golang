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
    // vitest 4 removed `poolOptions` — keep the previous "one fork process for
    // the whole run" behavior with `maxWorkers: 1` at the top of `test`. Don't
    // also set `isolate: false` (the migration guide's literal swap) — that
    // disables per-file module isolation, which breaks our cross-file mock
    // setup (App.test.jsx's `vi.mock('@react-oauth/google', ...)` would leak
    // into NavBar.test.jsx and others). https://vitest.dev/guide/migration#pool-rework
    maxWorkers: 1
  },
  resolve: {
    extensions: ['.js', '.jsx', '.ts', '.tsx'],
  },
});