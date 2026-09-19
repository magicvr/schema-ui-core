---
doc_type: evidence-matrix
id: r4-regression-matrix
parent: GOAL-005-r4-regression-and-closeout
status: verified
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# R4 回归证据矩阵

## 1 · 已完成自动化验证

| 面 | 命令/证据 | 结果 |
|----|-----------|------|
| API 全量 | `apps/api`: `go test ./... -count=1` | **PASS**（所有包 exit 0；commit 当前回归基线） |
| Web Vitest 全量 | `apps/web`: direct `node_modules/.bin/vitest run --reporter=dot` | **PASS** · 123 files / 1489 tests |
| Web TypeScript | direct `node_modules/typescript/bin/tsc -b --force` + `tsc -p e2e/tsconfig.json` | **PASS** |
| R2/R3 targeted | RuntimeBanner / VersionChip / AuthGate / AuthContext / catalog / App integration | **PASS**（此前 targeted suites） |

`pnpm run typecheck` / `pnpm run test` 的 package-manager 自检尝试因 Windows 现有 Node 进程锁定 `apps/web/node_modules/.bin/{vite,playwright}` 而在 install 自检阶段退出；未作为产品失败，改用同一已安装依赖的 direct binaries 验证并通过。

## 2 · 回归维度

| 维度 | 覆盖 | 证据 |
|------|------|------|
| runtime mode | normal / maintenance / degraded / read-only | `bootstrap_test.go`；`runtime-banner.test.tsx`；operational gate 全量 Go tests |
| permission | 有/无 `monitoring.read` | `version-chip.test.tsx`；`hasMonitoringRead`；Profile/Manifest tests |
| profile | mvp/admin/demo/custom 既有 Profile/Manifest 约束 | `apps/api/kernel/profile*` tests；Playwright profile harness；现有 profile matrix tests |
| locale | zh-CN / en-US | `catalog.test.ts`、`ui-bilingual.test.tsx`、全量 Vitest |
| theme | light/dark token + shell tests | `startup-config.test.tsx`、App/renderer suites；R4 Playwright visual surface |
| Host boundary | consumer maintenance terminal fixture保留；producer maintenance→degraded | `host/bootstrap.ts`、upstream host fixtures、`handler/bootstrap_test.go` |
| write gate | existing SERVICE_* errors + recovery allowlist | `operational_test.go`、`r5_operational_gate_test.go`、Go full suite |

## 3 · Playwright

| 命令 | 状态 | 说明 |
|------|------|------|
| `apps/web`: `node node_modules/@playwright/test/cli.js test --project=chromium` | **bounded residual** | Full mvp run: 17 passed / 5 skipped / 1 failed at sign-in helper fresh-seed ordering; isolated failed spec: 1 passed; existing roadmap residual, not VP-039 product failure |
| `APP_PROFILE=admin` core slice | **PASS** | 7 passed / 1 skipped; shell, console-health, host-failure, localization |

## 4 · 关门边界

- 不改 `kernel/profile.go` 默认集。
- 不改 `apps/web/src/protocol/upstream/**`。
- 不改 VP-012 operational gate 语义。
- 不实现 VP-040 / Redis / MQ / 多实例 / 搜索引擎。
- Root/VP 关门仍等待用户书面确认。
