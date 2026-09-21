---
id: E-003-guard-implementation-and-mutation
doc: execution-entry
status: recorded
goal_id: GOAL-008-typecheck-evidence-convention
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# E-003 · 实施防复发守卫（结构断言 + CI 门禁）并通过变异验证

2026-09-18，按 `D-002` 完成 C3：

## 实施

1. **结构守卫测试**：新增 `apps/web/src/typecheck-convention.guard.test.ts`（5 个断言），沿用本仓 `.guard.test.ts` 约定。断言内容：根 `tsconfig` 保持 solution-style；`typecheck` 脚本存在且为检查型；`build` 仍含检查型 `tsc`；所有可执行面无非检查型 `tsc` 调用；README 保留该约定。
2. **CI 显式门禁**：`.github/workflows/r6-basic-matrix.yml` 的 `web` job 在 `npm test` 之后新增 `- run: npm run typecheck`，并附注释说明「守卫只断言约定，本步才真正检查」以及不得使用裸 `tsc --noEmit` 的原因。

守卫的检测逻辑区分三类情形，避免误报：
- **命令位置**：`tsc` 必须处于命令位置（行首、shell 分隔符后、`npx`/`npm exec` 后，或 argv 元素起始的 `"tsc"`/`['tsc'`）；正文/日志中的提及（如 `… = tsc 全产物`）不视为调用。
- **命令词边界**：`tsc-built` 这类标识符不算 `tsc` 命令。
- **检查型标志**：`-b`/`--build`/`-p`/`--project` 须为独立 token（允许尾随引号/逗号/分号等标点），且可出现在 argv 数组的后续元素中（`["tsc", "-p", …]`）。

## 变异验证（证明守卫非空转）

仅「测试通过」不足以证明守卫有效，故对 4 种真实失效模式做变异测试：

| # | 变异 | 期望 | 实测 |
|---|------|------|------|
| 1 | `"typecheck": "tsc -b"` → `"tsc --noEmit"` | 失败 | ✅ 2 项断言失败（含 `typecheck must use -b or -p`） |
| 2 | `"build": "tsc -b && vite build"` → `"vite build"` | 失败 | ✅ `keeps the build script type-checked` 失败 |
| 3 | CI 工作流插入 `npx tsc --noEmit` | 失败 | ✅ 非检查型调用断言失败（复现历史失效模式） |
| 4 | `tsconfig.json` 加 `"include": ["src"]` | 失败 | ✅ `now selects sources; revisit the type-check convention` |

全部变异后均已还原文件（逐次校验内容一致），工作树最终只含预期新增/修改。

## 验证事实

- `npm exec -- vitest run src/typecheck-convention.guard.test.ts`：5/5 通过。
- 全量 Web Vitest：**112 文件 / 1425 测试**通过（+1 文件、+5 测试）。
- `npm run typecheck`（`tsc -b`）：exit 0。
- 变异测试：4/4 捕获。

C3 完成。C4 自审待执行。
