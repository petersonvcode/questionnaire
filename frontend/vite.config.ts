import { defineConfig } from 'vite';

export default defineConfig({
  plugins: [],
  root: 'src',
  server: {
    port: 3000,
  },
  build: {
    outDir: '../dist',
    target: 'es2015',
    minify: 'terser',
  }
});
