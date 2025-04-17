import path from "path";

import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// https://vite.dev/config/
export default defineConfig({
  plugins: [svelte()],
  resolve: {
    alias: {
      "$style": path.resolve("./src/style"),
      "$lib": path.resolve("./src/lib"),
      "$components": path.resolve("./src/components"),
      "$modules": path.resolve("./src/modules"),
    }
  }
})
