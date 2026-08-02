import path from "node:path";
import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
      "@schemas": path.resolve(__dirname, "../../docs/schemas"),
    },
  },
  server: {
    host: "127.0.0.1",
    // Default stays 5173 (CI/Linux). WEB_PORT lets Windows hosts outside the
    // Hyper-V excluded range (e.g. 9999) run Playwright without forking config.
    port: Number(process.env.WEB_PORT || 5173),
    strictPort: true,
    proxy: {
      "/api": {
        target: "http://127.0.0.1:8080",
        changeOrigin: true,
      },
    },
  },
});
