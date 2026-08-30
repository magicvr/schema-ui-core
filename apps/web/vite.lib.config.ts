import path from "node:path";
import { defineConfig } from "vite";

/**
 * Web 包分发试点库构建（B 路径 · GOAL-004 S2）。
 * 产出 @schema-ui/protocol（ESM + index.d.ts），依赖照常 bundle（自包含包）。
 * 与主应用构建（vite.config.ts）正交；不引 react/tailwind 插件（protocol 纯 TS）。
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
      entry: path.resolve(__dirname, "src/protocol/index.ts"),
      formats: ["es"],
      fileName: () => "index.js",
      name: "SchemaUIProtocol",
    },
    outDir: "dist-lib/@schema-ui/protocol",
    emptyOutDir: false,
    sourcemap: false,
    minify: false,
  },
});