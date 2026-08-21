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
    // Default 25173 (>25000) avoids Windows Hyper-V excluded ranges; WEB_PORT
    // still overrides when another port is needed.
    port: Number(process.env.WEB_PORT || 25173),
    strictPort: true,
    proxy: {
      "/.well-known/schema-ui/app-manifest.json": {
        target: "http://127.0.0.1:25080",
        changeOrigin: true,
      },
      "/.well-known/schema-ui/host-bootstrap.json": {
        target: "http://127.0.0.1:25080",
        changeOrigin: true,
      },
      "/api": {
        target: "http://127.0.0.1:25080",
        changeOrigin: true,
      },
    },
  },
  preview: {
    host: "127.0.0.1",
    port: Number(process.env.WEB_PORT || 25173),
    strictPort: true,
    // Production preview proxies the same same-origin protocol endpoints so
    // the built SPA can be verified headlessly against the real API.
    proxy: {
      "/.well-known/schema-ui/app-manifest.json": {
        target: "http://127.0.0.1:25080",
        changeOrigin: true,
      },
      "/.well-known/schema-ui/host-bootstrap.json": {
        target: "http://127.0.0.1:25080",
        changeOrigin: true,
      },
      "/api": {
        target: "http://127.0.0.1:25080",
        changeOrigin: true,
      },
    },
  },
});
