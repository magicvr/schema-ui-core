import path from "node:path";
import { defineConfig } from "vite";

/**
 * @schema-ui/renderer 库构建（B 路径 · GOAL-004 S3）。
 * 粗粒度单包：内部 bundle components/i18n/lib/protocol；
 * React 体系外部化（peer 依赖，防双实例）。
 */
export default defineConfig({
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
      "@schemas": path.resolve(__dirname, "../../docs/schemas"),
    },
  },
  build: {
    lib: {
      entry: path.resolve(__dirname, "src/renderer/index.ts"),
      formats: ["es"],
      fileName: () => "index.js",
      name: "SchemaUIRenderer",
    },
    outDir: "dist-lib/@schema-ui/renderer",
    emptyOutDir: false,
    sourcemap: false,
    minify: false,
    rollupOptions: {
      external: [
        "react",
        "react-dom",
        "react/jsx-runtime",
        "react-dom/client",
        "react-dom/server",
      ],
    },
  },
});