---
status: done
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-004-r3-six-packages
version: 0.1.0
---

# E-001 · 六包构建/发布/消费事实（2026-08-29）

## S1/S2 · 聚合入口 + 构建 + d.ts 自动化（TS5056 根治）

- 四个聚合入口：`src/{lib,theme}/index.ts`、`src/components/ui/index.ts`（+data-table）、`src/app/index.ts`（+host 面）。
- `scripts/build-lib-packages.mjs`：Vite JS API 循环构建（createRequire 从 apps/web 解析 vite；shell 引用误删 user-menu 源（交互测试内嵌）已修）。
- **d.ts 自动化修复**：TS5056 根因 = `render.ts`/`render.tsx`、`form-controls.ts`/`form-controls.tsx` 同名对 → **改名 `render.types.ts` / `form-controls.types.ts`**（22 文件引用更新）→ 五包 tsc declaration 全 0（补 `allowImportingTsExtensions`；四新 tsconfig 补 outDir 并清理 46 个 src 污染 d.ts）→ **F-006（workspace-022 挂账）核销**。
- 产物：`dist-lib/@schema-ui/{lib,theme,ui,shell}`（index.js + 全量 d.ts）+ renderer 0.2.0（d.ts 管线版）。

## S3 · 发布与消费

- 发布（存在性跳过防重发；pack 目录清理修复）：`@magicvr/schema-ui-{lib,theme,ui,shell}@0.1.0` + `@magicvr/schema-ui-renderer@0.2.0` + protocol 0.2.0（已有）。
- golden-field：六包 registry 拉取 → 旧三探针 PASS + **probe-six PASS**（cn/formatDisplayTime/resolveTheme/DataTable/App 核心导出）。renderer 0.1.0→0.2.0 = registry 升级又一实证。

## 知识/观察

- `formatDisplayTime` 仅接受 ISO 字符串；`resolveTheme` 输入 = `ThemeInput{stored,prefersDark,systemDefault}`（探针断言按真实语义核对）。
- GH Packages 对不存在版本 view = 404（本地语义），发布脚本以 try/catch 处理。