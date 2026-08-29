#!/usr/bin/env node
/**
 * build-lib-packages.mjs — 六包产物构建（R3 · VP-023 判据 #3）。
 * Vite JS API 循环构建：lib / theme / ui / shell（protocol/renderer 已有产物不动）。
 * 输出：apps/web/dist-lib/@schema-ui/<pkg>/{index.js, index.d.ts}
 * 用法：node scripts/build-lib-packages.mjs；d.ts = 紧跟的 tsc declaration（见 tsconfig.*.json 兄弟配置）。
 */
import { createRequire } from "node:module";
import path from "node:path";
import { fileURLToPath } from "node:url";

const here = path.dirname(fileURLToPath(import.meta.url));
const webRoot = path.resolve(here, "../apps/web");
// vite 位于 apps/web/node_modules（根无安装）——从 web 包解析。
const require = createRequire(path.join(webRoot, "package.json"));
const { build } = require("vite");
const src = (p) => path.join(webRoot, "src", p);
const alias = {
  "@": path.join(webRoot, "src"),
  "@schemas": path.resolve(webRoot, "../../docs/schemas"),
};

const packages = [
  {
    name: "@schema-ui/lib",
    entry: src("lib/index.ts"),
    external: ["react", "react-dom", "react/jsx-runtime"],
  },
  {
    name: "@schema-ui/theme",
    entry: src("theme/index.ts"),
    external: [],
  },
  {
    name: "@schema-ui/ui",
    entry: src("components/ui/index.ts"),
    external: ["react", "react-dom", "react/jsx-runtime"],
  },
  {
    name: "@schema-ui/shell",
    entry: src("app/index.ts"),
    external: ["react", "react-dom", "react/jsx-runtime"],
  },
];

for (const pkg of packages) {
  const outDir = path.join(webRoot, "dist-lib", pkg.name);
  await build({
    configFile: false,
    resolve: { alias },
    build: {
      lib: { entry: pkg.entry, formats: ["es"], fileName: () => "index.js", name: pkg.name },
      outDir,
      emptyOutDir: false,
      sourcemap: false,
      minify: false,
      rollupOptions: { external: pkg.external },
    },
  });
  console.log(`built ${pkg.name} → ${outDir}`);
}
console.log("六包构建完成（protocol/renderer 沿用既有产物）");