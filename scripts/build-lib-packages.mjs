#!/usr/bin/env node
/**
 * build-lib-packages.mjs — 六包产物构建（R3 · VP-023 判据 #3；R5 · VP-024 判据 #5 扩展）。
 *
 * R5（VP-024）：renderer 依赖图 external 化——构建时把内部面导入重写为包子路径并 external：
 *   @/i18n/*  → @schema-ui/lib/i18n/*
 *   @/protocol/* → @schema-ui/protocol/*
 *   @/lib/*   → @schema-ui/lib/lib/*
 *   @/components/ui/* → @schema-ui/ui/components/ui/*
 * 产物 index.js 的 import 语句指向 @schema-ui/*（消费端解析），renderer 不再自包含这些面。
 *
 * 输出：apps/web/dist-lib/@schema-ui/<pkg>/{index.js, <pkg>/...}
 * 用法：node scripts/build-lib-packages.mjs（d.ts = 紧随的 rewrite-lib-aliases + tsc declaration，见 R5 流程）。
 */
import { createRequire } from "node:module";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { rmSync, readFileSync, writeFileSync, existsSync } from "node:fs";

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

// R5：renderer 内部面 → 包子路径重写表（D-001-r5-granularity §1）。
const rendererFaceRewrite = {
  "@/i18n/": "@schema-ui/lib/i18n/",
  "@/protocol/": "@schema-ui/protocol/",
  "@/lib/": "@schema-ui/lib/lib/",
  "@/components/ui/": "@schema-ui/ui/components/ui/",
  "@/theme/": "@schema-ui/theme/",
};

// 对 renderer 依赖图内源码做 import 字符串级重写；被重写后的 @schema-ui/* id
// 命中 external 函数 → 不进 bundle，产物保留包子路径 import（消费端解析）。
function rewriteRenderFaces() {
  return {
    name: "rewrite-internal-faces",
    enforce: "pre",
    transform(code, id) {
      const faces = ["/src/i18n", "/src/protocol", "/src/lib", "/src/components/", "/src/renderer"];
      if (!faces.some((f) => id.includes(f))) return null;
      let out = code;
      for (const [k, v] of Object.entries(rendererFaceRewrite)) {
        if (!out.includes(k)) continue;
        out = out.split(k).join(v);
      }
      return out;
    },
  };
}

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
  {
    name: "@schema-ui/renderer",
    entry: src("renderer/index.ts"),
    external: (id) =>
      id.startsWith("@schema-ui/") ||
      id === "react" || id === "react-dom" || id === "react/jsx-runtime",
    plugins: [rewriteRenderFaces()],
    clean: true, // R5：清残留目录后重建
  },
];

for (const pkg of packages) {
  const outDir = path.join(webRoot, "dist-lib", pkg.name);
  if (pkg.clean) {
    rmSync(outDir, { recursive: true, force: true });
    console.log(`cleaned ${pkg.name} dist`);
  }
  await build({
    configFile: false,
    resolve: { alias },
    plugins: pkg.plugins ?? [],
    build: {
      lib: { entry: pkg.entry, formats: ["es"], fileName: () => "index.js", name: pkg.name },
      outDir,
      emptyOutDir: false,
      sourcemap: false,
      minify: false,
      rollupOptions: { external: pkg.external },
    },
  });
  console.log(`built ${pkg.name} -> ${outDir}`);
}

// R5 终版（VP-024 判据 #5 闭环）：protocol/lib/ui 改 tsc 全产物（js+d.ts 树）——
// renderer external 化后的子路径 import 需要目标包具备子路径 JS（bundle 型包不满足），
// tsc 逐模块产物天然满足：conformance/component-format、lib/datetime、components/ui/card 等。
import { execFileSync } from "node:child_process";
const tscEntries = {
  "@schema-ui/protocol": { include: "src/protocol/index.ts", entryJs: "protocol/index.js", entryTypes: "protocol/index.d.ts" },
  "@schema-ui/lib": { include: "src/lib/index.ts", entryJs: "lib/index.js", entryTypes: "lib/index.d.ts" },
  "@schema-ui/ui": { include: "src/components/ui/index.ts", entryJs: "components/ui/index.js", entryTypes: "components/ui/index.d.ts" },
};
for (const [name, cfg] of Object.entries(tscEntries)) {
  const outDir = path.join(webRoot, "dist-lib", name);
  rmSync(outDir, { recursive: true, force: true });
  execFileSync(process.platform === "win32" ? "npx.cmd" : "npx", ["tsc", "-p", `tsconfig.${name.split("/")[1]}.json`, "--emitDeclarationOnly", "false", "--allowImportingTsExtensions", "false"], {
    cwd: webRoot,
    stdio: "inherit",
    shell: process.platform === "win32",
  });
  // 修正包入口（tsc 产物入口 = src 内入口的镜像路径）；package.json 可能随 rmSync 被删 → 模板兜底
  const pkgPath = path.join(outDir, "package.json");
  let pkg = existsSync(pkgPath) ? JSON.parse(readFileSync(pkgPath, "utf8")) : {
    name, version: "0.0.0", type: "module",
    main: cfg.entryJs, types: cfg.entryTypes,
    exports: { ".": { types: cfg.entryTypes, import: cfg.entryJs }, "./*": "./*" },
    files: ["index.js"],
    license: "UNLICENSED",
  };
  pkg.main = cfg.entryJs;
  pkg.types = cfg.entryTypes;
  pkg.exports = { ".": { types: cfg.entryTypes, import: cfg.entryJs }, "./*": "./*" };
  writeFileSync(pkgPath, JSON.stringify(pkg, null, 2));
  console.log(`tsc-built ${name} -> ${outDir}（js+d.ts 树 · 子路径可解析）`);
}
console.log("六包构建完成：protocol/lib/ui = tsc 全产物 · renderer = vite external + tsc 声明 · theme/shell = vite bundle");