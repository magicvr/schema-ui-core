---
id: E-004-e2e-typecheck-gap-closure
doc: execution-entry
status: recorded
goal_id: GOAL-008-typecheck-evidence-convention
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# E-004 · 关闭同类的第二层缺口：e2e 项目类型检查

2026-09-18，在 C4 自审前的复核中发现 F-005 的**同类第二层缺口**并修复：

## 发现

`apps/web/e2e/tsconfig.json` 存在且 `include: ["./**/*.ts"]`，但**不在**根 `tsconfig.json` 的 `references` 图内（图内只有 `tsconfig.app.json` 与 `tsconfig.node.json`）。因此单独的 `tsc -b` 不检查 e2e 规格文件——与 F-005 完全同型，只是下沉一层。

证明（注入 `const __e2eProbe: number = "str";` 于 `e2e/sign-in.ts`）：

| 命令 | 结果 |
|------|------|
| `npm run typecheck`（当时为 `tsc -b`） | 无输出，exit **0**（未检查 e2e） |
| `npx tsc -p e2e/tsconfig.json --noEmit` | `error TS2322` + `TS6133`，exit **2** |

即：若不在 C4 前发现，本目标会以「已修复类型检查证据」结项，而 e2e 仍未被检查。

## 修复

- `apps/web/package.json`：`"typecheck": "tsc -b && tsc -p e2e/tsconfig.json"`，显式覆盖两个项目。
- `apps/web/src/typecheck-convention.guard.test.ts`：新增断言「`typecheck` 必须包含 `-p e2e/tsconfig.json`」，并断言 e2e 项目当前**不在** references 图内（若将来并入图内，该断言会提示简化脚本而非重复检查）。
- `apps/web/README.md`：类型检查口径小节补充第二层陷阱与各命令的覆盖范围对照。
- `.github/workflows/r6-basic-matrix.yml`：web job 的 `npm run typecheck` 步骤自动继承新覆盖范围。

## 验证事实

| 场景 | 结果 |
|------|------|
| 干净工作树 `npm run typecheck` | `tsc -b && tsc -p e2e/tsconfig.json`，exit 0 |
| 注入 e2e 类型错误 | `e2e/sign-in.ts(64,7): error TS2322` + `TS6133`，非零退出 |
| e2e 项目独立类型检查 | `tsc -p e2e/tsconfig.json` 干净，exit 0 |
| 变异测试（从 `typecheck` 移除 e2e 项目） | 守卫第 3 条断言失败 ✅ |
| 全量 Web Vitest | **112 文件 / 1426 测试**通过 |

注入后均已还原（校验内容一致）。

该发现与修复记录为 `A-001 F-001`（med required，已 `fixed`）。C3/C4 完成。
