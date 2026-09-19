---
id: E-002-r4-regression
doc: execution-entry
parent: GOAL-005-r4-regression-and-closeout
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# E-002 · R4 回归结果

## API

- `apps/api`: `go test ./... -count=1` → **PASS**，所有包 exit 0。

## Web unit/integration

- direct `node_modules/.bin/vitest run --reporter=dot` → **PASS**，123 test files / 1489 tests。
- direct forced TypeScript: `node_modules/typescript/bin/tsc -b --force` + `tsc -p e2e/tsconfig.json` → **PASS**。
- package-manager wrappers `pnpm run typecheck` / `pnpm run test` 未进入测试阶段：pnpm 自检尝试重写 locked `apps/web/node_modules/.bin`，Windows EPERM；同一已安装依赖的 direct binaries 已完成验证。

## Browser

- default/mvp managed Playwright full run → **17 passed / 5 skipped / 1 failed**。
- 唯一失败：`e2e/list-visual-surface.spec.ts:82` 在共享 fresh-seed suite 中，`sign-in.ts` fallback 等待 Sign in 按钮启用超时；不是页面断言。该机制性顺序契约已有 roadmap bounded residual（`roadmap.md` §未决项统一登记，继承 workspace-038/010）；隔离复验：同一 test 单独运行 → **1 passed**。
- admin Profile core slice (`console-health`, `host-failure`, `shell`, `localization`) → **7 passed / 1 skipped**（`APP_PROFILE=admin`）。

## Browser command

`$env:APP_PROFILE='admin'; node node_modules/@playwright/test/cli.js test e2e/shell.spec.ts e2e/console-health.spec.ts e2e/host-failure.spec.ts e2e/localization.spec.ts --project=chromium` → exit 0。

## R4 解释

全量浏览器失败属于既有 harness fresh-seed 顺序 residual，不是 VP-039 产品失败；R2/R3 相关 unit/integration、Go full、isolated browser regression 均通过。R4 C3 将保留该 bounded residual，不扩大为本 VP required finding。
